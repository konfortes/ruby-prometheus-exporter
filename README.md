# ruby-prometheus-exporter-go

a Golang port of [discourse/prometheus_exporter's](https://github.com/discourse/prometheus_exporter) exporter web server.  
The origin project's web server is written with WEBRick, make it hard to use in production under high throughput.

## Configuration

Listen address can be configured by setting the `SERVER_HOST` and `SERVER_PORT` environment variables.

## Logging

This app logs using the zap logger. To log with json encoder, set `GO_ENV` environment variable to `production`

## API

- `/metrics`: Returns all metrics in a text format (Prometheus format).
- `/send-metrics` Receives a metric with action, registers it (if not already registered) and acts on it. (_currently only increment of counter is supported_).  
**Important**: currently only http's `Transfer-Encoding: chunked` requests are supported. see `client.rb` for example.

## Releases

a Github action of build and push image is triggered on tag release. images can be pulled by:

```bash
docker pull konfortes/ruby-prometheus-exporter-go:v1.0.0
```
