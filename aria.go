package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var commands = map[string]func(args []string){
	"new":       newMicroService,
	"selfbuild": selfBuild,
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		return
	}
	if cmd, ok := commands[args[0]]; ok {
		cmd(args[1:])
	} else {
		fmt.Printf("Command not support: %s\n", args[0])
		printUsage()
		return
	}
}

func printUsage() {
	fmt.Println(`Usage: 
    aria command [arguments]

Available commands:
    new    Create a new project of aria
`)
}

func newMicroService(args []string) {
	printLogo()
	fmt.Println("Start creating a micro service project ...")
	time.Sleep(2 * time.Second)
	if len(args) < 1 {
		exitWithError(fmt.Errorf("Error: project name not assigned. Use:\n    aria new <project_name>"))
	}
	projectName := args[0]
	gopath := os.Getenv("GOPATH")
	fileInfo, err := os.Stat(gopath)
	if err != nil {
		exitWithError(fmt.Errorf("Error: stat GOPATH( %s ) error: %s", gopath, err))
	}
	if !fileInfo.IsDir() {
		exitWithError(fmt.Errorf("Error: can't open GOPATH(%s) which is not a directory."))
	}
	err = UnpackAssets(MICROSERVICE_GzFile, projectName, filepath.Join(gopath, "src"), "microservice")
	if err != nil {
		exitWithError(fmt.Errorf("Error: %s", err))
	}
	fmt.Printf("\nSuccessfully create new project [%s] in your GOPATH(%s).\n", projectName, gopath)
}

// not open for user, just for the maintainer
func selfBuild(args []string) {
	var err error
	if len(args) < 1 {
		exitWithError(fmt.Errorf("No argument assigned!"))
	}
	switch args[0] {
	case "inject":
		err = InjectAssets()
	case "restore":
		err = RestoreAssets()
	default:
		exitWithError(fmt.Errorf("Unsupported argument: %s", args[0]))
	}
	if err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func printLogo() {
	fmt.Println(`
   _____             .___          
  /  _  \   _______  |   | _____   
 /  /_\  \  \_  __ \ |   | \__  \  
/    |    \  |  | \/ |   |  / __ \_
\____|__  /  |__|    |___| (____  /
        \/                      \/
`)
}