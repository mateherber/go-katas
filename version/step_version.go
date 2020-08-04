package version

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func PrintStepVersion() {
	fmt.Println("Please enter a step's name to find out it's latest version:")
	reader := bufio.NewReader(os.Stdin)
	stepName, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Couldn't read step's name %s\n", err)
	}
	stepName = strings.TrimSpace(stepName)
	resp, err := http.Get("https://bitrise-steplib-collection.s3.amazonaws.com/spec.json")
	if err != nil {
		log.Fatalf("Couldn't execute http request %s\n", err)
	}
	defer resp.Body.Close()
	byteResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Couldn't read request body %s\n", err)
	}
	steps, err := UnmarshalSteps(byteResponse)
	if err != nil {
		log.Fatalf("Couldn't parse steps json %s\n", err)
	}
	for key, step := range steps.Steps {
		if key == stepName {
			fmt.Println(step.LatestVersion)
			return
		}
	}
	log.Fatalf("Couldn't find step %s\n", stepName)
}
