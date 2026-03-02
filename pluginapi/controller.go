package pluginapi

// State defines the current playback state.
type State string

const (
	StatePlaying State = "playing"
	StatePaused  State = "paused"
	StateStopped State = "stopped"
)

// NowPlaying groups metadata for the currently playing track.
type NowPlaying struct {
	Status     State  `json:"status"`
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Album      string `json:"album"`
	ArtworkURL string `json:"artwork_url"`
	Format     string `json:"format"`
	SampleRate int    `json:"sample_rate"`
	BitDepth   int    `json:"bit_depth,omitempty"`
	Bitrate    int    `json:"bitrate,omitempty"`
	TimeMS     int    `json:"time_ms"`
	DurationMS int    `json:"duration_ms"`
}

// Manifest contains metadata about the plugin.
type Manifest struct {
	Type    string `json:"type"` // e.g., "audio_source"
	Name    string `json:"name"`
	Repo    string `json:"repo"`
	Author  string `json:"author"`
	Image   string `json:"image"`
	Version string `json:"version"`
}

// Controller is the interface that every Zenos music plugin must implement.
type Controller interface {
	GetManifest() Manifest
	IsUp() bool
	GetNowPlaying() NowPlaying
	Play() error
	Pause() error
	Stop() error
	Next() error
	Previous() error
	Seek(offsetMS int) error
	SetVolume(level float64) error // level range: 0.0 to 100.0
}
