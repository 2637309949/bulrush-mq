// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mq

import (
	"github.com/2637309949/bulrush-addition/logger"
)

// RushLogger just for console log
var RushLogger = logger.CreateLogger(logger.SILLYLevel, nil,
	[]*logger.Transport{
		&logger.Transport{
			Level: logger.SILLYLevel,
		},
	},
)
