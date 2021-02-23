/*
Copyright © 2021 koh23<53048847+koh23@users.noreply.github.com >

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
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stream called")
		consumerKey := viper.GetString("consumer_key")
		consumerSecret := viper.GetString("consumer_secret")
		oauthToken := viper.GetString("oauth_token")
		oauthTokenSecret := viper.GetString("oauth_token_secret")
		config := oauth1.NewConfig(consumerKey, consumerSecret)
		token := oauth1.NewToken(oauthToken, oauthTokenSecret)
		httpClient := config.Client(oauth1.NoContext, token)
		client := twitter.NewClient(httpClient)

		params := &twitter.StreamSampleParams{
			StallWarnings: twitter.Bool(true),
		}
		stream, err := client.Streams.Sample(params)
		if err != nil {
			fmt.Println("Client Error!")
		}
		// receive message
		demux := twitter.NewSwitchDemux()
		demux.Tweet = func(tweet *twitter.Tweet) {
			//fmt.Println("[INFO]", tweet.Text)
			//fmt.Println("[INFO]", tweet.Lang)
			ext := tweet.ExtendedEntities
			if ext != nil {
				for _, media := range ext.Media {
					fmt.Println(media.MediaURLHttps)
				}
			}
		}
		go demux.HandleChan(stream.Messages)

		// Wait for SIGINT and SIGTERM (HIT CTRL-C)
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		fmt.Println(<-ch)
		stream.Stop()
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// streamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// streamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
