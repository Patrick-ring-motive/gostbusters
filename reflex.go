package main

import "reflect"

func MakeSliceOf(value interface{}, len, cap int) any {
	switch v := value.(type) {
	case reflect.Type:
		return reflect.MakeSlice(reflect.SliceOf(v),len,cap).Interface()
	case reflect.Value:
		return reflect.MakeSlice(reflect.SliceOf(v.Type()),len,cap).Interface()
	default:
		return reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(v)),len,cap).Interface()
	}

}
