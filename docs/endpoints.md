# Endpoints

| Method | Path | Purpose | Authentication |
|---|---|---|---|
| POST | `/api/auth/login` | Exchange credentials for JWT | Public |
| GET | `/api/probes` | List probes | JWT |
| GET | `/api/endpoints` | List endpoints | JWT |
| POST | `/api/endpoints` | Create endpoint | JWT |
| GET | `/api/grid/endpoints` | Endpoint latency grid | JWT |
| GET | `/api/grid/probes` | Probe latency grid | JWT |
| GET | `/api/agent/targets` | List active agent targets | Agent token |
| POST | `/api/agent/samples` | Upsert minute samples | Agent token |
