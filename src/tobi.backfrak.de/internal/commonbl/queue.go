package commonbl

// Copyright 2022 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"container/list"
	"sync"
)

// StringQueue, a type that represents a FI/FI Stack for strings
type StringQueue struct {
	mList  *list.List
	mMutex sync.Mutex
}

// Get a new empty instance of the StringQueue
func NewStringQueue() *StringQueue {
	var queue StringQueue
	queue.mList = list.New()

	return &queue
}

// Push (add) a string to the StringQueue
func (queue *StringQueue) Push(value string) {
	queue.mMutex.Lock()
	defer queue.mMutex.Unlock()
	queue.mList.PushBack(value)
}

// Pull (get and remove) a string from the StringQueue.
// Returns an error when the Queue is empty
func (queue *StringQueue) Pull() (string, error) {
	if queue.mList.Len() <= 0 {
		return "", NewEmptyStringQueueError()
	}
	var valueString string
	queue.mMutex.Lock()
	defer queue.mMutex.Unlock()

	e := queue.mList.Front()
	valueString = e.Value.(string)
	queue.mList.Remove(e)
	return valueString, nil
}

// Tell if the StringQueue is empty
func (queue *StringQueue) IsEmpty() bool {
	return queue.mList.Len() <= 0
}
