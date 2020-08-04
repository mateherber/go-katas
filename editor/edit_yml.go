package editor

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type ParseState string

const (
	LookingForWorkflow       ParseState = "LookingForWorkflow"
	LookingForCommonWorkflow            = "LookingForCommonWorkflow"
	LookingForSecondScript              = "LookingForSecondScript"
	Skipping                            = "Skipping"
	Finished                            = "Finished"
)

func ModifyYml() {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "*.yml")
	if err != nil {
		log.Fatalf("Cannot create temporary file %s\n", err)
	}
	defer tmpFile.Close()
	resp, err := http.Get("https://raw.githubusercontent.com/bitrise-steplib/steps-ios-auto-provision/master/bitrise.yml")
	if err != nil {
		log.Fatalf("Couldn't execute http request %s\n", err)
	}
	var result strings.Builder
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	parseState := LookingForWorkflow
	foundScript := 0
	indent := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) != 0 {
			switch parseState {
			case LookingForWorkflow:
				if matchesToken(line, "workflows") {
					parseState = LookingForCommonWorkflow
					indent = getIndent(line)
				}
			case LookingForCommonWorkflow:
				currentIndent := getIndent(line)
				if currentIndent <= indent {
					parseState = LookingForWorkflow
				} else {
					if matchesToken(line, "_common") {
						parseState = LookingForSecondScript
						indent = currentIndent
					}
				}
			case LookingForSecondScript:
				currentIndent := getIndent(line)
				if currentIndent <= indent {
					parseState = LookingForCommonWorkflow
				} else {
					if matchesToken(line, "script") {
						foundScript += 1
						if foundScript >= 2 {
							parseState = Skipping
							indent = currentIndent
							continue
						}
					}
				}
			case Skipping:
				currentIndent := getIndent(line)
				if currentIndent <= indent {
					parseState = Finished
				} else {
					continue
				}
			}
		}
		result.WriteString(line)
		result.WriteByte('\n')
	}
	_, _ = tmpFile.Write([]byte(result.String()))
	fmt.Println(tmpFile.Name())
}

func getIndent(line string) int {
	r := regexp.MustCompile("^[\\s]*(.*)$")
	return len(line) - len(r.ReplaceAllString(line, "$1"))
}

func matchesToken(line string, token string) bool {
	match, _ := regexp.MatchString("^[\\s-]*"+token+":.*$", line)
	return match
}
