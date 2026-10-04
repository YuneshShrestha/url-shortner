package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"os"
)

/*
	Keyed permutation using a Feistel network.

	Split the ID into two halves (left, right). Each round:

		left, right = right, left XOR F(key, round, right)

	F is HMAC-SHA256, so without the key the output looks random.
	The trick: XOR undoes itself, so running the rounds backwards
	recovers the original ID, even though F itself is one-way.

		with key "demo":
		100 → "CHIEFc"
		101 → "iEJUafQ"
		102 → "iirRkSP"

	Every ID maps to exactly one code (it's a permutation), so there
	are no collisions, and decoding is O(1): no extra DB lookup.
*/

// 40 bits gives us about 1.1 trillion possible IDs.
// Every 40-bit value fits into at most 7 Base62 characters.
// So Redis can store up to 1.1 trillion short URLs.
const (
	idBits        = 40
	halfBits      = idBits / 2
	halfMask      = 1<<halfBits - 1
	maxID         = 1<<idBits - 1
	feistelRounds = 10 // same round count as NIST FF1
)

var obfuscationKey = []byte(os.Getenv("URL_SHORTENER_SECRET"))

/*
Feistel Round:

	  SECRET KEY
	      │
	      ▼
	┌───────────┐

round ────►│           │
right ────►│ HMAC-SHA256

	│           │
	└─────┬─────┘
	      │
	      ▼
	huge random-looking
	    number
	      │
	      ▼
	 keep 20 bits
*/
func feistelRound(round int, half uint64) uint64 {
	mac := hmac.New(sha256.New, obfuscationKey)
	var buf [9]byte
	buf[0] = byte(round)
	binary.BigEndian.PutUint64(buf[1:], half)
	mac.Write(buf[:])
	return binary.BigEndian.Uint64(mac.Sum(nil)) & halfMask
}

/*
Why XOR?
Because XOR can undo itself.
Remember:
A XOR B XOR B = A

For example:
1011 XOR 1001 = 0010

Then:
0010 XOR 1001 = 1011

We recovered the original value.
That's why we don't need to reverse HMAC.
*/
func Obfuscate(id uint64) (uint64, error) {
	if id > maxID {
		return 0, fmt.Errorf("id %d exceeds %d-bit range", id, idBits)
	}
	left, right := id>>halfBits, id&halfMask
	for r := 0; r < feistelRounds; r++ {
		left, right = right, left^feistelRound(r, right)
	}
	return left<<halfBits | right, nil
}

func Deobfuscate(value uint64) (uint64, error) {
	if value > maxID {
		return 0, fmt.Errorf("value %d exceeds %d-bit range", value, idBits)
	}
	left, right := value>>halfBits, value&halfMask
	for r := feistelRounds - 1; r >= 0; r-- {
		left, right = right^feistelRound(r, left), left
	}
	return left<<halfBits | right, nil
}

// EncodeID: internal ID → public short code.
func EncodeID(id uint64) (string, error) {
	obfuscated, err := Obfuscate(id)
	if err != nil {
		return "", err
	}
	return EncodeBase62(obfuscated), nil
}

// DecodeID: public short code → internal ID.
func DecodeID(code string) (uint64, error) {
	value, err := DecodeBase62(code)
	if err != nil {
		return 0, err
	}
	return Deobfuscate(value)
}
