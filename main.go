package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
)

func main() {
	// Define CLI flags
	mode := flag.String("mode", "", "Mode: encode or decode")
	inputFile := flag.String("inputFile", "", "Path to an input image file")
	messageOutput := flag.String("messageOutput", "", "Path to the decoded message file")
	outputFile := flag.String("outputFile", "", "Path to the output image file (this or output required for encoding)")
	message := flag.String("message", "", "Message to encode (this or messageFile required for encoding)")
	messageFile := flag.String("messageFile", "", "Message file to pull message from(this or message required for encoding)")
	key := flag.String("key", "", "Encryption key (32 bytes for AES-256)")

	flag.Parse()

	// Basic error validation
	if *mode != "encode" && *mode != "decode" {
		log.Fatal("Error: Invalid mode. use 'encode' or 'decode'.")
	}

	if *key != "" && len(*key) != 32 {
		log.Fatal("Encryption key must be exactly 32 bytes.")
	}

	if *mode == "encode" {
		if *inputFile == "" || *outputFile == "" {
			log.Fatal("Error: 'encode' mode requires -input, -output.")
		}

		if *message == "" && *messageFile == "" {
			log.Fatal("Error: -message or -messageFile is required.")
		}

		fmt.Println("encoding message...")
		if *message != "" {
			if *key != "" {
				// Encrypt the message
				encryptedMessage, err := EncryptMessage([]byte(*key), []byte(*message))
				if err != nil {
					log.Fatal(err)
				}

				err = encodeMessage(*inputFile, *outputFile, encryptedMessage)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				err := encodeMessage(*inputFile, *outputFile, *message)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			// dealing with message file
			fh, err := os.Open(*messageFile)
			if err != nil {
				log.Fatal(err)
			}
			defer fh.Close()

			var messageContent string
			scanner := bufio.NewScanner(fh)
			for scanner.Scan() {
				messageContent += scanner.Text() + "\n"
			}

			if *key != "" {
				// Encrypt the message
				encryptedMessage, err := EncryptMessage([]byte(*key), []byte(messageContent))
				if err != nil {
					log.Fatal(err)
				}

				if err = encodeMessage(*inputFile, *outputFile, encryptedMessage); err != nil {
					log.Fatal(err)
				}
			} else {
				if err = encodeMessage(*inputFile, *outputFile, messageContent); err != nil {
					log.Fatal(err)
				}
			}
		}
	} else if *mode == "decode" {
		if *inputFile == "" {
			fmt.Println("Error: 'decode' mode requires -input.")
			os.Exit(1)
		}
		fmt.Println("decoding message...")
		message, err := decodeMessage(*inputFile)
		if err != nil {
			log.Fatal(err)
		}

		if *key != "" {
			// Decrypt the message
			decodedMessage, err := DecryptMessage([]byte(*key), message)
			if err != nil {
				log.Fatal(err)
			}

			if *messageOutput == "" {
				fmt.Println("Decoded message: ", decodedMessage)
			} else {
				// output message contents in file
				fh, err := os.Create(*messageOutput)
				if err != nil {
					log.Fatal(err)
				}
				defer fh.Close()

				writer := bufio.NewWriter(fh)
				n, err := writer.WriteString(decodedMessage)
				if err != nil {
					log.Fatal(err)
				}
				defer writer.Flush()
				log.Printf("Written %d bytes to disk.", n)
			}
		} else {
			if *messageOutput == "" {
				fmt.Println("Decoded message: ", message)
			} else {
				// output message contents in file
				fh, err := os.Create(*messageOutput)
				if err != nil {
					log.Fatal(err)
				}
				defer fh.Close()

				writer := bufio.NewWriter(fh)
				n, err := writer.WriteString(message)
				if err != nil {
					log.Fatal(err)
				}
				defer writer.Flush()
				log.Printf("Written %d bytes to disk.", n)
			}
		}
	} else {
		log.Fatal("Error: specify mode as either 'encode' or 'decode'.")
	}
}

func encodeMessage(inputPath, outputPath, message string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Decode image into image.Image to access pixel data
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Convert message to binary
	messageBytes := []byte(message)
	messageLength := len(messageBytes)

	// Convert length to binary(32 bits)
	lengthBinary := fmt.Sprintf("%032b", messageLength)

	// Combine length and message binary
	var binaryData []byte
	for _, bit := range lengthBinary {
		binaryData = append(binaryData, byte(bit-'0'))
	}
	binaryData = append(binaryData, bytesToBinary(messageBytes)...)

	// Check if message fits in the image
	if len(binaryData) > width*height {
		return fmt.Errorf("message is too large to fit in the image")
	}

	// Create a new image with the message encoded
	outImage := image.NewRGBA(bounds)
	binaryIndex := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if binaryIndex < len(binaryData) {
				r = (r & 0xFE) | uint32(binaryData[binaryIndex])
				binaryIndex++
			}
			outImage.Set(x, y, color.RGBA{
				uint8(r),
				uint8(g),
				uint8(b),
				uint8(a),
			})
		}
	}

	// Save the encoded image
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, outImage)
	if err != nil {
		return err
	}

	return nil
}

func bytesToBinary(data []byte) []byte {
	binary := make([]byte, 0, len(data)*8) // Preallocate enough space for 8 bits per byte
	for _, b := range data {
		for i := 7; i >= 0; i-- { // Extract bits from most to least significant
			binary = append(binary, (b>>i)&1)
		}
	}
	return binary
}

func decodeMessage(imagePath string) (string, error) {
	fh, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer fh.Close()

	img, _, err := image.Decode(fh)
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Extract binary data
	var binaryData []byte
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, _, _, _ := img.At(x, y).RGBA()
			binaryData = append(binaryData, byte(r&1)) // Get the LSB
		}
	}

	// Extract message length (first 32 bits)
	if len(binaryData) < 32 {
		return "", fmt.Errorf("image does not contain enough data")
	}

	lengthBinary := binaryData[:32]
	messageLength := 0
	for _, bit := range lengthBinary {
		messageLength = (messageLength << 1) | int(bit)
	}

	// Extract the message
	binaryMessage := binaryData[32 : 32+(messageLength*8)]
	decodedBytes, err := binaryToBytes(binaryMessage)
	if err != nil {
		return "", fmt.Errorf("failed to convert binary to bytes: %s", err)
	}

	return string(decodedBytes), nil
}

func binaryToBytes(binaryData []byte) ([]byte, error) {
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

func setLSB(value, bit byte) byte {
	return (value & 0xFE) | bit
}
