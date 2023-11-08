package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"wsfuzz/internal/conf"
	"wsfuzz/internal/core"

	"github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "wsfuzz",
	Short: "ws爆破工具",
}

func Execute() {
	coloredcobra.Init(&coloredcobra.Config{
		RootCmd:         RootCmd,
		Headings:        coloredcobra.HiGreen + coloredcobra.Underline,
		Commands:        coloredcobra.Cyan + coloredcobra.Bold,
		Example:         coloredcobra.Italic,
		ExecName:        coloredcobra.Bold,
		Flags:           coloredcobra.Cyan + coloredcobra.Bold,
		NoExtraNewlines: true,
	})
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var localFile string

func init() {
	RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	RootCmd.PersistentFlags().BoolVarP(&conf.Options.Debug, "debug", "d", false, "enable debug")

	RootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVarP(&localFile, "file", "f", "req.txt", "HTTP请求保存所在文件")
}

var testCmd = &cobra.Command{
	Use:   "fuzz",
	Short: "Run fuzz",
	Run: func(cmd *cobra.Command, args []string) {
		ws := core.DefaultWebSocket()
		ws.ParseFile(localFile)
		ws.Url.Scheme = "wss"
		thisuri := *ws.Url

		for i := 470; i < 480; i++ {
			ws = core.DefaultWebSocket()
			nuri := thisuri
			ws.Url = &nuri
			ws.Url.RawQuery = strings.Replace(ws.Url.RawQuery, "{CG}", fmt.Sprintf("%v", i), 1)
			err := ws.Connect()
			if err != nil {
				log.Printf("connect error: %v", err)
				return
			}
			log.Print(ws.Url)
			ws.WriteMessage(2, nil)

			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Fatal("read:", err)
			}
			log.Printf("Request: %v  Received %s\n", ws.Url.String(), message)
		}

	},
}
