package steganography

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

func BytesToBinary(data []byte) []byte {
	binary := make([]byte, 0, len(data)*8) // Preallocate enough space for 8 bits per byte
	for _, b := range data {
		for i := 7; i >= 0; i-- { // Extract bits from most to least significant
			binary = append(binary, (b>>i)&1)
		}
	}
	return binary
}

func BinaryToBytes(binaryData []byte) ([]byte, error) {
	if len(binaryData)%8 != 0 {
		return nil, fmt.Errorf("binary data length is not a multiple of 8")
	}

	var buffer bytes.Buffer
	for i := 0; i < len(binaryData); i += 8 {
		var b byte
		for j := 0; j < 8; j++ {
			b = (b << 1) | binaryData[i+j]
		}
		buffer.WriteByte(b)
	}

	return buffer.Bytes(), nil
}

func SetLSB(value, bit byte) byte {
	return (value & 0xFE) | bit
}

// Compress the input string using zlib
func CompressString(input string) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	_, err := writer.Write([]byte(input))
	if err != nil {
		return nil, err
	}
	writer.Close()
	return buf.Bytes(), nil
}

// Decompress the input byte slice using zlib
func DecompressString(input []byte) (string, error) {
	reader, err := zlib.NewReader(bytes.NewReader(input))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
