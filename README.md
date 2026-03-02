# Zenos Core Plugin

Zenos Core Plugin is the foundational library for developing music plugins within the **Zenos** ecosystem. It defines standard interfaces and provides a generic API server to facilitate the integration of new audio sources.

## Installation

```bash
go get github.com/0ddityspace/zenos-core-plugin
```

## Features

- **`Controller` Interface**: Defines essential actions a plugin must perform (play, pause, skip, volume, etc.).
- **`NowPlaying` Model**: A standardized data structure for reporting playback metadata to the frontend.
- **API Server**: A ready-to-use HTTP wrapper to quickly expose any `Controller` implementation.

## Project Structure

- `pluginapi/`: Universal interface definitions and data models.
- `server/`: Generic HTTP server implementation to expose a `Controller`.

## Usage

Implement the `Controller` interface:

```go
type MyPlugin struct { ... }

func (p *MyPlugin) Play() error { ... }
// ... implement other methods
```

Then register it with the provided server:

```go
ctrl := &MyPlugin{}
logger := slog.Default()
api := server.NewAPI(ctrl, logger)

mux := http.NewServeMux()
api.RegisterRoutes(mux)

http.ListenAndServe(":8080", mux)
```

## License

This project is licensed under the MIT License.
