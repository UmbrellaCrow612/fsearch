package args

import (
	"fmt"
	"reflect"
)

func printArgsMapValues(args *ArgsMap) {
	// Dereference the pointer
	v := reflect.ValueOf(args).Elem()
	t := v.Type()

	fmt.Println("----- ArgsMap Values -----")

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.Kind() == reflect.Slice {
			fmt.Printf("%-15s | %-12s | %v\n", field.Name, value.Type(), sliceToString(value))
			continue
		}

		fmt.Printf("%-15s | %-12s | %v\n", field.Name, value.Type(), value.Interface())
	}

	fmt.Println("---------------------------")
}

func sliceToString(v reflect.Value) string {
	if v.Len() == 0 {
		return "[]"
	}
	s := "["
	for i := 0; i < v.Len(); i++ {
		s += fmt.Sprintf("%v", v.Index(i))
		if i < v.Len()-1 {
			s += ", "
		}
	}
	s += "]"
	return s
}
