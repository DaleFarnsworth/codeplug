// Copyright 2017-2018 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This file is part of UserDB.
//
// UserDB is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU Lesser General Public
// License as published by the Free Software Foundation.
//
// UserDB is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with UserDB.  If not, see <http://www.gnu.org/licenses/>.

package userdb

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var specialUsersURL = "http://registry.dstar.su/api/node.php"
var fixedUsersURL = "https://raw.githubusercontent.com/travisgoodspeed/md380tools/master/db/fixed.csv"
var radioidUsersURL = "https://www.radioid.net/static/users_quoted.csv"
var hamdigitalUsersURL = "https://ham-digital.org/status/users_quoted.csv"
var reflectorUsersURL = "http://registry.dstar.su/reflector.db"
var cachedSpecialUsersURL = "https://www.farnsworth.org/dale/md380tools/userdb/cachedSpecialUsers.bin"
var cachedAmendedUsersURL = "https://www.farnsworth.org/dale/md380tools/userdb/cachedAmendedUsers.bin"

var transportTimeout = 20
var clientTimeout = 300

var tr = &http.Transport{
	TLSHandshakeTimeout:   time.Duration(transportTimeout) * time.Second,
	ResponseHeaderTimeout: time.Duration(transportTimeout) * time.Second,
}

var client = &http.Client{
	Transport: tr,
	Timeout:   time.Duration(clientTimeout) * time.Second,
}

type Options uint32

const (
	_ Options = 1 << iota
	AbbrevCountries
	AbbrevDirections
	AbbrevStates
	CheckTitleCase
	FixRomanNumerals
	FixStateCountries
	MiscChanges
	RemoveCallFromNick
	RemoveDupSurnames
	RemoveMatchingNick
	RemoveRepeats
	TitleCase
)

type User struct {
	ID       string
	Callsign string
	Name     string
	City     string
	State    string
	Nick     string
	Country  string
}

type UsersDB struct {
	filename          string
	options           Options
	uris              []string
	verbatimURIs      []string
	printFunc         func(*User) string
	progressCallback  func(progressCounter int) bool
	progressFunc      func() error
	progressIncrement int
	progressCounter   int
}

var defaultOptions = []Options{
	FixRomanNumerals,
	FixStateCountries,
	MiscChanges,
	RemoveCallFromNick,
	RemoveDupSurnames,
	RemoveMatchingNick,
	RemoveRepeats,
	TitleCase,
}

func New() *UsersDB {
	db := &UsersDB{
		progressFunc: func() error { return nil },
	}

	db.SetOptions(defaultOptions...)

	return db
}

func (db *UsersDB) SetOptions(options ...Options) {
	for _, opt := range options {
		db.options = Options(uint32(db.options) | uint32(opt))
	}
}

func (db *UsersDB) ClearOptions(options ...Options) {
	for _, opt := range options {
		db.options = Options(uint32(db.options) & ^uint32(opt))
	}
}

func (db *UsersDB) SetURIs(uris ...string) {
	db.uris = append(db.uris, uris...)
}

func (db *UsersDB) SetVerbatimURIs(uris ...string) {
	db.verbatimURIs = append(db.verbatimURIs, uris...)
}

func (db *UsersDB) setMaxProgressCount(max int) {
	db.progressFunc = func() error { return nil }
	if db.progressCallback != nil {
		db.progressIncrement = MaxProgress / max
		db.progressCounter = 0
		db.progressFunc = func() error {
			db.progressCounter += db.progressIncrement
			curProgress := db.progressCounter
			if curProgress > MaxProgress {
				curProgress = MaxProgress
			}

			if !db.progressCallback(db.progressCounter) {
				return errors.New("")
			}

			return nil
		}
		db.progressCallback(db.progressCounter)
	}
}

func (db *UsersDB) finalProgress() {
	//fmt.Fprintf(os.Stderr, "\nprogressMax %d\n", db.progressCounter/db.progressIncrement)
	if db.progressCallback != nil {
		db.progressCallback(MaxProgress)
	}
}

const (
	MinProgress = 0
	MaxProgress = 1000000
)

func (u *User) amend(options Options) {
	u.removeBlanksFromCallsigns()

	if options&RemoveDupSurnames != 0 {
		u.Name = removeDupSurnames(u.Name)
	}
	if options&RemoveRepeats != 0 {
		u.Name = removeRepeats(u.Name)
		u.City = removeRepeats(u.City)
		u.State = removeRepeats(u.State)
		u.Nick = removeRepeats(u.Nick)
		u.Country = removeRepeats(u.Country)
	}
	if options&TitleCase != 0 {
		u.Name = titleCase(u.Name)
		u.City = titleCase(u.City)
		u.State = titleCase(u.State)
		u.Country = titleCase(u.Country)
	}
	if options&RemoveMatchingNick != 0 {
		u.removeMatchingNicks()
	} else {
		u.addNicks()
	}
	if options&FixStateCountries != 0 {
		u.fixStateCountries()
	}
	if options&AbbrevCountries != 0 {
		u.Country = countryAbbreviations[u.Country]
	}
	if options&AbbrevStates != 0 {
		u.State = stateAbbreviations[u.State]
	}
	if options&AbbrevDirections != 0 {
		u.City = abbreviateDirections(u.City)
		u.State = abbreviateDirections(u.State)
		u.Nick = abbreviateDirections(u.Nick)
		u.Country = abbreviateDirections(u.Country)
	}
	if options&RemoveCallFromNick != 0 {
		u.Nick = removeSubstr(u.Nick, u.Callsign)
	}
	if options&MiscChanges != 0 {
		if strings.HasSuffix(u.City, " (B,") {
			length := len(u.City) - len(" (B,")
			u.City = u.City[:length]
		}
	}
	if options&FixRomanNumerals != 0 {
		u.Name = fixRomanNumerals(u.Name)
	}

	u.normalize()
}

func (u *User) normalize() {
	u.Callsign = normalizeString(u.Callsign)
	u.Name = normalizeString(u.Name)
	u.City = normalizeString(u.City)
	u.State = normalizeString(u.State)
	u.Nick = normalizeString(u.Nick)
	u.Country = normalizeString(u.Country)
}

func normalizeString(field string) string {
	field = asciify(field)
	field = strings.TrimSpace(field)
	field = strings.Replace(field, ",", ";", -1)

	for strings.Index(field, "  ") >= 0 {
		field = strings.Replace(field, "  ", " ", -1)
	}

	return field
}

func asciify(field string) string {
	runes := []rune(field)
	strs := make([]string, len(runes))
	for i, r := range runes {
		strs[i] = transliterations[r]
	}

	return strings.Join(strs, "")
}

func (u *User) removeBlanksFromCallsigns() {
	id64, err := strconv.ParseUint(u.ID, 10, 24)
	if err != nil {
		return
	}
	id := int(id64)
	if id < 1000000 {
		return
	}
	u.Callsign = strings.Replace(u.Callsign, " ", "", -1)
}

var stateAbbreviations = func() map[string]string {
	abbrevs := make(map[string]string)
	for _, stateAbbreviations := range stateAbbreviationsByCountry {
		for state, abbrev := range stateAbbreviations {
			abbrevs[state] = abbrev
		}
	}
	return abbrevs
}()

var titleCaseMap = func() map[string]string {
	m := make(map[string]string)
	for _, word := range titleCaseWords {
		m[word] = strings.Title(word)
	}
	return m
}()

func abbreviateDirections(field string) string {
	words := strings.Split(field, " ")
	words[0] = directionAbbreviations[words[0]]
	return strings.Join(words, " ")
}

func removeDupSurnames(field string) string {
	words := strings.Split(field, " ")
	length := len(words)
	if length < 3 || words[length-2] != words[length-1] {
		return field
	}

	return strings.Join(words[:length-1], " ")
}

func removeRepeats(field string) string {
	words := strings.Split(field, " ")
	if len(words) < 4 || len(words)%2 != 0 {
		return field
	}

	halfLen := len(words) / 2
	for i := 0; i < halfLen; i++ {
		if words[i] != words[i+halfLen] {
			return field
		}
	}
	return strings.Join(words[:halfLen], " ")
}

func titleCase(field string) string {
	words := strings.Split(field, " ")
	for i, word := range words {
		title := titleCaseMap[word]
		if title != "" {
			words[i] = title
		}
	}

	return strings.Join(words, " ")
}

func checkTitleCase(users []*User) {
	upperCaseMap := make(map[string]bool)
	for _, word := range upperCaseWords {
		upperCaseMap[word] = true
	}

	fmt.Println("new upper-case words:")
	for _, u := range users {
		fields := []string{
			u.Name,
			u.City,
			u.State,
			u.Nick,
			u.Country,
		}
		for _, field := range fields {
		nextWord:
			for _, word := range strings.Split(field, " ") {
				if len(word) < 2 {
					continue
				}

				for r := range word {
					if r < 'A' || r > 'Z' {
						continue nextWord
					}
				}

				if titleCaseMap[word] != "" {
					continue
				}

				if upperCaseMap[word] {
					continue
				}

				fmt.Println(word)
			}
		}
	}

	fmt.Println("end of new upper-case words")
}

func (u *User) removeMatchingNicks() {
	firstName := strings.SplitN(u.Name, " ", 2)[0]
	if u.Nick == firstName {
		u.Nick = ""
	}
}

func (u *User) addNicks() {
	firstName := strings.SplitN(u.Name, " ", 2)[0]
	if u.Nick == "" {
		u.Nick = firstName
	}
}

func removeSubstr(field string, subf string) string {
	index := strings.Index(strings.ToUpper(field), strings.ToUpper(subf))
	if index >= 0 {
		field = field[:index] + field[index+len(subf):]
	}

	return field
}

func fixRomanNumerals(field string) string {
	fLen := len(field)
	if fLen < 3 {
		return field
	}

	if strings.HasSuffix(field, "i") {
		if strings.HasSuffix(field, " Ii") {
			field = field[:fLen-1] + "I"
		} else if strings.HasSuffix(field, " Iii") {
			field = field[:fLen-2] + "II"
		}
	} else if strings.HasSuffix(field, " Iv") {
		field = field[:fLen-1] + "V"
	}

	return field
}

func (u *User) usCallsign() bool {
	runes := []rune(u.Callsign)
	if strings.ContainsRune("KNW", runes[0]) {
		return true
	}

	if runes[0] == 'A' && runes[1] >= 'A' && runes[1] <= 'L' {
		return true
	}

	return false
}

func (u *User) fixStateCountries() {
	for country, stateAbbrevs := range stateAbbreviationsByCountry {
		for state := range stateAbbrevs {
			if u.Country == state {
				if state == "Georgia" && !u.usCallsign() {
					continue
				}
				if u.State == "" {
					u.State = state
				}
				u.Country = country
			}
		}
	}
}

func getUrlBytes(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func getUrlLines(url string) ([]string, error) {
	bytes, err := getUrlBytes(url)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(bytes), "\n")

	return lines[:len(lines)-1], nil
}

func getRadioidUsers() ([]*User, error) {
	lines, err := getUrlLines(radioidUsersURL)
	if err != nil {
		errFmt := "error getting radioid users database: %s: %s"
		err = fmt.Errorf(errFmt, radioidUsersURL, err.Error())
		return nil, err
	}

	if len(lines) < 50000 {
		errFmt := "too few radioid users database entries: %s: %d"
		err = fmt.Errorf(errFmt, radioidUsersURL, len(lines))
		return nil, err
	}

	users := make([]*User, len(lines))
	for i, line := range lines {
		line = strings.Trim(line, `"`)
		fields := strings.Split(line, `","`)

		users[i] = &User{
			ID:       fields[0],
			Callsign: fields[1],
			Name:     fields[2],
			City:     fields[3],
			State:    fields[4],
			Country:  fields[5],
		}
	}
	return users, nil
}

func getHamdigitalUsers() ([]*User, error) {
	lines, err := getUrlLines(hamdigitalUsersURL)
	if err != nil {
		errFmt := "error getting hamdigital users database: %s: %s"
		err = fmt.Errorf(errFmt, hamdigitalUsersURL, err.Error())
		return nil, err
	}

	if len(lines) < 50000 {
		errFmt := "too few hamdigital users database entries: %s: %d"
		err = fmt.Errorf(errFmt, hamdigitalUsersURL, len(lines))
		return nil, err
	}

	users := make([]*User, len(lines))
	for i, line := range lines {
		line = strings.Trim(line, `"`)
		fields := strings.Split(line, `","`)

		users[i] = &User{
			ID:       fields[0],
			Callsign: fields[1],
			Name:     fields[2],
			City:     fields[3],
			State:    fields[4],
			Country:  fields[5],
		}
	}
	return users, nil
}

func getCachedSpecialUsers() ([]*User, error) {
	lines, err := getUrlLines(cachedSpecialUsersURL)
	if err != nil {
		return nil, err
	}

	users := make([]*User, len(lines))
	for i, line := range lines {
		fields := strings.Split(line, ",")
		users[i] = &User{
			ID:       fields[0],
			Callsign: fields[1],
			Name:     fields[2],
			City:     fields[3],
			State:    fields[4],
			Nick:     fields[5],
			Country:  fields[6],
		}
	}
	return users, nil
}

func linesToUsers(url string, lines []string) ([]*User, error) {
	users := make([]*User, 0, len(lines))
	errStrs := make([]string, 0)
	for i, line := range lines {
		fmtStr := ""
		fields := strings.Split(line, ",")
		id, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil || id > 16777215 {
			fmtStr = "%s%d invalid DMR ID value: %s"
			if err != nil {
				fmtStr = "%s:%d non-numeric DMR ID: %s"
			}
			err := fmt.Sprintf(fmtStr, url, i, line)
			errStrs = append(errStrs, err)
			continue
		}
		if len(fields) != 7 {
			fmtStr := "%s:%d too many fields: %s"
			if len(fields) < 7 {
				fields = append(fields, []string{
					"", "", "", "", "", "", "",
				}...)
				fmtStr = "%s:%d too few fields: %s"
			}
			err := fmt.Sprintf(fmtStr, url, i, line)
			errStrs = append(errStrs, err)
			fields = fields[:7]
		}
		user := &User{
			ID:       fields[0],
			Callsign: fields[1],
			Name:     fields[2],
			City:     fields[3],
			State:    fields[4],
			Nick:     fields[5],
			Country:  fields[6],
		}
		users = append(users, user)
	}

	err := errors.New(strings.Join(errStrs, "\n"))
	if len(errStrs) == 0 {
		err = nil
	}

	return users, err
}

func newFileUsersFuncs(path string) (func() ([]*User, error), error) {
	return func() ([]*User, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		lines := make([]string, 0)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		err = scanner.Err()
		if err != nil {
			return nil, err
		}
		return linesToUsers(path, lines)
	}, nil
}

func newUrlUsersFuncs(uri string) (func() ([]*User, error), error) {
	return func() ([]*User, error) {
		lines, err := getUrlLines(uri)
		if err != nil {
			return nil, err
		}

		return linesToUsers(uri, lines)
	}, nil
}

func getFixedUsers() ([]*User, error) {
	lines, err := getUrlLines(fixedUsersURL)
	if err != nil {
		errFmt := "getting fixed users: %s: %s"
		err = fmt.Errorf(errFmt, fixedUsersURL, err.Error())
		return nil, err
	}

	users := make([]*User, len(lines))
	for i, line := range lines {
		fields := strings.Split(line, ",")
		users[i] = &User{
			ID:       fields[0],
			Callsign: fields[1],
		}
	}
	return users, nil
}

type special struct {
	ID      string
	Country string
	Address string
}

func getSpecialURLs() ([]string, error) {
	bytes, err := getUrlBytes(specialUsersURL)
	if err != nil {
		return nil, err
	}

	var specials []special
	err = json.Unmarshal(bytes, &specials)

	var urls []string
	for _, s := range specials {
		url := "http://" + s.Address + "/md380tools/special_IDs.csv"
		urls = append(urls, url)
	}

	return urls, nil
}

func getSpecialUsers(url string) ([]*User, error) {
	lines, err := getUrlLines(url)
	if err != nil {
		errFmt := "getting special users: %s: %s"
		err = fmt.Errorf(errFmt, url, err.Error())
		return nil, nil // Ignore erros on special users
	}

	users := make([]*User, len(lines))
	for i, line := range lines {
		fields := strings.Split(line, ",")
		if len(fields) < 7 {
			continue
		}
		users[i] = &User{
			ID:       fields[0],
			Callsign: fields[1],
			Name:     fields[2],
			Country:  fields[6],
		}
	}
	return users, nil
}

func getReflectorUsers() ([]*User, error) {
	lines, err := getUrlLines(reflectorUsersURL)
	if err != nil {
		errFmt := "getting reflector users: %s: %s"
		err = fmt.Errorf(errFmt, reflectorUsersURL, err.Error())
		return nil, err
	}

	users := make([]*User, len(lines))
	for i, line := range lines[1:] {
		line := strings.Replace(line, "@", ",", 2)
		fields := strings.Split(line, ",")
		users[i] = &User{
			ID:       fields[0],
			Callsign: fields[1],
		}
	}
	return users, nil
}

func mergeAndSort(users []*User) ([]*User, error) {
	idMap := make(map[int]*User)
	for _, u := range users {
		if u == nil || u.ID == "" {
			continue
		}
		u.ID = strings.TrimPrefix(u.ID, "#")
		id, err := strconv.ParseUint(u.ID, 10, 24)
		if err != nil {
			return nil, err
		}
		existing := idMap[int(id)]
		if existing == nil {
			idMap[int(id)] = u
			continue
		}
		// non-empty fields in later entries replace fields in earlier
		if u.Callsign != "" {
			existing.Callsign = u.Callsign
		}
		if u.Name != "" {
			existing.Name = u.Name
		}
		if u.City != "" {
			existing.City = u.City
		}
		if u.State != "" {
			existing.State = u.State
		}
		if u.Country != "" {
			existing.Country = u.Country
		}
	}

	ids := make([]int, 0, len(idMap))
	for id := range idMap {
		ids = append(ids, id)
	}

	users = make([]*User, len(ids))
	sort.Ints(ids)
	for i, id := range ids {
		users[i] = idMap[id]
	}

	return users, nil
}

type result struct {
	index int
	users []*User
	err   error
}

func do(index int, f func() ([]*User, error), resultChan chan result) {
	var r result

	r.index = index
	r.users, r.err = f()
	resultChan <- r
}

func (db *UsersDB) Users() ([]*User, error) {
	getUsersFuncs := []func() ([]*User, error){
		getFixedUsers,
		getHamdigitalUsers,
		getRadioidUsers,
		getReflectorUsers,
	}

	var users []*User
	resultCount := len(getUsersFuncs)
	resultChan := make(chan result, resultCount)

	for i, f := range getUsersFuncs {
		go do(i, f, resultChan)
	}

	db.setMaxProgressCount(resultCount)

	results := make([]result, resultCount)
	for done := 0; done < resultCount; {
		select {
		case r := <-resultChan:
			if r.err != nil {
				return nil, r.err
			}
			results[r.index] = r
			done++
			err := db.progressFunc()
			if err != nil {
				return nil, err
			}
		}
	}
	for _, r := range results {
		users = append(users, r.users...)
	}

	var err error
	users, err = mergeAndSort(users)
	if err != nil {
		return nil, err
	}

	db.finalProgress()

	return users, nil
}

func (db *UsersDB) writeSized() (err error) {
	file, err := os.Create(db.filename)
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

	users, err := db.Users()
	if err != nil {
		return err
	}

	strs := make([]string, len(users))
	for i, u := range users {
		strs[i] = db.printFunc(u)
	}

	length := 0
	for _, s := range strs {
		length += len(s)
	}
	fmt.Fprintf(file, "%d\n", length)

	for _, s := range strs {
		fmt.Fprint(file, s)
	}

	return nil
}

func mergeUsers(userMap map[int][]*User, opts Options, verbatimIndex int) ([]*User, error) {
	idMap := make(map[int]*User)
	for i := 0; i < len(userMap); i++ {
		if i == verbatimIndex {
			for _, user := range idMap {
				user.amend(opts)
			}
		}
		for _, u := range userMap[i] {
			if u == nil || u.ID == "" {
				continue
			}
			id64, err := strconv.ParseUint(u.ID, 10, 24)
			if err != nil {
				return nil, err
			}
			id := int(id64)
			existing, exists := idMap[id]
			if exists {
				idMap[id] = mergeUser(existing, u)
				continue
			}
			idMap[id] = u
		}
	}

	if verbatimIndex < 0 {
		for _, user := range idMap {
			user.amend(opts)
		}
	}

	ids := make([]int, 0, len(idMap))
	for id := range idMap {
		ids = append(ids, id)
	}

	users := make([]*User, len(ids))
	sort.Ints(ids)
	for i, id := range ids {
		users[i] = idMap[id]
	}

	return users, nil
}

func mergeUser(existing, u *User) *User {
	if u.Callsign != "" {
		existing.Callsign = u.Callsign
	}
	if u.Name != "" {
		existing.Name = u.Name
	}
	if u.City != "" {
		existing.City = u.City
	}
	if u.State != "" {
		existing.State = u.State
	}
	if u.Nick != "" {
		existing.Nick = u.Nick
	}
	if u.Country != "" {
		existing.Country = u.Country
	}

	return existing
}

func (db *UsersDB) write(header bool) (err error) {
	file, err := os.Create(db.filename)
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

	if header {
		fmt.Sprintln("Radio ID,CallSign,Name,City,State,Firstname,Country")
	}

	users, err := db.Users()
	if err != nil {
		return err
	}

	for _, u := range users {
		fmt.Fprint(file, db.printFunc(u))
	}

	return nil
}

func (db *UsersDB) WriteMD380ToolsFile(filename string, progress func(cur int) bool) error {
	db.filename = filename
	db.progressCallback = progress
	db.printFunc = func(u *User) string {
		return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n",
			u.ID, u.Callsign, u.Name, u.City, u.State, u.Nick, u.Country)
	}

	return db.writeSized()
}

func (db *UsersDB) WriteMD2017File(filename string, progress func(cur int) bool) error {
	db.filename = filename
	db.progressCallback = progress
	db.printFunc = func(u *User) string {
		return fmt.Sprintf("%s,%s,%s,,%s,%s,%s\n",
			u.ID, u.Callsign, u.Name, u.City, u.State, u.Country)
	}

	header := true
	return db.write(header)
}
