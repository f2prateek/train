package httpevents

import (
	"net/http"

	"github.com/f2prateek/train"
	"github.com/segmentio/events"
)

func NewInterceptor() train.Interceptor {
	return NewInterceptorWith(events.DefaultLogger)
}

func NewInterceptorWith(logger *events.Logger) train.Interceptor {
	return &eventsInterceptor{logger}
}

type eventsInterceptor struct {
	logger *events.Logger
}

func (e *eventsInterceptor) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()
	r := makeRequest(req, "*")

	res, err := chain.Proceed(req)
	if res != nil {
		r.status = res.StatusCode
		r.log(e.logger, 1)
	}

	return res, err
}
