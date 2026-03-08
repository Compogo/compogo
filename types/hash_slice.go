package types

import (
	"errors"
	"fmt"
	"iter"

	"github.com/Compogo/compogo/tools"
)

const NoneIndex = -1

var (
	AlreadyExistsError = errors.New("already exists")
)

// HashSlice combines an ordered slice with a hash map for O(1) lookups by key.
// It maintains insertion order while providing fast index retrieval and
// element uniqueness.
//
// This is useful when you need to:
//   - Keep elements in a stable order (like insertion order)
//   - Quickly find an element's position (IndexOf)
//   - Check element existence (Contains)
//   - Ensure uniqueness of elements
//   - Iterate in order (All)
//
// Example:
//
//	hs, _ := NewHashSlice(StatusPending, StatusActive, StatusClosed)
//
//	idx := hs.IndexOf(StatusActive)      // returns 1 (O(1))
//	item, _ := hs.GetByIndex(2)          // returns StatusClosed
//
//	for idx, status := range hs.All() {
//	    fmt.Printf("%d: %v\n", idx, status)
//	}
//
// Trying to add duplicate returns AlreadyExistsError
//
//	_, err := hs.Add(StatusActive)
type HashSlice[T LinkKey] struct {
	linker    *Linker[T, int]
	items     []T
	typeName  string
	zeroValue T
}

func NewHashSlice[T LinkKey](items ...T) (*HashSlice[T], error) {
	hs := &HashSlice[T]{
		linker:    NewLinker[T, int](),
		items:     make([]T, 0, len(items)),
		typeName:  fmt.Sprintf("HashSlice[%s]", tools.TypeName[T]()),
		zeroValue: *new(T),
	}

	var err error
	for _, item := range items {
		if _, err = hs.Add(item); err != nil {
			return nil, err
		}
	}

	return hs, nil
}

func (hs *HashSlice[T]) Len() int {
	return len(hs.items)
}

func (hs *HashSlice[T]) Add(item T) (int, error) {
	if hs.linker.HasByKey(item) {
		return 0, fmt.Errorf("%s item %v %w", hs.typeName, item, AlreadyExistsError)
	}

	hs.items = append(hs.items, item)
	index := len(hs.items) - 1

	hs.linker.Add(item, index)

	return index, nil
}

func (hs *HashSlice[T]) GetByIndex(index int) (T, error) {
	if hs.Len()-1 < index {
		return hs.zeroValue, fmt.Errorf("%s index %d out of range", hs.typeName, index)
	}

	return hs.items[index], nil
}

func (hs *HashSlice[T]) Items() []T {
	return hs.items
}

func (hs *HashSlice[T]) IndexOf(item T) int {
	if !hs.linker.HasByKey(item) {
		return NoneIndex
	}

	index, _ := hs.linker.Get(item)
	return index
}

func (hs *HashSlice[T]) Remove(item T) {
	index, _ := hs.linker.Get(item)
	if index != NoneIndex {
		hs.RemoveByIndex(index)
	}
}

func (hs *HashSlice[T]) RemoveByIndex(index int) {
	key, err := hs.GetByIndex(index)
	if err != nil {
		return
	}

	hs.linker.Remove(key)

	hs.items = append(hs.items[:index], hs.items[index+1:]...)

	for i := index; i < hs.Len(); i++ {
		hs.linker.Add(hs.items[i], i)
	}
}

func (hs *HashSlice[T]) Contains(item T) bool {
	return hs.linker.HasByKey(item)
}

func (hs *HashSlice[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, item := range hs.items {
			if !yield(i, item) {
				return
			}
		}
	}
}
