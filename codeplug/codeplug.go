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
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// FileType tells whether the codeplug is an rdt file or a bin file.
type FileType int

const (
	FileTypeNone FileType = iota
	FileTypeRdt
	FileTypeBin
	FileTypeNew
	FileTypeText
)

// A Codeplug represents a codeplug file.
type Codeplug struct {
	filename      string
	textFilename  string
	fileType      FileType
	rdtSize       int
	fileSize      int
	fileOffset    int
	id            string
	bytes         []byte
	hash          [sha256.Size]byte
	rDesc         map[RecordType]*rDesc
	changed       bool
	lowFrequency  float64
	highFrequency float64
	connectChange func(*Change)
	changeList    []*Change
	changeIndex   int
	codeplugInfo  *CodeplugInfo
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

	case FileTypeText:
		cp.textFilename = filename
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
					log.Fatal(err.Error())
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
		log.Fatal("unknown file type")
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

func (cp *Codeplug) Load(model string, frequencyRange string, ignoreWarning bool) error {
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
	case FileTypeNew, FileTypeText:
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

	case FileTypeRdt, FileTypeBin:
		err := cp.read(cp.filename)
		if err != nil {
			return err
		}
	}

	err := cp.Revert(ignoreWarning)
	if err != nil {
		return err
	}

	switch cp.fileType {
	case FileTypeText:
		for _, rType := range cp.RecordTypes() {
			if cp.MaxRecords(rType) == 1 {
				continue
			}
			records := cp.records(rType)
			for i := len(records) - 1; i >= 0; i-- {
				cp.RemoveRecord(records[i])
			}
		}

		err := cp.importFrom(cp.textFilename, ignoreWarning)
		if _, warning := err.(Warning); warning && ignoreWarning {
			err = nil
		}

		if err != nil {
			return err
		}
	}

	codeplugs = append(codeplugs, cp)

	return nil
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

// ModelsFrequencyRanges returns the potential codeplug model and
// frequencyRange
func (cp *Codeplug) ModelsFrequencyRanges() (models []string, frequencyRanges map[string][]string) {
	models = make([]string, 0)
	frequencyRanges = make(map[string][]string)
	var model string
	var frequencyRange string

	switch cp.fileType {
	case FileTypeRdt:

	case FileTypeText:
		model, frequencyRange = parseModelFrequencyRange(cp.textFilename)
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
	fDescs := cp.rDesc[RecordType("BasicInformation")].records[0].fDesc
	return (*fDescs)[FieldType("Model")].fields[0].String()
}

func (cp *Codeplug) FrequencyRange() string {
	fDescs := cp.rDesc[RecordType("BasicInformation")].records[0].fDesc
	return (*fDescs)[FieldType("FrequencyRange")].fields[0].String()
}

func (cp *Codeplug) frequencyRanges() []string {
	for _, rInfo := range cp.codeplugInfo.RecordInfos {
		if rInfo.rType == "BasicInformation" {
			for _, fInfo := range rInfo.fieldInfos {
				if fInfo.fType == "FrequencyRange" {
					return *fInfo.strings
				}
			}
		}
	}
	return nil
}

func (cp *Codeplug) newFilenames() []string {
	for _, rInfo := range cp.codeplugInfo.RecordInfos {
		if rInfo.rType == "BasicInformation" {
			for _, fInfo := range rInfo.fieldInfos {
				if fInfo.fType == "NewFilename" {
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
		log.Fatal(err)
	}

	tarfile := tar.NewReader(archive)

	var bytes []byte
	for {
		hdr, err := tarfile.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		if hdr.Name != filename {
			continue
		}
		bytes, err = ioutil.ReadAll(tarfile)
		if err != nil {
			log.Fatal(err)
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
func (cp *Codeplug) Revert(ignoreWarning bool) error {
	cp.clearCachedListNames()

	cp.load()

	if err := cp.valid(); err != nil && !ignoreWarning {
		return err
	}

	cp.changed = false
	cp.hash = sha256.Sum256(cp.bytes)

	cp.changeList = []*Change{&Change{}}
	cp.changeIndex = 0

	return nil
}

// Save stores the state of the Codeplug into its file
// An error may be returned if the codeplug state is invalid.
func (cp *Codeplug) Save(ignoreWarning bool) error {
	return cp.SaveAs(cp.filename, ignoreWarning)
}

// SaveAs saves the state of the Codeplug into a named file.
// An error will be returned if the codeplug state is invalid.
// The named file becomes the current file associated with the codeplug.
func (cp *Codeplug) SaveAs(filename string, ignoreWarning bool) error {
	err := cp.SaveToFile(filename, ignoreWarning)
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
func (cp *Codeplug) SaveToFile(filename string, ignoreWarning bool) error {
	if err := cp.valid(); err != nil {
		return Warning{err}
	}

	cp.setLastProgrammedTime(time.Now())

	cp.store()

	dir, base := filepath.Split(filename)
	tmpFile, err := ioutil.TempFile(dir, base)
	if err != nil {
		return err
	}

	if err = cp.write(tmpFile); err != nil {
		os.Remove(tmpFile.Name())
		return err
	}

	if err := os.Rename(tmpFile.Name(), filename); err != nil {
		return err
	}

	return err
}

func (cp *Codeplug) setLastProgrammedTime(t time.Time) {
	r := cp.rDesc[RecordType("BasicInformation")].records[0]
	f := r.Field(FieldType("LastProgrammedTime"))
	dprint(f)
	f.setString(t.Format("02-Jan-06 15:04:05"))
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

// fecord returns the first record of a codeplug's given RecordType.
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
	strs := make([]string, 0, len(cp.rDesc))

	for rType := range cp.rDesc {
		strs = append(strs, string(rType))
	}
	sort.Strings(strs)

	rTypes := make([]RecordType, len(strs))
	for i, str := range strs {
		rTypes[i] = RecordType(str)
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
		log.Fatal("removeRecord: bad record")
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
	for _, ri := range cp.codeplugInfo.RecordInfos {
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
func (cp *Codeplug) valid() error {
	errStr := ""
	for _, rType := range cp.RecordTypes() {
		for _, r := range cp.records(rType) {
			if err := r.valid(); err != nil {
				errStr += err.Error()
			}
		}
	}

	for _, f := range deferredValidFields {
		if err := f.valid(); err != nil {
			errStr += fmt.Sprintf("%s %s\n", f.FullTypeName(), err.Error())
		}
	}

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

// write writes the codeplug's byte slice into the given file.
func (cp *Codeplug) write(file *os.File) (err error) {
	defer func() {
		err = file.Close()
		return
	}()

	cpi := cp.codeplugInfo
	fileSize := cpi.RdtSize
	fileOffset := 0

	bytes := cp.bytes[fileOffset : fileOffset+fileSize]
	bytesWritten, err := file.Write(bytes)
	if err != nil {
		return err
	}

	if bytesWritten != fileSize {
		return fmt.Errorf("write to %s failed", cp.filename)
	}

	return nil
}

// frequencyValid returns nil if the given frequency is valid for the
// codeplug.
func (cp *Codeplug) frequencyValid(freq float64) error {
	if cp.lowFrequency == 0 {
		fDescs := cp.record(RecordType("BasicInformation")).fDesc
		s := (*fDescs)[FtBiLowFrequency].fields[0].String()
		cp.lowFrequency, _ = strconv.ParseFloat(s, 64)
		s = (*fDescs)[FtBiHighFrequency].fields[0].String()
		cp.highFrequency, _ = strconv.ParseFloat(s, 64)
	}

	if freq >= cp.lowFrequency && freq <= cp.highFrequency {
		return nil
	}

	return fmt.Errorf("frequency out of range %+v", freq)
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

func PrintRecord(w io.Writer, r *Record) {
	rType := r.Type()
	ind := ""
	if r.max > 1 {
		ind = fmt.Sprintf("[%d]", r.rIndex+1)
	}
	fmt.Fprintf(w, "%s%s:\n", string(rType), ind)

	for _, fType := range r.FieldTypes() {
		if string(rType) == "BasicInformation" {
			switch string(fType) {
			case "CpsVersion", "NewFilename":
				continue
			case "LowFrequency", "HighFrequency":
				continue
			}
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

var nameToRt map[string]RecordType
var nameToFt map[RecordType]map[string]FieldType

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

func parseModelFrequencyRange(filename string) (model string, frequencyRange string) {
	file, err := os.Open(filename)
	if err != nil {
		return model, frequencyRange
	}
	defer file.Close()

	rdr := NewReader(file)

	rdr.ReadWhile(unicode.IsSpace)

	if rdr.pos.column != 0 {
		return model, frequencyRange
	}

newRecord:
	for {
		rName, _, err := parseName(rdr)
		if err != nil {
			break
		}
		for {
			if rdr.pos.column == 0 {
				break
			}
			fName, _, err := parseName(rdr)
			if err != nil {
				break newRecord
			}

			value, err := parseValue(rdr)
			if err != nil {
				break newRecord
			}
			if rName == "BasicInformation" {
				switch fName {
				case "Model":
					model = value
				case "FrequencyRange":
					frequencyRange = value
				default:
					continue
				}
				if model != "" && frequencyRange != "" {
					return model, frequencyRange
				}
			}
		}
	}

	return model, frequencyRange
}

func (cp *Codeplug) ParseRecords(iRdr io.Reader) ([]*Record, error) {
	var warning error
	var err error
	var pos position

	rdr := NewReader(iRdr)
	records := []*Record{}

	if len(nameToRt) == 0 {
		nameToRt = make(map[string]RecordType)
		for _, rType := range cp.RecordTypes() {
			nameToRt[string(rType)] = rType
		}

		nameToFt = make(map[RecordType]map[string]FieldType)
		for _, rType := range cp.RecordTypes() {
			m := make(map[string]FieldType)
			for _, fi := range cp.rDesc[rType].fieldInfos {
				fType := fi.fType
				name := string(fType)
				m[name] = fType
			}
			nameToFt[rType] = m
		}
	}

	rdr.ReadWhile(unicode.IsSpace)

parseRecord:
	for {
		var name string
		var index int
		pos = rdr.pos
		name, index, err = parseName(rdr)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if len(name) == 0 {
			err = fmt.Errorf("no record name")
			break
		}
		var r *Record
		r, err = cp.nameToRecord(name, index)
		if err != nil {
			break
		}
		for {
			if rdr.pos.column == 0 {
				break
			}
			pos = rdr.pos
			name, index, err = parseName(rdr)
			if err != nil {
				break parseRecord
			}
			if len(name) == 0 {
				err = fmt.Errorf("no field name")
				break parseRecord
			}
			pos = rdr.pos
			fType, ok := nameToFt[r.rType][name]
			if !ok {
				err = fmt.Errorf("bad field name")
				warning = appendWarningMsgs(warning, pos, err)
			}

			var str string
			pos = rdr.pos
			str, err = parseValue(rdr)
			if err != nil {
				err = fmt.Errorf("bad value: %s: %s: %s", name, str, err.Error())
				break parseRecord
			}
			var f *Field
			f, err = r.NewFieldWithValue(fType, index, str)
			if err != nil && ok {
				warning = appendWarningMsgs(warning, pos, err)
			}
			dValue, ok := f.value.(deferredValue)
			if ok {
				dValue.str = str
				dValue.pos = pos
				f.value = dValue
			}
			err = r.addField(f)
			if err != nil {
				break parseRecord
			}
		}

		records = append(records, r)
	}

	if err != nil {
		pErr, ok := err.(PositionError)
		if !ok {
			pErr = PositionError{
				position: &pos,
				error:    err,
			}
		}
		return records, pErr
	}

	if warning == nil {
		return records, nil
	}

	return records, warning
}

func appendWarningMsgs(oldError error, pos position, newError error) error {
	var oldMsg string
	if oldError != nil {
		oldMsg = oldError.Error()
	}
	err := fmt.Errorf("%s%d:%d: %s\n", oldMsg, pos.line+1,
		pos.column+1, newError.Error())
	return Warning{err}
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

	if err != nil {
		if err == io.EOF {
			return name, index, err
		}

		err = fmt.Errorf("bad %s name", nType)
		goto returnError
	}
	if len(name) == 0 || !unicode.IsLetter([]rune(name)[0]) {
		err = fmt.Errorf("bad %s name", nType)
		goto returnError
	}

	pos = rdr.pos
	r, _, err = rdr.ReadRune()
	if err != nil {
		err = fmt.Errorf("bad %s name", nType)
		goto returnError
	}
	switch r {
	case ':':
	case '[':
		pos = rdr.pos
		index, err = rdr.ReadInt()
		if err != nil {
			err = fmt.Errorf("bad %s index", nType)
			goto returnError
		}
		index--
		pos = rdr.pos
		r, _, err = rdr.ReadRune()
		if r != ']' {
			err = fmt.Errorf("bad %s index", nType)
			goto returnError
		}
		pos = rdr.pos
		r, _, err = rdr.ReadRune()
		if r != ':' {
			err = fmt.Errorf("bad %s name", nType)
			goto returnError
		}
	default:
		err = fmt.Errorf("bad %s index", nType)
		goto returnError
	}
	rdr.ReadWhile(unicode.IsSpace)

	return name, index, nil

returnError:
	pErr := PositionError{
		position: &pos,
		error:    err,
	}
	return name, 0, pErr
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

func (cp *Codeplug) nameToRecord(name string, index int) (*Record, error) {
	rType, ok := nameToRt[name]
	if !ok {
		return nil, fmt.Errorf("unknown record type: %s", name)
	}

	found := false
	for rt := range cp.rDesc {
		if rType == rt {
			found = true
			break
		}
	}
	if found {
		return cp.newRecord(rType, index), nil
	}

	name = string(rType)
	return nil, fmt.Errorf("codeplug has no record: %s", name)
}

func (cp *Codeplug) ExportTo(filename string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
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

func (cp *Codeplug) clearCachedListNames() {
	for _, rd := range cp.rDesc {
		rd.cachedListNames = nil
	}
}

func (cp *Codeplug) importFrom(filename string, ignoreWarning bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	records, err := cp.ParseRecords(file)
	if err != nil && !ignoreWarning {
		return err
	}

	for _, r := range records {
		r.rIndex = len(cp.records(r.rType))
		cp.InsertRecord(r)
	}

	derr := updateDeferredFields(records)
	if derr != nil {
		err = Warning{fmt.Errorf("%s%s", err.Error(), derr.Error())}
	}

	if err != nil && !ignoreWarning {
		return err
	}

	for _, rd := range cp.rDesc {
		if len(rd.records) == 0 {
			rtName := string(rd.rType)
			err := Warning{fmt.Errorf("no %s records found", rtName)}
			return err
		}
	}
	cp.store()
	cp.changeList = []*Change{&Change{}}
	cp.changeIndex = 0
	cp.changed = true

	return nil
}
