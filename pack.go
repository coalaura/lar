package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/coalaura/arguments"
)

func pack() {
	in := arguments.GetNamedAs("i", "input", "*")

	info("Collecting files...")

	var (
		files       []File
		directories []File
	)

	paths, err := filepath.Glob(in)
	if err != nil {
		fatalf(1, "failed to glob: %v", err)
	}

	out, err := arguments.NamedFile("o", "output", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644, os.Stdout)
	if err != nil {
		fatalf(2, "failed to open output: %v", err)
	}

	for _, path := range paths {
		inf, err := os.Stat(path)
		if err != nil {
			fatalf(3, "failed to stat: %v", err)
		}

		path = strings.ReplaceAll(path, "\\", "/")
		parms := uint16(inf.Mode().Perm())
		size := uint64(inf.Size())

		if inf.IsDir() {
			directories = append(directories, File{
				Path:  path,
				Perms: parms,
			})
		} else {
			files = append(files, File{
				Path:  path,
				Perms: parms,
				Size:  size,
			})
		}

		info(" - Added %s", path)
	}

	info("Found %d files, %d directories", len(files), len(directories))

	threads := arguments.GetNamedAs("t", "threads", runtime.NumCPU())

	info("Using %d threads", threads)

	jobs := make(chan File, threads)
	results := make(chan []byte, threads)

	var (
		swg sync.WaitGroup
		wwg sync.WaitGroup
	)

	wwg.Add(1)

	buffer := NewBuffer(out)

	// Write magic number (LAR in hex plus magic number 0x93)
	_, err = buffer.Write([]byte{0x4C, 0x41, 0x52, 0x93})
	if err != nil {
		fatalf(4, "failed to write header: %v", err)
	}

	info("Processing directories...")
	for _, directory := range directories {
		err = WriteDirectoryHeader(buffer, directory)
		if err != nil {
			fatalf(5, "failed to write directory header: %v", err)
		}
	}

	go func() {
		defer wwg.Done()

		for chunk := range results {
			_, err = buffer.Write(chunk)
			if err != nil {
				fatalf(4, "failed to write output: %v", err)
			}
		}

		err = buffer.Close()
		if err != nil {
			fatalf(4, "failed to close output: %v", err)
		}
	}()

	for i := 0; i < threads; i++ {
		swg.Add(1)

		go func() {
			defer swg.Done()

			var buf bytes.Buffer

			for file := range jobs {
				err := WriteFileHeader(&buf, file)
				if err != nil {
					fatalf(6, "failed to create file header: %v", err)
				}

				f, err := os.OpenFile(file.Path, os.O_RDONLY, 0)
				if err != nil {
					fatalf(6, "failed to open file: %v", err)
				}

				_, err = buf.ReadFrom(f)
				if err != nil {
					fatalf(6, "failed to read file: %v", err)
				}

				results <- buf.Bytes()

				buf.Reset()
			}
		}()
	}

	info("Processing files...")
	for _, path := range files {
		jobs <- path
	}

	close(jobs)
	swg.Wait()

	close(results)
	wwg.Wait()
}
