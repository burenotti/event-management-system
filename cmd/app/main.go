package main

import (
	"flag"
	"fmt"
	"github.com/burenotti/rtu-it-lab-recruit/handler"
	"github.com/burenotti/rtu-it-lab-recruit/pkg/httpserver"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func getLogger() *logrus.Logger {
	logger := logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.JSONFormatter{
			PrettyPrint: true,
		},
		Level: logrus.InfoLevel,
	}
	return &logger
}

func main() {
	logger := getLogger()

	var host, port string
	flag.StringVar(&host, "host", "0.0.0.0", "Server host")
	flag.StringVar(&port, "port", "80", "Server port")
	flag.Parse()
	addr := fmt.Sprintf("%s:%s", host, port)

	h := handler.New(&handler.Config{
		Name: "RTU ITLab",
	})
	logger.Infof("Server run on %s", addr)
	srv := httpserver.New(addr, h, logger)

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-srv.Notify():
		if err != nil {
			logger.WithError(err).Errorf("Server exited with error: %v", err)
		} else {
			logger.Info("Server exited without errors")
		}
	case s := <-interrupt:
		logger.WithField("signal", s.String()).Infof("%s Signal caught. Shutdown", s.String())
		if err := srv.Shutdown(); err != nil {
			logger.WithError(err).Errorf("Server shutdowned with error: %v", err)
		}
	}

}
