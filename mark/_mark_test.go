package mark

import (
	"testing"
)

func TestFormat(t *testing.T) {
	want := []string{
		OpenItalic + "Italic" + CloseItalic,
		OpenBold + "Bold" + CloseBold,
		OpenBoldItalic + "BoldItalic" + CloseBoldItalic,
		"<" + OpenUnderline + "bracketed" + CloseUnderline + ">",
	}
	args := []string{"*Italic*", "**Bold**", "***BoldItalic***", "<bracketed>"}
	for i, arg := range args {
		t.Logf("testing: %v\n", arg)
		got := Format(arg)
		if got != want[i] {
			t.Errorf("\nwant: %q\ngot:  %q\n", want[i], got)
		}
	}
}

func TestFormatWrapped(t *testing.T) {
	text := `
    Something *easy* to write here that can be indented however you like
		and wrapped and have each line indented and with <code>:

        This will not be messed with.
        Nor this.

    So it's a lot like a **simple** version of Markdown that only supports
    what is likely going to be used in stuff similar to man pages.

    Let's try a hard  
    return.`

	want := "     Something " + OpenItalic + "easy" + CloseItalic + " to write here that can be indented however\n     you like and wrapped and have each line indented and with\n     <" + OpenUnderline + "code" + CloseUnderline + ">:\n     \n         This will not be messed with.\n         Nor this.\n     \n     So it's a lot like a " + OpenBold + "simple" + CloseBold + " version of Markdown that only\n     supports what is likely going to be used in stuff similar to\n     man pages.\n     \n     Let's try a hard\n     return."

	got := FormatWrapped(text, 5, 70)
	t.Log("\n" + got)
	if want != got {
		t.Errorf("\nwant:\n%q\ngot:\n%q\n", want, got)
	}
}

func TestWrapped(t *testing.T) {
	text := `
    Something *easy* to write here that can be indented however you like
		and wrapped and have each line indented and with <code>:

        This will not be messed with.
        Nor this.

    So it's a lot like a **simple** version of Markdown that only supports
    what is likely going to be used in stuff similar to man pages.

    Let's try a hard  
    return.`

	want := "     Something *easy* to write here that can be indented however\n     you like and wrapped and have each line indented and with\n     <code>:\n     \n         This will not be messed with.\n         Nor this.\n     \n     So it's a lot like a **simple** version of Markdown that only\n     supports what is likely going to be used in stuff similar to\n     man pages.\n     \n     Let's try a hard\n     return."

	got := Wrapped(text, 5, 70)
	t.Log("\n" + got)
	if want != got {
		t.Errorf("\nwant:\n%q\ngot:\n%q\n", want, got)
	}
}

func TestPeekWord(t *testing.T) {
	var buf []rune
	var word string
	buf = []rune(`some thing`)
	word = string(peekWord(buf, 0))
	t.Logf("%q", word)
	if word != "some" {
		t.Fail()
	}
	word = string(peekWord(buf, 5))
	t.Logf("%q", word)
	if word != "thing" {
		t.Fail()
	}
	word = string(peekWord(buf, 4))
	t.Logf("%q", word)
	if word != "" {
		t.Fail()
	}
}

func TestWrap(t *testing.T) {
	buf := "Here's a string that's not long."
	want := "Here's a\nstring\nthat's not\nlong."
	got := Wrap(buf, 10)
	if want != got {
		t.Errorf("\nwant: %q\ngot:  %q\n", want, got)
	}
}

func TestWrap_none(t *testing.T) {
	if Wrap("some thing", 0) != "some thing" {
		t.Fail()
	}
}
