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
	"fmt"
	"strings"
	"time"

	"context"
	"os"

	"github.com/ctava/github-teamwork/github"
	"github.com/spf13/cobra"
)

var discussionCmdName = "teamdiscussion"

// discussionCmd prints out contributions to dicussions
var discussionCmd = &cobra.Command{
	Use:   discussionCmdName,
	Short: discussionCmdName + " org team user startDay endDay",
	Long:  discussionCmdName + ` org team user startDay endDay: prints out team discussion comments by date, user. includes reactions (total count, :+1:, :-1:, :laugh:, :confused:, :heart: and :hooray:)`,
	Run: func(cmd *cobra.Command, args []string) {

		githubAuthToken := os.Getenv("GITHUB_ACCESS_TOKEN")
		if githubAuthToken == "" {
			fmt.Println("warning: will be limited to 60 calls per hour without a token")
		}
		ctx := context.Background()
		fetcher := github.NewFetcher(ctx, githubAuthToken)

		team := getFlagString(cmd, "team")
		values := strings.Split(team, "/")
		if len(values) < 2 {
			fmt.Println("error: team name needs to be owner/teamname")
			return
		}
		org, teamName := values[0], values[1]

		user := getFlagString(cmd, "user")
		start := getFlagString(cmd, "start")
		end := getFlagString(cmd, "end")
		startTime, sterr := time.Parse("2006-01-02", start)
		if sterr != nil {
			fmt.Println("an error occurred while parsing the start time. err:", sterr)
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

		discussionComments, err := fetcher.FetchTeamDiscussionComments(ctx, org, teamName)
		if err != nil {
			fmt.Println("an error occurred while fetching PR Comments err:", err)
			return
		}
		var timeSeriesDataSet []byte
		fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s \n", "created_date", "handle", "body", "reaction_total_count", "reaction_plusone", "reaction_minusone", "reaction_laugh", "reaction_confused", "reaction_heart", "reaction_hooray")
		for _, c := range discussionComments {
			if strings.Compare(user, c.Handle) == 0 {
				if strings.Compare(c.CreatedAt, start) != -1 {
					if strings.Compare(c.CreatedAt, end) != 1 {
						fmt.Printf("%s,%s,%s,%v,%v,%v,%v,%v,%v,%v \n", c.CreatedAt, c.Handle, c.Body, c.ReactionTotalCount, c.ReactionPlusOne, c.ReactionMinusOne, c.ReactionLaugh, c.ReactionConfused, c.ReactionHeart, c.ReactionHooray)
						timeSeriesDataSet = append(timeSeriesDataSet, c.CreatedAt...)
						timeSeriesDataSet = append(timeSeriesDataSet, "\n"...)
					}
				}
			}
		}
		fileRoot := start + "-" + user + "-" + discussionCmdName
		writeDataSetToFile(fileRoot+".csv", timeSeriesDataSet)
		derr := drawChart(startYear, endYear, startMonth, endMonth, discussionCmdName, fileRoot+".csv", fileRoot+".png")
		if derr != nil {
			fmt.Println("an error occurred while drawing the chart. err:", derr)
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(discussionCmd)
	discussionCmd.Flags().StringP("team", "T", "", "team to search for discussion threads")
	discussionCmd.Flags().StringP("user", "U", "", "commenter to search for")
	discussionCmd.Flags().StringP("start", "S", "", "comment start day")
	discussionCmd.Flags().StringP("end", "E", "", "comment end day")
	discussionCmd.MarkFlagRequired("team")
	discussionCmd.MarkFlagRequired("user")
	discussionCmd.MarkFlagRequired("start")
	discussionCmd.MarkFlagRequired("end")
}
