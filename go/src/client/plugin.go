package client

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"fs"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/textproto"
	"net/url"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

// Client export
type Client struct {
	Server            string
	User              string
	Password          string
	httpClient        *http.Client
	WsConn            *websocket.Conn
	interrupt         chan int
	OnMessageCallback func([]byte) string
	ClientInfo        string
}

type message struct {
	Type   string          `json:"type"`
	Status int             `json:"status,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

type genericMessage struct {
	Type   string      `json:"type"`
	Status int         `json:"status,omitempty"`
	Data   interface{} `json:"data"`
}

// NewClient export
func NewClient(url, user, password string) *Client {
	c := Client{}
	c.Server = url
	c.User = user
	c.Password = password
	cookieJar, _ := cookiejar.New(nil)
	c.httpClient = &http.Client{Jar: cookieJar}
	return &c
}

// sends message with status code 200 ("ok")
func (c *Client) sendResponseMessage(msgType string, data interface{}) error {
	return c.WsConn.WriteJSON(genericMessage{Type: msgType, Status: 200, Data: data})
}

// sends error message
func (c *Client) sendErrorMessage(msgType string, data string) error {
	return c.WsConn.WriteJSON(genericMessage{Type: msgType, Status: 500, Data: data})
}

func (c *Client) readTextData(data string) (string, error) {
	var msg message
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		return "", err
	}
	var value string
	if err := json.Unmarshal(msg.Data, &value); err != nil {
		return "", err
	}
	return value, nil
}

type messageHandler interface {
	Accept(msgType string) bool
	Process(msg message) error
}

type statusHandler struct {
	client *Client
}

func (h *statusHandler) Accept(msgType string) bool {
	return msgType == "PluginStatus"
}

func (h *statusHandler) Process(msg message) error {
	info := map[string]string{"status": "connected", "client": h.client.ClientInfo}
	return h.client.sendResponseMessage("PluginStatus", info)
}

type filesHandler struct {
	client *Client
}

func (h *filesHandler) Accept(msgType string) bool {
	return msgType == "ProjectFiles"
}

func (h *filesHandler) Process(msg message) error {
	type filesMsg struct {
		Directory string    `json:"directory"`
		Files     []fs.File `json:"files"`
	}

	// TODO: better error handling
	r := h.client.OnMessageCallback([]byte(`{"type":"ProjectDirectory"}`))
	directory, _ := h.client.readTextData(r)
	if directory == "" {
		return errors.New("Project is not open")
	}
	files, err := fs.ListDir(directory, true)

	if err != nil {
		return err
	}
	for i, f := range *files {
		(*files)[i].Path = filepath.ToSlash(f.Path)
	}
	resp := filesMsg{Directory: directory, Files: *files}
	return h.client.sendResponseMessage(msg.Type, resp)
}

type uploadHandler struct {
	client *Client
	cancel context.CancelFunc
}

func (h *uploadHandler) Accept(msgType string) bool {
	return msgType == "UploadFiles" || msgType == "AbortUpload"
}

func (h *uploadHandler) Process(msg message) error {
	if msg.Type == "AbortUpload" {
		if h.cancel != nil {
			h.cancel()
			h.cancel = nil
		}
		return nil
	}
	type Params struct {
		Project string    `json:"project"`
		Files   []fs.File `json:"files"`
	}
	var params Params
	if err := json.Unmarshal(msg.Data, &params); err != nil {
		return err
	}

	go func() {
		readBody, writeBody := io.Pipe()
		defer readBody.Close()

		writer := multipart.NewWriter(writeBody)
		errChan := make(chan error, 1)

		go func() {
			compressRegex := regexp.MustCompile("(?i).*\\.(qgs|svg|json|sqlite|gpkg|geojson)$")
			defer writeBody.Close()
			writer.WriteField("changes", string(msg.Data))

			// TODO: better error handling
			r := h.client.OnMessageCallback([]byte(`{"type":"ProjectDirectory"}`))
			directory, _ := h.client.readTextData(r)

			for _, f := range params.Files {
				// ext := filepath.Ext(f.Path)
				fileOsPath := filepath.FromSlash(f.Path)
				useCompression := compressRegex.Match([]byte(f.Path))
				if useCompression {
					mh := make(textproto.MIMEHeader)
					mh.Set("Content-Type", "application/octet-stream")
					mh.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s.gz"`, f.Path, f.Path))
					part, _ := writer.CreatePart(mh)
					gzpart := gzip.NewWriter(part)
					err := fs.CopyFile(gzpart, filepath.Join(directory, fileOsPath))
					gzpart.Close()
					if err != nil {
						errChan <- err
						return
					}
				} else {
					part, err := writer.CreateFormFile(f.Path, f.Path)
					if err != nil {
						errChan <- err
						return
					}
					if err = fs.CopyFile(part, filepath.Join(directory, fileOsPath)); err != nil {
						errChan <- err
						return
					}
				}
			}
			errChan <- writer.Close()
		}()

		url := fmt.Sprintf("%s/api/project/upload/%s", h.client.Server, params.Project)
		req, _ := http.NewRequest("POST", url, readBody)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		ctx, cancel := context.WithCancel(context.Background())
		req = req.WithContext(ctx)
		h.cancel = cancel

		resp, err := h.client.httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
			h.client.WsConn.WriteJSON(message{Type: "UploadError"})
			return
		}
		defer resp.Body.Close()
		h.cancel = nil

		fmt.Println("Status", resp.StatusCode)

		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		if resp.StatusCode >= 400 {
			fmt.Println("Send UploadError")
			if err = h.client.sendErrorMessage("UploadError", string(respData)); err != nil {
				fmt.Printf("Failed to send error message: %s\n", err)
			}
		}
		fmt.Println(string(respData))
		err = <-errChan
		if err != nil {
			fmt.Println(err)
		}
	}()
	return nil
}

func (c *Client) login() error {
	form := url.Values{"username": {c.User}, "password": {c.Password}}
	url := fmt.Sprintf("%s/login/", c.Server)
	resp, err := c.httpClient.PostForm(url, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("Authentication failed")
	}
	return nil
}

func (c *Client) logout() error {
	url := fmt.Sprintf("%s/logout/", c.Server)
	_, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	return nil
}

// Start export
func (c *Client) Start() error {
	err := c.login()
	if err != nil {
		return err
	}
	defer c.logout()

	c.interrupt = make(chan int, 1)

	u, _ := url.Parse(c.Server)
	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	u.Path = fmt.Sprintf("/ws/plugin")
	log.Printf("connecting to %s", u.String())

	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 30 * time.Second,
		Jar:              c.httpClient.Jar,
	}
	header := make(http.Header, 1)
	header.Set("User-Agent", c.ClientInfo)
	wsConn, _, err := dialer.Dial(u.String(), header)
	if err != nil {
		return err
	}
	c.WsConn = wsConn
	defer wsConn.Close()

	done := make(chan struct{})

	var messageHandlers []messageHandler
	messageHandlers = append(messageHandlers, &statusHandler{c})
	messageHandlers = append(messageHandlers, &filesHandler{c})
	messageHandlers = append(messageHandlers, &uploadHandler{client: c})

	go func() {
		defer close(done)

	MESSAGE:
		for {
			_, rawMessage, err := wsConn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var msg message
			if err = json.Unmarshal(rawMessage, &msg); err != nil {
				log.Printf("Invalid message: %s\n", rawMessage)
				continue
			}
			// log.Println("Msg type: ", msg.Type)
			// log.Printf("Received: %s\n", message)
			for _, h := range messageHandlers {
				if h.Accept(msg.Type) {
					if err := h.Process(msg); err != nil {
						log.Println(err)
						c.sendErrorMessage(msg.Type, err.Error())
					}
					continue MESSAGE
				}
			}
			// possible issue if executed in different thread?
			resp := c.OnMessageCallback(rawMessage)
			if resp != "" {
				c.WsConn.WriteMessage(websocket.TextMessage, []byte(resp))
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-c.interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return nil
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

// Stop export
func (c *Client) Stop() {
	c.interrupt <- 1
}
