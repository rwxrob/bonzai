/*
Package json contains the PrintAsJSON interface and functions for
marshaling and unmarshaling that are exact port of the Go json package
without reflection and the completely unnecessary escaping of every
single Unicode character making the output readable by anyone speaking
or writing languages that depend heavily on Unicode The standard package
output is completely unreadable, which is really unfortunate since the
JSON standard fully supports Unicode characters as-is within strings.
The Bonzai json package overcomes these limitations and is used to
marshal everything in the module and fulfill all `fmt.Stringer`
interfaces as 2-space indented JSON instead of Go's virtually unusable
default string marshaling format.
*/
package json

// PrintAsJSON provides a consistent representation of any structure
// such that it an easily be read and compared as JSON whenever printed
// and test. Sadly, the default string representations for most types in
// Go are virtually unusable for consistent representations of any
// structure. And while it is true that JSON data should be supported in
// any way that is it presented, some consistent output makes for more
// consistent debugging, documentation, and testing.
//
// PrintAsJSON implementations must fulfill their marshaling to any
// textual/printable representation by producing compressed,
// single-line, no-spaces, parsable, shareable not-unnecessarily-escaped
// JSON data. When indented, long-form JSON is wanted utility functions
// and utilities (such as jq) can be used to expand the compressed,
// default JSON.
//
// All implementations must also implement the fmt.Stringer interface by
// calling JSON(), which most closely approximates the Go standard
// string marshalling. Thankfully, Go does ensure that the order of
// elements in any type will appear consistently in that same order
// during testing even though they should never be relied upon for such
// ordering other than in testing.
//
// All implementations must promise they will never escape any string
// character that is not specifically required to be escaped by the JSON
// standard as described in this PEGN specification:
//
//     String  <-- DQ (Escaped / [x20-x21] / [x23-x5B]
//                 / [x5D-x10FFFF])* DQ
//     Escaped  <- BKSLASH ("b" / "f" / "n" / "r" / "t" / "u" hex{4}
//                 / DQ / BKSLASH / SLASH)
//
// This means that binary data will never be represented as Unicode
// escapes of any kind and should always be converted to base64 encoding
// and expressed as a string value.
//
// For some, this means that MarshalJSON might need to be be implemented
// since the standard json.Marshal package unnecessarily escapes HTML
// and other characters unnecessarily even though they are a perfectly
// safe and accepted standard (and always have been). PrintAsJSON
// implementers promise those using their implementations will produce
// JSON output that is highly consistent. To facilitate this,
// implementations may implement MarshalJSON methods that delegate to
// bonzai/json.Marshal as a quick fix.
//
// In general, implementers of PrintAsJSON should not depend on
// structure tagging and reflection for unmarshaling instead
// implementing their own consistent UnmarshalJSON method. This allows
// for better error checking as the default does nothing to ensure that
// unmarshaled values are within acceptable ranges. Errors are only
// generated if the actual JSON syntax itself is incorrect.
type PrintAsJSON interface {
	// MarshalJSON() ([]byte, error) // encouraged, but not required
	JSON() string   // compressed, single line, no spaces, no extra escapes
	String() string // must return s.JSON()
	Print() string  // must call fmt.Println(s.JSON())
	Log() string    // must call log.Print(s.JSON())
}

// WS contains all JSON valid whitespace values.
const WS uint64 = 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' '

const (
	EOF = -(iota + 1)
	TYPE
	PSTRING
	LBRACKET
	RBRACKET
	ESC
	COMMA
	INVALID
)

var Token = map[rune]string{
	EOF:      `EOF`,
	TYPE:     `TYPE`,
	PSTRING:  `PSTRING`,
	LBRACKET: `LBRACKET`,
	RBRACKET: `RBRACKET`,
	ESC:      `ESC`,
	COMMA:    `COMMA`,
	INVALID:  `INVALID`,
}

// Escape returns an escaped version of the string suitable for
// inclusion as a JSON string value. Unlike Go's standard MarshalJSON,
// this function leaves Unicode as is without escaping it. It also does
// not dubiously escape HTML characters unnecessarily.
//

// Unicode and SLASH are not escaped here because they are considered
// valid JSON string data without the escape.
//
// Note that this function follows the mapf format and can be used with
// fn.Map and other mapping functions (see json_test.go for examples).
func Escape(in string) string {
	out := ``
	for _, r := range in {
		switch r {
		case '\t':
			out += `\t`
		case '\b':
			out += `\b`
		case '\f':
			out += `\f`
		case '\n':
			out += `\n`
		case '\r':
			out += `\r`
		case '\\':
			out += `\\`
		case '"':
			out += `\"`
		default:
			out += string(r)
		}
	}
	return out
}

/*
func isnum(r rune) bool { return '0' <= r && r <= '9' }

// Scanner scans JSON similar to the standard JSON package.
type Scanner struct {
	buf  []byte // input buffer of data
	beg  int    // beg position in buf
	end  int    // end position in buf
	tbuf []byte // output buffer of tokens
	tbeg int    // token text beg in sbuf
	tend int    // token text end in sbuf
	r    rune   // last scanned rune (before sbeg)
	rlen int    // current rune length
	line int    // line for error reporting
	col  int    // col for error reporting
}

// NewScanner returns a newly initialized scanner from any string,
// []byte, or io.Reader. See Init.
func NewScanner(in any) (*Scanner, error) {
	s := new(Scanner)
	err := s.Init(in)
	return s, err
}

// Init (re)initializes scanner and loads the data from a string,
// []byte, or io.Reader returning nil and an error if any are
// encountered during loading of the input data buffer.
func (s *Scanner) Init(in any) error {
	switch v := in.(type) {
	case string:
		s.buf = []byte(v)
	case []byte:
		s.buf = v
	case io.Reader:
		var err error
		if s.buf, err = io.ReadAll(v); err != nil {
			return err
		}
	}
	s.beg = 0
	s.end = len(s.buf) - 1
	s.rlen = 0
	s.tbeg = -1
	s.r = -2
	s.line = 0
	s.col = 0
	return nil
}

const error_form = `scanner: %s (line:%v col:%v)`

// Errorf returns an error updated with the parser position.
func (s *Scanner) Errorf(e string, a ...any) error {
	return fmt.Errorf(error_form, fmt.Sprintf(e, a...), s.line, s.col)
}
*/

/*


// next returns the next rune in the source byte array.
func (j *JsonParser) next() rune {
	if j.srcPos > j.srcEnd {
		return EOF
	}

	ch, width := rune(j.srcBuf[j.srcPos]), 1
	if ch >= utf8.RuneSelf {
		ch, width = utf8.DecodeRune(j.srcBuf[j.srcPos:j.srcEnd])
		if ch == utf8.RuneError && width == 1 {
			j.srcPos += width
			j.lastCharLen = width
			return ch
		}
	}

	j.srcPos += width
	j.lastCharLen = width

	j.column++

	switch ch {
	case '\n':
		j.line++
		j.column = 0
	}

	return ch
}

// peek returns the rune at the current source byte array position.
func (j *JsonParser) peek() rune {
	if j.ch == -2 {
		j.ch = j.next()
	}
	return j.ch
}

// Parses the byte array into a pegn Node using string types
func (j *JsonParser) ParseTypes(types []string) error {
	j.types = types
	j.typesMap = make(map[string]int)
	for i := 0; i < len(types); i++ {
		j.typesMap[types[i]] = i
	}
	n := new(Node)
	j.Parse(n)
	return nil
}

// Parse parses the byte array into a pegn Node
func (j *JsonParser) Parse(n *Node) error {
	tok, err := j.token()
	if err != nil {
		return err
	}
	switch tok {
	case LeftBracket:
		err = j.parseNode(n)
		if err != nil {
			return err
		}
		tok, err := j.token()
		if err != nil {
			return err
		}
		switch tok {
		case EOF:
			return nil
		default:
			return j.Errorf("extra characters after json string")
		}

	default:
		return j.Errorf("unsupported input")
	}
}

// token captures the next token in the source byte array and returns
// the token type
func (j *JsonParser) token() (rune, error) {
	ch := j.peek()

	j.tokPos = -1

	for WS&(1<<uint(ch)) != 0 {
		ch = j.next()
	}

	j.tokPos = j.srcPos - j.lastCharLen

	var tok rune
	var err error

	switch {
	case isnum(ch):
		tok = Type
		ch = j.tokenizeType()
	default:
		switch ch {
		case EOF:
			return EOF, nil
		case '[':
			ch = j.next()
			tok = LBRACKET
		case ']':
			ch = j.next()
			tok = RBRACKET
		case '"':
			ch, err = j.tokenizeString()
			if err != nil {
				return Invalid, err
			}
			tok = PSTRING
		case ',':
			ch = j.next()
			tok = COMMA
		default:
			return INVALID, j.Errorf("invalid token")
		}
	}

	j.tokEnd = j.srcPos - j.lastCharLen

	j.ch = ch

	return tok, nil
}

func (j *JsonParser) parseNode(n *Node) error {

	tok, err := j.token()

	switch tok {
	// Empty Node
	case RightBracket:
		return nil
	case Type:
		n.Type, err = j.tokenValueAsType()
		if err != nil {
			return err
		}
		return j.parseValueOrChildren(n)
	case PString:
		n.Type, err = j.tokenStringValueAsType()
		if err != nil {
			return err
		}
		return j.parseValueOrChildren(n)

	default:
		return j.Errorf("expecting node type")
	}
}

func (j *JsonParser) parseValueOrChildren(n *Node) error {
	tok, err := j.token()
	if err != nil {
		return err
	}
	switch tok {
	// No Value
	case RightBracket:
		return nil
	case Comma:
		tok, err := j.token()
		if err != nil {
			return err
		}
		switch tok {
		//Node Value
		case PString:
			n.Value = j.tokenValueAsString()

			tok, err = j.token()
			if err != nil {
				return err
			}
			switch tok {
			//End of Node
			case RightBracket:
				return nil
			default:
				return j.Errorf("expecting ']' after node value")
			}
		case LeftBracket:
			tok, err = j.token()
		nextChild:
			if err != nil {
				return err
			}
			switch tok {
			case LeftBracket:
				cn := new(Node)
				err := j.parseNode(cn)
				if err != nil {
					return err
				}
				n.AppendChild(cn)
				tok, err = j.token()
				if err != nil {
					return err
				}
				switch tok {
				// Another Child
				case Comma:
					tok, err = j.token()
					if err != nil {
						return err
					}
					switch tok {
					case LeftBracket:
						goto nextChild
					default:
						return j.Errorf("expecting child")
					}
				// End of Children
				case RightBracket:
					break
				default:
					return j.Errorf("expecting comma or end of children array")
				}
			default:
				return j.Errorf("expecting child node")
			}
		}
	default:
		return j.Errorf("invalid input, expecting comma after node type")
	}

	tok, err = j.token()
	if err != nil {
		return err
	}
	switch tok {
	// End of Node
	case RightBracket:
		return nil
	default:
		return j.Errorf("missing ending ']'")
	}

	return nil
}

func (j *JsonParser) tokenizeType() rune {
	ch := j.next()

	for isnum(ch) {
		ch = j.next()
	}
	return ch
}

func (j *JsonParser) tokenizeString() (rune, error) {
	tokStart := j.srcPos

	j.tokBuf.Reset()
	for ch := j.next(); ch != '"'; ch = j.next() {
		switch ch {
		case EOF:
			return Invalid, j.Errorf("unterminated string")
		case '\n':
			return Invalid, j.Errorf("invalid newline in string")
		case '\\':
			width := 2
			var es string
			ch = j.next()
			switch ch {
			case '"':
				es = "\""
			case '\\':
				es = "\\"
			case '/':
				es = "/"
			case 'b':
				es = "\b"
			case 'f':
				es = "\f"
			case 'n':
				es = "\n"
			case 'r':
				es = "\r"
			case 't':
				es = "\t"
			case 'u':
				return Invalid, j.Errorf("unicode escape not supported in ast json")
			default:
			}
			j.tokBuf.Write(j.srcBuf[tokStart : j.srcPos-width])
			j.tokBuf.WriteString(es)
			tokStart = j.srcPos
		}
	}
	j.tokBuf.Write(j.srcBuf[tokStart : j.srcPos-1])
	return j.next(), nil
}

// tokenValueAsString returns the most recently tokenzied value as
// a string
func (j *JsonParser) tokenValueAsString() string {
	switch {
	case j.tokPos < 0:
		return ""
	default:
		if j.tokBuf.Len() > 0 {
			return j.tokBuf.String()
		}
		return string(j.srcBuf[j.tokPos+1 : j.tokEnd-1])
	}
}

func (j *JsonParser) tokenStringValueAsType() (int, error) {
	return 1, nil
}

// tokenValueAsType returns the most recently tokenized value as
// an integer representing the node type
func (j *JsonParser) tokenValueAsType() (int, error) {
	s := string(j.srcBuf[j.tokPos:j.tokEnd])
	r, err := strconv.Atoi(s)
	return r, err
}
*/
