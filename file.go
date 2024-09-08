package main

import (
	"bytes"
	"encoding/binary"
)

func WriteFileHeader(buffer *bytes.Buffer, file File) error {
	// Write a 2 to indicate a file
	err := buffer.WriteByte(2)
	if err != nil {
		return err
	}

	// Then the path followed by a null byte
	buffer.WriteString(file.Path)
	buffer.WriteByte(0)

	// Then the permissions (2 bytes)
	err = binary.Write(buffer, binary.LittleEndian, file.Perms)
	if err != nil {
		return err
	}

	// Then the size (8 bytes)
	err = binary.Write(buffer, binary.LittleEndian, file.Size)
	if err != nil {
		return err
	}

	return nil
}

func WriteDirectoryHeader(buffer *BufferedWriter, file File) error {
	// Write a 1 to indicate a directory
	buffer.WriteByte(1)

	// Then the path
	buffer.WriteString(file.Path)

	// Followed by a null byte
	buffer.WriteByte(0)

	// Then the permissions
	return binary.Write(buffer, binary.LittleEndian, file.Perms)
}
