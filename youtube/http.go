package youtube

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ernesto-jimenez/httplogger"
	"github.com/spf13/viper"
)

func addLoggingTransportIfNeeded(client *http.Client) *http.Client {
	if viper.GetBool("Debug") {
		client.Transport = httplogger.NewLoggedTransport(http.DefaultTransport, newLogger())
	}
	return client
}

type httpLogger struct {
	out *log.Logger
	err *log.Logger
}

func newLogger() *httpLogger {
	const prefix = "[Http] "
	return &httpLogger{
		out: log.New(os.Stdout, prefix, log.LstdFlags),
		err: log.New(os.Stderr, prefix, log.LstdFlags),
	}
}

func (l *httpLogger) LogRequest(req *http.Request) {
	l.out.Printf(
		"Request %s %s",
		req.Method,
		req.URL.String(),
	)
}

func (l *httpLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	duration /= time.Millisecond
	if err != nil {
		l.err.Println(err)
		return
	}

	l.out.Printf(
		"Response method=%s status=%d durationMs=%d %s",
		req.Method,
		res.StatusCode,
		duration,
		req.URL.String(),
	)
}
