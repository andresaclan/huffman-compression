# Huffman Compression Tool

This project implements a text file lossless compression tool using Huffman encoding in Go.

## Features

- Compress text files using Huffman encoding
- Decompress files back to the original text
- Command-line interface for easy usage

## Usage

### Compression

To compress a file:

```sh
go run cmd/app/main.go <action (compress or decompress)> <input file> <output file>
