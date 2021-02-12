package pkg

import (
	"regexp"
	"time"

	"github.com/araddon/dateparse"
	"github.com/mniak/ytlive/internal"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func CleanupKeys(since time.Time) ([]string, error) {

	config := internal.GetOAuthConfig()
	ctx, tokenSource := internal.GetTokenSource(config)

	svc, err := youtube.NewService(ctx, option.WithTokenSource(tokenSource))
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

		if date.Before(since) {
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
