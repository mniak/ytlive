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

	"github.com/araddon/dateparse"
	"github.com/mniak/generate-streams/youtube"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// youtubeCmd represents the youtube command
var youtubeCmd = &cobra.Command{
	Use:     "youtube <title> <date> <time>",
	Aliases: []string{"yt"},
	Short:   "Schedule a new youtube live stream",
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]

		parsedDateTime, err := dateparse.ParseAny(args[1] + " " + args[2])
		if err != nil {
			log.Fatalf("invalid datetime: %s %s\n", args[1], args[2])
		}

		keyName, err := cmd.Flags().GetString("key-name")
		if err != nil {
			log.Fatalln(err)
		}
		newKey, err := cmd.Flags().GetBool("new-key")
		if err != nil {
			log.Fatalln(err)
		}
		autoStart, err := cmd.Flags().GetBool("auto-start")
		if err != nil {
			log.Fatalln(err)
		}
		autoStop, err := cmd.Flags().GetBool("auto-stop")
		if err != nil {
			log.Fatalln(err)
		}
		dvr, err := cmd.Flags().GetBool("dvr")
		if err != nil {
			log.Fatalln(err)
		}

		request := youtube.GenerateRequest{
			Title:         title,
			Date:          parsedDateTime,
			StreamKeyName: keyName,
			NewStreamKey:  newKey,
			AutoStart:     autoStart,
			AutoStop:      autoStop,
			DVR:           dvr,
		}
		response, err := youtube.Generate(request)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Youtube Stream Generated:")
		fmt.Printf("  ID:         %s\n", response.ID)
		fmt.Printf("  Title:      %s\n", response.Title)
		fmt.Printf("  Date:       %s\n", response.Date)
		fmt.Printf("  Link:       %s\n", response.Link)
		fmt.Printf("  Key Name:   %s\n", response.StreamKeyName)
		fmt.Printf("  Stream Key: %s\n", response.StreamKey)
		fmt.Printf("  Stream URL: %s\n", response.StreamURL)

	},
}

func init() {
	rootCmd.AddCommand(youtubeCmd)

	youtubeCmd.AddCommand(youtubeCleanupCmd)

	youtubeCmd.PersistentFlags().String("client-id", "", "the Youtube Client ID")
	viper.BindPFlag("Youtube.ClientID", youtubeCmd.PersistentFlags().Lookup("client-id"))

	youtubeCmd.PersistentFlags().String("client-secret", "", "the Youtube Client Secret")
	viper.BindPFlag("Youtube.ClientSecret", youtubeCmd.PersistentFlags().Lookup("client-secret"))

	youtubeCmd.Flags().String("key-name", "", "select stream key by name")
	youtubeCmd.Flags().Bool("new-key", false, "create new stream key")
	youtubeCmd.Flags().Bool("auto-start", false, "enable auto-start")
	youtubeCmd.Flags().Bool("auto-stop", false, "enable auto-stop")
	youtubeCmd.Flags().Bool("dvr", false, "enable DVR")
}
