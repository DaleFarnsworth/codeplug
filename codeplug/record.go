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
	"strconv"
	"strings"
)

// A Record represents a record within a Codeplug.
type Record struct {
	*rDesc
	fDesc  *map[FieldType]*fDesc
	rIndex int
}

// An rDesc contains a record type's dynamic information.
type rDesc struct {
	*recordInfo
	codeplug        *Codeplug
	records         []*Record
	cachedListNames *[]string
}

// A recordInfo contains a record type's static information.
type recordInfo struct {
	rType         RecordType
	typeName      string
	max           int
	offset        int
	size          int
	delDescs      []delDesc
	fieldInfos    []*fieldInfo
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

// loadRecords loads all the records in rDesc from the codeplug's file.
func (rd *rDesc) loadRecords() {
	records := rd.records
	records = make([]*Record, rd.max)

	cp := rd.codeplug
	length := 0
	for rIndex := range records {
		if !rd.recordIsDeleted(cp, rIndex) {
			r := records[rIndex]
			if r == nil {
				r = cp.newRecord(rd.rType, rIndex)
			}

			r.load()
			nameField := r.NameField()
			if nameField != nil && nameField.String() == "" {
				continue
			}

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
		if fd == nil {
			for _, fi := range r.rDesc.fieldInfos {
				if fi.fType == fType {
					fd = &fDesc{fi, r, make([]*Field, 0)}
					break
				}
			}
			fd.record = r
			fd.fields = make([]*Field, 0)
			(*r.fDesc)[fType] = fd
		}
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
		return fmt.Errorf("too many fields: %s", string(f.fType))
	}

	fd := (*r.fDesc)[f.fType]
	f.fIndex = len(fd.fields)
	fd.fields = append(fd.fields, f)

	return nil
}

// load replaces the record's contents with the fields found in
// the codeplug.
func (r *Record) load() {
	ri := r.rDesc.recordInfo

	for i := range ri.fieldInfos {
		fi := ri.fieldInfos[i]
		if fi.max == 0 {
			fi.max = 1
		}
		fi.recordInfo = ri
		fd := &fDesc{fieldInfo: fi}
		(*r.fDesc)[fi.fType] = fd
		if fi.valueType == VtName || fi.valueType == VtUniqueName {
			ri.nameFieldType = fi.fType
		}
		fd.record = r
	}

	for _, fd := range *r.fDesc {
		fi := fd.fieldInfo
		fields := make([]*Field, fi.max)

		length := 0
		for fIndex := range fields {
			if !fd.fieldDeleted(r, fIndex) {
				f := fields[fIndex]
				if f == nil {
					f = &Field{}
				}

				f.fDesc = fd
				f.fIndex = fIndex
				f.value = newValue(fi.valueType)

				f.load()

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
func (r *Record) store() {
	for _, fd := range *r.fDesc {
		for fIndex := 0; fIndex < fd.max; fIndex++ {
			if fIndex < len(fd.fields) {
				fd.fields[fIndex].store()
			} else {
				fd.deleteField(r, fIndex)
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
	fDesc := (*r.fDesc)[fType]
	if fDesc == nil {
		return nil
	}
	return fDesc.fields
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

func (r *Record) hasUniqueNames() bool {
	nameField := r.NameField()
	if nameField == nil {
		return false
	}

	_, unique := r.NameField().value.(*uniqueName)
	return unique
}

// makeNameUnique renames the record to make it different than all of
// the passed names.
func (r *Record) makeNameUnique(namesp *[]string) error {
	if namesp == nil {
		return nil
	}

	names := *namesp
	if len(names) == 0 {
		return nil
	}

	nameField := r.NameField()
	name := nameField.String()

	if !stringInSlice(name, names) {
		return nil
	}

	baseName := strings.TrimRight(strings.TrimSpace(name), "0123456789")
	suffix := strings.TrimPrefix(name, baseName)
	if suffix == "" {
		suffix = "2"
	}
	n64, err := strconv.ParseInt(suffix, 10, 32)
	if err != nil {
		log.Fatal("trailing digits not numeric")
	}
	n := int(n64)

	maxNameLen := nameField.bitSize / 16

	for len(baseName) > 0 {
		suffix := fmt.Sprintf("%d", n)
		for len(baseName)+len(suffix) > maxNameLen {
			baseName = baseName[:len(baseName)-1]
		}
		newName := strings.TrimSpace(baseName) + fmt.Sprintf("%d", n)
		if !stringInSlice(newName, names) {
			nameField.value = newValue(nameField.ValueType())
			nameField.value.setString(nameField, newName)
			return nil
		}
		n += 1
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

// MemberListNames returns a slice of possible member list names
func (rd *rDesc) MemberListNames() *[]string {
	if rd.rType != RecordType("DigitalContacts") {
		return rd.ListNames()
	}
	names := make([]string, 0, len(rd.records))
	for _, r := range rd.records {
		typeField := r.Field(FieldType("CallType"))
		if typeField.String() == "Group" {
			names = append(names, r.NameField().String())
		}
	}

	return &names
}

// recordIsDeleted returns true if the record at rIndex is deleted.
func (rd *rDesc) recordIsDeleted(cp *Codeplug, rIndex int) bool {
nextDelDesc:
	for _, dd := range rd.delDescs {
		offset := rd.offset + rIndex*rd.size + int(dd.offset)

		for i := 0; i < int(dd.size); i++ {
			if cp.bytes[offset+i] != dd.value {
				continue nextDelDesc
			}
		}
		return true
	}

	return false
}

// deleteRecord marks the record at rIndex as deleted.
func (rd *rDesc) deleteRecord(cp *Codeplug, rIndex int) {
	for _, dd := range rd.delDescs {
		offset := rd.offset + rIndex*rd.size + int(dd.offset)

		for i := 0; i < int(dd.size); i++ {
			cp.bytes[offset+i] = dd.value
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
	f := r.NewField(fType)
	f.fIndex = index
	switch f.valueType {
	case VtListIndex, VtMemberListIndex:
		if len(r.codeplug.rDesc[f.listRecordType].records) == 0 {
			f.value = deferredValue{value: f.value}
			return f, nil
		}
	}
	err := f.setString(str)
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

func (or *Record) Copy() *Record {
	r := new(Record)
	r.rDesc = or.rDesc
	r.rIndex = 0

	rfDesc := make(map[FieldType]*fDesc)
	for fType, fd := range *or.fDesc {
		fields := fd.CopyFields()
		fDesc := new(fDesc)
		fDesc.fieldInfo = fd.fieldInfo
		if len(fields) > 0 {
			fDesc = fields[0].fDesc
		}
		fDesc.fields = fields
		fDesc.record = r

		rfDesc[fType] = fDesc
	}
	r.fDesc = &rfDesc

	return r
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
					err := f.setString(dValue.str)
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
