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
)

type ChangeType string

type Change struct {
	cType        ChangeType
	records      []*Record
	fields       []*Field
	strings      []string
	afterStrings []string
	changes      []*Change
}

func (change *Change) Type() ChangeType {
	return change.cType
}

func (change *Change) Field() *Field {
	switch change.cType {
	case FieldChange:
		return change.fields[0]
	}
	return nil
}

func (change *Change) previousValue() string {
	switch change.cType {
	case FieldChange:
		return change.strings[0]
	}
	return ""
}

func (change *Change) Codeplug() *Codeplug {
	return change.records[0].codeplug
}

func (change *Change) Record() *Record {
	return change.records[0]
}

func (change *Change) Records() []*Record {
	return change.records
}

func (change *Change) FieldType() FieldType {
	if len(change.fields) == 0 {
		return ""
	}
	return change.fields[0].fType
}

func (change *Change) RecordType() RecordType {
	return change.records[0].rType
}

func (change *Change) refStrings() []string {
	var strings []string
	switch change.cType {
	case MoveRecordsChange, InsertRecordsChange, RemoveRecordsChange:
		records := change.records
		strings = make([]string, len(records))
		r := records[0]
		allRecords := r.codeplug.rDesc[r.rType].records
		for i, r := range records {
			strings[i] = ""
			if r.rIndex > 0 {
				strings[i] = allRecords[r.rIndex-1].Name()
			}
		}

	case MoveFieldsChange, InsertFieldsChange, RemoveFieldsChange:
		fields := change.fields
		strings = make([]string, len(fields))
		if len(fields) == 0 {
			return strings
		}
		f := fields[0]
		allFields := (*f.record.fDesc)[f.fType].fields
		for i, f := range fields {
			strings[i] = ""
			if f.fIndex > 0 && f.fIndex <= len(allFields) {
				strings[i] = allFields[f.fIndex-1].String()
			}
		}
	}

	return strings
}

func (change *Change) fieldStrings() []string {
	var strings []string
	strings = make([]string, len(change.fields))
	for i, f := range change.fields {
		if f != nil {
			strings[i] = f.String()
		}
	}

	return strings
}

func (change *Change) sIndex(i int) int {
	name := ""
	switch change.cType {
	case MoveRecordsChange, InsertRecordsChange, RemoveRecordsChange:
		r := change.records[i]
		name = change.strings[i]
		if name == "" {
			return 0
		}
		r = r.codeplug.FindRecordByName(r.rType, name)
		return r.rIndex + 1

	case MoveFieldsChange, InsertFieldsChange, RemoveFieldsChange:
		f := change.fields[i]
		name = change.strings[i]
		if name == "" {
			return 0
		}
		f = f.record.FindFieldByName(f.fType, name)
		return f.fIndex + 1

	}

	return -1
}

const (
	FieldChange         ChangeType = "FieldChange"
	MoveRecordsChange   ChangeType = "MoveRecordsChange"
	InsertRecordsChange ChangeType = "InsertRecordsChange"
	RemoveRecordsChange ChangeType = "RemoveRecordsChange"
	MoveFieldsChange    ChangeType = "MoveFieldsChange"
	InsertFieldsChange  ChangeType = "InsertFieldsChange"
	RemoveFieldsChange  ChangeType = "RemoveFieldsChange"
	ListIndexChange     ChangeType = "ListIndexChange"
	RecordsFieldChange  ChangeType = "RecordsFieldChange"
)

func fieldChange(f *Field, previousValue string) *Change {
	change := Change{
		cType:   FieldChange,
		records: []*Record{f.record},
		fields:  []*Field{f},
		strings: []string{previousValue},
	}

	return &change
}

func fieldsChange(t ChangeType, r *Record, fields []*Field) *Change {
	change := Change{
		cType:   t,
		records: []*Record{r},
		fields:  fields,
	}
	change.strings = change.refStrings()

	return &change
}

func recordsFieldChange(recs []*Record) *Change {
	change := Change{
		cType:   RecordsFieldChange,
		records: recs,
	}

	return &change
}

func recordsChange(t ChangeType, records []*Record) *Change {
	change := Change{
		cType:   t,
		records: records,
	}
	change.strings = change.refStrings()

	return &change
}

func listIndexChange(r *Record, fields []*Field) *Change {
	change := Change{
		cType:   ListIndexChange,
		records: []*Record{r},
		fields:  fields,
	}
	change.strings = change.fieldStrings()

	return &change
}

func (f *Field) Change(previousValue string) *Change {
	cp := f.record.codeplug

	change := cp.currentChange()
	if change != nil && change.cType == FieldChange &&
		change.Field() == f && f.String() != change.previousValue() &&
		change.previousValue() != invalidValueString {

		return change
	}

	return fieldChange(f, previousValue)
}

func (r *Record) MoveFieldsChange(fields []*Field) *Change {
	return fieldsChange(MoveFieldsChange, r, fields)
}

func (r *Record) InsertFieldsChange(fields []*Field) *Change {
	return fieldsChange(InsertFieldsChange, r, fields)
}

func (r *Record) RemoveFieldsChange(fields []*Field) *Change {
	change := fieldsChange(RemoveFieldsChange, r, fields)

	fType := change.FieldType()
	change.changes = rDescChanges(r.rDesc, fields[0].listRecordType, fType)

	return change
}

func (cp *Codeplug) MoveRecordsChange(records []*Record) *Change {
	change := recordsChange(MoveRecordsChange, records)

	change.changes = cp.listIndexChanges(change)

	return change
}

func (cp *Codeplug) InsertRecordsChange(records []*Record) *Change {
	change := recordsChange(InsertRecordsChange, records)

	change.changes = cp.listIndexChanges(change)

	return change
}

func (cp *Codeplug) RemoveRecordsChange(records []*Record) *Change {
	change := recordsChange(RemoveRecordsChange, records)

	change.changes = cp.listIndexChanges(change)

	return change
}

func (cp *Codeplug) RecordsFieldChange(recs []*Record) *Change {
	return recordsFieldChange(recs)
}

func (cp *Codeplug) listIndexChanges(change *Change) []*Change {
	rType := change.RecordType()
	fType := change.FieldType()
	changes := []*Change{}
	for _, rd := range cp.rDesc {
		changes = append(changes, rDescChanges(rd, rType, fType)...)
	}
	return changes
}

func rDescChanges(rd *rDesc, rType RecordType, fType FieldType) []*Change {
	changes := []*Change{}
	for _, fi := range rd.fieldInfos {
		if fi.listRecordType != rType || fi.fType == fType {
			continue
		}
		switch fi.valueType {
		case VtListIndex, VtMemberListIndex:
			rChanges := recordChanges(rd.records, fi.fType)
			changes = append(changes, rChanges...)
		}
	}
	return changes
}

func recordChanges(records []*Record, fType FieldType) []*Change {
	changes := []*Change{}
	for _, r := range records {
		refChange := listIndexChange(r, r.Fields(fType))
		changes = append(changes, refChange)
	}
	return changes
}

func deleteChange(changes *[]*Change, i int) {
	copy((*changes)[i:], (*changes)[i+1:])
	(*changes)[len(*changes)-1] = nil
	*changes = (*changes)[:len(*changes)-1]
}

func (cp *Codeplug) updateListIndexChanges(changes []*Change) []*Change {
	if changes == nil {
		return changes
	}
	for i := 0; i < len(changes); i++ {
		change := changes[i]
		changeStrings := change.strings
		if stringsEqual(change.fieldStrings(), changeStrings) {
			deleteChange(&changes, i)
			i--
			continue
		}
		r := change.Record()
		fType := change.FieldType()
		fields := r.Fields(fType)
		for i := len(fields) - 1; i >= 0; i-- {
			r.RemoveField(fields[i])
		}
		for i, str := range changeStrings {
			f, err := r.NewFieldWithValue(fType, i, str)
			if err == nil {
				r.addField(f)
			}
		}
		if len(r.Fields(fType)) == 0 {
			fd := (*r.fDesc)[fType]
			indexedStrings := fd.indexedStrings
			if indexedStrings != nil {
				iStrs := *indexedStrings
				str := iStrs[len(iStrs)-1].String
				f, _ := r.NewFieldWithValue(fType, 0, str)
				r.addField(f)
			}
		}
		fields = r.Fields(fType)
		strings := make([]string, len(fields))
		for i, f := range fields {
			strings[i] = f.String()
		}
		change.afterStrings = strings
	}
	return changes
}

func (cp *Codeplug) addChange(change *Change) {
	cp.changed = true

	switch change.cType {
	case MoveRecordsChange, InsertRecordsChange, RemoveRecordsChange,
		RemoveFieldsChange:
		change.changes = cp.updateListIndexChanges(change.changes)
	}

	i := cp.changeIndex + 1
	cp.changeList = append(cp.changeList[:i], change)
	cp.changeIndex = len(cp.changeList) - 1
}

func (change *Change) Complete() {
	switch change.cType {
	case InsertRecordsChange, InsertFieldsChange:
		change.strings = change.refStrings()
	}

	r := change.Record()
	cp := r.codeplug
	if change != cp.currentChange() {
		cp.addChange(change)
	}
	cp.publishChange(change)
}

func (cp *Codeplug) FindRecordByName(rType RecordType, name string) *Record {
	allRecords := cp.rDesc[rType].records
	for _, r := range allRecords {
		if r.Name() == name {
			return r
		}
	}
	return nil
}

func (cp *Codeplug) currentChange() *Change {
	if cp.changeList != nil {
		return cp.changeList[cp.changeIndex]
	}
	return nil
}

func (change *Change) undoReference() string {
	switch change.cType {
	case MoveRecordsChange, InsertRecordsChange, RemoveRecordsChange:
		r := change.records[0]
		r = r.codeplug.FindRecordByName(r.rType, r.Name())
		if r == nil {
			return ""
		}
		index := r.rIndex - 1

		if index < 0 {
			return "to top"
		}

		allRecords := r.codeplug.rDesc[r.rType].records
		return "after " + allRecords[index].Name()

	case MoveFieldsChange, InsertFieldsChange, RemoveFieldsChange:
		f := change.fields[0]
		r := f.record
		f = r.FindFieldByName(f.fType, f.String())
		if f == nil {
			return ""
		}
		index := f.fIndex - 1

		if index < 0 {
			return "to top"
		}

		allFields := (*r.fDesc)[f.fType].fields
		return "after " + allFields[index].String()
	}
	return ""
}

func (change *Change) Changes() []*Change {
	return change.changes
}

func (change *Change) redoReference() string {
	name := change.strings[0]
	if name == "" {
		return "to top"
	}

	return "after " + name
}

func (cp *Codeplug) UndoString() string {
	var str string

	changeList := cp.changeList
	changeIndex := cp.changeIndex

	if changeList == nil {
		return ""
	}

	if changeIndex < 0 ||
		changeIndex >= len(changeList) {
		logFatal("UndoString: bad changeIndex")
	}

	if changeIndex == 0 || len(changeList) == 1 {
		return ""
	}

	change := changeList[changeIndex]
	r := change.Record()
	rTypeName := r.TypeName()
	rName := r.Name()

	cType := change.cType
	switch cType {
	case FieldChange:
		f := change.Field()
		fTypeName := f.TypeName()
		value := f.String()
		prevVal := change.previousValue()
		str = fmt.Sprintf("%s.%s: %s: '%s' -> '%s'",
			rTypeName, rName, fTypeName, prevVal, value)

	case RecordsFieldChange:
		rCount := len(change.records)
		fChange := change.changes[0]
		f := fChange.Field()
		fTypeName := f.TypeName()
		value := f.String()
		howMany := "all"
		if rCount != len(f.record.records) {
			howMany = fmt.Sprintf("%d", rCount)
		}
		str = fmt.Sprintf("Set '%s' to '%s' in %s %s",
			fTypeName, value, howMany, rTypeName)

	case MoveRecordsChange:
		names := maxNamesString(recordNames(change.records), 5)
		ref := change.undoReference()
		str = fmt.Sprintf("%s: move %s %s", rTypeName, names, ref)

	case InsertRecordsChange:
		names := maxNamesString(recordNames(change.records), 5)
		ref := change.undoReference()
		str = fmt.Sprintf("%s: insert %s %s", rTypeName, names, ref)

	case RemoveRecordsChange:
		names := maxNamesString(recordNames(change.records), 5)
		str = fmt.Sprintf("%s: delete %s", rTypeName, names)

	case MoveFieldsChange:
		names := maxNamesString(fieldNames(change.fields), 5)
		ref := change.undoReference()
		str = fmt.Sprintf("%s.%s: move %s %s",
			rTypeName, rName, names, ref)

	case InsertFieldsChange:
		names := maxNamesString(fieldNames(change.fields), 5)
		ref := change.undoReference()
		str = fmt.Sprintf("%s.%sinsert %s %s",
			rTypeName, rName, names, ref)

	case RemoveFieldsChange:
		names := maxNamesString(fieldNames(change.fields), 5)
		str = fmt.Sprintf("%s.%sdelete %s", rTypeName, rName, names)

	default:
		logFatal("undoString: unexpected change type:", cType)
	}

	return str
}

func (cp *Codeplug) RedoString() string {
	var str string

	changeList := cp.changeList
	changeIndex := cp.changeIndex

	if changeList == nil {
		return ""
	}

	if changeIndex < 0 ||
		changeIndex >= len(changeList) {
		logFatal("UndoString: bad changeIndex")
	}

	if changeIndex == len(changeList)-1 {
		return ""
	}

	change := changeList[changeIndex+1]
	r := change.Record()
	rTypeName := r.TypeName()
	rName := r.Name()

	cType := change.cType
	switch cType {
	case FieldChange:
		f := change.Field()
		fTypeName := f.TypeName()
		value := f.String()
		prevValue := change.previousValue()
		str = fmt.Sprintf("%s.%s: %s: '%s' -> '%s'",
			rTypeName, rName, fTypeName, value, prevValue)

	case RecordsFieldChange:
		rCount := len(change.records)
		fChange := change.changes[0]
		f := fChange.Field()
		fTypeName := f.TypeName()
		prevValue := fChange.previousValue()
		howMany := "all"
		if rCount != len(f.record.records) {
			howMany = fmt.Sprintf("%d", rCount)
		}
		str = fmt.Sprintf("Set '%s' to '%s' in %s %s",
			fTypeName, prevValue, howMany, rTypeName)

	case MoveRecordsChange:
		names := maxNamesString(recordNames(change.records), 5)
		ref := change.redoReference()
		str = fmt.Sprintf("%s: move %s %s", rTypeName, names, ref)

	case InsertRecordsChange:
		names := maxNamesString(recordNames(change.records), 5)
		ref := change.redoReference()
		str = fmt.Sprintf("%s: insert %s %s", rTypeName, names, ref)

	case RemoveRecordsChange:
		names := maxNamesString(recordNames(change.records), 5)
		str = fmt.Sprintf("%s: delete %s", rTypeName, names)

	case MoveFieldsChange:
		names := maxNamesString(fieldNames(change.fields), 5)
		ref := change.redoReference()
		str = fmt.Sprintf("%s.%s: move %s %s",
			rTypeName, rName, names, ref)

	case InsertFieldsChange:
		names := maxNamesString(fieldNames(change.fields), 5)
		ref := change.redoReference()
		str = fmt.Sprintf("%s.%s: insert %s %s",
			rTypeName, rName, names, ref)

	case RemoveFieldsChange:
		names := maxNamesString(fieldNames(change.fields), 5)
		str = fmt.Sprintf("%s.%s: delete %s", rTypeName, rName, names)

	default:
		logFatal("undoString: unexpected change type:", cType)
	}

	return str
}

func (cp *Codeplug) UndoNameChangeCount() (string, int) {
	changeList := cp.changeList
	index := cp.changeIndex

	if len(changeList) == 1 {
		return "", 0
	}

	if index == 0 || index >= len(changeList) {
		logFatal("Undo: bad changeIndex")
	}

	change := changeList[index]

	return change.Record().TypeName(), len(change.changes)
}

func (cp *Codeplug) UndoChange(progFunc func(int)) {
	changeList := cp.changeList
	index := cp.changeIndex

	if len(changeList) == 1 {
		return
	}

	if index == 0 || index >= len(changeList) {
		logFatal("Undo: bad changeIndex")
	}

	cp.changeIndex--
	cp.changed = true

	change := changeList[index]
	change = cp.undoChange(change, progFunc)
	cp.publishChange(change)
}

func (cp *Codeplug) undoChange(change *Change, progFunc func(int)) *Change {
	cType := change.cType

	switch cType {
	case FieldChange:
		f := change.Field()
		previousValue := f.String()
		err := f.setString(change.previousValue())
		if err != nil {
			logFatal("UndoChange: FieldChange error ", err.Error())
		}

		*change = *fieldChange(f, previousValue)

	case RecordsFieldChange:
		for i, change := range change.changes {
			cp.undoChange(change, nil)
			progFunc(i)
		}

	case MoveRecordsChange:
		strings := change.refStrings()
		for i, r := range change.records {
			r = cp.FindRecordByName(r.rType, r.Name())
			dIndex := change.sIndex(i)
			cp.MoveRecord(dIndex, r)
		}
		change.strings = strings

		cp.undoListIndexChanges(change)

	case InsertRecordsChange:
		records := change.records
		rIndexes := make([]int, len(records))
		for i, r := range records {
			rIndexes[i] = r.rIndex
		}

		for i := len(records) - 1; i >= 0; i-- {
			cp.RemoveRecord(records[i])
			records[i].rIndex = rIndexes[i]
		}

		cp.undoListIndexChanges(change)

		newChange := *change
		change = &newChange
		change.cType = RemoveRecordsChange

	case RemoveRecordsChange:
		for i, r := range change.records {
			r.rIndex = change.sIndex(i)
			cp.InsertRecord(r)
		}

		cp.undoListIndexChanges(change)

		newChange := *change
		change = &newChange
		change.cType = InsertRecordsChange

	case MoveFieldsChange:
		fields := change.fields
		r := fields[0].record
		strings := change.refStrings()
		for i, f := range fields {
			f = r.FindFieldByName(f.fType, f.String())
			dIndex := change.sIndex(i)
			r.MoveField(dIndex, f)
		}
		change.strings = strings

	case InsertFieldsChange:
		r := change.Record()
		fields := change.fields
		for i := len(fields) - 1; i >= 0; i-- {
			r.RemoveField(fields[i])
		}

		newChange := *change
		change = &newChange
		change.cType = RemoveFieldsChange

	case RemoveFieldsChange:
		r := change.Record()
		for i, f := range change.fields {
			f.fIndex = change.sIndex(i)
			r.InsertField(f)
		}

		cp.undoListIndexChanges(change)

		newChange := *change
		change = &newChange
		change.cType = InsertFieldsChange

	case ListIndexChange:
		fType := change.FieldType()
		r := change.Record()
		fields := r.Fields(fType)
		for i := len(fields) - 1; i >= 0; i-- {
			r.RemoveField(fields[i])
		}
		for i, str := range change.strings {
			f, err := r.NewFieldWithValue(fType, i, str)
			if err != nil {
				logFatal("NewField error on ", str)
			}
			r.addField(f)
		}
		c := change
		c.strings, c.afterStrings = c.afterStrings, c.strings

	default:
		logFatal("Undo: unexpected change type:", cType)
	}

	return change
}

func (cp *Codeplug) RedoNameChangeCount() (string, int) {
	changeList := cp.changeList
	index := cp.changeIndex

	if len(changeList) <= index+1 {
		return "", 0
	}

	if index < 0 {
		logFatal("Redo: bad changeIndex")
	}

	change := changeList[index+1]

	return change.Record().TypeName(), len(change.changes)
}

func (cp *Codeplug) RedoChange(progFunc func(int)) {
	changeList := cp.changeList
	index := cp.changeIndex

	if len(changeList) <= index+1 {
		return
	}

	if index < 0 {
		logFatal("Redo: bad changeIndex")
	}

	index++
	cp.changeIndex = index
	cp.changed = true

	change := changeList[index]
	change = cp.redoChange(change, progFunc)
	cp.publishChange(change)
}

func (cp *Codeplug) redoChange(change *Change, progFunc func(int)) *Change {
	cType := change.cType

	switch cType {
	case FieldChange:
		f := change.Field()
		previousValue := f.String()
		err := f.setString(change.previousValue())
		if err != nil {
			logFatal("RedoChange: FieldChange error ", err.Error())
		}

		*change = *fieldChange(f, previousValue)

	case RecordsFieldChange:
		for i, change := range change.changes {
			cp.redoChange(change, nil)
			progFunc(i)
		}

	case MoveRecordsChange:
		strings := change.refStrings()
		for i, r := range change.records {
			r = cp.FindRecordByName(r.rType, r.Name())
			dIndex := change.sIndex(i)
			cp.MoveRecord(dIndex, r)
		}
		change.strings = strings

		cp.redoListIndexChanges(change)

	case InsertRecordsChange:
		for _, r := range change.records {
			cp.InsertRecord(r)
		}

		cp.redoListIndexChanges(change)

		newChange := *change
		change = &newChange
		change.cType = RemoveRecordsChange

	case RemoveRecordsChange:
		records := change.records
		for i := len(records) - 1; i >= 0; i-- {
			cp.RemoveRecord(records[i])
		}

		cp.redoListIndexChanges(change)

		newChange := *change
		change = &newChange
		change.cType = InsertRecordsChange

	case MoveFieldsChange:
		fields := change.fields
		r := fields[0].record
		strings := change.refStrings()
		for i, f := range fields {
			f = r.FindFieldByName(f.fType, f.String())
			dIndex := change.sIndex(i)
			r.MoveField(dIndex, f)
		}
		change.strings = strings

	case InsertFieldsChange:
		r := change.Record()
		for _, f := range change.fields {
			r.InsertField(f)
		}

		newChange := *change
		change = &newChange
		change.cType = RemoveFieldsChange

	case RemoveFieldsChange:
		r := change.Record()
		fields := change.fields
		for i := len(fields) - 1; i >= 0; i-- {
			r.RemoveField(fields[i])
		}

		cp.redoListIndexChanges(change)

		newChange := *change
		change = &newChange
		change.cType = InsertFieldsChange

	case ListIndexChange:
		fType := change.FieldType()
		r := change.Record()
		fields := r.Fields(fType)
		for i := len(fields) - 1; i >= 0; i-- {
			r.RemoveField(fields[i])
		}
		for i, str := range change.strings {
			f, _ := r.NewFieldWithValue(fType, i, str)
			r.addField(f)
		}
		c := change
		c.strings, c.afterStrings = c.afterStrings, c.strings

	default:
		logFatal("Redo: unexpected change type:", cType)
	}

	return change
}

func (cp *Codeplug) undoListIndexChanges(change *Change) {
	for _, change = range change.changes {
		cp.undoChange(change, nil)
	}
}

func (cp *Codeplug) redoListIndexChanges(change *Change) {
	for _, change = range change.changes {
		cp.redoChange(change, nil)
	}
}
