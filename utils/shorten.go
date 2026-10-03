package utils

import (
	"fmt"
	"math"
)

const base62Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

/*
	Redis ID
		│
		▼
		12510
		│
		│ % 62 = 48
		▼
		W
		│
		│ / 62
		▼
		201
		│
		│ % 62 = 15
		▼
		p
		│
		│ / 62
		▼
		3
		│
		│ % 62 = 3
		▼
		d
		│
		│ / 62
		▼
		0
*/
func EncodeBase62(number uint64) string {
	if number == 0 {
		return string(base62Alphabet[0])
	}

	var encoded []byte
	for number > 0 {
		encoded = append(encoded, base62Alphabet[number%uint64(len(base62Alphabet))])
		number /= uint64(len(base62Alphabet))
	}

	// Reverse the encoded string Because % gives us the rightmost digit first.
	// Eg. encoded = ['W', 'p', 'd'] should become ['d', 'p', 'W']
	
	for left, right := 0, len(encoded)-1; left < right; left, right = left+1, right-1 {
		encoded[left], encoded[right] = encoded[right], encoded[left]
	}
	return string(encoded)
}

/*
	For Decode:

	Why number * 62 + digit?

	This is the most important concept to understand.
	It's the same principle as normal decimal numbers.
	Suppose you read:

	123

	Start:
	0

	Read 1:
	0 × 10 + 1 = 1

	Read 2:
	1 × 10 + 2 = 12

	Read 3:
	12 × 10 + 3 = 123

	For Base62, instead of multiplying by 10, we multiply by 62.

	dpW

	d = 3
	p = 15
	W = 48

	Therefore:

	0 × 62 + 3
	= 3

	3 × 62 + 15
	= 201

	201 × 62 + 48
	= 12510

	That's the whole decoding algorithm.
*/
func DecodeBase62(value string) (uint64, error) {
	if value == "" {
		return 0, fmt.Errorf("Base62 value cannot be empty")
	}

	var number uint64
	for _, character := range value {
		index := -1
		for i, candidate := range base62Alphabet {
			if character == candidate {
				index = i
				break
			}
		}
		if index == -1 {
			return 0, fmt.Errorf("invalid Base62 character: %q", character)
		}

		digit := uint64(index)
		if number > (math.MaxUint64-digit)/uint64(len(base62Alphabet)) {
			return 0, fmt.Errorf("Base62 value overflows uint64")
		}
		number = number*uint64(len(base62Alphabet)) + digit
	}
	return number, nil

}
