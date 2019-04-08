package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/viper"
	"github.com/stoewer/go-strcase"
)

func main() {
	home := os.Getenv("HOME")

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("usage: newpost <postname>")
	}

	postname := strings.Join(args, " ")

	viper.SetConfigFile(home + "/.newpost/config.json")
	err := viper.ReadInConfig()
	if err != nil {
		if _, err := os.Stat(home + "/.newpost/config.json"); os.IsNotExist(err) {
			os.MkdirAll(home+"/.newpost", 0777)
			os.Create(home + "/.newpost/config.json")
		}
	}

	projectPath := viper.GetString("projectPath")
	postsFolder := viper.GetString("postsFolder")
	postsFolderPath := filepath.Join(projectPath, postsFolder)

	if projectPath == "" {
		fmt.Printf("missing projectPath value in  %s/.newpost/config.json\n", home)
		os.Exit(0)
	}

	if postsFolder == "" {
		fmt.Printf("missing postsFolder value in  %s/.newpost/config.json\n", home)
		os.Exit(0)
	}

	slug := strcase.KebabCase(postname)
	postFolderPath := filepath.Join(postsFolderPath, slug)
	os.MkdirAll(postFolderPath, 0777)
	postFilePath := createPostFile(postFolderPath, postname)

	openInEditor(projectPath, postFilePath)

	fmt.Printf("created post %s", postFolderPath)
}

func createPostFile(postFolderPath, postname string) string {
	title := titleize(postname)
	date := time.Now().Format("2006-01-02")
	fileContent := fmt.Sprintf(`
---
title: %s
date: %s
tags:
	- TAG_ONE
	- TAG_TWO
---

# title
		`, title, date)

	postFilePath := filepath.Join(postFolderPath, "index.md")
	err := ioutil.WriteFile(postFilePath, []byte(strings.Trim(fileContent, "\n")), 0777)
	if err != nil {
		panic(err)
	}

	return postFilePath
}

func titleize(input string) (titleized string) {
	isToUpper := false
	for k, v := range input {
		if k == 0 {
			titleized = strings.ToUpper(string(input[0]))
		} else {
			if isToUpper || unicode.IsUpper(v) {
				titleized += " " + strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if (v == '_') || (v == ' ') {
					isToUpper = true
				} else {
					titleized += string(v)
				}
			}
		}
	}
	return

}

func openInEditor(projectPath, filePath string) {
	editor := os.Getenv("EDITOR")

	var cmd *exec.Cmd
	if editor == "code" || editor == "code-insiders" {
		cmd = exec.Command(editor, projectPath, "--goto", filePath)
	} else {
		cmd = exec.Command(editor, projectPath)
	}

	cmd.Start()
}
