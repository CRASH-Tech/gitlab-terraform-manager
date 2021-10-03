package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func WriteFile(data, path string) {

	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func ReadFile(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data)
}

type StateFileName struct {
	Filename  string
	StateName string
	Timestamp int
}

func GetSavedStates(path string) (result []string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	tmpResult := make(map[string]StateFileName)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			parsedName := ParseStateName(f.Name())
			if tmpResult[parsedName.StateName].Timestamp < parsedName.Timestamp {
				tmpResult[parsedName.StateName] = parsedName
			}
		}
	}
	for _, file := range tmpResult {
		result = append(result, file.Filename)
	}
	return
}

func ParseStateName(name string) (result StateFileName) {
	r := regexp.MustCompile(`(^.+)\.(\d+)`)
	match := r.FindStringSubmatch(name)
	if len(match) != 3 {
		panic(fmt.Sprintf("Wrong filename: %s", name))
	}

	timestamp, err := strconv.Atoi(match[2])
	if err != nil {
		panic(err)
	}
	result.Filename = name
	result.StateName = match[1]
	result.Timestamp = timestamp

	return
}
