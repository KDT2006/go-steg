package main

import (
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
	encodeFlag := flag.Bool("encode", false, "Encode a message into an image")
	decodeFlag := flag.Bool("decode", false, "Decode a message from an image")
	inputFile := flag.String("input", "", "Path to an input image file")
	outputFile := flag.String("output", "", "Path to the output image file (required for encoding)")
	message := flag.String("message", "", "Message to encode (required for encoding)")

	flag.Parse()

	// Basic error validation
	if *encodeFlag && *decodeFlag {
		fmt.Println("Error: Use either -encode or -decode, not both.")
		os.Exit(1)
	}

	if *encodeFlag {
		if *inputFile == "" || *outputFile == "" || *message == "" {
			fmt.Println("Error: -encode requires -input, -output, and -message.")
			os.Exit(1)
		}
		fmt.Println("encoding message...")
		err := encodeMessage(*inputFile, *outputFile, *message)
		if err != nil {
			log.Fatal(err)
		}
	} else if *decodeFlag {
		if *inputFile == "" {
			fmt.Println("Error: -decode requires -input.")
			os.Exit(1)
		}
		fmt.Println("decoding message...")
	} else {
		fmt.Println("Error: specify either -encode or -decode.")
		os.Exit(1)
	}
}

func encodeMessage(inputPath, outputPath, message string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Deocode image into image.Image to access pixel data
	img, _, err := image.Decode(inputFile)
	if err != nil {
		return err
	}

	// Convert image to RGBA(Provides easier access to individual color channels)
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	// Convert message to binary
	messageBytes := []byte(message)
	binaryMessage := bytesToBinary(messageBytes)

	fmt.Println(binaryMessage)

	// Embed the binary message
	index := 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if index >= len(binaryMessage) {
				break
			}

			r, g, b, a := rgba.At(x, y).RGBA()
			rgba.Set(x, y, color.RGBA{
				R: setLSB(uint8(r>>8), binaryMessage[index]),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			})
			index++
		}
	}

	// Save the modified image
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, rgba)
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

func setLSB(value, bit byte) byte {
	return (value & 0xFE) | bit
}
