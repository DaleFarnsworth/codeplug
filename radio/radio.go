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
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage %s read | write", os.Args[0])
	}
	cmd := os.Args[1]
	filename := os.Args[2]

	dfu, err := dfu.NewDFU()
	if err != nil {
		log.Fatalf("NewDFU: %s", err.Error())
	}
	defer dfu.Close()

	switch cmd {
	case "read":
		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("os.Create: %s", err.Error())
		}

		prefix := "Connecting to radio."
		bytes, err := dfu.ReadCodeplug(func(min, max, cur int) bool {
			if cur == min {
				fmt.Printf("%s... %3d%%\n", prefix, 100)
				prefix = "Reading codeplug from radio."
			}

			fmt.Printf("%s... %3d%%\r", prefix, cur*100/max)
			return true
		})
		fmt.Println()
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

		prefix := "Connecting to radio."
		err = dfu.WriteCodeplug(bytes, func(min, max, cur int) bool {
			if cur == min {
				fmt.Printf("%s... %3d%%\n", prefix, 100)
				prefix = "Writing codeplug to radio."
			}
			fmt.Printf("%s... %3d%%\r", prefix, 100)
			return true
		})
		fmt.Println()
		if err != nil {
			log.Fatalf("dfu.WriteCodeplug: %s", err.Error())
		}
	}
}
