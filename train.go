package train

import "net/http"

type Chain interface {
	Request() *http.Request
	Proceed(*http.Request) (*http.Response, error)
}

type Interceptor interface {
	Intercept(Chain) (*http.Response, error)
}

type InterceptorFunc func(Chain) (*http.Response, error)

func (f InterceptorFunc) Intercept(c Chain) (*http.Response, error) {
	return f(c)
}

func Transport(interceptors ...Interceptor) http.RoundTripper {
	return &interceptorRoundTripper{
		interceptors: append([]Interceptor{}, interceptors...),
	}
}

type interceptorRoundTripper struct {
	interceptors []Interceptor
}

func (i *interceptorRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	chain := &interceptorChain{
		index:        0,
		request:      req,
		interceptors: i.interceptors,
	}
	return chain.Proceed(req)
}

type interceptorChain struct {
	index        int
	request      *http.Request
	interceptors []Interceptor
}

func (c *interceptorChain) Request() *http.Request {
	return c.request
}

func (c *interceptorChain) Proceed(req *http.Request) (*http.Response, error) {
	// If there's another interceptor in the chain, call that.
	if c.index < len(c.interceptors) {
		chain := &interceptorChain{
			index:        c.index + 1,
			request:      req,
			interceptors: c.interceptors,
		}
		interceptor := c.interceptors[c.index]
		return interceptor.Intercept(chain)
	}

	// No more interceptors. Do HTTP.
	return http.DefaultTransport.RoundTrip(req)
}
