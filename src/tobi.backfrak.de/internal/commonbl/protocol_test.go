package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetIdFromRequest(t *testing.T) {
	id, err := GetIdFromRequest("bal: 23")

	if err != nil {
		t.Errorf("Got error \"%s\" but expected none", err)
	}

	if id != 23 {
		t.Errorf("The id \"%d\" is not the expected", id)
	}

	id, err = GetIdFromRequest("bal: 23: sert")
	if err == nil {
		t.Errorf("Got no error but expected one")
	}

	if id != 0 {
		t.Errorf("The id \"%d\" is not the expected", id)
	}

	id, err = GetIdFromRequest("bal: 23  sert")
	if err == nil {
		t.Errorf("Got no error but expected one")
	}

	if id != 0 {
		t.Errorf("The id \"%d\" is not the expected", id)
	}
}

func TestGetRequest(t *testing.T) {
	id := 23
	rType := RequestType("bal:")
	request := GetRequest(rType, id)

	if strings.Contains(request, string(rType)) == false {
		t.Errorf("The request does not contain the expected request type")
	}

	if strings.Contains(request, fmt.Sprintf("%d", id)) == false {
		t.Errorf("The request does not contain the expected id")
	}

}

func TestGetTestResponseHeader(t *testing.T) {
	id := 23
	rType := RequestType("bal:")
	request := GetTestResponseHeader(rType, id)

	if strings.Contains(request, string(rType)) == false {
		t.Errorf("The request does not contain the expected request type")
	}

	if strings.Contains(request, fmt.Sprintf("%d", id)) == false {
		t.Errorf("The request does not contain the expected id")
	}

}

func TestGetResponseHeader(t *testing.T) {
	id := 23
	rType := RequestType("bal:")
	request := GetResponseHeader(rType, id)

	if strings.Contains(request, string(rType)) == false {
		t.Errorf("The request does not contain the expected request type")
	}

	if strings.Contains(request, fmt.Sprintf("%d", id)) == false {
		t.Errorf("The request does not contain the expected id")
	}

}

func TestGetResponse(t *testing.T) {
	id := 23
	rType := RequestType("bal:")
	header := GetResponseHeader(rType, id)
	data := "my data\nis hot\nor not"

	response := GetResponse(header, data)

	rHeader, rData, err := SplitResponse(response)

	if err != nil {
		t.Errorf("Got error \"%s\" but expected none", err)
	}

	if rHeader != header {
		t.Errorf("The header is not the expected")
	}

	if rData != data {
		t.Errorf("The data is not the expected")
	}
}

func TestSplitResponse(t *testing.T) {

	response := fmt.Sprintf("%s Response for id my data is hot or not", SHARE_REQUEST)

	header, data, err := SplitResponse(response)

	if header != "" {
		t.Errorf("The header is not the expected")
	}
	if data != "" {
		t.Errorf("The data is not the expected")
	}

	if err == nil {
		t.Errorf("Got no error but expected one")
	}
}

func TestCheckResponseHeader(t *testing.T) {
	id := 23
	rType := RequestType("bal:")
	header := GetResponseHeader(rType, id)

	if CheckResponseHeader(header, rType, id) == false {
		t.Errorf("CheckResponseHeader is false, but expected true")
	}

	if CheckResponseHeader("my header", rType, id) == true {
		t.Errorf("CheckResponseHeader is true, but expected false")
	}

}
