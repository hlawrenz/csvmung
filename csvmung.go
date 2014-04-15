package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"github.com/hlawrenz/csvmung/filters"
	"io"
	"os"
	"regexp"
	"strconv"
	"unicode/utf8"
)

var inputFn string
var outputFn string
var inputSep string
var outputSep string
var lazyQuotes bool
var strictLen bool

func init() {
	flag.StringVar(&inputFn, "i", "", "Input file. STDIN used if unspecified.")
	flag.StringVar(&outputFn, "o", "", "Output file. STDOUT used if unspecified.")
	flag.StringVar(&inputSep, "is", ",", "Input separator. Defaults to comma.")
	flag.StringVar(&outputSep, "os", ",", "Output separator. Defaults to comma.")
	flag.BoolVar(&lazyQuotes, "lazy-quote", false, "Allow lazy quotes in input.")
	flag.BoolVar(&strictLen, "strict-len", false, "Different length rows treates as an error.")
}

func readCsv(ch chan []string) {
	var reader *csv.Reader
	if inputFn == "" {
		reader = csv.NewReader(os.Stdin)
	} else {
		file, err := os.Open(inputFn)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = csv.NewReader(file)
	}
	if !strictLen {
		reader.FieldsPerRecord = -1
	}
	r, _ := utf8.DecodeRuneInString(inputSep)
	reader.Comma = r
	reader.LazyQuotes = lazyQuotes

	for {
		record, err := reader.Read()
		if err == io.EOF {
			close(ch)
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			close(ch)
			break
		}
		ch <- record
	}
}

func writeCsv(ch chan []string) {
	var writer *csv.Writer
	if outputFn == "" {
		writer = csv.NewWriter(os.Stdout)
	} else {
		file, err := os.Create(outputFn)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		defer file.Close()
		writer = csv.NewWriter(file)
	}

	r, _ := utf8.DecodeRuneInString(outputSep)
	writer.Comma = r

	for row := range ch {
		err := writer.Write(row)
		if err != nil {
			fmt.Println("Error:", err)
			close(ch)
			return
		}
	}
	writer.Flush()
}

func buildFilters() (f []filters.Filterer, err error) {
	var fs []filters.Filterer

	re := regexp.MustCompile(":")

	args := flag.Args()
	for a := range args {
		filterArg := re.Split(args[a], -1)
		filterName := filterArg[0]
		switch filterName {
		case "re":
			column, err := strconv.Atoi(filterArg[1])
			if err != nil {
				return nil, errors.New("Bad column argument " + filterArg[1])
			}
			newFilter := filters.RegexFilterer{column, regexp.MustCompile(filterArg[2])}
			fs = append(fs, newFilter)
		case "split":
			column, err := strconv.Atoi(filterArg[1])
			if err != nil {
				return nil, errors.New("Bad column argument " + filterArg[1])
			}
			newFilter := filters.SplitFilterer{column, regexp.MustCompile(filterArg[2])}
			fs = append(fs, newFilter)
		case "uniq":
			column, err := strconv.Atoi(filterArg[1])
			if err != nil {
				return nil, errors.New("Bad column argument " + filterArg[1])
			}
			newFilter := filters.UniqFilterer{column}
			fs = append(fs, newFilter)
		case "cols":
			var cols []interface{}
			for _, col := range filterArg[1:] {
				v, err := strconv.Atoi(col)
				if err != nil {
					cols = append(cols, col)
				} else {
					cols = append(cols, v)
				}
			}
			newFilter := filters.ColFilterer{cols}
			fs = append(fs, newFilter)
		default:
			fmt.Println("Unknown filter type")
		}
	}

	return fs, nil
}

func main() {
	flag.Parse()

	inCh := make(chan []string)
	go readCsv(inCh)

	filters, err := buildFilters()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(filters) < 1 {
		fmt.Println("You must specify at least one filter")
		os.Exit(1)
	}

	var outCh chan []string
	for filter := range filters {
		outCh = filters[filter].Filter(inCh)
		inCh = outCh
	}

	writeCsv(outCh)
}
