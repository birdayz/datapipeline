package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var outputPathTests = []struct {
	inputPath  string
	outputPath string
	extension  string
}{

	{inputPath: "test.csv", extension: "nfo", outputPath: "test.nfo"},
	{inputPath: "test..csv", extension: "nfo", outputPath: "test..nfo"},
}

func TestOutputFilename(t *testing.T) {
	for _, test := range outputPathTests {
		outputPath := getOutputFilename(test.inputPath, "nfo")
		assert.Equal(t, test.outputPath, outputPath)
	}
}
