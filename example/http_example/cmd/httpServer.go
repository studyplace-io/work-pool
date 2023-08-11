package cmd

import (
	"github.com/spf13/cobra"
	"github.com/studyplace-io/work-pool-framework/example/http_example/pkg/common"
	"github.com/studyplace-io/work-pool-framework/example/http_example/pkg/server"
)

func httpServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "httpServer",
		Short: "run http server",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := &common.ServerConfig{
				Debug: debug,
				Port:  serverPort,
			}
			// 启动http server
			server.HttpServer(cfg)
		},
	}
	return cmd
}
