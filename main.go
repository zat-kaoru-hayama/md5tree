package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var flagNoCr = flag.Bool("nocr", false, "not count CR code")

var flagAll = flag.Bool("a", false, "do not ignore dot files")

func getHash(thePath string) (string, error) {
	h := md5.New()
	fd, err := os.Open(thePath)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	if *flagNoCr {
		var buffer [1024]byte
		for {
			n, err := fd.Read(buffer[:])
			for i := 0; i < n; i++ {
				if buffer[i] != '\r' {
					h.Write([]byte{buffer[i]})
				}
			}
			if err != nil {
				break
			}
		}
	} else {
		io.Copy(h, fd)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func walker(thePath string, info fs.FileInfo, err error) error {
	name := info.Name()
	if len(name) >= 2 && name[0] == '.' && !*flagAll {
		if info.IsDir() {
			return filepath.SkipDir
		} else {
			return nil
		}
	}
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
			if name := stat.Name(); len(name) >= 2 && name[0] == '.' && !*flagAll {
				continue
			}
			if stat.IsDir() {
				err = filepath.Walk(root, walker)
			} else {
				err = walker(root, stat, nil)
			}
			if err != nil && err != filepath.SkipDir {
				return err
			}
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
