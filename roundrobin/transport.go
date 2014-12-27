// Copyright (c) 2014-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by the MIT license.
// See LICENSE file for details.

// Most of the code here is taken from the Google OAuth2 client library
// at https://github.com/golang/oauth2,
// especially https://github.com/golang/oauth2/blob/master/transport.go.
package roundrobin

import (
	"io"
	"net/http"
	"sync"

	"github.com/olivere/balancers"
)

// Transport implements a http Transport for a round-robin http balancer.
type Transport struct {
	Base http.RoundTripper

	balancer *Balancer

	mu     sync.Mutex
	modReq map[*http.Request]*http.Request
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	conn, err := t.balancer.Get()
	if err != nil {
		return nil, err
	}

	rc := cloneRequest(r)
	if err := modifyRequest(rc, conn); err != nil {
		return nil, err
	}
	t.setModReq(r, rc)

	res, err := t.base().RoundTrip(rc)
	if err != nil {
		t.setModReq(r, nil)
		return nil, err
	}
	res.Body = &onEOFReader{
		rc: res.Body,
		fn: func() { t.setModReq(rc, nil) },
	}
	return res, nil
}

func (t *Transport) CancelRequest(r *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}
	if cr, ok := t.base().(canceler); ok {
		t.mu.Lock()
		modReq := t.modReq[r]
		delete(t.modReq, r)
		t.mu.Unlock()
		cr.CancelRequest(modReq)
	}
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

func modifyRequest(r *http.Request, conn balancers.Connection) error {
	url := conn.URL()
	if url.Scheme != "" {
		r.URL.Scheme = url.Scheme
	}
	if url.Host != "" {
		r.URL.Host = url.Host
	}
	if url.User != nil {
		r.URL.User = url.User
	}
	return nil
}

func cloneRequest(r *http.Request) *http.Request {
	rc := new(http.Request)
	*rc = *r
	rc.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		rc.Header[k] = append([]string(nil), s...)
	}
	return rc
}

func (t *Transport) setModReq(orig, mod *http.Request) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.modReq == nil {
		t.modReq = make(map[*http.Request]*http.Request)
	}
	if mod == nil {
		delete(t.modReq, orig)
	} else {
		t.modReq[orig] = mod
	}
}

type onEOFReader struct {
	rc io.ReadCloser
	fn func()
}

func (r *onEOFReader) Read(p []byte) (n int, err error) {
	n, err = r.rc.Read(p)
	if err == io.EOF {
		r.runFunc()
	}
	return
}

func (r *onEOFReader) Close() error {
	err := r.rc.Close()
	r.runFunc()
	return err
}

func (r *onEOFReader) runFunc() {
	if fn := r.fn; fn != nil {
		fn()
		r.fn = nil
	}
}

type errorTransport struct{ err error }

func (t errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, t.err
}
