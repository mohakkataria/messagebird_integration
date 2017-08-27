package util

import (
    "reflect"
    "regexp"
    "unicode/utf8"
    "fmt"
)

var IsAlphanumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func IsString(value interface{}) bool {
    if (reflect.ValueOf(value).Kind() == reflect.String) {
        return true
    }
    return false
}

func IsFloat64(value interface{}) bool {
    fmt.Println(reflect.ValueOf(value).Kind())
    if (reflect.ValueOf(value).Kind() == reflect.Float64) {
        return true
    }
    return false
}

func IsUnicode(value string) bool {
    return len(value) != utf8.RuneCountInString(value)
}