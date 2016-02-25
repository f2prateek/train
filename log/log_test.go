package log_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/f2prateek/train"
	"github.com/f2prateek/train/log"
	"github.com/gohttp/response"
)

func ExampleNew() {
	var buf bytes.Buffer
	client := &http.Client{
		Transport: train.Train(log.New(&buf)),
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		response.OK(w, "Hello World!")
	}))
	defer ts.Close()

	client.Get(ts.URL)

	fmt.Println(buf.String())
}
