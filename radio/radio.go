// Copyright 2017 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of Radio.
//
// Radio is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU Lesser General Public
// License as published by the Free Software Foundation.
//
// Radio is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Radio.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"

	"github.com/dalefarnsworth/codeplug/dfu"
	"github.com/dalefarnsworth/codeplug/stdfu"
	"github.com/dalefarnsworth/codeplug/userdb"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage %s <command> args\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "commands:\n")
	fmt.Fprintf(os.Stderr, "\twriteCodeplug <filename>\n")
	fmt.Fprintf(os.Stderr, "\tdumpSPIFlash <filename>\n")
	fmt.Fprintf(os.Stderr, "\tdumpUsers <filename>\n")
	fmt.Fprintf(os.Stderr, "\twriteUsers <filename>\n")
	fmt.Fprintf(os.Stderr, "\tgetUsersFile <filename>\n")
	fmt.Fprintf(os.Stderr, "\tgetEuroUsersFile filename\n")
	fmt.Fprintf(os.Stderr, "\twriteFirmware filename\n")
	fmt.Fprintf(os.Stderr, "\tstdfu\n")
	os.Exit(1)
}

func fatal(s string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			logFatal(r)
		}
	}()

	if len(os.Args) < 2 {
		usage()
	}
	cmd := os.Args[1]

	switch cmd {
	case "readCodeplug":
		if len(os.Args) != 3 {
			usage()
		}
		filename := os.Args[2]

		prefixes := []string{
			"Preparing to read codeplug",
			"Reading codeplug from radio.",
		}

		dfu, err := dfu.New(progressFunc(prefixes))
		if err != nil {
			fatal(err.Error())
		}
		defer dfu.Close()

		file, err := os.Create(filename)
		if err != nil {
			logFatalf("os.Create: %s", err.Error())
		}
		defer func() {
			err := file.Close()
			if err != nil {
				logFatalf("os.Create: %s", err.Error())
			}
			return
		}()

		bytes := make([]byte, 256*1024)

		err = dfu.ReadCodeplug(bytes)
		if err != nil {
			logFatalf("dfu.ReadCodeplug: %s", err.Error())
		}

		bytesWritten, err := file.Write(bytes)
		if err != nil {
			logFatalf("file.Write: %s", err.Error())
		}
		if bytesWritten != len(bytes) {
			logFatalf("write to %s failed", filename)
		}

	case "writeCodeplug":
		if len(os.Args) != 3 {
			usage()
		}
		filename := os.Args[2]

		prefixes := []string{
			"Preparing to write codeplug",
			"Writing codeplug to radio.",
		}

		dfu, err := dfu.New(progressFunc(prefixes))
		if err != nil {
			fatal(err.Error())
		}
		defer dfu.Close()

		fileInfo, err := os.Stat(filename)
		if err != nil {
			logFatalf("os.Stat: %s", err.Error())
		}

		file, err := os.Open(filename)
		if err != nil {
			logFatalf("os.Open: %s", err.Error())
		}

		bytes := make([]byte, fileInfo.Size())

		bytesRead, err := file.Read(bytes)
		if err != nil {
			logFatalf("file.Read: %s", err.Error())
		}
		if bytesRead != len(bytes) {
			logFatalf("Failed to read all of %s", filename)
		}

		err = dfu.WriteCodeplug(bytes)
		if err != nil {
			logFatalf("dfu.WriteCodeplug: %s", err.Error())
		}

	case "dumpSPIFlash":
		if len(os.Args) != 3 {
			usage()
		}
		filename := os.Args[2]

		prefixes := []string{
			"Preparing to dump flash",
			"Dumping flash",
		}

		dfu, err := dfu.New(progressFunc(prefixes))
		if err != nil {
			fatal(err.Error())
		}
		defer dfu.Close()

		file, err := os.Create(filename)
		if err != nil {
			fatal(err.Error())
		}
		defer func() {
			err := file.Close()
			if err != nil {
				fatal(err.Error())
			}
			return
		}()

		err = dfu.DumpSPIFlash(file)
		if err != nil {
			fatal(err.Error())
		}

	case "dumpUsers":
		if len(os.Args) != 3 {
			usage()
		}
		filename := os.Args[2]

		prefixes := []string{
			"Preparing to dump users",
			fmt.Sprintf("Dumping users to %s", filename),
		}

		dfu, err := dfu.New(progressFunc(prefixes))
		if err != nil {
			fatal(err.Error())
		}
		defer dfu.Close()

		file, err := os.Create(filename)
		if err != nil {
			logFatalf("os.Create: %s", err.Error())
		}
		defer func() {
			err := file.Close()
			if err != nil {
				logFatalf("os.Create: %s", err.Error())
			}
			return
		}()

		err = dfu.DumpUsers(file)
		if err != nil {
			logFatalf(err.Error())
		}

	case "writeUsers":
		if len(os.Args) != 3 {
			usage()
		}
		filename := os.Args[2]

		prefixes := []string{
			"Erasing flash memory",
			"Writing users",
		}

		dfu, err := dfu.New(progressFunc(prefixes))
		if err != nil {
			fatal(err.Error())
		}
		defer dfu.Close()

		err = dfu.WriteUsers(filename)
		fmt.Println()
		if err != nil {
			logFatalf("writeUsers: %s", err.Error())
		}

	case "getUsersFile":
		if len(os.Args) != 3 {
			usage()
		}
		filename := os.Args[2]

		prefixes := []string{
			"Retrieving Users file",
		}

		euro := false
		err := userdb.WriteMD380ToolsFile(filename, euro, progressFunc(prefixes))
		if err != nil {
			logFatalf("getUsersFile: %s", err.Error())
		}

	case "getEuroUsersFile":
		if len(os.Args) != 3 {
			usage()
		}
		filename := os.Args[2]

		prefixes := []string{
			"Retrieving European Users file",
		}

		euro := true
		err := userdb.WriteMD380ToolsFile(filename, euro, progressFunc(prefixes))
		if err != nil {
			logFatalf("getEuroUsersFile: %s", err.Error())
		}

	case "writeFirmware":
		if len(os.Args) != 3 {
			usage()
		}
		filename := os.Args[2]

		prefixes := []string{
			"Erasing flash memory",
			"Writing firmware",
		}

		dfu, err := dfu.New(progressFunc(prefixes))
		if err != nil {
			fatal(err.Error())
		}
		defer dfu.Close()

		err = dfu.WriteFirmware(filename)
		fmt.Println()
		if err != nil {
			logFatalf("writeFirmware: %s", err.Error())
		}
	case "stdfu":
		if len(os.Args) != 2 {
			usage()
		}

		dfu, err := stdfu.New()
		if err != nil {
			logFatalf("writeFirmware: %s", err.Error())
		}

	default:
		usage()
	}
}

func progressFunc(aPrefixes []string) func(cur int) bool {
	var prefixes []string
	if aPrefixes != nil {
		prefixes = aPrefixes
	}
	prefixIndex := 0
	prefix := prefixes[prefixIndex]
	maxProgress := userdb.MaxProgress
	return func(cur int) bool {
		if cur == 0 {
			if prefixIndex != 0 {
				fmt.Println()
			}
			prefix = prefixes[prefixIndex]
			prefixIndex++
		}
		fmt.Printf("%s... %3d%%\r", prefix, cur*100/maxProgress)
		return true
	}
}
