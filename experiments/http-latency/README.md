# HTTP Latency / Prometheus Histogram Lab

## Goal

Verify the API's Prometheus request metrics and use its histogram to estimate
latency percentiles for a small local workload.

## Workload

```yaml
Endpoint: GET /health
Requests: 1,000
Concurrency: 20
Instrumentation: Prometheus Histogram
```

The workload was generated locally. Latency was measured server-side by the Go
HTTP middleware, so this was not an end-to-end network benchmark.

## Results

| Percentile | Estimated latency |
| --- | ---: |
| p50 | 0.0258 ms |
| p95 | 0.0491 ms |
| p99 | 0.1333 ms |

`http_request_duration_seconds_count{method="GET"}` was `1000`.

These local results are estimates derived from Prometheus histogram buckets,
not exact request timings or production performance measurements.

## PromQL

p50:

```promql
histogram_quantile(0.50, sum by (le) (rate(http_request_duration_seconds_bucket{method="GET"}[5m])))
```

p95:

```promql
histogram_quantile(0.95, sum by (le) (rate(http_request_duration_seconds_bucket{method="GET"}[5m])))
```

p99:

```promql
histogram_quantile(0.99, sum by (le) (rate(http_request_duration_seconds_bucket{method="GET"}[5m])))
```

## Notes

The API exposes:

- Database pool state: max open, open, in-use, and idle connections
- Database pool contention: wait count and cumulative wait duration
- HTTP request count by method and status
- HTTP request duration histogram by method

Prometheus runs from Docker Compose and scrapes the API's `/metrics` endpoint
using [`../../monitoring/prometheus.yml`](../../monitoring/prometheus.yml).
