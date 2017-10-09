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
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

// randomString returns a a random hex string of the given length.
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
	return float64(bcdToBinary(bytesToInt(b))) / 100000
}

// frequencyToBytes converts a floating point freqency to a byte slice.
func frequencyToBytes(f float64) []byte {
	return intToBytes(binaryToBcd(int(f*100000)), 4)
}

// frequencyToString produces a string from a floating point frequency.
func frequencyToString(f float64) string {
	return fmt.Sprintf("%3.5f", f)
}

// stringToFrequency converts a string to a floating point frequency.
func stringToFrequency(s string) (float64, error) {
	freq, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("bad frequency")
	}

	return float64(freq), nil
}

// bytesToInt converts a byte slice into an integer.
func bytesToInt(bytes []byte) int {
	var i int

	for j, b := range bytes {
		i |= int(b) << (uint(j) * 8)
	}

	return i
}

// intToBytes converts an integer into a byte slice of length len.
func intToBytes(i int, len int) []byte {
	bytes := make([]byte, len)

	for j := range bytes {
		bytes[j] = byte(i & 0xff)
		i >>= 8
	}

	return bytes
}

// reverse4Bytes returns an integer containing the low four bytes of
// the input in reverse order.
func reverse4Bytes(in int) int {
	out := (in & 0x000000ff) << 24
	out |= (in & 0x0000ff00) << 8
	out |= (in & 0x00ff0000) >> 8
	out |= (in & 0xff000000) >> 24

	return out
}

// bcdToBinary converts an BCD integer (in standard order) to a binary integer.
func bcdToBinary(bcd int) int {
	binary := 0
	mult := 1

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

// revBcdToBinary converts an BCD integer (in reverse order) to a binary
// integer.
func revBcdToBinary(bcd int) int {
	return bcdToBinary(reverse4Bytes(bcd))
}

// binaryToBcd converts an integer to BCD (in standard order).
func binaryToBcd(binary int) int {
	bcd := 0

	for i := 0; i < 8; i++ {
		bcd |= (binary % 10) << uint(4*i)
		binary /= 10
	}

	return bcd
}

// binaryToRevBcd converts an integer to BCD (in reverse order).
func binaryToRevBcd(binary int) int {
	return reverse4Bytes(binaryToBcd(binary))
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
