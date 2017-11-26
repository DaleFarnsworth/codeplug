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
	dev            *gousb.Device
	iface          *gousb.Interface
	ifaceDone      func()
	ctx            *gousb.Context
	blockSize      int
	eraseBlockSize int
	progress       func(min, max, cur int) bool
	tickFunc       func() bool
	minTick        int
	maxTick        int
	curTick        int
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

	dfu.blockSize = 1024
	dfu.eraseBlockSize = 64 * 1024

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

	err := dfu.md380Cmd(
		0x91, 0x01, // Programming Mode
		0xa2, 0x08, // Access clock memory
	)
	if err != nil {
		goto errRet
	}

	timeBytes, err = dfu.readBlock(0, timeBytes) // Read BCD time bytes
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

	err := dfu.md380Cmd(0x91, 0x02)
	if err != nil {
		return err
	}

	err = dfu.writeBlock(0, bytes)
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

func (dfu *DFU) writeBlock(blockNumber int, bytes []byte) error {
	_, err := dfu.dev.Control(0x21, reqWrite, uint16(blockNumber), 0, bytes)
	if err != nil {
		return fmt.Errorf("write error: %s", err.Error())
	}

	return nil
}

func (dfu *DFU) setAddress(address int) error {
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

func (dfu *DFU) eraseBlocks(address int, size int) error {
	count := (size + dfu.eraseBlockSize - 1) / dfu.eraseBlockSize
	addr := uint32(address)
	for i := 0; i < count; i++ {
		err := dfu.eraseBlock(addr)
		if err != nil {
			return err
		}
		addr += uint32(dfu.eraseBlockSize)
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

type md380Cmd struct {
	a int
	b int
}

const CmdSleep = -2

func (dfu *DFU) md380Cmd(commands []md380Cmd) error {
	for _, cmd := range commands {
		switch cmd.a {
		case CmdSleep:
			dfu.sleepMilliseconds(cmd.b)
			continue
		}

		err := dfu.md380Custom(cmd.a, cmd.b)
		if err != nil {
			return fmt.Errorf("md380Cmd: %s", err.Error())
		}
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

func (dfu *DFU) md380Custom(a int, b int) error {
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

func (dfu *DFU) readBlock(blockNumber int, data []byte) error {
	_, err := dfu.dev.Control(0xa1, reqRead, uint16(blockNumber), 0, data)
	if err != nil {
		return fmt.Errorf("read: %s", err.Error())
	}

	return nil
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

func (dfu *DFU) setTickCounts(min, max, cur int) {
	dfu.tickFunc = func() bool { return true }
	if dfu.progress != nil {
		dfu.minTick = min
		dfu.maxTick = max
		dfu.curTick = cur
		dfu.tickFunc = func() bool {
			dfu.curTick++
			minMax := dfu.curTick * 50 / 49
			if dfu.maxTick < minMax {
				dfu.maxTick = minMax
				//fmt.Fprintf(os.Stderr, "maxTick %d\n", dfu.maxTick)
			}
			return dfu.progress(dfu.minTick, dfu.maxTick, dfu.curTick-1)
		}
		dfu.progress(dfu.minTick, dfu.maxTick, dfu.curTick)
	}
}

func (dfu *DFU) read(address, offset int, data []byte) error {
	if offset%dfu.blockSize != 0 {
		return fmt.Errorf("dfu.read: offset is not a multiple of blockSize")
	}

	if len(data)%dfu.blockSize != 0 {
		return fmt.Errorf("dfu.read: data size is not a multiple of blockSize")
	}

	blockNumber := offset / dfu.blockSize
	blockCount := len(data) / dfu.blockSize

	err := dfu.setAddress(address)
	if err != nil {
		goto errRet
	}

	dfu.setTickCounts(0, blockCount, 0)

	offset = 0
	for i := 0; i < blockCount; i++ {
		if !dfu.tickFunc() {
			return nil
		}
		endOffset := offset + dfu.blockSize
		err = dfu.readBlock(blockNumber, data[offset:endOffset])
		if err != nil {
			goto errRet
		}

		blockNumber++
		offset += dfu.blockSize
	}

	if dfu.progress != nil {
		dfu.curTick = dfu.maxTick
		dfu.progress(dfu.minTick, dfu.maxTick, dfu.curTick)
	}

	return nil

errRet:
	return fmt.Errorf("dfu.read: %s", err.Error())
}

func (dfu *DFU) write(address, offset int, data []byte) error {
	blockNumber := offset / dfu.blockSize
	blockCount := len(data) / dfu.blockSize

	err := dfu.eraseBlocks(0x00000000, len(data))
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

	dfu.setTickCounts(0, blockCount, 0)

	offset = 0
	for i := 0; i < blockCount; i++ {
		if !dfu.tickFunc() {
			return nil
		}
		endOffset := offset + dfu.blockSize
		err = dfu.writeBlock(blockNumber, data[offset:endOffset])
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
		offset += dfu.blockSize
	}

	if dfu.progress != nil {
		dfu.curTick = dfu.maxTick
		dfu.progress(dfu.minTick, dfu.maxTick, dfu.curTick)
	}

	return nil

errRet:
	return fmt.Errorf("dfu.write: %s", err.Error())
}

func (dfu *DFU) ReadCodeplug(data []byte, progress func(min, max, val int) bool) error {
	dfu.progress = progress

	dfu.setTickCounts(0, 62, 1)

	err := dfu.md380Cmd([]md380Cmd{
		md380Cmd{0x91, 0x01}, // Programming Mode
		md380Cmd{0xa2, 0x02},
		md380Cmd{0xa2, 0x02},
		md380Cmd{0xa2, 0x03},
		md380Cmd{0xa2, 0x04},
		md380Cmd{0xa2, 0x07},
	})
	if err != nil {
		goto errRet
	}

	err = dfu.read(0, 2048, data)
	if err != nil {
		goto errRet
	}

	dfu.progress = nil
	return nil

errRet:
	dfu.progress = nil
	return fmt.Errorf("ReadCodeplug: %s", err.Error())
}

func (dfu *DFU) WriteCodeplug(data []byte, progress func(min, max, val int) bool) error {
	dfu.progress = progress
	dfu.setTickCounts(0, 276, 1)

	if len(data)%dfu.blockSize != 0 {
		return fmt.Errorf("WriteCodeplug: codeplug data size is not a multiple of %d", dfu.blockSize)
	}

	err := dfu.md380Cmd([]md380Cmd{
		md380Cmd{0x91, 0x01}, // Programming Mode
		md380Cmd{0x91, 0x01}, // Programming Mode
		md380Cmd{0xa2, 0x02},
		md380Cmd{CmdSleep, 2000},
		md380Cmd{0xa2, 0x02},
		md380Cmd{0xa2, 0x03},
		md380Cmd{0xa2, 0x04},
		md380Cmd{0xa2, 0x07},
	})
	if err != nil {
		goto errRet
	}

	dfu.write(0, 2048, data)
	dfu.progress = nil

	return nil

errRet:
	dfu.progress = nil
	return fmt.Errorf("WriteCodeplug: %s", err.Error())
}
