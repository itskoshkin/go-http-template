# Go HTTP Template

Project template for a Go web application

## What's included

- Clean project structure with constructor-based dependency injection
- [Gin](https://github.com/gin-gonic/gin) HTTP server with graceful shutdown
- [Viper](https://github.com/spf13/viper) config with YAML file, ENV overrides, validation and defaults
- Custom logger with console/file output and text/JSON formats
- Embedded static files, templates, fonts, and icons

## Project structure

```text
.
├── cmd/main.go                        # Entrypoint
├── internal/
│   ├── app/app.go                     # Bootstrap and DI wiring
│   ├── config/config.go               # Config loading and validation
│   ├── logger/logger.go               # Logger setup
│   ├── utils/
│   │   ├── colors/                    # Color codes for console
│   │   ├── gin/                       # Gin logging middlewares and formatters
│   │   └── text/                      # Colored text helpers
│   └── webserver/
│       ├── server.go                  # Server setup, middlewares, routes, graceful shutdown
│       └── handlers/handlers.go       # HTTP handlers
├── static/
│   ├── embed.go                       # Embedded FS for templates and public assets
│   ├── templates/                     # Go HTML templates
│   ├── styles/                        # CSS
│   ├── scripts/                       # JS
│   └── assets/
│       ├── fonts/                     # Fonts
│       ├── icons/                     # Favicons and icons
│       └── images/                    # Images
├── example_config.yaml                # Config template
└── go.mod
```

## Config

Config is loaded from `config.yaml` (override path via `CONFIG_PATH` env).

Values can also be set via environment variables (take precedence over config file):

| Config key                 | Env variable                  | Default           | Description                                                 |
|----------------------------|-------------------------------|-------------------|-------------------------------------------------------------|
| `app.web.host`             | `APP_HOST`                    | `0.0.0.0`         | Bind address                                                |
| `app.web.port`             | `APP_PORT`                    | `8080`            | Listening port                                              |
| `app.web.gin_release_mode` | `GIN_RELEASE_MODE`            | `true`            | Hide Gin debug output                                       |
| `app.web.shutdown_timeout` | `WEB_SERVER_SHUTDOWN_TIMEOUT` | `5s`              | Graceful shutdown timeout                                   |
| `app.log.level`            | `LOG_LEVEL`                   | `INFO`            | `DEBUG`, `INFO`, `WARN`, `ERROR`                            |
| `app.log.log_format`       | `LOG_FORMAT`                  | `text`            | `text` or `json`                                            |
| `app.log.log2console`      | `LOG_TO_CONSOLE`              | `true`            | Log to stdout                                               |
| `app.log.log2file`         | `LOG_TO_FILE`                 | `true`            | Log to file                                                 |
| `app.log.file_path`        | `LOG_FILE_PATH`               | `application.log` | Log file path                                               |
| `app.log.file_mode`        | `LOG_FILE_MODE`               | `append`          | `append`, `overwrite`, or `rotate`                          |
| `app.log.old_logs_folder`  | `LOG_FILES_FOLDER`            | —                 | Folder for rotated logs (required when `file_mode: rotate`) |


## Build & Run

### Prerequisites

- [Go 1.21+](https://gist.github.com/itskoshkin/5f45fca15c30f859955dc146080a00d9)

### Local

1. Build
   ```bash
   go build -o go-http-template ./cmd/main.go
   ```
2. Prepare and edit `config.yaml`
   ```bash
   cp example_config.yaml config.yaml && nano config.yaml
   ```
3. Run
   ```bash
   ./go-http-template
   ```
   or just
   ```bash
   go run ./cmd
   ```

### Docker

1. Build image
   ```bash
   docker build -t go-http-template .
   ```
2. Prepare and edit `config.yaml` (or skip and pass ENV later)
   ```bash
   cp example_config.yaml config.yaml && nano config.yaml
   ```
3. Run container
   ```bash
   docker run -d --name go-http-template \
		-v $(pwd)/config.yaml:/app/config.yaml:ro \
		-v $(pwd)/logs:/app/logs \
		go-http-template
   ```
   Or without config (uses defaults baked into image, override via ENV)
   ```bash
   docker run -d --name go-http-template \
		-v $(pwd)/logs:/app/logs \
		-e APP_PORT=8080 \
		go-http-template
   ```
