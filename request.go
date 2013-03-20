package gadget

import (
	"encoding/json"
	"fmt"
	"github.com/redneckbeard/gadget/env"
	"io/ioutil"
	"net/http"
	"time"
)

type Request struct {
	*http.Request
	Path      string
	Payload   map[string]interface{}
	Method    string
	UrlParams map[string]string
}

func NewRequest(raw *http.Request) *Request {
	r := &Request{Request: raw, Path: raw.URL.Path[1:], Method: raw.Method}
	r.setPayload()
	return r
}

func (r *Request) ContentType() string {
	accept := r.Request.Header.Get("Accept")
	if accept != "" {
		return accept
	}
	return r.Request.Header.Get("Content-Type")
}

func (r *Request) setPayload() {
	payload := make(map[string]interface{})
	switch r.Request.Header.Get("Content-Type") {
	case "":
		err := r.ParseForm()
		if err != nil {
			return
		}
		for k, v := range r.Form {
			if len(v) == 1 {
				payload[k] = v[0]
			} else {
				payload[k] = v
			}
		}
	case "application/json":
		if r.Request.Body != nil {
			raw, err := ioutil.ReadAll(r.Request.Body)
			if err != nil {
				return
			}
			err = json.Unmarshal(raw, payload)
			if err != nil {
				return
			}
		}
	}
	r.Payload = payload
}

func (r *Request) Log(status, contentLength int) {
	raw := r.Request
	env.Log(fmt.Sprintf(`[%s] "%s %s %s" %d %d`, time.Now().Format(time.RFC822), r.Method, raw.URL.Path, raw.Proto, status, contentLength))
}
