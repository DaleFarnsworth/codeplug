// Copyright 2017 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of GenLibTypes.
//
// GenLibTypes is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU General Public License
// as published by the Free Software Foundation.
//
// GenLibTypes is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with GenLibTypes.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Codeplugs struct {
	Codeplugs []*Codeplug `json:"codeplugs"`
}

type Codeplug struct {
	Name    string    `json:"name"`
	Records []*Record `json:"records"`
}

type Record struct {
	TypeName string    `json:"typeName"`
	Type     string    `json:"type"`
	Offset   int       `json:"offset"`
	Size     int       `json:"size"`
	Max      int       `json:"max"`
	DelDescs []DelDesc `json:"delDescs"`
	Fields   []*Field  `json:"fields"`
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
	Strings        *Strings        `json:"strings"`
	Span           *Span           `json:"span"`
	IndexedStrings *IndexedStrings `json:"indexedStrings"`
	Enabling       *Enabling       `json:"enabling"`
	EnablingValue  string
	Enabler        string
	Disabler       string
	ListType       *string `json:"listType"`
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

type Enabling struct {
	Value    string   `json:"value"`
	Enables  []string `json:"enables"`
	Disables []string `json:"disables"`
}

var codeplugs Codeplugs

type CodeplugMap map[string]int
type RecordMap map[string]int
type FieldMap map[string]int
type ValueTypeMap map[string]int

type TemplateVars struct {
	CodeplugMap  CodeplugMap
	RecordMap    RecordMap
	FieldMap     FieldMap
	ValueTypeMap ValueTypeMap
	Codeplugs    []*Codeplug
	Capitalize   func(string) string
	SliceAfter2  func(string) string
}

var templateVars = TemplateVars{
	RecordMap:    RecordMap{},
	FieldMap:     FieldMap{},
	ValueTypeMap: ValueTypeMap{},
	Capitalize:   strings.Title,
	SliceAfter2:  sliceAfter2,
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func unCapitalize(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

func sliceAfter2(s string) string {
	return s[2:]
}

func sortAndEnumTypes(m map[string]int) map[string]int {
	types := []string{}
	for k := range m {
		types = append(types, k)
	}
	sort.Strings(types)

	m = map[string]int{}

	for i, t := range types {
		m[t] = i
	}

	return m
}

func doEnables(r *Record, f *Field) {
	enabling := f.Enabling
	f.EnablingValue = enabling.Value
	for _, fType := range enabling.Enables {
		for _, f2 := range r.Fields {
			if f2.Type == fType {
				f2.Enabler = f.Type
			}
		}
	}
	for _, fType := range enabling.Disables {
		for _, f2 := range r.Fields {
			if f2.Type == fType {
				f2.Disabler = f.Type
			}
		}
	}
}

func readCodeplugJson(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &codeplugs)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range codeplugs.Codeplugs {
		for _, r := range c.Records {
			if r.Max == 0 {
				r.Max = 1
			}
			for _, f := range r.Fields {
				if f.Max == 0 {
					f.Max = 1
				}
				span := f.Span
				if span != nil {
					if span.MinString != "" {
						span.Min = 0
					}
				}
				if f.Enabling != nil {
					doEnables(r, f)
				}
			}
		}
	}

	templateVars.Codeplugs = codeplugs.Codeplugs

	codeplugMap := map[string]int{}
	recordMap := map[string]int{}
	fieldMap := map[string]int{}
	valueTypeMap := map[string]int{}

	for _, c := range templateVars.Codeplugs {
		codeplugMap[c.Name]++
		for _, r := range c.Records {
			recordMap[r.Type]++
			for _, f := range r.Fields {
				fieldMap[f.Type]++
				valueTypeMap[f.ValueType]++
			}
		}
	}

	templateVars.CodeplugMap = sortAndEnumTypes(codeplugMap)
	templateVars.RecordMap = sortAndEnumTypes(recordMap)
	templateVars.FieldMap = sortAndEnumTypes(fieldMap)
	templateVars.ValueTypeMap = sortAndEnumTypes(valueTypeMap)
}

func writeTypesFile(codeFilename string, inFilenames []string) {
	for _, filename := range inFilenames {
		readCodeplugJson(filename)
	}

	template, err := template.ParseFiles("template")
	if err != nil {
		log.Fatal(err)
	}

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
	codeFilename := "genTypes.code"
	linesFilename := "genTypes.lines"

	filenames := os.Args[1:]
	if len(filenames) > 0 {
		writeTypesFile(codeFilename, filenames)
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
	tmpFilename := "genTypes-tmp.go"

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
