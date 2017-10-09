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
	"log"
	"strings"

	"github.com/dalefarnsworth/codeplug/codeplug"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type RecordList struct {
	window         *Window
	qListView      *widgets.QListView
	recordToInsert *codeplug.Record
}

func (parent *HBox) AddRecordList(rType codeplug.RecordType) *RecordList {
	rl := new(RecordList)
	rl.window = parent.window
	rl.window.recordList = rl
	view := widgets.NewQListView(nil)
	rl.qListView = view
	model := rl.window.recordModel
	view.SetModel(model)
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

	rl.qListView.ConnectCurrentChanged(func(selected *core.QModelIndex, deSelected *core.QModelIndex) {
		if rl.window.recordFunc != nil {
			rl.window.recordFunc()
		}
	})

	parent.layout.AddWidget(view, 0, 0)

	return rl
}

func (rl *RecordList) SetCurrent(i int) {
	index := rl.qListView.Model().CreateIndex(i, 0, nil)
	rl.qListView.SetCurrentIndex(index)
	rl.qListView.ScrollTo(index, widgets.QAbstractItemView__EnsureVisible)
}

func (rl *RecordList) Current() int {
	current := rl.qListView.CurrentIndex().Row()
	records := rl.window.records()
	if current >= len(records) {
		current = len(records) - 1
		rl.SetCurrent(current)
	}

	return current
}

func (rl *RecordList) SelectedRecords() []*codeplug.Record {
	w := rl.window
	cp := w.mainWindow.codeplug
	allRecords := cp.Records(w.recordType)
	indexes := rl.qListView.SelectedIndexes()
	records := make([]*codeplug.Record, len(indexes))
	for i, index := range indexes {
		records[i] = allRecords[index.Row()]
	}
	return records
}

func (rl *RecordList) ClearSelection() {
	rl.qListView.ClearSelection()
}

func (rl *RecordList) Update() {
	rowCount := len(*rl.window.record().ListNames())
	row := rl.Current()
	if row < 0 {
		rl.SetCurrent(0)
	} else if row >= rowCount {
		rl.SetCurrent(rowCount - 1)
	}
	topLeft := rl.qListView.Model().CreateIndex(0, 0, nil)
	bottomRight := rl.qListView.Model().CreateIndex(rowCount-1, 0, nil)
	rl.qListView.DataChanged(topLeft, bottomRight, []int{})
}

func (rl *RecordList) AddSelected() error {
	w := rl.window
	rType := w.recordType
	cp := w.mainWindow.codeplug

	records := rl.SelectedRecords()
	if len(records) == 0 {
		return fmt.Errorf("no records selected")
	}

	if len(cp.Records(rType))+len(records) >= cp.MaxRecords(rType) {
		return fmt.Errorf("too many records")
	}

	model := rl.qListView.Model()
	qModelIndex := core.NewQModelIndex()

	row := len(cp.Records(rType))

	change := cp.InsertRecordsChange(records)
	for i, r := range records {
		r = r.Copy()
		records[i] = r
		r.SetIndex(row)

		model.BeginInsertRows(qModelIndex, row, row)
		cp.InsertRecord(r)
		model.EndInsertRows()
		row++
	}
	change.Complete()

	rl.Update()
	rl.SetCurrent(len(cp.Records(rType)) - 1)
	w.recordFunc()

	return nil
}

func (rl *RecordList) RemoveSelected() error {
	w := rl.window
	cp := w.mainWindow.codeplug
	allRecords := cp.Records(w.recordType)

	records := rl.SelectedRecords()
	if len(records) == 0 {
		return fmt.Errorf("no records selected")
	}

	if len(allRecords) <= len(records) {
		return fmt.Errorf("can't delete last record")
	}

	model := rl.qListView.Model()
	qModelIndex := core.NewQModelIndex()

	change := cp.RemoveRecordsChange(records)
	for _, r := range records {
		row := r.Index()

		model.BeginRemoveRows(qModelIndex, row, row)
		cp.RemoveRecord(r)
		model.EndRemoveRows()
	}
	change.Complete()

	rl.Update()
	w.recordFunc()

	return nil
}

func (w *Window) initRecordModel() {
	record := w.record()
	if record.NameField() == nil {
		return
	}
	model := core.NewQAbstractListModel(nil)
	w.recordModel = model

	model.ConnectRowCount(func(parent *core.QModelIndex) int {
		return len(*record.ListNames())
	})

	model.ConnectData(func(idx *core.QModelIndex, role int) *core.QVariant {
		row := idx.Row()
		if role == int(core.Qt__DisplayRole) && idx.IsValid() {
			names := *record.ListNames()
			if row >= 0 && row < len(names) {
				return core.NewQVariant14(names[row])
			}
		}

		return core.NewQVariant()
	})

	model.ConnectMoveRows(func(sParent *core.QModelIndex, sRow int, count int, dParent *core.QModelIndex, dRow int) bool {
		if count != 1 {
			log.Fatal("ConnectMoveRows: count != 1")
		}

		cp := w.mainWindow.codeplug
		r := cp.Records(w.recordType)[sRow]
		model.BeginMoveRows(sParent, sRow, sRow, dParent, dRow)
		cp.MoveRecord(dRow, r)
		model.EndMoveRows()

		return true
	})

	model.ConnectInsertRows(func(row int, count int, parent *core.QModelIndex) bool {
		rType := w.recordType
		rl := w.recordList
		cp := w.mainWindow.codeplug
		r := rl.recordToInsert
		if count != 1 {
			log.Fatal("bad insert records count")
		}

		if len(cp.Records(rType)) >= cp.MaxRecords(rType) {
			WarningPopup("Insert Records", "too many records")
			return false
		}

		r.SetIndex(row)

		model.BeginInsertRows(core.NewQModelIndex(), row, row)
		err := cp.InsertRecord(r)
		model.EndInsertRows()

		if err != nil {
			log.Fatal("ConnectInsertRows: InsertRecord failed")
		}

		rl.recordToInsert = nil
		return true
	})

	model.ConnectFlags(func(index *core.QModelIndex) core.Qt__ItemFlag {
		flags := core.Qt__ItemIsSelectable |
			core.Qt__ItemIsEnabled |
			core.Qt__ItemIsDragEnabled |
			core.Qt__ItemIsDropEnabled
		return flags
	})

	model.ConnectCanDropMimeData(func(data *core.QMimeData, action core.Qt__DropAction, row int, column int, parent *core.QModelIndex) bool {
		if !data.HasFormat("application/x.codeplug.record.list") {
			return false
		}

		if column > 0 {
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

		str := data.Data("application/x.codeplug.record.list").Data()
		reader := bufio.NewReader(strings.NewReader(str))
		id, err := reader.ReadString('\n')
		if err != nil {
			msg := fmt.Sprintf("Data read error: %s", err.Error())
			WarningPopup("Drop Error", msg)
			return false
		}

		id = strings.TrimSuffix(id, "\n")
		records, err := w.mainWindow.codeplug.ParseRecords(reader)
		if err != nil {
			msg := fmt.Sprintf("data format error: %s", err.Error())
			WarningPopup("Drop Error", msg)
			return false
		}

		if len(records) == 0 {
			WarningPopup("Drop Error", "No records received")
			return false
		}

		if records[0].Type() != w.recordType {
			WarningPopup("Drop Error", "Wrong data type")
			return false
		}

		sParent := core.NewQModelIndex()
		dParent := core.NewQModelIndex()

		cp := w.mainWindow.codeplug
		rv := true
	actionSwitch:
		switch {
		case action == core.Qt__MoveAction && id == cp.ID():
			change := cp.MoveRecordsChange(records)
			for i, r := range records {
				r := cp.FindRecordByName(r.Type(), r.Name())
				sRow := r.Index()

				if dRow == sRow || dRow == sRow+1 {
					continue
				}

				model.MoveRows(sParent, sRow, 1, dParent, dRow)
				if sRow < dRow {
					dRow--
				}
				// publish the moved record
				records[i] = r
				dRow++
			}
			change.Complete()

		default:
			change := cp.InsertRecordsChange(records)
			for _, r := range records {
				w.recordList.recordToInsert = r
				rv = model.InsertRows(dRow, 1, dParent)
				if !rv {
					break actionSwitch
				}
				dRow++
			}
			change.Complete()
		}

		w.recordList.Update()
		w.recordList.SetCurrent(dRow - 1)

		return rv
	})

	model.ConnectSupportedDropActions(func() core.Qt__DropAction {
		return core.Qt__CopyAction | core.Qt__MoveAction
	})

	model.ConnectMimeTypes(func() []string {
		return []string{"application/x.codeplug.record.list"}
	})

	model.ConnectMimeData(func(indexes []*core.QModelIndex) *core.QMimeData {
		var buf bytes.Buffer
		writer := bufio.NewWriter(&buf)

		cp := w.mainWindow.codeplug
		fmt.Fprintln(writer, cp.ID())
		for _, index := range indexes {
			r := cp.Records(w.recordType)[index.Row()]
			codeplug.PrintRecordWithIndex(writer, r)
		}
		writer.Flush()

		str := buf.String()
		byteArray := core.NewQByteArray2(str, len(str))
		mimeData := core.NewQMimeData()
		mimeData.SetData("application/x.codeplug.record.list", byteArray)
		return mimeData
	})
}
