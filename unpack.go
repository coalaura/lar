package main

import (
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/coalaura/arguments"
)

func unpack() {
	in, err := arguments.NamedFile("i", "input", os.O_RDONLY, 0, os.Stdin)
	if err != nil {
		fatalf(1, "failed to open input: %v", err)
	}

	out := arguments.GetNamedAs("o", "output", ".")

	inf, err := os.Stat(out)
	if os.IsNotExist(err) {
		err = os.MkdirAll(out, 0755)
		if err != nil {
			fatalf(2, "failed to create output directory: %v", err)
		}
	} else if err != nil {
		fatalf(2, "failed to stat output directory: %v", err)
	} else if !inf.IsDir() {
		fatalf(2, "output is not a directory")
	}

	// First 4 bytes are the header 0x4C 0x41 0x52 0x93 (LAR in hex plus magic number)
	var header [4]byte

	_, err = in.Read(header[:])
	if err != nil {
		fatalf(3, "failed to read header: %v", err)
	}

	if header[0] != 0x4C || header[1] != 0x41 || header[2] != 0x52 || header[3] != 0x93 {
		fatalf(4, "input is not a valid LAR file")
	}

	// The a list of directories and files
	var (
		typ   [1]byte
		perms uint16
		size  uint64

		chunk   [1]byte
		builder strings.Builder
		path    string

		files       int
		directories int
	)

	info("Unpacking...")

	for {
		// Read the type (1=directory, 2=file)
		_, err = in.Read(typ[:])
		if err != nil {
			if err == io.EOF {
				break
			}

			fatalf(5, "failed to read file-type: %v", err)
		}

		// Read the path (terminated by a null byte)
		builder.Reset()

		for {
			_, err = in.Read(chunk[:])
			if err != nil {
				fatalf(5, "failed to read path: %v", err)
			}

			if chunk[0] == 0 {
				break
			}

			builder.WriteByte(chunk[0])
		}

		path = filepath.Join(out, builder.String())

		info("Unpacking %s...", path)

		// Read the permissions
		err = binary.Read(in, binary.LittleEndian, &perms)
		if err != nil {
			fatalf(5, "failed to read permissions: %v", err)
		}

		// If its a directory, we create it (if it doesn't exist)
		if typ[0] == 1 {
			_, err := os.Stat(path)
			if err != nil {
				err = os.MkdirAll(path, os.FileMode(perms))
				if err != nil {
					fatalf(6, "failed to create directory: %v", err)
				}
			}

			directories++
		} else if typ[0] == 2 {
			// If its a file, we first read its size
			err = binary.Read(in, binary.LittleEndian, &size)
			if err != nil {
				fatalf(5, "failed to read size: %v", err)
			}

			// Then we create the file
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(perms))
			if err != nil {
				fatalf(6, "failed to create file: %v", err)
			}

			// Then we copy the data
			_, err = io.CopyN(file, in, int64(size))
			if err != nil {
				fatalf(6, "failed to copy data: %v", err)
			}

			err = file.Close()
			if err != nil {
				fatalf(6, "failed to close file: %v", err)
			}

			files++
		}
	}

	info("Unpacked %d files and %d directories", files, directories)
}
