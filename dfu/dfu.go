// Copyright 2017 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of DFU.
//
// DFU is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU Lesser General Public
// License as published by the Free Software Foundation.
//
// DFU is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with DFU.  If not, see <http://www.gnu.org/licenses/>.

// This code began as a transliteration of the python code found in
// https://github.com/travisgoodspeed/md380tools.

// Package dfu implements reading/writing from/to the md380 radio via usb.
package dfu

import (
	"fmt"
	"os"
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

type state int

const (
	appIdle state = iota
	appDetach
	dfuIdle
	dfuWriteSync
	dfuWriteBusy
	dfuWriteIdle
	dfuManifestSync
	dfuManifest
	dfuManifestWaitReset
	dfuReadIdle
	dfuError
)

var stateStrings = map[state]string{
	appIdle:              "appIdle",
	appDetach:            "appDetach",
	dfuIdle:              "dfuIdle",
	dfuWriteSync:         "dfuWriteSync",
	dfuWriteBusy:         "dfuWriteBusy",
	dfuWriteIdle:         "dfuWriteIdle",
	dfuManifestSync:      "dfuManifestSync",
	dfuManifest:          "dfuManifest",
	dfuManifestWaitReset: "dfuManifestWaitReset",
	dfuReadIdle:          "dfuReadIdle",
	dfuError:             "dfuError",
}

func (s state) String() string {
	return stateStrings[s]
}

type status int

const (
	statusOk status = iota
	errTarget
	errFile
	errWrite
	errErase
	errCheckErased
	errProgram
	errVerify
	errAddress
	errNotDone
	errFirmware
	errVendor
	errUsbR
	errPOR
	errUnknown
	errStalledPkt
)

var statusStrings = map[status]string{
	statusOk:       "ok",
	errTarget:      "errTarget",
	errFile:        "errFile",
	errWrite:       "errWrite",
	errErase:       "errErase",
	errCheckErased: "errCheckErased",
	errProgram:     "errProgram",
	errVerify:      "errVerify",
	errAddress:     "errAddress",
	errNotDone:     "errNotDone",
	errFirmware:    "errFirmware",
	errVendor:      "errVendor",
	errUsbR:        "errUsbR",
	errPOR:         "errPOR",
	errUnknown:     "errUnknown",
	errStalledPkt:  "errStalledPkt",
}

func (s status) String() string {
	return statusStrings[s]
}

const (
	md380Vendor  = 0x0483
	md380Product = 0xdf11
)

type DFU struct {
	dev       *gousb.Device
	iface     *gousb.Interface
	ifaceDone func()
	ctx       *gousb.Context
	tickFunc  func() bool
	minTick   int
	maxTick   int
	curTick   int
}

func NewDFU() (*DFU, error) {
	ctx := gousb.NewContext()

	dfu := &DFU{
		ctx:      ctx,
		tickFunc: func() bool { return true },
	}

	dev, err := ctx.OpenDeviceWithVIDPID(md380Vendor, md380Product)
	if err != nil {
		dfu.Close()
		return nil, fmt.Errorf("OpenDevice failed: %v", err)
	}
	if dev == nil {
		dfu.Close()
		return nil, fmt.Errorf("No Radio found on USB")
	}
	dfu.dev = dev

	iface, ifaceDone, err := dev.DefaultInterface()
	if err != nil {
		dfu.Close()
		return nil, fmt.Errorf("%s.DefaultInterface failed: %v", dev, err)
	}
	dfu.iface = iface
	dfu.ifaceDone = ifaceDone

	dev.ControlTimeout = time.Duration(3 * time.Second)

	err = dfu.enterDFUMode()
	if err != nil {
		dfu.Close()
		if err == gousb.ErrorPipe {
			return nil, fmt.Errorf("Failed to enter DFU mode.\nIs bootloader running?")
		}
		return nil, err
	}

	return dfu, nil
}

func (dfu *DFU) Close() {
	if dfu.ifaceDone != nil {
		dfu.ifaceDone()
	}
	if dfu.dev != nil {
		dfu.dev.Close()
	}
	if dfu.ctx != nil {
		dfu.ctx.Close()
	}
}

func (dfu *DFU) detach() error {
	_, err := dfu.dev.Control(0x21, reqDetach, 0, 0, nil)
	if err != nil {
		return fmt.Errorf("detach: %s", err.Error())
	}

	return nil
}

/* This commented-out code is untested and unused.

func (dfu *DFU) GetString(index int) (string, error) {
	str, err := dfu.dev.GetStringDescriptor(index)
	if err != nil {
		goto errRet
	}

	return str, nil

errRet:
	return str, fmt.Errorf("GetStringDescriptor failed: %s", err.Error())
}

func (dfu *DFU) toDecimal(b byte) int {
	return int(b&0xf + (b>>4)*10)
}

func (dfu *DFU) toBCD(i int) byte {
	return byte(i/10<<4 | i%10)
}

func (dfu *DFU) GetTime() (time.Time, error) {
	var year, day, hours, minutes, seconds int
	var month time.Month
	var timeBytes []byte
	location := time.Local

	err := dfu.md380Custom(
		0x91, 0x01, // Programming Mode
		0xa2, 0x08, // Access clock memory
	)
	if err != nil {
		goto errRet
	}

	timeBytes, err = dfu.read(0, 7) // Read the time bytes as BCD
	if err != nil {
		goto errRet
	}

	year = dfu.toDecimal(timeBytes[0])*100 + dfu.toDecimal(timeBytes[1])
	month = time.Month(dfu.toDecimal(timeBytes[2]))
	day = dfu.toDecimal(timeBytes[3])
	hours = dfu.toDecimal(timeBytes[4])
	minutes = dfu.toDecimal(timeBytes[5])
	seconds = dfu.toDecimal(timeBytes[6])

	return time.Date(year, month, day, hours, minutes, seconds, 0, location), nil

errRet:
	return time.Now(), fmt.Errorf("GetTime: %s", err.Error())
}

func (dfu *DFU) SetTime(t time.Time) error {
	year, month, day := t.Date()
	hours, minutes, seconds := t.Clock()
	t.Format("20060102150405")
	bytes := make([]byte, 7)
	bytes[0] = dfu.toBCD(year / 100)
	bytes[1] = dfu.toBCD(year % 100)
	bytes[2] = dfu.toBCD(int(month))
	bytes[3] = dfu.toBCD(day)
	bytes[4] = dfu.toBCD(hours)
	bytes[5] = dfu.toBCD(minutes)
	bytes[6] = dfu.toBCD(seconds)

	err := dfu.md380Custom(0x91, 0x02)
	if err != nil {
		return err
	}

	err = dfu.write(0, bytes)
	if err != nil {
		return err
	}

	err = dfu.waitUntilReady()
	if err != nil {
		return err
	}

	err = dfu.md380Reboot()
	if err != nil {
		return err
	}

	return nil
}

func (dfu *DFU) md380Reboot() error {
	bytes := []byte{byte(0x91), byte(0x05)}
	_, err := dfu.dev.Control(0x21, reqWrite, 0, 0, bytes)
	if err != nil {
		goto errRet
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		goto errRet
	}

	return nil

errRet:
	return fmt.Errorf("md380Reboot: %s", err.Error())
}

func (dfu *DFU) waitUntilReady() error {
	var err error

	for {
		state, err := dfu.getStatus()
		if err != nil {
			goto errRet
		}
		if state == dfuIdle {
			break
		}

		err = dfu.clearStatus()
		if err != nil {
			goto errRet
		}
	}

	return nil

errRet:
	return fmt.Errorf("waitUntilReady %s", err.Error())
}

*/

func (dfu *DFU) write(blockNumber int, bytes []byte) error {
	bn := uint16(blockNumber)
	_, err := dfu.dev.Control(0x21, reqWrite, bn, 0, bytes)
	if err != nil {
		return fmt.Errorf("write error: %s", err.Error())
	}

	return nil
}

func (dfu *DFU) setAddress(address uint32) error {
	var state state

	a := byte(address)
	b := byte((address >> 8))
	c := byte((address >> 16))
	d := byte((address >> 24))
	bytes := []byte{0x21, a, b, c, d}

	_, err := dfu.dev.Control(0x21, reqWrite, 0, 0, bytes)
	if err != nil {
		goto errRet
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		goto errRet
	}

	state, err = dfu.getStatus() // this actually gets the state
	if err != nil {
		goto errRet
	}

	if state != dfuWriteIdle {
		goto errRet
	}

	err = dfu.enterDFUMode()
	if err != nil {
		goto errRet
	}

	return nil

errRet:
	return fmt.Errorf("setAddress: %s", err.Error())
}

func (dfu *DFU) eraseBlocks(addresses ...uint32) error {
	for _, addr := range addresses {
		err := dfu.eraseBlock(addr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dfu *DFU) eraseBlock(address uint32) error {
	var state state

	a := byte(address)
	b := byte((address >> 8))
	c := byte((address >> 16))
	d := byte((address >> 24))
	bytes := []byte{0x41, a, b, c, d}

	_, err := dfu.dev.Control(0x21, reqWrite, 0, 0, bytes)
	if err != nil {
		goto errRet
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		goto errRet
	}
	state, err = dfu.getStatus() // this actually gets the state
	if err != nil {
		goto errRet
	}
	if state != dfuWriteIdle {
		goto errRet
	}

	err = dfu.enterDFUMode()
	if err != nil {
		goto errRet
	}

	return nil

errRet:
	return fmt.Errorf("EraseBlock: %s", err.Error())
}

func (dfu *DFU) md380Custom(args ...int) error {
	for len(args) > 0 {
		a := args[0]
		b := args[1]

		err := dfu.md380Custom1(a, b)
		if err != nil {
			return fmt.Errorf("md380Custom: %s", err.Error())
		}

		args = args[2:]
	}

	return nil
}

func (dfu *DFU) sleepMilliseconds(millis int) {
	tens := (millis + 9) / 10
	for i := 0; i < tens; i++ {
		dfu.tickFunc()
		time.Sleep(time.Duration(10 * time.Millisecond))
	}
}

func (dfu *DFU) md380Custom1(a int, b int) error {
	bytes := []byte{byte(a), byte(b)}

	_, err := dfu.dev.Control(0x21, reqWrite, 0, 0, bytes)
	if err != nil {
		return err
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		return err
	}

	dfu.sleepMilliseconds(100)

	state, err := dfu.getStatus() // this actually gets the state
	if err != nil {
		return err
	}

	if state != dfuWriteIdle {
		return err
	}

	err = dfu.enterDFUMode()
	if err != nil {
		return err
	}

	return nil
}

func (dfu *DFU) read(blockNumber int, length int) ([]byte, error) {
	bn := uint16(blockNumber)
	bytes := make([]byte, length)

	_, err := dfu.dev.Control(0xa1, reqRead, bn, 0, bytes)
	if err != nil {
		return nil, fmt.Errorf("read: %s", err.Error())
	}

	return bytes, nil
}

func (dfu *DFU) getCommand() ([]byte, error) {
	bytes := make([]byte, 32)

	_, err := dfu.dev.Control(0xa1, reqRead, 0, 0, bytes)
	if err != nil {
		goto errRet
	}

	_, err = dfu.getStatus()
	if err != nil {
		goto errRet
	}

	return bytes, nil

errRet:
	return nil, fmt.Errorf("getCommand: %s", err.Error())
}

func (dfu *DFU) getStatus() (state, error) {
	bytes := make([]byte, 6)
	_, err := dfu.dev.Control(0xa1, reqGetStatus, 0, 0, bytes)
	if err != nil {
		return 0, fmt.Errorf("getStatus: ", err.Error())
	}
	state := state(bytes[4])

	debug := false
	if debug {
		status := status(bytes[0]).String()
		timeout := (((bytes[1] << 8) | bytes[2]) << 8) | bytes[3]
		discarded := bytes[5]
		fmt.Fprintf(os.Stderr, status, timeout, state, discarded)
	}

	return state, nil
}

func (dfu *DFU) clearStatus() error {
	_, err := dfu.dev.Control(0x21, reqClearStatus, 0, 0, nil)
	if err != nil {
		return fmt.Errorf("clearStatus: %s", err.Error())
	}

	return nil
}

func (dfu *DFU) getState() (state, error) {
	bytes := make([]byte, 1)

	_, err := dfu.dev.Control(0xa1, reqGetState, 0, 0, bytes)
	if err != nil {
		return 0, fmt.Errorf("GetStatus: %s", err.Error())
	}

	return state(bytes[0]), nil
}

func (dfu *DFU) abort() error {
	_, err := dfu.dev.Control(0x21, reqAbort, 0, 0, nil)
	if err != nil {
		return fmt.Errorf("ClearStatus: %s", err.Error())
	}

	return nil
}

func (dfu *DFU) enterDFUMode() error {
	var err error
	var stat state
	actionMap := map[state]func() error{
		dfuWriteSync:         dfu.abort,
		dfuWriteIdle:         dfu.abort,
		dfuManifestSync:      dfu.abort,
		dfuReadIdle:          dfu.abort,
		dfuError:             dfu.clearStatus,
		appIdle:              dfu.detach,
		appDetach:            dfu.wait,
		dfuWriteBusy:         dfu.wait,
		dfuManifest:          dfu.abort,
		dfuManifestWaitReset: dfu.wait,
		dfuIdle:              dfu.wait,
	}

	for {
		stat, err = dfu.getState()
		if err != nil {
			goto errRet
		}
		if stat == dfuIdle {
			break
		}
		err = actionMap[stat]()
		if err != nil {
			goto errRet
		}
	}

	return nil

errRet:
	return fmt.Errorf("enterDFUMode %s", err.Error())
}

func (dfu *DFU) wait() error {
	dfu.sleepMilliseconds(100)
	return nil
}

func (dfu *DFU) ReadCodeplug(progress func(min, max, val int) bool) ([]byte, error) {
	dfu.tickFunc = func() bool { return true }
	if progress != nil {
		dfu.minTick = 0
		dfu.maxTick = 62
		dfu.curTick = 1
		dfu.tickFunc = func() bool {
			dfu.curTick++
			minMax := dfu.curTick * 50 / 49
			if dfu.maxTick < minMax {
				dfu.maxTick = minMax
				//fmt.Fprintf(os.Stderr, "maxTick %d\n", dfu.maxTick)
			}
			return progress(dfu.minTick, dfu.maxTick, dfu.curTick-1)
		}
	}

	var bytes []byte

	blockSize := 1024
	blockNumber := 2
	blockCount := 256

	data := make([]byte, 0, blockSize*blockCount)

	err := dfu.md380Custom(
		0x91, 0x01, // Programming Mode
		0xa2, 0x02,
		0xa2, 0x02,
		0xa2, 0x03,
		0xa2, 0x04,
		0xa2, 0x07,
	)
	if err != nil {
		goto errRet
	}

	err = dfu.setAddress(0x00000000)
	if err != nil {
		goto errRet
	}

	if progress != nil {
		dfu.curTick = dfu.maxTick - 1
		dfu.tickFunc()

		dfu.minTick = 0
		dfu.maxTick = blockCount
		dfu.curTick = 0
		dfu.tickFunc = func() bool {
			dfu.curTick++
			return progress(dfu.minTick, dfu.maxTick, dfu.curTick-1)
		}
	}

	for i := 0; i < blockCount; i++ {
		if !dfu.tickFunc() {
			return nil, nil
		}

		bytes, err = dfu.read(blockNumber, blockSize)
		if err != nil {
			goto errRet
		}

		blockNumber++
		if len(bytes) != blockSize {
			return nil, fmt.Errorf("bad read size: %d bytes", len(bytes))
		}
		data = append(data, bytes...)
	}

	dfu.curTick = dfu.maxTick
	dfu.tickFunc()

	return data, nil

errRet:
	return nil, fmt.Errorf("ReadCodeplug: %s", err.Error())
}

func (dfu *DFU) WriteCodeplug(data []byte, progress func(min, max, val int) bool) error {
	dfu.tickFunc = func() bool { return true }
	if progress != nil {
		dfu.minTick = 0
		dfu.maxTick = 276
		dfu.curTick = 1
		dfu.tickFunc = func() bool {
			dfu.curTick++
			minMax := dfu.curTick * 50 / 49
			if dfu.maxTick < minMax {
				dfu.maxTick = minMax
				//fmt.Fprintf(os.Stderr, "maxTick %d\n", dfu.maxTick)
			}
			return progress(dfu.minTick, dfu.maxTick, dfu.curTick-1)
		}
	}

	//var cmd []byte
	blockSize := 1024
	blockNumber := 2
	blockCount := len(data) / blockSize

	if len(data)%blockSize != 0 {
		return fmt.Errorf("codeplug data size is not a multiple of %d", blockSize)
	}

	err := dfu.md380Custom(
		0x91, 0x01, // Programming Mode
		0x91, 0x01, // Programming Mode
		0xa2, 0x02,
	)
	if err != nil {
		goto errRet
	}

	//cmd, err = dfu.getCommand()
	//if err != nil {
	//	goto errRet
	//}
	//fmt.Fprintf(os.Stderr, "%+v", cmd)

	dfu.sleepMilliseconds(2000)

	err = dfu.md380Custom(
		0xa2, 0x02,
		0xa2, 0x03,
		0xa2, 0x04,
		0xa2, 0x07,
	)
	if err != nil {
		goto errRet
	}

	err = dfu.eraseBlocks(0x00000000, 0x00010000, 0x00020000, 0x00030000)
	if err != nil {
		goto errRet
	}

	err = dfu.setAddress(0x00000000)
	if err != nil {
		goto errRet
	}

	_, err = dfu.getStatus()
	if err != nil {
		goto errRet
	}

	if progress != nil {
		dfu.curTick = dfu.maxTick - 1
		dfu.tickFunc()

		dfu.minTick = 0
		dfu.maxTick = blockCount
		dfu.curTick = 0
		dfu.tickFunc = func() bool {
			dfu.curTick++
			return progress(dfu.minTick, dfu.maxTick, dfu.curTick-1)
		}
	}

	for i := 0; i < blockCount; i++ {
		if !dfu.tickFunc() {
			return nil
		}

		offset := i * blockSize
		bytes := data[offset : offset+blockSize]
		err = dfu.write(blockNumber, bytes)
		if err != nil {
			goto errRet
		}

		for {
			state, err := dfu.getStatus()
			if err != nil {
				goto errRet
			}

			if state == dfuWriteIdle {
				break
			}
		}
		blockNumber++
	}

	dfu.curTick = dfu.maxTick
	dfu.tickFunc()

	return nil

errRet:
	return fmt.Errorf("WriteCodeplug: %s", err.Error())
}
