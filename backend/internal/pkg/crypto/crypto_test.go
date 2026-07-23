package crypto

import (
	"os"
	"testing"
)

const testKey = "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2"

func TestMain(m *testing.M) {
	// Set valid test key before all tests
	os.Setenv(EnvEncryptionKey, testKey)
	defer os.Unsetenv(EnvEncryptionKey)

	m.Run()
}

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	plaintext := "s3cret!Passw0rd_MySQL_2026"

	encrypted, err := EncryptString(plaintext)
	if err != nil {
		t.Fatalf("EncryptString failed: %v", err)
	}

	if encrypted == "" {
		t.Fatal("encrypted result should not be empty")
	}

	if encrypted == plaintext {
		t.Fatal("encrypted result should differ from plaintext")
	}

	decrypted, err := DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString failed: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf("roundtrip failed: got '%s', want '%s'", decrypted, plaintext)
	}
}

func TestEncryptDecrypt_EmptyString(t *testing.T) {
	plaintext := ""

	encrypted, err := EncryptString(plaintext)
	if err != nil {
		t.Fatalf("EncryptString failed: %v", err)
	}

	decrypted, err := DecryptString(encrypted)
	if err != nil {
		t.Fatalf("DecryptString failed: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf("roundtrip failed for empty string: got '%s'", decrypted)
	}
}

func TestEncrypt_UniqueCiphertext(t *testing.T) {
	// Encrypt same password twice should produce different ciphertext (due to random nonce)
	plaintext := "same_password_123"

	c1, err := EncryptString(plaintext)
	if err != nil {
		t.Fatalf("first EncryptString failed: %v", err)
	}

	c2, err := EncryptString(plaintext)
	if err != nil {
		t.Fatalf("second EncryptString failed: %v", err)
	}

	if c1 == c2 {
		t.Fatal("encrypting same plaintext twice should produce different output (random nonce)")
	}
}

func TestDecrypt_InvalidCiphertext(t *testing.T) {
	_, err := DecryptString("invalid-hex-string")
	if err == nil {
		t.Fatal("expected error for invalid hex string")
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	plaintext := "my_secret_password"

	encrypted, err := EncryptString(plaintext)
	if err != nil {
		t.Fatalf("EncryptString failed: %v", err)
	}

	// Change key before decrypting (simulate wrong key)
	os.Setenv(EnvEncryptionKey, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	defer os.Setenv(EnvEncryptionKey, testKey)

	_, err = DecryptString(encrypted)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestEncrypt_NoKey(t *testing.T) {
	os.Unsetenv(EnvEncryptionKey)
	defer os.Setenv(EnvEncryptionKey, testKey)

	_, err := EncryptString("test")
	if err == nil {
		t.Fatal("expected error when encryption key is not set")
	}
}

func TestEncrypt_InvalidKeyLength(t *testing.T) {
	os.Setenv(EnvEncryptionKey, "tooshort")
	defer os.Setenv(EnvEncryptionKey, testKey)

	_, err := EncryptString("test")
	if err == nil {
		t.Fatal("expected error for invalid key length")
	}
}
