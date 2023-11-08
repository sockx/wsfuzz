package core

import (
	"bufio"
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	Debug  bool
	WS     websocket.Dialer
	Url    *url.URL
	Header http.Header
	Conn   *websocket.Conn
}

func DefaultWebSocket() *WebSocket {
	dialer := *websocket.DefaultDialer
	return &WebSocket{
		WS: dialer,
	}
}

func (w *WebSocket) SetUnsafe() {
	w.WS.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func (w *WebSocket) SetHeaders(headers http.Header) {
	w.Header = headers
}

func (w *WebSocket) Connect() (err error) {
	w.Conn, _, err = w.WS.Dial(w.Url.String(), w.Header)
	if err != nil {
		return err
	}
	return
}

func (w *WebSocket) Close() error {
	return w.Conn.Close()
}

func (w *WebSocket) WriteMessage(messageType int, data []byte) error {
	if w.Debug {
		log.Printf("[tp] %v | [data] %v", messageType, string(data))
	}
	return w.Conn.WriteMessage(messageType, data)
}
func (w *WebSocket) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = w.Conn.ReadMessage()
	if w.Debug {
		log.Printf("[tp] %v | [data] %v", messageType, string(p))
	}
	return
}

func (w *WebSocket) ParseUri(uri string) (err error) {
	w.Url, err = url.Parse(uri)
	return
}

func (w *WebSocket) ParseFile(filepath string) {

	// 从文件中读取HTTP请求内容
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var httpRequest string
	for scanner.Scan() {
		httpRequest += scanner.Text() + "\n"
	}

	// 解析HTTP请求中的请求方式、path、header、body等内容
	parts := strings.Split(httpRequest, "\n\n")
	if len(parts) < 2 {
		log.Fatal("Invalid HTTP request format")
	}

	// 提取请求方式和path
	requestLine := parts[0]
	requestParts := strings.Split(requestLine, " ")
	if len(requestParts) != 3 {
		log.Fatal("Invalid request line format")
	}

	// 解析请求的 URL
	urlObj, err := url.Parse(requestParts[1])
	if err != nil {
		log.Fatal("Invalid URL: ", err)
	}
	// 分别获取 path 和 query 参数
	path := urlObj.Path
	query := urlObj.RawQuery

	// 提取header和body内容
	headers := http.Header{}
	var host string
	headerLines := strings.Split(parts[1], "\n")
	for _, headerLine := range headerLines {
		parts := strings.SplitN(headerLine, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "Host" {
			host = value
		} else if !(key == "Upgrade" || key == "Connection" || key == "Sec-WebSocket-Key" || key == "Sec-WebSocket-Extensions" || key == "Sec-WebSocket-Version") {
			headers.Add(key, value)
		}
	}

	// 使用WebSocket库将HTTP请求内容发送到指定的WebSocket地址
	if host == "" {
		log.Print("Host 获取失败")
		return
	}

	w.Url = &url.URL{
		Scheme:   "ws",
		Host:     host,
		Path:     path,
		RawQuery: query,
	}
	w.Header = headers
}
