// Copyright (c) 2014-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by the MIT license.
// See LICENSE file for details.
package balancers

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Connection is a single connection to a host. It is defined by a URL.
// It also maintains state in the form that a connection can be broken.
// TODO(oe) Not sure if this abstraction is necessary.
type Connection interface {
	// URL to the host.
	URL() *url.URL
	// IsBroken must return true if the connection to URL is currently not available.
	IsBroken() bool
}

// HttpConnection is a HTTP connection to a host.
// It implements the Connection interface and can be used by balancer
// implementations.
type HttpConnection struct {
	sync.Mutex
	url               *url.URL
	broken            bool
	heartbeatDuration time.Duration
	heartbeatStop     chan bool
}

// NewHttpConnection creates a new HTTP connection to the given URL.
func NewHttpConnection(url *url.URL) *HttpConnection {
	c := &HttpConnection{
		url:               url,
		heartbeatDuration: DefaultHeartbeatDuration,
		heartbeatStop:     make(chan bool),
	}
	c.checkBroken()
	go c.heartbeat()
	return c
}

// Close this connection.
func (c *HttpConnection) Close() error {
	c.Lock()
	defer c.Unlock()
	c.heartbeatStop <- true // wait for heartbeat ticker to stop
	c.broken = false
	return nil
}

// HeartbeatDuration sets the duration in which the connection is checked.
func (c *HttpConnection) HeartbeatDuration(d time.Duration) *HttpConnection {
	c.Lock()
	defer c.Unlock()
	c.heartbeatStop <- true // wait for heartbeat ticker to stop
	c.broken = false
	c.heartbeatDuration = d
	go c.heartbeat()
	return c
}

// heartbeat periodically checks if the connection is broken.
func (c *HttpConnection) heartbeat() {
	ticker := time.NewTicker(c.heartbeatDuration)
	for {
		select {
		case <-ticker.C:
			c.checkBroken()
		case <-c.heartbeatStop:
			return
		}
	}
}

// checkBroken checks if the HTTP connection is alive.
func (c *HttpConnection) checkBroken() {
	c.Lock()
	defer c.Unlock()

	// TODO(oe) Can we use HEAD?
	req, err := http.NewRequest("GET", c.url.String(), nil)
	if err != nil {
		c.broken = true
		return
	}
	// Add UA to heartbeat requests.
	req.Header.Add("User-Agent", UserAgent)

	// Use a standard HTTP client with a timeout of 5 seconds.
	cl := &http.Client{Timeout: 5 * time.Second}
	res, err := cl.Do(req)
	if err == nil {
		defer res.Body.Close()
		if res.StatusCode == http.StatusOK {
			c.broken = false
		} else {
			c.broken = true
		}
	} else {
		c.broken = true
	}
}

// URL returns the URL of the HTTP connection.
func (c *HttpConnection) URL() *url.URL {
	return c.url
}

// IsBroken returns true if the HTTP connection is currently broken.
func (c *HttpConnection) IsBroken() bool {
	return c.broken
}
