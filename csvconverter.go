package main

import (
	"datapipeline/core"
	"datapipeline/core/filter"
	"datapipeline/core/marshaler"
	"datapipeline/core/unmarshaler"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Commandline flag values
var (
	paramSortField        string
	paramSortOrder        string
	paramDeduplicateField string
	paramInputFormat      string
	paramOutputFormat     string
)

const (
	formatText = "txt"
	formatJSON = "json"
	formatCSV  = "csv"
	formatYAML = "yaml"

	extensionCSV  = "csv"
	extensionText = "txt"
	extensionJSON = "json"
	extensionYAML = "yaml"
	extensionYML  = "yml"

	orderAsc  = "asc"
	orderDesc = "desc"
)

// Mappings for format -> extension and extension -> format. These mappings are
// necessary, because a format (eg YAML) does not necessarily map to one extension and
// vice-versa.

// Used to infer the output file extension if the format is provided via
// paramOutputFormat
var formatToExtension = map[string]string{
	formatText: extensionText,
	formatJSON: extensionJSON,
	formatCSV:  extensionCSV,
	formatYAML: extensionYAML, // We use .yaml as default extension for YAML: http://www.yaml.org/faq.html

}

// Used to infer the output format if only the output path is provided
var extensionToFormat = map[string]string{
	extensionText: formatText,
	extensionJSON: formatJSON,
	extensionCSV:  formatCSV,
	extensionYAML: formatYAML,
	extensionYML:  formatYAML,
}

var fields = []string{"name", "address", "stars", "contact", "phone", "uri"}

func init() {
	flag.StringVar(&paramSortField, "sortField", "", "Field name which is used for sorting. (optional)")
	flag.StringVar(&paramSortOrder, "sortOrder", "asc", "Sort order. Eligible values: [asc, desc]. (optional)")
	flag.StringVar(&paramDeduplicateField, "dedupField", "", "Remove duplicates of given field (optional).")
	flag.StringVar(&paramOutputFormat, "outputFormat", "", "Output format. Eligible values: [txt,json,yaml].")
	flag.StringVar(&paramInputFormat, "inputFormat", "csv", "Input format. Eligible values: [csv].")

	flag.Usage = func() {
		fmt.Println("Usage: ", os.Args[0], "[OPTION]... INPUTFILE [OUTPUTFILE]")
		fmt.Println("If INPUTFILE is -, read from standard input.")
		fmt.Println("If OUTPUTFILE is -, write to standard output.")
		fmt.Println("If no OUTPUTFILE is given, OUTPUTFILE will be placed in the folder of the INPUTFILE,\nwith its corresponding file suffix.")
		flag.PrintDefaults()
	}
}

// Construct the output filename by cutting off the extension (if it exists),
// and adding the output extension. inputFilename must not be empty.
func getOutputFilename(inputFilename, outputExtension string) string {
	lastDot := strings.LastIndex(inputFilename, ".")
	if lastDot != -1 {
		inputFilename = inputFilename[:lastDot]
	}
	return inputFilename + "." + outputExtension
}

func abortError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)

	// Cleanup if program does not terminate successfully
	if outputFile != nil {
		_ = outputFile.Close()
		_ = os.Remove(outputFile.Name())
	}
	os.Exit(1)
}

func getMarshaler(format string, writer io.Writer) (m core.Marshaler, success bool) {
	switch format {
	case formatJSON:
		return marshaler.NewJSON(writer), true
	case formatText:
		return marshaler.NewASCIITable(fields, writer), true
	case formatYAML:
		return marshaler.NewYAML(writer), true
	default:
		return nil, false
	}
}

var outputFile *os.File

func main() {
	flag.Parse()

	// Determine input/output files
	var reader io.Reader
	var writer io.Writer
	var outputFilePath string
	{
		switch flag.NArg() {
		case 1:
			inputFilePath := flag.Arg(0)
			if inputFilePath == "-" {
				reader = os.Stdin

				// Writer must be stdout, because we can't infer
				// a target file if there is no real input file
				// path
				writer = os.Stdout
			} else {
				inputFile, err := os.Open(inputFilePath)
				if err != nil {
					abortError(err)
				}
				defer inputFile.Close()
				reader = inputFile

				// Infer output path from input path
				dir, inputFilename := filepath.Split(inputFilePath)
				if inputFilename == "" {
					abortError(errors.New("invalid input filename"))
				}

				outputFileExtension, ok := formatToExtension[paramOutputFormat]
				if !ok {
					outputFileExtension = extensionText
				}

				outputFilePath = filepath.Join(dir, getOutputFilename(inputFilename, outputFileExtension))
				output, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
				if err != nil {
					abortError(err)
				}
				defer output.Close()
				writer = output
				outputFile = output
			}

		case 2:
			inputFilePath := flag.Arg(0)
			outputFilePath = flag.Arg(1)
			if inputFilePath == "-" {
				reader = os.Stdin
			} else {
				inputFile, err := os.Open(flag.Arg(0))
				if err != nil {
					abortError(err)
				}
				defer inputFile.Close()
				reader = inputFile
			}
			if outputFilePath == "-" {
				writer = os.Stdout
			} else {
				output, err := os.OpenFile(flag.Arg(1), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
				if err != nil {
					abortError(err)
				}
				defer outputFile.Close()
				writer = output
				outputFile = output
			}
		default:
			flag.Usage()
			os.Exit(1)
		}
	}

	// Determine marshaler/unmarshaler
	var u core.Unmarshaler
	switch paramInputFormat {
	case formatCSV:
		u = unmarshaler.NewCSV(reader)
	default:
		abortError(fmt.Errorf("Invalid input format: %v", paramInputFormat))
	}

	var m core.Marshaler
	if marshaler, ok := getMarshaler(paramOutputFormat, writer); ok {
		m = marshaler
	} else if outputFilePath != "" && outputFilePath != "-" {
		// Try to infer format from extension of output path
		ext := strings.TrimLeft(filepath.Ext(outputFilePath), ".")
		if inferredFormat, ok := extensionToFormat[ext]; ok {
			if marshaler, ok := getMarshaler(inferredFormat, writer); ok {
				m = marshaler
			}
		}
	}
	if m == nil {
		// Fallback to default - text
		m = marshaler.NewASCIITable(fields, writer)
	}

	// Construct pipeline
	var pipes []core.Pipe
	if paramDeduplicateField != "" {
		dedupFilter := filter.NewDuplicate(paramDeduplicateField)
		pipes = append(pipes, dedupFilter)
	}

	if paramSortField != "" {
		var order filter.SortOrder
		switch paramSortOrder {
		case orderAsc:
			order = filter.OrderAscending
		case orderDesc:
			order = filter.OrderDescending
		default:
			order = filter.OrderAscending
		}
		sortFilter := filter.NewSort(paramSortField, order)
		pipes = append(pipes, sortFilter)
	}

	pipeline := core.NewPipeline(pipes...)

	// Run pipeline
	go func() {
		err := u.Unmarshal(pipeline.In)
		if err != nil {
			abortError(err)
		}
	}()
	err := m.Marshal(pipeline.Out)
	if err != nil {
		abortError(err)
	}

	// Only inform user about output file if we didn't write to stdout
	if writer != os.Stdout {
		fmt.Fprintf(os.Stderr, "Wrote output to %v\n", outputFile.Name())
	}
}
