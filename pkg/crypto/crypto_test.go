package crypto

import "testing"

func TestCrypto(t *testing.T) {
	message := "Hello there, this is a test message"
	key := "d3Q8fjQ9CQzosiIvHeQz3Lo0L09MkEXW"

	encrypted, err := EncryptMessage([]byte(key), []byte(message))
	if err != nil {
		t.Errorf("EncryptMessage() error: %v", err)
	}

	decrypted, err := DecryptMessage([]byte(key), encrypted)
	if err != nil {
		t.Errorf("DecryptMessage() error: %v", err)
	}

	if decrypted != message {
		t.Error("decrypted message does not match the original message")
	}
}
