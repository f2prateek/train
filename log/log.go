package log

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/f2prateek/train"
)

func New(out io.Writer) train.Interceptor {
	return &loggingInterceptor{out}
}

type loggingInterceptor struct {
	out io.Writer
}

func (interceptor *loggingInterceptor) Intercept(chain train.Chain) (*http.Response, error) {
	req := chain.Request()

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(interceptor.out, "%s", requestDump)

	resp, err := chain.Proceed(req)

	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(interceptor.out, "%s", responseDump)

	return resp, err
}
