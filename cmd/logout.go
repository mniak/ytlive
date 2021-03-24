package cmd

import (
	"log"

	"github.com/mniak/ytlive/pkg"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:     "logout",
	Aliases: []string{"deauth", "deauthenticate"},
	Short:   "Erase the Youtube token",
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.Logout()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
