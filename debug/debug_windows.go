// Copyright 2017-2018 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of Editcp.
//
// Editcp is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU General Public License
// as published by the Free Software Foundation.
//
// Editcp is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Editcp.  If not, see <http://www.gnu.org/licenses/>.

package l

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

func setStdHandle(stdhandle int32, handle syscall.Handle) error {
	r0, _, e1 := syscall.Syscall(procSetStdHandle.Addr(),
		2, uintptr(stdhandle), uintptr(handle), 0)
	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}

// redirectStderr to the file passed in
func redirectStderr(f *os.File) {
	err := setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(f.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}

	os.Stderr = f
}

func display(title string, args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "l.display panicked")
			os.Exit(1)
		}
	}()

	var mod = syscall.NewLazyDLL("user32.dll")
	var proc = mod.NewProc("MessageBoxW")
	var MB_OK = 0x00000000

	strs := make([]string, len(args))
	for i, arg := range args {
		strs[i] = fmt.Sprintf("%v", arg)
	}
	str := strings.Join(strs, " ")

	_, _, _ = proc.Call(0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(str))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(MB_OK))
}

func copyPanicString(source, dest string) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return
	}
	defer sourceFile.Close()

	var line string
	scanner := bufio.NewScanner(sourceFile)

	for scanner.Scan() {
		line = scanner.Text()

		if strings.HasPrefix(line, "panic: ") {
			break
		}
		if strings.HasPrefix(line, "fatal error: ") {
			break
		}
		line = ""
	}

	if line == "" {
		return
	}

	var destFile *os.File
	destFile, err = os.Create(dest)
	if err != nil {
		return
	}
	defer destFile.Close()

	fmt.Fprintln(destFile, line)

	for scanner.Scan() {
		fmt.Fprintln(destFile, scanner.Text())
	}
}

func init() {
	copyPanicString(StderrFilename, PreviousPanicFilename)

	stderr, err := os.Create(StderrFilename)
	if err != nil {
		return
	}

	redirectStderr(stderr)
}

func PreviousPanicString() string {
	bytes, _ := ioutil.ReadFile(PreviousPanicFilename)
	return string(bytes)
}

func RemovePreviousPanicFile() {
	os.Remove(PreviousPanicFilename)
}
