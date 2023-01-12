package incapsula

import (
	"math/rand"
)

func RandomString(letters []rune, n int) string {
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func RandomLetterAndNumberString(n int) string {
	return RandomString([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"), n)
}

func RandomLowLetterAndNumberString(n int) string {
	return RandomString([]rune("abcdefghijklmnopqrstuvwxyz0123456789"), n)
}

func RandomCapitalLetterAndNumberString(n int) string {
	return RandomString([]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"), n)
}

func RandomLowLetterString(n int) string {
	return RandomString([]rune("abcdefghijklmnopqrstuvwxyz"), n)
}

func RandomCapitalLetterString(n int) string {
	return RandomString([]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), n)
}

func RandomNumbersExcludingZeroString(n int) string {
	return RandomString([]rune("123456789"), n)
}
