// Copyright 2017-2018 Dale Farnsworth. All rights reserved.

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
	"sort"
	"strconv"
	"strings"

	l "github.com/dalefarnsworth/codeplug/debug"
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
	delDesc       *delDesc
	fieldInfos    []*fieldInfo
	nameFieldType FieldType
	index         int
	namePrefix    string
	names         []string
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

type fieldRef struct {
	rType RecordType
	fType FieldType
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

	for _, r := range rd.records {
		r.makeNameUnique()
	}
}

// newField creates and returns the address of a new field of the given type.
func (r *Record) NewField(fType FieldType) *Field {
	if r.fDesc == nil {
		m := make(map[FieldType]*fDesc)
		r.fDesc = &m
	}
	fd := (*r.fDesc)[fType]
	f := new(Field)
	if fd == nil {
		for _, fi := range r.rDesc.fieldInfos {
			if fi.fType == fType {
				fd = &fDesc{fi, r, make([]*Field, 0)}
				break
			}
		}
		if fd == nil {
			// bad field type
			fd = &fDesc{r.rDesc.fieldInfos[0], r, make([]*Field, 0)}
		}
		fd.record = r
		fd.fields = make([]*Field, 0)
		(*r.fDesc)[fType] = fd
	}
	f.fDesc = fd
	f.value = newValue(fd.valueType)
	f.SetDefault()

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
		fi.index = i
		if fi.max == 0 {
			fi.max = 1
		}
		fi.recordInfo = ri
		fd := &fDesc{fieldInfo: fi}
		(*r.fDesc)[fi.fType] = fd
		if fi.valueType == VtName || fi.valueType == VtContactName {
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

	indexedStrs := make(map[int]string)
	indexes := make([]int, 0, len(fds))

	for fType, fd := range fds {
		index := fd.fieldInfo.index
		indexes = append(indexes, index)
		indexedStrs[index] = string(fType)
	}
	sort.Ints(indexes)

	fTypes := make([]FieldType, len(indexes))
	for i, index := range indexes {
		fTypes[i] = FieldType(indexedStrs[index])
	}

	return fTypes
}

func (r *Record) AllFieldTypes() []FieldType {
	fTypes := make([]FieldType, 0)
	for _, fi := range r.rDesc.fieldInfos {
		fTypes = append(fTypes, fi.fType)
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

func (r *Record) AllFields() []*Field {
	fields := make([]*Field, 0)
	for _, fType := range r.FieldTypes() {
		for _, f := range r.Fields(fType) {
			fields = append(fields, f)
		}
	}
	return fields
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

func (r *Record) FullTypeName() string {
	s := r.typeName

	if r.max > 1 {
		name := r.Name()
		if r.names != nil {
			name = r.names[r.rIndex]
		}
		if name == "" {
			name = fmt.Sprintf("%d", r.rIndex)
		}
		s += fmt.Sprintf("[%s]", name)
	}

	return s
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
	f := r.NameField()
	if f != nil {
		dValue, deferred := f.value.(deferredValue)
		if deferred {
			return dValue.str
		}
		return r.NameField().String()
	}
	prefix := r.NamePrefix()
	if prefix != "" {
		return fmt.Sprintf("%s %d", prefix, r.rIndex+1)
	}
	names := r.Names()
	if len(names) > 0 && r.rIndex < len(names) {
		return names[r.rIndex]
	}
	return ""
}

func (r *Record) NamePrefix() string {
	return r.rDesc.recordInfo.namePrefix
}

func (r *Record) Names() []string {
	return r.rDesc.recordInfo.names
}

func (r *Record) MaxRecords() int {
	return r.rDesc.recordInfo.max
}

// makeNameUnique renames the record to make it different than all of
// the passed names.
func (r *Record) makeNameUnique() error {
	nameField := r.NameField()
	if nameField == nil {
		return nil
	}

	if r.Name() == "" {
		return nil
	}

	name := r.Name()

	existingRecordWithName := r.codeplug.FindRecordByName(r.rType, name)
	if existingRecordWithName == nil {
		return nil
	}

	if r == existingRecordWithName {
		return nil
	}

	name = strings.TrimSpace(name)
	baseName := strings.TrimRight(name, "0123456789")
	suffix := strings.TrimPrefix(name, baseName)
	if suffix == "" {
		suffix = "2"
	}
	n64, err := strconv.ParseInt(suffix, 10, 32)
	if err != nil {
		l.Fatal("trailing digits not numeric")
	}
	n := int(n64)

	maxNameLen := nameField.bitSize / 16

	for len(baseName) > 0 {
		suffix := fmt.Sprintf("%d", n)
		for len(baseName)+len(suffix) > maxNameLen {
			baseName = baseName[:len(baseName)-1]
		}
		newName := baseName + fmt.Sprintf("%d", n)
		if r.codeplug.FindRecordByName(r.rType, newName) == nil {
			nameField.value = newValue(nameField.ValueType())
			nameField.value.setString(nameField, newName, false)
			return nil
		}
		n += 1
	}

	return fmt.Errorf("too many record copies")
}

// ListNames returns a slice of the names of all records in the rDesc.
func (rd *rDesc) ListNames() *[]string {
	lenCachedListNames := 0
	if rd.cachedListNames != nil {
		lenCachedListNames = len(*rd.cachedListNames)
	}
	recordsLen := len(rd.records)
	if lenCachedListNames == 0 && recordsLen > 0 {
		names := make([]string, recordsLen)
		for i, r := range rd.records {
			name := r.Name()
			if name == "" {
				name = rd.namePrefix + fmt.Sprintf("%d", i+1)
			}
			names[i] = name
		}

		rd.cachedListNames = &names
	}

	return rd.cachedListNames
}

// MemberListNames returns a slice of possible member list names
func (rd *rDesc) MemberListNames(filter func(r *Record) bool) *[]string {
	if rd.rType != RtContacts {
		if filter == nil {
			return rd.ListNames()
		}

		names := make([]string, 0)
		for i, r := range rd.records {
			if !filter(r) {
				continue
			}
			name := r.Name()
			if name == "" {
				name = rd.namePrefix + fmt.Sprintf("%d", i+1)
			}
			names = append(names, name)
		}

		return &names
	}
	names := make([]string, 0, len(rd.records))
	for _, r := range rd.records {
		if filter != nil && !filter(r) {
			continue
		}
		typeField := r.Field(FtDcCallType)
		if typeField.String() == "Group" {
			names = append(names, r.Name())
		}
	}

	return &names
}

// recordIsDeleted returns true if the record at rIndex is deleted.
func (rd *rDesc) recordIsDeleted(cp *Codeplug, rIndex int) bool {
	dd := rd.delDesc
	if dd == nil {
		return false
	}

	offset := rd.offset + rIndex*rd.size + int(dd.offset)

	for i := 0; i < int(dd.size); i++ {
		if cp.bytes[offset+i] != dd.value {
			return false
		}
	}

	return true
}

// deleteRecord marks the record at rIndex as deleted.
func (rd *rDesc) deleteRecord(cp *Codeplug, rIndex int) {
	dd := rd.delDesc
	if dd == nil {
		l.Fatal("can't delete record %s", rd.records[rIndex])
	}

	offset := rd.offset + rIndex*rd.size + int(dd.offset)

	for i := 0; i < int(dd.size); i++ {
		cp.bytes[offset+i] = dd.value
	}
}

func (r *Record) NewFieldWithValue(fType FieldType, index int, str string) (*Field, error) {
	f := r.NewField(fType)
	f.fIndex = index

	if f.mustDeferValue(str) {
		f.deferValue(str)
		return f, nil
	}

	err := f.setString(str)
	if err != nil {
		return f, err
	}

	return f, nil
}

func (r *Record) NewFieldWithDeferredValue(fType FieldType, index int, str string) *Field {
	f := r.NewField(fType)
	f.fIndex = index
	f.deferValue(str)
	return f
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
		l.Fatal("RemoveField: bad field")
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

	for _, fType := range or.FieldTypes() {
		for _, of := range or.Fields(fType) {
			str := of.String()
			str = removeSuffix(of, str)
			str = AddSuffix(of, str)
			f, _ := r.NewFieldWithValue(of.fType, of.fIndex, str)
			f.resolveDeferredValue()
			r.addField(f)
		}
	}

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
	fields := (*r.fDesc)[fType].fields
	for _, f := range fields {
		if f.String() == name {
			return f
		}
	}
	return nil
}

func (r *Record) HasFieldType(fType FieldType) bool {
	return r.Field(fType) != nil
}

func DependentRecords(records []*Record) (newRecords []*Record, depRecords []*Record) {
	dRecsMap := make(map[string]bool)

	for _, r := range records {
		dRecsMap[r.FullTypeName()] = true
	}

	depRecords = make([]*Record, 0)
	for _, r := range records {
		depRecords = append(depRecords, r.dependentRecords(dRecsMap)...)
	}

	return records, depRecords
}

func (r *Record) dependentRecords(dRecsMap map[string]bool) []*Record {
	dRecs := make([]*Record, 0)
	for _, fType := range r.FieldTypes() {
		fields := r.Fields(fType)
		if len(fields) == 0 {
			continue
		}

		rType := fields[0].listRecordType
		if rType == "" {
			continue
		}

		for _, f := range fields {
			dr := r.codeplug.FindRecordByName(rType, f.String())
			if dr == nil {
				continue
			}

			drTypeName := dr.FullTypeName()
			if dRecsMap[drTypeName] {
				continue
			}
			dRecsMap[drTypeName] = true

			dRecs = append(dRecs, dr.dependentRecords(dRecsMap)...)
			dRecs = append(dRecs, dr)
		}
	}

	return dRecs
}

func (r *Record) InCodeplug() bool {
	records := r.codeplug.records(r.rType)
	for _, rec := range records {
		if rec == r {
			return true
		}
	}

	return false
}

func (r *Record) NameExists() bool {
	name := r.Name()
	rv := r.codeplug.FindRecordByName(r.rType, name) != nil
	return rv
}

func fieldRefFields(cp *Codeplug, fieldRefs []fieldRef) []*Field {
	fields := make([]*Field, 0)
	for _, fRef := range fieldRefs {
		fields = append(fields, cp.fields(fRef.rType, fRef.fType)...)
	}

	return fields
}

func RecordsRemoved(change *Change) {
	cp := change.Codeplug()
	for _, r := range change.records {
		name := r.Name()
		fieldRefs := rTypeFieldRefs[r.rType]
		for _, f := range fieldRefFields(cp, fieldRefs) {
			if f.listRecordType != r.rType {
				continue
			}

			if f.String() != name {
				continue
			}

			if f.max > 1 {
				r := f.record
				change := r.RemoveFieldsChange([]*Field{f})
				r.RemoveField(f)
				change.Complete()
				continue
			}

			f.SetString(f.defaultValue)
		}
	}
}
