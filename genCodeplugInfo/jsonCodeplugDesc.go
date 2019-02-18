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
