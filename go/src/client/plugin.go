package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"fs"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

// Client export
type Client struct {
	Server            string
	User              string
	WsConn            *websocket.Conn
	interrupt         chan os.Signal
	OnMessageCallback func([]byte) string
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
	Accept(msg []byte) bool
	Process(msg []byte) error
}

type pingHandler struct {
	client *Client
}

func (h pingHandler) Accept(msg []byte) bool {
	return bytes.Compare(msg, []byte("PingPlugin")) == 0
}

func (h pingHandler) Process(msg []byte) error {
	return h.client.WsConn.WriteMessage(websocket.TextMessage, []byte("PongPlugin"))
}

type filesHandler struct {
	client *Client
}

func (h filesHandler) Accept(msg []byte) bool {
	return bytes.Compare(msg, []byte("Files")) == 0
}

func (h filesHandler) Process(msg []byte) error {
	type filesMsg struct {
		Directory string    `json:"directory"`
		Files     []fs.File `json:"files"`
	}

	directory := h.client.OnMessageCallback([]byte("projectDirectory"))
	files, err := fs.ListDir(directory)
	if err != nil {
		return err
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
}

func (h uploadHandler) Accept(msg []byte) bool {
	return bytes.HasPrefix(msg, []byte("SendFiles:"))
}

func (h uploadHandler) Process(msg []byte) error {
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

	readBody, writeBody := io.Pipe()
	defer readBody.Close()

	writer := multipart.NewWriter(writeBody)
	errChan := make(chan error, 1)

	go func() {
		compressRegex := regexp.MustCompile("(?i).*\\.(qgs|svg|json|sqlite|gpkg|geojson)$")
		defer writeBody.Close()
		writer.WriteField("changes", string(jsonData))
		directory := h.client.OnMessageCallback([]byte("projectDirectory"))

		for _, f := range params.Files {
			// ext := filepath.Ext(f.Path)
			useCompression := compressRegex.Match([]byte(f.Path))
			if useCompression {
				mh := make(textproto.MIMEHeader)
				mh.Set("Content-Type", "application/octet-stream")
				mh.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"%s\"; filename=\"%s.gz\"", f.Path, f.Path))
				part, _ := writer.CreatePart(mh)
				gzpart := gzip.NewWriter(part)
				err := fs.CopyFile(gzpart, filepath.Join(directory, f.Path))
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
				if err = fs.CopyFile(part, filepath.Join(directory, f.Path)); err != nil {
					errChan <- err
					return
				}
			}
		}
		errChan <- writer.Close()
	}()

	client := &http.Client{}
	url := fmt.Sprintf("%s/api/project/upload/%s", h.client.Server, params.Project)
	req, err := http.NewRequest("POST", url, readBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// token, ok := os.LookupEnv("AUTH")
	// if ok {
	// 	req.Header.Set("Authorization", token)
	// }

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Status", resp.Status)
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(string(respData))
		return err
	}
	return <-errChan
}

// Start export
func (c *Client) Start() error {
	c.interrupt = make(chan os.Signal, 1)
	signal.Notify(c.interrupt, os.Interrupt)

	u, _ := url.Parse(c.Server)
	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	u.Path = fmt.Sprintf("/ws/plugin/%s", c.User)
	log.Printf("connecting to %s", u.String())

	wsConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	c.WsConn = wsConn
	defer wsConn.Close()

	done := make(chan struct{})

	var messageHandlers []messageHandler
	messageHandlers = append(messageHandlers, pingHandler{c})
	messageHandlers = append(messageHandlers, filesHandler{c})
	messageHandlers = append(messageHandlers, uploadHandler{c})

	go func() {
		defer close(done)

	MESSAGE:
		for {
			_, message, err := wsConn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			for _, h := range messageHandlers {
				if h.Accept(message) {
					if err := h.Process(message); err != nil {
						fmt.Println(err)
					}
					continue MESSAGE
				}
			}
			// possible issue if executed in different thread?
			c.OnMessageCallback(message)
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
	c.interrupt <- os.Interrupt
}
