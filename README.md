## Datapipeline
![Build Status](https://codeship.com/projects/23b8cf00-0e76-0135-7de1-5ac1425fdf96/status?branch=master)

This is a small tool to convert documents from various source formats to
different output formats. Currently, CSV is supported for input and YAML, Ascii
Text Tables and JSON are supported for output.
There are several filters to process input data before writing the output file:
- Removal of duplicates
- Sorting
- Validation

The tool can read from standard input and write to standard output,
but using files is also possible. See Usage for details.
## Performance
The decision towards channels has a cost: each pipeline stage costs around 300ns
(on my I5-3570K CPU) of CPU time per Entry (CSV row). This is ~3,3 million
entries per second. However, unless very rigid performance constraints are
given, i feel that this cost is justified. With channels, multithreaded
processing may be added easily: either multiple independent pipelines are used
concurrently (multiple "documents"), or specific pipeline stages/filters are
introduced that use multiple goroutines. Other stages of the pipeline do not
have to be touched in this case.

## Building
```
go get github.com/birdayz/datapipeline
```
### build (default)
builds the actual commandline tool. `make` or `make build`
### clean
do some cleanup `make clean`
### test
run tests `make test`
### bench
run benchmarks `make bench`
### lint
run linters (golint, go vet, ineffassign) `make lint`

Tested platform: Ubuntu Linux 16.04 amd64, Go 1.7.4 linux/amd64
Dependencies are vendored, no need to download anything.

To build, run `make`. The binary `csvconverter` is written to the current directory.

## Usage
```
Usage:  ./csvconverter [OPTION]... INPUTFILE [OUTPUTFILE]
If INPUTFILE is -, read from standard input.
If OUTPUTFILE is -, write to standard output.
If no OUTPUTFILE is given, OUTPUTFILE will be placed in the folder of the INPUTFILE,
with its corresponding file suffix.
  -dedupField string
        Remove duplicates of given field (optional).
  -inputFormat string
        Input format. Eligible values: [csv]. (default "csv")
  -outputFormat string
        Output format. Eligible values: [txt,json,yaml].
  -sortField string
        Field name which is used for sorting. (optional)
  -sortOrder string
        Sort order. Eligible values: [asc, desc]. (optional) (default "asc")
```
### Examples

`./csvconverter testdata/sample_data.csv` converts the csv file to a fancy ascii table. The output file is written to the folder of the input file.

`./csvconverter -outputFormat=json testdata/sample_data.csv` converts the json

`./csvconverter testdata/sample_data.csv testdata/sample_data.json` does the same as above. the tool can infer the output format from the file extension, if an output path is given.

csvconverter can read from stdin and write to stdout:

`cat testdata/sample_data.csv | ./csvconverter -` runs the csvconverter (with the default output format: ascii table) and prints output to stdout.

filter duplicates of a specific field:

`./csvconverter -dedupField=name testdata/sample_data.csv` removes duplicate rows with identical values in the "name" column. Rows occuring first have priority.

sort output according to a specific field (column):

`./csvconverter -sortField=name testdata/sample_data.csv` sorts the rows by the name column/field. The order may be changed with the -sortOrder flag.
