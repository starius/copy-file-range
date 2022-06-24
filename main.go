package main

import (
	"flag"
	"os"

	"golang.org/x/sys/unix"
)

var (
	inputFile    = flag.String("input-file", "", "Input file")
	outputFile   = flag.String("output-file", "", "Output file")
	inputOffset  = flag.Int("input-offset", 0, "Input file offset")
	outputOffset = flag.Int("output-offset", 0, "Output file offset")
	length       = flag.Int("length", 0, "Length of fragment to copy")
)

func main() {
	flag.Parse()

	input, err := os.Open(*inputFile)
	if err != nil {
		panic(err)
	}
	defer input.Close()

	// Get permissions to copy to output file if it does not exist.
	stat, err := input.Stat()
	if err != nil {
		panic(err)
	}

	output, err := os.OpenFile(*outputFile, os.O_WRONLY|os.O_CREATE, stat.Mode())
	if err != nil {
		panic(err)
	}
	defer output.Close()

	err = unix.IoctlFileCloneRange(int(output.Fd()), &unix.FileCloneRange{
		Src_fd:      int64(input.Fd()),
		Src_offset:  uint64(*inputOffset),
		Src_length:  uint64(*length),
		Dest_offset: uint64(*outputOffset),
	})
	if err != nil {
		panic(err)
	}
}
