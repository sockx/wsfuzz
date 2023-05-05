package cmd

import (
	"os"

	"wsfuzz/internal/ws"

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
	RootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVarP(&localFile, "file", "f", "req.txt", "HTTP请求保存所在文件")
}

var testCmd = &cobra.Command{
	Use:   "fuzz",
	Short: "Run fuzz",
	Run: func(cmd *cobra.Command, args []string) {
		// 读文件
		req := ws.ParseFile(localFile)
		if req == nil {
			return
		}

		for i := 470; i < 480; i++ {
			ws.SendData(req, i)
		}
	},
}
