// Copyright (c) 2014-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by the MIT license.
// See LICENSE file for details.
package roundrobin

import (
	"net/url"
	"sync"

	"github.com/olivere/balancers"
)

// Balancer implements a round-robin balancer.
type Balancer struct {
	sync.Mutex // guards the following variables
	conns      []balancers.Connection
	idx        int // index into conns
}

// NewBalancer creates a new round-robin balancer. It can be initializes by
// a variable number of connections. To use plain URLs instead of
// connections, use NewBalancerFromURL.
func NewBalancer(conns ...balancers.Connection) (balancers.Balancer, error) {
	b := &Balancer{
		conns: make([]balancers.Connection, 0),
	}
	if len(conns) > 0 {
		b.conns = append(b.conns, conns...)
	}
	return b, nil
}

// NewBalancerFromURL creates a new round-robin balancer for the
// given list of URLs. It returns an error if any of the URLs is invalid.
func NewBalancerFromURL(urls ...string) (*Balancer, error) {
	b := &Balancer{
		conns: make([]balancers.Connection, 0),
	}
	for _, rawurl := range urls {
		if u, err := url.Parse(rawurl); err != nil {
			return nil, err
		} else {
			b.conns = append(b.conns, balancers.NewHttpConnection(u))
		}
	}
	return b, nil
}

// Get returns a connection from the balancer that can be used for the next request.
// ErrNoConn is returns when no connection is available.
func (b *Balancer) Get() (balancers.Connection, error) {
	b.Lock()
	defer b.Unlock()

	if len(b.conns) == 0 {
		return nil, balancers.ErrNoConn
	}

	var conn balancers.Connection
	for i := 0; i < len(b.conns); i++ {
		candidate := b.conns[b.idx]
		b.idx = (b.idx + 1) % len(b.conns)
		if !candidate.IsBroken() {
			conn = candidate
			break
		}
	}

	if conn == nil {
		return nil, balancers.ErrNoConn
	}
	return conn, nil
}

// Connections returns a list of all connections.
func (b *Balancer) Connections() []balancers.Connection {
	b.Lock()
	defer b.Unlock()
	conns := make([]balancers.Connection, len(b.conns))
	for i, c := range b.conns {
		if oc, ok := c.(*balancers.HttpConnection); ok {
			// Make a clone
			cr := new(balancers.HttpConnection)
			*cr = *oc
			conns[i] = cr
		}
	}
	return conns
}
