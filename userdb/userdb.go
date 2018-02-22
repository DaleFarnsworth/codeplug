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
var marcUsersURL = "http://www.dmr-marc.net/cgi-bin/trbo-database/datadump.cgi?table=users&format=csv&header=0"
var reflectorUsersURL = "http://registry.dstar.su/reflector.db"

var timeoutSeconds = 20

var tr = &http.Transport{
	TLSHandshakeTimeout:   time.Duration(timeoutSeconds) * time.Second,
	ResponseHeaderTimeout: time.Duration(timeoutSeconds) * time.Second,
}

var client = &http.Client{
	Transport: tr,
	Timeout:   time.Duration(timeoutSeconds) * time.Second,
}

type User struct {
	ID       string
	Callsign string
	Name     string
	City     string
	State    string
	Country  string
}

type UsersDB struct {
	filename          string
	userFunc          func(*User) string
	progressCallback  func(progressCounter int) bool
	progressFunc      func() error
	progressIncrement int
	progressCounter   int
}

func newUserDB() *UsersDB {
	db := &UsersDB{
		progressFunc: func() error { return nil },
	}

	return db
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

func (u *User) normalize() {
	u.Callsign = normalizeString(u.Callsign)
	u.Name = normalizeString(u.Name)
	u.City = normalizeString(u.City)
	u.State = normalizeString(u.State)
	u.Country = normalizeString(u.Country)
}

func normalizeString(s string) string {
	s = asciify(s)
	s = strings.TrimSpace(s)
	s = strings.Replace(s, ",", ";", -1)

	for strings.Index(s, "  ") >= 0 {
		s = strings.Replace(s, "  ", " ", -1)
	}

	return s
}

func asciify(s string) string {
	runes := []rune(s)
	strs := make([]string, len(runes))
	for i, r := range runes {
		strs[i] = transliterations[r]
	}

	return strings.Join(strs, "")
}

func getBytes(url string) ([]byte, error) {
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

func getLines(url string) ([]string, error) {
	bytes, err := getBytes(url)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(bytes), "\n")

	return lines[:len(lines)-1], nil
}

func getMarcUsers() ([]*User, error) {
	lines, err := getLines(marcUsersURL)
	if err != nil {
		errFmt := "error getting MARC users database: %s: %s"
		err = fmt.Errorf(errFmt, marcUsersURL, err.Error())
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
			Country:  fields[5],
		}
	}
	return users, nil
}

func getFixedUsers() ([]*User, error) {
	lines, err := getLines(fixedUsersURL)
	if err != nil {
		errFmt := "error getting fixed users: %s: %s"
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
	bytes, err := getBytes(specialUsersURL)
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
	lines, err := getLines(url)
	if err != nil {
		errFmt := "error getting special users: %s: %s"
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
	lines, err := getLines(reflectorUsersURL)
	if err != nil {
		errFmt := "error getting reflector users: %s: %s"
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

func deDupAndSort(users []*User) ([]*User, error) {
	idMap := make(map[int]*User)
	for _, u := range users {
		if u == nil || u.ID == "" {
			continue
		}
		id, err := strconv.ParseUint(u.ID, 10, 24)
		if err != nil {
			return nil, err
		}
		idMap[int(id)] = u
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
	users []*User
	err   error
}

func do(f func() ([]*User, error), resultChan chan result) {
	var r result

	r.users, r.err = f()
	resultChan <- r
}

func (db *UsersDB) Users() ([]*User, error) {
	getUsersFuncs := []func() ([]*User, error){
		getFixedUsers,
		getMarcUsers,
		getReflectorUsers,
	}

	specialURLs, err := getSpecialURLs()
	if err != nil {
		return nil, err
	}
	for i := range specialURLs {
		url := specialURLs[i]
		f := func() ([]*User, error) {
			return getSpecialUsers(url)
		}
		getUsersFuncs = append(getUsersFuncs, f)
	}

	var users []*User
	resultChan := make(chan result, len(getUsersFuncs))

	for _, f := range getUsersFuncs {
		go do(f, resultChan)
	}

	db.setMaxProgressCount(len(getUsersFuncs))

	for done := 0; done < len(getUsersFuncs); {
		select {
		case r := <-resultChan:
			if r.err != nil {
				return nil, r.err
			}
			users = append(users, r.users...)
			done++

			err := db.progressFunc()
			if err != nil {
				return nil, err
			}
		}
	}

	users, err = deDupAndSort(users)
	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].normalize()
	}

	db.finalProgress()

	return users, nil
}

func (db *UsersDB) writeSizedUsersFile() (err error) {
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
		strs[i] = db.userFunc(u)
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

func (db *UsersDB) writeUsersFile() (err error) {
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

	fmt.Sprintln("Radio ID,CallSign,Name,NickName,City,State,Country")

	users, err := db.Users()
	if err != nil {
		return err
	}

	for _, u := range users {
		fmt.Fprint(file, db.userFunc(u))
	}

	return nil
}

func WriteMD380ToolsFile(filename string, euro bool, progress func(cur int) bool) error {
	db := newUserDB()
	db.filename = filename
	db.progressCallback = progress
	db.userFunc = func(u *User) string {
		return fmt.Sprintf("%s,%s,%s,%s,%s,,%s\n",
			u.ID, u.Callsign, u.Name, u.City, u.State, u.Country)
	}
	if euro {
		db.userFunc = func(u *User) string {
			return fmt.Sprintf("%s,%s,,%s,%s,,%s\n",
				u.ID, u.Callsign, u.City, u.State, u.Country)
		}
	}
	return db.writeSizedUsersFile()
}

func WriteMD2017File(filename string, euro bool, progress func(cur int) bool) error {
	db := newUserDB()
	db.filename = filename
	db.progressCallback = progress
	db.userFunc = func(u *User) string {
		return fmt.Sprintf("%s,%s,%s,,%s,%s,%s\n",
			u.ID, u.Callsign, u.Name, u.City, u.State, u.Country)
	}
	if euro {
		db.userFunc = func(u *User) string {
			return fmt.Sprintf("%s,%s,,,%s,%s,%s\n",
				u.ID, u.Callsign, u.City, u.State, u.Country)
		}
	}
	return db.writeUsersFile()
}
