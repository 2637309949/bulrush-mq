// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mq

import (
	"math/rand"
	"sort"

	"github.com/thoas/go-funk"
)

// Memo defined memory model
type (
	Memo struct {
		Model
		mess []Message
	}
	// Model defined model store
	Model interface {
		Save(Message)
		Find(string, string) *Message
		Count(string, string) uint
		Update(*Message, string) error
	}
)

// Save defined store message
func (m *Memo) Save(mes Message) {
	mes.ID = rand.Int()
	m.mess = append(m.mess, mes)
}

// Find defined find message
func (m *Memo) Find(mtype string, status string) *Message {
	iTask := funk.Filter(m.mess, func(mes Message) bool {
		return mes.Type == mtype && mes.Status == status
	}).([]Message)
	sort.Sort(sortByMsAt(iTask))
	if len(iTask) >= 1 {
		return &iTask[0]
	}
	return nil
}

// Count defined count message
func (m *Memo) Count(mtype string, status string) uint {
	mess := funk.Filter(m.mess, func(mes Message) bool {
		return mes.Type == mtype && mes.Status == status
	}).([]Message)
	return uint(len(mess))
}

// Update defined update message
func (m *Memo) Update(ms *Message, status string) error {
	for i, mes := range m.mess {
		if mes.ID == ms.ID {
			mes.Status = status
			m.mess[i].Status = status
		}
	}
	return nil
}
