package util

import (
    "reflect"
    "regexp"
    "unicode/utf8"
    "fmt"
    "math/rand"
    "time"
)

var IsAlphanumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func IsString(value interface{}) bool {
    if (reflect.ValueOf(value).Kind() == reflect.String) {
        return true
    }
    return false
}

func IsFloat64(value interface{}) bool {
    if (reflect.ValueOf(value).Kind() == reflect.Float64) {
        return true
    }
    return false
}

func IsUnicode(value string) bool {
    return len(value) != utf8.RuneCountInString(value)
}

func GetRandom2DigitHex() string {
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    n := r1.Intn(256)
    return ConvertIntToHex(n)
}

func ConvertIntToHex(i int) string {
    return fmt.Sprintf("%02X", i)
}