package commonbl

// Copyright 2021 by tobi@backfrak.de. All
// rights reserved. Use of this source code is governed
// by a BSD-style license that can be found in the
// LICENSE file.

// Parmeters - Data structure that stores the common paramters for the executables in this appalication
type Parmeters struct {
	PrintVersion bool
	Verbose      bool
	Help         bool
	Test         bool
	LogFilePath  string
}
