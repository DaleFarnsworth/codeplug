// Copyright 2017-2018 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of GenCodeplugInfo.
//
// GenCodeplugInfo is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU General Public License
// as published by the Free Software Foundation.
//
// GenCodeplugInfo is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with GenCodeplugInfo.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

type top struct {
	Codeplugs []*Codeplug `json:"codeplugs"`
	Records   []*Record   `json:"records"`
	Fields    []*Field    `json:"fields"`
}

type Codeplug struct {
	Models        []string `json:"models"`
	Type          string   `json:"type"`
	Ext           string   `json:"ext"`
	RdtSize       int      `json:"rdtSize"`
	HeaderSize    int      `json:"headerSize"`
	TrailerOffset int      `json:"trailerOffset"`
	TrailerSize   int      `json:"trailerSize"`
	RecordTypes   []string `json:"recordTypes"`
}

type Record struct {
	TypeName   string   `json:"typeName"`
	Type       string   `json:"type"`
	Offset     int      `json:"offset"`
	Size       int      `json:"size"`
	Max        int      `json:"max"`
	DelDesc    *DelDesc `json:"delDesc"`
	FieldTypes []string `json:"fieldTypes"`
	NamePrefix string   `json:"namePrefix"`
	Names      []string `json:"names"`
}

type DelDesc struct {
	Offset int `json:"offset"`
	Size   int `json:"size"`
	Value  int `json:"value"`
}

type Field struct {
	TypeName       string          `json:"typeName"`
	Type           string          `json:"type"`
	BitOffset      int             `json:"bitOffset"`
	BitSize        int             `json:"bitSize"`
	Max            int             `json:"max"`
	ValueType      string          `json:"valueType"`
	DefaultValue   string          `json:"defaultValue"`
	Strings        *Strings        `json:"strings"`
	Span           *Span           `json:"span"`
	IndexedStrings *IndexedStrings `json:"indexedStrings"`
	ExtOffset      int             `json:"extOffset"`
	ExtSize        int             `json:"extSize"`
	ExtIndex       int             `json:"extIndex"`
	ExtBitOffset   int             `json:"extBitOffset"`
	ListType       *string         `json:"listType"`
	EnablesIn      []*EnableIn     `json:"enables"`
	EnableIn       *EnableIn       `json:"enable"`
	EnablerType    string
	Enablers       []Enabler
	Enables        []string
}

type Strings []string

type IndexedString struct {
	Index  int    `json:"index"`
	String string `json:"string"`
}

type IndexedStrings []IndexedString

type Span struct {
	Min       int    `json:"min"`
	Max       int    `json:"max"`
	Scale     int    `json:"scale"`
	Interval  int    `json:"interval"`
	MinString string `json:"minString"`
}

type EnableIn struct {
	Value    string   `json:"value"`
	Enables  []string `json:"enables"`
	Disables []string `json:"disables"`
}

type Enabler struct {
	Value  string
	Enable bool
}

type FieldRef struct {
	RType string
	FType string
}

type ValueTypeMap map[string]int

type TemplateVars struct {
	Codeplugs        []*Codeplug
	Records          []*Record
	SortedRecords    []*Record
	Fields           []*Field
	SortedFields     []*Field
	ValueTypes       []string
	ListRecordTypes  []string
	FieldRefsMap     map[string][]FieldRef
	Sanitize         func(string) string
	RecordTypeString func(string) string
	FieldTypeString  func(string) string
}

func sanitize(s string) string {
	return strings.Title(strings.Replace(s, "-", "", -1))
}

func RecordTypeString(s string) string {
	index := strings.LastIndex(s, "_")
	if index > 0 {
		s = s[:index]
	}
	return s
}

func FieldTypeString(s string) string {
	index := strings.LastIndex(s, "_")
	if index > 0 {
		s = s[:index]
	}
	return s[2:]
}

var seenEnable = make(map[string]bool)

func doEnables(r *Record, fieldMap map[string]*Field) {
	fieldEnables := make(map[*Field]map[string]*Enabler)

	for _, fType := range r.FieldTypes {
		f := fieldMap[fType]
		if f == nil {
			fmt.Fprintf(os.Stderr, "1 found no field type: %s\n", fType)
			os.Exit(1)
		}

		enables := f.EnablesIn
		if f.EnableIn != nil {
			if len(enables) != 0 {
				fmt.Fprintf(os.Stderr, "both f.enable & f.enables found: %s\n", fType)
				os.Exit(1)
			}
			enables = []*EnableIn{f.EnableIn}
		}

		enablesMap := make(map[string]bool)

		for _, enable := range enables {
			for _, enable := range enable.Enables {
				enablesMap[enable] = true
			}
			for _, enable := range enable.Disables {
				enablesMap[enable] = true
			}

			for _, fTypeEn := range enable.Enables {
				f := fieldMap[fTypeEn]
				if f == nil {
					fmt.Fprintf(os.Stderr, "2 found no field type: %s\n", fTypeEn)
					os.Exit(1)
				}

				f.EnablerType = fType

				fieldEnable := fieldEnables[f]
				if fieldEnable == nil {
					fieldEnable = make(map[string]*Enabler)
					fieldEnables[f] = fieldEnable
				}
				enabler := fieldEnable[enable.Value]
				if enabler == nil {
					enabler = new(Enabler)
				}
				enabler.Value = enable.Value
				enabler.Enable = true
				fieldEnable[enable.Value] = enabler
				fieldEnables[f] = fieldEnable
			}

			for _, fTypeDis := range enable.Disables {
				f := fieldMap[fTypeDis]
				if f == nil {
					fmt.Fprintf(os.Stderr, "3 found no field type: %s\n", fTypeDis)
					os.Exit(1)
				}

				f.EnablerType = fType

				fieldEnable := fieldEnables[f]
				if fieldEnable == nil {
					fieldEnable = make(map[string]*Enabler)
					fieldEnables[f] = fieldEnable
				}
				enabler := fieldEnable[enable.Value]
				if enabler == nil {
					enabler = new(Enabler)
				}
				enabler.Value = enable.Value
				enabler.Enable = false
				fieldEnable[enabler.Value] = enabler
				fieldEnables[f] = fieldEnable
			}
		}

		f.Enables = make([]string, 0)
		for enable := range enablesMap {
			f.Enables = append(f.Enables, enable)
		}
		sort.Strings(f.Enables)
	}

	for f, enables := range fieldEnables {
		if f.Enablers == nil {
			f.Enablers = make([]Enabler, 0)
		}

		values := make([]string, 0)
		for value := range enables {
			values = append(values, value)
		}
		sort.Strings(values)

		for _, value := range values {
			enable := enables[value]
			seenKey := f.Type + ":" + value
			if !seenEnable[seenKey] {
				f.Enablers = append(f.Enablers, *enable)
				seenEnable[seenKey] = true
			}
		}
	}
}

func fieldRefsMap(records []*Record, fieldMap map[string]*Field) map[string][]FieldRef {
	rfMapExists := make(map[string]map[FieldRef]bool)
	rfMap := make(map[string][]FieldRef)
	for _, r := range records {
		for _, fType := range r.FieldTypes {
			f := fieldMap[fType]
			if f.ListType == nil {
				continue
			}
			rType := *f.ListType
			if rType == "" {
				continue
			}

			rtStr := RecordTypeString(rType)
			if rfMap[rtStr] == nil {
				rfMap[rtStr] = make([]FieldRef, 0)
				rfMapExists[rtStr] = make(map[FieldRef]bool)
			}
			fieldRef := FieldRef{
				RType: RecordTypeString(r.Type),
				FType: FieldTypeString(f.Type),
			}
			if !rfMapExists[rtStr][fieldRef] {
				rfMap[rtStr] = append(rfMap[rtStr], fieldRef)
			}
			rfMapExists[rtStr][fieldRef] = true
		}
	}

	return rfMap
}

func sortRecords(records []*Record) {
	recordTypes := make([]string, len(records))
	recordMap := make(map[string]*Record)
	for i, r := range records {
		recordTypes[i] = r.Type
		recordMap[r.Type] = r
	}
	sort.Strings(recordTypes)
	for i, typ := range recordTypes {
		records[i] = recordMap[typ]
	}
}

func sortFields(fields []*Field) {
	fieldTypes := make([]string, len(fields))
	fieldMap := make(map[string]*Field)
	for i, f := range fields {
		fieldTypes[i] = f.Type
		fieldMap[f.Type] = f
	}
	sort.Strings(fieldTypes)
	for i, typ := range fieldTypes {
		fields[i] = fieldMap[typ]
	}
}

func readCodeplugJson(filename string) TemplateVars {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var top top
	err = json.Unmarshal(bytes, &top)
	if err != nil {
		log.Fatal(err)
	}

	var templateVars TemplateVars

	templateVars.Codeplugs = top.Codeplugs

	fieldMap := make(map[string]*Field)
	valueTypeMap := make(map[string]int)

	sortedFields := make([]*Field, len(top.Fields))
	for i, f := range top.Fields {
		if f.Max == 0 {
			f.Max = 1
		}
		span := f.Span
		if span != nil {
			if span.MinString != "" {
				span.Min = 0
			}
		}
		existingField := fieldMap[f.Type]
		if existingField != nil {
			fmt.Fprintf(os.Stderr, "Duplicate field type %s\n", f.Type)
			os.Exit(1)
		}
		fieldMap[f.Type] = f
		valueTypeMap[f.ValueType]++
		sortedFields[i] = f
	}
	sortFields(sortedFields)
	templateVars.Fields = top.Fields
	templateVars.SortedFields = sortedFields

	sortedRecords := make([]*Record, len(top.Records))
	for i, r := range top.Records {
		if r.Max == 0 {
			r.Max = 1
		}
		doEnables(r, fieldMap)
		sortedRecords[i] = r
	}

	sortRecords(sortedRecords)
	templateVars.Records = top.Records
	templateVars.SortedRecords = sortedRecords

	valueTypes := make([]string, 0, len(valueTypeMap))
	for valueType := range valueTypeMap {
		valueTypes = append(valueTypes, valueType)
	}
	sort.Strings(valueTypes)
	templateVars.ValueTypes = valueTypes

	templateVars.Sanitize = sanitize
	templateVars.RecordTypeString = RecordTypeString
	templateVars.FieldTypeString = FieldTypeString

	fieldRefsMap := fieldRefsMap(top.Records, fieldMap)
	templateVars.FieldRefsMap = fieldRefsMap

	lrtStrings := make([]string, 0)
	for rtString := range fieldRefsMap {
		lrtStrings = append(lrtStrings, rtString)
	}
	sort.Strings(lrtStrings)
	templateVars.ListRecordTypes = lrtStrings

	return templateVars
}

func writeTypesFile(codeFilename string, filename string) {
	templateVars := readCodeplugJson(filename)

	template, err := template.ParseFiles("template")
	if err != nil {
		log.Fatal(err)
	}

	codeFile, err := os.Create(codeFilename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		fErr := codeFile.Close()
		if err == nil && fErr != nil {
			log.Fatal(fErr)
		}
	}()

	err = template.Execute(codeFile, templateVars)
	if err != nil {
		log.Fatal(err)
	}
}

type InsertData struct {
	LineNumber   int
	DeleteToLine int
	Filename     string
}

func main() {
	log.SetFlags(log.Lshortfile)
	codeFilename := "genCodeplugInfo.code"
	linesFilename := "genCodeplugInfo.lines"

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
		fErr := linesFile.Close()
		if err == nil && fErr != nil {
			log.Fatal(fErr)
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
	tmpFilename := "genCodeplugInfo-tmp.go"

	tmpFile, err := os.Create(tmpFilename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		fErr := tmpFile.Close()
		if err == nil && fErr != nil {
			log.Fatal(fErr)
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
