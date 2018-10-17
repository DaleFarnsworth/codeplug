// Copyright 2017-2018 Dale Farnsworth. All rights reserved.

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
	"math"
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
	view.SetMinimumWidth(view.SizeHintForColumn(0) + 20)
	view.SetMaximumWidth(view.SizeHintForColumn(0) + 20)
	view.SetDefaultDropAction(core.Qt__MoveAction)
	view.SetAcceptDrops(true)
	view.SetDropIndicatorShown(true)
	view.SetDragEnabled(true)
	view.SetSelectionBehavior(widgets.QAbstractItemView__SelectRows)

	rl.qListView.ConnectCurrentChanged(func(selected *core.QModelIndex, deSelected *core.QModelIndex) {
		w := rl.window
		if w.recordFunc != nil && !w.updating {
			DelayedCall(func() {
				dprint("ConnectCurrentChanged")
				w.recordFunc()
			})
		}
	})

	parent.layout.AddWidget(view, 0, 0)

	return rl
}

func (rl *RecordList) Model() *core.QAbstractItemModel {
	return rl.qListView.Model()
}

func (rl *RecordList) SelectionModel() *core.QItemSelectionModel {
	return rl.qListView.SelectionModel()
}

func (rl *RecordList) SetCurrent(i int) {
	index := rl.Model().CreateIndex(i, 0, nil)
	rl.qListView.ScrollTo(index, widgets.QAbstractItemView__EnsureVisible)
	rl.qListView.SetCurrentIndex(index)
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
	topLeft := rl.Model().CreateIndex(0, 0, nil)
	bottomRight := rl.Model().CreateIndex(rowCount-1, 0, nil)
	rl.qListView.DataChanged(topLeft, bottomRight, []int{})
}

func (rl *RecordList) AddDupSelected(add bool) error {
	w := rl.window
	rType := w.recordType
	cp := w.mainWindow.codeplug

	allRecords := cp.Records(rType)
	records := rl.SelectedRecords()
	if len(records) == 0 {
		records = []*codeplug.Record{allRecords[len(allRecords)-1]}
	}

	if len(allRecords)+len(records) > cp.MaxRecords(rType) {
		return fmt.Errorf("too many records")
	}

	model := rl.Model()
	qModelIndex := core.NewQModelIndex()

	row := records[len(records)-1].Index() + 1
	if add {
		row = len(cp.Records(rType))
	}

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
	if add {
		rl.SetCurrent(len(cp.Records(rType)) - 1)
	}
	w.recordFunc()

	return nil
}
func (rl *RecordList) AddSelected() error {
	add := true
	return rl.AddDupSelected(add)
}
func (rl *RecordList) DupSelected() error {
	add := false
	return rl.AddDupSelected(add)
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

	model := rl.Model()
	qModelIndex := core.NewQModelIndex()

	change := cp.RemoveRecordsChange(records)
	for _, r := range records {
		row := r.Index()
		model.BeginRemoveRows(qModelIndex, row, row)
		cp.RemoveRecord(r)
		model.EndRemoveRows()
	}
	change.Complete()

	allRecords = cp.Records(w.recordType)
	row := records[len(records)-1].Index()
	if row >= len(allRecords) {
		row = len(allRecords) - 1
	}
	rl.SetCurrent(row)
	rl.SelectRecords(allRecords[row])

	rl.Update()
	w.recordFunc()

	return nil
}

func (rl *RecordList) SelectRecords(records ...*codeplug.Record) {
	rl.ClearSelection()
	for _, r := range records {
		index := rl.Model().CreateIndex(r.Index(), 0, nil)
		rl.SelectionModel().Select(index, core.QItemSelectionModel__Select)
	}
}

func (w *Window) dataRecords(data *core.QMimeData, drop bool) (records []*codeplug.Record, depRecords []*codeplug.Record, id string, err error) {
	str := data.Data("application/x.codeplug.record.list").Data()
	reader := bufio.NewReader(strings.NewReader(str))

	var line string
	line, err = reader.ReadString('\n')
	if err != nil {
		err := fmt.Errorf("Data read error: %s", err.Error())
		return nil, nil, "", err
	}
	id = strings.TrimSuffix(line, "\n")

	line, err = reader.ReadString('\n')
	if err != nil {
		err := fmt.Errorf("Data read error: %s", err.Error())
		return nil, nil, "", err
	}
	rType := codeplug.RecordType(strings.TrimSuffix(line, "\n"))

	if rType != w.recordType {
		err := fmt.Errorf("Wrong data type")
		return nil, nil, "", err
	}

	line, err = reader.ReadString('\n')
	if err != nil {
		err := fmt.Errorf("Data read error: %s", err.Error())
		return nil, nil, "", err
	}

	strs := strings.Fields(strings.TrimSuffix(line, "\n"))
	cpType := strs[0]
	models := strs[1:]

	cp := w.mainWindow.codeplug

	if !drop && id != cp.ID() {
		return nil, nil, id, nil
	}

	deferValues := true
	records, depRecords, err = w.mainWindow.codeplug.ParseRecords(reader, deferValues)
	if err != nil && drop {
		for _, cpModel := range cp.Models() {
			for _, model := range models {
				if model != cpModel {
					continue
				}
				err := fmt.Errorf("syntax error: %s", err.Error())
				return nil, nil, "", err
			}
		}
		err = fmt.Errorf("Cannot copy from model %s to model %s", cpType, cp.Type())
	}

	if len(records) == 0 && drop {
		err := fmt.Errorf("no new %s", cp.RecordTypeName(rType))
		return nil, nil, "", err
	}

	return records, depRecords, id, nil
}

func (w *Window) initRecordModel(writable bool) {
	record := w.record()
	if record.MaxRecords() == 1 {
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

	if !writable {
		return
	}

	model.ConnectMoveRows(func(sParent *core.QModelIndex, sRow int, count int, dParent *core.QModelIndex, dRow int) bool {
		if count != 1 {
			logFatal("ConnectMoveRows: count != 1")
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
			logFatal("bad insert records count")
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
			logFatal("ConnectInsertRows: InsertRecord failed")
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

		drop := false
		records, _, id, err := w.dataRecords(data, drop)
		if err != nil {
			return false
		}

		cp := w.mainWindow.codeplug
		if id != cp.ID() {
			return true
		}

		firstSRow := math.MaxInt32
		lastSRow := -1
		for _, r := range records {
			r := cp.FindRecordByName(r.Type(), r.Name())
			sRow := 0
			if r != nil {
				sRow = r.Index()
			}
			if sRow < firstSRow {
				firstSRow = sRow
			}
			if sRow > lastSRow {
				lastSRow = sRow
			}
		}

		var dRow int
		if row != -1 {
			dRow = row
		} else if parent.IsValid() {
			dRow = parent.Row()
		} else {
			dRow = model.RowCount(nil)
		}

		if dRow >= firstSRow && dRow <= lastSRow+1 {
			return false
		}

		return true
	})

	model.ConnectDropMimeData(func(data *core.QMimeData, action core.Qt__DropAction, row int, column int, parent *core.QModelIndex) bool {
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

		drop := true
		records, depRecords, id, err := w.dataRecords(data, drop)
		if err != nil {
			InfoPopup("Data Drop", err.Error())
			return false
		}

		sParent := core.NewQModelIndex()
		dParent := core.NewQModelIndex()

		cp := w.mainWindow.codeplug
		rv := true
	actionSwitch:
		switch {
		case action == core.Qt__MoveAction && id == cp.ID():
			w.updating = true
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
			w.updating = false

		case action == core.Qt__CopyAction && id == cp.ID():
			w.updating = true
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
			w.updating = false

		case id != cp.ID():
			rTypeString := "Records"
			if len(records) > 0 {
				rTypeString = string(records[0].Type())
			}

			if action == core.Qt__MoveAction {
				newRecords := make([]*codeplug.Record, 0)
				for _, r := range records {
					if r.NameExists() {
						continue
					}
					newRecords = append(newRecords, r)
				}
				records = newRecords
			}

			if len(records) == 0 {
				msg := fmt.Sprintf("no new %s", rTypeString)
				InfoPopup("Data Drop", msg)
				return false
			}

			w.updating = true
			change := cp.InsertRecordsChange(records)
			for _, r := range depRecords {
				if r.NameExists() {
					continue
				}
				records := []*codeplug.Record{r}
				subChange := cp.InsertRecordsChange(records)
				change.AddChange(subChange)
				cp.AppendRecord(r)
			}

			for _, r := range records {
				w.recordList.recordToInsert = r
				rv = model.InsertRows(dRow, 1, dParent)
				if !rv {
					break actionSwitch
				}
				dRow++
			}

			cp.ResolveDeferredValueFields()
			change.Complete()
			w.updating = false
		}

		rl := w.recordList
		rl.SetCurrent(dRow - 1)
		rl.SelectRecords(records...)

		rl.Update()
		dprint("ConnectDropMimeData")
		w.recordFunc()

		if id != cp.ID() && !cp.Valid() {
			fmtStr := `
%d records with invalid field values were found in the codeplug.

Select "Menu->Edit->Show Invalid Fields" to view them.`
			msg := fmt.Sprintf(fmtStr, len(cp.Warnings()))
			DelayedCall(func() {
				InfoPopup("codeplug warning", msg)
			})
		}

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
		fmt.Fprintln(writer, w.recordType)

		fmt.Fprint(writer, cp.Type())
		for _, model := range cp.Models() {
			fmt.Fprintf(writer, " %s", model)
		}
		fmt.Fprintln(writer)

		records := make([]*codeplug.Record, 0)
		for _, index := range indexes {
			r := cp.Records(w.recordType)[index.Row()]
			records = append(records, r)
		}

		var depRecords []*codeplug.Record
		records, depRecords = codeplug.DependentRecords(records)
		for _, r := range depRecords {
			codeplug.PrintDependentRecordWithIndex(writer, r)
		}
		for _, r := range records {
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
