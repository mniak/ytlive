package pkg

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/mniak/ytlive/internal"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type ScheduleRequest struct {
	Title string
	// Visibility   string
	Description string
	// Category     string
	Date        time.Time
	MadeForKids bool

	AutoStart bool
	AutoStop  bool
	DVR       bool
}

type ScheduleResponse struct {
	ID            string
	Title         string
	Description   string
	Date          time.Time
	Link          string
	StreamKeyName string
	StreamKey     string
	StreamURL     string
}

func Schedule(options ScheduleRequest) (result ScheduleResponse, err error) {

	config := internal.GetConfig("", "")
	ctx, tokenSource := internal.GetTokenSource(config)

	svc, err := youtube.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		errors.Wrap(err, "could not create Youtube API client")
		return
	}

	stream, err := svc.LiveStreams.Insert(
		[]string{
			"snippet",
			"cdn",
			"contentDetails",
		},
		&youtube.LiveStream{
			Snippet: &youtube.LiveStreamSnippet{
				Title: fmt.Sprintf("[%s] Generated Key (ytlive)", options.Date.Format("2006-01-02")),
			},
			Cdn: &youtube.CdnSettings{
				FrameRate:     "30fps",
				IngestionType: "rtmp",
				Resolution:    "1080p",
			},
			ContentDetails: &youtube.LiveStreamContentDetails{
				IsReusable: false,
			},
		},
	).Do()

	if err != nil {
		err = errors.Wrap(err, "could not create a new stream")
		return
	}

	result.StreamURL = stream.Cdn.IngestionInfo.IngestionAddress
	result.StreamKey = stream.Cdn.IngestionInfo.StreamName
	result.StreamKeyName = stream.Snippet.Title

	broadcast, err := svc.LiveBroadcasts.Insert(
		[]string{
			"snippet",
			"status",
			"contentDetails",
		},
		&youtube.LiveBroadcast{
			Snippet: &youtube.LiveBroadcastSnippet{
				Title:              options.Title,
				Description:        options.Description,
				ScheduledStartTime: options.Date.Format(time.RFC3339),
			},
			Status: &youtube.LiveBroadcastStatus{
				PrivacyStatus:           "public",
				SelfDeclaredMadeForKids: false,
			},
			ContentDetails: &youtube.LiveBroadcastContentDetails{
				EnableAutoStart: options.AutoStart,
				EnableAutoStop:  options.AutoStop,
				EnableDvr:       options.DVR,
			},
		},
	).Do()

	if err != nil {
		err = errors.Wrap(err, "could not create a new stream")
		return
	}

	scheduledDate, _ := dateparse.ParseAny(broadcast.Snippet.ScheduledStartTime)
	result.ID = broadcast.Id
	result.Title = broadcast.Snippet.Title
	result.Link = fmt.Sprintf("https://youtu.be/%s", broadcast.Id)
	result.Date = scheduledDate
	return
}
