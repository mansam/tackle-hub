package addon

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

//
// Client provides a REST client.
// TODO: REPLACE WITH CLIENT IN controller/rest (TBD).
type Client struct {
	// baseURL for the nub.
	baseURL string
	// http client.
	http *http.Client
}

//
// Get a resource.
func (r *Client) Get(path string, object interface{}) (err error) {
	request := &http.Request{
		Method: http.MethodGet,
		URL: r.join(path),
	}
	reply, err := r.http.Do(request)
	if err != nil {
		return
	}
	defer func() {
		_ = reply.Body.Close()
	}()
	status := reply.StatusCode
	if status != http.StatusOK {
		err = errors.New(http.StatusText(status))
		return
	}
	body, err := io.ReadAll(reply.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, object)
	return
}

//
// Post a resource.
func (r *Client) Post(path string, object interface{}) (err error) {
	err = r.post(http.MethodPost, path, object)
	return
}

//
// Put a resource.
func (r *Client) Put(path string, object interface{}) (err error) {
	err = r.post(http.MethodPut, path, object)
	return
}

//
// Delete a resource.
func (r *Client) Delete(path string) (err error) {
	request := &http.Request{
		Method: http.MethodDelete,
		URL: r.join(path),
	}
	reply, err := r.http.Do(request)
	if err != nil {
		return
	}
	defer func() {
		_ = reply.Body.Close()
	}()
	status := reply.StatusCode
	if status != http.StatusNoContent {
		err = errors.New(http.StatusText(status))
		return
	}
	return
}

//
// Post a resource.
func (r *Client) post(method string, path string, object interface{}) (err error) {
	bfr, err := json.Marshal(object)
	if err != nil {
		return
	}
	reader := bytes.NewReader(bfr)
	request := &http.Request{
		Method: method,
		Body: ioutil.NopCloser(reader),
		URL: r.join(path),
	}
	reply, err := r.http.Do(request)
	if err != nil {
		return
	}
	status := reply.StatusCode
	switch status {
	case http.StatusOK,
		http.StatusCreated:
		var content []byte
		content, err = ioutil.ReadAll(reply.Body)
		if err != nil {
			return
		}
		err = json.Unmarshal(content, object)
		if err != nil {
			return
		}
	case http.StatusNoContent:
		return
	default:
		err = errors.New(http.StatusText(status))
		return
	}

	return
}

func (r *Client) join(path string) (parsedURL *url.URL) {
	parsedURL, _ = url.Parse(r.baseURL)
	parsedURL.Path = path
	return
}
