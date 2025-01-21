# Steganography Tool in Go

A robust command-line tool built in **Go** for encoding and decoding secret messages within images. This project was developed as a learning exercise to explore working with images, cryptography, and building CLI tools in Go.

---

## Features

- **Message Encoding**: Embed secret messages into image files securely.
- **AES Encryption**: Protect embedded messages using encryption.
- **Message Decoding**: Extract and decrypt messages from encoded images.
- **Multiple Image Format Support**: Works with PNG and JPEG image formats for now.
- **User-Friendly CLI**: Intuitive command-line interface for encoding, decoding, and configuring options.
- **File Input and Output**: Encode messages directly from files and output decoded messages to files.
- **Input Validation**: Graceful handling of invalid or missing inputs.
- **Error Handling and Logging**: Detailed logs for troubleshooting.

---

## Installation

Clone the repository and build the project:

```bash
git clone https://github.com/KDT2006/go-steg.git
cd go-steg
go build -o gosteg
