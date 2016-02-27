package train

import "net/http"

type Chain interface {
	// Request returns the `http.Request` for this chain.
	Request() *http.Request
	// Proceed returns the chain with a given request and returns the result.
	Proceed(*http.Request) (*http.Response, error)
}

// Observes, modifies, and potentially short-circuits requests going out and the corresponding
// requests coming back in. Typically interceptors will be used to add, remove, or transform headers
// on the request or response. Interceptors must return either a response or an error.
type Interceptor interface {
	// Intercept the chain and return a result.
	Intercept(Chain) (*http.Response, error)
}

// The InterceptorFunc type is an adapter to allow the use of ordinary functions as interceptors.
// If f is a function with the appropriate signature, InterceptorFunc(f) is a Interceptor that calls f.
type InterceptorFunc func(Chain) (*http.Response, error)

// Intercept calls f(c).
func (f InterceptorFunc) Intercept(c Chain) (*http.Response, error) {
	return f(c)
}

// Return a new `http.RoundTripper` with the given interceptors. Interceptors will be called in the order
// they are provided.
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
