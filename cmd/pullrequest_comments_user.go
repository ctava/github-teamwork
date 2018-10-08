// Copyright © 2018 Chris Tava <chris1tava@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ctava/github-teamwork/github"
	"github.com/spf13/cobra"
)

var pullrequestCommentsCmdName = "prcomments"

// pullrequestCommentsCmd prints out pull request comments and reactions
var pullrequestCommentsCmd = &cobra.Command{
	Use:   pullrequestCommentsCmdName,
	Short: pullrequestCommentsCmdName + " repo user start_day end_day",
	Long:  pullrequestCommentsCmdName + ` repo user start_day end_day: prints out pull request comments by date, user. includes reactions (total count, :+1:, :-1:, :laugh:, :confused:, :heart: and :hooray:)`,
	Run: func(cmd *cobra.Command, args []string) {

		githubAuthToken := os.Getenv("GITHUB_ACCESS_TOKEN")
		if githubAuthToken == "" {
			fmt.Println("warning: without a token, you will be limited to 60 calls per hour")
		}
		ctx := context.Background()
		fetcher := github.NewFetcher(ctx, githubAuthToken)

		repo := getFlagString(cmd, "repo")
		user := getFlagString(cmd, "user")

		start := getFlagString(cmd, "start")
		end := getFlagString(cmd, "end")
		startTime, sterr := time.Parse("2006-01-02", start)
		if sterr != nil {
			fmt.Println("an error occurred while parsing the end time. err:", sterr)
			return
		}
		startYear := startTime.Year()
		startMonth := startTime.Month()
		endTime, eterr := time.Parse("2006-01-02", end)
		if eterr != nil {
			fmt.Println("an error occurred while parsing the end time. err:", eterr)
			return
		}
		endYear := endTime.Year()
		endMonth := endTime.Month()

		var prComments []github.PullComment
		prComments, ferr := fetcher.FetchPullRequestComments(ctx, repo)
		if ferr != nil {
			fmt.Println("an error occurred while fetching PR Comments. err:", ferr)
			return
		}
		var filteredPRComments []github.PullComment
		var timeSeriesDataSet []byte
		fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s \n", "created_date", "handle", "body", "reaction_total_count", "reaction_plusone", "reaction_minusone", "reaction_laugh", "reaction_confused", "reaction_heart", "reaction_hooray")
		for _, c := range prComments {
			if strings.Compare(user, c.Handle) == 0 {
				if strings.Compare(c.CreatedAt, start) != -1 {
					if strings.Compare(c.CreatedAt, end) != 1 {
						fmt.Printf("%s,%s,%s,%v,%v,%v,%v,%v,%v,%v \n", c.CreatedAt, c.Handle, c.Body, c.ReactionTotalCount, c.ReactionPlusOne, c.ReactionMinusOne, c.ReactionLaugh, c.ReactionConfused, c.ReactionHeart, c.ReactionHooray)
						filteredPRComments = append(filteredPRComments, c)
						timeSeriesDataSet = append(timeSeriesDataSet, c.CreatedAt...)
						timeSeriesDataSet = append(timeSeriesDataSet, "\n"...)
					}
				}
			}
		}
		fileRoot := start + "-" + user + "-" + pullrequestCommentsCmdName
		writeDataSetToFile(fileRoot+".csv", timeSeriesDataSet)
		derr := drawChart(startYear, endYear, startMonth, endMonth, pullrequestCommentsCmdName, fileRoot+".csv", fileRoot+".png")
		if derr != nil {
			fmt.Println("an error occurred while drawing the chart. err:", derr)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(pullrequestCommentsCmd)
	pullrequestCommentsCmd.Flags().StringP("repo", "R", "", "repo to search for pull request comments")
	pullrequestCommentsCmd.Flags().StringP("user", "U", "", "pull request commenter to search for")
	pullrequestCommentsCmd.Flags().StringP("start", "S", "", "pull request comment start day")
	pullrequestCommentsCmd.Flags().StringP("end", "E", "", "pull request comment end day")
	pullrequestCommentsCmd.MarkFlagRequired("repo")
	pullrequestCommentsCmd.MarkFlagRequired("user")
	pullrequestCommentsCmd.MarkFlagRequired("start")
	pullrequestCommentsCmd.MarkFlagRequired("end")
}
