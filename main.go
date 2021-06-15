package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func getHash(thePath string) (string, error) {
	h := md5.New()
	fd, err := os.Open(thePath)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	io.Copy(h, fd)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func walker(thePath string, info fs.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	sum, err := getHash(thePath)
	if err != nil {
		return err
	}
	fmt.Printf("%s %s\n", sum, thePath)
	return nil
}

func mains(args []string) error {
	for _, arg1 := range args {
		files, err := filepath.Glob(arg1)
		if err != nil {
			files = []string{arg1}
		}
		for _, root := range files {
			stat, err := os.Stat(root)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}
			if stat.IsDir() {
				err = filepath.Walk(root, walker)
			} else {
				err = walker(root, stat, nil)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
