package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julinserg/go_home_project/internal/app"
	"github.com/julinserg/go_home_project/internal/logger"
	internalhttp "github.com/julinserg/go_home_project/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/previewer/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	err := config.Read(configFile)
	if err != nil {
		log.Println("error read config: " + err.Error())
		return
	}

	f, err := os.OpenFile("previewer.logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		log.Println("error opening logfile: " + err.Error())
		return
	}
	defer f.Close()

	logg := logger.New(config.Logger.Level, f)

	dname, err := os.MkdirTemp("", "previewercachedir")
	if err != nil {
		log.Println("error create temp dir for cache: " + err.Error())
		return
	}
	defer os.RemoveAll(dname)

	logg.Info("temp dir for cache pictures is created: " + dname)

	previewer := app.New(logg, config.LRUCache.Size, dname)

	endpoint := net.JoinHostPort(config.HTTP.Host, config.HTTP.Port)
	server := internalhttp.NewServer(logg, previewer, endpoint)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("previewer pictures is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
