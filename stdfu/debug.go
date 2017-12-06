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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
)

func dprint(v ...interface{}) {
	skip := 1

	_, filename, line, ok := runtime.Caller(skip)

	str := ""
	if ok {
		dir := filepath.Base(filepath.Dir(filename))
		filename = filepath.Base(filename)
		filename = filepath.Join(dir, filename)
		str = fmt.Sprintf("%s:%d", filename, line)
	}

	v = append([]interface{}{str}, v...)

	fmt.Fprintln(os.Stderr, v...)
}

func printStack() {
	fmt.Fprintln(os.Stderr, "start stack trace")
	debug.PrintStack()
	fmt.Fprintln(os.Stderr)
}

/*
func display(args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			dprint("display panicked", r)
			os.Exit(1)
		}
	}()

	var mod = syscall.NewLazyDLL("user32.dll")
	var proc = mod.NewProc("MessageBoxW")
	var MB_YESNOCANCEL = 0x00000003

	str := " "
	for _, arg := range args {
		str += fmt.Sprintf("%#v\n", arg)
	}

	_, _, _ = proc.Call(0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(str))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Done Title"))),
		uintptr(MB_YESNOCANCEL))
}

func dprint(v ...interface{}) {
	skip := 1

	_, filename, line, ok := runtime.Caller(skip)

	str := ""
	if ok {
		dir := filepath.Base(filepath.Dir(filename))
		filename = filepath.Base(filename)
		filename = filepath.Join(dir, filename)
		str = fmt.Sprintf("%s:%d", filename, line)
	}

	v = append([]interface{}{str}, v...)

	display(v...)
}
*/
