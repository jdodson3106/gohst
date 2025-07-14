package cmd

import (
	"os"

	"github.com/jdodson3106/gohst"
	"github.com/spf13/cobra"
)

var runPlayground bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gohst",
	Short: "Used to view and manage your env parameters and secrets in your local environment",
	Long: `Gohst is a env and secrets manager for managing key rotations both locally 
         and in your AWS account. 

         The tool facilitates the generation of env files and secure sharing of secrets with 
         other users on your team.`,

	Run: func(cmd *cobra.Command, args []string) {
		if runPlayground {
			gohst.Playground()
		} else {
			gohst.Start()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&runPlayground, "playground", "p", false, "Runs test code for development purposes")
}
