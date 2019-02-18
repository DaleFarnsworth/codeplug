// Copyright 2017-2019 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of DocCodeplug.
//
// DocCodeplug is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU General Public License
// as published by the Free Software Foundation.
//
// DocCodeplug is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with DocCodeplug.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func readCodeplugJson(filename string) (top *top) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &top)
	if err != nil {
		log.Fatal(err)
	}

	return top
}

func writeField(of *os.File, f *Field) {
	offset := f.BitOffset / 8
	bOffset := f.BitOffset % 8
	size := f.BitSize / 8
	if size == 0 {
		size = 1
	}
	bSize := f.BitSize

	strs := f.Strings
	span := f.Span
	iStrs := f.IndexedStrings

	fmt.Fprintf(of, "\n\tField: %s\n", f.TypeName)
	fmt.Fprintf(of, "\tOffset: 0x%06x\n", offset)
	fmt.Fprintf(of, "\tSize:     0x%04x\n", size)
	if f.ExtSize > 0 {
		extOffset := f.ExtOffset + f.ExtBitOffset/8
		fmt.Fprintf(of, "\tExtended offset: 0x%06x\n", extOffset)
		fmt.Fprintf(of, "\tExtended span:   0x%02x\n", f.ExtSize)
		fmt.Fprintf(of, "\tFirst extended index: %d\n", f.ExtIndex)
	}

	if bSize < 8 {
		str := ""
		for i := 0; i < bOffset; i++ {
			str += "-"
		}
		for i := bOffset; i < bOffset+bSize; i++ {
			str += "X"
		}
		for i := bOffset + bSize; i < 8; i++ {
			str += "-"
		}

		fmt.Fprintf(of, "\tBit Offset: %d, Bit Width: %d\t|%s|\n", bOffset, bSize, str)
	}
	if f.Max > 1 {
		fmt.Fprintf(of, "\tMax number of fields: %d\n", f.Max)
	}
	fmt.Fprintf(of, "\tValue Type: %s\n", f.ValueType)
	if strs != nil {
		fmt.Fprintf(of, "\tValue indexes into strings:\n")
		for _, s := range *strs {
			fmt.Fprintf(of, "\t\t\"%s\"\n", s)
		}
	}
	if span != nil {
		fmt.Fprintf(of, "\tRange of values:\n")
		scale := span.Scale
		if scale < 1 {
			scale = 1
		}
		step := span.Interval
		if step < 1 {
			step = 1
		}
		fmt.Fprintf(of, "\t\tminimum: %d, maximum: %d, step: %d, scale: %d\n", span.Min, span.Max, step, scale)
	}

	if iStrs != nil {
		for _, is := range *iStrs {
			fmt.Fprintf(of, "\t\t%d -> \"%s\"\n", is.Index, is.String)
		}
	}

	if f.ValueType == "onOff" {
		fmt.Fprintf(of, "\t\t0 -> on\n")
		fmt.Fprintf(of, "\t\t1 -> off\n")
	}

	if f.ValueType == "offOn" {
		fmt.Fprintf(of, "\t\t0 -> off\n")
		fmt.Fprintf(of, "\t\t1 -> on\n")
	}
}

func writeRecord(of *os.File, r *Record, top *top) {
	fmt.Fprintf(of, "\nRecord: %s\n", r.TypeName)
	fmt.Fprintf(of, "Offset: 0x%06x\n", r.Offset)
	fmt.Fprintf(of, "Size:     0x%04x\n", r.Size)
	if r.Max > 1 {
		fmt.Fprintf(of, "Max number of records: %d\n", r.Max)
	}

	fieldMap := make(map[string]*Field)
	for _, field := range top.Fields {
		fieldMap[field.Type] = field
	}

	for _, ft := range r.FieldTypes {
		writeField(of, fieldMap[ft])
	}
}

func writeDoc(cp *Codeplug, top *top) {
	outfilename := fmt.Sprintf("%s.txt", cp.Type)
	of, err := os.Create(outfilename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := of.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Fprintf(of, "Radio Type: %s\n", cp.Type)
	models := strings.Join(cp.Models, ", ")
	plural := ""
	if len(cp.Models) > 1 {
		plural = "s"
	}
	fmt.Fprintf(of, "Model name%s in header: %s\n", plural, models)
	fmt.Fprintf(of, "Codeplug size: %d\n", cp.RdtSize)
	fmt.Fprintf(of, "Header size: %d\n", cp.HeaderSize)

	recordMap := make(map[string]*Record)
	for _, record := range top.Records {
		recordMap[record.Type] = record
	}

	for _, rt := range cp.RecordTypes {
		writeRecord(of, recordMap[rt], top)
	}
}

func writeDocs(top *top) {
	for _, cp := range top.Codeplugs {
		writeDoc(cp, top)
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	infile := os.Args[1]

	top := readCodeplugJson(infile)

	writeDocs(top)
}
