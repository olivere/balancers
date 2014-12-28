// Copyright (c) 2014-2015 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by the MIT license.
// See LICENSE file for details.
package balancers

import (
	"runtime"
	"time"
)

const (
	// Version is the current version of this package.
	Version = "1.0.0"
)

var (
	// UserAgent is sent with all heartbeat requests.
	UserAgent = "balancers/" + Version + " (" + runtime.GOOS + "-" + runtime.GOARCH + ")"

	// DefaultHeartbeatDuration is the default time between heartbeat messages.
	DefaultHeartbeatDuration = 30 * time.Second
)
