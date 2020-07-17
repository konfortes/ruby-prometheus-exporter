package main

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metricsRegistry struct {
	sync.Map
}

func (mr *metricsRegistry) getOrStore(c counter) (prometheus.Counter, error) {
	promCounter, found := mr.Load(c.Name)

	if found {
		assertedCounter, ok := promCounter.(prometheus.Counter)
		if !ok {
			return nil, fmt.Errorf("unable to assert type of %+v", promCounter)
		}
		return assertedCounter, nil
	}

	return mr.register(c), nil
}

func (mr *metricsRegistry) register(c counter) prometheus.Counter {
	promCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name:        c.Name,
		Help:        c.Help,
		ConstLabels: prometheus.Labels(c.ConstLabels),
	})

	mr.Store(c.Name, promCounter)

	return promCounter
}

var (
	registry = metricsRegistry{}
)
