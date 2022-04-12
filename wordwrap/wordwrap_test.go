package wordwrap

import (
	"io"
	"testing"
)

func TestWordWrap(t *testing.T) {
	tt := []struct {
		Input        string
		Expected     string
		Limit        int
		KeepNewlines bool
		NewLineFunc  NewLineFunc
	}{
		// No-op, should pass through, including trailing whitespace:
		{
			"foobar\n ",
			"foobar\n ",
			0,
			true,
			nil,
		},
		// Nothing to wrap here, should pass through:
		{
			"foo",
			"foo",
			4,
			true,
			nil,
		},
		// A single word that is too long passes through.
		// We do not break long words:
		{
			"foobarfoo",
			"foobarfoo",
			4,
			true,
			nil,
		},
		// Lines are broken at whitespace:
		{
			"foo bar foo",
			"foo\nbar\nfoo",
			4,
			true,
			nil,
		},
		// A hyphen is a valid breakpoint:
		{
			"foo-foobar",
			"foo-\nfoobar",
			4,
			true,
			nil,
		},
		// Space buffer needs to be emptied before breakpoints:
		{
			"foo --bar",
			"foo --bar",
			9,
			true,
			nil,
		},
		// Lines are broken at whitespace, even if words
		// are too long. We do not break words:
		{
			"foo bars foobars",
			"foo\nbars\nfoobars",
			4,
			true,
			nil,
		},
		// A word that would run beyond the limit is wrapped:
		{
			"foo bar",
			"foo\nbar",
			5,
			true,
			nil,
		},
		// Whitespace that trails a line and fits the width
		// passes through, as does whitespace prefixing an
		// explicit line break. A tab counts as one character:
		{
			"foo\nb\t a\n bar",
			"foo\nb\t a\n bar",
			4,
			true,
			nil,
		},
		// Trailing whitespace is removed if it doesn't fit the width.
		// Runs of whitespace on which a line is broken are removed:
		{
			"foo    \nb   ar   ",
			"foo\nb\nar",
			4,
			true,
			nil,
		},
		// An explicit line break at the end of the input is preserved:
		{
			"foo bar foo\n",
			"foo\nbar\nfoo\n",
			4,
			true,
			nil,
		},
		// Explicit break are always preserved:
		{
			"\nfoo bar\n\n\nfoo\n",
			"\nfoo\nbar\n\n\nfoo\n",
			4,
			true,
			nil,
		},
		// Unless we ask them to be ignored:
		{
			"\nfoo bar\n\n\nfoo\n",
			"foo\nbar\nfoo",
			4,
			false,
			nil,
		},
		// Complete example:
		{
			" This is a list: \n\n\t* foo\n\t* bar\n\n\n\t* foo  \nbar    ",
			" This\nis a\nlist: \n\n\t* foo\n\t* bar\n\n\n\t* foo\nbar",
			6,
			true,
			nil,
		},
		// ANSI sequence codes don't affect length calculation:
		{
			"\x1B[38;2;249;38;114mfoo\x1B[0m\x1B[38;2;248;248;242m \x1B[0m\x1B[38;2;230;219;116mbar\x1B[0m",
			"\x1B[38;2;249;38;114mfoo\x1B[0m\x1B[38;2;248;248;242m \x1B[0m\x1B[38;2;230;219;116mbar\x1B[0m",
			7,
			true,
			nil,
		},
		// ANSI control codes don't get wrapped:
		{
			"\x1B[38;2;249;38;114m(\x1B[0m\x1B[38;2;248;248;242mjust another test\x1B[38;2;249;38;114m)\x1B[0m",
			"\x1B[38;2;249;38;114m(\x1B[0m\x1B[38;2;248;248;242mjust\nanother\ntest\x1B[38;2;249;38;114m)\x1B[0m",
			3,
			true,
			nil,
		},
		// Complete example with custom newline:
		{
			"* this is a long list item that gets wrapped.",
			"* this is a long\n   list item that\n   gets wrapped.",
			17,
			true,
			func(w io.Writer) int { w.Write([]byte("\n   ")); return 2 },
		},
	}

	for i, tc := range tt {
		f := NewWriter(tc.Limit)
		f.KeepNewlines = tc.KeepNewlines
		f.NewLineFunc = tc.NewLineFunc

		_, err := f.Write([]byte(tc.Input))
		if err != nil {
			t.Error(err)
		}
		f.Close()

		if f.String() != tc.Expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n\nActual Output:\n\n`%s`", i, tc.Expected, f.String())
		}
	}
}

func TestWordWrapString(t *testing.T) {
	actual := String("foo bar", 3)
	expected := "foo\nbar"
	if actual != expected {
		t.Errorf("expected:\n\n`%s`\n\nActual Output:\n\n`%s`", expected, actual)
	}
}
