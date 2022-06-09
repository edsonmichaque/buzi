package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/edsonmichaque/buzi/internal/providers/python"
	"github.com/edsonmichaque/buzi/types"
	"gopkg.in/yaml.v3"
)

func main() {
	f, err := os.Open("testdata/buzi.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, f); err != nil {
		log.Fatal(err)
	}

	var manifest types.Manifest
	if err := yaml.Unmarshal(buf.Bytes(), &manifest); err != nil {
		log.Fatal(err)
	}

	params := map[string]string{
		"module":  "github.com/author/buzi",
		"package": "buzi",
	}

	pipeline := python.Pipeline()

	files := make([]types.File, 0)
	for _, p := range pipeline {
		f, err := p.Apply(params, &manifest)
		if err != nil {
			log.Fatal(err)
		}

		files = append(files, f...)
	}

	fmt.Println(files)

	if err := write(files); err != nil {
		log.Fatal(err)
	}
}

func write(files []types.File) error {
	target := "tmp"

	for _, f := range files {
		fullpath := filepath.Join(target, f.Path)

		if strings.Contains(f.Path, "/") {
			dir := filepath.Dir(fullpath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		fd, err := os.Create(fullpath)
		if err != nil {
			return err
		}
		defer fd.Close()

		if _, err := io.Copy(fd, bytes.NewReader(f.Content)); err != nil {
			return err
		}

	}

	return nil
}
