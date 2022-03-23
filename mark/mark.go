package mark

import (
	"github.com/rwxrob/scan"
	z "github.com/rwxrob/scan/is"
	"github.com/rwxrob/scan/pg"
	"github.com/rwxrob/structs/tree"
)

func Parse(in any) (*tree.E[string], error) {

	names := []string{
		`Grammar`,  // 1
		`Bulleted`, // 2
		`Numbered`, // ...
		`Verbatim`,
		`Paragraph`,
		`BoldItalic`,
		`Bold`,
		`Italic`,
		`Bracketed`,
		`Plain`,
	}
	s := scan.New(names)
	s.X(Grammar)
	s.Print()
	s.Tree.Root.Print()

	return nil, nil
}

func Grammar(s *scan.R) bool { return s.X(z.P{1, z.X{Block, pg.EndLine}}) }

func Block(s *scan.R) bool {
	return s.X(z.X{z.I{Bulleted, Numbered, Paragraph, Verbatim}, EndBlock})
}

func EndBlock(s *scan.R) bool { return s.X(z.I{z.M{2, pg.EndLine}, scan.EOD}) }

func Bulleted(s *scan.R) bool {
	return s.X(z.P{2, z.M1{z.X{z.N{pg.EndLine}, '*', ' ', z.M0{Span}}}})
}

func Numbered(s *scan.R) bool {
	return s.X(z.P{3, z.M1{z.X{z.N{pg.EndLine}, "1.", ' ', z.M0{Span}}}})
}

func Verbatim(s *scan.R) bool {
	return s.X(z.P{4, z.M1{z.X{z.N{pg.EndLine}, z.C{4, ' '}, z.M1{pg.UGraphic}}}})
}

func Paragraph(s *scan.R) bool {
	return s.X(z.P{4, z.M1{z.X{z.N{z.I{z.C{2, pg.EndLine}, scan.EOD}}, Span}}})
}

func Span(s *scan.R) bool { return s.X(z.I{BoldItalic, Bold, Italic, Bracketed, Plain}) }

func BoldItalic(s *scan.R) bool {
	return s.X(z.P{6, z.X{"***", z.N{pg.WS}, z.M1{Plain}, z.N{pg.WS}, "***"}})
}

func Bold(s *scan.R) bool {
	return s.X(z.P{7, z.X{"**", z.N{pg.WS}, z.M1{Plain}, z.N{pg.WS}, "**"}})
}

func Italic(s *scan.R) bool {
	return s.X(z.P{8, z.X{"*", z.N{pg.WS}, z.M1{Plain}, z.N{pg.WS}, "*"}})
}

func Bracketed(s *scan.R) bool {
	//return s.X(z.P{9, z.X{'<', z.N{pg.WS}, z.M1{Plain}, z.N{pg.WS}, '>'}})
	return s.X('<', z.P{9, z.X{z.M1{z.X{z.N{'>'}, pg.UGraphic}}}}, '>')
}

func Plain(s *scan.R) bool { return s.X(z.P{10, z.M1{pg.UGraphic}}) }
