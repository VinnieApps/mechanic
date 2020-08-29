package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/vinnieapps/mechanic/internal/tarutil"
	"github.com/vinnieapps/mechanic/internal/work"
)

func main() {
	start := time.Now()
	abs, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current dir is: %s\n", abs)
	flag.Parse()

	args := flag.Args()
	if len(args) > 1 {
		printUsage()
		os.Exit(1)
	}

	packageName := "package.tar.gz"
	if len(args) == 1 {
		packageName = args[0]
	}

	if _, err := os.Stat(packageName); os.IsNotExist(err) {
		fmt.Printf("Sorry, I can't find package: %s\n", packageName)
		fmt.Println("Exiting now...Bye!")
		os.Exit(2)
	}

	fmt.Printf("Using package: %s\n", packageName)
	if err := processPackage(packageName); err != nil {
		fmt.Printf("Error while processing package. %s\n", err)
		os.Exit(3)
	}

	fmt.Printf("I'm done! Finished in %dms.\n", time.Now().Sub(start).Milliseconds())
	fmt.Println("See you next time!")
}

// getTaskFileNames returns a list with all the files containing tasks to
// be executed, sorted by order they should be processed
func getTaskFileNames(fileNames []string) []string {
	tasks := make([]string, 0)
	for _, fileName := range fileNames {
		if strings.HasPrefix(fileName, "tasks/") {
			tasks = append(tasks, fileName)
		}
	}
	sort.Strings(tasks)
	return tasks
}

func printUsage() {
	fmt.Println("usage:")
	fmt.Println()
	fmt.Println("  mechanic [package]")
	fmt.Println()
	fmt.Println("  package   Name of the package to process. Defaults to packages.tar.gz")
	fmt.Println()
}

func processPackage(packageName string) error {
	outputDir := "./output"
	files, untarErr := tarutil.UntarFile(packageName, outputDir)
	if untarErr != nil {
		return untarErr
	}

	for _, taskFile := range getTaskFileNames(files) {
		tasksToExecute, deserializeErr := work.Deserialize(filepath.Join(outputDir, taskFile))
		if deserializeErr != nil {
			return deserializeErr
		}

		for _, task := range tasksToExecute {
			if err := task.Task.Execute(outputDir); err != nil {
				return err
			}
		}
	}

	return nil
}
