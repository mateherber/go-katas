package main

import . "github.com/mateherber/go-katas/artifact"
import . "github.com/mateherber/go-katas/editor"
import . "github.com/mateherber/go-katas/version"
import . "github.com/mateherber/go-katas/simulator"

func main() {
	DownloadBuildArtifact()
	ModifyYml()
	PrintStepVersion()
	StartIosSimulator()
}
