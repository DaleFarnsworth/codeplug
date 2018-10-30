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
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

var StderrFilename = filepath.Join(os.TempDir(), "editcp_stderr.txt")
var PreviousPanicFilename = filepath.Join(os.TempDir(), "editcp_crash.txt")
var logBuffer strings.Builder

func init() {
	log.SetOutput(&logBuffer)
	log.SetFlags(log.Lshortfile)
}

func DisplayLog() {
	if logBuffer.Len() > 0 {
		display("Debugging Information", logBuffer.String())
	}
	logBuffer.Reset()
}

func P(v ...interface{}) {
	str := " "
	for _, arg := range v {
		str += fmt.Sprintf("%#v ", arg)
	}

	log.Output(2, str)
	DisplayLog()
}

func PrintStack() {
	log.Println(string(debug.Stack()))
	DisplayLog()
}

func Fatalf(s string, v ...interface{}) {
	log.Output(2, fmt.Sprintf(s, v...))

	log.Println(string(debug.Stack()))

	log.Output(2, fmt.Sprintf(s, v...))

	DisplayLog()
	os.Exit(1)
}

func Fatal(v ...interface{}) {
	log.Output(2, fmt.Sprintln(v...))

	log.Println(string(debug.Stack()))

	log.Output(2, fmt.Sprintln(v...))

	DisplayLog()
	os.Exit(1)
}

func Print(v ...interface{}) {
	display("Debugging Information", fmt.Sprint(v...))
}

func Println(v ...interface{}) {
	display("Debugging Information", fmt.Sprintln(v...))
}

func Printf(s string, v ...interface{}) {
	display("Debugging Information", fmt.Sprintf(s, v...))
}
