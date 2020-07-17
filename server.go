package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: error handling

// RequestBody ...
type RequestBody struct {
	MetricType string `json:"type"`
	Help       string `json:"help"`
	Name       string `json:"name"`
	Keys       struct {
		Link string `json:"link"`
	} `json:"keys"`
	Value                    int    `json:"value"`
	PrometheusExporterAction string `json:"prometheus_exporter_action"`
	CustomLabels             struct {
		App string `json:"app"`
		Env string `json:"env"`
	} `json:"custom_labels"`
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
		if unmarshalError := json.Unmarshal(buf[:n], &requestBody); unmarshalError != nil {
			fmt.Printf("error unmarshaling: %s\n", unmarshalError)
			http.Error(w, unmarshalError.Error(), http.StatusBadRequest)
			return
		}

		if incErr := metrics.increment(requestBody.Name); incErr != nil {
			fmt.Printf("error unmarshaling: %s\n", incErr)
			http.Error(w, incErr.Error(), http.StatusBadRequest)
			return
		}
	}
}

func main() {
	host := envOr("HOST", "0.0.0.0")
	port := envOr("PORT", "9394")
	addr := fmt.Sprintf("%s:%s", host, port)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/send-metrics", sendMetrics)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func envOr(key, defaultValue string) string {
	val, found := os.LookupEnv(key)

	if !found {
		val = defaultValue
	}

	return val
}
