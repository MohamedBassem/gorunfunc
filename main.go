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
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func checkError(err error) {
	if err != nil {
		log.Fatal("ERR: " + err.Error())
	}
}

func checkRemError(err error, filename string) {
	if err != nil {
		os.Remove(filename)
		checkError(err)
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

func generateAndFmtFile(testFileName, fileContent string) {
	err := ioutil.WriteFile(testFileName, []byte(fileContent), 0644)
	checkError(err)

	cmd := exec.Command("goimports", "-w=true", testFileName)
	output, err := cmd.CombinedOutput()
	if len(output) == 0 && err != nil {
		checkRemError(err, testFileName)
	}
	fmt.Print(string(output))
}

func runFile(testname, filename string) {
	cmd := exec.Command("go", "test", "--run", testname)
	output, err := cmd.CombinedOutput()
	stdout := string(output)
	stdoutLines := strings.Split(stdout, "\n")
	if len(output) == 0 && err != nil {
		checkRemError(err, filename)
	} else if err != nil {
		stdout = strings.Join(stdoutLines[:len(stdoutLines)-2], "\n")
		fmt.Println(stdout)
	} else {
		stdout = strings.Join(stdoutLines[:len(stdoutLines)-3], "\n")
		fmt.Println(stdout)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	dryRun := flag.Bool("dry-run", false, "Print the test file")
	flag.Parse()

	packageName := extractPackage()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Err: A function call must be given")
	}

	functionCall := args[0]

	testRandomName := randString(10)
	testFileName := testRandomName + "_test.go"
	fileContent := fmt.Sprintf(templateString, packageName, testRandomName, functionCall)
	generateAndFmtFile(testFileName, fileContent)
	defer os.Remove(testFileName)

	if *dryRun {
		content, _ := ioutil.ReadFile(testFileName)
		fmt.Println(string(content))
	} else {
		runFile(testRandomName, testFileName)
	}
}

var templateString = `package %v
import "testing"

func Test_%v(t *testing.T) {
	  %v
}
`
