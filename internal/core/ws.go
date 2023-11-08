package core

import (
	"bufio"
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"wsfuzz/internal/model"

	"github.com/gorilla/websocket"
)

func SendData(reqd *model.RequestData, id int) {
	// 创建一个自定义的 Dialer
	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// 进行 WebSocket 连接
	u := url.URL{Scheme: "wss", Host: reqd.Host, Path: reqd.Path, RawQuery: strings.Replace(reqd.Query, "{CG}", strconv.Itoa(id), 1)}
	// log.Printf("Request  %v", u.String())

	c, _, err := dialer.Dial(u.String(), reqd.Headers)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	if err := c.WriteMessage(websocket.BinaryMessage, reqd.Body); err != nil {
		log.Fatal("write:", err)
	}

	// 接收WebSocket服务器返回的数据，并输出到控制台
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Fatal("read:", err)
	}

	log.Printf("Request: %v  Received %s\n", u.String(), message)
}

func ParseFile(filepath string) (reqd *model.RequestData) {
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

	var body []byte
	if len(parts) > 2 {
		body = []byte(parts[2])
	}

	// 使用WebSocket库将HTTP请求内容发送到指定的WebSocket地址
	if host == "" {
		log.Print("Host 获取失败")
		return
	}

	reqd = &model.RequestData{
		Host:    host,
		Path:    path,
		Query:   query,
		Body:    body,
		Headers: headers,
	}
	return
}
