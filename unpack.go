package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/coalaura/arguments"
	"github.com/coalaura/binary"
)

func unpack() {
	in, err := arguments.File("i", "input", os.O_RDONLY, 0, os.Stdin)
	if err != nil {
		fatalf(1, "failed to open input: %v", err)
	}

	out := arguments.String("o", "output", ".")

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

	var (
		directories uint32
		files       int
	)

	info("Unpacking...")

	// Read the number of directories
	err = binary.Read(in, binary.LittleEndian, &directories)
	if err != nil {
		fatalf(5, "failed to read number of directories: %v", err)
	}

	// Read directories
	var directory DirectoryHeader

	for i := uint32(0); i < directories; i++ {
		err = binary.Read(in, binary.LittleEndian, &directory)
		if err != nil {
			fatalf(5, "failed to read directory header: %v", err)
		}

		_, err := os.Stat(directory.Path)
		if err != nil {
			err = os.MkdirAll(directory.Path, os.FileMode(directory.Perms))
			if err != nil {
				fatalf(6, "failed to create directory: %v", err)
			}
		}
	}

	// Then read all the files
	var (
		file FileHeader
		data []byte
	)

	for {
		err = binary.Read(in, binary.LittleEndian, &file)
		if err == io.EOF {
			break
		} else if err != nil {
			fatalf(5, "failed to read file header: %v", err)
		}

		data = make([]byte, file.Size)

		_, err = io.ReadFull(in, data)
		if err != nil {
			fatalf(5, "failed to read file data: %v", err)
		}

		err = os.WriteFile(filepath.Join(out, file.Path), data, os.FileMode(file.Perms))
		if err != nil {
			fatalf(6, "failed to write file: %v", err)
		}

		files++
	}

	info("Unpacked %d files and %d directories", files, directories)
}
