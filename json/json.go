package json

// String returns an escaped version of the string suitable for
// inclusion as a JSON string value. Unlike Go's standard MarshalJSON,
// this function leaves Unicode as is without escaping it. It also does
// not dubiously escape HTML characters unnecessarily.
//
// The following PEGN describes the formal specification of a JSON
// string:
//
//     String  <-- DQ (Escaped / [x20-x21] / [x23-x5B]
//               / [x5D-x10FFFF])* DQ
//     Escaped  <- BKSLASH ("b" / "f" / "n" / "r"
//               / "t" / "u" hex{4} / DQ / BKSLASH / SLASH)
//
// Unicode and SLASH are not escaped here because they are considered
// valid JSON string data without the escape.
//
// Note that this function follows the mapf format and can be used with
// fn.Map and other mapping functions.
func Escape(in string) string {
	out := ""
	for _, r := range in {
		switch r {
		case '\t':
			out += "\\t"
		case '\b':
			out += "\\b"
		case '\f':
			out += "\\f"
		case '\n':
			out += "\\n"
		case '\r':
			out += "\\r"
		case '\\':
			out += "\\\\"
		case '"':
			out += "\\\""
		default:
			out += string(r)
		}
	}
	return out
}

/*
// Whitespace value selects the white space characters
const Whitespace uint64 = 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' '

// The result of Token is one of these tokens
const (
	EOF = -(iota + 1)
	Type
	PString
	LeftBracket
	RightBracket
	Escape
	Comma
	Invalid
)

var tokenTypeString = map[rune]string{
	EOF:          "EOF",
	Type:         "Type",
	PString:      "PString",
	LeftBracket:  "LeftBracket",
	RightBracket: "RightBracket",
	Escape:       "Escape",
	Comma:        "Comma",
	Invalid:      "Invalid",
}

func isNumber(ch rune) bool { return '0' <= ch && ch <= '9' }

// tokenTypeAsString returns a printable string for a token or Unicode character.
func tokenTypeAsString(tok rune) string {
	if s, found := tokenTypeString[tok]; found {
		return s
	}
	return fmt.Sprintf("%q", string(tok))
}

// A JsonParser implements the parsing of a JSON string to a pegn/Node
type JsonParser struct {
	srcBuf []byte // Input
	srcPos int    // The current position in the source buffer
	srcEnd int    // ending position of the source buffer

	lastCharLen int //length of last character in bytes

	tokPos int          // token text tail position (srcBuf index)
	tokEnd int          // token text tail end (srcBuf index)
	tokBuf bytes.Buffer // token text

	ch rune // character before current srcPos

	line   int // current line for error reportin
	column int // current column for error reporting

	types    []string
	typesMap map[string]int
}

// Report error with position
func (j *JsonParser) Errorf(str string, a ...interface{}) error {
	return fmt.Errorf("jsonparser: %s line: %v, col: %v", fmt.Sprintf(str, a...), j.line, j.column)
}

// Init initializes the JSON parser and takes the byte array to be parsed.
func (j *JsonParser) Init(b []byte) *JsonParser {
	j.srcBuf = b
	j.srcPos = 0
	j.srcEnd = len(j.srcBuf) - 1

	j.lastCharLen = 0

	j.tokPos = -1

	j.ch = -2

	j.line = 0
	j.column = 0

	return j
}

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

	for Whitespace&(1<<uint(ch)) != 0 {
		ch = j.next()
	}

	j.tokPos = j.srcPos - j.lastCharLen

	var tok rune
	var err error

	switch {
	case isNumber(ch):
		tok = Type
		ch = j.tokenizeType()
	default:
		switch ch {
		case EOF:
			return EOF, nil
		case '[':
			ch = j.next()
			tok = LeftBracket
		case ']':
			ch = j.next()
			tok = RightBracket
		case '"':
			ch, err = j.tokenizeString()
			if err != nil {
				return Invalid, err
			}
			tok = PString
		case ',':
			ch = j.next()
			tok = Comma
		default:
			return Invalid, j.Errorf("invalid token")
		}
	}

	j.tokEnd = j.srcPos - j.lastCharLen

	j.ch = ch

	// fmt.Printf("%v - %v\n", tokenTypeAsString(tok), string(ch))

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

	for isNumber(ch) {
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
