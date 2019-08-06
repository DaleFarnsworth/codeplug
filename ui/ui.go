// Copyright 2017-2019 Dale Farnsworth. All rights reserved.

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
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dalefarnsworth/codeplug/codeplug"
	"github.com/dalefarnsworth/codeplug/debug"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var mainWindows = []*MainWindow{}

type App struct {
	qApp     widgets.QApplication
	filename string
}

func NewApp() (*App, error) {
	app := new(App)
	qApp := widgets.NewQApplication(len(os.Args), os.Args)

	if qApp == nil {
		return nil, fmt.Errorf("Not built with QT")
	}

	app.qApp = *qApp

	return app, nil
}

func (app *App) SetOrganizationName(str string) {
	app.qApp.SetOrganizationName(str)
}

func (app *App) SetOrganizationDomain(str string) {
	app.qApp.SetOrganizationDomain(str)
}

func (app *App) SetApplicationName(str string) {
	app.qApp.SetApplicationName(str)
}

func (app *App) NewSettings() *AppSettings {
	return &AppSettings{core.NewQSettings5(nil)}
}

func (app *App) Exec() {
	app.qApp.Exec()
}

func (app *App) ProcessEvents() {
	app.qApp.ProcessEvents(core.QEventLoop__AllEvents)
}

func (app *App) Quit() {
	app.qApp.Quit()
}

func ApplicationFontPointSize() int {
	font := widgets.QApplication_Font()
	return font.PointSize()
}

func SetApplicationFontPointSize(pointSize int) {
	font := widgets.QApplication_Font()
	font.SetPointSize(pointSize)
	widgets.QApplication_SetFont(font, "")
}

type AppSettings struct {
	qSettings *core.QSettings
}

func (as *AppSettings) Sync() {
	as.qSettings.Sync()
}

func (as *AppSettings) SetBool(str string, b bool) {
	as.qSettings.SetValue(str, core.NewQVariant9(b))
}

func (as *AppSettings) Bool(str string, deflt bool) bool {
	return as.qSettings.Value(str, core.NewQVariant9(deflt)).ToBool()
}

func (as *AppSettings) SetInt(str string, i int) {
	as.qSettings.SetValue(str, core.NewQVariant5(i))
}

func (as *AppSettings) Int(str string, deflt int) int {
	var ok *bool
	return as.qSettings.Value(str, core.NewQVariant5(deflt)).ToInt(ok)
}

func (as *AppSettings) SetString(str string, s string) {
	as.qSettings.SetValue(str, core.NewQVariant15(s))
}

func (as *AppSettings) String(str string, deflt string) string {
	return as.qSettings.Value(str, core.NewQVariant15(deflt)).ToString()
}

func (as *AppSettings) BeginWriteArray(prefix string, size int) {
	as.qSettings.BeginWriteArray(prefix, size)
}

func (as *AppSettings) BeginReadArray(prefix string) int {
	return as.qSettings.BeginReadArray(prefix)
}

func (as *AppSettings) EndArray() {
	as.qSettings.EndArray()
}

func (as *AppSettings) SetArrayIndex(i int) {
	as.qSettings.SetArrayIndex(i)
}

type DelayedCallStruct struct {
	core.QObject
	_ func() `slot:"create"`
}

func DelayedCall(f func()) {
	delayedCall := NewDelayedCallStruct(nil)
	delayedCall.ConnectCreate(f)
	go delayedCall.Create() // go routine, so it's called in event loop
}

type MainWindow struct {
	qMainWindow      widgets.QMainWindow
	codeplug         *codeplug.Codeplug
	recordWindows    map[codeplug.RecordType]*Window
	altRecordWindows map[codeplug.RecordType]*Window
	connectClose     func() bool
	connectChange    func(*codeplug.Change)
	changing         bool
}

func (mw *MainWindow) SetCodeplug(cp *codeplug.Codeplug) {
	mw.codeplug = cp
	mw.recordWindows = make(map[codeplug.RecordType]*Window)
	mw.altRecordWindows = make(map[codeplug.RecordType]*Window)

	mw.codeplug.ConnectChange(func(change *codeplug.Change) {
		mw := mainWindow(change.Codeplug())
		windows := make([]*Window, 0)
		for _, w := range mw.recordWindows {
			windows = append(windows, w)
		}
		for _, w := range mw.altRecordWindows {
			windows = append(windows, w)
		}

		for _, w := range windows {
			w.enableAllWidgets()
			if w.recordType == change.RecordType() {
				w.handleChange(change)
			}
		}

		mw.CodeplugChanged(change)
		if mw.connectChange != nil {
			mw.connectChange(change)
		}
	})
}

func NewMainWindow() *MainWindow {
	mw := new(MainWindow)
	mainWindows = append(mainWindows, mw)

	qmw := widgets.NewQMainWindow(nil, 0)
	mw.qMainWindow = *qmw

	qmw.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		if mw.connectClose != nil {
			if !mw.connectClose() {
				event.Ignore()
				return
			}

		}

		for _, w := range mw.recordWindows {
			w.Close()
		}
		for _, w := range mw.altRecordWindows {
			w.Close()
		}

		for i, mainWindow := range mainWindows {
			if mainWindow == mw {
				mainWindows = append(mainWindows[:i], mainWindows[i+1:]...)
				break
			}
		}

		event.Accept()
	})

	return mw
}

func (mW *MainWindow) Resize(w, h int) {
	mW.qMainWindow.Resize2(w, h)
}

func (parent *MainWindow) AddVbox() *VBox {
	box := NewVbox()
	box.SetContentsMargins(5, 5, 5, 5)
	parent.qMainWindow.SetCentralWidget(&box.qWidget)

	return box
}

func (parent *MainWindow) AddHbox() *HBox {
	box := NewHbox()
	box.SetContentsMargins(5, 5, 5, 5)
	parent.qMainWindow.SetCentralWidget(&box.qWidget)

	return box
}

func (mw *MainWindow) SetTitle(title string) {
	mw.qMainWindow.SetWindowTitle(title)
}

func (mw *MainWindow) Title() string {
	return mw.qMainWindow.WindowTitle()
}

func (mw *MainWindow) ConnectClose(fn func() bool) {
	mw.connectClose = fn
}

func (mw *MainWindow) Codeplug() *codeplug.Codeplug {
	return mw.codeplug
}

func (mw *MainWindow) ConnectChange(fn func(*codeplug.Change)) {
	mw.connectChange = fn
}

func MainWindows() []*MainWindow {
	return mainWindows
}

func (parent *Window) AddHbox() *HBox {
	box := NewHbox()
	parent.layout.AddWidget(&box.qWidget, 0, 0)
	box.window = parent
	return box
}

func (parent *Window) AddVbox() *VBox {
	box := NewVbox()
	parent.layout.AddWidget(&box.qWidget, 0, 0)
	box.window = parent
	return box
}

func (w *Window) SetTitle(title string) {
	w.qWidget.SetWindowTitle(title)
}

func (w *Window) Title() string {
	return w.qWidget.WindowTitle()
}

func (w *Window) AddMenuBar() *MenuBar {
	mb := new(MenuBar)
	mb.qMenuBar = widgets.NewQMenuBar(nil)
	w.layout.SetMenuBar(mb.qMenuBar)
	w.menuBar = mb
	return mb
}

func (w *Window) MenuBar() *MenuBar {
	return w.menuBar
}

func (enablerWidget *FieldWidget) enableWidgets() {
	enablerField := enablerWidget.field
	window := enablerWidget.window
	if window == nil {
		l.Fatal("nil window", enablerField.FullTypeName())
		return
	}
	for _, enablerType := range enablerField.Enables() {
		enableWidgetMap := window.widgets[enablerType]
		if enableWidgetMap == nil {
			// Not an error because some fields are not
			// displayed for some codeplug types, slso
			// some fields are not always enabled.
			continue
		}
		for enabledField, enabledWidget := range enableWidgetMap {
			if enablerField.Record() != enabledField.Record() {
				// Not an error because some fields are not
				// always enabled.
				continue
			}
			enabledWidget.enable(enablerWidget)
		}
	}
}

func (w *Window) enableWidgets() {
	for _, widget := range w.enablerWidgets {
		widget.enableWidgets()
	}
}

func (w *Window) AddEnablerWidget(widget *FieldWidget) {
	if w.enablerWidgets == nil {
		w.enablerWidgets = make([]*FieldWidget, 0)
	}
	w.enablerWidgets = append(w.enablerWidgets, widget)
}

func (w *Window) Show() {
	w.qWidget.Show()
	w.qWidget.ActivateWindow()
	w.qWidget.Raise()
	w.enableWidgets()
}

func clear(w Widget) {
	widget := w.qWidget_ITF().QWidget_PTR()

	for _, obj := range widget.Children() {
		if obj.Pointer() == widget.Layout().Pointer() {
			continue
		}
		obj.DeleteLater()
	}
}

func (w *Window) DeleteLater() {
	w.qWidget.DeleteLater()
}

func (w *Window) MainWindow() *MainWindow {
	return w.mainWindow
}

func (w *Window) RecordType() codeplug.RecordType {
	return w.recordType
}

func (w *Window) Records(rType codeplug.RecordType) []*codeplug.Record {
	mw := w.mainWindow
	if mw.recordWindows[rType] == nil {
		mw.altRecordWindows[rType] = w
		w.altRecordTypes[rType] = true
	}
	return mw.codeplug.Records(rType)
}

func (form *Form) RemoveWidget(widget *FieldWidget) {
	form.layout.RemoveRow2(widget.qWidget)
}

func (box *HBox) SetEnabled(enable bool) {
	box.qWidget.SetEnabled(enable)
}

func (box *VBox) SetEnabled(enable bool) {
	box.qWidget.SetEnabled(enable)
}

func (mw *MainWindow) Show() {
	mw.qMainWindow.Show()
	mw.qMainWindow.ActivateWindow()
	mw.qMainWindow.Raise()
}

func (mw *MainWindow) Close() {
	mw.qMainWindow.Close()
}

func (mw *MainWindow) BeginChange(change *codeplug.Change) {
	mw.changing = true
}

func (mw *MainWindow) EndChange(change *codeplug.Change) {
	mw.changing = false
	mw.CodeplugChanged(change)
}

func (mw *MainWindow) CodeplugChanged(change *codeplug.Change) {
	if mw == nil {
		return
	}

	if mw.changing {
		return
	}

	for _, w := range mw.recordWindows {
		rl := w.RecordList()
		if rl == nil {
			return
		}

		/*
			if w.recordType != change.RecordType() {
				rl.SetCurrent(0)
			}
		*/

		w.recordFunc()
		rl.Update()
	}
}

func (mw *MainWindow) RecordWindows() map[codeplug.RecordType]*Window {
	return mw.recordWindows
}

type Window struct {
	qWidget         widgets.QWidget
	layout          *widgets.QHBoxLayout
	menuBar         *MenuBar
	window          *Window
	mainWindow      *MainWindow
	recordType      codeplug.RecordType
	altRecordTypes  map[codeplug.RecordType]bool
	recordFunc      func()
	widgets         map[codeplug.FieldType]map[*codeplug.Field]*FieldWidget
	enablerWidgets  []*FieldWidget
	recordModel     *core.QAbstractListModel
	recordList      *RecordList
	connectClose    func() bool
	handleChange    func(*codeplug.Change)
	settingMultiple bool
}

func (mw *MainWindow) NewWindow() *Window {
	w := new(Window)
	w.qWidget = *widgets.NewQWidget(&mw.qMainWindow, core.Qt__Window)
	w.layout = widgets.NewQHBoxLayout2(&w.qWidget)
	w.mainWindow = mw
	w.window = w
	w.qWidget.Resize2(500, 500)

	w.qWidget.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		if w.connectClose != nil {
			if !w.connectClose() {
				event.Ignore()
				return
			}

		}
		event.Accept()
	})

	return w
}

func (mw *MainWindow) NewRecordWindow(rType codeplug.RecordType, writable bool) *Window {
	w := new(Window)
	mw.recordWindows[rType] = w
	w.qWidget = *widgets.NewQWidget(&mw.qMainWindow, core.Qt__Window)
	w.layout = widgets.NewQHBoxLayout2(&w.qWidget)
	w.mainWindow = mw
	w.window = w
	w.recordType = rType
	w.altRecordTypes = make(map[codeplug.RecordType]bool)
	w.widgets = make(map[codeplug.FieldType]map[*codeplug.Field]*FieldWidget)

	w.initRecordModel(writable)

	w.qWidget.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		if w.connectClose != nil {
			if !w.connectClose() {
				event.Ignore()
				return
			}

		}
		delete(mw.recordWindows, rType)
		delete(mw.altRecordWindows, rType)
		event.Accept()
	})

	w.handleChange = func(change *codeplug.Change) {
		changeType := change.Type()
		switch changeType {
		case codeplug.FieldChange:
			f := change.Field()
			if f == f.Record().NameField() {
				codeplug.NameFieldChanged(change)
			}

		case codeplug.RecordsFieldChange:
			changes := change.Changes()
			for _, change := range changes {
				w.handleChange(change)
			}

		case codeplug.RemoveRecordsChange:
			codeplug.RecordsRemoved(change)
		}
	}

	return w
}

func (w *Window) fieldWidget(f *codeplug.Field) *FieldWidget {
	widgets := w.widgets[f.Type()]
	if widgets == nil {
		return nil
	}

	return widgets[f]
}

func (w *Window) qWidget_ITF() widgets.QWidget_ITF {
	return &w.qWidget
}

func (w *Window) Window() *Window {
	return w.window
}

func (w *Window) AddWidget(widget Widget) {
	l.Fatal("cannot add widget to Window")
}

func (w *Window) SetContentsMargins(left int, right int, top int, bottom int) {
	l.Fatal("cannot set contents margin of Window")
}

func (w *Window) records() []*codeplug.Record {
	return w.mainWindow.codeplug.Records(w.recordType)
}

func (w *Window) record() *codeplug.Record {
	return w.mainWindow.codeplug.Record(w.recordType)
}

func (w *Window) Close() {
	w.qWidget.Close()
}

func (w *Window) ConnectClose(fn func() bool) {
	w.connectClose = fn
}

func (w *Window) RecordList() *RecordList {
	return w.recordList
}

type VBox struct {
	qWidget widgets.QWidget
	layout  *widgets.QVBoxLayout
	window  *Window
}

func NewVbox() *VBox {
	box := new(VBox)

	box.qWidget = *widgets.NewQWidget(nil, 0)
	box.layout = widgets.NewQVBoxLayout2(&box.qWidget)
	box.layout.SetContentsMargins(0, 0, 0, 0)

	return box
}

func (vBox *VBox) qWidget_ITF() widgets.QWidget_ITF {
	return &vBox.qWidget
}

func (vBox *VBox) Clear() {
	clear(vBox)
}

func (vBox *VBox) Window() *Window {
	return vBox.window
}

func (vBox *VBox) setWindow(window *Window) {
	vBox.window = window
}

func (parent *VBox) AddGroupbox(label string) *HBox {
	qgb := widgets.NewQGroupBox2(label, nil)
	layout := widgets.NewQHBoxLayout2(qgb)
	layout.SetContentsMargins(0, 0, 0, 0)

	box := NewHbox()
	layout.AddWidget(&box.qWidget, 0, 0)
	box.layout.SetContentsMargins(0, 0, 0, 0)
	parent.layout.AddWidget(qgb, 0, 0)
	box.window = parent.window

	return box
}

func (parent *HBox) AddGroupbox(label string) *HBox {
	qgb := widgets.NewQGroupBox2(label, nil)
	layout := widgets.NewQHBoxLayout2(qgb)
	layout.SetContentsMargins(0, 0, 0, 0)

	box := NewHbox()
	box.layout.SetContentsMargins(0, 0, 0, 0)
	layout.AddWidget(&box.qWidget, 0, 0)

	parent.layout.AddWidget(qgb, 0, 0)
	box.window = parent.window

	return box
}

func (parent *VBox) AddHbox() *HBox {
	box := NewHbox()
	parent.AddWidget(box)
	return box
}

func (parent *VBox) AddExistingHbox(box *HBox) *HBox {
	parent.AddWidget(box)
	return box
}

func (parent *HBox) AddHbox() *HBox {
	box := NewHbox()
	parent.AddWidget(box)
	return box
}

func (parent *HBox) AddExistingHbox(box *HBox) *HBox {
	parent.AddWidget(box)
	return box
}

func (parent *VBox) AddWidget(widget Widget) {
	parent.layout.AddWidget(widget.qWidget_ITF(), 0, 0)
	widget.setWindow(parent.Window())
}

func (parent *HBox) AddWidget(widget Widget) {
	parent.layout.AddWidget(widget.qWidget_ITF(), 0, 0)
	widget.setWindow(parent.Window())
}

func (vBox *VBox) SetContentsMargins(left int, right int, top int, bottom int) {
	vBox.layout.SetContentsMargins(left, right, top, bottom)
}

func (hBox *HBox) SetContentsMargins(left int, right int, top int, bottom int) {
	hBox.layout.SetContentsMargins(left, right, top, bottom)
}

func (parent *HBox) AddForm() *Form {
	form := NewForm()
	parent.AddWidget(form)
	return form
}

func NewForm() *Form {
	form := new(Form)

	form.qWidget = *widgets.NewQWidget(nil, 0)
	form.layout = widgets.NewQFormLayout(&form.qWidget)
	form.layout.SetContentsMargins(0, 0, 0, 0)
	// form.layout.SetSpacing(6)

	return form
}

func (parent *VBox) AddForm() *Form {
	form := NewForm()
	parent.AddWidget(form)
	return form
}

type Form struct {
	qWidget widgets.QWidget
	layout  *widgets.QFormLayout
	window  *Window
}

type HBox struct {
	qWidget widgets.QWidget
	layout  *widgets.QHBoxLayout
	window  *Window
}

func NewHbox() *HBox {
	box := new(HBox)

	box.qWidget = *widgets.NewQWidget(nil, 0)
	box.layout = widgets.NewQHBoxLayout2(&box.qWidget)
	box.layout.SetContentsMargins(0, 0, 0, 0)

	return box
}

func (hBox *HBox) Clear() {
	clear(hBox)
}

func (hBox *HBox) qWidget_ITF() widgets.QWidget_ITF {
	return &hBox.qWidget
}

func (hBox *HBox) Window() *Window {
	return hBox.window
}

func (hBox *HBox) setWindow(window *Window) {
	hBox.window = window
}

func (parent *HBox) AddVbox() *VBox {
	box := NewVbox()
	parent.AddWidget(box)
	return box
}

func (parent *HBox) AddExistingVbox(box *VBox) *VBox {
	parent.AddWidget(box)
	return box
}

func (parent *VBox) AddVbox() *VBox {
	box := NewVbox()
	parent.AddWidget(box)
	return box
}

func (parent *VBox) AddExistingVbox(box *VBox) *VBox {
	parent.AddWidget(box)
	return box
}

func (parent *HBox) AddButton(text string) *Button {
	b := NewButton(text)
	parent.layout.AddWidget(&b.qWidget, 0, 0)
	return b
}

func (parent *VBox) AddButton(text string) *Button {
	b := NewButton(text)
	parent.layout.AddWidget(&b.qWidget, 0, 0)
	return b
}

func (hBox *HBox) SetExpand() {
	hPolicy := widgets.QSizePolicy__Expanding
	vPolicy := widgets.QSizePolicy__Expanding
	hBox.qWidget.SetSizePolicy2(hPolicy, vPolicy)
}

func (hBox *HBox) AddSeparator() {
	frame := widgets.NewQFrame(nil, core.Qt__Widget)
	frame.SetFrameShape(widgets.QFrame__VLine)
	frame.SetFrameShadow(widgets.QFrame__Plain)
	hBox.layout.AddWidget(frame, 0, 0)
}

func (vBox *VBox) SetExpand() {
	hPolicy := widgets.QSizePolicy__Expanding
	vPolicy := widgets.QSizePolicy__Expanding
	vBox.qWidget.SetSizePolicy2(hPolicy, vPolicy)
}

func (vBox *VBox) AddSeparator() {
	frame := widgets.NewQFrame(nil, core.Qt__Widget)
	frame.SetFrameShape(widgets.QFrame__HLine)
	frame.SetFrameShadow(widgets.QFrame__Plain)
	vBox.layout.AddWidget(frame, 0, 0)
}

func (parent *HBox) AddLabel(str string) {
	qLabel := widgets.NewQLabel2(str, nil, 0)
	parent.layout.AddWidget(qLabel, 0, 0)
}

func (parent *VBox) AddLabel(str string) {
	qLabel := widgets.NewQLabel2(str, nil, 0)
	parent.layout.AddWidget(qLabel, 0, 0)
}

func (parent *HBox) AddSpace(width int) {
	w := gui.NewQFontMetrics(gui.NewQFont()).AverageCharWidth() * width
	h := 0
	hPolicy := widgets.QSizePolicy__Fixed
	vPolicy := widgets.QSizePolicy__Fixed
	spacer := widgets.NewQSpacerItem(w, h, hPolicy, vPolicy)
	parent.layout.AddItem(spacer)
}

func (parent *VBox) AddSpace(height int) {
	w := 0
	h := gui.NewQFontMetrics(gui.NewQFont()).AverageCharWidth() * height
	hPolicy := widgets.QSizePolicy__Fixed
	vPolicy := widgets.QSizePolicy__Fixed
	spacer := widgets.NewQSpacerItem(w, h, hPolicy, vPolicy)
	parent.layout.AddItem(spacer)
}

func (parent *HBox) AddFiller() {
	w := 0
	h := 0
	hPolicy := widgets.QSizePolicy__Expanding
	vPolicy := widgets.QSizePolicy__Expanding
	for i := 0; i < 30; i++ {
		filler := widgets.NewQSpacerItem(w, h, hPolicy, vPolicy)
		parent.layout.AddItem(filler)
	}
}

func (parent *VBox) AddFiller() {
	w := 0
	h := 0
	hPolicy := widgets.QSizePolicy__Expanding
	vPolicy := widgets.QSizePolicy__Expanding
	for i := 0; i < 30; i++ {
		filler := widgets.NewQSpacerItem(w, h, hPolicy, vPolicy)
		parent.layout.AddItem(filler)
	}
}

func (box *HBox) SetFixedHeight() {
	sizePolicy := box.qWidget.SizePolicy()
	sizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Fixed)
	box.qWidget.SetSizePolicy(sizePolicy)
}

func (box *VBox) SetFixedWidth() {
	sizePolicy := box.qWidget.SizePolicy()
	sizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Fixed)
	box.qWidget.SetSizePolicy(sizePolicy)
}

func (form *Form) qWidget_ITF() widgets.QWidget_ITF {
	return &form.qWidget
}

func (form *Form) Window() *Window {
	return form.window
}

func (form *Form) setWindow(window *Window) {
	form.window = window
}

func (parent *Form) AddWidget(w Widget) {
	switch v := w.(type) {
	case *FieldWidget:
		parent.AddFieldWidget(v)
	default:
		l.Fatalf("cannot add %T to form", v)
	}
}

func (parent *Form) AddFieldWidget(w *FieldWidget) {
	if w.label != nil {
		parent.layout.AddRow(w.label, w.qWidget_ITF())
		return
	}

	parent.layout.AddWidget(w.qWidget)
	w.window = parent.window
}

func (form *Form) SetContentsMargins(left int, right int, top int, bottom int) {
	form.layout.SetContentsMargins(left, right, top, bottom)
}

func (form *Form) QWidget() widgets.QWidget {
	return form.qWidget
}

func (parent *Form) AddRow(label string, w *FieldWidget) {
	w.SetLabel(label)
	parent.AddWidget(w)
}

func (parent *Form) AddFieldTypeRows(r *codeplug.Record, fTypes ...codeplug.FieldType) {
	fields := make([]*codeplug.Field, len(fTypes))
	for i, fType := range fTypes {
		fields[i] = r.Field(fType)
	}
	labelFunc := func(f *codeplug.Field) string {
		return f.TypeName()
	}
	recordNames := r.Names()
	if recordNames != nil {
		labelFunc = func(f *codeplug.Field) string {
			return recordNames[r.Index()]
		}
	}
	parent.AddFieldRows(labelFunc, fields...)
}

func (parent *Form) AddFieldRows(labelFunc func(*codeplug.Field) string, fields ...*codeplug.Field) {
	for _, f := range fields {
		parent.AddFieldRow(labelFunc, f)
	}
}

func (window *Window) AddFieldWidget(w *FieldWidget) {
	w.window = window

	widgets := window.widgets
	f := w.field
	fType := f.Type()
	if widgets[fType] == nil {
		window.widgets[fType] = make(map[*codeplug.Field]*FieldWidget)
	}
	window.widgets[fType][f] = w
}

func (window *Window) enableAllWidgets() {
	for _, widgetMap := range window.widgets {
		for _, w := range widgetMap {
			w.setEnabled()
		}
	}
}

func (window *Window) NewFieldWidget(label string, f *codeplug.Field) *FieldWidget {
	newFieldWidgetFunc := newFieldWidget[f.ValueType()]
	if newFieldWidgetFunc == nil {
		l.Fatalf("No %s entry in newFieldWidget slice", f.ValueType())
	}

	w := newFieldWidgetFunc(f)
	w.label = widgets.NewQLabel2(label, nil, 0)

	window.AddFieldWidget(w)

	if f.Enables() != nil {
		window.AddEnablerWidget(w)
	}

	enablerType := f.EnablerType()

	w.enable = func(enabler *FieldWidget) {
		ef := enabler.field
		if ef.Record().Type() != f.Record().Type() {
			l.Fatal("enabler record type", ef.Record().Type(), " enabled record type", f.Record().Type())
		}
		if ef.Record().Index() != f.Record().Index() {
			l.Fatal("enabler record index", ef.Record().Index(), " enabled record index", f.Record().Index())
		}
		if ef.Index() != f.Index() {
			l.Fatal("enabler field index", ef.Index(), " enabled field index", f.Index())
		}

		eft := ef.Type()
		if eft != enablerType {
			l.Fatal("enable(): unexpected field type:", eft)
		}

		w.setEnabled()
	}

	return w
}

func (parent *Form) AddFieldRow(labelFunc func(*codeplug.Field) string, f *codeplug.Field) {
	if f == nil {
		// This is not an error because some forms are used for
		// multiple models and some models do not include all fields.
		// We just ignore non-existent fields.
		return
	}
	w := parent.window.NewFieldWidget(labelFunc(f), f)
	parent.layout.AddRow(w.label, w.qWidget)
}

func (parent *Form) AddReadOnlyFieldTypeRows(r *codeplug.Record, fTypes ...codeplug.FieldType) {
	fields := make([]*codeplug.Field, len(fTypes))
	for i, fType := range fTypes {
		fields[i] = r.Field(fType)
	}
	labelFunc := func(f *codeplug.Field) string {
		return f.TypeName()
	}
	parent.AddReadOnlyFieldRows(labelFunc, fields...)
}

func (parent *Form) AddReadOnlyFieldRows(labelFunc func(*codeplug.Field) string, fields ...*codeplug.Field) {
	for _, f := range fields {
		parent.AddReadOnlyFieldRow(labelFunc, f)
	}
}

func (parent *Form) AddReadOnlyFieldRow(labelFunc func(*codeplug.Field) string, f *codeplug.Field) {
	if f == nil {
		// This is not an error because some forms are used for
		// multiple models and some models do not include all fields.
		// We just ignore non-existent fields.
		return
	}

	w := newFieldLineEdit(f)
	w.label = widgets.NewQLabel2(labelFunc(f), nil, 0)
	w.SetReadOnly(true)
	parent.layout.AddRow(w.label, w.qWidget)
}

func setFieldString(f *codeplug.Field, s string) error {
	err := f.TestSetString(s)
	if err != nil {
		return err
	}

	if !setMultipleRecords(f, s) {
		f.SetString(s)
	}
	return nil
}

type Table struct {
	qWidget widgets.QTableWidget
	window  *Window
}

func NewTable() *Table {
	t := new(Table)
	t.qWidget = *widgets.NewQTableWidget(nil)

	return t
}

func (parent *HBox) AddTable() *Table {
	t := NewTable()
	parent.AddWidget(t)
	t.window = parent.window
	return t
}

func (parent *VBox) AddTable() *Table {
	t := NewTable()
	parent.AddWidget(t)
	t.window = parent.window
	return t
}

func (t *Table) qWidget_ITF() widgets.QWidget_ITF {
	return &t.qWidget
}

func (t *Table) AddWidget(w Widget) {
	l.Fatal("cannot add widget directly to Table")
}

func (t *Table) SetContentsMargins(left int, right int, top int, bottom int) {
	l.Fatal("cannot set contents margin of Table")
}

func (t *Table) setWindow(window *Window) {
	t.window = window
}

func (t *Table) Window() *Window {
	return t.window
}

func (t *Table) AddRow(cells ...Widget) {
	qw := t.qWidget

	row := qw.RowCount()
	qw.SetRowCount(row + 1)

	if qw.ColumnCount() < len(cells) {
		qw.SetColumnCount(len(cells))
	}

	for i, cell := range cells {
		qw.SetCellWidget(row, i, cell.qWidget_ITF())
	}
}

func (t *Table) RowCount() int {
	return t.qWidget.RowCount()
}

func (t *Table) ColumnCount() int {
	return t.qWidget.ColumnCount()
}

func (t *Table) AddTopLabels(labels []string) {
	t.qWidget.SetHorizontalHeaderLabels(labels)
}

func (t *Table) AddLeftLabels(labels []string) {
	t.qWidget.SetVerticalHeaderLabels(labels)
}

func (t *Table) SetFixedSize() {
	sizePolicy := t.qWidget.SizePolicy()
	sizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Fixed)
	sizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Fixed)
	t.qWidget.SetSizePolicy(sizePolicy)
}

func (t *Table) ResizeToContents() {
	qw := t.qWidget
	hh := qw.HorizontalHeader()
	vh := qw.VerticalHeader()

	qw.ResizeColumnsToContents()
	qw.ResizeRowsToContents()
	hh.SetSectionResizeMode(widgets.QHeaderView__ResizeToContents)
	vh.SetSectionResizeMode(widgets.QHeaderView__ResizeToContents)

	width := 2 + vh.Width() + hh.Length()
	// vh.Width() apparently isn't updated to account for the label
	// We'll approximate it by the following.
	label := qw.VerticalHeaderItem(qw.RowCount()-1).Text() + "M"
	metrics := gui.NewQFontMetrics(widgets.QApplication_Font())
	vhWidth := metrics.HorizontalAdvance(label, -1)

	width = 2 + vhWidth + hh.Length()
	height := 2 + hh.Height() + vh.Length()

	qw.SetMinimumHeight(height)
	qw.SetMaximumHeight(height)
	qw.SetMinimumWidth(width)
	qw.SetMaximumWidth(width)
}

func setMultipleRecords(f *codeplug.Field, str string) bool {
	if f.MaxFields() > 1 {
		return false
	}
	r := f.Record()
	if r.MaxRecords() <= 1 {
		return false
	}

	recs := selectedRecords(r)
	if recs == nil {
		return false
	}

	if len(recs) <= 1 {
		return false
	}
	found := false
	for _, rec := range recs {
		if rec == r {
			found = true
			break
		}
	}
	if !found {
		return false
	}

	rw := recordWindow(r)
	if rw.settingMultiple {
		return false
	}
	rw.settingMultiple = true

	DelayedCall(func() {
		cp := r.Codeplug()
		howmany := "all"
		if len(recs) < len(cp.Records(r.Type())) {
			howmany = fmt.Sprintf("%d selected", len(recs))
		}

		typeName := f.Record().TypeName()

		msg := fmt.Sprintf(`Set "%s" to "%s" in %s %s?`, f.TypeName(), str, howmany, typeName)
		ans := YesNoPopup(fmt.Sprintf("Change multiple %s", typeName), msg)
		if ans != PopupYes {
			rw.settingMultiple = false
			f.SetString(str)
			return
		}

		rCount := len(recs)
		pd := NewProgressDialog("Updating " + typeName)
		pd.SetRange(0, rCount)

		progFunc := func(i int) {
			pd.SetValue(i)
		}

		cp.SetRecordsField(recs, f.Type(), str, progFunc)

		pd.Close()

		rw.settingMultiple = false
	})

	return true
}

type Widget interface {
	qWidget_ITF() widgets.QWidget_ITF
	AddWidget(Widget)
	SetContentsMargins(left int, right int, top int, bottom int)
	setWindow(*Window)
	Window() *Window
}

func (w *FieldWidget) setEnabled() {
	f := w.field
	enabled := f.IsEnabled()
	invalid := f.IsInvalidValue()
	if enabled && invalid {
		f.SetDefault()
	}

	w.SetEnabled(enabled)
	w.label.SetEnabled(enabled)
	if invalid {
		f.SetInvalidValue()
	}
	w.update()
	w.enableWidgets()
	if enabled && w.stacker != nil {
		w.stacker.enableOverlappingWidget(w)
	}
}

type FieldWidget struct {
	qWidget widgets.QWidget_ITF
	label   *widgets.QLabel
	field   *codeplug.Field
	enable  func(enabler *FieldWidget)
	stacker *StackedWidget
	window  *Window
}

func (w *FieldWidget) qWidget_ITF() widgets.QWidget_ITF {
	return w.qWidget
}

func (w *FieldWidget) Window() *Window {
	return w.window
}

func (w *FieldWidget) setWindow(window *Window) {
	w.window = window
}

func (fw *FieldWidget) AddWidget(w Widget) {
	l.Fatal("cannot add widget to FieldWidget")
}

func (fw *FieldWidget) SetContentsMargins(left int, right int, top int, bottom int) {
	l.Fatal("cannot set contents margin of FieldWidget")
}

func (w *FieldWidget) SetLabel(text string) {
	if w.label == nil {
		w.label = widgets.NewQLabel2(text, nil, 0)
		return
	}
	w.label.SetText(text)
}

func (w *FieldWidget) update() {
	f := w.field
	qw := w.qWidget

	switch qw.(type) {
	case *widgets.QCheckBox:
		setQCheckBox(qw.(*widgets.QCheckBox), f)

	case *widgets.QSpinBox:
		setQSpinBox(qw.(*widgets.QSpinBox), f)

	case *widgets.QLineEdit:
		str := f.String()
		str = codeplug.RemoveSuffix(f, str)
		qw.(*widgets.QLineEdit).SetText(str)

	case *widgets.QComboBox:
		qwc := qw.(*widgets.QComboBox)
		if f.IsInvalidValue() {
			ivs := codeplug.InvalidValueString
			qwc.InsertItems(0, []string{ivs})
			qwc.SetCurrentText(ivs)
			qwc.ConnectHighlighted(func(index int) {
				qwc.DisconnectHighlighted()
				if qwc.ItemData(0, 0).ToString() == ivs {
					if len(f.Strings()) > 0 {
						qwc.RemoveItem(0)
					}
				}
			})
		}
		qwc.SetCurrentText(codeplug.RemoveSuffix(f, f.String()))

	default:
		l.Fatal("update(): unexpected widget type")
	}
}

func (w *FieldWidget) SetEnabled(b bool) {
	qw := w.qWidget

	switch qw.(type) {
	case *widgets.QComboBox:
		qwc := qw.(*widgets.QComboBox)
		qwc.SetEnabled(b)
		f := w.field
		if f != nil {
			strs := codeplug.RemoveSuffixes(f.Strings())
			str := f.String()
			if f.IsInvalidValue() {
				ivs := codeplug.InvalidValueString
				strs = append([]string{ivs}, strs...)
				str = ivs
				qwc.ConnectHighlighted(func(index int) {
					qwc.DisconnectHighlighted()
					if qwc.ItemData(0, 0).ToString() == ivs {
						if len(f.Strings()) > 0 {
							qwc.RemoveItem(0)
						}
					}
				})
			}
			UpdateComboboxWidget(w, codeplug.RemoveSuffix(f, str), strs)
		}

	case *widgets.QPushButton:
		qw.(*widgets.QPushButton).SetEnabled(b)

	case *widgets.QCheckBox:
		qw.(*widgets.QCheckBox).SetEnabled(b)

	case *widgets.QSpinBox:
		qw.(*widgets.QSpinBox).SetEnabled(b)

	case *widgets.QLineEdit:
		qw.(*widgets.QLineEdit).SetEnabled(b)

	default:
		l.Fatal("SetEnabled(): unexpected widget type")
	}
}

func (w *FieldWidget) SetReadOnly(b bool) {
	qw := w.qWidget

	switch qw.(type) {
	case *widgets.QSpinBox:
		qw.(*widgets.QSpinBox).SetReadOnly(b)

	case *widgets.QLineEdit:
		qw.(*widgets.QLineEdit).SetReadOnly(b)

	default:
		l.Fatal("SetReadOnly(): unexpected widget type")
	}
}

func (w *FieldWidget) SetVisible(b bool) {
	qw := w.qWidget

	switch qw.(type) {
	case *widgets.QComboBox:
		qw.(*widgets.QComboBox).SetVisible(b)

	case *widgets.QPushButton:
		qw.(*widgets.QPushButton).SetVisible(b)

	case *widgets.QCheckBox:
		qw.(*widgets.QCheckBox).SetVisible(b)

	case *widgets.QSpinBox:
		qw.(*widgets.QSpinBox).SetVisible(b)

	case *widgets.QLineEdit:
		qw.(*widgets.QLineEdit).SetVisible(b)

	default:
		l.Fatal("SetVisible(): unexpected widget type")
	}

	if w.label != nil {
		w.label.SetVisible(b)
	}
}

func (w *FieldWidget) Width() int {
	qw := w.qWidget

	var width int

	switch qw.(type) {
	case *widgets.QSpinBox:
		width = qw.(*widgets.QSpinBox).Width()

	case *widgets.QLineEdit:
		width = qw.(*widgets.QLineEdit).Width()

	default:
		l.Fatal("SetMinimumWidth(): unexpected widget type")
	}

	return width
}

func (w *FieldWidget) SetMinimumWidth(width int) {
	qw := w.qWidget

	switch qw.(type) {
	case *widgets.QSpinBox:
		qw.(*widgets.QSpinBox).SetMinimumWidth(width)

	case *widgets.QLineEdit:
		qw.(*widgets.QLineEdit).SetMinimumWidth(width)

	default:
		l.Fatal("SetMinimumWidth(): unexpected widget type")
	}
}

func setQCheckBox(cb *widgets.QCheckBox, f *codeplug.Field) {
	checkState := core.Qt__Unchecked
	if f.String() == "On" {
		checkState = core.Qt__Checked
	}
	cb.SetCheckState(checkState)
}

func newFieldCheckbox(f *codeplug.Field) *FieldWidget {
	qw := widgets.NewQCheckBox(nil)
	w := new(FieldWidget)
	w.qWidget = qw
	w.field = f

	setQCheckBox(qw, f)

	qw.ConnectClicked(func(checked bool) {
		str := "Off"
		if checked {
			str = "On"
		}
		err := setFieldString(f, str)
		if err != nil {
			l.Fatal(err.Error())
		}
	})

	return w
}

func newFieldLineEdit(f *codeplug.Field) *FieldWidget {
	s := f.String()
	s = codeplug.RemoveSuffix(f, s)
	qw := widgets.NewQLineEdit2(s, nil)
	widget := new(FieldWidget)
	widget.qWidget = qw
	widget.field = f
	metrics := gui.NewQFontMetrics(widgets.QApplication_Font())
	widget.SetMinimumWidth(metrics.HorizontalAdvance(s, -1) + 12)

	var finished func()
	finished = func() {
		str := strings.TrimSpace(qw.Text())
		str = codeplug.AddSuffix(f, str)
		err := setFieldString(f, str)
		if err != nil {
			msg := f.TypeName() + " " + err.Error()
			qw.DisconnectEditingFinished()
			ErrorPopup("Value error", msg)
			qw.ConnectEditingFinished(finished)
		}
	}

	qw.ConnectEditingFinished(finished)

	return widget
}

func newFieldCombobox(f *codeplug.Field) *FieldWidget {
	qw := widgets.NewQComboBox(nil)
	widget := new(FieldWidget)
	widget.qWidget = qw
	widget.field = f

	strings := f.Strings()
	span := f.Span()
	if span != nil {
		strings = f.SpanStrings()
	}

	if len(strings) == 0 {
		l.Fatal("Combobox has no Strings()")
	}

	qw.InsertItems(0, codeplug.RemoveSuffixes(strings))
	qw.SetCurrentText(codeplug.RemoveSuffix(f, f.String()))

	if f.IsInvalidValue() {
		ivs := codeplug.InvalidValueString
		qw.InsertItems(0, []string{ivs})
		qw.SetCurrentText(ivs)
		qw.ConnectHighlighted(func(index int) {
			qw.DisconnectHighlighted()
			if qw.ItemData(0, 0).ToString() == ivs {
				if len(f.Strings()) > 0 {
					qw.RemoveItem(0)
				}
			}
		})
	}

	qw.ConnectActivated(func(index int) {
		str := f.Strings()[index]
		err := setFieldString(f, str)
		if err != nil {
			msg := f.TypeName() + " " + err.Error()
			ErrorPopup("Value error", msg)
		}
	})

	return widget
}

func setQSpinBox(sb *widgets.QSpinBox, f *codeplug.Field) {
	str := f.String()
	value := 0
	span := f.Span()
	if str != span.MinString() {
		i, err := strconv.ParseUint(str, 10, 32)
		if err == nil {
			value = int(i)
		}
	}
	sb.SetValue(value)
}

func NewLineEditWidget(str string) *FieldWidget {
	qw := widgets.NewQLineEdit(nil)
	widget := new(FieldWidget)
	widget.qWidget = qw
	qw.SetText(str)
	qw.SetReadOnly(true)

	return widget
}

func NewSpinboxWidget(value, min, max int, changedFunc func(int)) *FieldWidget {
	qw := widgets.NewQSpinBox(nil)
	widget := new(FieldWidget)
	widget.qWidget = qw
	qw.SetRange(min, max)
	qw.SetValue(value)

	qw.ConnectEditingFinished(func() {
		changedFunc(qw.Value())
	})

	return widget
}

func NewCheckboxWidget(checked bool, clickedFunc func(bool)) *FieldWidget {
	qw := widgets.NewQCheckBox(nil)
	widget := new(FieldWidget)
	widget.qWidget = qw
	qw.SetChecked(checked)
	qw.ConnectClicked(clickedFunc)

	return widget
}

func (widget *FieldWidget) SetChecked(checked bool) {
	checkbox, ok := widget.qWidget.(*widgets.QCheckBox)
	if ok {
		checkbox.SetChecked(checked)
	}
}

type StackedWidget struct {
	qStackedWidget widgets.QStackedWidget
	widgets        []Widget
	window         *Window
}

func NewStackedWidget(window *Window) *StackedWidget {
	sw := new(StackedWidget)
	sw.qStackedWidget = *widgets.NewQStackedWidget(nil)
	sw.window = window
	return sw
}

func (sw *StackedWidget) qWidget_ITF() widgets.QWidget_ITF {
	return widgets.QWidget_ITF(&sw.qStackedWidget)
}

func (sw *StackedWidget) Window() *Window {
	return sw.window
}

func (sw *StackedWidget) setWindow(window *Window) {
	sw.window = window
}

func (sw *StackedWidget) SetContentsMargins(left int, right int, top int, bottom int) {
	l.Fatal("cannot set contents margin of StackedWidget")
}

func (sw *StackedWidget) AddWidget(w Widget) {
	sw.qStackedWidget.AddWidget(w.qWidget_ITF())
	sw.widgets = append(sw.widgets, w)

	switch v := w.(type) {
	case *FieldWidget:
		v.stacker = sw
	}
}

func (sw *StackedWidget) enableOverlappingWidget(w *FieldWidget) {
	sw.qStackedWidget.SetCurrentWidget(w.qWidget)

	for _, widget := range sw.widgets {
		widget.(*FieldWidget).field.SetStore(widget == w)
	}
}

func (sw *StackedWidget) setCurrentWidget(w *FieldWidget) {
	sw.qStackedWidget.SetCurrentWidget(w.qWidget)
}

type TabWidget struct {
	qTabWidget widgets.QTabWidget
	widgets    []*FieldWidget
	window     *Window
}

func NewTabWidget(window *Window) *TabWidget {
	tw := new(TabWidget)
	tw.qTabWidget = *widgets.NewQTabWidget(nil)
	tw.window = window
	return tw
}

func (tw *TabWidget) AddTab(w *FieldWidget, label string) {
	tw.qTabWidget.AddTab(w.qWidget, label)
	tw.widgets = append(tw.widgets, w)
}

func (tw *TabWidget) ConnectChange(f func(w *FieldWidget)) {
	tw.qTabWidget.ConnectCurrentChanged(func(index int) {
		f(tw.widgets[index])
	})
}

func NewComboboxWidget(opt string, opts []string, changed func(int)) *FieldWidget {
	qw := widgets.NewQComboBox(nil)
	widget := new(FieldWidget)
	widget.qWidget = qw
	qw.InsertItems(0, opts)
	qw.SetCurrentText(opt)

	qw.ConnectCurrentIndexChanged(changed)

	return widget
}

func UpdateComboboxWidget(widget *FieldWidget, opt string, opts []string) {
	qcb := widget.qWidget.(*widgets.QComboBox)
	qcb.Clear()
	qcb.InsertItems(0, opts)
	qcb.SetCurrentText(opt)
}

func NewButtonWidget(text string, clicked func()) *FieldWidget {
	w := new(FieldWidget)
	b := widgets.NewQPushButton2(text, nil)
	b.SetSizePolicy2(widgets.QSizePolicy__Fixed,
		widgets.QSizePolicy__Fixed)
	b.ConnectClicked(func(checked bool) {
		clicked()
	})
	w.qWidget = b

	return w
}

type TextEdit struct {
	qWidget widgets.QTextEdit
}

func NewTextEdit() *TextEdit {
	t := new(TextEdit)
	t.qWidget = *widgets.NewQTextEdit(nil)

	return t
}

func (parent *HBox) AddTextEdit() *TextEdit {
	t := NewTextEdit()
	parent.layout.AddWidget(&t.qWidget, 0, 0)

	return t
}

func (parent *VBox) AddTextEdit() *TextEdit {
	t := NewTextEdit()
	parent.layout.AddWidget(&t.qWidget, 0, 0)

	return t
}

func (t *TextEdit) SetPlainText(str string) {
	t.qWidget.SetPlainText(str)
}

func (t *TextEdit) SetNoLineWrap() {
	t.qWidget.SetLineWrapMode(widgets.QTextEdit__NoWrap)
}

func (t *TextEdit) SetReadOnly(ro bool) {
	t.qWidget.SetReadOnly(ro)
}

func newFieldSpinbox(f *codeplug.Field) *FieldWidget {
	qw := widgets.NewQSpinBox(nil)
	widget := new(FieldWidget)
	widget.qWidget = qw
	widget.field = f

	span := f.Span()
	qw.SetRange(span.Minimum(), span.Maximum())
	qw.SetSingleStep(span.Step())
	qw.SetWrapping(true)
	qw.SetSpecialValueText(span.MinString())

	setQSpinBox(qw, f)

	if f.IsInvalidValue() {
		qw.SetSpecialValueText(f.String())
		qw.SetValue(span.Minimum())
		qw.ConnectFocusInEvent(func(event *gui.QFocusEvent) {
			setFieldString(f, strconv.Itoa(span.Minimum()))
			qw.SetSpecialValueText(span.MinString())
			qw.DisconnectFocusInEvent()
		})
	}

	qw.ConnectEditingFinished(func() {
		val := qw.Value()
		str := span.MinString()
		if str == "" || val != span.Minimum() {
			str = qw.TextFromValue(val)
		}
		err := setFieldString(f, str)
		if err != nil {
			msg := f.TypeName() + " " + err.Error()
			ErrorPopup("Value error", msg)
		}
	})

	return widget
}

func newFieldTextEdit(f *codeplug.Field) *FieldWidget {
	l.Fatal("newTextEdit: not implemented")
	return nil
}

var newFieldWidget = map[codeplug.ValueType]func(*codeplug.Field) *FieldWidget{
	codeplug.VtAscii:              newFieldLineEdit,
	codeplug.VtBandwidth:          newFieldCombobox,
	codeplug.VtBiFrequency:        newFieldLineEdit,
	codeplug.VtCallID:             newFieldLineEdit,
	codeplug.VtCallType:           newFieldCombobox,
	codeplug.VtContactListIndex:   newFieldCombobox,
	codeplug.VtContactName:        newFieldLineEdit,
	codeplug.VtCtcssDcs:           newFieldCombobox,
	codeplug.VtDerefListIndex:     newFieldCombobox,
	codeplug.VtFrequency:          newFieldLineEdit,
	codeplug.VtFrequencyOffset:    newFieldLineEdit,
	codeplug.VtGpsListIndex:       newFieldCombobox,
	codeplug.VtGpsReportInterval:  newFieldSpinbox,
	codeplug.VtHexadecimal32:      newFieldLineEdit,
	codeplug.VtHexadecimal4:       newFieldLineEdit,
	codeplug.VtIndexedStrings:     newFieldCombobox,
	codeplug.VtIntroLine:          newFieldLineEdit,
	codeplug.VtIStrings:           newFieldCombobox,
	codeplug.VtListIndex:          newFieldCombobox,
	codeplug.VtMemberListIndex:    newFieldCombobox,
	codeplug.VtName:               newFieldLineEdit,
	codeplug.VtNkContactListIndex: newFieldCombobox,
	codeplug.VtOffOn:              newFieldCheckbox,
	codeplug.VtOnOff:              newFieldCheckbox,
	codeplug.VtPcPassword:         newFieldLineEdit,
	codeplug.VtPrivacyNumber:      newFieldCombobox,
	codeplug.VtRadioButton:        newFieldCombobox,
	codeplug.VtRadioName:          newFieldLineEdit,
	codeplug.VtRadioPassword:      newFieldLineEdit,
	codeplug.VtRadioProgPassword:  newFieldLineEdit,
	codeplug.VtSpanList:           newFieldCombobox,
	codeplug.VtSpan:               newFieldSpinbox,
	codeplug.VtTextMessage:        newFieldLineEdit,
	codeplug.VtTimeStamp:          newFieldLineEdit,
}

type MenuBar struct {
	qMenuBar *widgets.QMenuBar
}

func (mw *MainWindow) MenuBar() *MenuBar {
	mb := new(MenuBar)
	mb.qMenuBar = mw.qMainWindow.MenuBar()

	return mb
}

type Menu struct {
	qMenu *widgets.QMenu
}

func (mb *MenuBar) AddMenu(name string) *Menu {
	menu := new(Menu)
	menu.qMenu = mb.qMenuBar.AddMenu2(name)

	return menu
}

func (mb *MenuBar) Clear() {
	mb.qMenuBar.Clear()
}

type Action struct {
	qAction *widgets.QAction
}

func (action *Action) SetVisible(b bool) {
	action.qAction.SetVisible(b)
}

func (menu *Menu) AddAction(name string, fn func()) *Action {
	action := new(Action)
	action.qAction = menu.qMenu.AddAction(name)

	action.qAction.ConnectTriggered(func(checked bool) {
		fn()
	})

	return action
}

func (menu *Menu) AddMenu(name string) *Menu {
	subMenu := new(Menu)
	subMenu.qMenu = menu.qMenu.AddMenu2(name)

	return subMenu
}

func (menu *Menu) ConnectAboutToShow(fn func()) {
	menu.qMenu.ConnectAboutToShow(fn)
}

func (menu *Menu) AddSeparator() {
	menu.qMenu.AddSeparator()
}

func (menu *Menu) Clear() {
	menu.qMenu.Clear()
}

func (menu *Menu) SetEnabled(enable bool) {
	menu.qMenu.SetEnabled(enable)
}

func (a *Action) SetText(s string) {
	a.qAction.SetText(s)
}

func (a *Action) SetEnabled(enable bool) {
	a.qAction.SetEnabled(enable)
}

type RadioButton struct {
	qWidget *widgets.QRadioButton
}

func (parent *HBox) AddRadioButton(option string) *RadioButton {
	b := new(RadioButton)
	b.qWidget = widgets.NewQRadioButton2(option, nil)
	parent.layout.AddWidget(b.qWidget, 0, 0)

	return b
}

func (b *RadioButton) ConnectClicked(fn func(checked bool)) {
	b.qWidget.ConnectClicked(func(checked bool) {
		fn(checked)
	})
}

func (b *RadioButton) SetChecked(bo bool) {
	b.qWidget.SetChecked(bo)
}

func (b *RadioButton) IsChecked() bool {
	return b.qWidget.IsChecked()
}

func (b *RadioButton) Text() string {
	return b.qWidget.Text()
}

type Button struct {
	qWidget widgets.QPushButton
}

func NewButton(text string) *Button {
	b := new(Button)
	b.qWidget = *widgets.NewQPushButton2(text, nil)
	b.qWidget.SetSizePolicy2(widgets.QSizePolicy__Fixed,
		widgets.QSizePolicy__Fixed)

	return b
}

func (b *Button) ConnectClicked(fn func()) {
	b.qWidget.ConnectClicked(func(checked bool) {
		fn()
	})
}

func (b *Button) SetText(str string) {
	b.qWidget.SetText(str)
}

func (b *Button) SetEnabled(enable bool) {
	b.qWidget.SetEnabled(enable)
}

func (w *Window) SetRecordFunc(fn func()) {
	w.recordFunc = fn
}

func (w *Window) RecordFunc() func() {
	return w.recordFunc
}

func InfoPopup(title string, msg string) {
	button := widgets.QMessageBox__Ok
	defaultButton := widgets.QMessageBox__Ok
	widgets.QMessageBox_Information(nil, title, msg, button, defaultButton)
}

func WarningPopup(title string, msg string) PopupValue {
	maxLines := 12
	lines := strings.SplitN(msg, "\n", maxLines+1)
	if len(lines) > 1 {
		msg = strings.Join(lines[:maxLines], "\n") + "\n"
		if len(lines) > maxLines {
			msg += "...\n"
		}
	}
	buttons := widgets.QMessageBox__Cancel | widgets.QMessageBox__Ignore
	defButton := widgets.QMessageBox__Cancel
	rv := widgets.QMessageBox_Warning(nil, title, msg, buttons, defButton)
	switch rv {
	case widgets.QMessageBox__Cancel:
		return PopupCancel

	case widgets.QMessageBox__Ignore:
		return PopupIgnore

	default:
		return PopupCancel
	}
}

func ErrorPopup(title string, msg string) {
	if msg == "" {
		return
	}
	button := widgets.QMessageBox__Ok
	defaultButton := widgets.QMessageBox__Ok
	widgets.QMessageBox_Critical(nil, title, msg, button, defaultButton)
}

type PopupValue int

const (
	PopupCancel PopupValue = iota
	PopupDiscard
	PopupIgnore
	PopupNo
	PopupSave
	PopupYes
)

func SavePopup(title string, msg string) PopupValue {
	buttons := widgets.QMessageBox__Save |
		widgets.QMessageBox__Discard | widgets.QMessageBox__Cancel

	rv := widgets.QMessageBox_Warning(nil, title, msg, buttons, 0)
	switch rv {
	case widgets.QMessageBox__Save:
		break

	case widgets.QMessageBox__Discard:
		return PopupDiscard

	case widgets.QMessageBox__Cancel:
		return PopupCancel
	}

	return PopupSave
}

func YesNoPopup(title string, msg string) PopupValue {
	buttons := widgets.QMessageBox__Yes | widgets.QMessageBox__No

	rv := widgets.QMessageBox_Warning(nil, title, msg, buttons, 0)
	switch rv {
	case widgets.QMessageBox__Yes:
		return PopupYes
		break

	default:
		break
	}

	return PopupNo
}

type ProgressDialog struct {
	qWidget *widgets.QProgressDialog
}

func NewProgressDialog(str string) *ProgressDialog {
	pd := new(ProgressDialog)
	qpd := widgets.NewQProgressDialog(nil, core.Qt__Dialog)
	pd.qWidget = qpd

	qpd.SetWindowModality(core.Qt__ApplicationModal)
	qpd.SetMinimumDuration(250) // 250 mS
	qpd.SetLabelText(str)

	return pd
}

func (pd *ProgressDialog) Close() {
	pd.qWidget.Close()
}

func (pd *ProgressDialog) SetLabelText(str string) {
	pd.qWidget.SetLabelText(str)
}

func (pd *ProgressDialog) SetRange(min int, max int) {
	pd.qWidget.SetRange(min, max)
}

func (pd *ProgressDialog) SetValue(value int) {
	pd.qWidget.SetValue(value)
}

func (pd *ProgressDialog) WasCanceled() bool {
	return pd.qWidget.WasCanceled()
}

type Dialog struct {
	*VBox
	qDialog *widgets.QDialog
	layout  *widgets.QVBoxLayout
}

func NewDialog(title string) *Dialog {
	dialog := new(Dialog)
	dialog.VBox = NewVbox()
	dialog.qDialog = widgets.NewQDialog(nil, core.Qt__Dialog)
	dialog.layout = widgets.NewQVBoxLayout2(dialog.qDialog)
	dialog.layout.AddWidget(&dialog.VBox.qWidget, 0, 0)

	dialog.qDialog.SetWindowTitle(title)
	dialog.qDialog.SetWindowModality(core.Qt__ApplicationModal)

	return dialog
}

func (d *Dialog) Exec() bool {
	return d.qDialog.Exec() == int(widgets.QDialog__Accepted)
}

func (d *Dialog) Accept() {
	d.qDialog.Accept()
}

func (d *Dialog) Reject() {
	d.qDialog.Reject()
}

func (d *Dialog) Update() {
	d.qDialog.Update()
}

func OpenTextFilename(title string, dir string) string {
	selF := "(*.txt)"
	filter := "Text files " + selF + ";;All files (*)"
	return widgets.QFileDialog_GetOpenFileName(nil, title, dir, filter, selF, 0)
}

func OpenJSONFilename(title string, dir string) string {
	selF := "(*.json)"
	filter := "JSON files " + selF + ";;All files (*)"
	return widgets.QFileDialog_GetOpenFileName(nil, title, dir, filter, selF, 0)
}

func OpenXLSXFilename(title string, dir string) string {
	selF := "(*.xlsx)"
	filter := "Spreadsheet files " + selF + ";;All files (*)"
	return widgets.QFileDialog_GetOpenFileName(nil, title, dir, filter, selF, 0)
}

func OpenCPFilenames(title string, dir string, exts []string) []string {
	for i, ext := range exts {
		exts[i] = "*." + ext
	}
	selF := "(" + strings.Join(exts, " ") + ")"
	filter := "Codeplug files " + selF + ";;All files (*)"
	return widgets.QFileDialog_GetOpenFileNames(nil, title, dir, filter, selF, 0)
}

func SaveFilename(title string, dir string, extension string) string {
	selF := "(*." + extension + ")"
	filter := "Codeplug files " + selF + ";;All files (*)"
	return widgets.QFileDialog_GetSaveFileName(nil, title, dir, filter, selF, 0)
}

func mainWindow(cp *codeplug.Codeplug) *MainWindow {
	for _, mw := range mainWindows {
		if mw.codeplug == cp {
			return mw
		}
	}
	return nil
}

func recordWindow(r *codeplug.Record) *Window {
	cp := r.Codeplug()
	mw := mainWindow(cp)
	rType := r.Type()
	for _, w := range mw.recordWindows {
		if w.recordType == rType {
			return w
		}
	}
	for _, w := range mw.altRecordWindows {
		if w.altRecordTypes[rType] {
			return w
		}
	}
	return nil
}

func recordList(r *codeplug.Record) *RecordList {
	w := recordWindow(r)
	if w == nil {
		return nil
	}
	return w.recordList
}

func selectedRecords(r *codeplug.Record) []*codeplug.Record {
	rl := recordList(r)
	if rl == nil {
		return nil
	}
	return rl.SelectedRecords()
}
