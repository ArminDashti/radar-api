# Endpoints

| Method | Path | Purpose | Authentication |
|---|---|---|---|
| POST | `/api/auth/login` | Exchange credentials for JWT | Public |
| GET | `/api/probes` | List probes | JWT |
| GET | `/api/hosts` | List hosts | JWT |
| POST | `/api/hosts` | Create host | JWT |
| PUT | `/api/hosts/:id` | Update host | JWT |
| GET | `/api/grid/hosts` | Host latency grid | JWT |
| GET | `/api/grid/probes` | Probe latency grid | JWT |
| GET | `/api/agent/targets` | List active agent targets | Agent token |
| POST | `/api/agent/samples` | Upsert minute samples | Agent token |
