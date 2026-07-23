// Package crypto menyediakan enkripsi simetris AES-256-GCM
// untuk melindungi kredensial database tenant saat disimpan
// di platform database (encrypt at rest).
//
// Format penyimpanan:
//
//	[12-byte nonce][ciphertext]
//
// Nonce dibuat acak untuk setiap operasi enkripsi dan digabungkan
// dengan ciphertext untuk memudahkan dekripsi (nonce tidak perlu rahasia).
//
// Environment variable: HRIS_ENCRYPTION_KEY (32-byte hex-encoded key)
// Jika tidak diset, encryption/decryption akan melempar error.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	// EnvEncryptionKey adalah nama environment variable untuk kunci enkripsi.
	EnvEncryptionKey = "HRIS_ENCRYPTION_KEY"

	// keyHexLength adalah panjang hex-encoded 32-byte key = 64 karakter.
	keyHexLength = 64
)

var (
	ErrInvalidKeyLength = errors.New("encryption key must be 32 bytes (64 hex characters)")
	ErrEmptyKey         = errors.New("encryption key is not set")
	ErrCiphertextTooShort = errors.New("ciphertext too short")
	ErrInvalidHexKey    = errors.New("encryption key is not valid hex")
)

// Encrypt mengenkripsi plaintext menggunakan AES-256-GCM.
// Mengembalikan string hex-encoded yang berisi nonce + ciphertext.
func Encrypt(plaintext []byte) (string, error) {
	key, err := loadKey()
	if err != nil {
		return "", fmt.Errorf("encrypt: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("encrypt: failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encrypt: failed to create GCM: %w", err)
	}

	// Buat nonce 12-byte (ukuran standar untuk GCM)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encrypt: failed to generate nonce: %w", err)
	}

	// Seal: nonce + ciphertext + auth tag
	// Output: nonce || ciphertext || tag
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Encode ke hex
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt mendekripsi ciphertext (hex-encoded) yang dihasilkan oleh Encrypt.
// Mengembalikan plaintext asli.
func Decrypt(encoded string) ([]byte, error) {
	key, err := loadKey()
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	ciphertext, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("decrypt: invalid hex: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("decrypt: failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("decrypt: failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrCiphertextTooShort
	}

	// Pisahkan nonce dan ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Open: decrypt dan verifikasi auth tag
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptString adalah convenience wrapper untuk Encrypt dengan string input.
func EncryptString(plaintext string) (string, error) {
	return Encrypt([]byte(plaintext))
}

// DecryptString adalah convenience wrapper untuk Decrypt dengan output string.
func DecryptString(encoded string) (string, error) {
	plaintext, err := Decrypt(encoded)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// LooksEncrypted memeriksa apakah string tampaknya merupakan data terenkripsi AES-256-GCM.
// Mengembalikan true jika string adalah hex valid dengan panjang minimal untuk menampung nonce (12 bytes = 24 hex chars).
// Fungsi ini berguna untuk migration script agar tidak salah mengira data terenkripsi
// (dengan kunci berbeda) sebagai plaintext yang perlu dienkripsi ulang.
func LooksEncrypted(s string) bool {
	if len(s) < 24 { // minimum: 12 bytes nonce = 24 hex chars
		return false
	}
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return false
	}
	// Minimum: 12 bytes nonce + 16 bytes GCM auth tag = 28 bytes
	// Tapi kita hanya cek nonce size untuk deteksi awal yang aman
	return len(decoded) >= 12
}

// loadKey membaca dan memvalidasi encryption key dari environment variable.
func loadKey() ([]byte, error) {
	keyHex := os.Getenv(EnvEncryptionKey)
	if keyHex == "" {
		return nil, ErrEmptyKey
	}

	if len(keyHex) != keyHexLength {
		return nil, fmt.Errorf("%w: got %d characters, expected %d",
			ErrInvalidKeyLength, len(keyHex), keyHexLength)
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidHexKey, err.Error())
	}

	return key, nil
}
