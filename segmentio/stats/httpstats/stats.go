package httpstats

import (
	"net/http"
	"time"

	"github.com/f2prateek/train"
	"github.com/segmentio/stats"
)

func NewInterceptor() train.Interceptor {
	return &statsInterceptor{stats.DefaultEngine}
}

func NewInterceptorWith(eng *stats.Engine) train.Interceptor {
	return &statsInterceptor{eng}
}

type statsInterceptor struct {
	eng *stats.Engine
}

func (s *statsInterceptor) Intercept(chain train.Chain) (*http.Response, error) {
	start := time.Now()
	req := chain.Request()

	if req.Body == nil {
		req.Body = &nullBody{}
	}

	req.Body = &requestBody{
		eng:  s.eng,
		req:  req,
		body: req.Body,
		op:   "write",
	}

	res, err := chain.Proceed(req)
	req.Body.Close() // safe guard, the transport should have done it already

	if err != nil {
		m := metrics{s.eng}
		m.observeError(req, "write")
	}

	if res != nil {
		res.Body = &responseBody{
			eng:   s.eng,
			res:   res,
			body:  res.Body,
			op:    "read",
			start: start,
		}
	}

	return res, err
}
