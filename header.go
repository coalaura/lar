package main

import (
	"io"

	"github.com/coalaura/binary"
)

type FileHeader struct {
	Path  string
	Perms uint16
	Size  uint64
}

type DirectoryHeader struct {
	Path  string
	Perms uint16
}

func WriteFileHeader(writer io.Writer, header FileHeader) error {
	return binary.Write(writer, binary.LittleEndian, header)
}

func WriteDirectoryHeader(writer io.Writer, header DirectoryHeader) error {
	return binary.Write(writer, binary.LittleEndian, header)
}

func ReadFileHeader(reader io.Reader) (FileHeader, error) {
	var header FileHeader

	err := binary.Read(reader, binary.LittleEndian, &header)
	if err != nil {
		return header, err
	}

	return header, nil
}

func ReadDirectoryHeader(reader io.Reader) (DirectoryHeader, error) {
	var header DirectoryHeader

	err := binary.Read(reader, binary.LittleEndian, &header)
	if err != nil {
		return header, err
	}

	return header, nil
}
