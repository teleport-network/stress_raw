/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	filePath string
	csvDir   string
	startNum int
	times    int
	rootCmd  = &cobra.Command{
		Use:   "tst",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&filePath, "filePath", "f", "", "privateKey's csv file(required)")
	rootCmd.PersistentFlags().StringVarP(&csvDir, "csvDir", "d", "./rawTxns", "txn's csv folder")
	rootCmd.PersistentFlags().IntVarP(&startNum, "start", "s", 1, "privateKey's start num")
	rootCmd.PersistentFlags().IntVarP(&times, "times", "t", 100, "run times")
	// gwtCmd.MarkFlagRequired("filePath")
}
