// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mq

import (
	"sort"
	"time"

	"github.com/2637309949/bulrush"
	"github.com/thoas/go-funk"
)

type (
	// MQ defined message queue struct
	MQ struct {
		bulrush.PNBase
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

type sortByMsAt []Message

func (a sortByMsAt) Len() int           { return len(a) }
func (a sortByMsAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByMsAt) Less(i, j int) bool { return a[i].CreatedAt.Unix() < a[j].CreatedAt.Unix() }

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

// New defined obtain a MQ
func New() *MQ {
	mq := &MQ{}
	mq.TypeTactic = append(mq.TypeTactic, TypeTactic{
		Tactic: Tactic{
			Interval: 3,
			CTCount:  1,
		},
	})
	return mq
}

// SetModel set mq model
func (mq *MQ) SetModel(m Model) {
	mq.Model = m
}

// AddTactics add Tactics to system
func (mq *MQ) AddTactics(tp string, tac Tactic) *MQ {
	one := funk.Find(mq.TypeTactic, func(tc TypeTactic) bool {
		return tc.Type == tp
	})
	if one != nil {
		typeOne := one.(TypeTactic)
		typeOne.Tactic = tac
	} else {
		mq.TypeTactic = append(mq.TypeTactic, TypeTactic{
			Type:   tp,
			Tactic: tac,
		})
	}
	return mq
}

// loop events
func (mq *MQ) loop() {
	// stop all interval
	RushLogger.Info("stop all interval")
	for _, inv := range mq.Interval {
		inv <- true
	}
	RushLogger.Info("start all interval")
	// start all interval
	for _, tac := range mq.TypeTactic {
		RushLogger.Info("start tac %v", tac)
		interval := tac.Tactic.Interval
		ctCount := tac.Tactic.CTCount
		ttype := tac.Type
		setInterval(func() {
			var exector []Exector
			if ttype == "" {
				exector = funk.Filter(mq.Exector, func(exe Exector) bool {
					return exe.Type == ttype
				}).([]Exector)
			} else {
				exector = funk.Filter(mq.Exector, func(exe Exector) bool {
					return funk.Filter(mq.TypeTactic, func(ttc TypeTactic) bool {
						return ttc.Type == exe.Type
					}) != nil
				}).([]Exector)
			}
			for _, exec := range exector {
				handler := exec.Handler
				handlerType := exec.Type
				pTaskCount := mq.Model.Count(handlerType, PROCESSING)
				iTask := mq.Model.Find(handlerType, INIT)
				sort.Sort(sortByMsAt(iTask))
				task := iTask[0]
				if pTaskCount < ctCount && task.Type != "" {
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
		}, time.Duration(interval)*time.Second)
	}
}

// Push events
func (mq *MQ) Push(mess Message) {
	mess.CreatedAt = time.Now()
	mess.Status = INIT
	mq.Model.Save(mess)
}

// Register event handler
func (mq *MQ) Register(tp string, handler func(Message) error) {
	mq.Exector = append(mq.Exector, Exector{Type: tp, Handler: handler})
}

// Plugin defined Mq Plugin
func (mq *MQ) Plugin() bulrush.PNRet {
	return func() *MQ {
		return mq
	}
}
