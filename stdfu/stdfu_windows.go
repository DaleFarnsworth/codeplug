// Copyright 2017 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of Stdfu.
//
// Stdfu is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU General Public License
// as published by the Free Software Foundation.
//
// Stdfu is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Stdfu.  If not, see <http://www.gnu.org/licenses/>.

package stdfu

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type StDfu struct {
	handle uintptr
}

type spDeviceInterfaceData struct {
	cbSize             uint32
	InterfaceClassGuid UUID
	Flags              uint32
	Reserved           uint
}

const (
	StDfuErrorOffset = 0x12340000

	StDfuNoError      = StDfuErrorOffset
	StDfuMemory       = (StDfuErrorOffset + 1)
	StDfuBadParameter = (StDfuErrorOffset + 2)

	StDfuNotImplemented  = (StDfuErrorOffset + 3)
	StDfuEnumFinished    = (StDfuErrorOffset + 4)
	StDfuOpenDriverError = (StDfuErrorOffset + 5)

	StDfuErrorDescriptorBuilding = (StDfuErrorOffset + 6)
	StDfuPipeCreationError       = (StDfuErrorOffset + 7)
	StDfuPipeResetError          = (StDfuErrorOffset + 8)
	StDfuPipeAbortError          = (StDfuErrorOffset + 9)
	StDfuStringDescriptorError   = (StDfuErrorOffset + 0xA)

	StDfuDriverIsClosed      = (StDfuErrorOffset + 0xB)
	StDfuVendorRqPB          = (StDfuErrorOffset + 0xC)
	StDfuErrorWhileReading   = (StDfuErrorOffset + 0xD)
	StDfuErrorBeforeReading  = (StDfuErrorOffset + 0xE)
	StDfuErrorWhileWriting   = (StDfuErrorOffset + 0xF)
	StDfuErrorBeforeWriting  = (StDfuErrorOffset + 0x10)
	StDfuDeviceResetError    = (StDfuErrorOffset + 0x11)
	StDfuCantUseUnplugEvent  = (StDfuErrorOffset + 0x12)
	StDfuIncorrectBufferSize = (StDfuErrorOffset + 0x13)
	StDfuDescriptorNotFound  = (StDfuErrorOffset + 0x14)
	StDfuPipesAreClosed      = (StDfuErrorOffset + 0x15)
	StDfuPipesAreOpen        = (StDfuErrorOffset + 0x16)

	StDfuTimeoutWaitingForReset = (StDfuErrorOffset + 0x17)

	StateIdle                 = 0x00
	StateDetach               = 0x01
	StateDfuIdle              = 0x02
	StateDfuDownloadSync      = 0x03
	StateDfuDownloadBusy      = 0x04
	StateDfuDownloaDIdle      = 0x05
	StateDfuManifestSync      = 0x06
	StateDfuManifest          = 0x07
	StateDfuManifestWaitReset = 0x08
	StateDfuUploadIdle        = 0x09
	StateDfuError             = 0x0A

	StateDfuUploadSync = 0x91
	StateDfuUploadBusy = 0x92

	StatusOK            = 0x00
	StatusErrTarget     = 0x01
	StatusErrFile       = 0x02
	StatusErrWrite      = 0x03
	StatusErrErase      = 0x04
	StatusErrCheckerase = 0x05
	StatusErrProg       = 0x06
	StatusErrVerify     = 0x07
	StatusErrAddress    = 0x08
	StatusErrNotDone    = 0x09
	StatusErrFirmware   = 0x0A
	StatusErrVendor     = 0x0B
	StatusErrUSBR       = 0x0C
	StatusErrPOR        = 0x0D
	StatusErrUnknown    = 0x0E
	StatusErrStalledPkt = 0x0F

	AttrDnloadCapable         = 0x01
	AttrUploadCapable         = 0x02
	AttrManifestationTolerant = 0x04
	AttrWillDetach            = 0x08
	AttrSTCanAccelerate       = 0x80
)

var dfuErrors = []error{
	errors.New("no error"),
	errors.New("memory error"),
	errors.New("bad parameter"),
	errors.New("not implemented"),
	errors.New("enumeration finished"),
	errors.New("open driver error"),
	errors.New("descriptor building error"),
	errors.New("pipe creation error"),
	errors.New("pipe reset error"),
	errors.New("pipe abort error"),
	errors.New("string descriptor error"),
	errors.New("driver is closed"),
	errors.New("vendor request pb"),
	errors.New("error while reading"),
	errors.New("error before reading"),
	errors.New("error while writing"),
	errors.New("error before writing"),
	errors.New("device reset error"),
	errors.New("can't use unplug event"),
	errors.New("incorrect buffer size"),
	errors.New("descriptor not found"),
	errors.New("pipes are closed"),
	errors.New("pipes are open"),
	errors.New("timeout waiting for reset"),
}

func errorFromErrno(errno uintptr) error {
	var err error

	if errno != StDfuNoError {
		errno -= StDfuErrorOffset
		if errno < 0 || int(errno) >= len(dfuErrors) {
			err = fmt.Errorf("error %d", errno)
		}

		err = dfuErrors[errno]
	}

	return err
}

var setupapi = &windows.LazyDLL{Name: "setupapi.dll", System: true}

var setupDiGetClassDevsW = setupapi.NewProc("SetupDiGetClassDevsW")
var setupDiEnumDeviceInterfaces = setupapi.NewProc("SetupDiEnumDeviceInterfaces")
var setupDiGetDeviceInterfaceDetailW = setupapi.NewProc("SetupDiGetDeviceInterfaceDetailW")
var setupDiGetDeviceRegistryPropertyW = setupapi.NewProc("SetupDiGetDeviceRegistryPropertyW")
var setupDiDestroyDeviceInfoList = setupapi.NewProc("SetupDiDestroyDeviceInfoList")

var stdfuDLL = &windows.LazyDLL{Name: "STDFU.dll", System: false}

var stdfuAbort = stdfuDLL.NewProc("STDFU_Abort")
var stdfuClose = stdfuDLL.NewProc("STDFU_Close")
var stdfuClrStatus = stdfuDLL.NewProc("STDFU_Clrstatus")
var stdfuDetach = stdfuDLL.NewProc("STDFU_Detach")
var stdfuDnload = stdfuDLL.NewProc("STDFU_Dnload")
var stdfuGetstate = stdfuDLL.NewProc("STDFU_Getstate")
var stdfuGetstatus = stdfuDLL.NewProc("STDFU_Getstatus")
var stdfuGetStringDescriptor = stdfuDLL.NewProc("STDFU_GetStringDescriptor")
var stdfuOpen = stdfuDLL.NewProc("STDFU_Open")
var stdfuSelectCurrentConfiguration = stdfuDLL.NewProc("STDFU_SelectCurrentConfiguration")
var stdfuUpload = stdfuDLL.NewProc("STDFU_Upload")

func (dfu *StDfu) Abort() error {
	errno, _, _ := stdfuAbort.Call(uintptr(unsafe.Pointer(&dfu.handle)))
	err := errorFromErrno(errno)

	return err
}

func (dfu *StDfu) Close() error {
	errno, _, _ := stdfuClose.Call(uintptr(unsafe.Pointer(&dfu.handle)))
	err := errorFromErrno(errno)

	return err
}

func (dfu *StDfu) ClrStatus() error {
	errno, _, _ := stdfuClrStatus.Call(uintptr(unsafe.Pointer(&dfu.handle)))
	err := errorFromErrno(errno)

	return err
}

func (dfu *StDfu) Detach() error {
	errno, _, _ := stdfuDetach.Call(uintptr(unsafe.Pointer(&dfu.handle)))
	err := errorFromErrno(errno)

	return err
}

func (dfu *StDfu) Dnload(blockNumber int, buffer []byte) error {
	errno, _, _ := stdfuDnload.Call(uintptr(unsafe.Pointer(&dfu.handle)), ((*reflect.SliceHeader)(unsafe.Pointer(&buffer))).Data, uintptr(len(buffer)), uintptr(blockNumber))
	err := errorFromErrno(errno)

	return err
}

func (dfu *StDfu) GetState() (State, error) {
	var state8 uint8

	errno, _, _ := stdfuGetstate.Call(uintptr(unsafe.Pointer(&dfu.handle)), uintptr(unsafe.Pointer(&state8)))
	err := errorFromErrno(errno)

	return State(state8), err
}

func (dfu *StDfu) GetStatus() (DfuStatus, error) {
	var bytes [6]byte

	errno, _, _ := stdfuGetstatus.Call(uintptr(unsafe.Pointer(&dfu.handle)), uintptr(unsafe.Pointer(&bytes)))
	err := errorFromErrno(errno)

	dfuStatus := DfuStatus{
		Status:      Status(bytes[0]),
		PollTimeout: int(bytes[1]<<16 | bytes[2]<<8 | bytes[3]),
		State:       State(bytes[4]),
		IString:     int(bytes[5]),
	}

	return dfuStatus, err
}

func (dfu *StDfu) GetStringDescriptor(index int) (string, error) {
	var buf [512]byte

	errno, _, _ := stdfuGetStringDescriptor.Call(uintptr(unsafe.Pointer(&dfu.handle)), uintptr(index), uintptr(unsafe.Pointer(&buf)), uintptr(512))
	err := errorFromErrno(errno)

	len := bytes.IndexByte(buf[:], 0)
	if len < 0 {
		len = 512
	}

	return string(buf[:len]), err
}

func (dfu *StDfu) Upload(blockNumber int, buffer []byte) error {
	errno, _, _ := stdfuUpload.Call(uintptr(unsafe.Pointer(&dfu.handle)), ((*reflect.SliceHeader)(unsafe.Pointer(&buffer))).Data, uintptr(len(buffer)), uintptr(blockNumber))
	err := errorFromErrno(errno)

	return err
}

type UUID struct {
	a uint32
	b uint16
	c uint16
	d [8]byte
}

func New() (*StDfu, error) {
	devUUID := &UUID{
		a: 0x3fe809ab,
		b: 0xfb91,
		c: 0x4cb5,
		d: [8]byte{0xa6, 0x43, 0x69, 0x67, 0x0d, 0x52, 0x36, 0x6e},
	}

	stDfu := new(StDfu)
	devIndex := 0

	hdev, _, err := setupDiGetClassDevsW.Call(
		uintptr(unsafe.Pointer(devUUID)),
		0,
		0,
		2|16, // DIGCF_PRESENT|DIGCF_DEVICEINTERFACE
	)
	if windows.Handle(hdev) == windows.InvalidHandle {
		if err == nil {
			err = syscall.EINVAL
		}
		return nil, fmt.Errorf("setupDiGetClassDevsW %s", err.Error())
	}
	defer syscall.Syscall(setupDiDestroyDeviceInfoList.Addr(), 1, hdev, 0, 0)
	var did spDeviceInterfaceData
	did.cbSize = uint32(unsafe.Sizeof(did))
	r0, r1, err := setupDiEnumDeviceInterfaces.Call(
		hdev,
		0,
		uintptr(unsafe.Pointer(devUUID)),
		uintptr(devIndex),
		uintptr(unsafe.Pointer(&did)),
	)
	if r0 == 0 { // false
		return nil, ErrDevNotFound
	}

	r0, r1, err = setupDiEnumDeviceInterfaces.Call(
		hdev,
		1,
		uintptr(unsafe.Pointer(devUUID)),
		uintptr(devIndex),
		uintptr(unsafe.Pointer(&did)),
	)
	if r0 == 1 { // true
		return nil, ErrMultipleDevs
	}

	var cbRequired uint32
	r0, r1, err = setupDiGetDeviceInterfaceDetailW.Call(
		hdev,
		uintptr(unsafe.Pointer(&did)),
		0,
		0,
		uintptr(unsafe.Pointer(&cbRequired)),
		0,
	)

	if r0 != 0 || r1 != 0x7a {
		return nil, fmt.Errorf("setupDiGetDeviceInterfaceDetailW %s", err.Error())
	}
	// The struct with ANYSIZE_ARRAY of utf16 in it is crazy.
	// So... let's emulate it with array of uint16 ;-D.
	// Keep in mind that the first two elements are actually cbSize.
	didd := make([]uint16, cbRequired/2-1)
	cbSize := (*uint32)(unsafe.Pointer(&didd[0]))
	if unsafe.Sizeof(uint(0)) == 8 {
		*cbSize = 8
	} else {
		*cbSize = 6
	}

	devInfoData := make([]uint16, len(didd))
	copy(devInfoData, didd)

	r0, r1, err = setupDiGetDeviceInterfaceDetailW.Call(
		hdev,
		uintptr(unsafe.Pointer(&did)),
		uintptr(unsafe.Pointer(&didd[0])),
		uintptr(cbRequired),
		0,
		uintptr(unsafe.Pointer(&devInfoData)),
	)
	if r0 != 0 {
		return nil, fmt.Errorf("setupDiGetDeviceInterfaceDetailW %s", err.Error())
	}
	devicePath := didd[2:]
	devicePathString := windows.UTF16ToString(devicePath)
	devicePathBytePtr, err := windows.BytePtrFromString(devicePathString)
	if err != nil {
		return nil, err
	}

	errno, _, _ := stdfuOpen.Call(uintptr(unsafe.Pointer(devicePathBytePtr)), uintptr(unsafe.Pointer(&stDfu.handle)))
	err = errorFromErrno(errno)
	if err != nil {
		return nil, fmt.Errorf("stdfuOpen %s", err.Error())
	}

	return stDfu, nil
}
