package pkg

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/mniak/ytlive/internal"
	"github.com/pkg/errors"
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
	ID          string
	Title       string
	Description string
	Date        time.Time
	Link        string
	StreamName  string
	StreamKey   string
	StreamURL   string
}

func Schedule(options ScheduleRequest) (result ScheduleResponse, err error) {

	svc, err := internal.CreateYoutubeClient()
	if err != nil {
		return
	}

	stream, err := svc.LiveStreams.Insert(
		[]string{
			"Snippet",
			"Cdn",
			"ContentDetails",
		},
		&youtube.LiveStream{
			Snippet: &youtube.LiveStreamSnippet{
				Title: fmt.Sprintf("[%s] Generated Key (ytlive)", options.Date.Format("2006-01-02")),
			},
			Cdn: &youtube.CdnSettings{
				IngestionType: "rtmp",
				Resolution:    "variable",
				FrameRate:     "variable",
			},
			ContentDetails: &youtube.LiveStreamContentDetails{
				IsReusable: false,
				ForceSendFields: []string{
					"IsReusable",
				},
			},
		},
	).Do()

	if err != nil {
		err = errors.Wrap(err, "could not create a new stream")
		return
	}

	result.StreamURL = stream.Cdn.IngestionInfo.IngestionAddress
	result.StreamKey = stream.Cdn.IngestionInfo.StreamName
	result.StreamName = stream.Snippet.Title

	broadcast, err := svc.LiveBroadcasts.Insert(
		[]string{
			"Snippet",
			"Status",
			"ContentDetails",
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
				ForceSendFields: []string{
					"SelfDeclaredMadeForKids",
				},
			},
			ContentDetails: &youtube.LiveBroadcastContentDetails{
				EnableAutoStart: options.AutoStart,
				EnableAutoStop:  options.AutoStop,
				EnableDvr:       options.DVR,

				ForceSendFields: []string{
					"AutoStart",
					"AutoStop",
					"EnableDvr",
				},
			},
		},
	).Do()

	if err != nil {
		err = errors.Wrap(err, "could not create a new broadcast")
		return
	}

	broadcast, err = svc.LiveBroadcasts.Bind(broadcast.Id, []string{}).
		StreamId(stream.Id).
		Do()

	if err != nil {
		err = errors.Wrap(err, "could not bind broadcast to stream")
		return
	}

	scheduledDate, _ := dateparse.ParseAny(broadcast.Snippet.ScheduledStartTime)
	result.ID = broadcast.Id
	result.Title = broadcast.Snippet.Title
	result.Link = fmt.Sprintf("https://youtu.be/%s", broadcast.Id)
	result.Date = scheduledDate
	return
}
