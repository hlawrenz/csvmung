package filters_test

import (
	"testing"
	"github.com/hlawrenz/csvmung/filters"
	"regexp"
)

func slEqual(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}

	for i, v := range x {
		if v != y[i] {
			return false
		}
	}
	return true
}

func TestColFilterer(t *testing.T) {
	var input = [][]string {
		[]string{"a", "b", "c", "d"},
		[]string{"x", "y", "z", "splunge"},
	}
	var expected = [][]string {
		[]string{"b", "d", "x"},
		[]string{"y", "splunge", "x"},
	}
	inCh := make(chan []string)
	var cols []interface{}
	cols = append(cols, 1)
	cols = append(cols, 3)
	cols = append(cols, "x")
	colFilter := filters.ColFilterer{cols}
	outCh := colFilter.Filter(inCh)

	go func() {
		for _, v := range input {
			inCh <- v
		}
		close(inCh)
	}()

	i := 0
	for res := range outCh {
		if ! slEqual(res, expected[i]) {
			t.Errorf("ColFilterer: got %v, expected %v", res, expected[i])
		}

		i++
	}
}

func TestRegexFilterer(t *testing.T) {
	var input = [][]string {
		[]string{"a", "b", "c", "d"},
		[]string{"a", "ff", "c", "d"},
		[]string{"a", "b", "c", "d"},
		[]string{"a", "b", "c", "d"},
		[]string{"a", "b", "c", "d"},
		[]string{"a", "ff", "c", "d"},
		[]string{"x", "y", "z", "splunge"},
	}
	var expected = [][]string {
		[]string{"a", "ff", "c", "d"},
		[]string{"a", "ff", "c", "d"},
	}
	inCh := make(chan []string)
	reFilter := filters.RegexFilterer{1, regexp.MustCompile("^ff$")}
	outCh := reFilter.Filter(inCh)

	go func() {
		for _, v := range input {
			inCh <- v
		}
		close(inCh)
	}()

	i := 0
	for res := range outCh {
		if ! slEqual(res, expected[i]) {
			t.Errorf("ColFilterer: got %v, expected %v", res, expected[i])
		}

		i++
	}
}

func TestSplitFilterer(t *testing.T) {
	var input = [][]string {
		[]string{"a", "b-y", "c", "d"},
		[]string{"a", "x-k", "c", "d"},
	}
	var expected = [][]string {
		[]string{"a", "b", "y", "c", "d"},
		[]string{"a", "x", "k", "c", "d"},
	}
	inCh := make(chan []string)
	splitFilter := filters.SplitFilterer{1, regexp.MustCompile("-")}
	outCh := splitFilter.Filter(inCh)

	go func() {
		for _, v := range input {
			inCh <- v
		}
		close(inCh)
	}()

	i := 0
	for res := range outCh {
		if ! slEqual(res, expected[i]) {
			t.Errorf("ColFilterer: got %v, expected %v", res, expected[i])
		}

		i++
	}
}

func TestUniqFilterer(t *testing.T) {
	var input = [][]string {
		[]string{"a", "c", "d"},
		[]string{"x", "c", "d"},
		[]string{"a", "c", "d"},
		[]string{"a", "c", "d"},
		[]string{"j", "c", "d"},
		[]string{"a", "c", "d"},
		[]string{"j", "e", "n"},
	}
	var expected = [][]string {
		[]string{"a", "c", "d"},
		[]string{"x", "c", "d"},
		[]string{"j", "c", "d"},
	}
	inCh := make(chan []string)
	uniqFilter := filters.UniqFilterer{0}
	outCh := uniqFilter.Filter(inCh)

	go func() {
		for _, v := range input {
			inCh <- v
		}
		close(inCh)
	}()

	i := 0
	for res := range outCh {
		if ! slEqual(res, expected[i]) {
			t.Errorf("ColFilterer: got %v, expected %v", res, expected[i])
		}

		i++
	}
}


