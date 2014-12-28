// Copyright (c) 2014-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by the MIT license.
// See LICENSE file for details.
package balancers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCloseRequest(t *testing.T) {
	orig, _ := http.NewRequest("GET", "https://user:pwd@localhost:12345/path?query=1#hash", nil)
	dup := cloneRequest(orig)
	if dup.URL.User != orig.URL.User {
		t.Errorf("expected userinfo %v; got: %v", orig.URL.User, dup.URL.User)
	}
	if dup.URL.Scheme != "https" {
		t.Errorf("expected scheme %q; got: %q", "https", dup.URL.Scheme)
	}
	if dup.URL.Host != "localhost:12345" {
		t.Errorf("expected host %q; got: %q", "localhost:12345", dup.URL.Host)
	}
	if dup.URL.Path != "/path" {
		t.Errorf("expected path %q; got: %q", "/path", dup.URL.Path)
	}
	if dup.URL.RawQuery != "query=1" {
		t.Errorf("expected raw query %q; got: %q", "query=1", dup.URL.RawQuery)
	}
	if dup.URL.Fragment != "hash" {
		t.Errorf("expected fragment %q; got: %q", "hash", dup.URL.Fragment)
	}
}

func TestModifyRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tests := []struct {
		Req      string
		ConnURL  string
		Expected error
	}{
		{
			"http://localhost:12345/path?query=1#hash",
			server.URL,
			nil,
		},
	}

	for _, test := range tests {
		orig, err := url.Parse(test.Req)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("GET", test.Req, nil)
		if err != nil {
			t.Fatal(err)
		}

		url, err := url.Parse(test.ConnURL)
		if err != nil {
			t.Fatal(err)
		}
		conn := NewHttpConnection(url)

		err = modifyRequest(req, conn)
		if err != test.Expected {
			t.Errorf("expected err = %v; got: %v", test.Expected, err)
		} else {
			if url.Scheme != "" && req.URL.Scheme != url.Scheme {
				t.Errorf("expected scheme %q; got: %q", url.Scheme, req.URL.Scheme)
			}
			if url.Host != "" && req.URL.Host != url.Host {
				t.Errorf("expected host %q; got: %q", url.Scheme, req.URL.Scheme)
			}
			if url.User != nil && req.URL.User != url.User {
				t.Errorf("expected userinfo %v; got: %v", url.Scheme, req.URL.Scheme)
			}
			if req.URL.Path != orig.Path {
				t.Errorf("expected path %q; got: %q", orig.Path, req.URL.Path)
			}
			if req.URL.RawQuery != orig.RawQuery {
				t.Errorf("expected raw query %q; got: %q", orig.RawQuery, req.URL.RawQuery)
			}
			if req.URL.Fragment != orig.Fragment {
				t.Errorf("expected fragment %q; got: %q", orig.Fragment, req.URL.Fragment)
			}
		}
	}
}
