// Copyright (c) 2014-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by the MIT license.
// See LICENSE file for details.
package balancers

import (
	"net/http"
)

// NewClient returns a http Client that applies a certain scheduling algorithm
// (like round-robin) to load balance between several HTTP servers.
func NewClient(b Balancer) *http.Client {
	return &http.Client{
		Transport: &Transport{balancer: b},
	}
}
