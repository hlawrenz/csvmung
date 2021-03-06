package filters

import (
	"regexp"
)

type Filterer interface {
	Filter(inCh chan []string) chan []string
}

type ColFilterer struct {
	Columns []interface{}
}

func (f ColFilterer) Filter(inCh chan []string) chan []string {
	outCh := make(chan []string)
	go func() {
		for row := range inCh {
			var newRow []string
			var c interface{}
			for _, c = range f.Columns {
				switch c.(type) {
				case int:
					newRow = append(newRow, row[c.(int)])
				default:
					newRow = append(newRow, c.(string))
				}
			}
			outCh <- newRow
		}
		close(outCh)
	}()
	return outCh
}

type RegexFilterer struct {
	Col     int
	Pattern *regexp.Regexp
}

func (f RegexFilterer) Filter(inCh chan []string) chan []string {
	outCh := make(chan []string)
	go func() {
		for row := range inCh {
			if f.Pattern.MatchString(row[f.Col]) == true {
				outCh <- row
			}
		}
		close(outCh)
	}()
	return outCh
}

type SplitFilterer struct {
	Col     int
	Pattern *regexp.Regexp
}

func (f SplitFilterer) Filter(inCh chan []string) chan []string {
	outCh := make(chan []string)
	go func() {
		for row := range inCh {
			split := f.Pattern.Split(row[f.Col], -1)
			newRowLen := len(split) + len(row) - 1
			newRow := make([]string, f.Col, newRowLen)
			copy(newRow, row[0:f.Col])
			newRow = append(newRow, split...)
			newRow = append(newRow, row[f.Col+1:]...)
			outCh <- newRow
		}
		close(outCh)
	}()
	return outCh
}

type UniqFilterer struct {
	Col int
}

func (f UniqFilterer) Filter(inCh chan []string) chan []string {
	outCh := make(chan []string)
	go func() {
		seen := make(map[string]int)
		for row := range inCh {
			_, present := seen[row[f.Col]]
			if !present {
				outCh <- row
				seen[row[f.Col]] = 1
			}
		}
		close(outCh)
	}()
	return outCh
}

