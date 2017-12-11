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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dalefarnsworth/codeplug/stdfu"
)

const (
	MinProgress = 0
	MaxProgress = 1000000
)

const (
	controlBlock = 0
	spiBlock     = 1
	flashBlock   = 2
)

const spiEraseSPIFlashBlockDelay = 500 // milliseconds

type Dfu struct {
	stDfu             *stdfu.StDfu
	blockSize         int
	eraseBlockSize    int
	progressCallback  func(progressCounter int) bool
	progressFunc      func() error
	progressIncrement int
	progressCounter   int
}

func (dfu *Dfu) Close() {
	dfu.stDfu.Close()
	dfu.progressCallback = nil
}

/* This commented-out code is untested.

func (dfu *Dfu) toDecimal(b byte) int {
	return int(b&0xf + (b>>4)*10)
}

func (dfu *Dfu) toBCD(i int) byte {
	return byte(i/10<<4 | i%10)
}

func (dfu *Dfu) GetTime() (time.Time, error) {
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

	timeBytes, err = dfu.read(controlBlock, timeBytes) // Read BCD time bytes
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

func (dfu *Dfu) SetTime(t time.Time) error {
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

	err = dfu.write(controlBlock, bytes)
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

*/

func (dfu *Dfu) md380Reboot() error {
	err := dfu.waitUntilReady()
	if err != nil {
		return wrapError("md380Reboot", err)
	}

	stDfu := dfu.stDfu

	rebootCmd := []byte{byte(0x91), byte(0x05)}
	err = stDfu.Dnload(0, rebootCmd)
	if err != nil {
		return wrapError("md380Reboot", err)
	}

	_, err = stDfu.GetStatus() // this changes state
	if err != nil {
		return wrapError("md380Reboot", err)
	}

	return nil
}

func (dfu *Dfu) waitUntilReady() error {
	stDfu := dfu.stDfu

	for {
		dfuStatus, err := stDfu.GetStatus()
		if err != nil {
			return wrapError("waitUntilReady", err)
		}

		if dfuStatus.State == stdfu.DfuIdle {
			break
		}

		err = stDfu.ClrStatus()
		if err != nil {
			return wrapError("waitUntilReady", err)
		}
	}

	return nil
}

func (dfu *Dfu) setAddress(address int) error {
	a := byte(address)
	b := byte((address >> 8))
	c := byte((address >> 16))
	d := byte((address >> 24))
	addrCmd := []byte{0x21, a, b, c, d}

	stDfu := dfu.stDfu

	err := stDfu.Dnload(controlBlock, addrCmd)
	if err != nil {
		return wrapError("setAddress", err)
	}

	_, err = stDfu.GetStatus() // this changes state
	if err != nil {
		return wrapError("setAddress", err)
	}

	dfuStatus, err := stDfu.GetStatus() // this actually gets the state
	if err != nil {
		return wrapError("setAddress", err)
	}

	if dfuStatus.State != stdfu.DfuWriteIdle {
		return wrapError("setAddress", err)
	}

	err = dfu.enterDfuMode()
	if err != nil {
		return wrapError("setAddress", err)
	}

	return nil
}

func (dfu *Dfu) eraseBlocks(addr int, size int) error {
	count := (size + dfu.eraseBlockSize - 1) / dfu.eraseBlockSize
	for i := 0; i < count; i++ {
		err := dfu.eraseBlock(addr)
		if err != nil {
			return err
		}
		addr += dfu.eraseBlockSize
	}

	return nil
}

func (dfu *Dfu) eraseSPIFlashBlocks(addr int, size int) error {
	count := (size + dfu.eraseBlockSize - 1) / dfu.eraseBlockSize

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
		addr += dfu.eraseBlockSize
	}

	dfu.finalProgress()

	return nil
}

func (dfu *Dfu) eraseBlock(address int) error {
	addrCmd := []byte{
		0x41,
		byte(address),
		byte((address >> 8)),
		byte((address >> 16)),
		byte((address >> 24)),
	}

	stDfu := dfu.stDfu

	err := stDfu.Dnload(controlBlock, addrCmd)
	if err != nil {
		return wrapError("eraseBlock", err)
	}

	_, err = stDfu.GetStatus() // this changes state
	if err != nil {
		return wrapError("eraseBlock", err)
	}
	dfuStatus, err := stDfu.GetStatus() // this actually gets the state
	if err != nil {
		return wrapError("eraseBlock", err)
	}
	if dfuStatus.State != stdfu.DfuWriteIdle {
		return errors.New("eraseBlock: state != dfuWriteIdle")
	}

	err = dfu.enterDfuMode()
	if err != nil {
		return wrapError("eraseBlock", err)
	}

	return nil
}

func (dfu *Dfu) eraseSPIFlashBlock(address int) error {
	addrCmd := []byte{
		byte(0x03), // SPIFLASHWRITE
		byte(address),
		byte((address >> 8)),
		byte((address >> 16)),
		byte((address >> 24)),
	}

	stDfu := dfu.stDfu

	err := stDfu.Dnload(spiBlock, addrCmd)
	if err != nil {
		return wrapError("eraseSPIFlashBlock", err)
	}

	_, err = stDfu.GetStatus() // this changes state
	if err != nil {
		return wrapError("eraseSPIFlashBlock", err)
	}

	err = dfu.sleepMilliseconds(spiEraseSPIFlashBlockDelay)
	if err != nil {
		return err
	}

	_, err = stDfu.GetStatus() // this actually gets the state
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

func (dfu *Dfu) md380Cmd(commands []md380Cmd) error {
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

func (dfu *Dfu) sleepMilliseconds(millis int) error {
	for i := 0; i < millis; i++ {
		err := dfu.progressFunc()
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(time.Millisecond))
	}

	return nil
}

func (dfu *Dfu) wait() error {
	err := dfu.sleepMilliseconds(100)
	if err != nil {
		return err
	}
	return nil
}

func (dfu *Dfu) enterDfuMode() error {
	stDfu := dfu.stDfu

	actionMap := map[stdfu.State]func() error{
		stdfu.DfuWriteSync:         stDfu.Abort,
		stdfu.DfuWriteIdle:         stDfu.Abort,
		stdfu.DfuManifestSync:      stDfu.Abort,
		stdfu.DfuReadIdle:          stDfu.Abort,
		stdfu.DfuError:             stDfu.ClrStatus,
		stdfu.AppIdle:              stDfu.Detach,
		stdfu.AppDetach:            dfu.wait,
		stdfu.DfuWriteBusy:         dfu.wait,
		stdfu.DfuManifest:          stDfu.Abort,
		stdfu.DfuManifestWaitReset: dfu.wait,
		stdfu.DfuIdle:              dfu.wait,
	}

	for {
		state, err := stDfu.GetState()
		if err != nil {
			return wrapError("enterDfuMode", err)
		}
		if state == stdfu.DfuIdle {
			break
		}
		err = actionMap[state]()
		if err != nil {
			return wrapError("enterDfuMode", err)
		}
	}

	return nil
}

func (dfu *Dfu) md380Custom(acmd md380Cmd) error {
	cmd := []byte{byte(acmd.a), byte(acmd.b)}

	stDfu := dfu.stDfu

	err := stDfu.Dnload(controlBlock, cmd)
	if err != nil {
		return err
	}

	_, err = stDfu.GetStatus() // this changes state
	if err != nil {
		return err
	}

	err = dfu.sleepMilliseconds(100)
	if err != nil {
		return err
	}

	dfuStatus, err := stDfu.GetStatus() // this actually gets the state
	if err != nil {
		return err
	}

	if dfuStatus.State != stdfu.DfuWriteIdle {
		return err
	}

	err = dfu.enterDfuMode()
	if err != nil {
		return err
	}

	return nil
}

func (dfu *Dfu) readSPIFlashTo(address, size int, iWriter io.Writer) error {
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

func (dfu *Dfu) writeSPIFlashFrom(address, size int, iRdr io.Reader) error {
	rdr := bufio.NewReader(iRdr)
	buf := make([]byte, dfu.blockSize)

	flashSize, err := dfu.spiFlashSize()
	if err != nil {
		return wrapError("writeSPIFlashFrom", err)
	}

	if address+size > flashSize {
		return fmt.Errorf("writeSPIFlash: flash too small to write %d bytes at %d", size, address)
	}

	err = dfu.eraseSPIFlashBlocks(address, size)
	if err != nil {
		return wrapError("writeSPIFlashFrom", err)
	}

	err = dfu.setAddress(0x00000000)
	if err != nil {
		return wrapError("writeSPIFlashFrom", err)
	}

	stDfu := dfu.stDfu

	_, err = stDfu.GetStatus()
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
			for i := range buf[n:] {
				buf[n+i] = 0xff
			}
		}

		err = dfu.writeSPIFlash(addr, buf)
		if err != nil {
			return wrapError("writeSPIFlashFrom", err)
		}

		for {
			dfuStatus, err := stDfu.GetStatus()
			if err != nil {
				return wrapError("writeSPIFlashFrom", err)
			}

			if dfuStatus.State == stdfu.DfuWriteIdle {
				break
			}
		}
	}

	dfu.finalProgress()

	return nil
}

func (dfu *Dfu) DumpUsers(file *os.File) error {
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

	err = dfu.readSPIFlashTo(address, size, file)
	if err != nil {
		return err
	}

	return nil
}

func (dfu *Dfu) DumpSPIFlash(file *os.File) error {
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

func (dfu *Dfu) readSPIFlash(address int, bytes []byte) error {
	cmd := []byte{
		byte(0x01), // SPIFLASHREAD
		byte(address),
		byte(address >> 8),
		byte(address >> 16),
		byte(address >> 24),
	}

	stDfu := dfu.stDfu

	err := stDfu.Dnload(spiBlock, cmd)
	if err != nil {
		return wrapError("readSPIFlash", err)
	}

	_, err = stDfu.GetStatus() // this changes state
	if err != nil {
		return wrapError("readSPIFlash", err)
	}

	_, err = stDfu.GetStatus() // this actually gets the state
	if err != nil {
		return wrapError("readSPIFlash", err)
	}

	err = stDfu.Upload(spiBlock, bytes)
	if err != nil {
		return wrapError("readSPIFlash", err)
	}

	return nil
}

func (dfu *Dfu) writeSPIFlash(address int, bytes []byte) error {
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

	stDfu := dfu.stDfu

	cmd = append(cmd, bytes...)
	err := stDfu.Dnload(spiBlock, cmd)
	if err != nil {
		return wrapError("writeSPIFlash", err)
	}

	_, err = stDfu.GetStatus() // this changes state
	if err != nil {
		return wrapError("writeSPIFlash", err)
	}

	_, err = stDfu.GetStatus() // this actually gets the state
	if err != nil {
		return wrapError("writeSPIFlash", err)
	}

	return nil
}

func (dfu *Dfu) getCommand() ([]byte, error) {
	cmd := make([]byte, 32)

	stDfu := dfu.stDfu

	err := stDfu.Upload(controlBlock, cmd)
	if err != nil {
		return nil, wrapError("getCommand", err)
	}

	_, err = stDfu.GetStatus()
	if err != nil {
		return nil, wrapError("getCommand", err)
	}

	return cmd, nil
}

func (dfu *Dfu) spiFlashID() (string, error) {
	bytes := make([]byte, 4)

	stDfu := dfu.stDfu

	cmd := []byte{0x05} // SPIFLASHGETID
	err := stDfu.Dnload(spiBlock, cmd)
	if err != nil {
		return "", wrapError("spiFlashID", err)
	}

	_, err = stDfu.GetStatus() // this changes state
	if err != nil {
		return "", wrapError("spiFlashID", err)
	}

	_, err = stDfu.GetStatus() // this actually gets the state
	if err != nil {
		return "", wrapError("spiFlashID", err)
	}

	err = stDfu.Upload(spiBlock, bytes)
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

func (dfu *Dfu) spiFlashSize() (int, error) {
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

func (dfu *Dfu) setMaxProgressCount(max int) {
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

func (dfu *Dfu) readFlashTo(address, offset int, size int, iWriter io.Writer) error {
	if offset%dfu.blockSize != 0 {
		return fmt.Errorf("readFlashTo: offset is not a multiple of blockSize")
	}

	if size%dfu.blockSize != 0 {
		return fmt.Errorf("readFlashTo: data size is not a multiple of blockSize")
	}

	blockNumber := offset / dfu.blockSize
	blockCount := size / dfu.blockSize

	writer := bufio.NewWriter(iWriter)
	bytes := make([]byte, dfu.blockSize)

	err := dfu.setAddress(address)
	if err != nil {
		return wrapError("readFlashTo", err)
	}

	dfu.setMaxProgressCount(blockCount)

	stDfu := dfu.stDfu

	for i := 0; i < blockCount; i++ {
		err := dfu.progressFunc()
		if err != nil {
			return err
		}

		err = stDfu.Upload(blockNumber, bytes)
		if err != nil {
			return wrapError("readFlashTo", err)
		}

		n, err := writer.Write(bytes)
		if err != nil {
			return wrapError("readFlashTo", err)
		}

		if n != len(bytes) {
			err = errors.New("short write")
			return wrapError("readFlashTo", err)
		}

		blockNumber++
	}

	err = dfu.md380Reboot()
	if err != nil {
		return wrapError("readFlashTo", err)
	}

	dfu.finalProgress()

	return nil
}

func (dfu *Dfu) writeFlashFrom(address, offset int, size int, iRdr io.Reader) error {
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

	stDfu := dfu.stDfu

	_, err = stDfu.GetStatus()
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

		err = stDfu.Dnload(blockNumber, buf)
		if err != nil {
			return wrapError("writeFlashFrom", err)
		}

		for {
			dfuStatus, err := stDfu.GetStatus()
			if err != nil {
				return wrapError("writeFlashFrom", err)
			}

			if dfuStatus.State == stdfu.DfuWriteIdle {
				break
			}
		}
		blockNumber++
	}

	err = dfu.md380Reboot()
	if err != nil {
		return wrapError("writeFlashFrom", err)
	}

	dfu.finalProgress()

	return nil
}

type block struct {
	address int
	size    int
	end     int
}

func (dfu *Dfu) writeFirmwareFrom(iRdr io.Reader) error {
	blocks := []block{
		block{0x0800c000, 0x04000, 0x11},
		block{0x08010000, 0x10000, 0x41},
		block{0x08020000, 0x20000, 0x81},
		block{0x08040000, 0x20000, 0x81},
		block{0x08060000, 0x20000, 0x81},
		block{0x08080000, 0x20000, 0x81},
		block{0x080a0000, 0x20000, 0x81},
		block{0x080c0000, 0x20000, 0x81},
		block{0x080e0000, 0x20000, 0x81},
	}

	stDfu := dfu.stDfu

	mfg, err := stDfu.GetStringDescriptor(1)
	if err != nil {
		return wrapError("writeFirmware", err)
	}
	if mfg != "AnyRoad Technology" {
		msg := `
The radio is not in bootloader mode. Enter bootloader mode by holding
down the PTT button and the button above it while turning on the radio.
The radio's LED will blink green and red.`
		return errors.New(msg[1:])
	}

	dfuStatus, err := stDfu.GetStatus()
	if dfuStatus.State != stdfu.DfuIdle {
		return errors.New("writeFirmware: radio is not in the idle state")
	}

	err = dfu.md380Cmd([]md380Cmd{
		md380Cmd{0x91, 0x01}, // Programming Mode
		md380Cmd{0x91, 0x31},
	})

	dfu.setMaxProgressCount(len(blocks) + 1)

	totalBlocks := 0
	for _, block := range blocks {
		err := dfu.progressFunc()
		if err != nil {
			return wrapError("writeFirmware", err)
		}
		dfu.eraseBlock(block.address)

		totalBlocks += block.size / dfu.blockSize
	}

	dfu.finalProgress()

	rdr := bufio.NewReader(iRdr)
	header := "OutSecurityBin"
	headerBytes, err := rdr.Peek(len(header))
	if err != nil {
		return wrapError("writeFirmware", err)
	}
	if string(headerBytes) == header {
		_, err = rdr.Discard(0x100)
		if err != nil {
			return wrapError("writeFirmware", err)
		}
	}

	buf := make([]byte, dfu.blockSize)

	dfu.setMaxProgressCount(totalBlocks + 1)

	for _, block := range blocks {
		err = dfu.setAddress(block.address)
		if err != nil {
			return wrapError("writeFirmware", err)
		}

		blockCount := block.size / dfu.blockSize
		for blockNumber := 0; blockNumber < blockCount; blockNumber++ {
			err := dfu.progressFunc()
			if err != nil {
				return wrapError("writeFirmware", err)
			}

			err = fillBuffer(rdr, buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				return wrapError("writeFirmware", err)
			}

			err = stDfu.Dnload(flashBlock+blockNumber, buf)
			if err != nil {
				return wrapError("writeFirmware", err)
			}

			err = dfu.waitUntilReady()
			if err != nil {
				return wrapError("writeFirmware", err)
			}
		}
	}

	dfu.finalProgress()

	return nil
}

func fillBuffer(rdr io.Reader, buf []byte) error {
	for i := 0; i < len(buf); {
		n, err := rdr.Read(buf[i:])
		i += n
		if n == 0 {
			return err
		}
		if err != nil {
			if err == io.EOF {
				err = nil
				for j := range buf[i:] {
					buf[i+j] = 0xff
				}
			}
			return err
		}
	}
	return nil
}

func (dfu *Dfu) finalProgress() {
	//fmt.Fprintf(os.Stderr, "\nprogressMax %d\n", dfu.progressCounter/dfu.progressIncrement)
	if dfu.progressCallback != nil {
		dfu.progressCallback(MaxProgress)
	}
}

func (dfu *Dfu) ReadCodeplug(data []byte) error {
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

	err = dfu.readFlashTo(0, 2048, size, buffer)
	if err != nil {
		return wrapError("ReadCodeplug", err)
	}

	return nil
}

func (dfu *Dfu) WriteCodeplug(data []byte) error {
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

func (dfu *Dfu) WriteUsers(filename string) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		logFatalf("os.Stat: %s", err.Error())
	}

	file, err := os.Open(filename)
	if err != nil {
		logFatalf("WriteUsers: %s", err.Error())
	}
	defer file.Close()

	size := int(fileInfo.Size())

	return dfu.writeSPIFlashFrom(0x100000, size, file)
}

func (dfu *Dfu) WriteFirmware(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		logFatalf("WriteFirmware: %s", err.Error())
	}
	defer file.Close()

	return dfu.writeFirmwareFrom(file)
}

func wrapError(prefix string, err error) error {
	if err.Error() == "" {
		return err
	}
	return fmt.Errorf("%s: %s", prefix, err.Error())
}
