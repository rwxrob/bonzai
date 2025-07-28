// Copyright 2022 Robert S. Muhlestein
// SPDX-License-Identifier: Apache-2.0

/*
Package to contains a number of converters that take any number of types and return something transformed from them. It also contains a more granular approach to fmt.Stringer.
*/
package to

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/rwxrob/bonzai/ds/qstack"
	"github.com/rwxrob/bonzai/fn/maps"
	"github.com/rwxrob/bonzai/is"
	"github.com/rwxrob/bonzai/scanner"
)

type Text interface{ string | []rune }

// String first looks for string, []byte, []rune, and io.Reader types
// and if matched returns a string with their content and the string
// type.
//
// String converts whatever remains to that types fmt.Sprintf("%v")
// string version (but avoids calling it if possible). Be sure you use
// things with consistent string representations. Keep in mind that this
// is extremely high level for rapid tooling and prototyping.
func String(in any) string {
	switch v := in.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case []rune:
		return string(v)
	case io.Reader:
		buf, err := io.ReadAll(v)
		if err != nil {
			log.Println(err)
		}
		return string(buf)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Bytes converts whatever is passed into a []byte slice. Logs and
// returns nil if it cannot convert. Supports the following types:
// string, []byte, []rune, io.Reader.
func Bytes(in any) []byte {
	switch v := in.(type) {
	case string:
		return []byte(v)
	case []byte:
		return v
	case []rune:
		return []byte(string(v))
	case io.Reader:
		buf, err := io.ReadAll(v)
		if err != nil {
			log.Println(err)
		}
		return buf
	default:
		log.Printf("cannot convert %T to bytes", in)
		return nil
	}
}

// HumanFriend implementations have a human readable form that is even
// friendlier than fmt.Stringer.
type HumanFriend interface {
	Human() string
}

// Human returns a human-friendly string version of the item,
// specifically:
//
//   - single-quoted runes
//   - double-quoted strings
//   - numbers as numbers
//   - function names are looked up
//   - slices joined with "," and wrapped in []
//
// Anything else is rendered as its fmt.Sprintf("%v",it) form.
func Human(a any) string {
	switch v := a.(type) {

	case string:
		return fmt.Sprintf("%q", v)

	case rune:
		return fmt.Sprintf("%q", v)

	case []string:
		st := []string{}
		for _, r := range v {
			st = append(st, fmt.Sprintf("%q", r))
		}
		return "[" + strings.Join(st, ",") + "]"

	case []rune:
		st := []string{}
		for _, r := range v {
			st = append(st, fmt.Sprintf("%q", r))
		}
		return "[" + strings.Join(st, ",") + "]"

	case []any:
		st := []string{}
		for _, r := range v {
			st = append(st, Human(r))
		}
		return "[" + strings.Join(st, ",") + "]"

	case HumanFriend:
		return v.Human()

	default:
		typ := fmt.Sprintf("%v", reflect.TypeOf(a))
		if len(typ) > 3 && typ[0:4] == "func" {
			return FuncName(a)
		}
		return fmt.Sprintf("%v", a)

	}
}

// FuncName makes a best effort attempt to return the string name of the
// passed function. Anonymous functions are named "funcN" where N is the
// order of appearance within the current scope. Note that this function
// will panic if not passed a function.
func FuncName(i any) string {
	p := runtime.FuncForPC(reflect.ValueOf(i).Pointer())
	n := strings.Split(p.Name(), `.`)
	return n[len(n)-1]
}

// Lines transforms the input into a string and then divides that string
// up into lines (\r?\n) suitable for functional map operations.
func Lines(in any) []string {
	buf := String(in)
	lines := []string{}
	s := bufio.NewScanner(strings.NewReader(buf))
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return lines
}

// Indented returns a string with each line indented by the specified
// number of spaces. Carriage returns are stripped (if found) as
// a side-effect. Per [pkg/text/template] rules, the string input
// argument is last so that this function can be used as-is within
// template pipe constructs:
//
//	{{ .Some | indent 4 }}
func Indented(indent int, in string) string {
	prefix := strings.Repeat(" ", indent)
	return Prefixed(prefix, in)
}

// IndentWrapped adds the specified number of spaces to the beginning of
// every line ensuring that the wrapping is preserved to the specified
// width. See Wrapped and Indented.
func IndentWrapped(indent, width int, in string) string {
	wwidth := width - indent
	body, _ := Wrapped(wwidth, in)
	return Indented(indent, body)
}

// Prefixed returns a string where every line is prefixed. Carriage
// returns (if any) are dropped.
func Prefixed(pre, in string) string {
	lines := Lines(in)
	lines = maps.Prefix(lines, pre)
	return strings.Join(lines, "\n")
}

// PrefixTrimmed returns a string where every line has the passed prefix
// removed if found at the absolute beginning of the line. Carriage
// returns (if any) are dropped.
func PrefixTrimmed(pre, in string) string {
	lines := Lines(in)
	lines = maps.TrimPrefix(lines, pre)
	return strings.Join(lines, "\n")
}

// LinesChopped returns the string with lines chopped from the bottom
// (positive offset) or the top (negative offset).
func LinesChopped(offset int, in string) string {
	lines := Lines(in)
	if offset > 0 {
		return strings.Join(lines[:offset], "\n")
	}
	return strings.Join(lines[-offset:], "\n")
}

var isblank = regexp.MustCompile(`^\s*$`)

// Dedented discards any initial blank lines with nothing but whitespace in
// them and then detects the number and type of whitespace characters at
// the beginning of the first line to the first non-whitespace rune and
// then subsequently removes that number of runes from every
// following line treating empty lines as if they had only n number of
// spaces.  Note that if any line does not have n number of initial
// spaces it the initial runes will still be removed. It is, therefore,
// up to the content creator to ensure that all lines have the same
// space indentation.
func Dedented(in string) string {
	lines := Lines(in)
	for len(lines) == 1 && isblank.MatchString(lines[0]) {
		return ""
	}
	var n int
	for len(lines[n]) == 0 || isblank.MatchString(lines[n]) {
		n++
	}
	starts := n
	indent := Indentation(lines[n])
	for ; n < len(lines); n++ {
		if len(lines[n]) >= indent {
			lines[n] = lines[n][indent:]
		}
	}
	return strings.Join(lines[starts:], "\n")
}

// Indentation returns the number of whitespace runes (in bytes) between
// beginning of the passed string and the first non-whitespace rune.
func Indentation[T Text](in T) int {
	var n int
	var v rune
	for n, v = range []rune(in) {
		if !unicode.IsSpace(v) {
			break
		}
	}
	return n
}

// RuneCount returns the actual number of runes of the string only
// counting the unicode.IsGraphic runes. All others are ignored.  This
// is critical when calculating line lengths for terminal output where
// the string contains escape characters. Note that some runes will
// occupy two columns instead of one depending on the terminal. This
// includes omitting any ASCI terminal escape sequences.
func RuneCount[T string | []byte | []rune](in T) int {
	var c int
	s := scanner.New(in)
	var inesc bool
	for s.Scan() {

		if inesc {
			if s.R == 'm' {
				inesc = false
			}
			continue
		}

		// check for ansi terminal escapes
		if s.R == '\033' {
			m := s.Mark()
			s.Scan()
			if s.R != '[' {
				s.Goto(m)
				continue
			}
			inesc = true
			continue
		}

		if unicode.IsGraphic(s.R) {
			c++
		}
	}
	return c
}

// Words will return the string will all contiguous runs of
// unicode.IsSpace runes converted into a single space. All leading and
// trailing white space will also be trimmed.
func Words(it string) string {
	return strings.Join(qstack.Fields(it).Items(), " ")
}

// Wrapped returns a word wrapped string at the given boundary width
// (in bytes) and the count of words contained in the string.  All
// white space is compressed to a single space. Any width less than
// 1 will simply trim and crunch white space returning essentially the
// same string and the word count.  If the width is less than any given
// word at the start of a line than it will be the only word on the line
// even if the word length exceeds the width. No attempt at
// word-hyphenation is made. Note that white space is defined as
// unicode.IsSpace and does not include control characters. Anything
// that is not unicode.IsSpace or unicode.IsGraphic will be ignored in
// the column count. Any terminal escapes that begin with \033[ will
// also be kept automatically out of calculations. See Unescaped.
func Wrapped(width int, it string) (string, int) {
	words := qstack.Fields(it)
	if width < 1 {
		return strings.Join(words.Items(), " "), words.Len
	}
	var curwidth int
	var wrapped string
	var line []string
	for words.Scan() {
		cur := words.Current()
		count := RuneCount(cur)
		if len(line) == 0 {
			line = append(line, cur)
			curwidth += count
			continue
		}
		if curwidth+count+1 > width {
			wrapped += strings.Join(line, " ") + "\n"
			curwidth = count
			line = []string{cur}
			continue
		}
		line = append(line, cur)
		curwidth += RuneCount(cur) + 1
	}
	wrapped += strings.Join(line, " ")
	return wrapped, words.Len
}

// MergedMaps combines the maps with "last wins" priority. Always
// returns a new map of the given type, even if empty.
func MergedMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	combined := map[K]V{}
	for _, m := range maps {
		for k, v := range m {
			combined[k] = v
		}
	}
	return combined
}

// StopWatch converts a duration into a string that one would expect to
// see on a stopwatch.
func StopWatch(dur time.Duration) string {
	var out string

	sec := dur.Seconds()
	if sec < 0 {
		out += "-"
	}
	sec = math.Abs(sec)

	if sec >= 3600 {
		hours := sec / 3600
		sec = math.Mod(sec, 3600)
		out += fmt.Sprintf("%v:", int(hours))
	}

	if sec >= 60 {
		var form string
		mins := sec / 60
		sec = math.Mod(sec, 60)
		if len(out) == 0 {
			form = `%v:`
		} else {
			form = `%02v:`
		}
		out += fmt.Sprintf(form, int(mins))
	}

	var form string
	if len(out) == 0 {
		form = `%02v`
	} else {
		form = `%02v`
	}
	out += fmt.Sprintf(form, int(sec))

	return out
}

// EscReturns changes any actual carriage returns or line returns into
// their backslashed equivalents and returns a string. This is different
// than Sprintf("%q") since that escapes several other things.
func EscReturns[T string | []byte | []rune](in T) string {
	runes := []rune(string(in))
	var out string
	for _, r := range runes {
		switch r {
		case '\r':
			out += "\\r"
		case '\n':
			out += "\\n"
		default:
			out += string(r)
		}
	}
	return string(out)
}

// UnEscReturns changes any escaped carriage returns or line returns into
// their actual values.
func UnEscReturns[T string | []byte | []rune](in T) string {
	runes := []rune(string(in))
	var out string
	for n := 0; n < len(runes); n++ {
		if runes[n] == '\\' && runes[n+1] == 'r' {
			out += "\r"
			n++
			continue
		}
		if runes[n] == '\\' && runes[n+1] == 'n' {
			out += "\n"
			n++
			continue
		}
		out += string(runes[n])
	}
	return string(out)
}

// HTTPS simply adds the prefix "https://" if not found. Useful for
// allowing non-prefixed URLs and later converting them.
func HTTPS(url string) string {
	if len(url) < 8 || url[0:8] != "https://" {
		return "https://" + url
	}
	return url
}

// CrunchSpace crunches all unicode.IsSpace into a single space. It does
// not trim. See TrimCrunchSpace.
func CrunchSpace(in string) string {
	runes := make([]rune, 0)
	s := scanner.New(in)
	var inspace bool
	for s.Scan() {
		r := s.Rune()
		if unicode.IsSpace(r) {
			if inspace {
				continue
			}
			runes = append(runes, ' ')
			inspace = true
			continue
		}
		inspace = false
		runes = append(runes, r)
	}
	return string(runes)
}

// CrunchSpaceVisible crunches all unicode.IsSpace into a single space
// and filters out anything that is not unicode.IsPrint. It does not
// trim. See TrimCrunchSpaceVisible. Note that this requires three
// passes through the string in order to resolve any white space that
// might have been separated by escape and other characters.
func CrunchSpaceVisible(in string) string {
	in = CrunchSpace(in)
	in = Visible(in)
	in = CrunchSpace(in)
	return in
}

// Visible filters out any rune that is not unicode.IsPrint().
func Visible(in string) string {
	runes := make([]rune, 0)
	s := scanner.New(in)
	for s.Scan() {
		r := s.Rune()
		if unicode.IsPrint(r) {
			runes = append(runes, r)
		}
	}
	return string(runes)
}

// Type attempts to convert a string [from] into the type of the
// provided [this] fallback value. It uses the type of [this] to
// determine the conversion logic and returns the converted value if
// successful. If conversion fails, it returns [this] as a fallback.
func Type[T any](from string, this T) T {
	var result any = this
	var err error
	switch any(this).(type) {
	case bool:
		result = is.Truthy(from)
	case int:
		result, err = strconv.Atoi(from)
		if err != nil {
			return this
		}
	case float64:
		result, err = strconv.ParseFloat(from, 64)
		if err != nil {
			return this
		}
	case string:
		result = from
	}
	return result.(T)
}

// TrimVisible removes anything but unicode.IsPrint and then trims. It
// does not crunch spaces, however.
func TrimVisible(in string) string { return strings.TrimSpace(Visible(in)) }

// TrimCrunchSpace is same as CrunchSpace but trims initial and trailing
// space.
func TrimCrunchSpace(in string) string { return strings.TrimSpace(CrunchSpace(in)) }

// TrimCrunchSpaceVisible is same as CrunchSpaceVisible but trims initial and trailing
// space.
func TrimCrunchSpaceVisible(in string) string { return strings.TrimSpace(CrunchSpaceVisible(in)) }
