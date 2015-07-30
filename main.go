package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

func checkError(err error) {
	if err != nil {
		log.Fatal("ERR: " + err.Error())
	}
}

func randString(length int) string {
	var ret string
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPKRSTUVWXYZ0123456789")
	for i := 0; i < length; i++ {
		ret += string(runes[rand.Intn(len(runes))])
	}
	return ret
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
	rand.Seed(time.Now().Unix())
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

func main() {
	dryRun := flag.Bool("dry-run", false, "Print the test file")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Err: At least a function name should be given")
	}

	packageName := extractPackage()
	functionName := args[0]
	functionArgs := ""
	if len(args) > 1 {
		functionArgs = strings.Join(args[1:], ", ")
	}

	fileContent := fmt.Sprintf(templateString, packageName, randString(10), functionName, functionArgs)
	if *dryRun {
		fmt.Println(fileContent)
	}
}

var templateString = `package %v
import "testing"

func Test_%v(t *testing.T) {
	  %v(%v)
}
`
