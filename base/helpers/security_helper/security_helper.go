package security_helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/hkdf"
)

/*
GeneratePasswordHash generates a bcrypt hash for the given password using bcrypt's MinCost parameter.

Parameters:
- password (string): The password to generate the hash for.

Returns:
- string: The bcrypt hash generated for the password.

Example:

	hash := GeneratePasswordHash("mysecretpassword")
	fmt.Println(hash) // $2a$10$...

Note:

	The bcrypt.MinCost parameter is used to determine the cost factor of the bcrypt algorithm, which affects the computational complexity of hashing the password.
*/
func GeneratePasswordHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

/*
ValidatePassword compares a given password with a bcrypt hash to validate if the password matches the hash.

Parameters:
- password (string): The password to validate.
- hash (string): The bcrypt hash to compare the password against.

Returns:
- bool: True if the password matches the hash, false otherwise.

Example:

	valid := ValidatePassword("mysecretpassword", "$2a$10$...")

Note:

	The function uses bcrypt.CompareHashAndPassword to compare the provided password with the stored hash.
*/
func ValidatePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/*
EncryptByKey encrypts the given plaintext using the provided secret key.

Parameters:
- plaintext (string): The plaintext to be encrypted.
- secret (string): The secret key used for encryption.

Returns:
- []byte: The encrypted data.
- error: An error if encryption fails.

The function generates a random salt and initialization vector (IV), derives a key from the secret key and salt, pads the plaintext data, encrypts the padded data using AES in CBC mode, generates an authentication hash, and combines the salt, authentication hash, IV, and ciphertext into the final encrypted result.
*/
func EncryptByKey(plaintext string, secret string) ([]byte, error) {
	salt := generateRandomBytes(16)
	iv := generateRandomBytes(16)

	key, err := deriveKey([]byte(secret), salt, nil)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintextBytes := []byte(plaintext)
	padText := padData(plaintextBytes)

	ciphertext := make([]byte, len(padText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padText)

	sha, err := generateAuthHash(append(iv, ciphertext...), key)
	if err != nil {
		return nil, err
	}
	fmt.Println(sha)

	result := append(salt, []byte(sha)...)
	result = append(result, iv...)
	result = append(result, ciphertext...)

	return result, nil
}

/*
DecryptByKey decrypts the given encoded byte slice using the provided secret key.

Parameters:
- decodedByte: The byte slice to be decrypted, containing salt, expected hash, IV, and ciphertext.
- secret: The secret key used for decryption.

Returns:
- string: The decrypted data as a string.
- error: An error if decryption fails or if the HMAC verification fails.
*/
func DecryptByKey(decodedByte []byte, secret string) (string, error) {
	if len(decodedByte) < 96 {
		return "", fmt.Errorf("invalid input: expected at least 96 bytes, got %d", len(decodedByte))
	}

	salt := decodedByte[:16]
	expectedHash := string(decodedByte[16:80])
	iv := decodedByte[80:96]
	ciphertext := decodedByte[96:]

	fmt.Println(expectedHash)

	key, err := deriveKey([]byte(secret), salt, nil)
	if err != nil {
		return "", err
	}

	sha, err := generateAuthHash(append(iv, ciphertext...), key)
	if err != nil {
		return "", err
	}

	if sha != string(expectedHash) {
		return "", fmt.Errorf("invalid HMAC_256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	return string(unpadData(decrypted)), nil
}

/*
Hex2bin decodes a hexadecimal string into a byte slice.

Parameters:
- hexString (string): The hexadecimal string to be decoded.

Returns:
- []byte: The byte slice decoded from the hexadecimal string.
- error: An error if the decoding process fails.

Example:

	binData, err := Hex2bin("48656c6c6f20576f726c64")
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
	} else {
		fmt.Println(binData) // [72 101 108 108 111 32 87 111 114 108 100]
	}

Note:

	The function removes any "0x" or "0X" prefixes from the input hexadecimal string before decoding it into bytes using hex.DecodeString.
*/
func Hex2bin(hexString string) ([]byte, error) {
	// Remove any "0x" or "0X" prefixes if present
	hexString = strings.TrimPrefix(hexString, "0x")
	hexString = strings.TrimPrefix(hexString, "0X")

	// Decode the hexadecimal string to bytes
	binData, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}

	return binData, nil
}

/*
generateAuthHash takes in a byte slice of data and a key, derives an authentication key using the key, and then generates a SHA256 HMAC hash of the data using the derived authentication key. It returns the hash as a string and any error encountered during the process.
*/
func generateAuthHash(data []byte, key []byte) (string, error) {
	authKey, err := deriveKey(key, nil, []byte("AuthorizationKey"))
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, authKey)
	mac.Write(data)
	sha := hex.EncodeToString(mac.Sum(nil))
	return sha, nil
}

/*
generateRandomBytes generates a byte slice of random bytes with the specified size using the crypto/rand.Reader as the source of randomness.

Parameters:
- size (int): The size of the byte slice to generate.

Returns:
- []byte: A byte slice containing random bytes of the specified size.

Note:

	The function uses the io.ReadFull function to fill the byte slice with random bytes from the rand.Reader source of randomness.
*/
func generateRandomBytes(size int) []byte {
	bytes := make([]byte, size)
	io.ReadFull(rand.Reader, bytes)
	return bytes
}

/*
deriveKey takes a secret key, a salt, and optional info as input parameters. It creates a new HKDF instance using SHA256 hashing algorithm with the provided secret key, salt, and info. It then reads 16 bytes from the HKDF instance to derive a key and returns the derived key as a byte slice along with any error encountered during the process.
*/
func deriveKey(secret []byte, salt []byte, info []byte) ([]byte, error) {
	keyHkdf := hkdf.New(sha256.New, secret, salt, info)
	key := make([]byte, 16)
	_, err := keyHkdf.Read(key)
	return key, err
}

/*
padData pads the input byte slice with PKCS#7 padding to ensure its length is a multiple of the AES block size.

Parameters:
- data ([]byte): The input byte slice to be padded.

Returns:
- []byte: The padded byte slice.

Note:
PKCS#7 padding involves appending bytes such that each byte appended has the value of the number of bytes being added. For example, if 4 bytes are needed for padding, the bytes [4, 4, 4, 4] will be added to the input data.
*/
func padData(data []byte) []byte {
	padding := aes.BlockSize - (len(data) % aes.BlockSize)
	padText := append(data, byte(padding))
	return append(padText, bytes.Repeat([]byte{byte(padding)}, padding-1)...)
}

/*
unpadData removes PKCS#7 padding from the input byte slice to retrieve the original unpadded data.

Parameters:
- data ([]byte): The byte slice containing the padded data to be unpadded.

Returns:
- []byte: The unpadded byte slice representing the original data.

Note:
PKCS#7 padding involves appending bytes where each byte indicates the number of bytes added for padding. The unpadData function reverses this process by removing the padding based on the last byte value in the input data.
*/
func unpadData(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
