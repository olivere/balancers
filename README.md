# Balancers

Balancers provides implementations of HTTP load-balancers.

[![Build Status](https://travis-ci.org/olivere/balancers.svg?branch=master)](https://travis-ci.org/olivere/balancers)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/olivere/balancers)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/olivere/balancers/master/LICENSE)

## What does it do?

Balancer gives you a `http.Client` from [net/http](http://golang.org/pkg/net/http)
that rewrites your requests according to its rules.

## How does it work?

Let's start by an example. Suppose you have a cluster of two servers
(on two different URLs) and you want to load balance between.
If you want to do this in a round-robin fashion, you can use code like this:

```go
    balancer, err := roundrobin.NewBalancerFromURL("https://server1.com", "https://server2.com")
    ...
    client := balancer.Client() // Get a HTTP client for the round-robin balancer
    ...
    client.Get("http://example.com/path1?foo=bar") // will rewrite URL to https://server1.com/path1?foo=bar
    client.Get("http://example.com/path1?foo=bar") // will rewrite URL to https://server2.com/path1?foo=bar
    client.Get("http://example.com/path1?foo=bar") // will rewrite URL to https://server1.com/path1?foo=bar
    client.Get("/path1?foo=bar") // will rewrite URL to https://server2.com/path1?foo=bar
```

## Status

The current state of Balancers is a proof-of-concept.
It didn't touch production systems yet.

## Credits

Thanks a lot for the great folks working on [Go](http://www.golang.org/).

## LICENSE

MIT-LICENSE. See [LICENSE](http://olivere.mit-license.org/)
or the LICENSE file provided in the repository for details.
