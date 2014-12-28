# Balancers

Balancers provides implementations of HTTP load-balancers.

[![Build Status](https://travis-ci.org/olivere/balancers.svg?branch=master)](https://travis-ci.org/olivere/balancers)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/olivere/balancers)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/olivere/balancers/master/LICENSE)

## What does it do?

Balancers gives you a `http.Client` from [net/http](http://golang.org/pkg/net/http)
that rewrites your requests' scheme, host, and userinfo according to the
rules of a balancer. A balancer is simply an algorithm to pick the host for
the next request a `http.Client`.

## How does it work?

Suppose you have a cluster of two servers (on two different URLs) and you
want to load balance between them. A very simple implementation can be done
with the [round-robin scheduling algorithm](http://en.wikipedia.org/wiki/Round-robin_scheduling).
Round-robin iterates through the list of available hosts and restarts
at the first when the end is reached. Here's some code that illustrates that:

```go
// Get a balancer that performs round-robin scheduling between two servers.
balancer, err := roundrobin.NewBalancerFromURL("https://server1.com", "https://server2.com")

// Get a HTTP client based on that balancer.
client := balancers.NewClient(balancer)

// Now request some data. The scheme, host, and user info will be rewritten
// by the balancer; you'll never get data from http://example.com, only data
// from http://server1.com or http://server2.com.
client.Get("http://example.com/path1?foo=bar") // rewritten to https://server1.com/path1?foo=bar
client.Get("http://example.com/path1?foo=bar") // rewritten to https://server2.com/path1?foo=bar
client.Get("http://example.com/path1?foo=bar") // rewritten to https://server1.com/path1?foo=bar
client.Get("/path1?foo=bar")                   // rewritten to https://server2.com/path1?foo=bar
```

## Status

The current state of Balancers is a proof-of-concept.
It didn't touch production systems yet.

## Credits

Thanks a lot for the great folks working on [Go](http://www.golang.org/).

## LICENSE

MIT-LICENSE. See [LICENSE](http://olivere.mit-license.org/)
or the LICENSE file provided in the repository for details.
