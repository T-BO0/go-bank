package util

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

// RandomFloat generates a random float64 between min and max
func RandomFloat(min, max int64) float64 {
	return float64(float64(min)+float64(rand.Int64N(max-min+1))) * rand.Float64()
}

// RandomString generates a random string fo length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.IntN(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner will generate a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney will generate a random money
func RandomMoney() float64 {
	return RandomFloat(1000, 10000)
}

// RandomCurrency will generate a random currency
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "GEL"}

	n := len(currencies)

	return currencies[rand.IntN(n)]
}

func RandomEmail() string {
	domains := []string{"gmail", "yehu", "mail"}

	username := RandomString(7)
	domain := domains[rand.IntN(len(domains))]
	email := fmt.Sprintf("%s@%s.com", username, domain)

	return email
}
