// Copyright 2017 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of Ui.
//
// Ui is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU Lesser General Public
// License as published by the Free Software Foundation.
//
// Ui is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Ui.  If not, see <http://www.gnu.org/licenses/>.

package ui

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"

	"github.com/dalefarnsworth/codeplug/codeplug"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type fieldNameList struct {
	fieldNames    []string
	fieldMembers  *FieldMembers
	qListView     *widgets.QListView
	model         *core.QStringListModel
	fieldToInsert *codeplug.Field
}

func (fnl *fieldNameList) current() int {
	return fnl.qListView.CurrentIndex().Row()
}

func (fnl *fieldNameList) selectedNames() []string {
	indexes := fnl.qListView.SelectedIndexes()
	strs := make([]string, len(indexes))
	for i, index := range indexes {
		strs[i] = fnl.fieldNames[index.Row()]
	}
	return strs
}

func (fnl *fieldNameList) clearSelection() {
	fnl.qListView.ClearSelection()
}

func (fnl *fieldNameList) Width() int {
	return fnl.qListView.SizeHintForColumn(0) * 6 / 5
}

func (fnl *fieldNameList) SetWidth(width int) {
	fnl.qListView.SetMinimumWidth(width)
	fnl.qListView.SetMaximumWidth(width)
}

func addMemberList(parent *VBox, fm *FieldMembers, fieldNames []string) *fieldNameList {
	fnl := new(fieldNameList)
	view := widgets.NewQListView(nil)
	fnl.qListView = view
	fnl.fieldMembers = fm
	fnl.fieldNames = fieldNames
	fm.memberList = fnl

	model := core.NewQStringListModel2(fieldNames, nil)
	view.SetModel(model)
	fnl.model = model
	w := parent.window
	records := w.mainWindow.codeplug.Records(w.recordType)
	r := records[w.recordList.Current()]
	fnl.initMemberModel(model, r, fm.fType)
	view.Viewport().SetAcceptDrops(true)
	view.SetDragDropMode(widgets.QAbstractItemView__DragDrop)
	view.SetSelectionMode(widgets.QAbstractItemView__ExtendedSelection)
	view.SetMinimumWidth(view.SizeHintForColumn(0) * 6 / 5)
	view.SetMaximumWidth(view.SizeHintForColumn(0) * 6 / 5)
	view.SetDefaultDropAction(core.Qt__MoveAction)
	view.SetAcceptDrops(true)
	view.SetDropIndicatorShown(true)
	view.SetDragEnabled(true)
	view.SetSelectionBehavior(widgets.QAbstractItemView__SelectRows)

	parent.layout.AddWidget(view, 0, 0)

	return fnl
}

func addAvailableList(parent *VBox, fm *FieldMembers, fieldNames []string) *fieldNameList {
	fnl := new(fieldNameList)
	view := widgets.NewQListView(nil)
	fnl.qListView = view
	fnl.fieldMembers = fm
	fnl.fieldNames = fieldNames
	fm.availableList = fnl

	model := core.NewQStringListModel2(fieldNames, nil)
	view.SetModel(model)
	fnl.model = model
	w := parent.window
	records := w.mainWindow.codeplug.Records(w.recordType)
	r := records[w.recordList.Current()]
	fnl.initAvailableModel(model, r, fm.fType)
	view.Viewport().SetAcceptDrops(true)
	view.SetDragDropMode(widgets.QAbstractItemView__DragDrop)
	view.SetSelectionMode(widgets.QAbstractItemView__ExtendedSelection)
	view.SetMinimumWidth(view.SizeHintForColumn(0) * 6 / 5)
	view.SetMaximumWidth(view.SizeHintForColumn(0) * 6 / 5)
	view.SetDefaultDropAction(core.Qt__MoveAction)
	view.SetAcceptDrops(true)
	view.SetDropIndicatorShown(true)
	view.SetDragEnabled(true)
	view.SetSelectionBehavior(widgets.QAbstractItemView__SelectRows)

	parent.layout.AddWidget(view, 0, 0)

	return fnl
}

func (fnl *fieldNameList) newNames(names []string) {
	model := core.NewQStringListModel2(names, nil)
	fnl.qListView.SetModel(model)
	fnl.model = model
}

func (fnl *fieldNameList) initMemberModel(model *core.QStringListModel, r *codeplug.Record, fType codeplug.FieldType) {
	fnl.initFieldNameModel(model, r, fType)

	model.ConnectRowCount(func(parent *core.QModelIndex) int {
		return len(r.Fields(fType))
	})

	model.ConnectCanDropMimeData(func(data *core.QMimeData, action core.Qt__DropAction, row int, column int, parent *core.QModelIndex) bool {
		mimeType := "application/x.codeplug.fieldname.list"
		if !data.HasFormat(mimeType) {
			return false
		}

		if column > 0 {
			return false
		}

		id, _, _, err := fnl.parseMimeData(mimeType, data)
		if err != nil {
			return false
		}

		if id != fnl.fieldMembers.id {
			return false
		}

		return true
	})

	model.ConnectDropMimeData(func(data *core.QMimeData, action core.Qt__DropAction, row int, column int, parent *core.QModelIndex) bool {
		if !model.CanDropMimeData(data, action, row, column, parent) {
			return false
		}

		if action == core.Qt__IgnoreAction {
			return true
		}

		var dRow int
		if row != -1 {
			dRow = row
		} else if parent.IsValid() {
			dRow = parent.Row()
		} else {
			dRow = model.RowCount(nil)
		}

		mimeT := "application/x.codeplug.fieldname.list"
		id, sfl, fieldNames, err := fnl.parseMimeData(mimeT, data)
		if err != nil {
			WarningPopup("Drop Error", err.Error())
			return false
		}

		fields := make([]*codeplug.Field, len(fieldNames))
		for i, name := range fieldNames {
			f, err := r.NewFieldWithValue(fType, 0, name)
			if err != nil {
				WarningPopup("Drop Error", err.Error())
				return false
			}

			fields[i] = f
		}

		sParent := core.NewQModelIndex()
		dParent := core.NewQModelIndex()
		rv := true
	actionSwitch:
		switch {
		case action == core.Qt__MoveAction && id != fnl.fieldMembers.id:
			return false

		case action == core.Qt__MoveAction && sfl == fnl:
			for _, f := range fields {
				f2 := r.FindFieldByName(f.Type(), f.String())
				f.SetIndex(f2.Index())
			}
			change := r.MoveFieldsChange(fields)
			for i, f := range fields {
				f := r.FindFieldByName(f.Type(), f.String())
				sRow := f.Index()
				model.MoveRows(sParent, sRow, 1, dParent, dRow)
				if sRow < dRow {
					dRow--
				}
				// publish the moved field
				fields[i] = f
				dRow++
			}
			change.Complete()

		default:
			change := r.InsertFieldsChange(fields)
			for _, f := range fields {
				fnl.fieldToInsert = f
				rv = model.InsertRows(dRow, 1, dParent)
				if !rv {
					break actionSwitch
				}
				dRow++
			}
			change.Complete()
		}

		return rv
	})

	model.ConnectInsertRows(func(row int, count int, parent *core.QModelIndex) bool {
		f := fnl.fieldToInsert
		if count != 1 {
			log.Fatal("bad insert fields count")
		}

		if len(r.Fields(fType)) >= r.MaxFields(fType) {
			WarningPopup("Insert Fields", "too many fields")
			return false
		}

		f.SetIndex(row)

		model.BeginInsertRows(core.NewQModelIndex(), row, row)
		err := r.InsertField(f)
		model.EndInsertRows()

		if err != nil {
			log.Fatal("ConnectInsertRows: InsertField failed")
		}

		fnl.fieldToInsert = nil

		return true
	})
}

func (fnl *fieldNameList) initAvailableModel(model *core.QStringListModel, r *codeplug.Record, fType codeplug.FieldType) {
	fnl.initFieldNameModel(model, r, fType)

	model.ConnectCanDropMimeData(func(data *core.QMimeData, action core.Qt__DropAction, row int, column int, parent *core.QModelIndex) bool {
		mimeType := "application/x.codeplug.fieldname.list"
		if !data.HasFormat(mimeType) {
			return false
		}

		if column > 0 {
			return false
		}

		id, sfl, _, err := fnl.parseMimeData(mimeType, data)
		if err != nil {
			return false
		}

		fm := fnl.fieldMembers
		if id != fm.id || sfl != fm.memberList {
			return false
		}

		return true
	})

	model.ConnectDropMimeData(func(data *core.QMimeData, action core.Qt__DropAction, row int, column int, parent *core.QModelIndex) bool {
		if !model.CanDropMimeData(data, action, row, column, parent) {
			return false
		}

		if action == core.Qt__IgnoreAction {
			return true
		}

		mimeT := "application/x.codeplug.fieldname.list"
		id, sfl, fieldNames, err := fnl.parseMimeData(mimeT, data)
		if err != nil {
			WarningPopup("Drop Error", err.Error())
			return false
		}

		if id != fnl.fieldMembers.id || sfl == fnl {
			return false
		}

		if action != core.Qt__MoveAction {
			return false
		}

		fm := sfl.fieldMembers
		r := fm.record

		fields := make([]*codeplug.Field, len(fieldNames))
		for i, name := range fieldNames {
			f, err := r.NewFieldWithValue(fType, 0, name)
			if err != nil {
				WarningPopup("Drop Error", err.Error())
				return false
			}
			fields[i] = r.FindFieldByName(f.Type(), f.String())
		}

		model := sfl.qListView.Model()
		qModelIndex := core.NewQModelIndex()

		change := r.RemoveFieldsChange(fields)
		for _, f := range fields {
			row = f.Index()
			model.BeginRemoveRows(qModelIndex, row, row)
			r.RemoveField(f)
			model.EndRemoveRows()
		}
		change.Complete()
		return true
	})
}

func (fnl *fieldNameList) parseMimeData(mimeType string, data *core.QMimeData) (id string, rfnl *fieldNameList, fieldNames []string, err error) {
	str := data.Data(mimeType).Data()
	reader := bufio.NewReader(strings.NewReader(str))
	id, err = reader.ReadString('\n')
	if err != nil {
		msg := fmt.Sprintf("Data read error 1: %s", err.Error())
		return id, rfnl, fieldNames, fmt.Errorf(msg)
	}

	id = strings.TrimSuffix(id, "\n")
	fieldNames = make([]string, 0, 10)
	for {
		name, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			msg := fmt.Sprintf("Data read error 2: %s", err.Error())
			return id, rfnl, fieldNames, fmt.Errorf(msg)
		}
		name = strings.TrimSuffix(name, "\n")
		fieldNames = append(fieldNames, name)
	}

	if len(fieldNames) == 0 {
		err = fmt.Errorf("no fields received")
		return id, rfnl, fieldNames, err
	}

	otherList := fnl.fieldMembers.availableList
	if fnl == fnl.fieldMembers.availableList {
		otherList = fnl.fieldMembers.memberList
	}
	lists := []*fieldNameList{fnl, otherList}

	for _, rfnl = range lists {
		for _, str := range rfnl.fieldNames {
			if str == fieldNames[0] {
				return id, rfnl, fieldNames, err
			}
		}
	}
	log.Fatal("field name not found")

	return id, rfnl, fieldNames, err
}

func (fnl *fieldNameList) initFieldNameModel(model *core.QStringListModel, r *codeplug.Record, fType codeplug.FieldType) {
	model.ConnectSupportedDropActions(func() core.Qt__DropAction {
		return core.Qt__MoveAction
	})

	model.ConnectMoveRows(func(sParent *core.QModelIndex, sRow int, count int, dParent *core.QModelIndex, dRow int) bool {
		if dRow == sRow {
			return true
		}
		if count != 1 {
			log.Fatal("ConnectMoveRows count != 1")
		}
		fType := fnl.fieldMembers.fType
		allFields := r.Fields(fType)

		f := allFields[sRow]
		model.BeginMoveRows(sParent, sRow, sRow, dParent, dRow)
		r.MoveField(dRow, f)
		model.EndMoveRows()
		return true
	})

	model.ConnectRemoveRows(func(row int, count int, parent *core.QModelIndex) bool {
		return true
	})

	model.ConnectFlags(func(index *core.QModelIndex) core.Qt__ItemFlag {
		flags := core.Qt__ItemIsSelectable |
			core.Qt__ItemIsEnabled |
			core.Qt__ItemIsDragEnabled |
			core.Qt__ItemIsDropEnabled
		return flags
	})

	model.ConnectMimeTypes(func() []string {
		return []string{"application/x.codeplug.fieldname.list"}
	})

	model.ConnectMimeData(func(indexes []*core.QModelIndex) *core.QMimeData {
		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)

		fmt.Fprintln(writer, fnl.fieldMembers.id)
		for _, index := range indexes {
			fmt.Fprintln(writer, fnl.fieldNames[index.Row()])
		}
		writer.Flush()

		str := buf.String()
		byteArray := core.NewQByteArray2(str, len(str))
		mimeData := core.NewQMimeData()
		mimeData.SetData("application/x.codeplug.fieldname.list", byteArray)
		return mimeData
	})
}

type FieldMembers struct {
	memberList    *fieldNameList
	availableList *fieldNameList
	record        *codeplug.Record
	fType         codeplug.FieldType
	id            string
}

func (vBox *VBox) AddFieldMembers(r *codeplug.Record, sortAvailable *bool, nameType codeplug.FieldType, memberType codeplug.FieldType, name string) *FieldMembers {
	var err error
	fm := new(FieldMembers)
	fm.record = r
	fm.fType = memberType
	fm.id, err = randomString(64)
	if err != nil {
		log.Fatal("randomString failure")
	}

	listRecordType := r.NewField(memberType).ListRecordType()

	availableMap := make(map[string]bool)
	names := *r.Codeplug().Record(listRecordType).ListNames()
	for _, name := range names {
		availableMap[name] = true
	}

	fields := r.Fields(memberType)
	memberNames := make([]string, len(fields))
	for i, f := range fields {
		name := f.String()
		memberNames[i] = name
		delete(availableMap, name)
	}

	availableNames := make([]string, 0, len(availableMap))
	for _, name := range names {
		if availableMap[name] {
			availableNames = append(availableNames, name)
		}
	}
	names = availableNames
	if *sortAvailable {
		names = make([]string, len(availableNames))
		copy(names, availableNames)
		sort.Strings(names)
	}

	row := vBox.AddHbox()
	form := row.AddForm()
	form.AddFieldRows(r, nameType)

	row = vBox.AddHbox()

	group := row.AddGroupbox(name)
	groupColumn := group.AddVbox()
	groupRow := groupColumn.AddHbox()
	column := groupRow.AddVbox()
	column.AddLabel("Available")
	availableList := addAvailableList(column, fm, names)

	column = groupRow.AddVbox()
	column.AddFiller()
	add := column.AddButton("Add >>")
	delete := column.AddButton("<< Delete")
	column.AddFiller()

	column = groupRow.AddVbox()
	column.AddLabel("Members")
	membersList := addMemberList(column, fm, memberNames)

	width := availableList.Width()
	if membersList.Width() > width {
		width = membersList.Width()
	}
	availableList.SetWidth(width)
	membersList.SetWidth(width)

	groupRow = groupColumn.AddHbox()
	form = groupRow.AddForm()
	label := widgets.NewQLabel2("Sort", nil, 0)
	checkBox := widgets.NewQCheckBox(nil)
	checkState := core.Qt__Unchecked
	if *sortAvailable {
		checkState = core.Qt__Checked
	}
	checkBox.SetCheckState(checkState)
	checkBox.ConnectClicked(func(checked bool) {
		*sortAvailable = checked
		names := make([]string, len(availableNames))
		copy(names, availableNames)
		if *sortAvailable {
			sort.Strings(names)
		}
		availableList.newNames(names)
	})
	form.layout.AddRow(label, checkBox)

	add.ConnectClicked(func() {
		allFields := r.Fields(memberType)
		names := availableList.selectedNames()
		if len(names) == 0 {
			WarningPopup("Add Member", "no members selected")
			return
		}

		if len(allFields)+len(names) > r.MaxFields(memberType) {
			WarningPopup("Add Member", "too many fields")
			return
		}

		fields := make([]*codeplug.Field, len(names))
		for i, name := range names {
			f, _ := r.NewFieldWithValue(memberType, 0, name)
			fields[i] = f
		}

		model := membersList.model
		row := len(allFields)

		change := r.InsertFieldsChange(fields)
		for _, f := range fields {
			f.SetIndex(row)
			model.BeginInsertRows(core.NewQModelIndex(), row, row)
			err := r.InsertField(f)
			model.EndInsertRows()

			if err != nil {
				WarningPopup("Add Member(s)", err.Error())
				break
			}
			row++
		}
		change.Complete()
	})

	delete.ConnectClicked(func() {
		names := membersList.selectedNames()
		if len(names) == 0 {
			WarningPopup("Delete Member", "no members selected")
			return
		}

		fields := make([]*codeplug.Field, len(names))
		for i, name := range names {
			f, _ := r.NewFieldWithValue(memberType, 0, name)
			fields[i] = r.FindFieldByName(f.Type(), f.String())
		}

		model := membersList.model
		qModelIndex := core.NewQModelIndex()

		change := r.RemoveFieldsChange(fields)
		for _, f := range fields {
			row := f.Index()
			model.BeginRemoveRows(qModelIndex, row, row)
			r.RemoveField(f)
			model.EndRemoveRows()
		}
		change.Complete()
	})

	return fm
}
