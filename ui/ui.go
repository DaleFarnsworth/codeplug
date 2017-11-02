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
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dalefarnsworth/codeplug/codeplug"
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

func (app *App) Quit() {
	app.qApp.Quit()
}

type AppSettings struct {
	qSettings *core.QSettings
}

func (as *AppSettings) Sync() {
	as.qSettings.Sync()
}

func (as *AppSettings) SetBool(str string, b bool) {
	as.qSettings.SetValue(str, core.NewQVariant11(b))
}

func (as *AppSettings) Bool(str string, deflt bool) bool {
	return as.qSettings.Value(str, core.NewQVariant11(deflt)).ToBool()
}

func (as *AppSettings) SetInt(str string, i int) {
	as.qSettings.SetValue(str, core.NewQVariant7(i))
}

func (as *AppSettings) Int(str string, deflt int) int {
	var ok bool
	return as.qSettings.Value(str, core.NewQVariant7(deflt)).ToInt(ok)
}

func (as *AppSettings) SetString(str string, s string) {
	as.qSettings.SetValue(str, core.NewQVariant17(s))
}

func (as *AppSettings) String(str string, deflt string) string {
	return as.qSettings.Value(str, core.NewQVariant17(deflt)).ToString()
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

type MainWindow struct {
	qMainWindow   widgets.QMainWindow
	codeplug      *codeplug.Codeplug
	recordWindows map[codeplug.RecordType]*Window
	connectClose  func() bool
	connectChange func(*codeplug.Change)
}

func (mw *MainWindow) SetCodeplug(cp *codeplug.Codeplug) {
	mw.codeplug = cp
	mw.recordWindows = make(map[codeplug.RecordType]*Window)

	mw.codeplug.ConnectChange(func(change *codeplug.Change) {
		for _, mw := range mainWindows {
			if mw.codeplug != change.Codeplug() {
				continue
			}
			changes := append(change.Changes(), change)
			for _, change := range changes {
				w := mw.recordWindows[change.RecordType()]
				if w != nil {
					w.handleChange(change)
				}
			}

			mw.connectChange(change)
		}
	})
}

func NewMainWindow() *MainWindow {
	mw := new(MainWindow)
	mainWindows = append(mainWindows, mw)

	qmw := widgets.NewQMainWindow(nil, 0)

	qmw.Resize2(600, 100)
	//var aGeometry = widgets.QApplication_Desktop().AvailableGeometry2(qmw)
	//qmw.Move2((aGeometry.Width()-qmw.Width())/2,
	//(aGeometry.Height()-qmw.Height())/2)

	mw.qMainWindow = *qmw

	qmw.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		if mw.connectClose != nil {
			if !mw.connectClose() {
				event.Ignore()
				return
			}

		}
		for i, mainWindow := range mainWindows {
			if mainWindow == mw {
				mainWindows = append(mainWindows[:i], mainWindows[i+1:]...)
				break
			}
		}

		objs := qmw.Children()
		for _, obj := range objs {
			obj.DeleteLater()
		}

		event.Accept()
	})

	return mw
}

func (parent *MainWindow) AddVbox() *VBox {
	box := newVbox()
	parent.qMainWindow.SetCentralWidget(&box.qWidget)

	return box
}

func (parent *MainWindow) AddHbox() *HBox {
	box := newHbox()
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
	box := newHbox()

	parent.layout.AddWidget(&box.qWidget, 0, 0)
	box.window = parent

	return box
}

func (parent *Window) AddVbox() *VBox {
	box := newVbox()

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
	//w.layout.AddWidget(mb.qMenuBar, 0, 0)
	w.menuBar = mb
	return mb
}

func (w *Window) MenuBar() *MenuBar {
	return w.menuBar
}

func (w *Window) EnableWidgets() {
	widgets := w.widgets
	for senderType, subs := range w.subscriptions {
		for _, receiverType := range subs {
			widgets[receiverType].receive(widgets[senderType])
		}
	}
}
func (w *Window) Show() {
	w.qWidget.Show()
	w.qWidget.ActivateWindow()
	w.qWidget.Raise()
}

func (w *Window) Clear() {
	clear(w.qWidget)
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

func (box *HBox) Clear() {
	clear(box.qWidget)
}

func (box *VBox) Clear() {
	clear(box.qWidget)
}

func (box *HBox) SetEnabled(enable bool) {
	box.qWidget.SetEnabled(enable)
}

func (box *VBox) SetEnabled(enable bool) {
	box.qWidget.SetEnabled(enable)
}

func clear(widget widgets.QWidget) {
	objs := widget.Children()
	for _, obj := range objs {
		if obj.Pointer() == widget.Layout().Pointer() {
			continue
		}
		obj.DeleteLater()
	}
}

func (mw *MainWindow) Show() {
	mw.qMainWindow.Show()
	mw.qMainWindow.ActivateWindow()
	mw.qMainWindow.Raise()
}

func (mw *MainWindow) Close() {
	mw.qMainWindow.Close()
}

func (mw *MainWindow) RecordWindows() map[codeplug.RecordType]*Window {
	return mw.recordWindows
}

type Window struct {
	qWidget       widgets.QWidget
	layout        *widgets.QHBoxLayout
	menuBar       *MenuBar
	window        *Window
	mainWindow    *MainWindow
	recordType    codeplug.RecordType
	recordFunc    func()
	widgets       map[codeplug.FieldType]*Widget
	subscriptions map[codeplug.FieldType][]codeplug.FieldType
	recordModel   *core.QAbstractListModel
	recordList    *RecordList
	connectClose  func() bool
	handleChange  func(*codeplug.Change)
}

func (mw *MainWindow) NewWindow() *Window {
	w := new(Window)
	w.qWidget = *widgets.NewQWidget(&mw.qMainWindow, core.Qt__Window)
	w.layout = widgets.NewQHBoxLayout2(&w.qWidget)
	w.mainWindow = mw
	w.window = w

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

func (mw *MainWindow) NewRecordWindow(rType codeplug.RecordType) *Window {
	w := new(Window)
	mw.recordWindows[rType] = w
	w.qWidget = *widgets.NewQWidget(&mw.qMainWindow, core.Qt__Window)
	w.layout = widgets.NewQHBoxLayout2(&w.qWidget)
	w.mainWindow = mw
	w.window = w
	w.recordType = rType
	w.subscriptions = make(map[codeplug.FieldType][]codeplug.FieldType)
	w.widgets = make(map[codeplug.FieldType]*Widget)

	w.initRecordModel()

	w.qWidget.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		if w.connectClose != nil {
			if !w.connectClose() {
				event.Ignore()
				return
			}

		}
		delete(mw.recordWindows, rType)
		event.Accept()
	})

	w.handleChange = func(change *codeplug.Change) {
		rl := w.recordList

		updateRecordList := false
		updateRecord := false
		newCurrentRecord := -1

		changeType := change.Type()
		switch changeType {
		case codeplug.FieldChange:
			f := change.Field()
			for _, mw := range mainWindows {
				w := mw.recordWindows[w.recordType]
				if w != nil {
					widget := w.widgets[f.Type()]
					if widget != nil {
						widget.receive(widget)
					}
				}
			}

			if f.Type() == f.Record().NameFieldType() {
				updateRecordList = true
			}

		case codeplug.MoveRecordsChange, codeplug.InsertRecordsChange:
			newCurrentRecord = change.Record().Index()
			updateRecordList = true

		case codeplug.RemoveRecordsChange:
			updateRecordList = true

		case codeplug.MoveFieldsChange,
			codeplug.InsertFieldsChange,
			codeplug.RemoveFieldsChange,
			codeplug.ListIndexChange:
			updateRecord = true

		default:
			log.Fatal("Unknown change type", changeType)
		}

		if updateRecordList {
			rl.Update()
		}

		if newCurrentRecord >= 0 {
			rl.SetCurrent(newCurrentRecord)
		}

		if updateRecordList || updateRecord {
			if w.recordList.Current() == change.Record().Index() {
				w.recordFunc()
			}
		}
	}

	return w
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

func newVbox() *VBox {
	box := new(VBox)

	box.qWidget = *widgets.NewQWidget(nil, 0)
	box.layout = widgets.NewQVBoxLayout2(&box.qWidget)
	box.layout.SetContentsMargins(0, 0, 0, 0)

	return box
}

func (vBox *VBox) SetContentsMargins(left int, right int, top int, bottom int) {
	vBox.layout.SetContentsMargins(left, right, top, bottom)
}

func (parent *VBox) AddGroupbox(label string) *HBox {
	qgb := widgets.NewQGroupBox2(label, nil)
	layout := widgets.NewQHBoxLayout2(qgb)
	layout.SetContentsMargins(0, 0, 0, 0)

	box := newHbox()
	box.layout.SetContentsMargins(0, 0, 0, 0)
	layout.AddWidget(&box.qWidget, 0, 0)

	parent.layout.AddWidget(qgb, 0, 0)
	box.window = parent.window

	return box
}

func (vBox *VBox) Window() *Window {
	return vBox.window
}

func (parent *HBox) AddGroupbox(label string) *HBox {
	qgb := widgets.NewQGroupBox2(label, nil)
	layout := widgets.NewQHBoxLayout2(qgb)
	layout.SetContentsMargins(0, 0, 0, 0)

	box := newHbox()
	box.layout.SetContentsMargins(0, 0, 0, 0)
	layout.AddWidget(&box.qWidget, 0, 0)

	parent.layout.AddWidget(qgb, 0, 0)
	box.window = parent.window

	return box
}

func (hBox *HBox) Window() *Window {
	return hBox.window
}

func (parent *VBox) AddHbox() *HBox {
	box := newHbox()

	parent.layout.AddWidget(&box.qWidget, 0, 0)
	box.window = parent.window

	return box
}

func (parent *HBox) AddHbox() *HBox {
	box := newHbox()

	parent.layout.AddWidget(&box.qWidget, 0, 0)
	box.window = parent.window

	return box
}

func (parent *VBox) AddWidget(widget *Widget) {
	parent.layout.AddWidget(widget.qWidget, 0, 0)
}

func (parent *HBox) AddWidget(widget *Widget) {
	parent.layout.AddWidget(widget.qWidget, 0, 0)
}

func (parent *HBox) AddForm() *Form {
	form := new(Form)

	form.qWidget = *widgets.NewQWidget(nil, 0)
	form.layout = widgets.NewQFormLayout(&form.qWidget)
	form.layout.SetContentsMargins(0, 0, 0, 0)
	//form.layout.SetLabelAlignment(core.Qt__AlignRight)

	parent.layout.AddWidget(&form.qWidget, 0, 0)
	form.window = parent.window

	return form
}

func (parent *VBox) AddForm() *Form {
	form := new(Form)

	form.qWidget = *widgets.NewQWidget(nil, 0)
	form.layout = widgets.NewQFormLayout(&form.qWidget)
	//form.layout.SetLabelAlignment(core.Qt__AlignRight)

	parent.layout.AddWidget(&form.qWidget, 0, 0)
	form.window = parent.window

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

func newHbox() *HBox {
	box := new(HBox)

	box.qWidget = *widgets.NewQWidget(nil, 0)
	box.layout = widgets.NewQHBoxLayout2(&box.qWidget)
	box.layout.SetContentsMargins(0, 0, 0, 0)

	return box
}

func (hBox *HBox) SetContentsMargins(left int, right int, top int, bottom int) {
	hBox.layout.SetContentsMargins(left, right, top, bottom)
}

func (parent *HBox) AddVbox() *VBox {
	box := newVbox()

	parent.layout.AddWidget(&box.qWidget, 0, 0)
	box.window = parent.window

	return box
}

func (parent *HBox) AddButton(text string) *Button {
	b := NewButton(text)
	parent.layout.AddWidget(&b.qButton, 0, 0)

	return b
}

func (parent *VBox) AddButton(text string) *Button {
	b := NewButton(text)
	parent.layout.AddWidget(&b.qButton, 0, 0)

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
	filler := widgets.NewQSpacerItem(w, h, hPolicy, vPolicy)
	parent.layout.AddItem(filler)
}

func (parent *VBox) AddSpace(height int) {
	w := 0
	h := gui.NewQFontMetrics(gui.NewQFont()).AverageCharWidth() * height
	hPolicy := widgets.QSizePolicy__Fixed
	vPolicy := widgets.QSizePolicy__Fixed
	filler := widgets.NewQSpacerItem(w, h, hPolicy, vPolicy)
	parent.layout.AddItem(filler)
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

func (parent *Form) AddWidget(w *Widget) {
	if w.label != nil {
		parent.layout.AddRow(w.label, w.qWidget)
		return
	}

	parent.layout.AddWidget(w.qWidget)
}

func (parent *Form) AddRow(label string, w *Widget) {
	w.SetLabel(label)
	parent.AddWidget(w)
}

func (parent *Form) AddFieldRows(r *codeplug.Record, fTypes ...codeplug.FieldType) {
	for _, fType := range fTypes {
		parent.addFieldRow(r, fType)
	}
}

func (parent *Form) addFieldRow(r *codeplug.Record, fType codeplug.FieldType) {
	f := r.Field(fType)
	if f == nil {
		// This is not an error because some forms are used for
		// multiple models. We just ignore non-existent fields.
		return
	}

	w := newFieldWidget[f.ValueType()](f)
	w.label = widgets.NewQLabel2(f.TypeName(), nil, 0)
	parent.layout.AddRow(w.label, w.qWidget)

	widgets := parent.window.widgets
	widgets[fType] = w

	enablingFieldType := f.EnablingFieldType()

	w.receive = func(sender *Widget) {
		if sender.field.Record().Type() != w.field.Record().Type() {
			log.Fatal("sender record type", sender.field.Record().Type(), "receiver record type", w.field.Record().Type())
		}
		if sender.field.Record().Index() != w.field.Record().Index() {
			log.Fatal("sender record index", sender.field.Record().Index(), "receiver record index", w.field.Record().Index())
		}
		if sender.field.Index() != w.field.Index() {
			log.Fatal("sender field index", sender.field.Index(), "receiver field index", w.field.Index())
		}
		switch sender.field.Type() {
		case "":
			log.Fatal("receive(): invalid field type")

		case fType:
			w.update()
			subs := parent.window.subscriptions[fType]
			for _, sub := range subs {
				widgets[sub].receive(w)
			}

		case enablingFieldType:
			setEnabled(w, f)

		default:
			log.Fatal("receive(): unexpected field type")
		}
	}

	if enablingFieldType != "" {
		parent.subscribe(enablingFieldType, w.field.Type())
	}

}

func setEnabled(w *Widget, f *codeplug.Field) {
	enabled := f.IsEnabled()
	qWidget := w.qWidget.QWidget_PTR()
	if qWidget.IsEnabled() == enabled {
		return
	}

	if enabled && !f.IsValid() {
		f.SetString(f.DefaultValue())
	}

	qWidget.SetEnabled(enabled)
	w.label.SetEnabled(enabled)
	w.receive(w)
}

type Widget struct {
	qWidget widgets.QWidget_ITF
	label   *widgets.QLabel
	field   *codeplug.Field
	receive func(sender *Widget)
}

func (form *Form) subscribe(sender codeplug.FieldType, receiver codeplug.FieldType) {
	subs := form.window.subscriptions
	if subs[sender] == nil {
		subs[sender] = []codeplug.FieldType{}
	}
	subs[sender] = append(subs[sender], receiver)
}

func (w *Widget) SetLabel(label string) {
	w.label = widgets.NewQLabel2(label, nil, 0)
}

func (w *Widget) update() {
	f := w.field
	qw := w.qWidget

	switch qw.(type) {
	case *widgets.QCheckBox:
		setQCheckBox(qw.(*widgets.QCheckBox), f)

	case *widgets.QSpinBox:
		setQSpinBox(qw.(*widgets.QSpinBox), f)

	case *widgets.QLineEdit:
		qw.(*widgets.QLineEdit).SetText(f.String())

	case *widgets.QComboBox:
		qw.(*widgets.QComboBox).SetCurrentText(f.String())

	default:
		log.Fatal("update(): unexpected widget type")
	}
}

func (w *Widget) SetEnabled(b bool) {
	qw := w.qWidget

	switch qw.(type) {
	case *widgets.QComboBox:
		qw.(*widgets.QComboBox).SetEnabled(b)

	case *widgets.QPushButton:
		qw.(*widgets.QPushButton).SetEnabled(b)

	case *widgets.QCheckBox:
		qw.(*widgets.QCheckBox).SetEnabled(b)

	case *widgets.QSpinBox:
		qw.(*widgets.QSpinBox).SetEnabled(b)

	case *widgets.QLineEdit:
		qw.(*widgets.QLineEdit).SetEnabled(b)

	default:
		log.Fatal("SetEnabled(): unexpected widget type")
	}
}

func setQCheckBox(cb *widgets.QCheckBox, f *codeplug.Field) {
	checkState := core.Qt__Unchecked
	if f.String() == "On" {
		checkState = core.Qt__Checked
	}
	cb.SetCheckState(checkState)
}

func newFieldCheckbox(f *codeplug.Field) *Widget {
	qw := widgets.NewQCheckBox(nil)
	w := new(Widget)
	w.qWidget = qw
	w.field = f

	setQCheckBox(qw, f)

	qw.ConnectClicked(func(checked bool) {
		str := "Off"
		if checked {
			str = "On"
		}
		err := f.SetString(str)
		if err != nil {
			log.Fatal(err.Error())
		}
	})

	return w
}

func newFieldLineEdit(f *codeplug.Field) *Widget {
	qw := widgets.NewQLineEdit2(f.String(), nil)
	widget := new(Widget)
	widget.qWidget = qw
	widget.field = f

	var finished func()
	finished = func() {
		err := f.SetString(qw.Text())
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

func newFieldCombobox(f *codeplug.Field) *Widget {
	qw := widgets.NewQComboBox(nil)
	widget := new(Widget)
	widget.qWidget = qw
	widget.field = f

	strings := f.Strings()
	if len(strings) == 0 {
		log.Fatal("Combobox has no Strings()")
	}

	qw.InsertItems(0, strings)
	qw.SetCurrentText(f.String())

	qw.ConnectActivated2(func(str string) {
		err := f.SetString(str)
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
		if err != nil {
			log.Fatal("bad span string value")
		}
		value = int(i)
	}
	sb.SetValue(value)
}

func NewSpinboxWidget(value, min, max int, changedFunc func(int)) *Widget {
	qw := widgets.NewQSpinBox(nil)
	widget := new(Widget)
	widget.qWidget = qw
	qw.SetRange(min, max)
	qw.SetValue(value)

	qw.ConnectValueChanged(changedFunc)

	return widget
}

func NewComboboxWidget(opt string, opts []string, changed func(string)) *Widget {
	qw := widgets.NewQComboBox(nil)
	widget := new(Widget)
	widget.qWidget = qw
	qw.InsertItems(0, opts)
	qw.SetCurrentText(opt)

	qw.ConnectCurrentIndexChanged2(changed)

	return widget
}

func UpdateComboboxWidget(widget *Widget, opt string, opts []string) {
	qcb := widget.qWidget.(*widgets.QComboBox)
	qcb.Clear()
	qcb.InsertItems(0, opts)
	qcb.SetCurrentText(opt)
}

func NewButtonWidget(text string, clicked func()) *Widget {
	w := new(Widget)
	b := widgets.NewQPushButton2(text, nil)
	b.SetSizePolicy2(widgets.QSizePolicy__Fixed,
		widgets.QSizePolicy__Preferred)
	b.ConnectClicked(func(checked bool) {
		clicked()
	})
	w.qWidget = b

	return w
}

func newFieldSpinbox(f *codeplug.Field) *Widget {
	qw := widgets.NewQSpinBox(nil)
	widget := new(Widget)
	widget.qWidget = qw
	widget.field = f

	span := f.Span()
	qw.SetRange(span.Minimum(), span.Maximum())
	qw.SetSingleStep(span.Step())
	qw.SetWrapping(true)
	qw.SetSpecialValueText(span.MinString())

	setQSpinBox(qw, f)

	qw.ConnectValueChanged2(func(str string) {
		err := f.SetString(str)
		if err != nil {
			msg := f.TypeName() + " " + err.Error()
			ErrorPopup("Value error", msg)
		}
	})

	return widget
}

func newFieldTextEdit(f *codeplug.Field) *Widget {
	log.Fatal("newTextEdit: not implemented")
	return nil
}

var newFieldWidget = map[codeplug.ValueType]func(*codeplug.Field) *Widget{
	codeplug.VtAscii:           newFieldLineEdit,
	codeplug.VtCallID:          newFieldLineEdit,
	codeplug.VtCallType:        newFieldCombobox,
	codeplug.VtCtcssDcs:        newFieldCombobox,
	codeplug.VtFrequency:       newFieldLineEdit,
	codeplug.VtIndexedStrings:  newFieldCombobox,
	codeplug.VtIntroLine:       newFieldLineEdit,
	codeplug.VtIStrings:        newFieldCombobox,
	codeplug.VtListIndex:       newFieldCombobox,
	codeplug.VtMemberListIndex: newFieldCombobox,
	codeplug.VtName:            newFieldLineEdit,
	codeplug.VtOffOn:           newFieldCheckbox,
	codeplug.VtOnOff:           newFieldCheckbox,
	codeplug.VtPcPassword:      newFieldLineEdit,
	codeplug.VtPrivacyNumber:   newFieldLineEdit,
	codeplug.VtRadioName:       newFieldLineEdit,
	codeplug.VtRadioPassword:   newFieldLineEdit,
	codeplug.VtBiFrequency:     newFieldLineEdit,
	codeplug.VtSpan:            newFieldSpinbox,
	codeplug.VtTextMessage:     newFieldTextEdit,
	codeplug.VtTimeStamp:       newFieldLineEdit,
	codeplug.VtUniqueName:      newFieldLineEdit,
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
	qButton *widgets.QRadioButton
}

func (parent *HBox) AddRadioButton(option string) *RadioButton {
	b := new(RadioButton)
	b.qButton = widgets.NewQRadioButton2(option, nil)
	parent.layout.AddWidget(b.qButton, 0, 0)

	return b
}

func (b *RadioButton) ConnectClicked(fn func(checked bool)) {
	b.qButton.ConnectClicked(func(checked bool) {
		fn(checked)
	})
}

func (b *RadioButton) SetChecked(bo bool) {
	b.qButton.SetChecked(bo)
}

func (b *RadioButton) IsChecked() bool {
	return b.qButton.IsChecked()
}

func (b *RadioButton) Text() string {
	return b.qButton.Text()
}

type Button struct {
	qButton widgets.QPushButton
}

func NewButton(text string) *Button {
	b := new(Button)
	b.qButton = *widgets.NewQPushButton2(text, nil)
	b.qButton.SetSizePolicy2(widgets.QSizePolicy__Fixed,
		widgets.QSizePolicy__Preferred)

	return b
}

func (b *Button) ConnectClicked(fn func()) {
	b.qButton.ConnectClicked(func(checked bool) {
		fn()
	})
}

func (b *Button) SetText(str string) {
	b.qButton.SetText(str)
}

func (b *Button) SetEnabled(enable bool) {
	b.qButton.SetEnabled(enable)
}

func (w *Window) SetRecordFunc(fn func()) {
	w.recordFunc = fn
}

func InfoPopup(title string, msg string) {
	button := widgets.QMessageBox__Ok
	defaultButton := widgets.QMessageBox__Ok
	widgets.QMessageBox_Information(nil, title, msg, button, defaultButton)
}

func WarningPopup(title string, msg string) PopupValue {
	maxLines := 20
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

type Dialog struct {
	*VBox
	qDialog *widgets.QDialog
	layout  *widgets.QVBoxLayout
}

func NewDialog(title string) *Dialog {
	dialog := new(Dialog)
	dialog.VBox = newVbox()
	dialog.qDialog = widgets.NewQDialog(nil, core.Qt__Dialog)
	dialog.layout = widgets.NewQVBoxLayout2(dialog.qDialog)
	dialog.layout.AddWidget(&dialog.VBox.qWidget, 0, 0)

	dialog.qDialog.SetWindowTitle(title)

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

func OpenFilename(title string, dir string, exts []string) string {
	for i, ext := range exts {
		exts[i] = "*." + ext
	}
	selF := "(" + strings.Join(exts, " ") + ")"
	filter := "Codeplug files " + selF + ";;All files (*)"
	return widgets.QFileDialog_GetOpenFileName(nil, title, dir, filter, selF, 0)
}

func OpenFilenames(title string, dir string, exts []string) []string {
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

func ResetWindows(cp *codeplug.Codeplug) {
	for _, mw := range mainWindows {
		if mw.codeplug != cp {
			continue
		}

		for _, w := range mw.RecordWindows() {
			rl := w.RecordList()
			if rl != nil {
				rl.SetCurrent(0)
				rl.Update()
			}

			w.recordFunc()
		}
	}
}

func randomString(size int) (string, error) {
	rnd := make([]byte, size/2)
	_, err := rand.Read(rnd)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	writer := bufio.NewWriter(&buf)
	for _, b := range rnd {
		fmt.Fprintf(writer, "%02x", b)
	}
	writer.Flush()

	return buf.String(), nil
}
