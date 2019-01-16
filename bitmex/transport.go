package bitmex

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/nntaoli-project/GoEx"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
)

type Transport struct {
	*httptransport.Runtime
	Key        string
	Secret     string
	timeOffset int64
}

func NewTransport(host, basePath, key, secret string, schemes []string) (t *Transport) {
	t = new(Transport)
	t.Key = key
	t.Secret = secret
	t.Runtime = httptransport.New(host, basePath, schemes)
	t.Runtime.Producers["application/x-www-form-urlencoded"] = runtime.TextProducer()
	return
}

func (t *Transport) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	if operation.Method == "POST" || operation.Method == "DELETE" || operation.Method == "PUT" {
		operation.ConsumesMediaTypes = []string{"application/x-www-form-urlencoded", "application/json"}
	}
	var fn runtime.ClientAuthInfoWriterFunc
	fn = func(req runtime.ClientRequest, formats strfmt.Registry) error {
		expires, sign := t.signature(req, operation, formats)
		req.SetTimeout(30 * time.Second)
		req.SetHeaderParam("api-key", t.Key)
		req.SetHeaderParam("api-nonce", expires)
		req.SetHeaderParam("api-signature", sign)
		return nil
	}
	operation.AuthInfo = fn
	return t.Runtime.Submit(operation)
}

func (t *Transport) signature(req runtime.ClientRequest, operation *runtime.ClientOperation, formats strfmt.Registry) (expires, sign string) {
	method := strings.ToUpper(req.GetMethod())
	path := t.BasePath + req.GetPath()
	query := req.GetQueryParams()
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	nonce := time.Now().UnixNano()
	expires = strconv.FormatInt(nonce/1000000-t.timeOffset, 10)[:13]
	var data string
	switch method {
	case "GET":
		data = ""
	case "POST", "DELETE", "PUT":
		rb := NewRequestBuffer()
		err := operation.Params.WriteToRequest(rb, formats)
		if err != nil {
			log.Error("WriteToRequest error:", err.Error())
		}
		data = rb.GetFormParams().Encode()
	default:
		log.Info("unsupport method:", method)
	}
	content := method + path + expires + data
	sign, _ = GetParamHmacSHA256Sign(t.Secret, content)
	return
}

type RequestBuffer struct {
	pathPattern string
	method      string

	pathParams map[string]string
	header     http.Header
	query      url.Values
	formFields url.Values
	fileFields map[string][]runtime.NamedReadCloser
	payload    interface{}
	timeout    time.Duration
	buf        *bytes.Buffer
}

func NewRequestBuffer() (r *RequestBuffer) {
	r = new(RequestBuffer)
	return
}

func (r *RequestBuffer) SetHeaderParam(name string, values ...string) error {
	if r.header == nil {
		r.header = make(http.Header)
	}
	r.header[http.CanonicalHeaderKey(name)] = values
	return nil
}

func (r *RequestBuffer) SetQueryParam(name string, values ...string) error {
	if r.query == nil {
		r.query = make(url.Values)
	}
	r.query[name] = values
	return nil
}

func (r *RequestBuffer) SetFormParam(name string, values ...string) error {
	if r.formFields == nil {
		r.formFields = make(url.Values)
	}
	r.formFields[name] = values
	return nil
}

func (r *RequestBuffer) SetPathParam(name string, value string) error {
	if r.pathParams == nil {
		r.pathParams = make(map[string]string)
	}

	r.pathParams[name] = value
	return nil
}

func (r *RequestBuffer) GetFormParams() url.Values {
	var result = make(url.Values)
	for key, value := range r.formFields {
		result[key] = append([]string{}, value...)
	}
	return result
}

func (r *RequestBuffer) GetQueryParams() url.Values {
	var result = make(url.Values)
	for key, value := range r.query {
		result[key] = append([]string{}, value...)
	}
	return result
}

func (r *RequestBuffer) SetFileParam(name string, files ...runtime.NamedReadCloser) error {
	for _, file := range files {
		if actualFile, ok := file.(*os.File); ok {
			fi, err := os.Stat(actualFile.Name())
			if err != nil {
				return err
			}
			if fi.IsDir() {
				return fmt.Errorf("%q is a directory, only files are supported", file.Name())
			}
		}
	}

	if r.fileFields == nil {
		r.fileFields = make(map[string][]runtime.NamedReadCloser)
	}
	if r.formFields == nil {
		r.formFields = make(url.Values)
	}

	r.fileFields[name] = files
	return nil
}

func (r *RequestBuffer) SetBodyParam(payload interface{}) error {
	r.payload = payload
	return nil
}

func (r *RequestBuffer) SetTimeout(timeout time.Duration) error {
	r.timeout = timeout
	return nil
}
func (r *RequestBuffer) GetMethod() string {
	return r.method
}

func (r *RequestBuffer) GetPath() string {
	path := r.pathPattern
	for k, v := range r.pathParams {
		path = strings.Replace(path, "{"+k+"}", v, -1)
	}
	return path
}

func (r *RequestBuffer) GetBody() []byte {
	if r.buf == nil {
		return nil
	}
	return r.buf.Bytes()
}

func (r *RequestBuffer) GetBodyParam() interface{} {
	return r.payload
}

func (r *RequestBuffer) GetFileParam() map[string][]runtime.NamedReadCloser {
	return r.fileFields
}
