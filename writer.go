package main

import (
	"bytes"
	"os"
)

type BufferedWriter struct {
	file   *os.File
	buffer bytes.Buffer
}

func NewBuffer(file *os.File) *BufferedWriter {
	return &BufferedWriter{file: file}
}

func (b *BufferedWriter) Flush(force bool) error {
	if force || b.buffer.Len() >= 1024*1024 {
		_, err := b.file.Write(b.buffer.Bytes())
		if err != nil {
			return err
		}

		b.buffer.Reset()
	}

	return nil
}

func (b *BufferedWriter) Write(p []byte) (int, error) {
	n, _ := b.buffer.Write(p)

	return n, b.Flush(false)
}

func (b *BufferedWriter) WriteByte(bt byte) error {
	b.buffer.WriteByte(bt)

	return b.Flush(false)
}

func (b *BufferedWriter) WriteString(s string) error {
	b.buffer.WriteString(s)

	return b.Flush(false)
}

func (b *BufferedWriter) Close() error {
	err := b.Flush(true)
	if err != nil {
		return err
	}

	return b.file.Close()
}
