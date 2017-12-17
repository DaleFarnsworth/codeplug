// Copyright 2017-2018 Dale Farnsworth. All rights reserved.

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
	"github.com/dalefarnsworth/codeplug/dfu"
	"github.com/dalefarnsworth/codeplug/userdb"
)

func errorf(s string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, s, v...)
}

func usage() {
	errorf("Usage %s <subCommand> args\n", os.Args[0])
	errorf("subCommands:\n")
	errorf("\treadCodeplug -model <model> -freq <freqRange> <codeplugFilename>\n")
	errorf("\twriteCodeplug <codeplugFilename>\n")
	errorf("\twriteFirmware <firmwareFilename>\n")
	errorf("\treadUsers <usersFilename>\n")
	errorf("\twriteUsers <usersFilename>\n")
	errorf("\treadSPIFlash <filename>\n")
	errorf("\tgetUsers <usersFilename>\n")
	errorf("\tcodeplugToText <codeplugFilename> <textFilename>\n")
	errorf("\ttextToCodeplug <textFilename> <codeplugFilename>\n")
	errorf("\tcodeplugToJSON <codeplugFilename> <jsonFilename>\n")
	errorf("\tjsonToCodeplug <jsonFilename> <codeplugFilename>\n")
	errorf("\tcodeplugToXLSX <codeplugFilename> <xlsxFilename>\n")
	errorf("\txlsxToCodeplug <xlsxFilename> <codeplugFilename>\n")
	errorf("\tversion\n")
	errorf("Use '%s <subCommand> -h' for subCommand help\n", os.Args[0])
	os.Exit(1)
}

func allModelsFrequencyRanges() (models []string, freqRanges map[string][]string) {
	freqRanges = codeplug.AllFrequencyRanges()
	models = make([]string, 0, len(freqRanges))

	for model := range freqRanges {
		models = append(models, model)
	}

	sort.Slice(models, func(i, j int) bool {
		return models[i] < models[j]
	})

	return models, freqRanges
}

func loadCodeplug(fType codeplug.FileType, filename string) (*codeplug.Codeplug, error) {
	cp, err := codeplug.NewCodeplug(fType, filename)
	if err != nil {
		return nil, err
	}

	models, freqs := cp.ModelsFrequencyRanges()
	if len(models) == 0 {
		return nil, errors.New("unknown model in codeplug")
	}

	model := models[0]

	if len(freqs[model]) == 0 {
		return nil, errors.New("unknown frequency range in codeplug")
	}

	freq := freqs[model][0]

	ignoreWarnings := true
	err = cp.Load(model, freq, ignoreWarnings)
	if err != nil {
		return nil, err
	}

	return cp, nil
}

func progressFunc(aPrefixes []string) func(cur int) bool {
	var prefixes []string
	if aPrefixes != nil {
		prefixes = aPrefixes
	}
	prefixIndex := 0
	prefix := prefixes[prefixIndex]
	maxProgress := userdb.MaxProgress
	return func(cur int) bool {
		if cur == 0 {
			if prefixIndex != 0 {
				fmt.Println()
			}
			prefix = prefixes[prefixIndex]
			prefixIndex++
		}
		fmt.Printf("%s... %3d%%\r", prefix, cur*100/maxProgress)
		return true
	}
}

func readCodeplug() error {
	var model string
	var freq string

	flags := flag.NewFlagSet("writeCodeplug", flag.ExitOnError)
	flags.StringVar(&model, "model", "", "<model name>")
	flags.StringVar(&freq, "freq", "", "<frequency range>")

	flags.Usage = func() {
		errorf("Usage: %s %s -model <modelName> -freq <freqRange> codePlugFilename\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		errorf("modelName must be chosen from the following list,\n")
		errorf("and freqRange must be one of its associated values.\n")
		models, freqs := allModelsFrequencyRanges()
		for _, model := range models {
			errorf("\t%s\n", model)
			for _, freq := range freqs[model] {
				errorf("\t\t%s\n", "\""+freq+"\"")
			}
		}
		os.Exit(1)
	}

	modelFreqs := codeplug.AllFrequencyRanges()

	flags.Parse(os.Args[2:len(os.Args)])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	if modelFreqs[model] == nil {
		errorf("bad modelName\n\n")
		flags.Usage()
	}
	freqMap := make(map[string]bool)
	for _, freq := range modelFreqs[model] {
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

	ignoreWarnings := true
	err = cp.Load(model, freq, ignoreWarnings)
	if err != nil {
		return err
	}

	prefixes := []string{
		"Preparing to read codeplug",
		"Reading codeplug from radio.",
	}

	err = cp.ReadRadio(progressFunc(prefixes))
	if err != nil {
		return err
	}

	return cp.SaveAs(filename, ignoreWarnings)
}

func writeCodeplug() error {
	flags := flag.NewFlagSet("writeCodeplug", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <codeplugFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:len(os.Args)])
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

	return cp.WriteRadio(progressFunc(prefixes))
}

func readSPIFlash() (err error) {
	flags := flag.NewFlagSet("readSPIFlash", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <filename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:len(os.Args)])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	filename := args[0]

	prefixes := []string{
		"Preparing to read flash",
		"Reading flash",
	}

	dfu, err := dfu.New(progressFunc(prefixes))
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

func readUsers() (err error) {
	flags := flag.NewFlagSet("readUsers", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <usersFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:len(os.Args)])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	filename := args[0]

	prefixes := []string{
		"Preparing to read users",
		fmt.Sprintf("Reading users to %s", filename),
	}

	dfu, err := dfu.New(progressFunc(prefixes))
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

func writeUsers() error {
	flags := flag.NewFlagSet("writeUsers", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <usersFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:len(os.Args)])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	filename := args[0]

	prefixes := []string{
		"Erasing flash memory",
		"Writing users",
	}

	dfu, err := dfu.New(progressFunc(prefixes))
	if err != nil {
		return err
	}
	defer dfu.Close()

	return dfu.WriteUsers(filename)
}

func getUsers() error {
	flags := flag.NewFlagSet("getUsers", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <usersFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:len(os.Args)])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	filename := args[0]

	prefixes := []string{
		"Retrieving Users file",
	}

	db := userdb.New()
	return db.WriteMD380ToolsFile(filename, progressFunc(prefixes))
}

func writeFirmware() error {
	flags := flag.NewFlagSet("writeFirmware", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <firmwareFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:len(os.Args)])
	args := flags.Args()
	if len(args) != 1 {
		flags.Usage()
	}
	filename := args[0]

	prefixes := []string{
		"Erasing flash memory",
		"Writing firmware",
	}

	dfu, err := dfu.New(progressFunc(prefixes))
	if err != nil {
		return err
	}
	defer dfu.Close()

	return dfu.WriteFirmware(filename)
}

func textToCodeplug() error {
	flags := flag.NewFlagSet("textToCodeplug", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <textFilename> <codeplugFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:len(os.Args)])
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

	ignoreWarnings := true
	return cp.SaveAs(codeplugFilename, ignoreWarnings)
}

func codeplugToText() error {
	flags := flag.NewFlagSet("codeplugToText", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <codeplugFilename> <textFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:len(os.Args)])
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

	flags.Parse(os.Args[2:len(os.Args)])
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

	ignoreWarnings := true
	return cp.SaveAs(codeplugFilename, ignoreWarnings)
}

func codeplugToJSON() error {
	flags := flag.NewFlagSet("codeplugToJSON", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <codeplugFilename> <jsonFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:len(os.Args)])
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

	flags.Parse(os.Args[2:len(os.Args)])
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

	ignoreWarnings := true
	return cp.SaveAs(codeplugFilename, ignoreWarnings)
}

func codeplugToXLSX() error {
	flags := flag.NewFlagSet("codeplugToXLSX", flag.ExitOnError)

	flags.Usage = func() {
		errorf("Usage: %s %s <codeplugFilename> <xlsxFilename>\n", os.Args[0], os.Args[1])
		flags.PrintDefaults()
		os.Exit(1)
	}

	flags.Parse(os.Args[2:len(os.Args)])
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

	flags.Parse(os.Args[2:len(os.Args)])
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
		"readcodeplug":   readCodeplug,
		"writecodeplug":  writeCodeplug,
		"readspiflash":   readSPIFlash,
		"readusers":      readUsers,
		"writeusers":     writeUsers,
		"getusers":       getUsers,
		"writefirmware":  writeFirmware,
		"texttocodeplug": textToCodeplug,
		"codeplugtotext": codeplugToText,
		"jsontocodeplug": jsonToCodeplug,
		"codeplugtojson": codeplugToJSON,
		"xlsxtocodeplug": xlsxToCodeplug,
		"codeplugtoxlsx": codeplugToXLSX,
		"version":        printVersion,
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
