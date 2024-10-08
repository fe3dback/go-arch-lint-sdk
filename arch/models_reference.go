package arch

type (
	// Ref is generic structure.
	// This struct contains any value and Reference
	Ref[T any] struct {
		Value T
		Ref   Reference
	}

	// Reference contains pointer to specific file/line/column in file system
	// It can be used in errors / notices / linter warnings / etc
	Reference struct {
		File   PathAbsolute `json:"File"`
		Line   int          `json:"Line"`
		Column int          `json:"Column"`
		XPath  string       `json:"-"`
		Valid  bool         `json:"Valid"`
	}
)

type RefSlice[T comparable] []Ref[T]

func (rs RefSlice[T]) Values() []T {
	list := make([]T, 0, len(rs))

	for _, refValue := range rs {
		list = append(list, refValue.Value)
	}

	return list
}

func (rs RefSlice[T]) Contains(ref Ref[T]) bool {
	for _, refValue := range rs {
		if refValue.Value == ref.Value {
			return true
		}
	}

	return false
}

type RefMap[K comparable, V any] struct {
	values map[K]V
	refs   map[K]Reference
}

func (rf *RefMap[K, V]) Len() int {
	return len(rf.values)
}

func (rf *RefMap[K, V]) Set(key K, val V, ref Reference) {
	rf.values[key] = val
	rf.refs[key] = ref
}

func (rf *RefMap[K, V]) Get(key K) (V, Reference, bool) {
	value, hasValue := rf.values[key]
	ref, hasRef := rf.refs[key]

	return value, ref, hasValue && hasRef
}

func (rf *RefMap[K, V]) Has(key K) bool {
	_, hasValue := rf.values[key]
	return hasValue
}

func (rf *RefMap[K, V]) Each(fn func(K, V, Reference)) {
	for k, v := range rf.values {
		ref, exist := rf.refs[k]
		if !exist {
			continue
		}

		fn(k, v, ref)
	}
}

func NewRef[T any](value T, ref Reference) Ref[T] {
	return Ref[T]{
		Value: value,
		Ref:   ref,
	}
}

func NewInvalidRef[T any](value T) Ref[T] {
	return Ref[T]{
		Value: value,
		Ref:   NewInvalidReference(),
	}
}

func NewReference(file PathAbsolute, line int, column int, xpath string) Reference {
	return Reference{
		File:   file,
		Line:   line,
		Column: column,
		XPath:  xpath,
		Valid:  true,
	}
}

func NewFileReference(file PathAbsolute) Reference {
	return Reference{
		File:   file,
		Line:   1,
		Column: 0,
		Valid:  true,
	}
}

func NewInvalidReference() Reference {
	return Reference{
		Valid: false,
	}
}

func NewRefMap[K comparable, V any](size int) RefMap[K, V] {
	return RefMap[K, V]{
		values: make(map[K]V, size),
		refs:   make(map[K]Reference, size),
	}
}

func NewRefMapFrom[K comparable, V any](in map[K]Ref[V]) RefMap[K, V] {
	refMap := NewRefMap[K, V](len(in))

	for k, refValue := range in {
		refMap.Set(k, refValue.Value, refValue.Ref)
	}

	return refMap
}
