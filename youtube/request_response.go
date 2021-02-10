package youtube

import "time"

type GenerateRequest struct {
	Title string
	// Visibility   string
	Description string
	// Category     string
	Date        time.Time
	MadeForKids bool

	NewStreamKey  bool
	StreamKeyName string

	AutoStart bool
	AutoStop  bool
	DVR       bool
}

type GenerateResponse struct {
	ID            string
	Title         string
	Description   string
	Date          time.Time
	Link          string
	StreamKeyName string
	StreamKey     string
	StreamURL     string
}
