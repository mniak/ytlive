package cmd

import (
	"log"

	"github.com/mniak/ytlive/config"
	"github.com/mniak/ytlive/pkg"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"auth", "authenticate"},
	Short:   "Authenticates on Youtube",
	Run: func(cmd *cobra.Command, args []string) {
		clientID, err := cmd.Flags().GetString("client-id")
		if err != nil {
			log.Fatalln(err)
		}

		clientSecret, err := cmd.Flags().GetString("client-secret")
		if err != nil {
			log.Fatalln(err)
		}

		config.Root.Application.ClientID = clientID
		config.Root.Application.ClientSecret = clientSecret
		err = pkg.Login()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().String("client-id", "", "the Youtube Client ID")
	loginCmd.MarkFlagRequired("client-id")

	loginCmd.Flags().String("client-secret", "", "the Youtube Client Secret")
	loginCmd.MarkFlagRequired("client-secret")
}
