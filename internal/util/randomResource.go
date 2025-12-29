package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	Alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// random Integer Generator. It generates a random integer b/w min and max.
func RandomInt(min, max int64) int64 {
	return (min + rand.Int63n(max-min+1))
}

// random string Generator. It generates a random string of "n" length.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(Alphanumeric)

	for i := 0; i < n; i++ {
		c := Alphanumeric[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(8))
}
func RandomName(n int)string{
	return RandomString(n)
}
func RandomUsername()string{
	return fmt.Sprintf("%s%d", RandomString(5),time.Now().UnixNano())	
}
