package main

import "github.com/prometheus/client_golang/prometheus"

type Metric interface {
	prometheus.Metric
}

type metric struct{}

func (m *metric) increment(name string) {

}
