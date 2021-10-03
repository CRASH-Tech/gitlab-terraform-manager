package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type cfgImpl struct {
	ServerAddress string
	ProjectID     string
	ProjectPath   string
	AccessToken   string
}

var (
	Config cfgImpl
)

func init() {
	var exist bool
	if Config.ServerAddress, exist = os.LookupEnv("GITLAB_URL"); !exist {
		panic("No gitlab url provided")
	}
	if Config.AccessToken, exist = os.LookupEnv("GITLAB_TOKEN"); !exist {
		panic("No gitlab access token provided")
	}
	if Config.ProjectID, exist = os.LookupEnv("GITLAB_PROJECT_ID"); !exist {
		panic("No gitlab project id provided")
	}
	if Config.ProjectPath, exist = os.LookupEnv("GITLAB_PROJECT_PATH"); !exist {
		panic("No gitlab project path provided")
	}
}

func main() {
	gitlab := gitlabInit(Config.ServerAddress, Config.ProjectID, Config.ProjectPath, Config.AccessToken)

	if len(os.Args) < 2 {
		panic("You must set operation command!")
	}

	switch os.Args[1] {
	case "list":
		GetStatesList(gitlab)
	case "save":
		SaveState(gitlab)
	case "saveall":
		SaveAllStates(gitlab)
	case "remove":
		RemoveState(gitlab)
	case "restore":
		RestoreState(gitlab)
	case "restoreall":
		RestoreAllStates(gitlab)
	default:
		panic("Unknown operation")
	}
}

func GetStatesList(gitlab *GitlabImpl) {
	for _, name := range gitlab.GitlabGetStatesList() {
		fmt.Println(name)
	}
}

func SaveState(gitlab *GitlabImpl) {
	if len(os.Args) < 3 {
		panic("You must set statefile name")
	}

	state_name := os.Args[2]
	state_data := gitlab.GetState(state_name)

	WriteFile(state_data, fmt.Sprintf("%s.%d.json", state_name, time.Now().Unix()))
}

func SaveAllStates(gitlab *GitlabImpl) {
	if len(os.Args) < 3 {
		panic("You must set output folder")
	}

	path := os.Args[2]

	for _, state_name := range gitlab.GitlabGetStatesList() {
		fmt.Println("Processing:", state_name)

		state_data := gitlab.GetState(state_name)
		WriteFile(state_data, fmt.Sprintf("%s/%s.%d.json", path, state_name, time.Now().Unix()))
	}
}

func RemoveState(gitlab *GitlabImpl) {
	if len(os.Args) < 3 {
		panic("You must set statefile name")
	}

	state_name := os.Args[2]
	gitlab.RemoveState(state_name)
}

func RestoreState(gitlab *GitlabImpl) {
	if len(os.Args) < 3 {
		panic("You must set statefile to restore")
	}
	state_path := os.Args[2]
	filename := filepath.Base(state_path)
	state_name := filename[:strings.IndexByte(filename, '.')]

	gitlab.RestoreState(state_name, ReadFile(state_path))
}

func RestoreAllStates(gitlab *GitlabImpl) {
	if len(os.Args) < 3 {
		panic("You must set statefile to restore")
	}
	states_path := os.Args[2]

	for _, stateFilename := range GetSavedStates(states_path) {
		fmt.Println("Processing:", stateFilename)
		state_name := stateFilename[:strings.IndexByte(stateFilename, '.')]

		gitlab.RestoreState(state_name, ReadFile(fmt.Sprintf("%s/%s", states_path, stateFilename)))
	}
}
