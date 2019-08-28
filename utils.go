// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mq

import "time"

type sortByMsAt []Message

func (a sortByMsAt) Len() int           { return len(a) }
func (a sortByMsAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByMsAt) Less(i, j int) bool { return a[i].CreatedAt.Unix() < a[j].CreatedAt.Unix() }

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

// DEFAULTTYPETACTIC defined default Tactic
var DEFAULTTYPETACTIC = TypeTactic{
	Tactic: Tactic{
		Interval:   3,
		AsyncCount: 1,
	},
}
