package encoder

import (
	"errors"
	"math"
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length   = len(alphabet)
)

func GetRandomId(l int) string {
	rand.Seed(time.Now().UnixNano())
	randy := ""
	for i := 0; i < l; i++ {
		u := alphabet[rand.Intn(length)]
		randy += string(u)
	}
	return randy
}

func Encode(number uint64) string {
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(11)

	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(number % uint64(length))])
	}

	return encodedBuilder.String()
}

func Decode(encoded string) (uint64, error) {
	var number uint64

	for i, symbol := range encoded {
		alphabeticPosition := strings.IndexRune(alphabet, symbol)

		if alphabeticPosition == -1 {
			return uint64(alphabeticPosition), errors.New("invalid character: " + string(symbol))
		}
		number += uint64(alphabeticPosition) * uint64(math.Pow(float64(length), float64(i)))
	}

	return number, nil
}
