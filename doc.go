// Copyright (c) 2014-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by the MIT license.
// See LICENSE file for details.

/*
Package balancers provides implementations of HTTP load-balancers.

It has two key interfaces: A Balancer is the implementation of a load-balancer
that chooses from a set of Connections.

You can e.g. use the balancer from the roundrobin package to rewrite
HTTP requests and use URLs from a given set of HTTP connections.

Suppose you have a cluster of two servers (on two different URLs) and you
want to load-balance between the two in a round-robin fashion, you can use
code like this:

	balancer, err := roundrobin.NewBalancerFromURL("https://server1.com", "https://server2.com")
	...
	// Get a HTTP client for the roundrobin balancer.
	client := balancer.Client()
	...
	client.Get("http://example.com/path1?foo=bar") // will rewrite URL to https://server1.com/path1?foo=bar
	client.Get("http://example.com/path1?foo=bar") // will rewrite URL to https://server2.com/path1?foo=bar
	client.Get("http://example.com/path1?foo=bar") // will rewrite URL to https://server1.com/path1?foo=bar
	client.Get("/path1?foo=bar") // will rewrite URL to https://server2.com/path1?foo=bar
*/
package balancers
