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
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

const displaySuffixes = true

const contactSuffixLength = 40

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$@"
const (
	letterIdxBits = 6                    // 6 bits represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits
	letterIdxMax  = 63 / letterIdxBits   // # of indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits,
	// enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func AddSuffix(f *Field, str string) string {
	if f.fType != FtDcName {
		return str
	}

	if f.record.rType != RtContacts {
		return str
	}

	if str == "" {
		return str
	}

	if len(str) > contactSuffixLength {
		return str
	}

	str += "_" + RandomString(contactSuffixLength)

	return str
}

func removeSuffix(f *Field, str string) string {
	if f.fType != FtDcName {
		return str
	}

	if f.record.rType != RtContacts {
		return str
	}

	if len(str) < contactSuffixLength {
		return str
	}

	sepIndex := len(str) - contactSuffixLength - 1
	if str[sepIndex:sepIndex+1] != "_" {
		return str
	}

	str = str[0 : len(str)-contactSuffixLength-1]

	return str
}

func RemoveSuffix(f *Field, str string) string {
	if displaySuffixes {
		return str
	}

	return removeSuffix(f, str)
}

func RemoveSuffixes(strs []string) []string {
	if displaySuffixes {
		return strs
	}

	for i, str := range strs {
		if len(str) < contactSuffixLength {
			continue
		}

		sepIndex := len(str) - contactSuffixLength - 1
		if str[sepIndex:sepIndex+1] != "_" {
			continue
		}

		strs[i] = str[0 : len(str)-contactSuffixLength-1]
	}

	return strs
}

// stringInSlice returns true if the given string exists in the given
// string slice.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// mustBePrintableAscii returns an error if any of the characters in s
// is not an printable ascii character.
func mustBePrintableAscii(s string) error {
	errStr := fmt.Errorf("must contain only printable ASCII characters")

	for _, r := range s {
		if r == 0 || r >= unicode.MaxASCII || !unicode.IsPrint(r) {
			return errStr
		}
	}

	return nil
}

// mustBeNumeric returns an error if any of the characters in s
// is non-numeric.
func mustBeNumericAscii(s string) error {
	errStr := fmt.Errorf("must contain only the characters 0 - 9")

	for _, r := range s {
		if r == 0 || r >= unicode.MaxASCII || !unicode.IsNumber(r) {
			return errStr
		}
	}

	return nil
}

// bytesToFrequency converts a byte slice to a floating point frequency.
func bytesToFrequency(b []byte) float64 {
	return float64(bcdToInt64(bytesToInt64(b))) / 100000
}

// frequencyToBytes converts a floating point freqency to a byte slice.
func frequencyToBytes(f float64) []byte {
	return int64ToBytes(int64ToBcd(int64(math.Round(f*100000))), 4)
}

// frequencyToString produces a string from a floating point frequency.
func frequencyToString(f float64) string {
	return fmt.Sprintf("%3.5f", f)
}

// frequencyToSignedString produces a string from a floating point frequency.
func frequencyToSignedString(f float64) string {
	return fmt.Sprintf("%+3.5f", f)
}

// stringToFrequency converts a string to a floating point frequency.
func stringToFrequency(s string) (float64, error) {
	freq, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("bad frequency")
	}

	return float64(freq), nil
}

// bytesToInt64 converts a byte slice into an integer.
func bytesToInt64(bytes []byte) int64 {
	var i int64

	for j, b := range bytes {
		i |= int64(b) << (uint(j) * 8)
	}

	return i
}

// int64ToBytes converts an integer into a byte slice of length len.
func int64ToBytes(i int64, len int) []byte {
	bytes := make([]byte, len)

	for j := range bytes {
		bytes[j] = byte(i & 0xff)
		i >>= 8
	}

	return bytes
}

// reverse4Bytes returns an integer containing the low four bytes of
// the input in reverse order.
func reverse4Bytes(in int64) int64 {
	out := (in & 0x000000ff) << 24
	out |= (in & 0x0000ff00) << 8
	out |= (in & 0x00ff0000) >> 8
	out |= (in & 0xff000000) >> 24

	return out
}

// bcdToInt64 converts an BCD integer (in standard order) to an int64
func bcdToInt64(bcd int64) int64 {
	var binary int64 = 0
	mult := int64(1)

	for i := 0; i < 8; i++ {
		if (bcd & 0xf) > 9 {
			return -1
		}
		binary += (bcd & 0xf) * mult
		bcd >>= 4
		mult *= 10
	}

	return binary
}

// revBcdToInt64 converts an BCD integer (in reverse order) to an int64
// integer.
func revBcdToInt64(bcd int64) int64 {
	return bcdToInt64(reverse4Bytes(bcd))
}

// int64ToBcd converts an integer to BCD (in standard order).
func int64ToBcd(binary int64) int64 {
	var bcd int64

	for i := 0; i < 8; i++ {
		bcd |= (binary % 10) << uint(4*i)
		binary /= 10
	}

	return bcd
}

// binaryToRevBcd converts an integer to BCD (in reverse order).
func int64ToRevBcd(binary int64) int64 {
	return reverse4Bytes(int64ToBcd(binary))
}

// bcdBytesToString converts a BCD-encodd byte slice to a decimal string
func bcdBytesToString(bytes []byte) string {
	s := ""
	for _, b := range bytes {
		s += string('0' + ((int32(b) >> 4) & 0xf))
		s += string('0' + uint32(b)&0xf)
	}
	return s
}

// stringToBcdBytes converts a decimal string into a BCD-encoded byte slice
func stringToBcdBytes(s string) []byte {
	bytes := make([]byte, len(s)/2)
	for i := range bytes {
		bytes[i] = byte(((s[i*2] - '0') << 4) + s[i*2+1] - '0')
	}
	return bytes
}

// ucs2BytesToString converts a utf16 byte slice into a utf8 string.
func ucs2BytesToString(b []byte) string {
	ucs2S := make([]uint16, 1)
	byteS := make([]byte, 4)
	buf := &bytes.Buffer{}

	for i := 0; i < len(b); i += 2 {
		ucs2S[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)

		if ucs2S[0] == 0 {
			break
		}

		n := utf8.EncodeRune(byteS, utf16.Decode(ucs2S)[0])

		buf.Write([]byte(string(byteS[:n])))
	}

	return buf.String()
}

// stringToUcs2Bytes converts a utf8 string into a utf16 byte slice.
func stringToUcs2Bytes(s string, minLength int) ([]byte, error) {
	runes := make([]rune, 1)
	buf := &bytes.Buffer{}

	for _, r := range s {
		runes[0] = r
		ucs2S := utf16.Encode(runes)
		if len(ucs2S) != 1 {
			return []byte{0}, fmt.Errorf("cannot encode into UCS-2")
		}

		buf.WriteByte(byte(ucs2S[0] & 0xff))
		buf.WriteByte(byte((ucs2S[0] >> 8) & 0xff))
	}

	for buf.Len() < minLength {
		buf.WriteByte(0)
	}

	return buf.Bytes(), nil
}

func stringsEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func maxNamesString(names []string, max int) string {
	if len(names) > max {
		names = append(names[:max], "...")
	}
	return strings.Join(names, ", ")
}
