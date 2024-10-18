/*
Package uniq provides several universal unique identifiers.

Package uniq is a utility package that provides common random unique
identifiers in UUID, Base32, and n*2 random hexadecimal characters.

	    6c671957-2f39-4ce5-9f0e-e8d5ec53bfde (16 bytes, 36 chars, hex-)
	    H6M0STKP0MTSU0493GERQDCSJ5BMF3VO     (20 bytes, 32 chars, base32)
	    20060102150405                       (ISO8601 seconds without punch)
	    20060102T150405Z                     (ISO8601 with letters, not punch)
			2006-01-02T15:04:05Z                 (ISO8601 readable)
	    1561158139                           (8 bytes, 10+ chars, int64)
	    5b ...                               (n bytes, n*2 chars, hex)

When a simple random identifier is all that is needed Base32() provides
a better alternative to UUID(). It takes less space (32 characters), is safe
for use with all filesystems including case-insensitive ones, and provides
additional randomness increased from 2^128 (uuid) to 2^160 (base32).

This package includes the following convenience commands as well for use when
integrating with shell scripts:

	uuid
	uid32
	isosec
	isosect
	isodate
	isonan
	epoch [SECONDS]
	randhex [COUNT]
*/
package uniq

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"time"
)

// Bytes returns n number of cryptographically secure pseudo-random
// bytes or a zero length slice if unable to read them from the hosting
// device.
func Bytes(n int) []byte {
	byt := make([]byte, n)
	_, err := rand.Read(byt)
	if err != nil {
		return []byte{}
	}
	return byt
}

// Hex returns a random hexadecimal string that is n*2 characters in
// length.  Calling Hex(18) is superior to UUID() and fits into the same
// 36 character space. Base32() remains superior, but sometimes content
// limitations and validators require only hexadecimal characters. Hex()
// can also be used to generate random RGB colors with Hex(3). Returns
// empty string if unable to read random data from host device.
func Hex(n int) string {
	return hex.EncodeToString(Bytes(n))
}

// Base32 returns a base32 encoded 20 byte string. This has a greater
// range than UUID() and is safe for use with filesystems. Base32 is
// rendered in uppercase for clarity and because it is case insensitive.
// Base32 depends on 40 bit chunks. 20 bytes exceeds UUID() randomness
// and is the closest. (15 would be insufficient to cover the same
// range.) Base32() is therefore superior to UUID() both in range of
// randomness and practicality.  Returns an empty string if unable to
// read random data.
func Base32() string {
	byt := make([]byte, 20)
	_, err := rand.Read(byt)
	if err != nil {
		return ""
	}
	return base32.HexEncoding.EncodeToString(byt)
}

// UUID returns a standard UUID v4 according to RFC 4122. UUIDs have
// become deeply entrenched universally but are inferior to Base32 as
// the need for a greater range of randomness emerges (IPv6, etc.)
// Returns empty string if unable to read random data for any reason.
func UUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return ""
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// Second returns the current (time.Unix()) second since Jan 1, 1970
// UTC. This is frequently a very good unique suffix that has the added
// advantage of being chronologically sortable.
func Second() string {
	return fmt.Sprintf("%v", time.Now().Unix())
}

// Isosec returns the GMT current time in ISO8601 (RFC3339) without
// any punctuation or the T.  This is frequently a very good unique
// suffix that has the added advantage of being chronologically sortable
// and more readable than the epoch. (Also see Second())
func Isosec() string {
	return fmt.Sprintf("%v", time.Now().In(time.UTC).Format("20060102150405"))
}

// IsosecT is same as Isosec with but with the preferred T.
func IsosecT() string {
	return fmt.Sprintf("%v", time.Now().In(time.UTC).Format("20060102T150405"))
}

// Isodate is a human-friendly date and time with Z for UTC and T to
// avoid space (per ISO8601). This identifier does have a space in it,
// but is more compatible with databases.
func Isodate() string {
	return fmt.Sprintf("%v", time.Now().In(time.UTC).Format("2006-01-02T15:04:05Z"))
}

// Isonan returns the GMT current time in ISO8601 (RFC3339) but for
// nanoseconds without any punctuation or the T.  This is frequently
// a very good unique suffix that has the added advantage of being
// chronologically sortable and more readable than the epoch and
// provides considerably more granularity than just Second. Note that
// the length of these strings varies depending on the operating system.
func Isonan() string {
	t := time.Now()
	return fmt.Sprintf("%v%v",
		t.In(time.UTC).Format("20060102150405"),
		t.In(time.UTC).Format(".999999999")[1:],
	)
}
