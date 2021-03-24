package cmd

import (
	"fmt"
	"log"

	"github.com/araddon/dateparse"
	"github.com/mniak/ytlive/pkg"
	"github.com/spf13/cobra"
)

// scheduleCmd represents the schedule command
var scheduleCmd = &cobra.Command{
	Use:     "schedule <title> <date> <time>",
	Aliases: []string{"new"},
	Short:   "Schedule a new youtube live stream",
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]

		parsedDateTime, err := dateparse.ParseLocal(args[1] + " " + args[2])
		if err != nil {
			log.Fatalf("invalid datetime: %s %s\n", args[1], args[2])
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
		request := pkg.ScheduleRequest{
			Title:     title,
			Date:      parsedDateTime,
			AutoStart: autoStart,
			AutoStop:  autoStop,
			DVR:       dvr,
		}
		response, err := pkg.Schedule(request)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Youtube Stream Scheduled:")
		fmt.Printf("  ID:         %s\n", response.ID)
		fmt.Printf("  Title:      %s\n", response.Title)
		fmt.Printf("  Date:       %s\n", response.Date.Local())
		fmt.Printf("  Link:       %s\n", response.Link)
		fmt.Printf("  Stream Name:   %s\n", response.StreamName)
		fmt.Printf("  Stream Key: %s\n", response.StreamKey)
		fmt.Printf("  Stream URL: %s\n", response.StreamURL)
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)

	scheduleCmd.Flags().Bool("auto-start", false, "enable auto-start")
	scheduleCmd.Flags().Bool("auto-stop", false, "enable auto-stop")
	scheduleCmd.Flags().Bool("dvr", false, "enable DVR")
}
