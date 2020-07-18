package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

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
			handleRequestError(err, w, http.StatusBadRequest, "Invalid request")
			return
		}

		cntr := fromRequest(requestBody)
		var promCounter prometheus.Counter
		if promCounter, err = registry.getOrRegister(cntr); err != nil {
			msg := fmt.Sprintf("error getOrRegister metric: %s", cntr.Name)
			handleRequestError(err, w, http.StatusBadRequest, msg)
			return
		}

		promCounter.Inc()
	}
}
