# Autotask Datasource for Grafana

A Grafana backend datasource plugin that connects to the [Kaseya Autotask](https://www.autotask.net/) PSA REST API, allowing you to query and visualize PSA data in Grafana dashboards.

## Features

- **Entity types**: Query Tickets, Companies, Contacts, and Resources
- **Time range mapping**: Map Grafana's time picker to Autotask date fields (e.g. `createDate`, `dueDateTime`)
- **Filtering**: Pass Autotask query filter JSON to narrow results
- **Health check**: Validates API credentials via Zone Information endpoint
- **Secure credentials**: API secret and integration code stored in Grafana's encrypted secret store

## Requirements

- Grafana 10.0+
- An Autotask PSA account with API access
- API integration credentials (username, secret, integration code)

## Installation

### From release

1. Download the latest release archive
2. Extract to your Grafana plugins directory (e.g. `/var/lib/grafana/plugins/`)
3. Restart Grafana
4. Allow unsigned plugins by adding to `grafana.ini`:
   ```ini
   [plugins]
   allow_loading_unsigned_plugins = wyretech-autotask-datasource
   ```

### From source

```bash
# Frontend
npm install
npm run build

# Backend (for your target OS)
go build -o dist/autotask-datasource_linux_amd64 ./pkg
# or for macOS:
go build -o dist/autotask-datasource_darwin_arm64 ./pkg
```

## Configuration

1. In Grafana, go to **Connections → Data sources → Add data source**
2. Search for "Autotask"
3. Fill in:
   - **API URL**: Your Autotask API base URL (e.g. `https://webservices6.autotask.net`). If unsure, use `https://webservices.autotask.net` — the zone is auto-detected.
   - **Username**: Your Autotask API username (email)
   - **API Secret**: Your Autotask API secret
   - **Integration Code**: Your Autotask API integration code
4. Click **Save & Test** to verify the connection

## Query Editor

| Field | Description |
|-------|-------------|
| **Entity** | The Autotask entity type to query (Tickets, Companies, Contacts, Resources) |
| **Time Field** | Optional — map a date field to the Grafana time range picker for filtering |
| **Filter** | Optional — Autotask query filter as JSON |

### Filter examples

```json
{"op":"eq","field":"status","value":1}
```

```json
{"op":"and","items":[{"op":"eq","field":"status","value":1},{"op":"eq","field":"priority","value":2}]}
```

## Development

```bash
# Start development environment (Docker + Grafana)
npm run server

# Frontend dev mode with hot reload
npm run dev

# Build everything
npm run build
go build -o dist/autotask-datasource_darwin_arm64 ./pkg
```

## License

Apache License 2.0 — see [LICENSE](LICENSE).
