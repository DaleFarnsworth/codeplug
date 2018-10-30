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
	"fmt"
	"os"
)

func redirectStderr(f *os.File) {
	/*
		err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
		if err != nil {
			log.Fatalf("Failed to redirect stderr to file: %v", err)
		}
	*/
}

func display(title string, v ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "l.display panicked")
			os.Exit(1)
		}
	}()

	fmt.Fprint(os.Stderr, v...)
}

func PreviousPanicString() string {
	return ""
}

func RemovePreviousPanicFile() {
}
