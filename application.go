// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mq

import (
	"sort"
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
	// Model defined model store
	Model interface {
		Save(Message)
		Find(string, string) []Message
		Count(string, string) uint
		Update(Message, string) error
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
		Handler func(Message) error
	}
	// Tactic defined interval type
	Tactic struct {
		Interval int
		CTCount  uint
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
	mq.Model = &MemoModel{}
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
	typeTac := funk.Find(mq.TypeTactic, func(tc TypeTactic) bool {
		return tc.Type == tp
	})
	if typeTac != nil {
		rushLogger.Info("rewrite Tactic strategy %v", typeTac)
		typeOne := typeTac.(TypeTactic)
		typeOne.Tactic = tac
	} else {
		mq.TypeTactic = append(mq.TypeTactic, TypeTactic{
			Type:   tp,
			Tactic: tac,
		})
	}
	return mq
}

func (mq *MQ) stopTactic() *MQ {
	funk.ForEach(mq.Interval, func(inv chan bool) {
		inv <- true
	})
	return mq
}

func (mq *MQ) startTactic() *MQ {
	funk.ForEach(mq.TypeTactic, func(tac TypeTactic) {
		setInterval(func() {
			ctCount := tac.Tactic.CTCount
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
				iTask := mq.Model.Find(handlerType, INIT)
				sort.Sort(sortByMsAt(iTask))
				if len(iTask) >= 1 {
					task := iTask[0]
					if pTaskCount < ctCount {
						err := mq.Model.Update(task, PROCESSING)
						if err != nil {
							mq.Model.Update(task, FAILED)
						} else {
							err := handler(task)
							if err != nil {
								mq.Model.Update(task, FAILED)
							} else {
								mq.Model.Update(task, SUCCEED)
							}
						}
					}
				}
			})
		}, time.Duration(tac.Tactic.Interval)*time.Second)
	})
	return mq
}

// loop events
func (mq *MQ) loop() *MQ {
	// 1. stop all tactic
	mq.stopTactic()
	// 2. restart all tactic
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
func (mq *MQ) Register(tp string, handler func(Message) error) *MQ {
	mq.Exector = append(mq.Exector, Exector{Type: tp, Handler: handler})
	return mq
}

// Plugin defined Mq Plugin
func (mq *MQ) Plugin() *MQ {
	return mq
}
