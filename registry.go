package main

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type metricsRegistry struct {
	sync.Map
}

func (mr *metricsRegistry) get(c counter) (*prometheus.CounterVec, error) {
	promCounter, found := mr.Load(c.Name)

	if found {
		assertedCounter, ok := promCounter.(*prometheus.CounterVec)
		if !ok {
			return nil, fmt.Errorf("unable to assert type of %+v", promCounter)
		}
		return assertedCounter, nil
	}

	return nil, nil
}

func (mr *metricsRegistry) getOrRegister(c counter) (*prometheus.CounterVec, error) {
	val, _ := mr.LoadOrStore(c.Name, mr.registerCounter(c))

	assertedCounter, ok := val.(*prometheus.CounterVec)
	if !ok {
		return nil, fmt.Errorf("unable to assert type of %+v", val)
	}
	return assertedCounter, nil
}

func (mr *metricsRegistry) registerCounter(c counter) *prometheus.CounterVec {
	alreadyRegistered, _ := mr.get(c)
	if alreadyRegistered != nil {
		return alreadyRegistered
	}

	// promCounter := promauto.NewCounter(prometheus.CounterOpts{
	// 	Name:        c.Name,
	// 	Help:        c.Help,
	// 	ConstLabels: prometheus.Labels(c.ConstLabels),
	// })

	// using promauto.NewCounter returns already registered errors
	// although it should be thread safe because It is being called only on sync.Map.LoadOrStore
	// TODO: fix the race condition and change to promauto

	// TODO: find better way to do that
	promLabelKeys := []string{}
	for key := range c.Labels {
		promLabelKeys = append(promLabelKeys, key)
	}
	promCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        c.Name,
		Help:        c.Help,
		ConstLabels: prometheus.Labels(c.ConstLabels),
	}, promLabelKeys)
	prometheus.Register(promCounter)

	return promCounter
}

var (
	registry = metricsRegistry{}
)
