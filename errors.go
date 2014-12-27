// Copyright (c) 2014-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by the MIT license.
// See LICENSE file for details.
package balancers

import (
	"errors"
)

// ErrNoConn must be returned when a Balancer does not find a (non-broken) connection.
var ErrNoConn = errors.New("no connection")
