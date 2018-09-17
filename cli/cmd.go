package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

var serviceCmd *cobra.Command
var selfbuildCmd *cobra.Command
var gatewayCmd *cobra.Command
var webserverCmd *cobra.Command

func ServiceCmd() *cobra.Command {
	if serviceCmd != nil {
		return serviceCmd
	}
	var projectName string
	createSubCmd := &cobra.Command{
		Use:   "create",
		Short: `Create a micro service frame in your GOPATH. Use "-n" to assign your project name`,
		Long:  "Command used to create a micro service frame in your GOPATH",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectName == "" {
				return fmt.Errorf("Project name not assigned.")
			}
			return newMicroService(projectName)

		},
	}
	createSubCmd.Flags().StringVarP(&projectName, "name", "n", "", "The name of your project.")
	serviceCmd = &cobra.Command{
		Use:   "service",
		Short: "Command of micro service",
		Long:  "Command of micro service",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	serviceCmd.AddCommand(createSubCmd)
	return serviceCmd
}

func GatewayCmd() *cobra.Command {
	if gatewayCmd != nil {
		return gatewayCmd
	}
	var projectName string
	createSubCmd := &cobra.Command{
		Use:   "create",
		Short: `Create an API gateway frame in your GOPATH. Use "-n" to assign your project name`,
		Long:  "Command used to create an API gateway frame in your GOPATH",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectName == "" {
				return fmt.Errorf("Project name not assigned.")
			}
			return newApiGateway(projectName)

		},
	}
	createSubCmd.Flags().StringVarP(&projectName, "name", "n", "", "The name of your project.")
	gatewayCmd = &cobra.Command{
		Use:   "gateway",
		Short: "Command of API gateway",
		Long:  "Command of API gateway",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	gatewayCmd.AddCommand(createSubCmd)
	return gatewayCmd
}

func WebserverCmd() *cobra.Command {
	if webserverCmd != nil {
		return webserverCmd
	}
	var projectName string
	createSubCmd := &cobra.Command{
		Use:   "create",
		Short: `Create an Web Server frame in your GOPATH. Use "-n" to assign your project name`,
		Long:  "Command used to create a Web Server frame in your GOPATH",
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectName == "" {
				return fmt.Errorf("Project name not assigned.")
			}
			return newWebserver(projectName)

		},
	}
	createSubCmd.Flags().StringVarP(&projectName, "name", "n", "", "The name of your project.")
	webserverCmd = &cobra.Command{
		Use:   "webserver",
		Short: "Command of Web Server",
		Long:  "Command of Web Server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	webserverCmd.AddCommand(createSubCmd)
	return webserverCmd
}

func SelfbuildCmd(version string) *cobra.Command {
	if selfbuildCmd != nil {
		return selfbuildCmd
	}
	selfbuildCmd = &cobra.Command{
		Use: "selfbuild",
		RunE: func(cmd *cobra.Command, args []string) error {
			return selfBuild(version, args)
		},
	}
	return selfbuildCmd
}

func newMicroService(projectName string) error {
	printLogo()
	fmt.Println("Start creating a micro service project ...")
	time.Sleep(1 * time.Second)

	gopath := os.Getenv("GOPATH")
	fileInfo, err := os.Stat(gopath)
	if err != nil {
		return fmt.Errorf("Error: stat GOPATH( %s ) error: %s", gopath, err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("Error: can't open GOPATH(%s) which is not a directory.")
	}
	err = UnpackAssets(MICROSERVICE_GzFile, projectName, filepath.Join(gopath, "src"), "microservice")
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	fmt.Printf("\nSuccessfully create new project [%s] in your GOPATH(%s).\n", projectName, gopath)
	return nil
}

func newApiGateway(projectName string) error {
	printLogo()
	fmt.Println("Start creating an API gateway project ...")
	time.Sleep(1 * time.Second)

	gopath := os.Getenv("GOPATH")
	fileInfo, err := os.Stat(gopath)
	if err != nil {
		return fmt.Errorf("Error: stat GOPATH( %s ) error: %s", gopath, err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("Error: can't open GOPATH(%s) which is not a directory.")
	}
	err = UnpackAssets(APIGATEWAY_GzFile, projectName, filepath.Join(gopath, "src"), "apigateway")
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	fmt.Printf("\nSuccessfully create new project [%s] in your GOPATH(%s).\n", projectName, gopath)
	return nil
}

func newWebserver(projectName string) error {
	printLogo()
	fmt.Println("Start creating an Web Server project ...")
	time.Sleep(1 * time.Second)

	gopath := os.Getenv("GOPATH")
	fileInfo, err := os.Stat(gopath)
	if err != nil {
		return fmt.Errorf("Error: stat GOPATH( %s ) error: %s", gopath, err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("Error: can't open GOPATH(%s) which is not a directory.")
	}
	err = UnpackAssets(WEBSERVER_GzFile, projectName, filepath.Join(gopath, "src"), "webserver")
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	fmt.Printf("\nSuccessfully create new project [%s] in your GOPATH(%s).\n", projectName, gopath)
	return nil
}

// not open for user, just for the maintainer
func selfBuild(version string, args []string) error {
	var err error
	if len(args) < 1 {
		return fmt.Errorf("No argument assigned!")
	}
	switch args[0] {
	case "inject":
		err = InjectAssets()
	case "restore":
		err = RestoreAssets()
	case "version":
		fmt.Println(version)
		return nil
	default:
		return fmt.Errorf("Unsupported argument: %s", args[0])
	}
	if err != nil {
		return err
	}
	return nil
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
