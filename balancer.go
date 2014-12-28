// Copyright (c) 2014-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by the MIT license.
// See LICENSE file for details.
package balancers

// Balancer holds a list of connections to hosts.
type Balancer interface {
	// Get returns a connection that can be used for the next request.
	Get() (Connection, error)

	// Connections is the list of available connections.
	Connections() []Connection
}
