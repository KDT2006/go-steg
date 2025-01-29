package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"os"

	steg "github.com/KDT2006/go-steg/pkg/steganography"
)

func main() {
	// Define CLI flags
	mode := flag.String("mode", "", "Mode: encode | decode | capacity")
	inputFile := flag.String("i", "", "Path to an input image file")
	outputFile := flag.String("o", "", "Path to the output image file")
	message := flag.String("m", "", "Message to encode (this or -mf required for encoding)")
	messageFile := flag.String("mf", "", "Message file to pull message from(this or -m required for encoding)")
	key := flag.String("k", "", "Encryption key (32 bytes for AES-256)")
	compress := flag.Bool("c", false, "Compress the message")

	flag.Usage = func() {
		usage :=
			`Usage:
gosteg -mode "encode" -m "message" -i "input.png" -o "output.png"
gosteg -mode "decode" -i "encoded.png"
gosteg -mode "capacity" -i "input.png"

Flags:
		 	`
		fmt.Fprintln(flag.CommandLine.Output(), usage)
		flag.PrintDefaults()
	}

	flag.Parse()

	// Basic error validation
	if *mode != "encode" && *mode != "decode" && *mode != "capacity" {
		log.Fatal("Error: Invalid mode. use 'encode' | 'decode' | 'capacity'.")
	}

	if *key != "" && len(*key) != 32 {
		log.Fatal("Encryption key must be exactly 32 bytes.")
	}

	switch *mode {
	case "encode":
		if *inputFile == "" || *outputFile == "" {
			log.Fatal("Error: 'encode' mode requires -i, -o.")
		}

		if *message == "" && *messageFile == "" {
			log.Fatal("Error: -m or -mf is required.")
		}

		fmt.Println("encoding message...")

		// Target message to encode(encrypted or decrypted)
		var targetMessage string

		if *message != "" {
			targetMessage = *message
		} else {
			// dealing with message file
			fh, err := os.Open(*messageFile)
			if err != nil {
				log.Fatal(err)
			}
			defer fh.Close()

			fc, err := io.ReadAll(fh)
			if err != nil {
				log.Fatal(err)
			}
			targetMessage = string(fc)
		}

		err := steg.EncodeMessage(*inputFile, *outputFile, targetMessage, *key, *compress)
		if err != nil {
			log.Fatal(err)
		}

	case "decode":
		if *inputFile == "" {
			log.Println("Error: 'decode' mode requires -i.")
			os.Exit(1)
		}
		fmt.Println("decoding message...")
		message, err := steg.DecodeMessage(*inputFile, *key)
		if err != nil {
			log.Fatal(err)
		}

		if *outputFile == "" {
			fmt.Println("Decoded message: ", message)
		} else {
			// output message contents in file
			fh, err := os.Create(*outputFile)
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

	case "capacity":
		if *inputFile == "" {
			log.Fatal("Error: 'capacity' mode requires -i.")
		}

		fh, err := os.Open(*inputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer fh.Close()

		img, _, err := image.Decode(fh)
		if err != nil {
			log.Fatal(err)
		}

		bounds := img.Bounds()
		// -3 for the header
		capacity := (bounds.Dx() * bounds.Dy() / 8) - 3 // 1 byte of encoding per 8 pixels(red channel)
		fmt.Printf("This image can store upto %d bytes of data.\n", capacity)
	default:
		log.Fatal("Error: specify mode as either 'encode' or 'decode'.")
	}
}
