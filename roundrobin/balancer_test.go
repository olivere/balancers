// Copyright (c) 2014-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by the MIT license.
// See LICENSE file for details.
package roundrobin

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/olivere/balancers"
)

func TestNewBalancer(t *testing.T) {
	url1, _ := url.Parse("http://127.0.0.1:12345")
	url2, _ := url.Parse("http://127.0.0.1:23456")

	balancer, err := NewBalancer(
		balancers.NewHttpConnection(url1),
		balancers.NewHttpConnection(url2))
	if err != nil {
		t.Fatal(err)
	}
	conns := balancer.Connections()
	if len(conns) != 2 {
		t.Errorf("expected %d connections; got: %v", 2, len(conns))
	}
	url := conns[0].URL()
	if url.String() != "http://127.0.0.1:12345" {
		t.Errorf("expected %q; got: %q", "http://127.0.0.1:12345", url.String())
	}
}

func TestBalancerErrNoConnWithoutConnections(t *testing.T) {
	balancer, err := NewBalancer()
	if err != nil {
		t.Fatal(err)
	}
	conns := balancer.Connections()
	if len(conns) != 0 {
		t.Errorf("expected %d connections; got: %v", 0, len(conns))
	}
	_, err = balancer.Get()
	if err != balancers.ErrNoConn {
		t.Fatalf("expected %v; got: %v", balancers.ErrNoConn, err)
	}
}

func TestBalancer(t *testing.T) {
	var visited []int

	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only count non-heartbeat requests
		if r.Header.Get("User-Agent") != balancers.UserAgent {
			visited = append(visited, 1)
		}
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only count non-heartbeat requests
		if r.Header.Get("User-Agent") != balancers.UserAgent {
			visited = append(visited, 2)
		}
	}))
	defer server2.Close()

	balancer, err := NewBalancerFromURL(server1.URL, server2.URL)
	if err != nil {
		t.Fatal(err)
	}

	client := balancers.NewClient(balancer)
	client.Get(server1.URL)
	client.Get(server1.URL)
	client.Get(server1.URL)

	if len(visited) != 3 {
		t.Fatalf("expected %d URLs to be visited; got: %d", 3, len(visited))
	}
	if visited[0] != 1 {
		t.Errorf("expected 1st URL to be %q", server1.URL)
	}
	if visited[1] != 2 {
		t.Errorf("expected 2nd URL to be %q", server2.URL)
	}
	if visited[2] != 1 {
		t.Errorf("expected 3rd URL to be %q", server1.URL)
	}
}

func TestBalancerWithBrokenConnections(t *testing.T) {
	var visited []int

	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only count non-heartbeat requests
		if r.Header.Get("User-Agent") != balancers.UserAgent {
			visited = append(visited, 1)
		}
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only count non-heartbeat requests
		if r.Header.Get("User-Agent") != balancers.UserAgent {
			visited = append(visited, 2)
		}
	}))
	defer server2.Close()

	balancer, err := NewBalancerFromURL(server1.URL, "http://localhost:12345", server2.URL, "http://localhost:12346")
	if err != nil {
		t.Fatal(err)
	}

	client := balancers.NewClient(balancer)
	client.Get(server1.URL)
	client.Get(server1.URL)
	client.Get(server1.URL)
	client.Get(server1.URL)
	client.Get(server1.URL)

	if len(visited) != 5 { // 5 requests
		t.Fatalf("expected %d URLs to be visited; got: %d", 5, len(visited))
	}
	if visited[0] != 1 {
		t.Errorf("expected 1st URL to be %q", server1.URL)
	}
	if visited[1] != 2 {
		t.Errorf("expected 2nd URL to be %q", server2.URL)
	}
	if visited[2] != 1 {
		t.Errorf("expected 3rd URL to be %q", server1.URL)
	}
	if visited[3] != 2 {
		t.Errorf("expected 4th URL to be %q", server2.URL)
	}
	if visited[4] != 1 {
		t.Errorf("expected 5th URL to be %q", server1.URL)
	}
}

func TestBalancerRewritesSchemeAndURLButNotPathOrQuery(t *testing.T) {
	var visited []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only count non-heartbeat requests
		if r.Header.Get("User-Agent") != balancers.UserAgent {
			visited = append(visited, r.URL.String())
		}
	}))
	defer server.Close()

	balancer, err := NewBalancerFromURL(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	client := balancers.NewClient(balancer)
	client.Get(server.URL + "/path?foo=bar&n=1")
	client.Get(server.URL + "/path?n=2")
	client.Get(server.URL + "/no/3")

	if len(visited) != 3 {
		t.Fatalf("expected %d URLs to be visited; got: %d", 3, len(visited))
	}
	if visited[0] != "/path?foo=bar&n=1" {
		t.Errorf("expected 1st URL to be %q; got: %q", "/path?foo=bar&n=1", visited[0])
	}
	if visited[1] != "/path?n=2" {
		t.Errorf("expected 2nd URL to be %q; got: %q", "/path?n=2", visited[1])
	}
	if visited[2] != "/no/3" {
		t.Errorf("expected 3rd URL to be %q; got: %q", "/no/3", visited[2])
	}
}
