package train_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/f2prateek/train"
	"github.com/gohttp/response"
	"github.com/stretchr/testify/mock"
)

// Interceptor that simply calls `chain.Proceed(chain.Request())`.
var fallThrough = train.InterceptorFunc(func(chain train.Chain) (*http.Response, error) {
	return chain.Proceed(chain.Request())
})

func TestFallThroughInterceptor(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		response.OK(w, "Hello World!")
	}))
	defer ts.Close()

	client := &http.Client{
		Transport: train.Train(fallThrough),
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
}

type MockInterceptor struct {
	mock.Mock
}

func (m *MockInterceptor) Intercept(chain train.Chain) (*http.Response, error) {
	args := m.Called(chain)
	return args.Get(0).(*http.Response), args.Error(1)
}

func NewMockInterceptor() *MockInterceptor {
	return new(MockInterceptor)
}

func TestInterceptorCanShortCircuit(t *testing.T) {
	m1 := NewMockInterceptor()
	{
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			response.OK(w, "Hello World!")
		}))
		defer ts.Close()
		resp, err := http.Get(ts.URL)
		m1.On("Intercept", mock.AnythingOfType("*train.interceptorChain")).Return(resp, err)
	}
	m2 := NewMockInterceptor()

	client := &http.Client{
		Transport: train.Train(fallThrough, m1, m2),
	}

	resp, err := client.Get("https://golang.org/")

	// Assert that the application sees our "shorted" response.
	assert.Equal(t, nil, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Hello World!\n", string(body))

	// Assert our mocks.
	m1.AssertExpectations(t)
	m2.AssertExpectations(t)
}
