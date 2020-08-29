package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/vinnieapps/mechanic/internal/fileutil"
	"github.com/vinnieapps/mechanic/internal/tarutil"
)

const examplesDir = "examples/"
const baseOutputDir = "build/output/examples"

var binary string

func main() {
	if err := os.MkdirAll(baseOutputDir, 0700); err != nil {
		log.Fatal("Error while creating output directory.", err)
	}

	cleanOuputDir()
	var absErr error
	binary, absErr = filepath.Abs(filepath.Join(baseOutputDir, "mechanic"))
	if absErr != nil {
		log.Fatal("Error while calculating absolute path for binary.", absErr)
	}

	compileMechanic()

	examples, readDirErr := ioutil.ReadDir(examplesDir)
	if readDirErr != nil {
		log.Fatal("Error while reading examples directory.", readDirErr)
	}

	for _, exampleDir := range examples {
		examplePath := filepath.Join(examplesDir, exampleDir.Name())
		outputDir := filepath.Join(baseOutputDir, exampleDir.Name())
		log.Printf("--- Running example: %s in %s\n", examplePath, outputDir)
		os.Mkdir(outputDir, 0700)

		if err := prepareExample(examplePath, outputDir); err != nil {
			log.Fatal("Error while preparing example directory.\n", err)
		}

		if err := buildDockerImage(examplePath, outputDir); err != nil {
			log.Fatal("Error while creating docker image.\n", err)
		}

		if err := runMechanic(outputDir); err != nil {
			log.Fatal("Error while running mechanic for example.\n", err)
		}

		log.Printf("--- Finished example: %s\n", examplePath)
	}
}

func cleanOuputDir() {
	os.RemoveAll(baseOutputDir)
}

func compileMechanic() {
	log.Print("Compiling mechanic...")
	environment := os.Environ()
	environment = append(environment, "GOOS=linux")
	environment = append(environment, "GOARCH=amd64")

	cmd := exec.Command("go", "build", "-o", binary, "cmd/mechanic/main.go")
	cmd.Env = environment
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if runErr := cmd.Run(); runErr != nil {
		log.Fatal("Error while compiling mechanic.\n", runErr)
	}
}

func buildDockerImage(exampleDir string, outputDir string) error {
	dockerFile := filepath.Join(exampleDir, "Dockerfile")
	if _, err := os.Stat(dockerFile); os.IsNotExist(err) {
		content := `
FROM ubuntu:latest
RUN mkdir -p /opt/test

COPY ./mechanic /usr/bin/mechanic
RUN chown root:root /usr/bin/mechanic && chmod +x /usr/bin/mechanic

WORKDIR /opt/test
CMD ["/usr/bin/mechanic"]
`

		dockerFile = filepath.Join(outputDir, "Dockerfile")
		ioutil.WriteFile(dockerFile, []byte(content), 0700)
	}

	log.Printf("Building Docker image from %s", dockerFile)
	cmd := exec.Command(
		"docker",
		"build",
		"--no-cache",
		"-t", filepath.Base(exampleDir),
		"-f", dockerFile,
		outputDir,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func prepareExample(exampleDir string, outputDir string) error {
	packageDir := filepath.Join(exampleDir, "package")
	log.Printf("Creating package for example from directory: %s", packageDir)

	log.Printf("Copying binary to output directory")
	fileutil.CopyFile(binary, filepath.Join(outputDir, "mechanic"))
	return tarutil.CreateTarFromDir(filepath.Join(outputDir, "package.tar.gz"), packageDir)
}

func runMechanic(outputDir string) error {
	exampleName := filepath.Base(outputDir)
	log.Printf("------------ Starting mechanic for %s", exampleName)
	defer log.Printf("------------ Finished mechanic for %s", exampleName)

	absOutput, absErr := filepath.Abs(outputDir)
	if absErr != nil {
		return absErr
	}

	cmd := exec.Command("docker", "run", fmt.Sprintf("--volume=%s:/opt/test", absOutput), exampleName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
