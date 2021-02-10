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

	"github.com/mniak/generate-streams/youtube"
	"github.com/spf13/cobra"
)

// cleanupCmd represents the youtube command
var cleanupCmd = &cobra.Command{
	Use: "cleanup",
	Aliases: []string{
		"clear",
		"purge",
	},
	Short: "Schedule a new youtube live stream",
	Run: func(cmd *cobra.Command, args []string) {

		since := time.Now().Add(7 * 24 * time.Hour * -1)
		cleaned, err := youtube.CleanupStreams(since)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Cleaned up old streams")
		for _, streamName := range cleaned {
			fmt.Printf("  - %s", streamName)
		}
	},
}

func init() {
	youtubeCmd.AddCommand(cleanupCmd)
}
