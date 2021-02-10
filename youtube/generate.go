package youtube

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/araddon/dateparse"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	yt "google.golang.org/api/youtube/v3"
)

func Generate(options GenerateRequest) (result GenerateResponse, err error) {

	ctx, tokenSource := getTokenSource()

	svc, err := yt.NewService(ctx, option.WithTokenSource(tokenSource))
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
		&yt.LiveStream{
			Snippet: &yt.LiveStreamSnippet{
				Title: fmt.Sprintf("[%s] Generated Key", options.Date.Format("2006-01-02")),
			},
			Cdn: &yt.CdnSettings{
				FrameRate:     "30fps",
				IngestionType: "rtmp",
				Resolution:    "1080p",
			},
			ContentDetails: &yt.LiveStreamContentDetails{
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

	// broadcast, err := svc.LiveBroadcasts.Insert(
	// 	[]string{
	// 		"snippet",
	// 		"cdn",
	// 		"contentDetails",
	// 	},
	// 	&yt.LiveBroadcast{
	// 		Snippet: &yt.LiveBroadcastSnippet{
	// 			Title: fmt.Sprintf("[%s] Generated Key", options.Date.Format("2006-01-02")),
	// 		},
	// 		Cdn: &yt.CdnSettings{
	// 			FrameRate:     "30fps",
	// 			IngestionType: "rtmp",
	// 			Resolution:    "1080p",
	// 		},
	// 		ContentDetails: &yt.LiveBroadcastContentDetails{
	// 			IsReusable: false,
	// 		},
	// 	},
	// ).Do()

	if err != nil {
		err = errors.Wrap(err, "could not create a new stream")
		return
	}

	return
}

func Cleanup() ([]string, error) {

	ctx, tokenSource := getTokenSource()

	svc, err := yt.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, errors.Wrap(err, "could not create Youtube API client")
	}

	streams, err := svc.LiveStreams.List(
		[]string{"snippet"},
	).Do()

	if err != nil {
		return nil, errors.Wrap(err, "could not create a new stream")
	}

	cleaned := make([]string, 0)
	r := regexp.MustCompile("\\[(.*)\\].*")
	for _, stream := range streams.Items {
		title := stream.Snippet.Title
		if !r.MatchString(title) {
			continue
		}
		sm := r.FindStringSubmatch(stream.Snippet.Title)
		if len(sm) < 2 {
			continue
		}

		date, err := dateparse.ParseAny(sm[1])
		if err != nil {
			continue
		}

		if date.Add(24 * 7 * time.Hour).After(time.Now()) {
			continue
		}

		err = svc.LiveStreams.Delete(stream.Id).Do()
		if err != nil {
			return cleaned, err
		}

		cleaned = append(cleaned, stream.Snippet.Title)
	}
	return cleaned, nil
}

func getTokenSource() (context.Context, CachedTokenSource) {
	ctx := context.Background()
	gts := NewGoogleTokenSource(ctx)
	tokenSource := NewCachedTokenSource(ctx, gts, gts.Config)
	return ctx, tokenSource
}
