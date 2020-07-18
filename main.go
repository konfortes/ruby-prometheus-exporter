package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/konfortes/go-server-utils/server"
	"github.com/konfortes/go-server-utils/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// RequestBody ...
type RequestBody struct {
	MetricType               string       `json:"type"`
	Help                     string       `json:"help"`
	Name                     string       `json:"name"`
	Keys                     metricLabels `json:"keys"`
	Value                    int          `json:"value"`
	PrometheusExporterAction string       `json:"prometheus_exporter_action"`
	CustomLabels             metricLabels `json:"custom_labels"`
}

var (
	logger *zap.SugaredLogger
)

func main() {
	logger = initLogger()
	defer logger.Sync()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/send-metrics", sendMetrics)

	addr := resolveAddress()

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		logger.Infof("Listeneing on " + addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	server.GracefulShutdown(srv)
}

func initLogger() *zap.SugaredLogger {
	var l *zap.Logger
	var err error

	if os.Getenv("GO_ENV") == "production" {
		l, err = zap.NewProduction()
	} else {
		l, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Panic(err)
	}
	return l.Sugar()
}

func resolveAddress() string {
	host := utils.GetEnvOr("SERVER_HOST", "0.0.0.0")
	port := utils.GetEnvOr("SERVER_PORT", "9394")
	return fmt.Sprintf("%s:%s", host, port)
}

func handleRequestError(err error, w http.ResponseWriter, httpStatus int, msg string) {
	message := fmt.Sprintf("%s: %s", msg, err)
	logger.Error(message)
	http.Error(w, msg, httpStatus)
}
