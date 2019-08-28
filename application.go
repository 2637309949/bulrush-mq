// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mq

import (
	"time"

	"github.com/thoas/go-funk"
)

const (
	// INIT status
	INIT = "INIT"
	// PROCESSING status
	PROCESSING = "PROCESSING"
	// SUCCEED status
	SUCCEED = "SUCCEED"
	// FAILED status
	FAILED = "FAILED"
)

type (
	// MQ defined message queue struct
	MQ struct {
		Model      Model
		Exector    []Exector
		TypeTactic []TypeTactic
		Interval   []chan bool
	}
	// Message defined message entity struct
	Message struct {
		ID        int
		Type      string
		Body      map[string]interface{}
		Status    string
		CreatedAt time.Time
	}
	// Exector defined loop handler
	Exector struct {
		Type    string
		Handler func(*Message) error
	}
	// Tactic defined interval type
	Tactic struct {
		// Interval every interval do work
		Interval int
		// AsyncCount defined all count in aount
		AsyncCount uint
	}
	// TypeTactic defined interval type
	TypeTactic struct {
		Type   string
		Tactic Tactic
	}
)

// New defined obtain a MQ
func New() *MQ {
	mq := &MQ{}
	mq.TypeTactic = append(mq.TypeTactic, DEFAULTTYPETACTIC)
	mq.Model = &Memo{}
	mq.loop()
	return mq
}

// SetModel set mq model
func (mq *MQ) SetModel(model Model) *MQ {
	mq.Model = model
	return mq
}

// AddTactics add Tactics to system
func (mq *MQ) AddTactics(tp string, tac Tactic) *MQ {
	target := funk.Find(mq.TypeTactic, func(ttc TypeTactic) bool {
		return ttc.Type == tp
	})
	if target != nil {
		rushLogger.Info("rewrite Tactic strategy %v", target)
		typeOne := target.(TypeTactic)
		typeOne.Tactic = tac
	} else {
		if tac.AsyncCount == 0 {
			tac.AsyncCount = 1
		}
		mq.TypeTactic = append(mq.TypeTactic, TypeTactic{
			Type:   tp,
			Tactic: tac,
		})
	}
	go mq.loop()
	return mq
}

// stopTactic defined stop all tick
func (mq *MQ) stopTactic() *MQ {
	funk.ForEach(mq.Interval, func(timer chan bool) {
		timer <- true
		close(timer)
	})
	mq.Interval = make([]chan bool, 0)
	return mq
}

// stopTactic defined start all tick
func (mq *MQ) startTactic() *MQ {
	funk.ForEach(mq.TypeTactic, func(tac TypeTactic) {
		timer := setInterval(func() {
			AsyncCount := tac.Tactic.AsyncCount
			ttype := tac.Type
			var exector []Exector
			if ttype == "" {
				exector = funk.Filter(mq.Exector, func(exe Exector) bool {
					return funk.Find(mq.TypeTactic, func(ttc TypeTactic) bool {
						return ttc.Type == exe.Type
					}) == nil
				}).([]Exector)
			} else {
				exector = funk.Filter(mq.Exector, func(exe Exector) bool {
					return exe.Type == ttype
				}).([]Exector)
			}
			funk.ForEach(exector, func(exec Exector) {
				handler := exec.Handler
				handlerType := exec.Type
				pTaskCount := mq.Model.Count(handlerType, PROCESSING)
				message := mq.Model.Find(handlerType, INIT)
				if message != nil && pTaskCount < AsyncCount {
					err := mq.Model.Update(message, PROCESSING)
					if err != nil {
						mq.Model.Update(message, FAILED)
					} else {
						err := handler(message)
						if err != nil {
							mq.Model.Update(message, FAILED)
						} else {
							mq.Model.Update(message, SUCCEED)
						}
					}
				}
			})
		}, time.Duration(tac.Tactic.Interval)*time.Second)
		mq.Interval = append(mq.Interval, timer)
	})
	return mq
}

// loop events
func (mq *MQ) loop() *MQ {
	mq.stopTactic()
	mq.startTactic()
	return mq
}

// Push events
func (mq *MQ) Push(mess Message) {
	mess.CreatedAt = time.Now()
	mess.Status = INIT
	mq.Model.Save(mess)
}

// Register event handler
func (mq *MQ) Register(tp string, handler func(*Message) error) *MQ {
	mq.Exector = append(mq.Exector, Exector{Type: tp, Handler: handler})
	return mq
}

// Plugin defined Mq Plugin
func (mq *MQ) Plugin() *MQ {
	return mq
}
