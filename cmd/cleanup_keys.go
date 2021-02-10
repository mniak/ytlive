/*
Copyright © 2021 Andre Soares

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
