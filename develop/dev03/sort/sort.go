package sort

import (
	"fmt"
	"os"
	"sort"

	"golang.org/x/exp/constraints"
)

//ExecuteCLI executes sort command
func ExecuteCLI(args []string) int {
	opt, err := newOptions(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	if err := run(opt); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func run(opt *options) error {
	outData := make([]string, 0)
	data, err := readData(opt.reader)
	if err != nil {
		return err
	}

	if opt.isUniqueOnly {
		data = removeDuplicates(data)
	}

	if opt.sortColEnabled {
		outData = doColumnSort(opt, data)
	} else {
		isDesc := opt.isDescOrder
		switch opt.sortType {
		case sortTypeDefautl:
			outData = doDefaultStringSort(data, isDesc)
		case sortTypeNumeric:
			outData = doStringNumericSort(data, isDesc)
		}
	}

	return writeData(outData, opt.outFilePath)
}

func doColumnSort(opt *options, data []string) []string {
	sep := opt.separator
	colIndex := opt.sortColIndex
	sortType := opt.sortType
	isDescending := opt.isDescOrder

	src := append([]string{}, data...)
	def := []string{}
	colStr := []columnSortItem[string]{}
	colInt := []columnSortItem[int]{}

	for i, v := range data {
		sp := splitStringAndCleanUp(v, sep)
		if colIndex >= len(sp) {
			def = append(def, v)
			continue
		}

		chunk := sp[colIndex]
		switch sortType {
		case sortTypeDefautl:
			colStr = append(colStr, newColumnSortItem(chunk, i))
		case sortTypeNumeric:
			if startsWithDigit(chunk) {
				d, _ := getNumberFromStringStart(chunk)
				colInt = append(colInt, newColumnSortItem(d, i))
				continue
			}
			def = append(def, v)
		default:
			panic("unknown sort type")
		}
	}

	switch sortType {
	case sortTypeDefautl:
		return sortAndMergeDefaultAndCol(src, def, colStr, isDescending)
	case sortTypeNumeric:
		return sortAndMergeDefaultAndCol(src, def, colInt, isDescending)
	}
	panic("unknown sort type")
}

func doDefaultStringSort(data []string, isDescending bool) []string {
	less := func(i, j int) bool { return data[i] < data[j] }
	if isDescending {
		sort.Slice(data, reverse(less))
		return data
	}
	sort.Slice(data, less)
	return data
}

func doStringNumericSort(data []string, isDescending bool) []string {
	src := append([]string{}, data...)
	noNumeric := []string{}
	numeric := []columnSortItem[int]{}

	for i, v := range data {
		if startsWithDigit(v) {
			d, _ := getNumberFromStringStart(v)
			numeric = append(numeric, newColumnSortItem(d, i))
			continue
		}
		noNumeric = append(noNumeric, v)
	}

	return sortAndMergeDefaultAndCol(src, noNumeric, numeric, isDescending)
}

func getSourceStringFromColumnItems[T constraints.Ordered](src []string, items []columnSortItem[T]) []string {
	out := make([]string, len(items))
	for i, v := range items {
		out[i] = src[v.index]
	}
	return out
}

func sortCol[T constraints.Ordered](data []columnSortItem[T], isDescending bool) {
	less := func(i, j int) bool { return data[i].data < data[j].data }
	if isDescending {
		sort.Slice(data, reverse(less))
		return
	}
	sort.Slice(data, less)
}

func sortAndMergeDefaultAndCol[T constraints.Ordered](src, def []string, col []columnSortItem[T], isDescending bool) []string {
	doDefaultStringSort(def, isDescending)
	sortCol(col, isDescending)

	var a1 = def
	var a2 = getSourceStringFromColumnItems(src, col)

	if isDescending {
		a1, a2 = a2, a1
	}
	return append(a1, a2...)
}

func reverse(less func(i, j int) bool) func(i, j int) bool {
	return func(i, j int) bool {
		return !less(i, j)
	}
}
