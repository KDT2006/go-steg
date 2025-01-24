package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

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
		return fmt.Errorf("message is too large to fit in the image (requires %d bytes but only %d bytes available)",
			len(binaryData)/8, width*height)
	}

	// Create a new image with the message encoded
	outImage := image.NewRGBA(bounds)
	binaryIndex := 0
	totalPixels := (width * height) + 2
	progressInterval := totalPixels / 100 // Update progress every 1%
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

			// Update progress for every pixel processed
			pixelIndex := y*width + x + 1 // +1 for saving the new image
			if pixelIndex%progressInterval == 0 || pixelIndex == totalPixels {
				progress := (pixelIndex * 100) / totalPixels
				fmt.Printf("\rEncoding progress: %d%%", progress)
			}
		}
	}
	fmt.Println()

	// Save the encoded image
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	outputPathSplit := strings.Split(outputPath, ".")
	outputImageFormat := outputPathSplit[len(outputPathSplit)-1]
	switch outputImageFormat {
	case "png":
		if err = png.Encode(outputFile, outImage); err != nil {
			return err
		}
	case "jpeg", "jpg":
		if err = jpeg.Encode(outputFile, outImage, nil); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Error: unsupported file format: %s", outputImageFormat)
	}

	return nil
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

	totalPixels := width * height
	progressInterval := totalPixels / 100 // Update progress every 1%

	// Extract binary data
	var binaryData []byte
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, _, _, _ := img.At(x, y).RGBA()
			binaryData = append(binaryData, byte(r&1)) // Get the LSB

			// Update progress for every pixel processed
			pixelIndex := y*width + x + 1
			if pixelIndex%progressInterval == 0 || pixelIndex == totalPixels {
				progress := (pixelIndex * 100) / totalPixels
				fmt.Printf("\rDecoding progress: %d%%", progress)
			}
		}
	}
	fmt.Println()

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
