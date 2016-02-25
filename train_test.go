package train_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/f2prateek/train"
	"github.com/gohttp/response"
)

type recorderInterceptor struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (interceptor *recorderInterceptor) Intercept(chain train.Chain) (*http.Response, error) {
	interceptor.request = chain.Request()
	interceptor.response, interceptor.err = chain.Proceed(interceptor.request)
	return interceptor.response, interceptor.err
}

func TestInterceptor(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		response.OK(w, "Hello World!")
	}))
	defer ts.Close()
	recorder := &recorderInterceptor{}

	client := &http.Client{
		Transport: train.Train(recorder),
	}

	resp, err := client.Get(ts.URL)
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Hello World!\n", string(body))

	assert.Equal(t, ts.URL, recorder.request.URL.String())
	assert.Equal(t, 200, recorder.response.StatusCode)
	assert.Equal(t, nil, recorder.err)
}

func TestInterceptorCanShortCircuit(t *testing.T) {
	recorder1 := &recorderInterceptor{}
	short := func(c train.Chain) (*http.Response, error) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			response.OK(w, "Hello World!")
		}))
		defer ts.Close()

		return http.Get(ts.URL)
	}
	recorder2 := &recorderInterceptor{}

	client := &http.Client{
		Transport: train.Train(recorder1, train.InterceptorFunc(short), recorder2),
	}

	resp, err := client.Get("https://golang.org/")
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Hello World!\n", string(body))

	assert.Equal(t, "https://golang.org/", recorder1.request.URL.String())
	assert.Equal(t, 200, recorder1.response.StatusCode)
	assert.Equal(t, nil, recorder1.err)

	if recorder2.request != nil {
		t.Errorf("recorder2 should not have been invoked with a request.")
	}
	if recorder2.response != nil {
		t.Errorf("recorder2 should not have received a response.")
	}
	assert.Equal(t, nil, recorder2.err)
}
