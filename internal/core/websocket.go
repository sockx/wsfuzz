package core

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"

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

// func (w *WebSocket) SetWs(host, path, rawQuery string) {
// 	w.Url = &url.URL{
// 		Scheme:   "ws",
// 		Host:     host,
// 		Path:     path,
// 		RawQuery: rawQuery,
// 	}
// }

// func (w *WebSocket) SetWss(host, path, rawQuery string) {
// 	w.Url = &url.URL{
// 		Scheme:   "wss",
// 		Host:     host,
// 		Path:     path,
// 		RawQuery: rawQuery,
// 	}
// }

func (w *WebSocket) ParseUri(uri string) (err error) {
	w.Url, err = url.Parse(uri)
	return
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
