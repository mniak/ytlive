package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/mniak/ytlive/pkg"
	"github.com/spf13/cobra"
)

// cleanupKeysCmd represents the cleanup-keys command
var cleanupKeysCmd = &cobra.Command{
	Use:   "cleanup-keys",
	Short: "Cleanup old stream keys",
	Run: func(cmd *cobra.Command, args []string) {

		since := time.Now().Add(7 * 24 * time.Hour * -1)
		cleaned, err := pkg.CleanupKeys(since)
		if err != nil {
			log.Fatalln(err)
		}

		if len(cleaned) > 0 {
			fmt.Println("Stream keys cleaned:")
			for _, streamName := range cleaned {
				fmt.Printf("  - %s", streamName)
			}
		} else {
			fmt.Println("There was none old stream keys")
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanupKeysCmd)
}
