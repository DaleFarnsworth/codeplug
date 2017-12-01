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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
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

const (
	MinProgress = 0
	MaxProgress = 1000000
)

const spiEraseSPIFlashBlockDelay = 500 // milliseconds

type DFU struct {
	dev               *gousb.Device
	iface             *gousb.Interface
	ifaceDone         func()
	ctx               *gousb.Context
	blockSize         int
	eraseBlockSize    int
	progressCallback  func(progressCounter int) bool
	progressFunc      func() error
	progressIncrement int
	progressCounter   int
}

func NewDFU(progressCallback func(progressCounter int) bool) (*DFU, error) {
	ctx := gousb.NewContext()

	dfu := &DFU{
		ctx:              ctx,
		progressCallback: progressCallback,
		progressFunc:     func() error { return nil },
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
	dfu.progressCallback = nil
}

func (dfu *DFU) detach() error {
	_, err := dfu.dev.Control(0x21, reqDetach, 0, 0, nil)
	if err != nil {
		return wrapError("detach", err)
	}

	return nil
}

/* This commented-out code is untested and unused.

func (dfu *DFU) GetString(index int) (string, error) {
	str, err := dfu.dev.GetStringDescriptor(index)
	if err != nil {
		return str, wrapError("GetString", err)
	}

	return str, nil
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
		return time.Now(), wrapError("GetTime", err)
	}

	timeBytes, err = dfu.read(0, timeBytes) // Read BCD time bytes
	if err != nil {
		return time.Now(), wrapError("GetTime", err)
	}

	year = dfu.toDecimal(timeBytes[0])*100 + dfu.toDecimal(timeBytes[1])
	month = time.Month(dfu.toDecimal(timeBytes[2]))
	day = dfu.toDecimal(timeBytes[3])
	hours = dfu.toDecimal(timeBytes[4])
	minutes = dfu.toDecimal(timeBytes[5])
	seconds = dfu.toDecimal(timeBytes[6])

	return time.Date(year, month, day, hours, minutes, seconds, 0, location), nil
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
	rebootCmd := []byte{byte(0x91), byte(0x05)}
	_, err := dfu.dev.Control(0x21, reqWrite, 0, 0, rebootCmd)
	if err != nil {
		return wrapError("md380Reboot", err)
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		return wrapError("md380Reboot", err)
	}

	return nil
}

func (dfu *DFU) waitUntilReady() error {
	var err error

	for {
		state, err := dfu.getStatus()
		if err != nil {
			return fmt.Errorf("waitUntilReady %s", err.Error())
		}
		if state == dfuIdle {
			break
		}

		err = dfu.clearStatus()
		if err != nil {
			return fmt.Errorf("waitUntilReady %s", err.Error())
		}
	}

	return nil
}

*/

func (dfu *DFU) write(blockNumber int, bytes []byte) error {
	_, err := dfu.dev.Control(0x21, reqWrite, uint16(blockNumber), 0, bytes)
	if err != nil {
		return wrapError("write error", err)
	}

	return nil
}

func (dfu *DFU) setAddress(address int) error {
	a := byte(address)
	b := byte((address >> 8))
	c := byte((address >> 16))
	d := byte((address >> 24))
	addrCmd := []byte{0x21, a, b, c, d}

	_, err := dfu.dev.Control(0x21, reqWrite, 0, 0, addrCmd)
	if err != nil {
		return wrapError("setAddress", err)
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		return wrapError("setAddress", err)
	}

	state, err := dfu.getStatus() // this actually gets the state
	if err != nil {
		return wrapError("setAddress", err)
	}

	if state != dfuWriteIdle {
		return wrapError("setAddress", err)
	}

	err = dfu.enterDFUMode()
	if err != nil {
		return wrapError("setAddress", err)
	}

	return nil
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

func (dfu *DFU) eraseSPIFlashBlocks(address int, size int) error {
	count := (size + dfu.eraseBlockSize - 1) / dfu.eraseBlockSize
	addr := uint32(address)

	dfu.setMaxProgressCount((count + 1) * spiEraseSPIFlashBlockDelay)

	for i := 0; i < count; i++ {
		err := dfu.progressFunc()
		if err != nil {
			return err
		}
		err = dfu.eraseSPIFlashBlock(addr)
		if err != nil {
			return err
		}
		addr += uint32(dfu.eraseBlockSize)
	}

	dfu.finalProgress()

	return nil
}

func (dfu *DFU) eraseBlock(address uint32) error {
	a := byte(address)
	b := byte((address >> 8))
	c := byte((address >> 16))
	d := byte((address >> 24))
	addrCmd := []byte{0x41, a, b, c, d}

	_, err := dfu.dev.Control(0x21, reqWrite, 0, 0, addrCmd)
	if err != nil {
		return wrapError("eraseBlock", err)
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		return wrapError("eraseBlock", err)
	}
	state, err := dfu.getStatus() // this actually gets the state
	if err != nil {
		return wrapError("eraseBlock", err)
	}
	if state != dfuWriteIdle {
		return errors.New("eraseBlock: state != dfuWriteIdle")
	}

	err = dfu.enterDFUMode()
	if err != nil {
		return wrapError("eraseBlock", err)
	}

	return nil
}

func (dfu *DFU) eraseSPIFlashBlock(address uint32) error {
	addrCmd := []byte{
		byte(0x03), // SPIFLASHWRITE
		byte(address),
		byte((address >> 8)),
		byte((address >> 16)),
		byte((address >> 24)),
	}

	_, err := dfu.dev.Control(0x21, reqWrite, 1, 0, addrCmd)
	if err != nil {
		return wrapError("eraseSPIFlashBlock", err)
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		return wrapError("eraseSPIFlashBlock", err)
	}

	err = dfu.sleepMilliseconds(spiEraseSPIFlashBlockDelay)
	if err != nil {
		return err
	}

	_, err = dfu.getStatus() // this actually gets the state
	if err != nil {
		return wrapError("eraseSPIFlashBlock", err)
	}

	return nil
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
			err := dfu.sleepMilliseconds(cmd.b)
			if err != nil {
				return err
			}
			continue
		}

		err := dfu.md380Custom(cmd)
		if err != nil {
			return wrapError("md380Cmd", err)
		}
	}

	return nil
}

func (dfu *DFU) sleepMilliseconds(millis int) error {
	for i := 0; i < millis; i++ {
		err := dfu.progressFunc()
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(time.Millisecond))
	}

	return nil
}

func (dfu *DFU) md380Custom(acmd md380Cmd) error {
	cmd := []byte{byte(acmd.a), byte(acmd.b)}

	_, err := dfu.dev.Control(0x21, reqWrite, 0, 0, cmd)
	if err != nil {
		return err
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		return err
	}

	err = dfu.sleepMilliseconds(100)
	if err != nil {
		return err
	}

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

func (dfu *DFU) read(blockNumber int, bytes []byte) error {
	_, err := dfu.dev.Control(0xa1, reqRead, uint16(blockNumber), 0, bytes)
	if err != nil {
		return wrapError("read", err)
	}

	return nil
}

func (dfu *DFU) readSPIFlashTo(address, size int, iWriter io.Writer) error {
	writer := bufio.NewWriter(iWriter)
	bytes := make([]byte, dfu.blockSize)

	dfu.setMaxProgressCount(size/dfu.blockSize + 1)

	endAddress := address + size
	for addr := address; addr < endAddress; addr += dfu.blockSize {
		err := dfu.progressFunc()
		if err != nil {
			return err
		}

		remaining := endAddress - addr
		if remaining < len(bytes) {
			bytes = make([]byte, remaining)
		}

		dfu.readSPIFlash(addr, bytes)

		n, err := writer.Write(bytes)
		if err != nil {
			return wrapError("readSPIFlashTo", err)
		}

		if n != len(bytes) {
			err = errors.New("short write")
			return wrapError("readSPIFlashTo", err)
		}
	}

	writer.Flush()

	dfu.finalProgress()

	return nil
}

func (dfu *DFU) writeSPIFlashFrom(address, size int, iRdr io.Reader) error {
	rdr := bufio.NewReader(iRdr)
	buf := make([]byte, dfu.blockSize)

	flashSize, err := dfu.spiFlashSize()
	if err != nil {
		return wrapError("writeSPIFlashFrom", err)
	}

	if address+size > flashSize {
		return fmt.Errorf("writeSPIFlash: flash too small to write %d bytes at %d", size, address)
	}

	err = dfu.eraseSPIFlashBlocks(0x00000000, size)
	if err != nil {
		return wrapError("writeSPIFlashFrom", err)
	}

	err = dfu.setAddress(0x00000000)
	if err != nil {
		return wrapError("writeSPIFlashFrom", err)
	}

	_, err = dfu.getStatus()
	if err != nil {
		return wrapError("writeSPIFlashFrom", err)
	}

	dfu.setMaxProgressCount(size/dfu.blockSize + 1)

	endAddress := address + size
	for addr := address; addr < endAddress; addr += dfu.blockSize {
		err := dfu.progressFunc()
		if err != nil {
			return err
		}

		n, err := rdr.Read(buf)
		if err != nil {
			return wrapError("writeSPIFlashFrom", err)
		}

		if n < len(buf) {
			buf = make([]byte, n)
		}

		err = dfu.writeSPIFlash(addr, buf)
		if err != nil {
			return wrapError("writeSPIFlashFrom", err)
		}

		for {
			state, err := dfu.getStatus()
			if err != nil {
				return wrapError("writeSPIFlashFrom", err)
			}

			if state == dfuWriteIdle {
				break
			}
		}
	}

	dfu.finalProgress()

	return nil
}

func (dfu *DFU) DumpUsers(file *os.File) error {
	dfu.setMaxProgressCount(100)

	size, err := dfu.spiFlashSize()
	if err != nil {
		return err
	}
	if size < 16*1024*1024 {
		return fmt.Errorf("Flash is only %d bytes", size)
	}
	address := 0x100000
	size -= address

	dfu.finalProgress()
	dfu.setMaxProgressCount(size / (dfu.blockSize + 1))

	err = dfu.readSPIFlashTo(address, size, file)
	if err != nil {
		return err
	}

	dfu.finalProgress()
	return nil
}

func (dfu *DFU) DumpSPIFlash(file *os.File) error {
	dfu.setMaxProgressCount(100)

	size, err := dfu.spiFlashSize()
	if err != nil {
		return err
	}

	dfu.setMaxProgressCount(size / dfu.blockSize)

	err = dfu.readSPIFlashTo(0, size, file)
	if err != nil {
		return err
	}

	dfu.finalProgress()

	return nil
}

func (dfu *DFU) readSPIFlash(address int, bytes []byte) error {
	cmd := []byte{
		byte(0x01), // SPIFLASHREAD
		byte(address),
		byte(address >> 8),
		byte(address >> 16),
		byte(address >> 24),
	}
	_, err := dfu.dev.Control(0x21, reqWrite, 1, 0, cmd)
	if err != nil {
		return wrapError("readSPIFlash", err)
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		return wrapError("readSPIFlash", err)
	}

	_, err = dfu.getStatus() // this actually gets the state
	if err != nil {
		return wrapError("readSPIFlash", err)
	}

	err = dfu.read(1, bytes)
	if err != nil {
		return wrapError("readSPIFlash", err)
	}

	return nil
}

func (dfu *DFU) writeSPIFlash(address int, bytes []byte) error {
	size := len(bytes)
	cmd := []byte{
		byte(0x04), // SPIFLASHWRITE_NEW
		byte(address),
		byte(address >> 8),
		byte(address >> 16),
		byte(address >> 24),
		byte(size),
		byte(size >> 8),
		byte(size >> 16),
		byte(size >> 24),
	}
	cmd = append(cmd, bytes...)
	_, err := dfu.dev.Control(0x21, reqWrite, 1, 0, cmd)
	if err != nil {
		return wrapError("writeSPIFlash", err)
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		return wrapError("writeSPIFlash", err)
	}

	_, err = dfu.getStatus() // this actually gets the state
	if err != nil {
		return wrapError("writeSPIFlash", err)
	}

	return nil
}

func (dfu *DFU) getCommand() ([]byte, error) {
	cmd := make([]byte, 32)

	_, err := dfu.dev.Control(0xa1, reqRead, 0, 0, cmd)
	if err != nil {
		return nil, wrapError("getCommand", err)
	}

	_, err = dfu.getStatus()
	if err != nil {
		return nil, wrapError("getCommand", err)
	}

	return cmd, nil
}

func (dfu *DFU) getStatus() (state, error) {
	bytes := make([]byte, 6)
	_, err := dfu.dev.Control(0xa1, reqGetStatus, 0, 0, bytes)
	if err != nil {
		return 0, wrapError("getStatus", err)
	}
	state := state(bytes[4])

	debug := false
	if debug {
		status := status(bytes[0]).String()
		timeout := (((bytes[1] << 8) | bytes[2]) << 8) | bytes[3]
		discarded := bytes[5]
		fmt.Fprintln(os.Stderr, status, timeout, state, discarded)
	}

	return state, nil
}

func (dfu *DFU) clearStatus() error {
	_, err := dfu.dev.Control(0x21, reqClearStatus, 0, 0, nil)
	if err != nil {
		return wrapError("clearStatus", err)
	}

	return nil
}

func (dfu *DFU) getState() (state, error) {
	bytes := make([]byte, 1)

	_, err := dfu.dev.Control(0xa1, reqGetState, 0, 0, bytes)
	if err != nil {
		return 0, wrapError("GetStatus", err)
	}

	return state(bytes[0]), nil
}

func (dfu *DFU) abort() error {
	_, err := dfu.dev.Control(0x21, reqAbort, 0, 0, nil)
	if err != nil {
		return wrapError("ClearStatus", err)
	}

	return nil
}

func (dfu *DFU) enterDFUMode() error {
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
		stat, err := dfu.getState()
		if err != nil {
			return wrapError("enterDFUMode", err)
		}
		if stat == dfuIdle {
			break
		}
		err = actionMap[stat]()
		if err != nil {
			return wrapError("enterDFUMode", err)
		}
	}

	return nil
}

func (dfu *DFU) wait() error {
	err := dfu.sleepMilliseconds(100)
	if err != nil {
		return err
	}
	return nil
}

func (dfu *DFU) spiFlashID() (string, error) {
	bytes := make([]byte, 4)
	blockNumber := 1

	cmd := []byte{0x05} // SPIFLASHGETID
	_, err := dfu.dev.Control(0x21, reqWrite, 1, 0, cmd)
	if err != nil {
		return "", wrapError("spiFlashID", err)
	}

	_, err = dfu.getStatus() // this changes state
	if err != nil {
		return "", wrapError("spiFlashID", err)
	}

	_, err = dfu.getStatus() // this actually gets the state
	if err != nil {
		return "", wrapError("spiFlashID", err)
	}

	err = dfu.read(blockNumber, bytes)
	if err != nil {
		return "", wrapError("spiFlashID", err)
	}

	var str string

	id := int(bytes[0])<<16 | int(bytes[1])<<8 | int(bytes[2])
	switch id {
	case 0xef4018, 0x10dc01:
		str = "W25Q128FV"

	case 0xef4014:
		str = "W25Q80BL"

	case 0x70f101:
		err = fmt.Errorf("Bad LibUSB connection.  Please see the advice from N6YN at https://github.com/travisgoodspeed/md380tools/issues/186")

	default:
		err = fmt.Errorf("Unknown SPI flash: %06x, please report", id)
	}

	if err != nil {
		return "", wrapError("spiFlashID", err)
	}

	return str, nil
}

func (dfu *DFU) spiFlashSize() (int, error) {
	id, err := dfu.spiFlashID()
	if err != nil {
		return 0, err
	}

	switch id {
	case "W25Q128FV":
		return 16 * 1024 * 1024, nil
	case "W25Q80BL":
		return 1 * 1024 * 1024, nil
	}

	return 0, fmt.Errorf("bad SPI Flash ID: %s", id)
}

func (dfu *DFU) setMaxProgressCount(max int) {
	dfu.progressFunc = func() error { return nil }
	if dfu.progressCallback != nil {
		dfu.progressIncrement = MaxProgress / max
		dfu.progressCounter = 0
		dfu.progressFunc = func() error {
			dfu.progressCounter += dfu.progressIncrement
			curProgress := dfu.progressCounter
			if curProgress > MaxProgress {
				curProgress = MaxProgress
			}

			if !dfu.progressCallback(dfu.progressCounter) {
				return errors.New("")
			}

			return nil
		}
		dfu.progressCallback(dfu.progressCounter)
	}
}

func (dfu *DFU) readTo(address, offset int, size int, iWriter io.Writer) error {
	if offset%dfu.blockSize != 0 {
		return fmt.Errorf("readTo: offset is not a multiple of blockSize")
	}

	if size%dfu.blockSize != 0 {
		return fmt.Errorf("readTo: data size is not a multiple of blockSize")
	}

	blockNumber := offset / dfu.blockSize
	blockCount := size / dfu.blockSize

	writer := bufio.NewWriter(iWriter)
	bytes := make([]byte, dfu.blockSize)

	err := dfu.setAddress(address)
	if err != nil {
		return wrapError("readTo", err)
	}

	dfu.setMaxProgressCount(blockCount)

	for i := 0; i < blockCount; i++ {
		err := dfu.progressFunc()
		if err != nil {
			return err
		}

		err = dfu.read(blockNumber, bytes)
		if err != nil {
			return wrapError("readTo", err)
		}

		n, err := writer.Write(bytes)
		if err != nil {
			return wrapError("readTo", err)
		}

		if n != len(bytes) {
			err = errors.New("short write")
			return wrapError("readTo", err)
		}

		blockNumber++
	}

	dfu.finalProgress()

	return nil
}

func (dfu *DFU) writeFlashFrom(address, offset int, size int, iRdr io.Reader) error {
	blockNumber := offset / dfu.blockSize
	blockCount := (size + dfu.blockSize - 1) / dfu.blockSize
	size = blockCount * dfu.blockSize

	rdr := bufio.NewReader(iRdr)
	buf := make([]byte, dfu.blockSize)

	err := dfu.eraseBlocks(0x00000000, size)
	if err != nil {
		return wrapError("writeFlashFrom", err)
	}

	err = dfu.setAddress(0x00000000)
	if err != nil {
		return wrapError("writeFlashFrom", err)
	}

	_, err = dfu.getStatus()
	if err != nil {
		return wrapError("writeFlashFrom", err)
	}

	dfu.setMaxProgressCount(blockCount)

	for i := 0; i < blockCount; i++ {
		err := dfu.progressFunc()
		if err != nil {
			return err
		}

		n, err := rdr.Read(buf)
		if err != nil {
			return wrapError("writeFlashFrom", err)
		}

		paddingSize := len(buf) - n
		if paddingSize != 0 {
			padding := bytes.Repeat([]byte{0xff}, paddingSize)
			copy(buf[n:], padding)
		}

		err = dfu.write(blockNumber, buf)
		if err != nil {
			return wrapError("writeFlashFrom", err)
		}

		for {
			state, err := dfu.getStatus()
			if err != nil {
				return wrapError("writeFlashFrom", err)
			}

			if state == dfuWriteIdle {
				break
			}
		}
		blockNumber++
	}

	dfu.finalProgress()

	return nil
}

func (dfu *DFU) finalProgress() {
	//fmt.Fprintf(os.Stderr, "\nprogressMax %d\n", dfu.progressCounter/dfu.progressIncrement)
	if dfu.progressCallback != nil {
		dfu.progressCallback(MaxProgress)
	}
}

func (dfu *DFU) ReadCodeplug(data []byte) error {
	dfu.setMaxProgressCount(620)

	size := len(data)
	buffer := bytes.NewBuffer(data[0:0])

	err := dfu.md380Cmd([]md380Cmd{
		md380Cmd{0x91, 0x01}, // Programming Mode
		md380Cmd{0xa2, 0x02},
		md380Cmd{0xa2, 0x02},
		md380Cmd{0xa2, 0x03},
		md380Cmd{0xa2, 0x04},
		md380Cmd{0xa2, 0x07},
	})
	if err != nil {
		return wrapError("ReadCodeplug", err)
	}

	dfu.finalProgress()

	err = dfu.readTo(0, 2048, size, buffer)
	if err != nil {
		return wrapError("ReadCodeplug", err)
	}

	return nil
}

func (dfu *DFU) WriteCodeplug(data []byte) error {
	dfu.setMaxProgressCount(2750)

	if len(data)%dfu.blockSize != 0 {
		return fmt.Errorf("WriteCodeplug: codeplug data size is not a multiple of blocksize %d", dfu.blockSize)
	}

	buffer := bytes.NewBuffer(data)

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
		return wrapError("WriteCodeplug", err)
	}

	dfu.finalProgress()

	return dfu.writeFlashFrom(0, 2048, len(data), buffer)
}

func (dfu *DFU) WriteUsers(filename string) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("os.Stat: %s", err.Error())
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("WriteUsers: %s", err.Error())
	}
	defer file.Close()

	size := int(fileInfo.Size())

	return dfu.writeSPIFlashFrom(0x100000, size, file)
}

func wrapError(prefix string, err error) error {
	if err.Error() == "" {
		return err
	}
	return fmt.Errorf("%s: %s", prefix, err.Error())
}
