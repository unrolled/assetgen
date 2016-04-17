package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var jsonReplaceConst = "FILE_NAME"
var jsonTemplate = `{
  "images" : [
    {
      "idiom" : "universal",
      "filename" : "FILE_NAME.png",
      "scale" : "1x"
    },
    {
      "idiom" : "universal",
      "filename" : "FILE_NAME@2x.png",
      "scale" : "2x"
    },
    {
      "idiom" : "universal",
      "filename" : "FILE_NAME@3x.png",
      "scale" : "3x"
    }
  ],
  "info" : {
    "version" : 1,
    "author" : "xcode"
  }
}`

func main() {
	// Grab the directory from the os args.
	searchDirectory, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Ensure we have a trailing slash.
	if !strings.HasSuffix(searchDirectory, "/") {
		searchDirectory += "/"
	}

	// Walk the directory and grab files that need to be converted.
	filesToConvert := make(map[string]string)
	filepath.Walk(searchDirectory, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) != ".png" ||
			strings.Contains(path, ".imageset/") ||
			strings.Contains(path, ".xcassets/") ||
			strings.HasSuffix(path, "@2x.png") ||
			strings.HasSuffix(path, "@3x.png") {
			return nil
		}

		tmpDir, tmpFile := filepath.Split(path)
		filesToConvert[strings.Replace(tmpFile, ".png", "", 1)] = tmpDir
		return nil
	})

	// Iterate over file list and generate the asset catalogs.
	for file, directory := range filesToConvert {
		newDir := fmt.Sprintf("%s%s.imageset", directory, file)
		err = os.Mkdir(newDir, 0755)
		if err != nil {
			log.Printf("Folder already exists, skipping. (%s)\n", err)
			continue
		}

		// Copy files into the imageset folder.
		err = copyFile(fmt.Sprintf("%s/%s.png", directory, file), fmt.Sprintf("%s/%s.png", newDir, file))
		if err != nil {
			log.Fatal(err)
		}
		err = copyFile(fmt.Sprintf("%s/%s@2x.png", directory, file), fmt.Sprintf("%s/%s@2x.png", newDir, file))
		if err != nil {
			log.Printf("File did not exist, skipping. (%s)\n", err)
			continue
		}
		err = copyFile(fmt.Sprintf("%s/%s@3x.png", directory, file), fmt.Sprintf("%s/%s@3x.png", newDir, file))
		if err != nil {
			log.Printf("File did not exist, skipping. (%s)\n", err)
			continue
		}

		// Write JSON manifest.
		data := strings.Replace(jsonTemplate, jsonReplaceConst, file, -1)
		err := ioutil.WriteFile(fmt.Sprintf("%s/Contents.json", newDir), []byte(data), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// Write success message.
		log.Printf("Created %s\n", newDir)
	}
}

func copyFile(src, dst string) error {
	src_file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer src_file.Close()

	src_file_stat, err := src_file.Stat()
	if err != nil {
		return err
	}

	if !src_file_stat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	dst_file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dst_file.Close()

	_, err = io.Copy(dst_file, src_file)
	return err
}
