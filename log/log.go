package log

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/f2prateek/train"
)

type Level uint8

const (
	None Level = iota
	Basic
	Body
)

func New(out io.Writer, level Level) train.Interceptor {
	return &loggingInterceptor{
		out:   out,
		level: level,
	}
}

type loggingInterceptor struct {
	out   io.Writer
	level Level
}

func (interceptor *loggingInterceptor) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()
	if interceptor.level == None {
		return chain.Proceed(req)
	}

	logBody := interceptor.level == Body

	requestDump, err := httputil.DumpRequestOut(req, logBody)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(interceptor.out, "%s", requestDump)

	resp, err := chain.Proceed(req)

	responseDump, err := httputil.DumpResponse(resp, logBody)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(interceptor.out, "%s", responseDump)

	return resp, err
}
