package steganography

import (
	"testing"

	"github.com/KDT2006/go-steg/pkg/crypto"
)

func TestEncoding(t *testing.T) {
	input := "../../testdata/am.png"
	output := "../../testdata/amEncoded.png"
	message := "Hello there, this is a test message"

	err := EncodeMessage(input, output, message)
	if err != nil {
		t.Errorf("EncodeMessage() error: %v", err)
	}

	decodedMessage, err := DecodeMessage(output)
	if err != nil {
		t.Errorf("DecodeMessage() error: %v", err)
	}

	if decodedMessage != message {
		t.Error("decoded message does not match the original message")
	}
}

func TestCrypto(t *testing.T) {
	message := "Hello there, this is a test message"
	key := "d3Q8fjQ9CQzosiIvHeQz3Lo0L09MkEXW"

	encrypted, err := crypto.EncryptMessage([]byte(key), []byte(message))
	if err != nil {
		t.Errorf("EncryptMessage() error: %v", err)
	}

	decrypted, err := crypto.DecryptMessage([]byte(key), encrypted)
	if err != nil {
		t.Errorf("DecryptMessage() error: %v", err)
	}

	if decrypted != message {
		t.Error("decrypted message does not match the original message")
	}
}
