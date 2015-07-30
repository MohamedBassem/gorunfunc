package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func checkError(err error) {
	if err != nil {
		log.Fatal("ERR: " + err.Error())
	}
}

func extractFilePackage(filename string) (string, error) {
	regex, err := regexp.Compile("^ *package (.*)$")
	checkError(err)
	fh, err := os.Open(filename)
	f := bufio.NewReader(fh)

	checkError(err)
	defer fh.Close()

	buf := make([]byte, 1024)
	for {
		buf, _, err = f.ReadLine()
		if err != nil {
			return "", errors.New("Package def in file " + filename + " not found.")
		}
		s := string(buf)
		if regex.MatchString(s) {
			return strings.Split(string(buf), " ")[1], nil
		}
	}
}

func extractPackage() string {
	workingDir, err := os.Getwd()
	checkError(err)
	dirFiles, err := ioutil.ReadDir(workingDir)
	checkError(err)
	goFiles := make([]string, 0)
	for _, file := range dirFiles {
		if !file.IsDir() {
			isGoFile, _ := regexp.Match("\\.go$", []byte(file.Name()))
			if isGoFile {
				goFiles = append(goFiles, file.Name())
			}
		}
	}
	if len(goFiles) == 0 {
		log.Fatal("Err: No go source files were found.")
	}
	packageNames := make(map[string]bool)
	for _, file := range goFiles {
		packageName, err := extractFilePackage(file)
		checkError(err)
		packageNames[packageName] = true
	}
	if len(packageNames) > 1 {
		log.Fatal("Err: Multiple package definitions in the same directory")
	}

	var packageName string
	for k := range packageNames {
		packageName = k
		break
	}
	return packageName
}
