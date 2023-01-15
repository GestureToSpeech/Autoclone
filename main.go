package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type Config struct {
	PushFolder string
	PullFolder string
	Repos      []struct {
		Ssh string
		Key string
	}
}

func main() {
	// Open the file
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Unmarshal the JSON data into the Config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Print the values of the MainFolder and Repos fields
	fmt.Println("Push folder:", config.PushFolder)
	fmt.Println("Pull folder:", config.PullFolder)
	fmt.Println("Repos:")
	for _, repo := range config.Repos {
		fmt.Println("  SSH:", repo.Ssh)
		fmt.Println("  Key:", repo.Key)
	}

	// remove pull and push folders
	err = executeCommand("", "rm", "-rf", fmt.Sprintf("\"%s\"", config.PushFolder))
	if err != nil {
		fmt.Println("Error removing push folder", err)
		return
	}
	err = executeCommand("", "rm", "-rf", fmt.Sprintf("\"%s\"", config.PullFolder))
	if err != nil {
		fmt.Println("Error removing pull folder", err)
		return
	}

	// create pull and push folders
	err = executeCommand("", "mkdir", fmt.Sprintf("\"%s\"", config.PushFolder))
	if err != nil {
		fmt.Println("Error creating push folder", err)
		return
	}
	err = executeCommand("", "mkdir", fmt.Sprintf("\"%s\"", config.PullFolder))
	if err != nil {
		fmt.Println("Error creating pull folder", err)
		return
	}
}

func executeCommand(dir string, commandName string, arg ...string) error {
	cmd := exec.Command(commandName, arg...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
