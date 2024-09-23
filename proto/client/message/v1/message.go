package v1

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/neo532/apitool/transport/http/xhttp/middleware"
)

func Demo() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(c context.Context, req, reply interface{}) (err error) {

			return handler(c, req, reply)
		}
	}
}

func RequestEncoder(c context.Context, contentType string, in interface{}) (body []byte, err error) {
	return json.Marshal(in)
}

func ResponseDecoder(c context.Context, res *http.Response, v interface{}) (body []byte, err error) {
	defer res.Body.Close()
	if body, err = io.ReadAll(res.Body); err != nil {
		return
	}
	err = json.Unmarshal(body, v)
	return
}

func ErrorDecoder(c context.Context, resp *http.Response) (err error) {
	if resp == nil {
		return errors.New("nil *http.Response")
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return
	}
	return errors.New(resp.Status)
}
