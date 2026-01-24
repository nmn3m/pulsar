# Pulsar Observability Stack

This folder contains the observability infrastructure for Pulsar, built on the VictoriaMetrics ecosystem with OpenTelemetry.

## Architecture

```
┌─────────────────┐
│   Pulsar App    │
│   (Backend)     │
└────────┬────────┘
         │ OTLP (gRPC/HTTP)
         ▼
┌─────────────────┐
│    OpenTelemetry│
│    Collector    │
└────────┬────────┘
         │
    ┌────┼────┬────────────┐
    ▼    ▼    ▼            ▼
┌──────┐┌──────┐┌────────┐┌───────┐
│Metrics││Logs ││ Traces ││Debug  │
└──┬───┘└──┬───┘└───┬────┘└───────┘
   ▼       ▼        ▼
┌──────────────┐┌──────────────┐┌──────────────┐
│VictoriaMetrics││VictoriaLogs ││VictoriaTraces│
└──────┬───────┘└──────┬───────┘└──────┬───────┘
       └───────────────┼───────────────┘
                       ▼
               ┌───────────────┐
               │    Grafana    │
               └───────────────┘
```

## Quick Start

### 1. Start the Stack

```bash
cd instrumentation
docker compose up -d
```

### 2. Enable Telemetry in Backend

Set the following environment variables in your backend:

```bash
OTEL_ENABLED=true
OTEL_SERVICE_NAME=pulsar-backend
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317  # Use localhost:4317 if running locally
OTEL_EXPORTER_OTLP_PROTOCOL=grpc
OTEL_ENVIRONMENT=development
```

Or update `docker-compose.yml` in the project root to set `OTEL_ENABLED=true`.

### 3. Access the Services

| Service | URL | Credentials |
|---------|-----|-------------|
| Grafana | http://localhost:3000 | admin / admin |
| VictoriaMetrics | http://localhost:8428 | - |
| VictoriaLogs | http://localhost:9428 | - |
| VictoriaTraces | http://localhost:10428 | - |

## Services

### Grafana (Port 3000)

The main visualization interface for all telemetry data.

**Access:** http://localhost:3000

**Default Credentials:**
- Username: `admin`
- Password: `admin`

**Pre-configured Datasources:**
- **VictoriaMetrics** - For metrics (Prometheus-compatible)
- **VictoriaLogs** - For logs
- **VictoriaTraces** - For traces (Jaeger-compatible)

**Usage:**
1. Navigate to **Explore** in the left sidebar
2. Select a datasource from the dropdown:
   - Use **VictoriaMetrics** to query metrics with PromQL
   - Use **VictoriaTraces** to search and view distributed traces
   - Use **VictoriaLogs** to query logs with LogsQL

---

### VictoriaMetrics (Port 8428)

Time-series database for storing metrics. Prometheus-compatible.

**Access:** http://localhost:8428

**Useful Endpoints:**
| Endpoint | Description |
|----------|-------------|
| `/vmui` | Built-in web UI for querying metrics |
| `/api/v1/query` | Instant query API |
| `/api/v1/query_range` | Range query API |
| `/api/v1/labels` | List all label names |
| `/metrics` | VictoriaMetrics own metrics |

**Example Queries (PromQL):**
```promql
# HTTP request rate
rate(http_server_request_duration_seconds_count[5m])

# Request latency P99
histogram_quantile(0.99, rate(http_server_request_duration_seconds_bucket[5m]))

# Error rate
sum(rate(http_server_request_duration_seconds_count{http_status_code=~"5.."}[5m]))
```

**Web UI:** http://localhost:8428/vmui

---

### VictoriaLogs (Port 9428)

Log aggregation and storage solution.

**Access:** http://localhost:9428

**Useful Endpoints:**
| Endpoint | Description |
|----------|-------------|
| `/select/vmui` | Built-in web UI for querying logs |
| `/select/logsql/query` | LogsQL query API |
| `/select/logsql/tail` | Live tail API |
| `/select/logsql/stats_query` | Stats query API |

**Example Queries (LogsQL):**
```logsql
# All logs from pulsar-backend
service.name:pulsar-backend

# Error logs only
service.name:pulsar-backend AND level:error

# Logs containing specific text
service.name:pulsar-backend AND "database connection"

# Logs from last 1 hour
service.name:pulsar-backend AND _time:1h
```

**Web UI:** http://localhost:9428/select/vmui

---

### VictoriaTraces (Port 10428)

Distributed tracing backend. Compatible with Jaeger Query API.

**Access:** http://localhost:10428

**Useful Endpoints:**
| Endpoint | Description |
|----------|-------------|
| `/vmui` | Built-in web UI for exploring traces |
| `/api/traces` | Jaeger-compatible trace query API |
| `/api/services` | List all traced services |
| `/api/services/{service}/operations` | List operations for a service |

**Features:**
- Search traces by service, operation, tags, duration
- View trace timelines and span details
- Analyze service dependencies
- Compare traces

**Web UI:** http://localhost:10428/vmui

---

### OpenTelemetry Collector (Ports 4317, 4318)

Receives, processes, and exports telemetry data.

**Endpoints:**
| Port | Protocol | Description |
|------|----------|-------------|
| 4317 | gRPC | OTLP gRPC receiver |
| 4318 | HTTP | OTLP HTTP receiver |
| 8888 | HTTP | Collector metrics |
| 13133 | HTTP | Health check |

**Health Check:** http://localhost:13133

**Send Test Data:**
```bash
# Using curl to send a test log via OTLP HTTP
curl -X POST http://localhost:4318/v1/logs \
  -H "Content-Type: application/json" \
  -d '{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"test"}}]},"scopeLogs":[{"logRecords":[{"body":{"stringValue":"Test log message"}}]}]}]}'
```

## Configuration

### Retention Period

All services are configured with 30-day retention. To modify:

**docker-compose.yml:**
```yaml
# VictoriaMetrics
command:
  - "--retentionPeriod=30d"  # Change to desired period

# VictoriaLogs
command:
  - "--retentionPeriod=30d"

# VictoriaTraces
command:
  - "--retentionPeriod=30d"
```

### Adding Custom Dashboards

Place dashboard JSON files in:
```
grafana/provisioning/dashboards/
```

Dashboards will be automatically loaded on Grafana startup.

### Modifying Collector Pipeline

Edit `otel-collector-config.yaml` to:
- Add new receivers (e.g., Prometheus scrape)
- Modify processors (e.g., filtering, sampling)
- Add exporters (e.g., additional backends)

## Troubleshooting

### Check Service Health

```bash
# Check all containers
docker compose ps

# Check collector health
curl http://localhost:13133

# Check VictoriaMetrics
curl http://localhost:8428/health

# Check VictoriaLogs
curl http://localhost:9428/health

# Check VictoriaTraces
curl http://localhost:10428/health
```

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f otel-collector
docker compose logs -f victoriametrics
docker compose logs -f victorialogs
docker compose logs -f victoriatraces
docker compose logs -f grafana
```

### No Data Appearing

1. **Check backend telemetry is enabled:**
   ```bash
   # Verify OTEL_ENABLED=true in backend environment
   ```

2. **Check collector is receiving data:**
   ```bash
   docker compose logs otel-collector | grep -i "traces\|metrics\|logs"
   ```

3. **Check network connectivity:**
   ```bash
   # From backend container, test collector connectivity
   docker exec pulsar-backend nc -zv otel-collector 4317
   ```

4. **Verify collector config:**
   ```bash
   docker compose exec otel-collector cat /etc/otelcol-contrib/config.yaml
   ```

### Reset All Data

```bash
docker compose down -v
docker compose up -d
```

## Useful Links

- [VictoriaMetrics Documentation](https://docs.victoriametrics.com/)
- [VictoriaLogs Documentation](https://docs.victoriametrics.com/victorialogs/)
- [VictoriaTraces Documentation](https://docs.victoriametrics.com/victoriatraces/)
- [OpenTelemetry Collector Documentation](https://opentelemetry.io/docs/collector/)
- [Grafana Documentation](https://grafana.com/docs/grafana/latest/)
- [PromQL Cheat Sheet](https://promlabs.com/promql-cheat-sheet/)
- [LogsQL Documentation](https://docs.victoriametrics.com/victorialogs/logsql/)
