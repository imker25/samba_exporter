package commonbl

import (
	"fmt"
	"reflect"
	"testing"
)

// Copyright 2022 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestNewStringQueue(t *testing.T) {
	qu := NewStringQueue()

	if qu.IsEmpty() == false {
		t.Errorf("The StringQueue is not empty right after creation")
	}

	_, err := qu.Pull()
	if err == nil {
		t.Errorf("Got no error when pull on empty queue")
	}

	switch err.(type) {
	case *EmptyStringQueueError:
		fmt.Println("OK")
	default:
		t.Errorf("Expected a EmptyStringQueueError, got a %s", reflect.TypeOf(err))
	}
}

func TestStringQueue(t *testing.T) {

	var pull string
	var err error
	qu := NewStringQueue()

	if qu.IsEmpty() == false {
		t.Errorf("The StringQueue is not empty right after creation")
	}

	qu.Push("a")

	if qu.IsEmpty() == true {
		t.Errorf("The StringQueue is empty right after 1. push")
	}

	qu.Push("b")
	qu.Push("c")
	qu.Push("d")

	if qu.IsEmpty() == true {
		t.Errorf("The StringQueue is empty right after 4. push")
	}

	pull, err = qu.Pull()
	if err != nil {
		t.Errorf("Got unexpected error while pull for \"a\"")
	}
	if pull != "a" {
		t.Errorf("Got \"%s\", but expected \"%s\"", pull, "a")
	}

	pull, err = qu.Pull()
	if err != nil {
		t.Errorf("Got unexpected error while pull for \"b\"")
	}
	if pull != "b" {
		t.Errorf("Got \"%s\", but expected \"%s\"", pull, "b")
	}

	pull, err = qu.Pull()
	if err != nil {
		t.Errorf("Got unexpected error while pull for \"c\"")
	}
	if pull != "c" {
		t.Errorf("Got \"%s\", but expected \"%s\"", pull, "c")
	}

	if qu.IsEmpty() == true {
		t.Errorf("The StringQueue is empty right after 4. push and 3. pull ")
	}

	pull, err = qu.Pull()
	if err != nil {
		t.Errorf("Got unexpected error while pull for \"d\"")
	}
	if pull != "d" {
		t.Errorf("Got \"%s\", but expected \"%s\"", pull, "d")
	}

	if qu.IsEmpty() == false {
		t.Errorf("The StringQueue is not empty right after pulling all elements ")
	}
	_, err = qu.Pull()
	if err == nil {
		t.Errorf("Got no error when pull on empty queue")
	}

}
