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
	"log"

	"github.com/mniak/ytlive/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticates on youtube",
	Run: func(cmd *cobra.Command, args []string) {
		clientID, err := cmd.Flags().GetString("client-id")
		if err != nil {
			log.Fatalln(err)
		}

		clientSecret, err := cmd.Flags().GetString("client-secret")
		if err != nil {
			log.Fatalln(err)
		}

		err = pkg.Login(clientID, clientSecret)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().String("client-id", "", "the Youtube Client ID")
	loginCmd.MarkFlagRequired("client-id")
	viper.BindPFlag("Application.ClientID", loginCmd.Flags().Lookup("client-id"))

	loginCmd.Flags().String("client-secret", "", "the Youtube Client Secret")
	loginCmd.MarkFlagRequired("client-secret")
	viper.BindPFlag("Application.ClientSecret", loginCmd.Flags().Lookup("client-secret"))
}
