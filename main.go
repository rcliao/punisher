package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func computeMd5(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	return hash.Sum(result), nil
}

func visit(count *int, md5Map map[string][]string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			*count++

			if b, err := computeMd5(path); err != nil {
				fmt.Printf("Err: %v", err)
			} else {
				checksum := hex.EncodeToString(b)
				_, okay := md5Map[checksum]

				if okay {
					md5Map[checksum] = append(md5Map[checksum], path)
				} else {
					md5Map[checksum] = make([]string, 0)
					md5Map[checksum] = append(md5Map[checksum], path)
				}
			}
		}
		return nil
	}
}

func main() {
	count := 0
	md5Map := map[string][]string{}
	flag.Parse()
	root := flag.Arg(0)
	err := filepath.Walk(root, visit(&count, md5Map))
	fmt.Printf("filepath.Walk() returned %v\n", err)
	fmt.Printf("Total processed files %v\n", count)
	for key, value := range md5Map {
		fmt.Printf("Checksum: %s, Name: %v\n", key, value)
	}
}
