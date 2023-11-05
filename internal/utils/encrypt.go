package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net"
)

func GetMachineIdentifier() ([]byte, error) {
	// Retrieve MAC address
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %v", err)
	}
	var macAddr string
	for _, ifa := range ifas {
			if ifa.Flags&net.FlagLoopback == 0 && ifa.Flags&net.FlagUp != 0 {
					macAddr = ifa.HardwareAddr.String()
					break
			}
	}
	
	// Hash the MAC address to create a unique identifier
	hash := sha256.Sum256([]byte(macAddr))
	return hash[:], nil
}

func GenerateEncryptionKey() ([]byte, error) {
	machineIdentifier, err := GetMachineIdentifier()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine identifier: %v", err)
	}
	
	// Use the machine identifier to generate an encryption key
	key := make([]byte, 32) // 256-bit key
	for i := 0; i < len(machineIdentifier); i += 4 {
			binary.LittleEndian.PutUint32(key[i:], binary.LittleEndian.Uint32(machineIdentifier[i:]))
	}
	
	return key, nil
}

func EncryptString(text string) ([]byte, error) {
	key, err := GenerateEncryptionKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption token: %v", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get chipher: %v", err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to random key for encryption: %v", err)
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))
	return ciphertext, nil
}

func DecryptString(ciphertext []byte) (string, error) {
	key, err := GenerateEncryptionKey()
	if err != nil {
		return "", fmt.Errorf("failed to get encryption token: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to get chipher: %v", err)
	}
	if len(ciphertext) < aes.BlockSize {
			return "", fmt.Errorf("ciphertext too short: %v", err)
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func DecodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}
