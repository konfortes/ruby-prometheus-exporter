package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: error handling

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

func sendMetrics(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	buf := make([]byte, 256)
	for {
		n, err := req.Body.Read(buf)
		if err == io.EOF {
			break
		}
		var requestBody RequestBody
		if err := json.Unmarshal(buf[:n], &requestBody); err != nil {
			fmt.Printf("error unmarshaling: %s\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cntr := fromRequest(requestBody)
		var promCounter *prometheus.CounterVec
		if promCounter, err = registry.getOrRegister(cntr); err != nil {
			fmt.Printf("error finding metric: %s\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: find better way to do that
		promLabels := prometheus.Labels{}
		for key, val := range cntr.Labels {
			promLabels[key] = val
		}

		// promCounter.With(promLabels).Inc()

		cardinalCounter, err := promCounter.GetMetricWith(promLabels)
		if err != nil {
			fmt.Println("ronen")
			fmt.Println(err)
		}
		cardinalCounter.Inc()
	}
}

func main() {
	host := envOr("SERVER_HOST", "0.0.0.0")
	port := envOr("SERVER_PORT", "9394")
	addr := fmt.Sprintf("%s:%s", host, port)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/send-metrics", sendMetrics)

	log.Println("Listeneing on " + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func envOr(key, defaultValue string) string {
	val, found := os.LookupEnv(key)

	if !found {
		val = defaultValue
	}

	return val
}
