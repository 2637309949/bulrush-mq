// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mq

import (
	"github.com/thoas/go-funk"
)

// MemoModel defined memory model
type MemoModel struct {
	Model
	mess []Message
}

// Save defined store message
func (m *MemoModel) Save(mes Message) {
	m.mess = append(m.mess, mes)
}

// Find defined find message
func (m *MemoModel) Find(mtype string, status string) []Message {
	return funk.Filter(m.mess, func(mes Message) bool {
		return mes.Type == mtype && mes.Status == status
	}).([]Message)
}

// Count defined count message
func (m *MemoModel) Count(mtype string, status string) uint {
	mess := funk.Filter(m.mess, func(mes Message) bool {
		return mes.Type == mtype && mes.Status == status
	}).([]Message)
	return uint(len(mess))
}

// Update defined update message
func (m *MemoModel) Update(ms Message, status string) error {
	mess := funk.Find(m.mess, func(mes Message) bool {
		return mes.Type == ms.Type && mes.Status == ms.Status
	}).(Message)
	mess.Status = status
	return nil
}
