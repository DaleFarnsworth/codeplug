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

type ChangeType string

type Change struct {
	cType   ChangeType
	records []*Record
	fields  []*Field
	strings []string
	changes []*Change
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

const (
	FieldChange         ChangeType = "FieldChange"
	MoveRecordsChange   ChangeType = "MoveRecordsChange"
	InsertRecordsChange ChangeType = "InsertRecordsChange"
	RemoveRecordsChange ChangeType = "RemoveRecordsChange"
	MoveFieldsChange    ChangeType = "MoveFieldsChange"
	InsertFieldsChange  ChangeType = "InsertFieldsChange"
	RemoveFieldsChange  ChangeType = "RemoveFieldsChange"
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
				if len(allRecords) <= r.rIndex {
					continue
				}
				record := allRecords[r.rIndex-1]
				if record != nil {
					strings[i] = record.Name()
				}
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

func (f *Field) Change(previousValue string) *Change {
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

	return change
}

func (cp *Codeplug) MoveRecordsChange(records []*Record) *Change {
	change := recordsChange(MoveRecordsChange, records)

	return change
}

func (cp *Codeplug) InsertRecordsChange(records []*Record) *Change {
	change := recordsChange(InsertRecordsChange, records)

	return change
}

func (cp *Codeplug) RemoveRecordsChange(records []*Record) *Change {
	change := recordsChange(RemoveRecordsChange, records)

	return change
}

func (cp *Codeplug) RecordsFieldChange(recs []*Record) *Change {
	return recordsFieldChange(recs)
}

func deleteChange(changes *[]*Change, i int) {
	copy((*changes)[i:], (*changes)[i+1:])
	(*changes)[len(*changes)-1] = nil
	*changes = (*changes)[:len(*changes)-1]
}

func (change *Change) Complete() {

	switch change.cType {
	case InsertRecordsChange, InsertFieldsChange:
		change.strings = change.refStrings()
	}

	cp := change.Codeplug()

	cp.changed = true

	cp.publishChange(change)
}

func (change *Change) Changes() []*Change {
	return change.changes
}

func (change *Change) AddChange(newChange *Change) {
	change.changes = append(change.changes, newChange)
}
