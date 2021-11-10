package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordWrap(t *testing.T) {
	type TestCase struct {
		InputIndent    int
		InputWidth     int
		InputString    string
		ExpectedOutput string
	}
	testcases := map[string]TestCase{
		"basic": {
			InputWidth:     10,
			InputString:    "foo bar baz",
			ExpectedOutput: "foo bar\nbaz",
		},
		"unwrap": {
			InputWidth:     80,
			InputString:    "foo\nbar\nbaz\n",
			ExpectedOutput: "foo bar baz",
		},
		"paragraphs": {
			InputWidth:     80,
			InputString:    "foo\n\nbar\n\nbaz",
			ExpectedOutput: "foo\n\nbar\n\nbaz",
		},
		"sentence": {
			InputWidth:     80,
			InputString:    " Sentence  one.    Sentence  two. ",
			ExpectedOutput: "Sentence one.  Sentence two.",
		},
		"sentences": {
			InputWidth:     80,
			InputString:    "Sentence one.\nSentence two.",
			ExpectedOutput: "Sentence one.  Sentence two.",
		},
		"nbsp": {
			InputWidth:     40,
			InputString:    "The style of dress is not the key Mr. Dobalina.",
			ExpectedOutput: "The style of dress is not the key\nMr. Dobalina.",
		},
		"nbsp-counter": {
			InputWidth:     45,
			InputString:    "The style of dress is not the key Mr.  Dobalina.",
			ExpectedOutput: "The style of dress is not the key Mr.\nDobalina.",
		},
		"longword": {
			InputWidth:     4,
			InputString:    "averylongword",
			ExpectedOutput: "averylongword",
		},
		"empty": {
			InputIndent:    4,
			InputWidth:     80,
			InputString:    "",
			ExpectedOutput: "",
		},
	}
	for tcName, tcData := range testcases {
		tcData := tcData
		t.Run(tcName, func(t *testing.T) {
			actualOutput := wordwrap(
				tcData.InputIndent,
				tcData.InputWidth,
				tcData.InputString)
			assert.Equal(t, tcData.ExpectedOutput, actualOutput)
		})
	}
}
