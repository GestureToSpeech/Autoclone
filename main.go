package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	PushFolder  string
	PullFolder  string
	SshPushBase string
	Repos       []struct {
		Ssh string
		Key string
	}
}

func main() {
	// Open the file
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Printf("Error reading file: %s", err)
		return
	}

	// Unmarshal the JSON data into the Config struct
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Printf("Error parsing JSON: %s", err)
		return
	}

	/*
		// remove pull and push folders
		err = executeCommand("", "rm", "-rf", fmt.Sprintf("\"%s\"", config.PushFolder))
		if err != nil {
			log.Printf("Error removing push folder %s", err)
			return
		}
		err = executeCommand("", "rm", "-rf", fmt.Sprintf("\"%s\"", config.PullFolder))
		if err != nil {
			log.Printf("Error removing pull folder %s", err)
			return
		}
	*/
	// create pull and push folders
	err = executeCommand("", "mkdir", "-p", config.PushFolder)
	if err != nil {
		log.Printf("Error creating push folder %s", err)
		return
	}
	err = executeCommand("", "mkdir", "-p", config.PullFolder)
	if err != nil {
		log.Printf("Error creating pull folder %s", err)
		return
	}

	// clone repos
	for _, repo := range config.Repos {
		log.Printf("Cloning repo %s", repo.Ssh)
		err = cloneRepo(config.PullFolder, repo.Ssh)
		if err != nil {
			log.Printf("Error cloning repo %s; error message: %s", repo.Ssh, err)
			return
		}

		gitLabRepo := config.SshPushBase + getRepoName(repo.Ssh) + ".git"
		err = cloneRepo(config.PushFolder, gitLabRepo)
		if err != nil {
			log.Printf("Error cloning repo %s; error message: %s", gitLabRepo, err)
			return
		}

		err = fetchOrigin(config.PullFolder, repo.Ssh)
		if err != nil {
			log.Printf("Error fetching origin repo %s; error message: %s", repo.Ssh, err)
			return
		}

		allBranches, err := getAllBranches(config.PullFolder, repo.Ssh)
		if err != nil {
			log.Printf("Couldn't get all branches from repo %s; error message: %s", repo.Ssh, err)
			return
		}

		for _, branch := range allBranches {
			log.Print(branch)
		}
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

func getRepoName(ssh string) string {
	repoSSHParts := strings.Split(ssh, "/")
	repoName := repoSSHParts[len(repoSSHParts)-1]
	return strings.TrimSuffix(repoName, ".git")
}

func getRepoFolder(ssh string, folder string) string {
	repoName := getRepoName(ssh)

	return folder + repoName + "/"
}

func copyFiles(pullDir string, pushDir string, repoSSH string) error {
	destDir := getRepoFolder(repoSSH, pushDir)
	originDir := getRepoFolder(repoSSH, pullDir)

	// remove old files
	err := executeCommand(
		"",
		"[",
		"-d",
		fmt.Sprintf("\"%s\"", destDir),
		"]",
		"&&",
		"find",
		fmt.Sprintf("\"%s\"", destDir),
		"-mindepth",
		"1",
		"-not",
		"-path",
		fmt.Sprintf("\"%s/.git/*\"", destDir),
		"-not",
		"-path",
		fmt.Sprintf("\"%s/.git\"", destDir),
		"-delete",
	)
	if err != nil {
		return err
	}

	// copy files
	err = executeCommand(
		"",
		"rsync",
		"-av",
		"--exclude=\".git\"",
		fmt.Sprintf("\"%s\"", originDir),
		fmt.Sprintf("\"%s\"", pushDir),
	)

	return err
}

func cloneRepo(dir string, repoSSH string) error {
	repoFolder := getRepoFolder(repoSSH, dir)

	_, err := os.Stat(repoFolder)
	if !os.IsNotExist(err) {
		log.Print("Repository already initialized")
		return nil
	}

	log.Printf("Initializing repository %s", repoSSH)
	err = executeCommand(dir, "git", "clone", repoSSH)
	return err
}

func getAllBranches(pullDir string, repoSSH string) ([]string, error) {
	cmd := exec.Command("git", "branch", "-a")
	cmd.Dir = getRepoFolder(repoSSH, pullDir)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var branches []string
	for _, line := range lines {
		if strings.HasPrefix(line, "*") {
			line = strings.TrimPrefix(line, "*")
		}
		line = strings.TrimSpace(line)
		if line != "" {
			branches = append(branches, line)
		}
	}

	return branches, nil
}

func fetchOrigin(dir string, repoSSH string) error {
	return executeCommand(getRepoFolder(repoSSH, dir), "git", "fetch", "origin")
}
