package pipecomunication

import (
	"fmt"
	"os"
	"testing"

	"tobi.backfrak.de/internal/commonbl"
)

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

func TestGetSambaStatusTimeout(t *testing.T) {
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	_, _, _, err := GetSambaStatus(requestHandler, responseHandler)

	if err == nil {
		t.Errorf("Exptected an error but got none")
	}

	switch err.(type) {
	case *SmbStatusTimeOutError:
		fmt.Fprintln(os.Stdout, "OK")
	default:
		t.Errorf("Got error of the wrong type")
	}
}
