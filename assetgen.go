package main

import (
	"os"
	"path/filepath"
	"log"
	"strings"
	"fmt"
	"io"
	"io/ioutil"
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
	searchDirectory, err := filepath.Abs(filepath.Dir(os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}

	// Walk the directory and grab files that need to be converted.
	filesToConvert := make(map[string]string)
	filepath.Walk(searchDirectory, func(path string, f os.FileInfo, err error) error {
		pathLower := strings.ToLower(path)
		if filepath.Ext(pathLower) != ".png" ||
		   strings.Contains(pathLower, ".imageset/") ||
		   strings.Contains(pathLower, ".xcassets/") ||
		   strings.HasSuffix(pathLower, "@2x.png") ||
		   strings.HasSuffix(pathLower, "@3x.png") {
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
			log.Fatal(err)
		}

		// Copy files in.
		err = copyFile(fmt.Sprintf("%s/%s.png",directory,file), fmt.Sprintf("%s/%s.png",newDir, file))
		if err != nil {
			log.Fatal(err)
		}
		err = copyFile(fmt.Sprintf("%s/%s@2x.png",directory,file), fmt.Sprintf("%s/%s@2x.png",newDir, file))
		if err != nil {
			log.Fatal(err)
		}
		err = copyFile(fmt.Sprintf("%s/%s@3x.png",directory,file), fmt.Sprintf("%s/%s@3x.png",newDir, file))
		if err != nil {
			log.Fatal(err)
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
		return fmt.Errorf("%s is not a regular file", src)
	}

	dst_file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dst_file.Close()

	_, err = io.Copy(dst_file, src_file)
	return err
}
