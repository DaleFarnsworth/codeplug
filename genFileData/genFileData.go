// Copyright 2017 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of GenFileData
//
// GenFileData is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU General Public License
// as published by the Free Software Foundation.
//
// GenFileData is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with GenFileData.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func writeTypesFile(codeFilename string, filename string) {
	infile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	codeFile, err := os.Create(codeFilename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = codeFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	varName := strings.Replace(filename, ".", "_", -1)

	fmt.Fprintf(codeFile, "var %s = []byte{", varName)

	bytes, err := ioutil.ReadAll(infile)
	if err != nil {
		log.Fatal(err)
	}

	for i, b := range bytes {
		if i%12 == 0 {
			fmt.Fprintf(codeFile, "\n\t")
		}
		fmt.Fprintf(codeFile, "%#02x,", b)
	}

	fmt.Fprintf(codeFile, "\n}\n")
}

type InsertData struct {
	LineNumber   int
	DeleteToLine int
	Filename     string
}

func main() {
	log.SetFlags(log.Lshortfile)
	codeFilename := "genFileData.code"
	linesFilename := "genFileData.lines"

	filenames := os.Args[1:]
	if len(filenames) > 0 {
		writeTypesFile(codeFilename, filenames[0])
	}

	goFilename := os.Getenv("GOFILE")
	goLineStr := os.Getenv("GOLINE")

	if goFilename == "" || goLineStr == "" {
		os.Exit(0)
	}

	linesFile, err := os.OpenFile(linesFilename, os.O_CREATE|os.O_RDWR, 0600)
	linesFile.Seek(0, io.SeekEnd)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = linesFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if len(filenames) == 0 {
		fmt.Fprintf(linesFile, "%s end\n", goLineStr)
	} else {
		fmt.Fprintf(linesFile, "%s %s\n", goLineStr, codeFilename)
		os.Exit(0)
	}

	insertDatas := []InsertData{}

	linesFile.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(linesFile)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			log.Fatal(fmt.Errorf("bad data in %s", linesFilename))
		}
		lineNumber, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}

		filename := fields[1]

		if iDataLen := len(insertDatas); iDataLen > 0 {
			insertDatas[iDataLen-1].DeleteToLine = lineNumber - 1
		}

		insertData := InsertData{lineNumber, 0, filename}
		insertDatas = append(insertDatas, insertData)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	insertFiles(goFilename, insertDatas)

	os.Remove(linesFilename)
	os.Remove(codeFilename)
}

func insertFiles(filename string, insertDatas []InsertData) {
	tmpFilename := "genFileData-tmp.go"

	tmpFile, err := os.Create(tmpFilename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = tmpFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	inScanner := bufio.NewScanner(file)
	lnum := 0
	for _, iData := range insertDatas {
		for ; lnum < iData.LineNumber && inScanner.Scan(); lnum++ {
			fmt.Fprintln(tmpFile, inScanner.Text())
		}
		if iData.Filename == "end" {
			for inScanner.Scan() {
				fmt.Fprintln(tmpFile, inScanner.Text())
			}
			break
		}

		for ; lnum < iData.DeleteToLine && inScanner.Scan(); lnum++ {
			continue
		}

		insertFile, err := os.Open(iData.Filename)
		if err != nil {
			log.Fatal(err)
		}
		defer insertFile.Close()

		insertScanner := bufio.NewScanner(insertFile)
		for insertScanner.Scan() {
			fmt.Fprintln(tmpFile, insertScanner.Text())
		}
	}

	exec.Command("gofmt", "-w", tmpFilename).Run()

	err = os.Rename(tmpFilename, filename)
	if err != nil {
		log.Fatal(err)
	}
}
