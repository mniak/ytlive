package cmd

import (
	"fmt"
	"os"

	"github.com/mniak/ytlive/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ytlive",
	Short: "A tool to manage YouTube live streams",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if _, err := config.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
