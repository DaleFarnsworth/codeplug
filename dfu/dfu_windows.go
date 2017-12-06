// Copyright 2017 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of Dfu.
//
// Dfu is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU Lesser General Public
// License as published by the Free Software Foundation.
//
// Dfu is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Dfu.  If not, see <http://www.gnu.org/licenses/>.

// This code began as a transliteration of the python code found in
// https://github.com/travisgoodspeed/md380tools.

// Package dfu implements reading/writing from/to the md380 radio via usb.
package dfu

import (
	"github.com/dalefarnsworth/codeplug/stdfu"
)

func New(progressCallback func(progressCounter int) bool) (*Dfu, error) {
	stDfu, err := stdfu.New()
	if err != nil {
		return nil, err
	}

	dfu := &Dfu{
		stDfu:            stDfu,
		progressCallback: progressCallback,
		progressFunc:     func() error { return nil },
	}

	err = dfu.enterDfuMode()
	if err != nil {
		dfu.Close()
		return nil, err
	}

	dfu.blockSize = 1024
	dfu.eraseBlockSize = 64 * 1024

	return dfu, nil
}
