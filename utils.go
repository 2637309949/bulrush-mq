// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mq

import "time"

func setInterval(what func(), delay time.Duration) chan bool {
	ticker := time.NewTicker(delay)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				go what()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}
