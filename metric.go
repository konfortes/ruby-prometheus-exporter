package main

type metricLabels map[string]string

type counter struct {
	Help        string
	Name        string
	ConstLabels metricLabels
	Labels      metricLabels
	Value       int
}

func (c counter) labelKeys() []string {
	keys := []string{}
	for key := range c.Labels {
		keys = append(keys, key)
	}

	return keys
}

func fromRequest(request RequestBody) counter {
	c := counter{}
	c.ConstLabels = make(metricLabels)
	c.Labels = make(metricLabels)

	for key, val := range request.CustomLabels {
		c.ConstLabels[key] = val
	}

	for key, val := range request.Keys {
		c.Labels[key] = val
	}

	c.Name = request.Name
	c.Help = request.Help
	c.Value = request.Value

	return c
}
