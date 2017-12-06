// Copyright 2017 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of StDFU.
//
// StDFU is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU General Public License
// as published by the Free Software Foundation.
//
// StDFU is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with StDFU.  If not, see <http://www.gnu.org/licenses/>.

package stdfu

import (
	"fmt"
	"time"

	"github.com/google/gousb"
)

const (
	reqDetach byte = iota
	reqWrite
	reqRead
	reqGetStatus
	reqClearStatus
	reqGetState
	reqAbort
)

type StDfu struct {
	dev       *gousb.Device
	iface     *gousb.Interface
	ifaceDone func()
	ctx       *gousb.Context
}

func New() (*StDfu, error) {
	ctx := gousb.NewContext()

	stDfu := &StDfu{
		ctx: ctx,
	}

	const (
		md380Vendor  = 0x0483
		md380Product = 0xdf11
	)

	dev, err := ctx.OpenDeviceWithVIDPID(md380Vendor, md380Product)
	if err != nil {
		stDfu.Close()
		return nil, fmt.Errorf("OpenDevice failed: %v", err)
	}
	if dev == nil {
		stDfu.Close()
		return nil, fmt.Errorf("No Radio found on USB")
	}
	stDfu.dev = dev

	iface, ifaceDone, err := dev.DefaultInterface()
	if err != nil {
		stDfu.Close()
		dprint()
		return nil, fmt.Errorf("%s: DefaultInterface failed: %v", dev, err)
	}
	stDfu.iface = iface
	stDfu.ifaceDone = ifaceDone

	dev.ControlTimeout = time.Duration(3 * time.Second)

	return stDfu, nil
}

func wrapError(prefix string, err error) error {
	if err.Error() == "" {
		return err
	}
	return fmt.Errorf("%s: %s", prefix, err.Error())
}

func (stDfu *StDfu) Close() {
	if stDfu.ifaceDone != nil {
		stDfu.ifaceDone()
	}
	if stDfu.dev != nil {
		stDfu.dev.Close()
	}
	if stDfu.ctx != nil {
		stDfu.ctx.Close()
	}
}

func (stDfu *StDfu) Abort() error {
	_, err := stDfu.dev.Control(0x21, reqAbort, 0, 0, nil)
	if err != nil {
		return wrapError("ClearStatus", err)
	}

	return nil
}

func (stDfu *StDfu) ClrStatus() error {
	_, err := stDfu.dev.Control(0x21, reqClearStatus, 0, 0, nil)
	if err != nil {
		return wrapError("clearStatus", err)
	}

	return nil
}

func (stDfu *StDfu) Detach() error {
	_, err := stDfu.dev.Control(0x21, reqDetach, 0, 0, nil)
	if err != nil {
		return wrapError("detach", err)
	}

	return nil
}

func (stDfu *StDfu) Dnload(blockNumber int, buffer []byte) error {
	_, err := stDfu.dev.Control(0x21, reqWrite, uint16(blockNumber), 0, buffer)
	if err != nil {
		return wrapError("write error", err)
	}

	return nil
}

func (stDfu *StDfu) GetState() (State, error) {
	bytes := make([]byte, 1)

	_, err := stDfu.dev.Control(0xa1, reqGetState, 0, 0, bytes)
	if err != nil {
		return 0, wrapError("GetStatus", err)
	}

	return State(bytes[0]), nil
}

func (stDfu *StDfu) GetStatus() (DfuStatus, error) {
	bytes := make([]byte, 6)

	_, err := stDfu.dev.Control(0xa1, reqGetStatus, 0, 0, bytes)
	if err != nil {
		err = wrapError("getStatus", err)
	}

	dfuStatus := DfuStatus{
		Status:      Status(bytes[0]),
		PollTimeout: int(bytes[1]<<16 | bytes[2]<<8 | bytes[3]),
		State:       State(bytes[4]),
		IString:     int(bytes[5]),
	}

	return dfuStatus, err
}

func (stDfu *StDfu) GetStringDescriptor(index int) (string, error) {
	str, err := stDfu.dev.GetStringDescriptor(index)
	if err != nil {
		err = wrapError("getString", err)
	}

	return str, err
}

func (stDfu *StDfu) Upload(blockNumber int, buffer []byte) error {
	_, err := stDfu.dev.Control(0xa1, reqRead, uint16(blockNumber), 0, buffer)
	if err != nil {
		err = wrapError("read", err)
	}

	return err
}
