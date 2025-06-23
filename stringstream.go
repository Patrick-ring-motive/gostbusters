package main

import "strings"

type StringStream []string

func NewStringStream[T string|[]string|Stream[string]](s T)StringStream{
  switch v:=any(s).(type){
  case string:
      return StringStream(strings.Split(v,""))
  case []string:
      return StringStream(v)
  case Stream[string]:
      return StringStream(v.Slice())
  }
  return any(s).(StringStream)
}

/* returns the underlying slice */
func (s StringStream) Slice() []string {
  return s
}

/* returns identity */
func (s StringStream) Stream() Stream[string] {
  return NewStream(s.Slice())
}

func (s StringStream) MapAny /*[X any]*/ (fna interface {
  /* func(t T) X */
},
  streamTypes ...AnyStream /* Stream[X] */) any {
  return s.Stream().MapAny(fna,streamTypes...)
}

func (s StringStream) FlatMapAny /*[X any]*/ (fna interface{/* Stream[X] */}) any {
  return s.Stream().FlatMapAny(fna)
}