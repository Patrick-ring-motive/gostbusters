package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unsafe"
)

// Stream example
type Stream[T any] []T

var streamTypeMap = map[string]interface{}{
	"int":            Stream[int]{},
	"int8":           Stream[int8]{},
	"int16":          Stream[int16]{},
	"int32":          Stream[int32]{},
	"int64":          Stream[int64]{},
	"uint":           Stream[uint]{},
	"uint8":          Stream[uint8]{},
	"uint16":         Stream[uint16]{},
	"uint32":         Stream[uint32]{},
	"uint64":         Stream[uint64]{},
	"uintptr":        Stream[uintptr]{},
	"byte":           Stream[byte]{},
	"rune":           Stream[rune]{},
	"float32":        Stream[float32]{},
	"float64":        Stream[float64]{},
	"complex64":      Stream[complex64]{},
	"complex128":     Stream[complex128]{},
	"string":         Stream[string]{},
	"bool":           Stream[bool]{},
	"error":          Stream[error]{},
	"interface":      Stream[interface{}]{},
	"any":            Stream[any]{},
	"reflect.Value":  Stream[reflect.Value]{},
	"io.Reader":      Stream[io.Reader]{},
	"unsafe.Pointer": Stream[unsafe.Pointer]{},
}

/* Create a stream from an slice */
func NewStream[T any](s Stream[T]) Stream[T] {
	return s
}

/*Example showing how to use and chain stream methods */
func main() {

	x := []int{1, 2, 3, 8}

	streamTest := NewStream(x).Map(func(i int) int {
		return i + 1
	}).ToStrings().MapAny(strconv.Atoi, Stream[int]{}).(Stream[int]).Reverse().ToAny(Stream[int]{}).(Stream[int])

	fmt.Println(streamTest.ToString())

	fmt.Println(len(NewStringStream("cheese")))
}

/* Filters out elements that result in a false when passed through the predicate */
func (s Stream[T]) Filter(predicate func(T) bool) Stream[T] {
	var result []T
	for _, element := range s {
		if predicate(element) {
			result = append(result, element)
		}
	}
	return result
}

func (s Stream[T]) Reduce(fn func(T, T) T) T {
	var t T
	if len(s) == 0 {
		return t
	}
	t = s[0]
	for i := 1; i < len(s); i++ {
		t = fn(t, s[i])
	}
	return t
}

/* Filters out elements that match the passed in element */
func (s Stream[T]) Not(t T) Stream[T] {
	return s.Filter(func(i T) bool {
		return any(i) != any(t) && !reflect.DeepEqual(i, t)
	})
}

func (s Stream[T]) Unique() Stream[T] {
	seen := make(map[any]bool)
	result := []T{}
	for _, val := range s {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}

func (s Stream[T]) Append(t T) Stream[T] {
	return append(s, t)
}

func (s Stream[T]) AppendAll(t ...T) Stream[T] {
	return append(s.Slice(), t...)
}

func (s Stream[T]) AppendStream(t Stream[T]) Stream[T] {
	return append(s.Slice(), t...)
}

func (s Stream[T]) Reverse() Stream[T] {
	return s.Apply(slices.Reverse)
}

/* Sets an element at the specified index*/
func (s Stream[T]) Set(i int, t T) Stream[T] {
	arr := s.Slice()
	arr[i] = t
	return arr
}

/*  */
func (s Stream[T]) Get(i int) T {
	return s.Slice()[i]
}

func (s Stream[T]) SubStream(start, end int) Stream[T] {
	return s.Slice()[start:end]
}

func (s Stream[T]) First() T {
	return s.Slice()[0]
}

func (s Stream[T]) Last() T {
	return s.Slice()[len(s)-1]
}

func (s Stream[T]) Map(fn func(T) T) Stream[T] {
	result := make([]T, len(s))
	for i, element := range s {
		result[i] = fn(element)
	}
	return result
}

func (s Stream[T]) ForEach(fn func(T)) {

	for _, element := range s {
		fn(element)
	}
}

/* To stream existing in-place transforms directly */
func (s Stream[T]) Apply(fn func([]T)) Stream[T] {
	arr := s.Slice()
	fn(arr)
	return arr
}

/* To stream existing in-place transforms directly */
func (s Stream[T]) Join(delim ...string) string {
	d := ""
	if len(delim) > 0 {
		d = delim[0]
	}
	return strings.Join(s.ToStrings().Slice(), d)
}

func (s Stream[T]) FlatMap(fn func(T) []T) Stream[T] {
	result := make([]T, 0, len(s))
	for _, element := range s {
		result = append(result, fn(element)...)
	}
	return result
}

type AnyStream interface {
	MapAny(function interface{}, returnStreamType AnyStream) AnyStream
	FlatMapAny(function interface{}) AnyStream
	ToAny(returnStreamType AnyStream) AnyStream
	sealed()
}

/*
MapAny and FlatMapAny are a bit dangerous.

Because there are no generics in methods,
you can only do this generically with typeless interfaces and reflection
better to create the typed version for what you want.
See Stream.ToStrings() for example
*/
func (s Stream[T]) MapAny /*[X any]*/ (function interface {
	/* func(t T) X */
},
	returnStreamType AnyStream /* Stream[X] */) AnyStream {

	resultType := reflect.TypeOf(function).Out(0)
	size := len(s)
	result := reflect.MakeSlice(reflect.SliceOf(resultType), 0, size)
	fn := reflect.ValueOf(function)
	for _, element := range s {
		result = reflect.Append(result, fn.Call([]reflect.Value{reflect.ValueOf(element)})[0])
	}
	return result.Convert(reflect.TypeOf(returnStreamType)).Interface().(AnyStream)
}

func (s Stream[T]) FlatMapAny /*[X any]*/ (function interface { /* func(T)Stream[X] */
}) AnyStream {
	streamType := reflect.TypeOf(function).Out(0)
	resultType := streamType.Elem()
	size := len(s)
	result := reflect.MakeSlice(reflect.SliceOf(resultType), 0, size)
	fn := reflect.ValueOf(function)
	for _, element := range s {
		result = reflect.AppendSlice(result, fn.Call([]reflect.Value{reflect.ValueOf(element)})[0])
	}
	return result.Convert(streamType).Interface().(AnyStream)
}

func (s Stream[T]) ToAny /*[X any]*/ (
	returnStreamType AnyStream /* Stream[X] */) AnyStream {
	return reflect.ValueOf(s.Slice()).Convert(reflect.TypeOf(returnStreamType)).Interface().(AnyStream)
}

func (s Stream[T]) sealed() {}

/* returns the underlying slice */
func (s Stream[T]) Slice() []T {
	return s
}

/* returns identity */
func (s Stream[T]) Stream() Stream[T] {
	return s
}

/*Stringifies individual elements and returns them as a Stream of strings */
func (ss Stream[T]) ToStrings() Stream[string] {
	out := make([]string, len(ss))
	for i, s := range ss {
		str := ""
		bits, err := json.Marshal(s)
		if err != nil {
			str = fmt.Sprintf("%+v", s)
		} else {
			str = string(bits)
		}
		out[i] = strings.ReplaceAll(str, "\"", "")
	}
	return out
}

/* Stringifies the whole stream. Attempts with JSON and fallsback to verbose print */
func (s Stream[T]) ToString() string {
	str := ""
	bits, err := json.Marshal(s)
	if err != nil {
		str = fmt.Sprintf("%+v", s)
	} else {
		str = string(bits)
	}
	return strings.ReplaceAll(str, "\"", "")
}

/*
Sorts in ascending order
	numbers by value,strings alphabetically,
	everything else alphabetically by stringified value
*/

type ordered interface{ cmp.Ordered }

func (s Stream[T]) Sort() Stream[T] {

	switch v := any(s.Slice()).(type) {
	case []int:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []int8:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []int16:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []int32:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []int64:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []uint:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []uint8:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []uint16:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []uint32:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []uint64:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []uintptr:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []float32:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []float64:
		slices.Sort(v)
		return any(v).(Stream[T])
	case []string:
		slices.Sort(v)
		return any(v).(Stream[T])
	default:
		m := make(map[string]T)
		st := s.ToStrings()
		for i, r := range s {
			m[st[i]] = r
		}
		slices.Sort(st)
		out := make([]T, len(st))
		for i, r := range st {
			out[i] = m[r]
		}
		return out
	}
}
