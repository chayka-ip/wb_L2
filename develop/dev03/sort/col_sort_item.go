package sort

import "golang.org/x/exp/constraints"

type columnSortItem[T constraints.Ordered] struct {
	data  T
	index int
}

func newColumnSortItem[T constraints.Ordered](data T, index int) columnSortItem[T] {
	return columnSortItem[T]{
		index: index,
		data:  data,
	}
}
