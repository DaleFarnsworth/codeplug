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
	"log"
	"os"

	"github.com/dalefarnsworth/codeplug/dfu"
	"github.com/dalefarnsworth/codeplug/userdb"
)

func usage() {
	log.Fatalf("Usage %s read|write|dumpSPIFlash|dumpUsers|writeUsers|getUsersFile|getEuroUsersFile filename", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}
	cmd := os.Args[1]
	filename := os.Args[2]

	switch cmd {
	case "read":
		prefixes := []string{
			"Connecting to radio.",
			"Reading codeplug from radio.",
		}

		dfu, err := dfu.NewDFU(progressFunc(prefixes))
		if err != nil {
			log.Fatalf("NewDFU: %s", err.Error())
		}
		defer dfu.Close()

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("os.Create: %s", err.Error())
		}

		bytes := make([]byte, 256*1024)

		err = dfu.ReadCodeplug(bytes)
		if err != nil {
			log.Fatalf("dfu.ReadCodeplug: %s", err.Error())
		}

		bytesWritten, err := file.Write(bytes)
		if err != nil {
			log.Fatalf("file.Write: %s", err.Error())
		}
		if bytesWritten != len(bytes) {
			log.Fatalf("write to %s failed", filename)
		}

	case "write":
		prefixes := []string{
			"Preparing to dump flash",
			"Writing codeplug to radio.",
		}

		dfu, err := dfu.NewDFU(progressFunc(prefixes))
		if err != nil {
			log.Fatalf("NewDFU: %s", err.Error())
		}
		defer dfu.Close()

		fileInfo, err := os.Stat(filename)
		if err != nil {
			log.Fatalf("os.Stat: %s", err.Error())
		}

		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("os.Open: %s", err.Error())
		}

		bytes := make([]byte, fileInfo.Size())

		bytesRead, err := file.Read(bytes)
		if err != nil {
			log.Fatalf("file.Read: %s", err.Error())
		}
		if bytesRead != len(bytes) {
			log.Fatalf("Failed to read all of %s", filename)
		}

		err = dfu.WriteCodeplug(bytes)
		if err != nil {
			log.Fatalf("dfu.WriteCodeplug: %s", err.Error())
		}

	case "dumpSPIFlash":
		prefixes := []string{
			"Preparing to dump flash",
			"Dumping flash",
		}

		dfu, err := dfu.NewDFU(progressFunc(prefixes))
		if err != nil {
			log.Fatalf("NewDFU: %s", err.Error())
		}
		defer dfu.Close()

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("os.Create: %s", err.Error())
		}

		err = dfu.DumpSPIFlash(file)
		if err != nil {
			log.Fatalf(err.Error())
		}

	case "dumpUsers":
		prefixes := []string{
			"Preparing to dump users",
			fmt.Sprintf("Dumping users to %s", filename),
		}

		dfu, err := dfu.NewDFU(progressFunc(prefixes))
		if err != nil {
			log.Fatalf("NewDFU: %s", err.Error())
		}
		defer dfu.Close()

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("os.Create: %s", err.Error())
		}

		err = dfu.DumpUsers(file)
		if err != nil {
			log.Fatalf(err.Error())
		}

	case "writeUsers":
		prefixes := []string{
			"Erasing flash memory",
			"Writing users",
		}

		dfu, err := dfu.NewDFU(progressFunc(prefixes))
		if err != nil {
			log.Fatalf("NewDFU: %s", err.Error())
		}
		defer dfu.Close()

		err = dfu.WriteUsers(filename)
		fmt.Println()
		if err != nil {
			log.Fatalf("writeUsers: %s", err.Error())
		}

	case "getUsersFile":
		prefixes := []string{
			"Retrieving Users file",
		}

		euro := false
		err := userdb.WriteMD380ToolsFile(filename, euro, progressFunc(prefixes))
		if err != nil {
			log.Fatalf("getUsersFile: %s", err.Error())
		}

	case "getEuroUsersFile":
		prefixes := []string{
			"Retrieving European Users file",
		}

		euro := true
		err := userdb.WriteMD380ToolsFile(filename, euro, progressFunc(prefixes))
		if err != nil {
			log.Fatalf("getEuroUsersFile: %s", err.Error())
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
