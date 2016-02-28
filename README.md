# train

Chainable HTTP client middleware. Borrowed from [OkHttp](https://github.com/square/okhttp/wiki/Interceptors).

Interceptors monitor outgoing requests, and incoming responses.
Here's an interceptor that logs all requests made and responses received by the client.

```go
import "github.com/f2prateek/train"

f := func(chain train.Chain) (*http.Response, error) {
  req := chain.Request()
  fmt.Println(httputil.DumpRequestOut(req, false))

  resp, err := chain.Proceed(req)
  fmt.Println(httputil.DumpRequestOut(req, false))

  return resp, err
})

transport := train.Transport(train.InterceptorFunc(f))
client := &http.Client{
  Transport: transport,
}

client.Get("https://golang.org")

// Output:
// GET / HTTP/1.1
// Host: 127.0.0.1:64598
//
// HTTP/1.1 200 OK
// Content-Length: 13
// Content-Type: text/plain; charset=utf-8
// Date: Thu, 25 Feb 2016 09:49:28 GMT
```

Train chains interceptors so that they can be plugged into any `http.Client`. For example, this chain will transparently:
* retry requests for temporary errors.
* log requests and responses.
* increment some stats based the requests and responses.

```go
import "github.com/f2prateek/log"
import "github.com/f2prateek/statsd"
import "github.com/f2prateek/hickson"
import "github.com/f2prateek/hickson/temporary"

errInterceptor := hickson.New(hickson.RetryMax(5, temporary.RetryErrors()))
logInterceptor := log.New(os.Stdout, log.Body)
statsInterceptor := statsd.New(statsdClient)

transport := train.Transport(errInterceptor, logInterceptor, statsInterceptor)
client := &http.Client{
  Transport: transport,
}

client.Get("https://golang.org")
```

Interceptors are consulted in the order they are provided. You'll need to decide what order you want your interceptors to be called in.

```go
// This chain will log and record all requests/responses (regardless of temporary errors).
transport := train.Transport(errInterceptor, logInterceptor, statsInterceptor)

// This chain will log and record requests/responses only after all retries for temporary errors have been exhausted.
transport := train.Transport(logInterceptor, statsInterceptor, errInterceptor)
```

Interceptors may modify requests and responses.

```go
func Intercept(chain train.Chain) (*http.Response, error) {
  req := chain.Request()
  req.Header.Add("User-Agent", "Train Example")

  resp, err := chain.Proceed(req)
  resp.Header.Add("Cache-Control", "max-age=60")

  return resp, err
})
```

Or interceptors can simply short circuit the chain.

```go
func Intercept(train.Chain) (*http.Response, error) {
  return cannedResponse, cannedError
})
```
