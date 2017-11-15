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
	"bytes"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const invalidValueString = "*INVALID*"

// A Field represents a field within a Record.
type Field struct {
	*fDesc
	fIndex int
	value
}

// A value represents the value a Field Contains.
type value interface {
	getString(*Field) string
	setString(*Field, string) error
	valid(*Field) error
	load(*Field)
	store(*Field)
}

// An fDesc contains a field type's dynamic information.
type fDesc struct {
	*fieldInfo
	record *Record
	fields []*Field
}

// A fieldInfo contains a field type's static information.
type fieldInfo struct {
	fType          FieldType
	typeName       string
	max            int
	bitOffset      int
	bitSize        int
	valueType      ValueType
	defaultValue   string
	span           *Span
	strings        *[]string
	indexedStrings *[]IndexedString
	enablingValue  string
	enabler        FieldType
	disabler       FieldType
	listRecordType RecordType
	recordInfo     *recordInfo
	extOffset      int
	extSize        int
	extIndex       int
}

// A FieldType represents a field's type
type FieldType string

// A ValueType represents the type of a field's value
type ValueType string

// A Span represents a range of values.
type Span struct {
	min       uint8
	max       uint8
	scale     uint8
	interval  uint8
	minString string
}

// An IndexedString represents a string corresponding to a field's
// integer value.
type IndexedString struct {
	Index  uint16
	String string
}

// Minimum returns the minimum value of a Span.
func (span *Span) Minimum() int {
	return int(span.min) * int(span.scale)
}

// Maximum returns the maximum value of a Span.
func (span *Span) Maximum() int {
	return int(span.max) * int(span.scale)
}

// Step returns the step (minimal increment) for a Span.
func (span *Span) Step() int {
	return int(span.interval) * int(span.scale)
}

// MinString return a special string represented by the span's minimum
// value.  If the span doesn't have such a special value, the empty string
// is returned.
func (span *Span) MinString() string {
	return span.minString
}

// String returns the fields value as a string.
func (f Field) String() string {
	if _, invalid := f.value.(invalidValue); invalid {
		return invalidValueString
	}
	return f.value.getString(&f)
}

// setString set the strings value from the given string, recording a change.
func (f *Field) SetString(str string) error {
	previousString := f.String()
	if str == previousString {
		return nil
	}

	err := f.setString(str)
	if err == nil {
		change := f.Change(previousString)
		change.Complete()
	}
	return err
}

// setString set the strings value from the given string.
func (f *Field) setString(s string) error {
	if s == invalidValueString {
		f.value = invalidValue{value: f.value}
		return nil
	}
	s = strings.TrimSpace(s)
	err := f.value.setString(f, s)
	if err != nil {
		return err
	}

	if invalidValue, invalid := f.value.(invalidValue); invalid {
		f.value = invalidValue.value
	}

	if f.fType == f.record.nameFieldType {
		f.record.rDesc.cachedListNames = nil
	}

	return err
}

// listNames returns a slice of the names of all the records in the
// field's record's rDesc.
func (f *Field) listNames() []string {
	return *f.record.codeplug.rDesc[f.listRecordType].ListNames()
}

// memberListNames returns a slice of the names of the field's member records.
func (f *Field) memberListNames() []string {
	r := f.record
	fieldInfos := r.fieldInfos
	fieldInfo := fieldInfos[len(fieldInfos)-1]
	fDesc := (*r.fDesc)[fieldInfo.fType]
	fields := fDesc.fields
	memberNames := make([]string, len(fields))
	for i, f := range fields {
		name := f.String()
		memberNames[i] = name
	}
	return memberNames
}

// Returns the fields Span struct, if any
func (f *Field) Span() Span {
	if f.span == nil {
		return Span{}
	}
	return *f.span
}

// Strings returns a slice of valid string values for the field.
func (f *Field) Strings() []string {
	var strs []string
	switch f.valueType {
	case VtListIndex:
		strs = []string{}
		if f.indexedStrings != nil {
			strs = append(strs, (*f.indexedStrings)[0].String)
		}

		strs = append(strs, f.listNames()...)

		if f.indexedStrings != nil && len(*f.indexedStrings) > 1 {
			strs = append(strs, (*f.indexedStrings)[1].String)
		}
	case VtMemberListIndex:
		strs = []string{}
		if f.indexedStrings != nil {
			strs = append(strs, (*f.indexedStrings)[0].String)
		}

		strs = append(strs, f.memberListNames()...)

		if f.indexedStrings != nil && len(*f.indexedStrings) > 1 {
			strs = append(strs, (*f.indexedStrings)[1].String)
		}

	case VtCtcssDcs:
		strs = ctcssDcsStrings()

	case VtIStrings, VtCallType:
		strs = *f.strings

	case VtIndexedStrings:
		strs = []string{}

		if f.indexedStrings != nil {
			for _, is := range f.IndexedStrings() {
				strs = append(strs, is.String)
			}
		}

	default:
		log.Fatal("unexpected f.valueType in f.Strings()")
	}

	return strs
}

// IndexedStrings returns the IndexedString struct for the field, if any.
func (f *Field) IndexedStrings() []IndexedString {
	if f.indexedStrings == nil {
		return nil
	}
	return *f.indexedStrings
}

// Type returns the field's type.
func (f *Field) Type() FieldType {
	return f.fType
}

// Record returns the record that the field is part of.
func (f *Field) Record() *Record {
	return f.record
}

// Index returns the field's slice index.
func (f *Field) Index() int {
	return f.fIndex
}

// SetIndex sets the field's slice index.
func (f *Field) SetIndex(index int) {
	f.fIndex = index
}

// FullTypeName returns a string containing the field's record's type name
// and index as well as the field's type name and index. The index is omitted
// if the MaxRecords or MaxFields is 1.
func (f *Field) FullTypeName() string {
	r := f.record
	s := r.typeName

	if r.max > 1 {
		s += fmt.Sprintf("[%s]", r.Name())
	}

	s += "." + f.typeName

	if f.max > 1 {
		s += fmt.Sprintf("[%s]", f.String())
	}

	return s
}

// valid returns nil if the field's value is valid.
func (f *Field) valid() error {
	err := f.value.valid(f)
	if err != nil {
		if _, invalid := f.value.(invalidValue); !invalid {
			f.value = invalidValue{value: f.value}
		}
	}
	if !f.IsEnabled() {
		return nil
	}
	return err
}

// IsValid returns false if the field has previously been determined
// to be invalid. The field can only be invalid if the value read from
// the codeplug file was invalid.
func (f *Field) IsValid() bool {
	_, invalid := f.value.(invalidValue)
	return !invalid
}

func (f *Field) DefaultValue() string {
	return f.defaultValue
}

// load sets the field's value from the field's part of cp.bytes.
func (f *Field) load() {
	f.value.load(f)
}

// store inserts the field's value into the field's part of cp.bytes.
func (f *Field) store() {
	if !f.IsEnabled() {
		if _, invalid := f.value.(invalidValue); invalid {
			// Leave invalid value in the codeplug as we loaded it.
			return
		}
	}

	f.value.store(f)
}

// bytes returns the field's part of cp.bytes.
func (f *Field) bytes() []byte {
	return f.fDesc.bytes(f.record, f.fIndex)
}

// storeBytes stores the field's value into the field's part of cp.bytes.
func (f *Field) storeBytes(bytes []byte) {
	f.fDesc.storeBytes(bytes, f.record, f.fIndex)
}

// TypeName returns the field's type's name.
func (f *Field) TypeName() string {
	return f.typeName
}

// ValueType returns the field's value's type.
func (f *Field) ValueType() ValueType {
	return f.valueType
}

// ListRecordType returns the field's list's record type.
func (f *Field) ListRecordType() RecordType {
	return f.listRecordType
}

// sibling returns the field's sibling field of the given type.
func (f *Field) sibling(fType FieldType) *Field {
	fields := (*f.record.fDesc)[fType].fields
	if len(fields) == 0 {
		return nil
	}
	return fields[0]
}

func (f *Field) IsEnabled() bool {
	enablingField := f.enablingField()
	if enablingField == nil {
		return true
	}

	if !enablingField.IsEnabled() {
		return false
	}

	if f.enabler != "" {
		return enablingField.String() == enablingField.enablingValue
	}

	if f.disabler != "" {
		return enablingField.String() != enablingField.enablingValue
	}

	return true
}

func (f *Field) enablingField() *Field {
	siblingType := f.EnablingFieldType()

	if siblingType == "" {
		return nil
	}

	return f.sibling(siblingType)
}

func (f *Field) EnablingFieldType() FieldType {
	var siblingType FieldType

	if f.enabler != "" {
		siblingType = f.enabler
	}

	if f.disabler != "" {
		siblingType = f.disabler
	}

	return siblingType
}

// fieldDeleted returns true if the field at fIndex is deleted.
func (fd *fDesc) fieldDeleted(r *Record, fIndex int) bool {
	if fd.max == 1 {
		return false
	}

	bytes := fd.bytes(r, fIndex)
	for i := range bytes {
		if bytes[i] != 0 {
			return false
		}
	}

	return true
}

// deleteField marks the field at fIndex as deleted.
func (fd *fDesc) deleteField(r *Record, fIndex int) {
	bytes := fd.bytes(r, fIndex)
	for i := range bytes {
		bytes[i] = 0
	}
	fd.storeBytes(bytes, r, fIndex)
}

func (fd *fDesc) fieldOffset(r *Record, fIndex int) int {
	var offset int
	if fd.extIndex == 0 || fIndex < fd.extIndex {
		offset = r.offset + r.rIndex*r.size + fd.offset(fIndex)
	} else {
		fExtOffset := (fIndex - fd.extIndex) * fd.size()
		offset = fd.extOffset + r.rIndex*fd.extSize + fExtOffset
	}

	return offset
}

// bytes returns the bytes of the field from cp.bytes.
func (fd *fDesc) bytes(r *Record, fIndex int) []byte {
	cp := r.codeplug
	offset := fd.fieldOffset(r, fIndex)
	fieldBytes := cp.bytes[offset : offset+fd.size()]

	bytes := make([]byte, len(fieldBytes))
	copy(bytes, fieldBytes)

	if fd.bitSize >= 8 {
		return bytes
	}

	rightOffset := (fd.bitOffset + fd.bitSize) % 8
	if rightOffset != 0 {
		bytes[0] >>= 8 - byte(rightOffset)
	}
	bytes[0] &= (1 << uint(fd.bitSize)) - 1

	return bytes
}

// storeBytes inserts bytes value into the field's bits in cp.bytes.
func (fd *fDesc) storeBytes(bytes []byte, r *Record, fIndex int) {
	if fd.size() != len(bytes) {
		panic(fmt.Sprintf("%s: storeBytes(%v) size mismatch",
			fd.typeName, bytes))
	}

	cp := r.codeplug
	offset := fd.fieldOffset(r, fIndex)
	if fd.bitSize >= 8 {
		copy(cp.bytes[offset:offset+fd.size()], bytes)
		return
	}

	value := int(bytes[0])
	mask := (1 << uint(fd.bitSize)) - 1

	rightOffset := uint((fd.bitOffset + fd.bitSize) % 8)
	if rightOffset != 0 {
		mask <<= 8 - rightOffset
		value <<= 8 - rightOffset
	}
	mask = ^mask
	if (value & mask) != 0 {
		panic("value wider than bitSize")
	}

	cp.bytes[offset] &= byte(mask)
	cp.bytes[offset] |= byte(value)
}

// offset returns the byte offset of the field at fIndex within the field's
// record bytes.
func (fi *fieldInfo) offset(fIndex int) int {
	return (fi.bitOffset + fIndex*fi.bitSize) / 8
}

// size returns the field's size in bytes
func (fi *fieldInfo) size() (fSize int) {
	return (fi.bitSize + 7) / 8
}

// deferredValidFields contains a list fields whose values could not be
// checked for validity because they have not yet been loaded. They will
// be checked after all fields are loaded.
var deferredValidFields []*Field

// frequency is a field value representing a frequency in Hertz.
type frequency float64

// String returns the frequency's value as a string.
func (v *frequency) getString(f *Field) string {
	return frequencyToString(float64(*v))
}

// setString sets the frequency's value from a string.
func (v *frequency) setString(f *Field, s string) error {
	freq, err := stringToFrequency(s)
	if err != nil {
		return err
	}

	err = f.record.codeplug.frequencyValid(freq)
	if err != nil {
		return err
	}

	*v = frequency(freq)

	return nil
}

// valid returns nil if the frequency's value is valid.
func (v *frequency) valid(f *Field) error {
	return f.record.codeplug.frequencyValid(float64(*v))
}

// load sets the frequency's value from its bits in cp.bytes.
func (v *frequency) load(f *Field) {
	*v = frequency(bytesToFrequency(f.bytes()))
}

// store stores the frequency's value into its bits in cp.bytes.
func (v *frequency) store(f *Field) {
	f.storeBytes(frequencyToBytes(float64(*v)))
}

// onOff is a field value representing a boolean value.
// It is used when a 1 bit in the codeplug represents false
type onOff bool

// String returns the onOff's value as a string.
func (v *onOff) getString(f *Field) string {
	s := "Off"
	if *v {
		s = "On"
	}

	return s
}

// setString sets the onOff's value from a string.
func (v *onOff) setString(f *Field, s string) error {
	switch s {
	case "Off":
		*v = false

	case "On":
		*v = true

	default:
		return fmt.Errorf("must be 'Off' or 'On'")
	}

	return nil
}

// valid returns nil if the onOff's value is valid.
func (v *onOff) valid(f *Field) error {
	return nil
}

// load sets the onOff's value from its bits in cp.bytes.
func (v *onOff) load(f *Field) {
	*v = false
	if f.bytes()[0] == 0 {
		*v = true
	}
}

// store stores the onOff's value into its bits in cp.bytes.
func (v *onOff) store(f *Field) {
	b := 1
	if *v {
		b = 0
	}
	f.storeBytes([]byte{byte(b)})
}

// offOn is a field value representing a boolean value.
// It is used when a 1 bit in the codeplug represents true
type offOn struct {
	onOff
}

// load sets the offOn's value from its bits in cp.bytes.
func (v *offOn) load(f *Field) {
	v.onOff = false
	if f.bytes()[0] != 0 {
		v.onOff = true
	}
}

// store stores the offOn's value into its bits in cp.bytes.
func (v *offOn) store(f *Field) {
	b := 0
	if v.onOff {
		b = 1
	}
	f.storeBytes([]byte{byte(b)})
}

// iStrings is a field value where an integer value is used to index
// into a slice of strings.
type iStrings int

// String returns the iStrings' value as a string.
func (v *iStrings) getString(f *Field) string {
	i := int(*v)
	strings := *f.strings
	if i >= len(strings) {
		return invalidValueString
	}

	return strings[i]
}

// setString sets the iStrings' value from a string.
func (v *iStrings) setString(f *Field, s string) error {
	fd := f.fDesc
	strings := *fd.strings
	for i, str := range strings {
		if s == str && str != "" {
			*v = iStrings(i)
			return nil
		}
	}

	strs := make([]string, 0, len(strings))

	for _, str := range strings {
		if str != "" {
			strs = append(strs, fmt.Sprintf(`"%s"`, str))
		}
	}

	return fmt.Errorf("must be one of %+v", strs)
}

// valid returns nil if the iStrings' value is valid.
func (v *iStrings) valid(f *Field) error {
	fd := f.fDesc
	strings := *fd.strings

	i := int(*v)
	if i >= len(strings) || (i == 0 && strings[0] == "") {
		return fmt.Errorf("%d: bad string index", i)
	}

	return nil
}

// valid returns nil if the iStrings' value is valid.
func (v *iStrings) load(f *Field) {
	*v = iStrings(f.bytes()[0])
}

// store stores the iStrings' value into its bits in cp.bytes.
func (v *iStrings) store(f *Field) {
	f.storeBytes([]byte{byte(*v)})
}

// span is a field value representing a range of integer values
type span uint8

// String returns the span's value as a string.
func (v *span) getString(f *Field) string {
	sp := *f.span
	i := uint8(*v)
	if sp.minString != "" && i == sp.min {
		return sp.minString
	}
	return fmt.Sprintf("%d", int(*v)*int(sp.scale))
}

// setString sets the span's value from a string.
func (v *span) setString(f *Field, s string) error {
	sp := *f.span

	if s == sp.minString && s != "" {
		*v = span(sp.min)
		return nil
	}

	value64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return err
	}
	value := int(value64)

	if err := v.validValue(f, value); err != nil {
		return err
	}

	value = value / int(sp.scale)

	*v = span(value)

	return nil
}

// valid returns nil if the span's value is valid.
func (v *span) valid(f *Field) error {
	sp := *f.span
	value := int(*v)
	return v.validValue(f, value*int(sp.scale))
}

// validValue returns nil if the specified value is valid for a span.
func (v *span) validValue(f *Field, value int) error {
	sp := *f.span

	if value == int(sp.min)*int(sp.scale) && sp.minString != "" {
		return nil
	}

	multiple := int(sp.interval) * int(sp.scale)
	if value%multiple != 0 {
		return fmt.Errorf("is not a multiple of %d", multiple)
	}

	if value < (int(sp.min)*int(sp.scale)) || value > (int(sp.max)*int(sp.scale)) {
		fmt.Errorf("must be between %d and %d", sp.min, sp.max)
	}

	return nil
}

// load sets the span's value from its bits in cp.bytes.
func (v *span) load(f *Field) {
	*v = span(f.bytes()[0])
}

// store stores the span's value into its bits in cp.bytes.
func (v *span) store(f *Field) {
	f.storeBytes([]byte{byte(*v)})
}

// findexedStrings is a field value where specific integer values
// represent specific strings.
type indexedStrings uint16

// String returns the indexedStrings's value as a string.
func (v *indexedStrings) getString(f *Field) string {
	for _, is := range *(*f.fDesc).indexedStrings {
		if is.Index == uint16(*v) {
			return is.String
		}
	}
	return ""
}

// setString sets the indexedStrings's value from a string.
func (v *indexedStrings) setString(f *Field, s string) error {
	fd := f.fDesc

	for _, is := range *fd.indexedStrings {
		if is.String == s {
			*v = indexedStrings(is.Index)
			return nil
		}
	}

	strs := make([]string, 0, len(*fd.indexedStrings))

	for _, is := range *fd.indexedStrings {
		strs = append(strs, is.String)
	}

	return fmt.Errorf("must be one of %#v", strs)
}

// valid returns nil if the indexedStrings's value is valid.
func (v *indexedStrings) valid(f *Field) error {
	fd := f.fDesc

	for _, is := range *fd.indexedStrings {
		if is.Index == uint16(*v) {
			return nil
		}
	}

	return fmt.Errorf("%d: invalid index", uint16(*v))
}

// load sets the indexedStrings's value from its bits in cp.bytes.
func (v *indexedStrings) load(f *Field) {
	*v = indexedStrings(f.bytes()[0])
}

// store stores the indexedStrings's value into its bits in cp.bytes.
func (v *indexedStrings) store(f *Field) {
	f.storeBytes([]byte{byte(*v)})
}

// biFrequency is a field value representing a frequency in Hertz.
type biFrequency float64

// String returns the biFrequency's value as a string.
func (v *biFrequency) getString(f *Field) string {
	return frequencyToString(float64(*v))
}

// setString sets the biFrequency's value from a string.
func (v *biFrequency) setString(f *Field, s string) error {
	return nil
}

// valid returns nil if the biFrequency's value is valid.
func (v *biFrequency) valid(f *Field) error {
	return nil
}

// load sets the biFrequency's value from its bits in cp.bytes.
func (v *biFrequency) load(f *Field) {
	*v = biFrequency(bcdToBinary(bytesToInt(f.bytes()))) / 10
}

// store stores the biFrequency's value into its bits in cp.bytes.
func (v *biFrequency) store(f *Field) {
	i := binaryToBcd(int(*v) * 10)
	f.storeBytes(intToBytes(i, f.size()))
}

// introLine is a field value representing a introductory line of text
type introLine string

// String returns the introLine's value as a string.
func (v *introLine) getString(f *Field) string {
	return string(*v)
}

// setString sets the introLine's value from a string.
func (v *introLine) setString(f *Field, s string) error {
	if utf8.RuneCountInString(s) > f.size()/2 {
		return fmt.Errorf("is too long")
	}

	if err := v.validValue(f, s); err != nil {
		return err
	}

	*v = introLine(s)

	return nil
}

// valid returns nil if the introLine's value is valid.
func (v *introLine) valid(f *Field) error {
	return v.validValue(f, string(*v))
}

// validValue returns nil if the specified value is valid for a introLine.
func (v *introLine) validValue(f *Field, value string) error {
	_, err := stringToUcs2Bytes(value, f.size())

	return err
}

// load sets the introLine's value from its bits in cp.bytes.
func (v *introLine) load(f *Field) {
	*v = introLine(ucs2BytesToString(f.bytes()))
}

// store stores the introLine's value into its bits in cp.bytes.
func (v *introLine) store(f *Field) {
	ucs2, _ := stringToUcs2Bytes(string(*v), f.size())
	f.storeBytes(ucs2)
}

type callType struct {
	iStrings
}

func (v *callType) setString(f *Field, s string) error {
	if s != "All" {
		return v.iStrings.setString(f, s)
	}

	for _, r := range f.record.records {
		field := r.Field(f.fType)
		if field.String() != "All" {
			continue
		}
		return fmt.Errorf("An \"All\" record already exists: %s",
			r.NameField().String())
	}

	err := v.iStrings.setString(f, s)
	if err != nil {
		return err
	}

	f.record.Field(FtDcCallID).SetString("16777215")

	return nil
}

// callID is a field value representing a DMR Call ID
type callID int32

// String returns the callID's value as a string.
func (v *callID) getString(f *Field) string {
	return fmt.Sprintf("%d", int(*v))
}

// setString sets the callID's value from a string.
func (v *callID) setString(f *Field, s string) error {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return fmt.Errorf("must be a positive integer")
	}

	if val >= 16777216 {
		return fmt.Errorf("must be less than 16777216")
	}

	*v = callID(val)

	return nil
}

// valid returns nil if the callID's value is valid.
func (v *callID) valid(f *Field) error {
	return nil
}

// load sets the callID's value from its bits in cp.bytes.
func (v *callID) load(f *Field) {
	*v = callID(bytesToInt(f.bytes()))
}

// store stores the callID's value into its bits in cp.bytes.
func (v *callID) store(f *Field) {
	f.storeBytes(intToBytes(int(*v), f.size()))
}

// radioPassword is a field value representing password for the radio.
type radioPassword string

// String returns the radioPassword's value as a string.
func (v *radioPassword) getString(f *Field) string {
	return string(*v)
}

// setString sets the radioPassword's value from a string.
func (v *radioPassword) setString(f *Field, s string) error {
	length := f.size() * 2
	if len(s) != length {
		return fmt.Errorf("password must be %d characters long", length)
	}
	if err := mustBeNumericAscii(s); err != nil {
		return err
	}

	*v = radioPassword(s)

	return nil
}

// valid returns nil if the radioPassword's value is valid.
func (v *radioPassword) valid(f *Field) error {
	return mustBeNumericAscii(string(*v))
}

// load sets the radioPassword's value from its bits in cp.bytes.
func (v *radioPassword) load(f *Field) {
	intValue := bytesToInt(f.bytes())
	if uint(intValue) == 0xffffffff {
		*v = radioPassword("00000000")
		return
	}
	*v = radioPassword(fmt.Sprintf("%08d", bcdToBinary(intValue)))
}

// store stores the radioPassword's value into its bits in cp.bytes.
func (v *radioPassword) store(f *Field) {
	val, _ := strconv.ParseUint(string(*v), 10, 32)
	bytes := intToBytes(binaryToBcd(int(val)), f.size())
	f.storeBytes(bytes)
}

// pcPassword is a field value representing a password for the computer.
type pcPassword string

// String returns the pcPassword's value as a string.
func (v *pcPassword) getString(f *Field) string {
	if string(*v) == "\xff\xff\xff\xff\xff\xff\xff\xff" {
		return ""
	}

	return strings.ToLower(string(*v))
}

// setString sets the pcPassword's value from a string.
func (v *pcPassword) setString(f *Field, s string) error {
	if s == "" {
		*v = pcPassword(s)
		return nil
	}

	length := f.size()
	if len(s) != length {
		return fmt.Errorf("password must be %d characters long", length)
	}

	if err := mustBePrintableAscii(s); err != nil {
		return err
	}

	*v = pcPassword(strings.ToLower(s))

	return nil
}

// valid returns nil if the pcPassword's value is valid.
func (v *pcPassword) valid(f *Field) error {
	if string(*v) == "\xff\xff\xff\xff\xff\xff\xff\xff" {
		return nil
	}

	return mustBePrintableAscii(string(*v))
}

// load sets the pcPassword's value from its bits in cp.bytes.
func (v *pcPassword) load(f *Field) {
	*v = pcPassword(f.bytes())
}

// store stores the pcPassword's value into its bits in cp.bytes.
func (v *pcPassword) store(f *Field) {
	if string(*v) == "" {
		bytes := []byte("\xff\xff\xff\xff\xff\xff\xff\xff")
		f.storeBytes(bytes)
		return
	}

	f.storeBytes([]byte(*v))
}

// radioName is a field value representing the name of the radio.
type radioName string

// String returns the radioName's value as a string.
func (v *radioName) getString(f *Field) string {
	return string(*v)
}

// setString sets the radioName's value from a string.
func (v *radioName) setString(f *Field, s string) error {
	if utf8.RuneCountInString(s) > f.size()/2 {
		return fmt.Errorf("name too long")
	}

	if err := v.validValue(f, s); err != nil {
		return err
	}

	*v = radioName(s)

	return nil
}

// valid returns nil if the radioName's value is valid.
func (v *radioName) valid(f *Field) error {
	return v.validValue(f, string(*v))
}

// validValue returns nil if the specified value is valid for a radioName.
func (v *radioName) validValue(f *Field, s string) error {
	_, err := stringToUcs2Bytes(s, f.size())
	if err != nil {
		return err
	}

	return mustBePrintableAscii(string(*v))
}

// load sets the radioName's value from its bits in cp.bytes.
func (v *radioName) load(f *Field) {
	*v = radioName(ucs2BytesToString(f.bytes()))
}

// store stores the radioName's value into its bits in cp.bytes.
func (v *radioName) store(f *Field) {
	ucs2, _ := stringToUcs2Bytes(string(*v), f.size())
	f.storeBytes(ucs2)
}

// textMessage is a field value representing a text message
type textMessage string

// String returns the textMessage's value as a string.
func (v *textMessage) getString(f *Field) string {
	return string(*v)
}

// valid returns nil if the textMessage's value is valid.
func (v *textMessage) setString(f *Field, s string) error {
	if utf8.RuneCountInString(s) >= f.size()/2 {
		return fmt.Errorf("line too long")
	}

	_, err := stringToUcs2Bytes(s, f.size())
	if err != nil {
		return err
	}

	*v = textMessage(s)

	return nil
}

// valid returns nil if the textMessage's value is valid.
func (v *textMessage) valid(f *Field) error {
	_, err := stringToUcs2Bytes(string(*v), f.size())
	return err
}

// load sets the textMessage's value from its bits in cp.bytes.
func (v *textMessage) load(f *Field) {
	*v = textMessage(ucs2BytesToString(f.bytes()))
}

// store stores the textMessage's value into its bits in cp.bytes.
func (v *textMessage) store(f *Field) {
	ucs2, _ := stringToUcs2Bytes(string(*v), f.size())
	f.storeBytes(ucs2)
}

// name is a field value representing a utf8 name that must be unique
// among all records containing the field
type uniqueName struct {
	name
}

// name is a field value representing a utf8 name.
type name string

// String returns the name's value as a string.
func (v *name) getString(f *Field) string {
	return string(*v)
}

// setString sets the name's value from a string.
func (v *name) setString(f *Field, s string) error {
	if utf8.RuneCountInString(s) > f.size()/2 {
		return fmt.Errorf("name too long")
	}

	_, err := stringToUcs2Bytes(string(*v), f.size())
	if err != nil {
		return err
	}

	*v = name(s)

	return nil
}

// valid returns nil if the name's value is valid.
func (v *name) valid(f *Field) error {
	return nil
}

// load sets the name's value from its bits in cp.bytes.
func (v *name) load(f *Field) {
	*v = name(ucs2BytesToString(f.bytes()))
}

// store stores the name's value into its bits in cp.bytes.
func (v *name) store(f *Field) {
	ucs2, _ := stringToUcs2Bytes(string(*v), f.size())
	f.storeBytes(ucs2)
}

// privacyNumber is a field value representing a privacy number.
type privacyNumber struct {
	span
}

// String returns the privacyNumber's value as a string.
func (v *privacyNumber) getString(f *Field) string {
	ss := f.sibling(FtCiPrivacy).String()

	value := int(v.span)
	if ss == "Enhanced" && value >= 8 {
		value = 7
	}

	s := v.span.getString(f)
	v.span = span(value)

	return s
}

// setString sets the privacyNumber's value from a string.
func (v *privacyNumber) setString(f *Field, s string) error {
	ss := f.sibling(FtCiPrivacy).String()

	if ss == "Enhanced" && int(v.span) >= 8 {
		return fmt.Errorf("must be less than 8 for enhanced privacy")
	}

	return v.span.setString(f, s)
}

// valid returns nil if the privacyNumber's value is valid.
func (v *privacyNumber) valid(f *Field) error {
	sibling := f.sibling(FtCiPrivacy)
	if sibling == nil {
		deferredValidFields = append(deferredValidFields, f)
		return nil
	}
	ss := sibling.String()

	if ss == "Enhanced" && int(v.span) >= 8 {
		return fmt.Errorf("must be less than 8 for enhanced privacy")
	}

	return v.span.valid(f)
}

// ctcssDcs is a field value representing a CTCSS or DCS tone.
type ctcssDcs int

// String returns the ctcssDcs's value as a string.
func (v *ctcssDcs) getString(f *Field) string {
	s, _ := ctcssDcsCode(int(*v))

	return s
}

// setString sets the ctcssDcs's value from a string.
func (v *ctcssDcs) setString(f *Field, s string) error {
	value := ctcssDcsStringToBinary(s)
	if value >= 0 {
		*v = ctcssDcs(value)
		return nil
	}

	return fmt.Errorf("bad tone designator")
}

// valid returns nil if the ctcssDcs's value is valid.
func (v *ctcssDcs) valid(f *Field) error {
	if _, err := ctcssDcsCode(int(*v)); err != nil {
		return err
	}

	return nil
}

// load sets the ctcssDcs's value from its bits in cp.bytes.
func (v *ctcssDcs) load(f *Field) {
	*v = ctcssDcs(bytesToInt(f.bytes()))
}

// store stores the ctcssDcs's value into its bits in cp.bytes.
func (v *ctcssDcs) store(f *Field) {
	f.storeBytes(intToBytes(int(*v), f.size()))
}

// memberListIndex is a field value representing an index into a slice
// of member records.
type memberListIndex struct {
	listIndex
}

// String returns the memberListIndex's value as a string.
func (v *memberListIndex) getString(f *Field) string {
	name := v.listIndex.getString(f)
	if v.listIndex > 0 && int(v.listIndex) <= len(f.fields) {
		for _, str := range f.memberListNames() {
			if str == name {
				return name
			}
		}
		isSlice := *f.fDesc.indexedStrings
		is := isSlice[len(isSlice)-1]
		v.listIndex = listIndex(is.Index)
		name = is.String
	}

	return name
}

// setString sets the memberListIndex's value from a string.
func (v *memberListIndex) setString(f *Field, s string) error {
	fd := f.fDesc

	if fd.indexedStrings != nil {
		for _, is := range *fd.indexedStrings {
			if is.String == s {
				v.listIndex = listIndex(is.Index)
				return nil
			}
		}
	}

	for _, name := range f.memberListNames() {
		if name == s {
			for i, name := range f.listNames() {
				if name == s {
					v.listIndex = listIndex(i + 1)
					return nil
				}
			}
		}
	}
	return fmt.Errorf("bad record name")
}

// listIndex is a field value representing an index into a slice of records
type listIndex uint16

// String returns the listIndex's value as a string.
func (v *listIndex) getString(f *Field) string {
	fd := f.fDesc

	if fd.indexedStrings != nil {
		for _, is := range *fd.indexedStrings {
			if is.Index == uint16(*v) {
				return is.String
			}
		}
	}

	ind := int(*v) - 1
	names := f.listNames()
	if ind < 0 || ind >= len(names) {
		return ""
	}

	return names[ind]
}

// setString sets the listIndex's value from a string.
func (v *listIndex) setString(f *Field, s string) error {
	fd := f.fDesc

	if fd.indexedStrings != nil {
		for _, is := range *fd.indexedStrings {
			if is.String == s {
				*v = listIndex(is.Index)
				return nil
			}
		}
	}

	for i, name := range f.listNames() {
		if name == s {
			*v = listIndex(i + 1)
			return nil
		}
	}
	return fmt.Errorf("bad record name")
}

// valid returns nil if the listIndex's value is valid.
func (v *listIndex) valid(f *Field) error {
	fd := f.fDesc

	if fd.indexedStrings != nil {
		for _, is := range *fd.indexedStrings {
			if is.Index == uint16(*v) {
				return nil
			}
		}
	}

	listNames := fd.record.codeplug.rDesc[fd.listRecordType].ListNames()
	if listNames == nil {
		deferredValidFields = append(deferredValidFields, f)
		return nil
	}

	if *v > 0 && int(*v) <= len(*listNames) {
		return nil
	}

	return fmt.Errorf("index out of range")
}

// load sets the listIndex's value from its bits in cp.bytes.
func (v *listIndex) load(f *Field) {
	*v = listIndex(bytesToInt(f.bytes()))
}

// store stores the listIndex's value into its bits in cp.bytes.
func (v *listIndex) store(f *Field) {
	f.storeBytes(intToBytes(int(*v), f.size()))
}

// model is a field value representing a codeplug model name
type model struct {
	ascii
}

// String returns the ascii's value as a string.
func (v *model) getString(f *Field) string {
	s := v.ascii.getString(f)

	cpi := f.record.codeplug.codeplugInfo
	if len(cpi.Models) > 1 && s == cpi.Models[1] {
		s = cpi.Models[0]
	}

	return s
}

// setString sets the ascii's value from a string.
func (v *model) setString(f *Field, s string) error {
	cpi := f.record.codeplug.codeplugInfo
	if len(cpi.Models) > 1 && s == cpi.Models[0] {
		s = cpi.Models[1]
	}

	return v.ascii.setString(f, s)
}

// ascii is a field value representing a ASCII string.
type ascii string

// String returns the ascii's value as a string.
func (v *ascii) getString(f *Field) string {
	return string(*v)
}

// setString sets the ascii's value from a string.
func (v *ascii) setString(f *Field, s string) error {
	if len(s) > f.size() {
		return fmt.Errorf("string too long")
	}

	*v = ascii(s)

	return nil
}

// valid returns nil if the ascii's value is valid.
func (v *ascii) valid(f *Field) error {
	return nil
}

// load sets the ascii's value from its bits in cp.bytes.
func (v *ascii) load(f *Field) {
	fBytes := f.bytes()
	nullIndex := bytes.IndexByte(fBytes, 0)
	if nullIndex >= 0 {
		fBytes = fBytes[:nullIndex]
	}
	*v = ascii(string(fBytes))
}

// store stores the ascii's value into its bits in cp.bytes.
func (v *ascii) store(f *Field) {
	bytes := bytes.Repeat([]byte{0xff}, f.size())
	str := string(*v) + "\000"
	copy(bytes, []byte(str))
	f.storeBytes(bytes)
}

// timeStamp is a field value representing a BCD-encoded time string
type timeStamp string

// String returns the timeStamp's value as a string.
func (v *timeStamp) getString(f *Field) string {
	t, _ := time.Parse("20060102150405", string(*v))
	return t.Format("02-Jan-06 15:04:05")
}

// SetString sets the timeStamp's value from a string.
func (v *timeStamp) setString(f *Field, s string) error {
	t, err := time.Parse("02-Jan-06 15:04:05", s)
	if err != nil {
		return err
	}

	*v = timeStamp(t.Format("20060102150405"))

	return nil
}

// valid returns nil if the timeStamp's value is valid.
func (v *timeStamp) valid(f *Field) error {
	for _, r := range string(*v) {
		if r < '0' && r > '9' {
			return fmt.Errorf("timeStamp is not a decimal value")
		}
	}
	return nil
}

// load sets the timeStamp's value from its bits in cp.bytes.
func (v *timeStamp) load(f *Field) {
	*v = timeStamp(bcdBytesToString(f.bytes()))
}

// store stores the timeStamp's value into its bits in cp.bytes.
func (v *timeStamp) store(f *Field) {
	f.storeBytes(stringToBcdBytes(string(*v)))
}

// timeStamp is a field value representing a BCD-encoded time string
type cpsVersion string

// String returns the cpsVersion's value as a string.
func (v *cpsVersion) getString(f *Field) string {
	return string(*v)
}

// SetString sets the cpsVersion's value from a string.
func (v *cpsVersion) setString(f *Field, s string) error {
	if len(s) != f.size() {
		return fmt.Errorf("bad string length")
	}

	*v = cpsVersion(s)

	return nil
}

// valid returns nil if the cpsVersion's value is valid.
func (v *cpsVersion) valid(f *Field) error {
	for _, r := range string(*v) {
		if r < '0' && r > '9' {
			return fmt.Errorf("cpsVersion is not a decimal value")
		}
	}
	return nil
}

// load sets the cpsVersion's value from its bits in cp.bytes.
func (v *cpsVersion) load(f *Field) {
	s := ""
	for _, b := range f.bytes() {
		s += string(int('0') + int(b))
	}
	*v = cpsVersion(s)
}

// store stores the cpsVersion's value into its bits in cp.bytes.
func (v *cpsVersion) store(f *Field) {
	bytes := make([]byte, len(*v))
	for i, r := range *v {
		bytes[i] = byte(int(r) - int('0'))
	}
	f.storeBytes(bytes)
}

type biFilename struct {
	iStrings
}

func (v *biFilename) setString(f *Field, s string) error {
	return nil
}

type deferredValue struct {
	value
	str string
	pos position
}

type invalidValue struct {
	value
}

var cachedCtcssDcsStrings []string

func ctcssDcsStrings() []string {
	if cachedCtcssDcsStrings != nil {
		return cachedCtcssDcsStrings
	}

	count := len(ctcssFrequencies) + 2*len(dcsCodes) + 1
	cachedCtcssDcsStrings = make([]string, count)

	i := 0

	cachedCtcssDcsStrings[i] = "None"
	i++

	for _, f := range ctcssFrequencies {
		cachedCtcssDcsStrings[i] = fmt.Sprintf("%d.%d", f/10, f%10)
		i++
	}

	for _, c := range dcsCodes {
		cachedCtcssDcsStrings[i] = fmt.Sprintf("D%03dN", c)
		i++
	}

	for _, c := range dcsCodes {
		cachedCtcssDcsStrings[i] = fmt.Sprintf("D%03dI", c)
		i++
	}

	return cachedCtcssDcsStrings
}

var dcsCodes = [...]int{
	23, 25, 26, 31, 32, 36, 43, 47, 51, 53, 54, 65, 71, 72,
	73, 74, 114, 115, 116, 122, 125, 131, 132, 134, 143, 145, 152, 155,
	156, 162, 165, 172, 174, 205, 212, 223, 225, 226, 243, 244, 245, 246,
	251, 252, 255, 261, 263, 265, 266, 271, 274, 306, 311, 315, 325, 331,
	332, 343, 346, 351, 356, 364, 365, 371, 411, 412, 413, 423, 431, 432,
	445, 446, 452, 454, 455, 462, 464, 465, 466, 503, 506, 516, 523, 526,
	532, 546, 565, 606, 612, 624, 627, 631, 632, 654, 662, 664, 703, 712,
	723, 731, 732, 734, 743, 754,
}

var ctcssFrequencies = [...]int{
	670, 693, 719, 744, 770, 797, 825, 854, 885, 915, 948, 974,
	1000, 1035, 1072, 1109, 1148, 1188, 1230, 1273, 1318, 1365, 1413, 1462,
	1514, 1567, 1598, 1622, 1655, 1679, 1713, 1738, 1773, 1799, 1835, 1862,
	1899, 1928, 1966, 1995, 2035, 2065, 2107, 2181, 2257, 2291, 2336, 2418,
	2503, 2541,
}

func goodDcsCode(code int) bool {
	i := sort.SearchInts(dcsCodes[:], code)
	if i < len(dcsCodes) && dcsCodes[i] == code {
		return true
	}

	return false
}

func goodCtcssFrequency(frequency int) bool {
	i := sort.SearchInts(ctcssFrequencies[:], frequency)
	if i < len(ctcssFrequencies) && ctcssFrequencies[i] == frequency {
		return true
	}

	return false
}

func ctcssDcsCode(v int) (string, error) {
	if v == 0 || v == 0xffff {
		return "None", nil
	}

	vType := v >> 14

	v = bcdToBinary(v & 0x03fff)
	if v < 0 {
		return "", fmt.Errorf("only decimal digits are permitted")
	}

	switch vType {
	case 0:
		if !goodCtcssFrequency(v) {
			return "", fmt.Errorf("bad ctcss frequency: %3d.%1d",
				v/10, v%10)
		}

		return fmt.Sprintf("%d.%1d", v/10, v%10), nil

	case 1:
		return "", fmt.Errorf("bad CtcssDcs type 0x%04x", v)

	case 2:
		if !goodDcsCode(v) {
			return "", fmt.Errorf("bad dcs code: %03d", v)
		}

		return fmt.Sprintf("D%03dN", v), nil

	case 3:
		if !goodDcsCode(v) {
			return "", fmt.Errorf("bad dcs code: %03d", v)
		}

		return fmt.Sprintf("D%03dI", v), nil
	}

	return "", fmt.Errorf("bad ctcssDcs value")

}

func ctcssDcsStringToBinary(s string) int {
	if s == "None" {
		return 0xffff
	}

	value := 0
	sType := 0

	if s[0] == 'D' {
		if len(s) != 5 {
			return -1
		}

		switch s[4] {
		case 'N':
			sType = 2

		case 'I':
			sType = 3

		default:
			return -1
		}

		v, err := strconv.ParseInt(s, 10, 16)
		value = int(v)
		if err != nil || !goodDcsCode(value) {
			return -1
		}
	} else {
		flt, err := strconv.ParseFloat(s, 16)
		value = int(flt * 10)
		if err != nil || !goodCtcssFrequency(value) {
			return -1
		}
	}

	value = binaryToBcd(value) | sType<<14

	return value
}

func deleteField(fields *[]*Field, i int) {
	copy((*fields)[i:], (*fields)[i+1:])
	(*fields)[len(*fields)-1] = nil
	*fields = (*fields)[:len(*fields)-1]
}

func fieldNames(fields []*Field) []string {
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.String()
	}

	return names
}

func fieldsToStrings(fields []*Field) []string {
	strings := make([]string, len(fields))
	for i, f := range fields {
		strings[i] = f.String()
	}

	return strings
}
