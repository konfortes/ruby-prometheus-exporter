package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metricsRegistry map[string]prometheus.Counter

func (mr metricsRegistry) exist(name string) bool {
	_, ok := mr[name]

	return ok
}

var (
	metrics = metricsRegistry{}
)

func (mr metricsRegistry) increment(name string) error {
	if !mr.exist(name) {
		mr.register(name, "")
	}
	mr[name].Inc()
	return nil
}

func (mr metricsRegistry) register(name, help string) {
	if mr.exist(name) {
		// TODO: error?
		return
	}

	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})

	mr[name] = counter
}
