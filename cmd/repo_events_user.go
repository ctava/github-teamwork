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

var repoEventsCmdName = "repoevents"

// repoEventsCmd prints out events in a repository associated with a user
var repoEventsCmd = &cobra.Command{
	Use:   repoEventsCmdName,
	Short: repoEventsCmdName + " repo user start_day end_day",
	Long:  repoEventsCmdName + ` repo user start_day end_day: prints out events by date, user)`,
	Run: func(cmd *cobra.Command, args []string) {

		githubAuthToken := os.Getenv("GITHUB_ACCESS_TOKEN")
		if githubAuthToken == "" {
			fmt.Println("warning: will be limited to 60 calls per hour without a token")
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

		var createBranchTimeSeriesDataSet []byte
		var pushesTimeSeriesDataSet []byte
		var pullrequestsTimeSeriesDataSet []byte
		//var issueCreateTimeSeriesDataSet []byte
		//var pullrequestreviewsTimeSeriesDataSet []byte
		//var pullrequestreviewcommentsTimeSeriesDataSet []byte
		var deleteBranchTimeSeriesDataSet []byte
		var events []github.RepoEvent
		var err error
		events, err = fetcher.FetchRepoEvents(ctx, repo)
		if err != nil {
			fmt.Println("an error occurred while fetching events. err:", err)
			return
		}
		fmt.Printf("%s,%s,%s \n", "created_date", "handle", "type")
		for _, e := range events {
			if strings.Compare(user, e.Handle) == 0 {
				if strings.Compare(e.CreatedAt, start) != -1 {
					if strings.Compare(e.CreatedAt, end) != 1 {
						fmt.Printf("%s,%s,%s \n", e.CreatedAt, e.Handle, e.Type)
						if e.Type == "CreateEvent" {
							createBranchTimeSeriesDataSet = append(createBranchTimeSeriesDataSet, e.CreatedAt...)
							createBranchTimeSeriesDataSet = append(createBranchTimeSeriesDataSet, "\n"...)
						}
						if e.Type == "PushEvent" {
							pushesTimeSeriesDataSet = append(pushesTimeSeriesDataSet, e.CreatedAt...)
							pushesTimeSeriesDataSet = append(pushesTimeSeriesDataSet, "\n"...)
						}
						if e.Type == "PullRequestEvent" {
							pullrequestsTimeSeriesDataSet = append(pullrequestsTimeSeriesDataSet, e.CreatedAt...)
							pullrequestsTimeSeriesDataSet = append(pullrequestsTimeSeriesDataSet, "\n"...)
						}
						if e.Type == "DeleteEvent" {
							deleteBranchTimeSeriesDataSet = append(deleteBranchTimeSeriesDataSet, e.CreatedAt...)
							deleteBranchTimeSeriesDataSet = append(deleteBranchTimeSeriesDataSet, "\n"...)
						}
					}
				}
			}
		}
		fileRoot := start + "-" + user + "-" + repoEventsCmdName
		fileRoot1 := start + "-" + user + "-" + "createbranch"
		writeDataSetToFile(fileRoot1+".csv", createBranchTimeSeriesDataSet)
		fileRoot2 := start + "-" + user + "-" + "pushes"
		writeDataSetToFile(fileRoot2+".csv", pushesTimeSeriesDataSet)
		fileRoot3 := start + "-" + user + "-" + "pullrequests"
		writeDataSetToFile(fileRoot3+".csv", pullrequestsTimeSeriesDataSet)
		fileRoot4 := start + "-" + user + "-" + "deletebranch"
		writeDataSetToFile(fileRoot4+".csv", deleteBranchTimeSeriesDataSet)
		derr := drawChartWithFourLines(startYear, endYear, startMonth, endMonth, "createbranch", "pushes", "pullrequests", "deletebranch", fileRoot1+".csv", fileRoot2+".csv", fileRoot3+".csv", fileRoot4+".csv", fileRoot+".png")
		if derr != nil {
			fmt.Println("an error occurred while drawing the chart. err:", derr)
			return
		}

	},
}

func init() {
	RootCmd.AddCommand(repoEventsCmd)
	repoEventsCmd.Flags().StringP("repo", "R", "", "repo to search")
	repoEventsCmd.Flags().StringP("user", "U", "", "user to search")
	repoEventsCmd.Flags().StringP("start", "S", "", "user events start day")
	repoEventsCmd.Flags().StringP("end", "E", "", "user events end day")
	repoEventsCmd.MarkFlagRequired("repo")
	repoEventsCmd.MarkFlagRequired("user")
	repoEventsCmd.MarkFlagRequired("start")
	repoEventsCmd.MarkFlagRequired("end")
}
