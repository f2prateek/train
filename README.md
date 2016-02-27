# train

[![GoDoc](https://godoc.org/github.com/f2prateek/train?status.svg)](https://godoc.org/github.com/f2prateek/train)

Chainable HTTP client middleware. Borrowed from [OkHttp](https://github.com/square/okhttp/wiki/Interceptors).

Train can be installed transparently on any `http.Client`.

Here's the logging interceptor in action.

```go
import "github.com/f2prateek/train"
import "github.com/f2prateek/train/log"
import "github.com/f2prateek/train/statsd"

transport := train.Transport(log.New(os.Stdout, log.Body))

client := &http.Client{
  Transport: transport,
}
client.Get("https://golang.org")
```

The interceptor will transparently log all requests and responses being made by the client.

```
GET / HTTP/1.1
Host: 127.0.0.1:64598

HTTP/1.1 200 OK
Content-Length: 13
Content-Type: text/plain; charset=utf-8
Date: Thu, 25 Feb 2016 09:49:28 GMT

Hello World!
```
