# Zenos Core Plugin

Zenos Core Plugin is the foundational library for developing music plugins within the **Zenos** ecosystem. It defines standard interfaces and provides a generic HTTP API server to facilitate the integration of new audio sources.

## Installation

```bash
go get github.com/0ddityspace/zenos-core-plugin
```

## Features

- **`Controller` Interface**: Defines essential actions a plugin must perform (play, pause, skip, volume, etc.).
- **`Manifest`**: A standardized struct to describe a plugin (name, type, author, repo, image, version).
- **`NowPlaying` Model**: A standardized data structure for reporting playback metadata.
- **API Server**: A ready-to-use HTTP wrapper to quickly expose any `Controller` implementation.

## API Endpoints

All endpoints are automatically registered by `server.RegisterRoutes`.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/v1/manifest` | Returns plugin metadata (name, type, repo, author, image, version) |
| `GET` | `/v1/status` | Returns engine health (`ok`, `engine_up`) |
| `GET` | `/v1/now-playing` | Returns current playback state and track metadata |
| `POST` | `/v1/player/play` | Resume playback |
| `POST` | `/v1/player/pause` | Pause playback |
| `POST` | `/v1/player/stop` | Stop playback |
| `POST` | `/v1/player/next` | Skip to next track |
| `POST` | `/v1/player/previous` | Go to previous track |
| `POST` | `/v1/player/seek?offset=<ms>` | Seek to position (milliseconds) |
| `POST` | `/v1/player/volume?level=<0-100>` | Set volume level |

## Project Structure

- `pluginapi/`: Universal interface definitions and data models (`Controller`, `Manifest`, `NowPlaying`).
- `server/`: Generic HTTP server implementation to expose a `Controller`.

## Usage

Implement the `Controller` interface:

```go
type MyPlugin struct{}

func (p *MyPlugin) GetManifest() pluginapi.Manifest {
    return pluginapi.Manifest{
        Type:    "audio_source",
        Name:    "My Plugin",
        Repo:    "github.com/you/my-plugin",
        Author:  "you",
        Image:   "https://example.com/icon.png",
        Version: "1.0.0",
    }
}

func (p *MyPlugin) IsUp() bool       { return true }
func (p *MyPlugin) Play() error      { return nil }
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

You can also add custom routes before registering:

```go
api.AddRoute("GET /v1/artwork", myArtworkHandler)
api.RegisterRoutes(mux)
```

## License

This project is licensed under the MIT License.
