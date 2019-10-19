package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
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
	"strings"
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

// SendJSONMessage export
func (c *Client) SendJSONMessage(msgType string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("%s:%s", msgType, jsonData)
	c.WsConn.WriteMessage(websocket.TextMessage, []byte(msg))
	return nil
}

type messageHandler interface {
	Accept(msgType string) bool
	Process(msg []byte) error
}

type pingHandler struct {
	client *Client
}

func (h *pingHandler) Accept(msgType string) bool {
	return msgType == "PingPlugin"
}

func (h *pingHandler) Process(msg []byte) error {
	return h.client.WsConn.WriteMessage(websocket.TextMessage, []byte("PongPlugin"))
}

type filesHandler struct {
	client *Client
}

func (h *filesHandler) Accept(msgType string) bool {
	return msgType == "Files"
}

func (h *filesHandler) Process(msg []byte) error {
	type filesMsg struct {
		Directory string    `json:"directory"`
		Files     []fs.File `json:"files"`
	}

	directory := h.client.OnMessageCallback([]byte("ProjectDirectory"))
	files, err := fs.ListDir(directory, true)

	if err != nil {
		return err
	}
	for i, f := range *files {
		(*files)[i].Path = filepath.ToSlash(f.Path)
	}
	resp := filesMsg{Directory: directory, Files: *files}
	if err = h.client.SendJSONMessage("Files", resp); err != nil {
		return err
	}
	// if err = connection.WriteJSON(resp); err != nil {
	// 	fmt.Println(err)
	// }
	return nil
}

type uploadHandler struct {
	client *Client
	cancel context.CancelFunc
}

func (h *uploadHandler) Accept(msgType string) bool {
	return msgType == "SendFiles" || msgType == "AbortUpload"
}

func (h *uploadHandler) Process(msg []byte) error {
	if bytes.HasPrefix(msg, []byte("AbortUpload")) {
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
	jsonData := bytes.TrimPrefix(msg, []byte("SendFiles:"))
	// jsonData := msg[len("SendFiles:"):]
	if err := json.Unmarshal(jsonData, &params); err != nil {
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
			writer.WriteField("changes", string(jsonData))
			directory := h.client.OnMessageCallback([]byte("ProjectDirectory"))

			for _, f := range params.Files {
				// ext := filepath.Ext(f.Path)
				fileOsPath := filepath.FromSlash(f.Path)
				useCompression := compressRegex.Match([]byte(f.Path))
				if useCompression {
					mh := make(textproto.MIMEHeader)
					mh.Set("Content-Type", "application/octet-stream")
					mh.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"%s\"; filename=\"%s.gz\"", f.Path, f.Path))
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
			h.client.WsConn.WriteMessage(websocket.TextMessage, []byte("UploadError"))
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
			msg := fmt.Sprintf("UploadError:%s", respData)
			err = h.client.WsConn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
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
		return fmt.Errorf("Authentication failed")
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
	wsConn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	c.WsConn = wsConn
	defer wsConn.Close()

	done := make(chan struct{})

	var messageHandlers []messageHandler
	messageHandlers = append(messageHandlers, &pingHandler{c})
	messageHandlers = append(messageHandlers, &filesHandler{c})
	messageHandlers = append(messageHandlers, &uploadHandler{client: c})

	go func() {
		defer close(done)

	MESSAGE:
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			strMsg := string(message)
			log.Printf("recv: %s", message)
			msgType := strings.SplitN(strMsg, ":", 2)[0]
			for _, h := range messageHandlers {
				if h.Accept(msgType) {
					if err := h.Process(message); err != nil {
						fmt.Println(err)
					}
					continue MESSAGE
				}
			}
			// possible issue if executed in different thread?
			resp := c.OnMessageCallback(message)
			reply := fmt.Sprintf("%s:%s", msgType, resp)
			c.WsConn.WriteMessage(websocket.TextMessage, []byte(reply))
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
