package main

import (
	"strings"
)

func countFiles(repoPath string) int {

	files := getFiles(repoPath)

	return len(files)

}

func countLanguages(repoPath string) map[string]int {

	// XXX tie-into languages defined in templates

	languages := map[string]int{}
	files := getFiles(repoPath)

	for _, file := range files {
		ext := strings.Split(file, ".")
		if _, ok := languages[ext[len(ext)-1]]; ok {
			languages[ext[len(ext)-1]] = languages[ext[len(ext)-1]]+1
		} else {
			languages[ext[len(ext)-1]] = 1
		}
	}

	return languages

}

func countLines(file string) int {

	data := getFileContents(file)
	lines := strings.Split(string(data), "\n")

	return len(lines)

}

func countLinesPerLanguage(repoPath string) map[string]int {

	counts := map[string]int{}
	files := getFiles(repoPath)
	for _, file := range files {
		ext := strings.Split(file, ".")
		if _, ok := counts[ext[len(ext)-1]]; ok {
			counts[ext[len(ext)-1]] = counts[ext[len(ext)-1]]+countLines(file)
		} else {
			counts[ext[len(ext)-1]] = countLines(file)
		}
	}

	return counts

}

func countCommits(repoPath string) int {

	commits, _ := getCommits(repoPath)

	return len(commits)

}

func countAuthors(repoPath string) int {

	authors := countAuthorCommits(repoPath)

	return len(authors)

}

func countAuthorsByCommits(repoPath string, commit string) int {

	counts := map[string]int{}
	commits, commitMap := getCommits(repoPath)
	index := -1
	for i, c := range commits {
		if strings.EqualFold(commit, c) {
			index = i
		}
	}
	if index != -1 {
		commits = commits[index:]
	}
	for k, commit := range commitMap {
		for _, c := range commits {
			if c == k {
				if _, ok := counts[commit["author"]]; ok {
					counts[commit["author"]] = counts[commit["author"]]+1
				} else {
					counts[commit["author"]] = 1
				}
			}
		}
	}

	return len(counts)

}

func countAuthorCommits(repoPath string) map[string]int {

	counts := map[string]int{}
	_, commits := getCommits(repoPath)
	for _, commit := range commits {
		if _, ok := counts[commit["author"]]; ok {
			counts[commit["author"]] = counts[commit["author"]]+1
		} else {
			counts[commit["author"]] = 1
		}
	}

	return counts

}

func countAuthorLines(repoPath string) map[string]int {

	// XXX stub, requires reading diffs
	counts := map[string]int{}

	return counts

}
