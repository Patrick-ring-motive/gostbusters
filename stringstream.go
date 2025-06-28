package main

import "strings"

type String string
type StringStream []String

func (s String) Str() string {
	return string(s)
}

func (s String) DoSplit1(fn func(string) []string) StringStream {
	return stringsToStrings(fn(string(s)))
}

func (s String) DoSplit2(fn func(string, string) []string, str String) StringStream {
	return stringsToStrings(fn(string(s), string(str)))
}

func (s String) Split(delim String) StringStream {
	return stringsToStrings(strings.Split(string(s), string(delim)))
}

func (s String) SplitAfter(delim String) StringStream {
	return stringsToStrings(strings.SplitAfter(string(s), string(delim)))
}

func stringsToStrings(s []string) []String {
	out := make([]String, len(s))
	for i, r := range s {
		out[i] = String(r)
	}
	return out
}

func StringsTostrings(s []String) []string {
	out := make([]string, len(s))
	for i, r := range s {
		out[i] = string(r)
	}
	return out
}

func NewStringStream[T String | []String | Stream[String] | string | []string | Stream[string]](s T) StringStream {
	switch v := any(s).(type) {
	case string:
		return StringStream(stringsToStrings(strings.Split(v, "")))
	case String:
		return StringStream(stringsToStrings(strings.Split(string(v), "")))
	case []string:
		return StringStream(stringsToStrings(v))
	case Stream[string]:
		return StringStream(stringsToStrings(v.Slice()))
	case Stream[String]:
		return StringStream((v.Slice()))
	case []String:
		return StringStream(v)
	}
	return any(s).(StringStream)
}

/* returns the underlying slice */
func (s StringStream) Slice() []String {
	return s
}

func (s StringStream) StrSlice() []string {
	return StringsTostrings(s)
}

/* returns identity */
func (s StringStream) Stream() Stream[String] {
	return NewStream(s.Slice())
}

func (s StringStream) StrStream() Stream[string] {
	return NewStream(StringsTostrings(s.Slice()))
}

func (s StringStream) MapAny /*[X any]*/ (function interface {
	/* func(t T) X */
},
	returnStreamType AnyStream /* Stream[X] */) AnyStream {
	return s.Stream().MapAny(function, returnStreamType)
}

func (s StringStream) FlatMapAny /*[X any]*/ (function interface { /* Stream[X] */
}) AnyStream {
	return s.Stream().FlatMapAny(function)
}

func (s StringStream) ToAny /*[X any]*/ (
	returnStreamType AnyStream /* Stream[X] */) AnyStream {
	return s.Stream().ToAny(returnStreamType)
}

func (s StringStream) sealed() {}
