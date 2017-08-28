// Package util contains all the util functions to be used across the project
package util

import (
    "reflect"
    "regexp"
    "unicode/utf8"
    "fmt"
    "math/rand"
    "time"
)

// IsAlphanumeric is used to check if the given string is alphanumeric
var IsAlphanumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

// IsString is used to check if the given interface{} is a string
func IsString(value interface{}) bool {
    if (reflect.ValueOf(value).Kind() == reflect.String) {
        return true
    }
    return false
}

// IsFloat64 is used to check if the given interface{} is a float64. By default the number inputs from JSON
// are mapped to float64, which need type assertion afterwards.
func IsFloat64(value interface{}) bool {
    if (reflect.ValueOf(value).Kind() == reflect.Float64) {
        return true
    }
    return false
}

// IsUnicode is used to check if the given string in unicode encoded.
func IsUnicode(value string) bool {
    return len(value) != utf8.RuneCountInString(value)
}

// GetRandom2DigitHex is used to generate a 2 digit random hexadecimal is XX format
func GetRandom2DigitHex() string {
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    n := r1.Intn(256)
    return ConvertIntToHex(n)
}

// ConvertIntToHex is used to convert an integer to XX format hexadecimal
func ConvertIntToHex(i int) string {
    return fmt.Sprintf("%02X", i)
}