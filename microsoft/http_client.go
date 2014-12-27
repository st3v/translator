package microsoft

import (
	"io"
	"net/http"

	"github.com/st3v/tracerr"
	"github.com/st3v/translator/microsoft/auth"
)

// The HTTPClient sends authenticated HTTP requests to Microsoft's Translation API
type HTTPClient interface {
	SendRequest(method, uri string, body io.Reader, contentType string) (*http.Response, error)
}

type httpClient struct {
	client        *http.Client
	authenticator auth.Authenticator
}

func newHTTPClient(authenticator auth.Authenticator) HTTPClient {
	return &httpClient{
		client:        &http.Client{},
		authenticator: authenticator,
	}
}

func (h *httpClient) SendRequest(method, uri string, body io.Reader, contentType string) (*http.Response, error) {
	request, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	request.Header.Add("Content-Type", contentType)

	err = h.authenticator.Authenticate(request)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	response, err := h.client.Do(request)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return response, nil
}
