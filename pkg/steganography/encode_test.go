package steganography

import (
	"io"
	"log"
	"os"
	"testing"
)

const ()

func TestEncoding(t *testing.T) {
	input := "../../testdata/am.png"
	output := "../../testdata/amEncoded.png"
	message := "Hello there, this is a test message"

	err := EncodeMessage(input, output, message, "", false)
	if err != nil {
		t.Errorf("EncodeMessage() error: %v", err)
	}

	decodedMessage, err := DecodeMessage(output, "")
	if err != nil {
		t.Errorf("DecodeMessage() error: %v", err)
	}

	if decodedMessage != message {
		t.Error("decoded message does not match the original message")
	}
}

func TestEncodingWithEncryption(t *testing.T) {
	input := "../../testdata/large.png"
	output := "../../testdata/largeEncoded.png"
	message := "Hello there, this is a test message"
	key := "3GjmkamDG8k4JLxeCJ58KB0ne65wmJFl"

	err := EncodeMessage(input, output, message, key, true)
	if err != nil {
		t.Errorf("EncodeMessage() error: %v", err)
	}

	decodedMessage, err := DecodeMessage(output, key)
	if err != nil {
		t.Errorf("DecodeMessage() error: %v", err)
	}

	if decodedMessage != message {
		t.Error("decoded message does not match the original message")
	}
}

func TestCompression(t *testing.T) {
	messageFile := "../../testdata/compMessage.txt"
	var message string
	fh, err := os.Open(messageFile)
	if err != nil {
		t.Errorf("error opening message file: %v", err)
	}
	fc, err := io.ReadAll(fh)
	if err != nil {
		t.Errorf("error reading message from file: %v", err)
	}
	message = string(fc)
	log.Printf("Original length: %d", len(message))

	compressed, err := CompressString(message)
	if err != nil {
		t.Errorf("CompressString() error: %v", err)
	}
	log.Printf("Compressed length: %d", len(compressed))

	decompressed, err := DecompressString(compressed)
	if err != nil {
		t.Errorf("DecompressString() error: %v", err)
	}

	if string(compressed) == decompressed || len(compressed) >= len(decompressed) {
		t.Errorf("compression didn't work properly")
	}
}
