package cmd

import (
	"bic-cd/internal"
	"bic-cd/internal/model"
	"bic-cd/pkg/config"
	"bic-cd/pkg/db"
	"bic-cd/pkg/jwt"
	"bic-cd/pkg/log"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	rootCmd = &cobra.Command{
		Use:   "help",
		Short: "this is a game server and gate, you can run with cmd: game or gate",
	}
	bicCmd = &cobra.Command{
		Use:   "bic",
		Short: "start a bic server",
		Long:  `bic server for code template`,
		PreRun: func(cmd *cobra.Command, args []string) {
			conf, err := cmd.Flags().GetString("conf")
			if err != nil {
				panic("need bic server config")
			}
			config.SetupYaml(conf)
			log.Setup()
			db.Setup()
			jwt.SetSecret([]byte(config.GlobalConf.App.JwtSecret))
			model.Setup()
		},
		Run: func(cmd *cobra.Command, args []string) {
			port, err := cmd.Flags().GetUint16("port")
			if err != nil {
				panic("need bic server port")
			}
			engine := internal.Setup()
			startHttpServer(engine, int16(port))
		},
	}
	server *http.Server
)

func init() {
	bicCmd.Flags().Uint16P("port", "p", 6996, "bic server port")
	bicCmd.Flags().StringP("conf", "c", "./conf/app.yaml", "bic server config")
}

func stop() {
	_ = server.Close()
	log.Stop()
}

func startHttpServer(engine *gin.Engine, port int16) {
	server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      engine,
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		defer stop()
		for {
			s := <-c
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				//终端主动退出。（Ctrl+C）、（Ctrl+/）、（KILL + PID）
				return
			case syscall.SIGHUP:
				//终端控制进程结束（终端连接断开）
				return
			}
		}
	}()
	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server closed under request")
		} else {
			panic("Server closed unexpect")
		}
	}
}

func init() {
	rootCmd.AddCommand(bicCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
