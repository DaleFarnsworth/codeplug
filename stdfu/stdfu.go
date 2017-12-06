package stdfu

import (
	"fmt"
)

var ErrDevNotFound = fmt.Errorf("No devices found")
var ErrMultipleDevs = fmt.Errorf("Multiple devices found")

type DfuStatus struct {
	Status      Status
	PollTimeout int
	State       State
	IString     int
}

type State int

const (
	AppIdle State = iota
	AppDetach
	DfuIdle
	DfuWriteSync
	DfuWriteBusy
	DfuWriteIdle
	DfuManifestSync
	DfuManifest
	DfuManifestWaitReset
	DfuReadIdle
	DfuError
)

var stateStrings = map[State]string{
	AppIdle:              "appIdle",
	AppDetach:            "appDetach",
	DfuIdle:              "dfuIdle",
	DfuWriteSync:         "dfuWriteSync",
	DfuWriteBusy:         "dfuWriteBusy",
	DfuWriteIdle:         "dfuWriteIdle",
	DfuManifestSync:      "dfuManifestSync",
	DfuManifest:          "dfuManifest",
	DfuManifestWaitReset: "dfuManifestWaitReset",
	DfuReadIdle:          "dfuReadIdle",
	DfuError:             "dfuError",
}

func (s State) String() string {
	return stateStrings[s]
}

type Status int

const (
	StatusOk Status = iota
	ErrTarget
	ErrFile
	ErrWrite
	ErrErase
	ErrCheckErased
	ErrProgram
	ErrVerify
	ErrAddress
	ErrNotDone
	ErrFirmware
	ErrVendor
	ErrUsbR
	ErrPOR
	ErrUnknown
	ErrStalledPkt
)

var statusStrings = map[Status]string{
	StatusOk:       "ok",
	ErrTarget:      "errTarget",
	ErrFile:        "errFile",
	ErrWrite:       "errWrite",
	ErrErase:       "errErase",
	ErrCheckErased: "errCheckErased",
	ErrProgram:     "errProgram",
	ErrVerify:      "errVerify",
	ErrAddress:     "errAddress",
	ErrNotDone:     "errNotDone",
	ErrFirmware:    "errFirmware",
	ErrVendor:      "errVendor",
	ErrUsbR:        "errUsbR",
	ErrPOR:         "errPOR",
	ErrUnknown:     "errUnknown",
	ErrStalledPkt:  "errStalledPkt",
}

func (s Status) String() string {
	return statusStrings[s]
}
