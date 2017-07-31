package main

import (
	"fmt"
	"github.com/AdamJacobMuller/archive"
	"io"
	"os"
)

func main() {
	fi, err := archive.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	for {
		file, err := fi.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s is %d\n", file.Name(), len(file.Bytes()))
	}
}
