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
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dalefarnsworth/codeplug/dfu"
	"github.com/tealeg/xlsx"
)

// FileType tells whether the codeplug is an rdt file or a bin file.
type FileType int

const (
	FileTypeNone FileType = iota
	FileTypeRdt
	FileTypeBin
	FileTypeNew
	FileTypeText
	FileTypeJSON
	FileTypeXLSX
)

const (
	MinProgress = dfu.MinProgress
	MaxProgress = dfu.MaxProgress
)

// A Codeplug represents a codeplug file.
type Codeplug struct {
	filename            string
	importFilename      string
	fileType            FileType
	rdtSize             int
	fileSize            int
	fileOffset          int
	id                  string
	bytes               []byte
	hash                [sha256.Size]byte
	rDesc               map[RecordType]*rDesc
	changed             bool
	lowFrequency        float64
	highFrequency       float64
	connectChange       func(*Change)
	changeList          []*Change
	changeIndex         int
	codeplugInfo        *CodeplugInfo
	loaded              bool
	cachedNameToRt      map[string]RecordType
	cachedNameToFt      map[RecordType]map[string]FieldType
	deferredValueFields []*Field
	errCount            int
	//warnings            []string
}

type CodeplugInfo struct {
	Type        string
	Models      []string
	Ext         string
	RdtSize     int
	BinSize     int
	BinOffset   int
	RecordInfos []*recordInfo
}

// NewCodeplug returns a Codeplug, given a filename and codeplug type.
func NewCodeplug(fType FileType, filename string) (*Codeplug, error) {
	cp := new(Codeplug)
	cp.fileType = fType

	switch fType {
	case FileTypeNone:
		err := cp.findFileType(filename)
		if err != nil {
			return nil, err
		}
		if cp.fileType == FileTypeRdt {
			err = cp.read(filename)
			if err != nil {
				cp.fileType = FileTypeNone
				return nil, err
			}
		}

	case FileTypeText, FileTypeJSON, FileTypeXLSX:
		cp.importFilename = filename
		fallthrough

	case FileTypeNew:
		baseName := "codeplug"
		for i := 1; ; i++ {
			filename = fmt.Sprintf("%s%d", baseName, i)

			found := false
			for _, cp := range codeplugs {
				if strings.HasPrefix(cp.filename, filename) {
					found = true
					break
				}
			}
			if !found {
				matches, err := filepath.Glob(filename + "*")
				if err != nil {
					logFatal(err.Error())
				}
				if len(matches) != 0 {
					found = true
					break
				}
			}
			if !found {
				break
			}
		}
	default:
		logFatal("unknown file type")
	}

	cp.filename = filename
	cp.rDesc = make(map[RecordType]*rDesc)
	cp.changeList = []*Change{&Change{}}
	cp.changeIndex = 0

	var err error
	cp.id, err = randomString(64)
	if err != nil {
		return nil, err
	}

	return cp, nil
}

type Warning struct {
	error
}

func (cp *Codeplug) Load(model string, frequencyRange string) error {
	notFound := fmt.Errorf("codeplug type not found: %s", model)

	found := false
findCodeplugInfo:
	for _, cpi := range codeplugInfos {
		for _, cpiModel := range cpi.Models {
			if cpiModel == model {
				cp.codeplugInfo = cpi
				found = true
				break findCodeplugInfo
			}
		}
	}
	if !found {
		return notFound
	}

	switch cp.fileType {
	case FileTypeNew, FileTypeBin, FileTypeText, FileTypeJSON, FileTypeXLSX:
		var filename string
		for i, v := range cp.frequencyRanges() {
			if v == frequencyRange {
				filename = cp.newFilenames()[i]
				break
			}
		}
		if filename == "" {
			return notFound
		}
		err := cp.readNew(filename)
		if err != nil {
			return err
		}
	}

	switch cp.fileType {
	case FileTypeRdt, FileTypeBin:
		err := cp.read(cp.filename)
		if err != nil {
			return err
		}
	}

	err := cp.Revert()
	if err != nil {
		return err
	}

	switch cp.fileType {
	case FileTypeText, FileTypeJSON, FileTypeXLSX:
		for _, rType := range cp.RecordTypes() {
			if cp.MaxRecords(rType) == 1 {
				continue
			}
			records := cp.records(rType)
			for i := len(records) - 1; i >= 0; i-- {
				cp.RemoveRecord(records[i])
			}
		}

		var err error
		switch cp.fileType {
		case FileTypeText:
			err = cp.importText(cp.importFilename)
		case FileTypeJSON:
			err = cp.importJSON(cp.importFilename)
		case FileTypeXLSX:
			err = cp.importXLSX(cp.importFilename)
		}
		if _, warning := err.(Warning); warning {
			err = nil
		}

		if err != nil {
			return err
		}
	}

	codeplugs = append(codeplugs, cp)
	cp.loaded = true

	return nil
}

func (cp *Codeplug) Loaded() bool {
	return cp.loaded
}

func (cp *Codeplug) AllExts() []string {
	extMap := make(map[string]bool)
	for _, cpi := range codeplugInfos {
		extMap[cpi.Ext] = true
	}
	exts := make([]string, 0, len(extMap))
	for ext := range extMap {
		exts = append(exts, ext)
	}

	return exts
}

func (cp *Codeplug) Ext() string {
	return cp.codeplugInfo.Ext
}

func AllFrequencyRanges() map[string][]string {
	models := make([]string, 0)
	freqRanges := make(map[string][]string)

	for _, cpi := range codeplugInfos {
		for _, model := range cpi.Models {
			for _, rInfo := range cpi.RecordInfos {
				if rInfo.rType == RtBasicInformation_md380 {
					for _, fInfo := range rInfo.fieldInfos {
						if fInfo.fType == FtBiFrequencyRange_md380 {
							freqRanges[model] = *fInfo.strings
						}
					}
				}
			}
			models = append(models, model)
		}
	}
	return freqRanges
}

// ModelsFrequencyRanges returns the potential codeplug model and
// frequencyRange
func (cp *Codeplug) ModelsFrequencyRanges() (models []string, frequencyRanges map[string][]string) {
	models = make([]string, 0)
	frequencyRanges = make(map[string][]string)
	var model string
	var frequencyRange string

	switch cp.fileType {
	case FileTypeRdt:

	case FileTypeText, FileTypeJSON, FileTypeXLSX:
		model, frequencyRange = cp.parseModelFrequencyRange()
		fallthrough
	default:
		cp.bytes = make([]byte, codeplugInfos[0].RdtSize)
	}

	for _, cpi := range codeplugInfos {
		cp.codeplugInfo = cpi
		cp.loadHeader()
		mainModel := cpi.Models[0]
		frequencyRanges[mainModel] = cp.frequencyRanges()

		if cp.fileType == FileTypeRdt {
			if cpi.RdtSize != cp.rdtSize {
				models = append(models, model)
				continue
			}
			model = cp.Model()
			frequencyRange = cp.FrequencyRange()
		}

		for _, cpiModel := range cpi.Models {
			if cpiModel == model {
				models = []string{mainModel}
				for _, v := range frequencyRanges[mainModel] {
					if v == frequencyRange {
						frequencyRanges[mainModel] = []string{v}
					}
				}
				return models, frequencyRanges
			}
		}
		models = append(models, mainModel)
	}

	if cp.fileType != FileTypeRdt {
		cp.bytes = nil
	}

	cp.codeplugInfo = nil

	return models, frequencyRanges
}

func (cp *Codeplug) Model() string {
	fDescs := cp.rDesc[RtBasicInformation_md380].records[0].fDesc
	return (*fDescs)[FtBiModel].fields[0].String()
}

func (cp *Codeplug) FrequencyRange() string {
	fDescs := cp.rDesc[RtBasicInformation_md380].records[0].fDesc
	return (*fDescs)[FtBiFrequencyRange_md380].fields[0].String()
}

func (cp *Codeplug) frequencyRanges() []string {
	for _, rInfo := range cp.codeplugInfo.RecordInfos {
		if rInfo.rType == RtBasicInformation_md380 {
			for _, fInfo := range rInfo.fieldInfos {
				if fInfo.fType == FtBiFrequencyRange_md380 {
					return *fInfo.strings
				}
			}
		}
	}
	return nil
}

func (cp *Codeplug) newFilenames() []string {
	for _, rInfo := range cp.codeplugInfo.RecordInfos {
		if rInfo.rType == RtBasicInformation_md380 {
			for _, fInfo := range rInfo.fieldInfos {
				if fInfo.fType == FtBiNewFilename_md380 {
					return *fInfo.strings
				}
			}
		}
	}
	return nil
}

func (cp *Codeplug) Type() string {
	return cp.codeplugInfo.Type
}

func (cp *Codeplug) ErrorCount() int {
	return cp.errCount
}

// Codeplugs returns a slice containing all currently open codeplugs.
func Codeplugs() []*Codeplug {
	return codeplugs
}

// Free frees a codeplug
func (cp *Codeplug) Free() {
	for i, codeplug := range codeplugs {
		if cp == codeplug {
			codeplugs = append(codeplugs[:i], codeplugs[i+1:]...)
			for _, rd := range cp.rDesc {
				rd.codeplug = nil
			}
			break
		}
	}
}

func (cp *Codeplug) readNew(filename string) error {
	gzipped := bytes.NewReader(new_tgz)

	archive, err := gzip.NewReader(gzipped)
	if err != nil {
		logFatal(err)
	}

	tarfile := tar.NewReader(archive)

	var bytes []byte
	for {
		hdr, err := tarfile.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			logFatal(err)
		}
		if hdr.Name != filename {
			continue
		}
		bytes, err = ioutil.ReadAll(tarfile)
		if err != nil {
			logFatal(err)
		}
		break
	}

	if len(bytes) == 0 {
		return fmt.Errorf("file %s not found", filename)
	}

	cp.bytes = bytes

	return nil
}

// read opens a file and reads its contents into cp.bytes.
func (cp *Codeplug) read(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if cp.bytes == nil {
		cp.bytes = make([]byte, cp.fileOffset+cp.fileSize)
	}
	bytes := cp.bytes[cp.fileOffset : cp.fileOffset+cp.fileSize]

	bytesRead, err := file.Read(bytes)
	if err != nil {
		return err
	}

	if bytesRead != cp.fileSize {
		err = fmt.Errorf("Failed to read all of %s", filename)
		return err
	}

	return nil
}

// Revert reverts the codeplug to its state after the most recent open or
// save operation.  An error is returned if the new codeplug state is
// invalid.
func (cp *Codeplug) Revert() error {
	cp.clearCachedListNames()

	cp.load()

	cp.Valid()

	cp.changed = false
	cp.hash = sha256.Sum256(cp.bytes)

	cp.changeList = []*Change{&Change{}}
	cp.changeIndex = 0

	return nil
}

// Save stores the state of the Codeplug into its file
// An error may be returned if the codeplug state is invalid.
func (cp *Codeplug) Save() error {
	return cp.SaveAs(cp.filename)
}

// SaveAs saves the state of the Codeplug into a named file.
// An error will be returned if the codeplug state is invalid.
// The named file becomes the current file associated with the codeplug.
func (cp *Codeplug) SaveAs(filename string) error {
	err := cp.SaveToFile(filename)
	if err != nil {
		return err
	}

	cp.filename = filename
	cp.changed = false
	cp.hash = sha256.Sum256(cp.bytes)

	return nil
}

// SaveToFile saves the state of the Codeplug into a named file.
// An error will be returned if the codeplug state is invalid.
// The state of the codeplug is not changed, so this
// is useful for use by an autosave function.
func (cp *Codeplug) SaveToFile(filename string) (err error) {
	if err = cp.Valid(); err != nil {
		_, warning := err.(Warning)
		if !warning {
			return err
		}
	}

	cp.setLastProgrammedTime(time.Now())

	cp.store()

	dir, base := filepath.Split(filename)
	tmpFile, err := ioutil.TempFile(dir, base)
	if err != nil {
		return err
	}
	tmpFilename := tmpFile.Name()

	defer func() {
		closeErr := tmpFile.Close()
		if err == nil {
			err = closeErr
		}

		if err != nil {
			os.Remove(tmpFilename)
			return
		}

		err = os.Rename(tmpFilename, filename)
	}()

	cpi := cp.codeplugInfo
	fileSize := cpi.RdtSize
	fileOffset := 0

	bytes := cp.bytes[fileOffset : fileOffset+fileSize]
	bytesWritten, err := tmpFile.Write(bytes)
	if err != nil {
		return err
	}

	if bytesWritten != fileSize {
		return fmt.Errorf("write to %s failed", cp.filename)
	}

	return err
}

func (cp *Codeplug) setLastProgrammedTime(t time.Time) {
	r := cp.rDesc[RtBasicInformation_md380].records[0]
	f := r.Field(FtBiLastProgrammedTime)
	f.setString(t.Format("02-Jan-2006 15:04:05"))
}

func (cp *Codeplug) getLastProgrammedTime() (time.Time, error) {
	r := cp.rDesc[RtBasicInformation_md380].records[0]
	f := r.Field(FtBiLastProgrammedTime)
	return time.Parse("02-Jan-2006 15:04:05", f.String())
}

// Filename returns the path name of the file associated with the codeplug.
// This is the file named in the most recent Open or SaveAs function.
func (cp *Codeplug) Filename() string {
	return cp.filename
}

// CurrentHash returns a cryptographic hash of the current (modified) codeplug
func (cp *Codeplug) CurrentHash() [sha256.Size]byte {
	if !cp.changed {
		return cp.hash
	}

	bytes := make([]byte, len(cp.bytes))
	copy(bytes, cp.bytes)
	saveBytes := cp.bytes
	cp.bytes = bytes
	cp.store()
	cp.bytes = saveBytes

	return sha256.Sum256(bytes)
}

// Changed returns false if the codeplug state is the same as that at
// the most recent Open or Save/SaveAs operation.
func (cp *Codeplug) Changed() bool {
	if cp.changed && cp.CurrentHash() != cp.hash {
		return true
	}

	return false
}

func (cp *Codeplug) SetChanged() {
	cp.changed = true
	for i := range cp.hash {
		cp.hash[i] = 0
	}
}

// FileType returns the type of codeplug file (rdt or bin).
func (cp *Codeplug) FileType() FileType {
	return cp.fileType
}

// Records returns all of a codeplug's records of the given RecordType.
func (cp *Codeplug) Records(rType RecordType) []*Record {
	records := cp.rDesc[rType].records
	if len(records) == 0 {
		rIndex := 0
		r := cp.newRecord(rType, rIndex)
		cp.InsertRecord(r)
		r.load()
		nameField := r.NameField()
		if nameField != nil {
			nameField.setString(string(rType) + "1")
		}
		records = cp.rDesc[rType].records
	}
	return records
}

func (cp *Codeplug) records(rType RecordType) []*Record {
	return cp.rDesc[rType].records
}

// Record returns the first record of a codeplug's given RecordType.
func (cp *Codeplug) Record(rType RecordType) *Record {
	return cp.Records(rType)[0]
}

// record returns the first record of a codeplug's given RecordType.
func (cp *Codeplug) record(rType RecordType) *Record {
	return cp.records(rType)[0]
}

// MaxRecords returns a codeplug's maximum number of records of the given
// Recordtype.
func (cp *Codeplug) MaxRecords(rType RecordType) int {
	return cp.rDesc[rType].max
}

// RecordTypes returns all of the record types of the codeplug except
// BasicInformation.  The BasicInformation record is omitted.
func (cp *Codeplug) RecordTypes() []RecordType {
	indexedStrs := make(map[int]string)
	indexes := make([]int, 0, len(cp.rDesc))

	for rType, rDesc := range cp.rDesc {
		index := rDesc.recordInfo.index
		indexes = append(indexes, index)
		indexedStrs[index] = string(rType)
	}
	sort.Ints(indexes)

	rTypes := make([]RecordType, len(indexes))
	for i, index := range indexes {
		rTypes[i] = RecordType(indexedStrs[index])
	}

	return rTypes
}

// ID returns a string unique to the codeplug.
func (cp *Codeplug) ID() string {
	return cp.id
}

// MoveRecord moves a record from its current slice index to the given index.
func (cp *Codeplug) MoveRecord(dIndex int, r *Record) {
	sIndex := r.rIndex
	cp.RemoveRecord(r)
	if sIndex < dIndex {
		dIndex--
	}
	r.rIndex = dIndex
	cp.InsertRecord(r)
}

// InsertRecord inserts the given record into the codeplug.
// The record's index determines the slice index at which it will be inserted.
// If the name of the record matches that of an existing record,
// the name is modified to make it unique.  An error will be returned if
// the codeplug's maximum records of that type would be exceeded.
func (cp *Codeplug) InsertRecord(r *Record) error {
	rType := r.rType
	records := cp.records(rType)
	if len(records) >= cp.MaxRecords(rType) {
		return fmt.Errorf("too many records")
	}

	if r.hasUniqueNames() {
		err := r.makeNameUnique(r.ListNames())
		if err != nil {
			return err
		}
	}

	i := r.rIndex
	records = append(records[:i], append([]*Record{r}, records[i:]...)...)

	for i, r := range records {
		r.rIndex = i
	}
	cp.rDesc[rType].records = records

	records[0].cachedListNames = nil
	return nil
}

// RemoveRecord removes the given record from the codeplug.
func (cp *Codeplug) RemoveRecord(r *Record) {
	rType := r.rType
	index := -1
	records := cp.records(rType)
	for i, record := range records {
		if record == r {
			index = i
			break
		}
	}
	if index < 0 || index >= len(records) {
		logFatal("removeRecord: bad record")
	}

	records[0].cachedListNames = nil

	deleteRecord(&records, index)

	for i, r := range records {
		r.rIndex = i
	}
	cp.rDesc[rType].records = records
}

// ConnectChange will cause the given function to be called passing
// the given change.
func (cp *Codeplug) ConnectChange(fn func(*Change)) {
	cp.connectChange = fn
}

// loadHeader loads the rdt header into the codeplug from its file.
func (cp *Codeplug) loadHeader() {
	cp.clearCachedListNames()
	ri := cp.codeplugInfo.RecordInfos[0]
	ri.max = 1

	rd := &rDesc{recordInfo: ri}
	cp.rDesc[ri.rType] = rd
	rd.codeplug = cp
	rd.loadRecords()
}

// load loads all the records into the codeplug from its file.
func (cp *Codeplug) load() {
	cp.clearCachedListNames()
	for i, ri := range cp.codeplugInfo.RecordInfos {
		ri.index = i
		if ri.max == 0 {
			ri.max = 1
		}

		rd := &rDesc{recordInfo: ri}
		cp.rDesc[ri.rType] = rd
		rd.codeplug = cp
		rd.loadRecords()
	}
}

// newRecord creates and returns the address of a new record of the given type.
func (cp *Codeplug) newRecord(rType RecordType, rIndex int) *Record {
	r := new(Record)
	r.rDesc = cp.rDesc[rType]
	r.rIndex = rIndex
	m := make(map[FieldType]*fDesc)
	r.fDesc = &m

	return r
}

// valid returns nil if all fields in the codeplug are valid.
func (cp *Codeplug) Valid() error {
	errStr := ""
	cp.errCount = 0
	for _, rType := range cp.RecordTypes() {
		for _, r := range cp.records(rType) {
			if err := r.valid(); err != nil {
				errStr += err.Error()
				cp.errCount++
			}
		}
	}

	for _, f := range cp.deferredValueFields {
		if err := f.valid(); err != nil {
			errStr += fmt.Sprintf("%s %s\n", f.FullTypeName(), err.Error())
			cp.errCount++
		}
	}
	cp.deferredValueFields = nil

	if errStr != "" {
		return Warning{fmt.Errorf("%s", errStr)}
	}

	return nil
}

// findFileType sets the codeplug type based on file size
func (cp *Codeplug) findFileType(filename string) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("%s: does not exist", filename)
		}
		cp.fileType = FileTypeNone
		return err
	}

	for _, cpi := range codeplugInfos {
		cp.rdtSize = cpi.RdtSize
		switch fileInfo.Size() {
		case int64(cpi.RdtSize):
			cp.fileType = FileTypeRdt
			cp.fileSize = cpi.RdtSize
			cp.fileOffset = 0
			return nil

		case int64(cpi.BinSize):
			cp.fileType = FileTypeBin
			cp.fileSize = cpi.BinSize
			cp.fileOffset = cpi.BinOffset
			return nil
		}
	}

	cp.fileType = FileTypeNone
	err = fmt.Errorf("%s is not a known codeplug file type", filename)
	return err
}

// store stores all all fields of the codeplug into its byte slice.
func (cp *Codeplug) store() {
	for _, rd := range cp.rDesc {
		for rIndex := 0; rIndex < rd.max; rIndex++ {
			if rIndex < len(rd.records) {
				rd.records[rIndex].store()
			} else {
				rd.deleteRecord(cp, rIndex)
			}
		}
	}
}

// frequencyValid returns nil if the given frequency is valid for the
// codeplug.
func (cp *Codeplug) frequencyValid(freq float64) error {
	if cp.lowFrequency == 0 {
		fDescs := cp.record(RtBasicInformation_md380).fDesc
		s := (*fDescs)[FtBiLowFrequency].fields[0].String()
		cp.lowFrequency, _ = strconv.ParseFloat(s, 64)
		s = (*fDescs)[FtBiHighFrequency].fields[0].String()
		cp.highFrequency, _ = strconv.ParseFloat(s, 64)
	}

	if freq >= cp.lowFrequency && freq <= cp.highFrequency {
		return nil
	}

	return fmt.Errorf("frequency out of range %v", freq)
}

// publishChange passes the given change (with any additional generated
// changes resulting from that change) to a registered function.
func (cp *Codeplug) publishChange(change *Change) {
	if cp.connectChange != nil {
		cp.connectChange(change)
	}
}

// codeplugs contains the list of open codeplugs.
var codeplugs []*Codeplug

func filterField(rType RecordType, fType FieldType) bool {
	switch rType {
	case RtBasicInformation_md380:
		switch fType {
		case FtBiCpsVersion, FtBiNewFilename_md380:
			return true
		case FtBiLowFrequency, FtBiHighFrequency:
			return true
		}
	}
	return false
}

func PrintRecord(w io.Writer, r *Record) {
	rType := r.Type()

	ind := ""
	if r.max > 1 {
		ind = fmt.Sprintf("[%d]", r.rIndex+1)
	}

	fmt.Fprintf(w, "%s%s:\n", string(rType), ind)

	for i, fType := range r.FieldTypes() {
		if filterField(rType, fType) {
			continue
		}

		if i == 0 {
		}

		name := string(fType)
		for _, f := range r.Fields(fType) {
			value := quoteString(f.String())
			ind := ""
			if f.max > 1 {
				ind = fmt.Sprintf("[%d]", f.fIndex+1)
			}
			fmt.Fprintf(w, "\t%s%s: %s\n", name, ind, value)
		}
	}
}

func PrintRecordWithIndex(w io.Writer, r *Record) {
	rType := r.Type()
	ind := ""
	if r.max > 1 {
		ind = fmt.Sprintf("[%d]", r.rIndex+1)
	}
	fmt.Fprintf(w, "%s%s:", string(rType), ind)

	for _, fType := range r.FieldTypes() {
		name := string(fType)
		for _, f := range r.Fields(fType) {
			value := quoteString(f.String())
			ind := ""
			if f.max > 1 {
				ind = fmt.Sprintf("[%d]", f.fIndex+1)
			}
			fmt.Fprintf(w, " %s%s:%s", name, ind, value)
		}
	}
	fmt.Fprintln(w)
}

func quoteString(str string) string {
	quote := false
	for _, c := range str {
		if unicode.IsSpace(c) {
			quote = true
			break
		}
	}
	if quote {
		runes := []rune{}
		for _, c := range str {
			switch c {
			case '"', '\\', '\n', '\t', '\r':
				runes = append(runes, '\\')
				switch c {
				case '\n':
					c = 'n'
				case '\t':
					c = 't'
				case '\r':
					c = 'r'
				}
			}
			runes = append(runes, c)
		}
		str = string(runes)
	}
	if quote || str == "" {
		str = `"` + str + `"`
	}

	return str
}

func printRecord(w io.Writer, r *Record, rFmt string, fFmt string) {
	rType := r.Type()
	ind := ""
	if r.max > 1 {
		ind = fmt.Sprintf("[%d]", r.Index()+1)
	}
	fmt.Fprintf(w, rFmt, string(rType), ind)

	for _, fType := range r.FieldTypes() {
		name := string(fType)
		for _, f := range r.Fields(fType) {
			value := quoteString(f.String())
			ind := ""
			if f.max > 1 {
				ind = fmt.Sprintf("[%d]", f.Index()+1)
			}
			fmt.Fprintf(w, fFmt, name, ind, value)
		}
	}
}
func (cp *Codeplug) nameToRt(rTypeName string) (RecordType, error) {
	if len(cp.cachedNameToRt) == 0 {
		cp.cachedNameToRt = make(map[string]RecordType)
		for _, rType := range cp.RecordTypes() {
			cp.cachedNameToRt[string(rType)] = rType
		}
	}

	var rType RecordType
	var ok bool

	rType, ok = cp.cachedNameToRt[rTypeName]
	if !ok {
		return rType, fmt.Errorf("unknown record type: %s", rTypeName)
	}

	return rType, nil
}

func (cp *Codeplug) nameToFt(rType RecordType, fTypeName string) (FieldType, error) {
	if len(cp.cachedNameToFt) == 0 {
		cp.cachedNameToFt = make(map[RecordType]map[string]FieldType)
		for _, rType := range cp.RecordTypes() {
			m := make(map[string]FieldType)
			for _, fi := range cp.rDesc[rType].fieldInfos {
				fType := fi.fType
				name := string(fType)
				m[name] = fType
			}
			cp.cachedNameToFt[rType] = m
		}
	}

	var fType FieldType
	var ok bool

	fType, ok = cp.cachedNameToFt[rType][fTypeName]
	if !ok {
		return "", fmt.Errorf("unknown field type: %s", fTypeName)
	}

	return fType, nil
}

type reader struct {
	*bufio.Reader
	pos     position
	prevPos position
}

func NewReader(ioReader io.Reader) *reader {
	rdr := new(reader)
	rdr.Reader = bufio.NewReader(ioReader)
	return rdr
}

func (rdr *reader) ReadRune() (rune, int, error) {
	r, size, err := rdr.Reader.ReadRune()
	if err != nil {
		return r, size, err
	}

	rdr.prevPos = rdr.pos

	switch r {
	case '\n':
		rdr.pos.line++
		rdr.pos.column = 0

	case '\t':
		rdr.pos.column = rdr.pos.column%8 + 8
	default:
		rdr.pos.column++
	}

	return r, size, err
}

func (rdr *reader) UnreadRune() error {
	rdr.pos = rdr.prevPos
	return rdr.Reader.UnreadRune()
}

func (rdr *reader) ReadUntil(f func(rune) bool) (string, error) {
	var err error
	var r rune
	runes := []rune{}

	for {
		r, _, err = rdr.ReadRune()
		if err != nil {
			break
		}
		if f(r) {
			rdr.UnreadRune()
			break
		}
		runes = append(runes, r)
	}

	return string(runes), err
}

func (rdr *reader) ReadWhile(f func(rune) bool) (string, error) {
	return rdr.ReadUntil(func(r rune) bool {
		return !f(r)
	})
}

func (rdr *reader) ReadEscapedUntil(f func(rune) bool) (string, error) {
	var err error
	var r rune
	runes := []rune{}

	for ; ; runes = append(runes, r) {
		r, _, err = rdr.ReadRune()
		if err != nil {
			break
		}
		if r == '\\' {
			r, _, err = rdr.ReadRune()
			if err != nil {
				break
			}
			switch r {
			case 'n':
				r = '\n'
			case 't':
				r = '\t'
			case 'r':
				r = '\r'
			}
			continue
		}
		if f(r) {
			rdr.UnreadRune()
			break
		}
	}

	return string(runes), err
}

func (rdr *reader) ReadInt() (int, error) {
	str, err := rdr.ReadWhile(unicode.IsDigit)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(n), nil
}

type position struct {
	line   int
	column int
}

type PositionError struct {
	error
	position *position
}

func (e *PositionError) Line() int {
	return e.position.line + 1
}

func (e *PositionError) Column() int {
	return e.position.column + 1
}

func (cp *Codeplug) parseModelFrequencyRange() (model string, frequencyRange string) {
	file, err := os.Open(cp.importFilename)
	if err != nil {
		return model, frequencyRange
	}
	defer file.Close()

	var pRecs []*parsedRecord

	switch cp.fileType {
	case FileTypeText:
		pRecs = cp.parseTextFile(file)

	case FileTypeJSON:
		pRecs = cp.parseJSONFile(file)

	case FileTypeXLSX:
		pRecs = cp.parseXLSXFile(file)
	}

	for _, pr := range pRecs {
		if pr.name != string(RtBasicInformation_md380) {
			continue
		}
		for _, pf := range pr.pFields {
			switch pf.name {
			case string(FtBiModel):
				model = pf.value
			case string(FtBiFrequencyRange_md380):
				frequencyRange = pf.value
			default:
				continue
			}
			if model != "" && frequencyRange != "" {
				return model, frequencyRange
			}
		}
	}

	return model, frequencyRange
}

type parsedField struct {
	name  string
	index int
	err   error
	pos   *position
	value string
}

type parsedRecord struct {
	name    string
	index   int
	err     error
	pos     *position
	pFields []*parsedField
}

func (cp *Codeplug) parseTextFile(iRdr io.Reader) []*parsedRecord {
	var index int
	var err error
	var pRecords []*parsedRecord

	rdr := NewReader(iRdr)

	rdr.ReadWhile(unicode.IsSpace)

	for {
		var pRecord parsedRecord
		pRecords = append(pRecords, &pRecord)
		pos := rdr.pos
		pRecord.pos = &pos
		pRecord.name, index, err = parseName(rdr)
		if err != nil {
			if err == io.EOF {
				err = nil
				pRecords = pRecords[:len(pRecords)-1]
				break
			}
			pRecord.err = err
			break
		}
		pRecord.index = index
		if len(pRecord.name) == 0 {
			pRecord.err = fmt.Errorf("syntax: no record name")
			break
		}
		var pFields []*parsedField
		for {
			if rdr.pos.column == 0 {
				break
			}

			var pField parsedField
			pFields = append(pFields, &pField)

			pos := rdr.pos
			pField.pos = &pos

			pField.name, index, err = parseName(rdr)
			pField.index = index
			if err != nil {
				pField.err = err
				break
			}
			if len(pField.name) == 0 {
				pField.err = fmt.Errorf("syntax: no field name")
				break
			}

			pos = rdr.pos
			pField.pos = &pos

			pField.value, err = parseValue(rdr)
			if err != nil {
				pField.err = fmt.Errorf("syntax: value: %s: %s: %s", pField.name, pField.value, err.Error())
				break
			}
		}
		pRecord.pFields = pFields
	}

	return pRecords
}

func parseName(rdr *reader) (string, int, error) {
	pos := rdr.pos
	nType := "record"
	if pos.column != 0 {
		nType = "field"
	}

	var r rune
	index := 0

	name, err := rdr.ReadWhile(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	})

	wrapError := func(err error) (string, int, error) {
		pErr := PositionError{
			position: &pos,
			error:    err,
		}
		return name, 0, pErr
	}

	if err != nil {
		if err == io.EOF {
			return name, index, err
		}

		err = fmt.Errorf("bad %s name", nType)
		return wrapError(err)
	}
	if len(name) == 0 || !unicode.IsLetter([]rune(name)[0]) {
		err = fmt.Errorf("bad %s name", nType)
		return wrapError(err)
	}

	pos = rdr.pos
	r, _, err = rdr.ReadRune()
	if err != nil {
		err = fmt.Errorf("bad %s name", nType)
		return wrapError(err)
	}
	switch r {
	case ':':
	case '[':
		pos = rdr.pos
		index, err = rdr.ReadInt()
		if err != nil {
			err = fmt.Errorf("bad %s index", nType)
			return wrapError(err)
		}
		index--
		pos = rdr.pos
		r, _, err = rdr.ReadRune()
		if r != ']' {
			err = fmt.Errorf("bad %s index", nType)
			return wrapError(err)
		}
		pos = rdr.pos
		r, _, err = rdr.ReadRune()
		if r != ':' {
			err = fmt.Errorf("bad %s name", nType)
			return wrapError(err)
		}
	default:
		err = fmt.Errorf("bad %s index", nType)
		return wrapError(err)
	}
	rdr.ReadWhile(unicode.IsSpace)

	return name, index, nil
}

func (cp *Codeplug) ParseRecords(rdr io.Reader) ([]*Record, error) {
	pRecs := cp.parseTextFile(rdr)
	records, err := cp.parsedFileToRecs(pRecs)

	return records, err
}

func (cp *Codeplug) parsedFileToRecs(pRecs []*parsedRecord) ([]*Record, error) {
	var warning error
	var err error
	var pos *position

	var records []*Record
	appendWarning := func(pr *parsedRecord, pf *parsedField, err error) {
		err = fmt.Errorf("%s.%s: %s", pr.name, pf.name, err.Error())
		appendWarningMsgs(&warning, pf.pos, err)
	}

	wrapError := func(err error) ([]*Record, error) {
		if err != nil {
			pErr, ok := err.(PositionError)
			if !ok {
				pErr = PositionError{
					position: pos,
					error:    err,
				}
			}
			return records, pErr
		}

		return records, err
	}

	for _, pr := range pRecs {
		if pr.err != nil {
			pos = pr.pos
			err = pr.err
			return wrapError(err)
		}

		var r *Record
		r, err = cp.rNameToRecord(pr.name, pr.index)
		if err != nil {
			pos = pr.pos
			return wrapError(err)
		}

		for _, pf := range pr.pFields {
			if pf.err != nil {
				pos = pr.pos
				err = pf.err
				return wrapError(err)
			}

			fType, err := cp.nameToFt(r.rType, pf.name)
			if err != nil {
				appendWarning(pr, pf, err)
				continue
			}

			var f *Field
			f, err = r.NewFieldWithValue(fType, pf.index, pf.value)
			if err != nil {
				appendWarning(pr, pf, err)
			}
			dValue, ok := f.value.(deferredValue)
			if ok {
				dValue.str = pf.value
				if pf.pos != nil {
					dValue.pos = pf.pos
				}
				f.value = dValue
			}
			err = r.addField(f)
			if err != nil {
				appendWarning(pr, pf, err)
			}
		}
		records = append(records, r)
	}

	return records, warning
}

func appendWarningMsgs(pWarning *error, ppos *position, warning error) {
	var oldMsg string
	if *pWarning != nil {
		oldMsg = (*pWarning).Error()
	}
	err := fmt.Errorf("%s%s\n", oldMsg, warning.Error())
	if ppos != nil {
		pos := *ppos
		err = fmt.Errorf("%sline %d:%d: %s\n", oldMsg,
			pos.line+1, pos.column+1, warning.Error())
	}
	*pWarning = Warning{err}
}

func parseValue(rdr *reader) (string, error) {
	pos := rdr.pos

	r, _, err := rdr.ReadRune()
	if err != nil {
		pErr := PositionError{
			position: &pos,
			error:    err,
		}
		return "", pErr
	}
	termFunc := unicode.IsSpace
	if r == '"' {
		termFunc = func(r rune) bool {
			return r == '"'
		}
	} else {
		rdr.UnreadRune()
	}

	value, err := rdr.ReadEscapedUntil(termFunc)
	if err != nil {
		pErr := PositionError{
			position: &pos,
			error:    err,
		}
		return value, pErr
	}
	rdr.ReadRune()

	rdr.ReadWhile(unicode.IsSpace)

	return value, nil
}

func (cp *Codeplug) rNameToRecord(name string, index int) (*Record, error) {
	rType, err := cp.nameToRt(name)
	if err != nil {
		return nil, err
	}

	found := false
	for rt := range cp.rDesc {
		if rt == rType {
			found = true
			break
		}
	}
	if found {
		return cp.newRecord(rType, index), nil
	}

	return nil, fmt.Errorf("codeplug has no record: %s", string(rType))
}

func (cp *Codeplug) ExportText(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		fErr := file.Close()
		if err == nil {
			err = fErr
		}
		return
	}()

	w := bufio.NewWriter(file)
	for i, rType := range cp.RecordTypes() {
		for j, r := range cp.records(rType) {
			if i != 0 || j != 0 {
				fmt.Fprintln(w)
			}
			PrintRecord(w, r)
		}
	}
	w.Flush()

	return nil
}

func (cp *Codeplug) importText(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	records, err := cp.ParseRecords(file)
	err = cp.storeParsedRecords(records)
	if err != nil {
		return err
	}

	return nil
}

func (cp *Codeplug) ExportJSON(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		fErr := file.Close()
		if err == nil {
			err = fErr
		}
		return
	}()

	recordTypes := cp.RecordTypes()
	recordMap := make(map[string]interface{})
	for _, rType := range recordTypes {
		records := cp.records(rType)
		recordSlice := make([]map[string]interface{}, len(records))
		for i, r := range records {
			fieldTypes := r.FieldTypes()
			fieldMap := make(map[string]interface{})
			for _, fType := range fieldTypes {
				if filterField(rType, fType) {
					continue
				}
				fields := r.Fields(fType)

				fieldSlice := make([]string, len(fields))
				for j, f := range fields {
					fieldSlice[j] = f.String()
				}

				fTypeString := string(fType)
				fieldMap[fTypeString] = fieldSlice
				if r.MaxFields(fType) == 1 {
					fieldMap[fTypeString] = fieldSlice[0]
				}
			}
			recordSlice[i] = fieldMap
		}
		rTypeString := string(rType)
		recordMap[rTypeString] = recordSlice
		if cp.MaxRecords(rType) == 1 {
			recordMap[rTypeString] = recordSlice[0]
		}
	}

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(recordMap)
	if err != nil {
		return err
	}
	writer.Flush()

	return nil
}

func (cp *Codeplug) importJSON(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	pRecs := cp.parseJSONFile(file)
	records, err := cp.parsedFileToRecs(pRecs)
	if err != nil {
		return err
	}

	err = cp.storeParsedRecords(records)
	if err != nil {
		return err
	}

	return nil
}

func (cp *Codeplug) parseJSONFile(iRdr io.Reader) []*parsedRecord {
	errorInvalidJSON := fmt.Errorf("Invalid codeplug JSON file")

	var parRecord parsedRecord
	parRecords := []*parsedRecord{
		&parRecord,
	}

	reader := bufio.NewReader(iRdr)
	decoder := json.NewDecoder(reader)
	var i interface{}
	err := decoder.Decode(&i)
	if err != nil {
		parRecord.err = errorInvalidJSON
		return parRecords
	}
	recordMap, ok := i.(map[string]interface{})
	if !ok {
		parRecord.err = errorInvalidJSON
		return parRecords
	}

	var pRecords []*parsedRecord

	for rName, i := range recordMap {
		var recordSlice []map[string]interface{}

		switch v := i.(type) {
		case []interface{}:
			recordSlice = make([]map[string]interface{}, len(v))
			for i := range v {
				record, ok := v[i].(map[string]interface{})
				if !ok {
					parRecord.err = errorInvalidJSON
					return parRecords
				}
				recordSlice[i] = record
			}

		case map[string]interface{}:
			recordSlice = []map[string]interface{}{v}

		default:
			parRecord.err = errorInvalidJSON
			return parRecords
		}

		for index, fMap := range recordSlice {
			var pRecord parsedRecord
			pRecord.index = index
			pRecord.name = rName
			pRecord.pFields = cp.parseJSONFields(fMap)
			pRecords = append(pRecords, &pRecord)
		}
	}

	return pRecords
}

func (cp *Codeplug) parseJSONFields(fMap map[string]interface{}) []*parsedField {
	errorInvalidJSON := fmt.Errorf("Invalid codeplug JSON file")
	var parField parsedField
	parFields := []*parsedField{
		&parField,
	}

	var pFields []*parsedField
	for fName, i := range fMap {
		var fieldSlice []string
		switch v := i.(type) {
		case []interface{}:
			fieldSlice = make([]string, len(v))
			for i := range v {
				field, ok := v[i].(string)
				if !ok {
					parField.err = errorInvalidJSON
					return parFields
				}
				fieldSlice[i] = field
			}

		case string:
			fieldSlice = []string{v}

		default:
			parField.err = errorInvalidJSON
			return parFields
		}

		for index, str := range fieldSlice {
			pField := parsedField{
				name:  fName,
				index: index,
				value: str,
			}
			pFields = append(pFields, &pField)
		}
	}

	return pFields
}

func (cp *Codeplug) ExportXLSX(filename string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()

	recordTypes := cp.RecordTypes()
	for _, rType := range recordTypes {
		sheet, err = file.AddSheet(string(rType))
		if err != nil {
			return err
		}
		records := cp.records(rType)
		headerRow := sheet.AddRow()
		r := records[0]
		for _, fType := range r.FieldTypes() {
			if filterField(rType, fType) {
				continue
			}

			for i := 0; i < (*r.fDesc)[fType].max; i++ {
				cell = headerRow.AddCell()
				cell.Value = string(fType)
			}
		}

		for _, r := range records {
			row = sheet.AddRow()
			for _, fType := range r.FieldTypes() {
				if filterField(rType, fType) {
					continue
				}

				fields := r.Fields(fType)
				for _, f := range fields {
					cell = row.AddCell()
					cell.Value = f.String()
				}
			}
		}
	}
	err = file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func (cp *Codeplug) importXLSX(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	pRecs := cp.parseXLSXFile(file)
	records, err := cp.parsedFileToRecs(pRecs)
	if err != nil {
		return err
	}

	err = cp.storeParsedRecords(records)
	if err != nil {
		return err
	}

	return nil
}

func (cp *Codeplug) parseXLSXFile(iRdr io.Reader) []*parsedRecord {
	errorInvalidXLSX := fmt.Errorf("Invalid codeplug spreadsheet file")

	var parRecord parsedRecord
	parRecords := []*parsedRecord{
		&parRecord,
	}

	bytes, err := ioutil.ReadAll(iRdr)
	if err != nil {
		parRecord.err = errorInvalidXLSX
		return parRecords
	}

	file, err := xlsx.OpenBinary(bytes)
	if err != nil {
		parRecord.err = errorInvalidXLSX
		return parRecords
	}

	var pRecords []*parsedRecord

	for _, sheet := range file.Sheets {
		rTypeName := sheet.Name

		headerCells := sheet.Rows[0].Cells
		fTypeNames := make([]string, len(headerCells))
		for i, cell := range headerCells {
			fTypeNames[i] = cell.String()
		}
		for index, row := range sheet.Rows[1:] {
			pRecord := parsedRecord{
				name:  rTypeName,
				index: index,
			}
			var parFields []*parsedField
			for index, cell := range row.Cells {
				strs := cell.String()
				for _, str := range strings.Split(strs, "\n") {
					parField := parsedField{
						name:  fTypeNames[index],
						value: str,
					}
					parFields = append(parFields, &parField)
				}
				pRecord.pFields = parFields
			}
			pRecords = append(pRecords, &pRecord)
		}
	}

	return pRecords
}

func (cp *Codeplug) storeParsedRecords(records []*Record) error {
	seenType := make(map[RecordType]bool)
	for _, r := range records {
		if !seenType[r.rType] {
			seenType[r.rType] = true
			cp.rDesc[r.rType].records = []*Record{}
		}
		r.rIndex = len(cp.records(r.rType))
		cp.InsertRecord(r)
	}

	err := cp.resolveDeferredValueFields()
	if err != nil {
		return err
	}

	for _, rd := range cp.rDesc {
		if len(rd.records) == 0 {
			return fmt.Errorf("no %s found", rd.rType)
		}
	}
	cp.store()
	cp.changeList = []*Change{&Change{}}
	cp.changeIndex = 0
	cp.changed = true

	return nil
}

func (cp *Codeplug) clearCachedListNames() {
	for _, rd := range cp.rDesc {
		rd.cachedListNames = nil
	}
}

func (cp *Codeplug) resolveDeferredValueFields() error {
	var warning error
	appendWarning := func(f *Field, pos *position, err error) {
		rName := f.record.TypeName()
		fName := f.TypeName()
		err = fmt.Errorf("%s.%s: %s", rName, fName, err.Error())
		appendWarningMsgs(&warning, pos, err)
	}

	for i := 0; len(cp.deferredValueFields) > 0 && i < 10; i++ {
		deferredValueFields := cp.deferredValueFields
		cp.deferredValueFields = nil
		for _, f := range deferredValueFields {
			dValue, deferred := f.value.(deferredValue)
			if !deferred {
				logFatal("not deferred", f.FullTypeName())
			}

			if f.isDeferredValue(dValue.str) {
				continue
			}

			f.value = dValue.value
			pos := dValue.pos
			err := f.setString(dValue.str)
			if err != nil {
				appendWarning(f, pos, err)
			}
		}
	}

	for _, f := range cp.deferredValueFields {
		err := fmt.Errorf("unresolved deferredValueField: %s", f.String())
		dValue, _ := f.value.(deferredValue)
		pos := dValue.pos
		appendWarning(f, pos, err)
	}

	return warning
}

func RadioExists() error {
	dfu, err := dfu.New(nil)
	if err != nil {
		return err
	}
	dfu.Close()

	return nil
}

func (cp *Codeplug) ReadRadio(progress func(cur int) error) error {
	cpi := cp.codeplugInfo
	binBytes := cp.bytes[cpi.BinOffset : cpi.BinOffset+cpi.BinSize]

	dfu, err := dfu.New(func(cur int) error {
		return progress(cur)
	})
	if err != nil {
		return err
	}
	defer dfu.Close()

	bytes := make([]byte, len(binBytes))
	err = dfu.ReadCodeplug(bytes)
	if err != nil {
		return err
	}

	copy(binBytes, bytes)

	cp.Revert()

	cp.SetChanged()

	return nil
}

func (cp *Codeplug) WriteRadio(progress func(cur int) error) error {
	savedTime, err := cp.getLastProgrammedTime()
	if err != nil {
		return err
	}
	cp.setLastProgrammedTime(time.Now())

	savedBytes := make([]byte, len(cp.bytes))
	copy(savedBytes, cp.bytes)

	cp.store()

	cpi := cp.codeplugInfo
	binBytes := cp.bytes[cpi.BinOffset : cpi.BinOffset+cpi.BinSize]

	cp.bytes = savedBytes
	cp.setLastProgrammedTime(savedTime)

	dfu, err := dfu.New(func(cur int) error {
		return progress(cur)
	})
	if err != nil {
		return err
	}
	defer dfu.Close()

	err = dfu.WriteCodeplug(binBytes)
	if err != nil {
		return err
	}

	return nil
}
