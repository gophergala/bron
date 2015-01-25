package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

//	"github.com/gophergala/bron/filters"
)

var (
	blessedPtr   string
	dashboardPtr string
	repoPtr      string
	repoPathPtr  string
	verbosePtr   int
	vizPtr       bool
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	flag.StringVar(&blessedPtr, "blessedPath", "/go/src/github.com/yaronn/blessed-contrib", "Path where blessed-contrib is installed")
	flag.StringVar(&dashboardPtr, "dashboard", "example", "Name of dashboard to use for visualization")
	flag.StringVar(&repoPtr, "repo", "", "Git repository to scan")
	flag.StringVar(&repoPathPtr, "path", "", "Git repository file path to scan")
	flag.IntVar(&verbosePtr, "v", 1, "verbosity level")
	flag.BoolVar(&vizPtr, "viz", false, "Visualize the results, requires blessed")

	flag.Parse()

	if repoPtr == "" && repoPathPtr == "" {
		fmt.Println("please specify either a repo or a path to a git repo to scan")
	} else if repoPtr != "" && repoPathPtr != "" {
		fmt.Println("please specify only either a repo or a path to a git repo to scan, not both")
	} else {
		if verbosePtr > 0 {
			fmt.Println("going to scan repository", repoPtr, "...")
		}
	}

	if repoPtr != "" {
		uuidRepo := cloneRepo(repoPtr)

		// XXX example calls through all commits
		x, _ := getCommits(uuidRepo)
		for _, commit := range x {
			checkoutCommit(uuidRepo, commit)
			// XXX simple channel starts, for now
			files := getFiles(uuidRepo)
			parse(files)
		}
		checkoutCommit(uuidRepo, x[0])

		// XXX test template parsing
		templates := templateParse("templates")
		fmt.Println(templates)

		if vizPtr {
			// get data for dashboard
			languages := "["
			languageLines := "["
			languageMap := countLinesPerLanguage(uuidRepo)
			for key := range languageMap {
				languages += "'"+key+"', "
				languageLines += "'"+strconv.Itoa(languageMap[key])+"', "
			}
			languages = languages[0:len(languages)-2]+"]"
			languageLines = languageLines[0:len(languageLines)-2]+"]"

			authorMap := countAuthorCommits(uuidRepo)
			authors := "["
			for key := range authorMap {
				authors += "['"+key+"', '"+strconv.Itoa(authorMap[key])+"'], "
			}
			authors = authors[0:len(authors)-2]+"]"

			x, _ := getCommits(uuidRepo)
			for _, commit := range x {
				checkoutCommit(uuidRepo, commit)
				fmt.Println("number of authors:", countAuthorsByCommits(uuidRepo, commit))
				fmt.Println("number of files;", countFiles(uuidRepo))
				fmt.Println("langs by files:", countLanguages(uuidRepo))
				files := getFiles(uuidRepo)
				for _, file := range files {
					fmt.Println("File:", file, ":", countLines(file))
				}
			}
			checkoutCommit(uuidRepo, x[0])

			chErr := os.Chdir(blessedPtr)
			check(chErr)
			updateData("dashboards/"+dashboardPtr+"/dashboard.js", "languages", languages)
			updateData("dashboards/"+dashboardPtr+"/dashboard.js", "languageLines", languageLines)
			updateData("dashboards/"+dashboardPtr+"/dashboard.js", "authors", authors)

			// XXX fill in '[]' with real data
			updateData("dashboards/"+dashboardPtr+"/dashboard.js", "numLanguagesData", "{x:[''],y:['']}")
			updateData("dashboards/"+dashboardPtr+"/dashboard.js", "numLinesData", "{x:[''],y:['']}")
			updateData("dashboards/"+dashboardPtr+"/dashboard.js", "numAuthorsData", "{x:[''],y:['']}")
			updateData("dashboards/"+dashboardPtr+"/dashboard.js", "numFilesData", "{x:[''],y:['']}")

			binary, lookErr := exec.LookPath("node")
			check(lookErr)
			args := []string{"node", "./dashboards/"+dashboardPtr+"/dashboard.js"}
			env := os.Environ()
			execErr := syscall.Exec(binary, args, env)
			check(execErr)
		}
	}

}
