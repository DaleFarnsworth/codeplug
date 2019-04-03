// Copyright 2017-2019 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of Radio.
//
// Radio is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU Lesser General Public
// License as published by the Free Software Foundation.
//
// Radio is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Radio.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dalefarnsworth/codeplug/codeplug"
	"github.com/dalefarnsworth/codeplug/debug"
	"github.com/dalefarnsworth/codeplug/dfu"
	"github.com/dalefarnsworth/codeplug/userdb"
)

func errorf(s string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, s, v...)
}

func usage() {
	errorf("Usage %s <subCommand> args\n", os.Args[0])
	errorf("subCommands:\n")
	errorf("\treadCodeplug -model <model> -freq <freqRange> <codeplugFile>\n")
	errorf("\twriteCodeplug <codeplugFile>\n")
	errorf("\twriteFirmware <firmwareFile>\n")
	errorf("\treadMD380Users <usersFile>\n")
	errorf("\twriteMD380Users <usersFile>\n")
	errorf("\twriteMD2017Users <usersFile>\n")
	errorf("\twriteUV380Users <usersFile>\n")
	errorf("\treadSPIFlash <filename>\n")
	errorf("\tgetUsers <usersFile>\n")
	errorf("\tgetInputUsers <usersFile>\n")
	errorf("\tcodeplugToText <codeplugFile> <textFile>\n")
	errorf("\ttextToCodeplug <textFile> <codeplugFile>\n")
	errorf("\tcodeplugToJSON <codeplugFile> <jsonFile>\n")
	errorf("\tjsonToCodeplug <jsonFile> <codeplugFile>\n")
	errorf("\tcodeplugToXLSX <codeplugFile> <xlsxFile>\n")
	errorf("\txlsxToCodeplug <xlsxFile> <codeplugFile>\n")
	errorf("\tversion\n")
	errorf("Use '%s <subCommand> -h' for subCommand help\n", os.Args[0])
	os.Exit(1)
}

func allTypesFrequencyRanges() (types []string, freqRanges map[string][]string) {
	freqRanges = codeplug.AllFrequencyRanges()
	types = make([]string, 0, len(freqRanges))

	for typ := range freqRanges {
		types = append(types, typ)
	}

	sort.Strings(types)

	return types, freqRanges
}

func loadCodeplug(fType codeplug.FileType, filename string) (*codeplug.Codeplug, error) {
	cp, err := codeplug.NewCodeplug(fType, filename)
	if err != nil {
		return nil, err
	}

	types, freqs := cp.TypesFrequencyRanges()
	if len(types) == 0 {
		return nil, errors.New("unknown model in codeplug")
	}

	typ := types[0]

	if len(freqs[typ]) == 0 {
		return nil, errors.New("unknown frequency range in codeplug")
	}

	freqRange := freqs[typ][0]

	err = cp.Load(typ, freqRange)
	if err != nil {
		return nil, err
	}

	return cp, nil
}

func progressCallback(aPrefixes []string) func(cur int) error {
	var prefixes []string
	if aPrefixes != nil {
		prefixes = aPrefixes
	}
	prefixIndex := 0
	prefix := prefixes[prefixIndex]
	maxProgress := userdb.MaxProgress
	return func(cur int) error {
		if cur == 0 {
			if prefixIndex != 0 {
				fmt.Println()
			}
			prefix = prefixes[prefixIndex]
			prefixIndex++
		}
		percent := cur * 100 / maxProgress
		fmt.Printf("%s... %3d%%\r", prefix, percent)
		return nil
	}
}

func readCodeplug() error {
	var typ string
	var freq string

	flags := flag.NewFlagSet("writeCodeplug", flag.ExitOnError)
	flags.StringVar(&typ, "model", "", "<model name>")
	flags.StringVar(&freq, "freq", "", "<frequency range>")

	flags.Usage = func() {
		errorf("Usage: %s %s -model <modelName> -freq <freqRange> codePlugFilename\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		errorf("modelName must be chosen from the following list,\n")
		errorf("and freqRange must be one of its associated values.\n")
		types, freqs := allTypesFrequencyRanges()
		for _, typ := range types {
			errorf("\t%s\n", typ)
			for _, freq := range freqs[typ] {
				errorf("\t\t%s\n", "\""+freq+"\"")
			}
		}
		os.Exit(1)
	}

	typeFreqs := codeplug.AllFrequencyRanges()

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	if typeFreqs[typ] == nil {
		errorf("bad modelName\n\n")
		flags.Usage()
	}
	freqMap := make(map[string]bool)
	for _, freq := range typeFreqs[typ] {
		freqMap[freq] = true
	}
	if !freqMap[freq] {
		errorf("bad freqRange\n\n")
		flags.Usage()
	}
	filename := args[0]

	cp, err := codeplug.NewCodeplug(codeplug.FileTypeNew, "")
	if err != nil {
		return err
	}

	err = cp.Load(typ, freq)
	if err != nil {
		return err
	}

	prefixes := []string{
		"Preparing to read codeplug",
		"Reading codeplug from radio.",
	}

	err = cp.ReadRadio(progressCallback(prefixes))
	if err != nil {
		return err
	}

	return cp.SaveAs(filename)
}

func writeCodeplug() error {
	flags := flag.NewFlagSet("writeCodeplug", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <codeplugFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	filename := args[0]

	cp, err := loadCodeplug(codeplug.FileTypeNone, filename)
	if err != nil {
		return err
	}

	prefixes := []string{
		"Preparing to write codeplug",
		"Writing codeplug to radio.",
	}

	return cp.WriteRadio(progressCallback(prefixes))
}

func readSPIFlash() (err error) {
	flags := flag.NewFlagSet("readSPIFlash", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <filename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	filename := args[0]

	prefixes := []string{
		"Preparing to read flash",
		"Reading flash",
	}

	dfu, err := dfu.New(progressCallback(prefixes))
	if err != nil {
		return err
	}
	defer dfu.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	return dfu.ReadSPIFlash(file)
}

func usersFilename() string {
	flags := flag.NewFlagSet("writeUsers", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <usersFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}

	return args[0]
}

func readMD380Users() (err error) {
	filename := usersFilename()

	prefixes := []string{
		"Preparing to read users",
		fmt.Sprintf("Reading users to %s", filename),
	}

	dfu, err := dfu.New(progressCallback(prefixes))
	if err != nil {
		return err
	}
	defer dfu.Close()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("os.Create: %s", err.Error())
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	return dfu.ReadUsers(file)
}

func writeMD380Users() error {
	filename := usersFilename()

	prefixes := []string{
		"Erasing flash memory",
		"Writing users",
	}

	dfu, err := dfu.New(progressCallback(prefixes))
	if err != nil {
		return err
	}
	defer dfu.Close()

	return dfu.WriteUsers(filename)
}

func writeMD2017Users() error {
	filename := usersFilename()

	prefixes := []string{
		"Erasing flash memory",
		"Writing users",
	}

	df, err := dfu.New(progressCallback(prefixes))
	if err != nil {
		return err
	}
	defer df.Close()

	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		return err
	}

	users := dfu.ParseUV380Users(file)
	return df.WriteUV380Users(users)
}

func writeUV380Users() error {
	filename := usersFilename()

	prefixes := []string{
		"Erasing flash memory",
		"Writing users",
	}

	df, err := dfu.New(progressCallback(prefixes))
	if err != nil {
		return err
	}
	defer df.Close()

	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		return err
	}

	users := dfu.ParseUV380Users(file)

	return df.WriteUV380Users(users)
}

func getUsers() error {
	flags := flag.NewFlagSet("getUsers", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <usersFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	filename := args[0]

	prefixes := []string{
		"Retrieving Users file",
	}

	db := userdb.New()
	return db.WriteMD380ToolsFile(filename, progressCallback(prefixes))
}

func getInputUsers() error {
	flags := flag.NewFlagSet("getUsers", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <usersFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	filename := args[0]

	prefixes := []string{
		"Retrieving Users file",
	}

	db := userdb.Input()
	return db.WriteMD380ToolsFile(filename, progressCallback(prefixes))
}

func writeFirmware() error {
	flags := flag.NewFlagSet("writeFirmware", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <firmwareFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	filename := args[0]

	prefixes := []string{
		"Erasing flash memory",
		"Writing firmware",
	}

	dfu, err := dfu.New(progressCallback(prefixes))
	if err != nil {
		return err
	}
	defer dfu.Close()

	file, err := os.Open(filename)
	if err != nil {
		l.Fatalf("writeFirmware: %s", err.Error())
	}

	defer file.Close()

	return dfu.WriteFirmware(file)
}

func textToCodeplug() error {
	flags := flag.NewFlagSet("textToCodeplug", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <textFilename> <codeplugFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 2 {
		flags.Usage()
	}
	textFilename := args[0]
	codeplugFilename := args[1]

	cp, err := loadCodeplug(codeplug.FileTypeText, textFilename)
	if err != nil {
		return err
	}

	return cp.SaveAs(codeplugFilename)
}

func codeplugToText() error {
	flags := flag.NewFlagSet("codeplugToText", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <codeplugFilename> <textFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 2 {
		flags.Usage()
	}
	codeplugFilename := args[0]
	textFilename := args[1]

	cp, err := loadCodeplug(codeplug.FileTypeNone, codeplugFilename)
	if err != nil {
		return err
	}

	return cp.ExportText(textFilename)
}

func jsonToCodeplug() error {
	flags := flag.NewFlagSet("jsonToCodeplug", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <jsonFilename> <codeplugFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 2 {
		flags.Usage()
	}
	jsonFilename := args[0]
	codeplugFilename := args[1]

	cp, err := loadCodeplug(codeplug.FileTypeJSON, jsonFilename)
	if err != nil {
		return err
	}

	return cp.SaveAs(codeplugFilename)
}

func codeplugToJSON() error {
	flags := flag.NewFlagSet("codeplugToJSON", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <codeplugFilename> <jsonFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 2 {
		flags.Usage()
	}
	codeplugFilename := args[0]
	jsonFilename := args[1]

	cp, err := loadCodeplug(codeplug.FileTypeNone, codeplugFilename)
	if err != nil {
		return err
	}

	return cp.ExportJSON(jsonFilename)
}

func xlsxToCodeplug() error {
	flags := flag.NewFlagSet("xlsxToCodeplug", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <xlsxFilename> <codeplugFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 2 {
		flags.Usage()
	}
	xlsxFilename := args[0]
	codeplugFilename := args[1]

	cp, err := loadCodeplug(codeplug.FileTypeXLSX, xlsxFilename)
	if err != nil {
		return err
	}

	return cp.SaveAs(codeplugFilename)
}

func codeplugToXLSX() error {
	flags := flag.NewFlagSet("codeplugToXLSX", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <codeplugFilename> <xlsxFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 2 {
		flags.Usage()
	}
	codeplugFilename := args[0]
	xlsxFilename := args[1]

	cp, err := loadCodeplug(codeplug.FileTypeNone, codeplugFilename)
	if err != nil {
		return err
	}

	return cp.ExportXLSX(xlsxFilename)
}

func printVersion() error {
	flags := flag.NewFlagSet("version", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:])
	args := flags.Args()
	if len(args) != 0 {
		flags.Usage()
	}

	fmt.Printf("%s\n", version)
	return nil
}

func main() {
	log.SetPrefix(filepath.Base(os.Args[0]) + ": ")
	log.SetFlags(log.Lshortfile)

	if len(os.Args) < 2 {
		usage()
	}

	subCommandName := strings.ToLower(os.Args[1])

	subCommands := map[string]func() error{
		"readcodeplug":     readCodeplug,
		"writecodeplug":    writeCodeplug,
		"readspiflash":     readSPIFlash,
		"readmd380users":   readMD380Users,
		"writemd380users":  writeMD380Users,
		"writemd2017users": writeMD2017Users,
		"writeuv380users":  writeUV380Users,
		"getusers":         getUsers,
		"getinputusers":    getInputUsers,
		"writefirmware":    writeFirmware,
		"texttocodeplug":   textToCodeplug,
		"codeplugtotext":   codeplugToText,
		"jsontocodeplug":   jsonToCodeplug,
		"codeplugtojson":   codeplugToJSON,
		"xlsxtocodeplug":   xlsxToCodeplug,
		"codeplugtoxlsx":   codeplugToXLSX,
		"version":          printVersion,
	}

	subCommand := subCommands[subCommandName]
	if subCommand == nil {
		usage()
	}

	err := subCommand()
	if err != nil {
		errorf("%s\n", err.Error())
		os.Exit(1)
	}
}
