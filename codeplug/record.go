// Copyright 2017 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of Codeplug.
//
// Codeplug is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU Lesser General Public
// License as published by the Free Software Foundation.
//
// Codeplug is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Codeplug.  If not, see <http://www.gnu.org/licenses/>.

// Package codeplug implements access to MD380-style codeplug files.
// It can read/update/write both .rdt files and .bin files.
package codeplug

import (
	"fmt"
	"log"
	"sort"
)

// A Record represents a record within a Codeplug.
type Record struct {
	*rDesc
	fDesc  *map[FieldType]*fDesc
	rIndex int
}

// An rDesc contains a record type's dynamic information.
type rDesc struct {
	*rInfo
	codeplug        *Codeplug
	records         []*Record
	cachedListNames *[]string
}

// An rInfo contains a record type's static information.
type rInfo struct {
	rType         RecordType
	typeName      string
	max           int
	offset        int
	size          int
	delDescs      []delDesc
	fInfos        []fInfo
	nameFieldType FieldType
}

// A RecordType represents a record's type
type RecordType string

// A delDesc contains the location and value of a deleted record indicator.
type delDesc struct {
	offset uint8
	size   uint8
	value  byte
}

// deleteRecord removes record i from the slice: records
func deleteRecord(records *[]*Record, i int) {
	copy((*records)[i:], (*records)[i+1:])
	(*records)[len(*records)-1] = nil
	*records = (*records)[:len(*records)-1]
}

// recordBytes returns the byte slice making up the record.
func (r *Record) recordBytes() []byte {
	recordBytes := make([]byte, r.size)
	r.store(recordBytes)

	return recordBytes
}

// loadRecords loads all the records in rDesc from the codeplug's file.
func (rd *rDesc) loadRecords(cpBytes []byte) {
	records := rd.records
	if len(records) > 0 {
		records = records[0:cap(records)]
	} else {
		records = make([]*Record, rd.max)
	}

	cp := rd.codeplug
	length := 0
	for rIndex := range records {
		if !rd.recordIsDeleted(rIndex, cpBytes) {
			offset := rd.offset + rd.size*rIndex
			recordBytes := cpBytes[offset : offset+rd.size]
			r := records[rIndex]
			if r == nil {
				r = cp.newRecord(rd.rType, rIndex)
			}

			r.load(recordBytes)

			records[length] = r
			length++
		}
	}
	rd.records = records[:length]
}

// newField creates and returns the address of a new field of the given type.
func (r *Record) NewField(fType FieldType) *Field {
	f := new(Field)
	fd := (*r.fDesc)[fType]
	if fd == nil {
		fd = f.fDesc
		(*r.fDesc)[fType] = fd
		fd.record = r
	}
	f.fDesc = fd
	f.value = newValue(fd.valueType)

	return f
}

// addField adds the given field to the record.
func (r *Record) addField(f *Field) error {
	if len(f.fields) >= f.max {
		return fmt.Errorf("too many fields: ", string(f.fType))
	}

	fd := (*r.fDesc)[f.fType]
	f.fIndex = len(fd.fields)
	fd.fields = append(fd.fields, f)

	return nil
}

// load replaces the record's contents with those corresponding
// to the byte slice.
func (r *Record) load(recordBytes []byte) {
	ri := r.rDesc.rInfo

	for i := range ri.fInfos {
		fi := &ri.fInfos[i]
		if fi.max == 0 {
			fi.max = 1
		}
		fi.rInfo = ri
		fd := &fDesc{fInfo: fi}
		(*r.fDesc)[fi.fType] = fd
		if fi.valueType == VtName {
			ri.nameFieldType = fi.fType
		}
		fd.record = r
	}

	for _, fd := range *r.fDesc {
		fi := fd.fInfo
		fields := fd.fields
		if len(fields) > 0 {
			fields = fields[0:cap(fields)]
		} else {
			fields = make([]*Field, fi.max)
		}

		length := 0
		for fIndex := range fields {
			if !fd.fieldDeleted(fIndex, recordBytes) {
				f := fields[fIndex]
				if f == nil {
					f = &Field{}
				}

				f.fDesc = fd
				f.fIndex = fIndex
				f.value = newValue(fi.valueType)

				f.load(recordBytes)

				span := f.span
				if span != nil {
					if span.scale == 0 {
						span.scale = 1
					}
					if span.interval == 0 {
						span.interval = 1
					}
				}

				fields[length] = f
				length++
			}
		}

		fd.fields = fields[:length]
	}
}

// valid returns nil if all fields in the record are valid.
func (r *Record) valid() error {
	errStr := ""
	for _, fType := range r.FieldTypes() {
		for _, f := range r.Fields(fType) {
			if err := f.valid(); err != nil {
				errStr += f.FullTypeName() + ": " + err.Error() + "\n"
			}
		}
	}

	if errStr != "" {
		return fmt.Errorf("%s", errStr)
	}

	return nil
}

// stores stores all all fields of the record into the given byte slice.
func (r *Record) store(recordBytes []byte) {
	for _, fd := range *r.fDesc {
		for fIndex := 0; fIndex < fd.max; fIndex++ {
			if fIndex < len(fd.fields) {
				fd.fields[fIndex].store(recordBytes)
			} else {
				fd.deleteField(fIndex, recordBytes)
			}
		}
	}
}

// FieldTypes return all valid FieldTypes for the record.
func (r *Record) FieldTypes() []FieldType {
	fds := *r.fDesc

	strs := make([]string, 0, len(fds))

	for fType := range fds {
		strs = append(strs, string(fType))
	}
	sort.Strings(strs)

	fTypes := make([]FieldType, len(strs))
	for i, str := range strs {
		fTypes[i] = FieldType(str)
	}

	return fTypes
}

// Fields returns a slice of all fields of the given type in the record.
func (r *Record) Fields(fType FieldType) []*Field {
	return (*r.fDesc)[fType].fields
}

// Field returns the first field of the given type in the record.
func (r *Record) Field(fType FieldType) *Field {
	fields := r.Fields(fType)
	if len(fields) == 0 {
		return nil
	}
	return fields[0]
}

// MaxFields returns the maximum number of fields of the given type for
// record.
func (r *Record) MaxFields(fType FieldType) int {
	return (*r.fDesc)[fType].max
}

// Type returns the record's type.
func (r *Record) Type() RecordType {
	return r.rType
}

// TypeName returns the name of the record's type.
func (r *Record) TypeName() string {
	return r.typeName
}

// Index returns the slice index of the record.
func (r *Record) Index() int {
	return r.rIndex
}

// Index set the index of the record.
func (r *Record) SetIndex(index int) {
	r.rIndex = index
}

// Codeplug returns the codeplug for the record.
func (r *Record) Codeplug() *Codeplug {
	return r.codeplug
}

// NameField returns the field containing the record's name.
func (r *Record) NameField() *Field {
	if (*r.fDesc)[r.nameFieldType] == nil {
		return nil
	}
	return r.Field(r.nameFieldType)
}

// NameType returns the fieldtype containing the record's name field.
func (r *Record) NameFieldType() FieldType {
	return r.nameFieldType
}

// Name returns the record's name.
func (r *Record) Name() string {
	if r.NameField() != nil {
		return r.NameField().String()
	}
	return ""
}

// makeNameUnique renames the record to make it different than all of
// the passed names.
func (r *Record) makeNameUnique(names []string) error {
	if len(names) == 0 {
		return nil
	}

	nameField := r.NameField()
	baseName := nameField.String()

	if !stringInSlice(baseName, names) {
		return nil
	}

	maxNameLen := nameField.bitSize / 16
	if len(baseName) >= maxNameLen {
		baseName = baseName[:maxNameLen-2]
	}

	runes := []rune(baseName)
	if runes[len(runes)-2] == '.' {
		runes = runes[:len(runes)-2]
		baseName = string(runes)
	}

	var newName string
	suffixRunes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for _, c := range suffixRunes {
		newName = baseName + "." + string(c)
		if !stringInSlice(newName, names) {
			nameField.value = newValue(nameField.ValueType())
			nameField.value.SetString(nameField, newName)
			return nil
		}
	}

	return fmt.Errorf("too many record copies")
}

// ListNames returns a slice of the names of all records in the rDesc.
func (rd *rDesc) ListNames() *[]string {
	if rd.cachedListNames == nil {
		names := make([]string, len(rd.records))
		for i, r := range rd.records {
			name := r.Name()
			if name == "" {
				return nil
			}
			names[i] = name
		}

		rd.cachedListNames = &names
	}

	return rd.cachedListNames
}

// recordIsDeleted returns true if the record at rIndex is deleted.
func (rd *rDesc) recordIsDeleted(rIndex int, cpBytes []byte) bool {
nextDelDesc:
	for _, dd := range rd.delDescs {
		offset := rd.offset + rIndex*rd.size + int(dd.offset)

		for i := 0; i < int(dd.size); i++ {
			if cpBytes[offset+i] != dd.value {
				continue nextDelDesc
			}
		}
		return true
	}

	return false
}

// deleteRecord marks the record at rIndex as deleted.
func (rd *rDesc) deleteRecord(rIndex int, cpBytes []byte) {
	for _, dd := range rd.delDescs {
		offset := rd.offset + rIndex*rd.size + int(dd.offset)

		for i := 0; i < int(dd.size); i++ {
			cpBytes[offset+i] = dd.value
		}
	}
}

func (r *Record) nameToField(name string, index int, value string) (*Field, error) {
	rType := r.rType
	fType, ok := nameToFt[rType][name]
	if !ok {
		return nil, fmt.Errorf("bad field name: %s", name)
	}

	f, err := r.NewFieldWithValue(fType, index, value)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (r *Record) NewFieldWithValue(fType FieldType, index int, str string) (*Field, error) {
	fd := (*r.fDesc)[fType]
	if fd == nil {
		for _, fi := range r.rDesc.fInfos {
			if fi.fType == fType {
				fd = &fDesc{&fi, r, make([]*Field, 0)}
				break
			}
		}
		fd.record = r
		fd.fields = make([]*Field, 0)
		(*r.fDesc)[fType] = fd
	}

	f := r.NewField(fType)
	f.fIndex = index
	switch f.valueType {
	case VtListIndex, VtMemberListIndex:
		if len(r.codeplug.rDesc[f.listRecordType].records) == 0 {
			f.value = deferredValue{value: f.value}
			return f, nil
		}
	}
	err := f.SetString(str)
	if err != nil {
		return f, err
	}

	return f, nil
}

func (r *Record) MoveField(dIndex int, f *Field) {
	sIndex := f.fIndex
	r.RemoveField(f)
	if sIndex < dIndex {
		dIndex--
	}

	f.fIndex = dIndex
	r.InsertField(f)
}

func (r *Record) InsertField(f *Field) error {
	fType := f.fType
	i := f.fIndex
	fields := (*r.fDesc)[fType].fields
	fields = append(fields[:i], append([]*Field{f}, fields[i:]...)...)

	for i, f := range fields {
		f.fIndex = i
	}

	(*r.fDesc)[fType].fields = fields

	return nil
}

func (r *Record) RemoveField(f *Field) {
	fType := f.fType
	index := -1
	fields := r.Fields(fType)
	for i, field := range fields {
		if field == f {
			index = i
			break
		}
	}
	if index < 0 {
		log.Fatal("RemoveField: bad field")
	}

	deleteField(&fields, index)

	for i, f := range fields {
		f.fIndex = i
	}
	(*r.fDesc)[fType].fields = fields
}

func (r *Record) Copy() *Record {
	copy := *r

	fDesc := make(map[FieldType]*fDesc)
	copy.fDesc = &fDesc

	for fType, fd := range *r.fDesc {
		fDesc := *fd
		(*copy.fDesc)[fType] = &fDesc

		fDesc.fields = make([]*Field, len(fd.fields))
		for i, f := range fd.fields {
			fDesc.fields[i] = f.Copy()
		}
	}

	return &copy
}

func recordNames(records []*Record) []string {
	names := make([]string, len(records))
	for i, r := range records {
		names[i] = r.Name()
	}

	return names
}

func (r *Record) FindFieldByName(fType FieldType, name string) *Field {
	allFields := (*r.fDesc)[fType].fields
	for _, f := range allFields {
		if f.String() == name {
			return f
		}
	}
	return nil
}

func updateDeferredFields(records []*Record) (error, *Field) {
	for _, r := range records {
		for _, fType := range r.FieldTypes() {
			for _, f := range r.Fields(fType) {
				dValue, deferred := f.value.(deferredValue)
				if deferred {
					f.value = dValue.value
					err := f.SetString(dValue.str)
					if err != nil {
						f.value = dValue
						return err, f
					}
				}
			}
		}
	}
	return nil, nil
}
