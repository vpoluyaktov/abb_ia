package utils

import (
	"bytes"
	"testing"
)

func TestGetMachineIdentifier(t *testing.T) {
	id1, err := GetMachineIdentifier()
	if err != nil {
		t.Fatalf("GetMachineIdentifier() failed: %v", err)
	}

	// Test length
	if len(id1) != 32 {
		t.Errorf("Expected identifier length of 32, got %d", len(id1))
	}

	// Test consistency
	id2, err := GetMachineIdentifier()
	if err != nil {
		t.Fatalf("Second GetMachineIdentifier() call failed: %v", err)
	}

	if !bytes.Equal(id1, id2) {
		t.Error("Machine identifiers should be consistent across calls")
	}
}

func TestGenerateEncryptionKey(t *testing.T) {
	key1, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey() failed: %v", err)
	}

	// Test length
	if len(key1) != 32 {
		t.Errorf("Expected key length of 32, got %d", len(key1))
	}

	// Test consistency
	key2, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("Second GenerateEncryptionKey() call failed: %v", err)
	}

	if !bytes.Equal(key1, key2) {
		t.Error("Encryption keys should be consistent across calls")
	}
}

func TestEncryptDecryptString(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{"empty string", ""},
		{"simple text", "Hello, World!"},
		{"special chars", "!@#$%^&*()_+"},
		{"unicode", "Hello, 世界"},
		{"long text", "This is a longer text that spans multiple blocks in the encryption process"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := EncryptString(tt.text)
			if err != nil {
				t.Fatalf("EncryptString() failed: %v", err)
			}

			// Verify not equal to original
			if string(encrypted) == tt.text {
				t.Error("Encrypted text should not be equal to original text")
			}

			// Decrypt
			decrypted, err := DecryptString(encrypted)
			if err != nil {
				t.Fatalf("DecryptString() failed: %v", err)
			}

			// Verify decrypted equals original
			if decrypted != tt.text {
				t.Errorf("Decrypted text does not match original.\nGot: %q\nWant: %q", decrypted, tt.text)
			}
		})
	}
}

func TestBase64EncodeDecode(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"empty", []byte("")},
		{"simple text", []byte("Hello, World!")},
		{"binary data", []byte{0, 1, 2, 3, 4, 5}},
		{"special chars", []byte("!@#$%^&*()_+")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			encoded := EncodeBase64(tt.data)

			// Verify not equal to original
			if string(encoded) == string(tt.data) && len(tt.data) > 0 {
				t.Error("Encoded data should not be equal to original data")
			}

			// Decode
			decoded, err := DecodeBase64(encoded)
			if err != nil {
				t.Fatalf("DecodeBase64() failed: %v", err)
			}

			// Verify decoded equals original
			if !bytes.Equal(decoded, tt.data) {
				t.Errorf("Decoded data does not match original.\nGot: %v\nWant: %v", decoded, tt.data)
			}
		})
	}
}

func TestDecryptStringErrors(t *testing.T) {
	tests := []struct {
		name       string
		ciphertext []byte
	}{
		{"empty ciphertext", []byte{}},
		{"too short", make([]byte, 15)}, // AES block size is 16
		{"invalid ciphertext", []byte("invalid")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecryptString(tt.ciphertext)
			if err == nil {
				t.Error("Expected error for invalid ciphertext, got nil")
			}
		})
	}
}

func TestDecodeBase64Errors(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
	}{
		{"invalid base64", "!invalid!"},
		{"incomplete base64", "SGVsbG8==="}, // Too many padding characters
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeBase64(tt.encoded)
			if err == nil {
				t.Error("Expected error for invalid base64, got nil")
			}
		})
	}
}
