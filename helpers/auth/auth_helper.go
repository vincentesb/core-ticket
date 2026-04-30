package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"math"
)

func GenerateRandomKeyBase64(maxLength int) (string, error) {
	if maxLength <= 0 {
		return "", errors.New("MaxLength must be greater than 0")
	}

	// Calculate the number of bytes needed to achieve the desired length
	numBytes := int(math.Ceil(float64(maxLength) * 3.0 / 4.0))

	// Use crypto/rand for secure random bytes
	key := make([]byte, numBytes)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to Base64
	encodedKey := base64.RawURLEncoding.EncodeToString(key)

	// Ensure the generated string does not exceed the maxLength
	if len(encodedKey) > maxLength {
		encodedKey = encodedKey[:maxLength]
	}

	return encodedKey, nil
}
