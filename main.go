package main

import (
	"aria/cli"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const APP_VERSION = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "aria",
		Short: "Generate micro service project frame",
		Long: `
   _____             .___          
  /  _  \   _______  |   | _____   
 /  /_\  \  \_  __ \ |   | \__  \  
/    |    \  |  | \/ |   |  / __ \_
\____|__  /  |__|    |___| (____  /
        \/                      \/
Aria is a tool for generating micro 
service project frame based on go-kit`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Name() == cli.SelfbuildCmd(APP_VERSION).Name() {
				if err := cli.SelfbuildCmd(APP_VERSION).RunE(cmd, args); err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
				os.Exit(0)
			} else {
				cmd.RemoveCommand(cli.SelfbuildCmd(APP_VERSION))
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		Version: APP_VERSION,
	}

	rootCmd.AddCommand(cli.ServiceCmd(), cli.GatewayCmd(), cli.SelfbuildCmd(APP_VERSION))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
