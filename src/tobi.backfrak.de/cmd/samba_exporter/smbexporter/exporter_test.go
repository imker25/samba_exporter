package smbexporter

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

import (
	"testing"

	"tobi.backfrak.de/internal/commonbl"
)

func TestNewSambaExporter(t *testing.T) {
	requestHandler := *commonbl.NewPipeHandler(true, commonbl.RequestPipe)
	responseHandler := *commonbl.NewPipeHandler(true, commonbl.ResposePipe)
	logger := *commonbl.NewLogger(true)
	exporter := NewSambaExporter(requestHandler, responseHandler, logger)

	if exporter.RequestHandler.PipeType != commonbl.RequestPipe {
		t.Errorf("The exporter.RequestHandler is not of the expected type")
	}

	if exporter.ResponseHander.PipeType != commonbl.ResposePipe {
		t.Errorf("The exporter.RequestHandler is not of the expected type")
	}

	if exporter.hostName == "" {
		t.Errorf("The exporter.hostName is empty")
	}

	if exporter.Descriptions == nil {
		t.Errorf("exporter.Descriptions are nil")
	}

	if logger.Verbose != exporter.Logger.Verbose {
		t.Errorf("The exporter uses the wrong logger")
	}
}
