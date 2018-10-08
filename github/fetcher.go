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

package github

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//NewFetcher public function to create client for interfacing with github.com API
func NewFetcher(ctx context.Context, token string) Fetcher {
	if token == "" {
		return &fetcher{
			client: github.NewClient(nil),
		}
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	newClient := oauth2.NewClient(ctx, ts)
	client := github.NewClient(newClient)
	return &fetcher{
		client: *client,
	}
}

type fetcher struct {
	client *github.Client
}

//Fetcher public functions interfacing with github.com API
type Fetcher interface {
	FetchPullRequestComments(ctx context.Context, repositoryURL string) ([]PullComment, error)
	FetchRepoEvents(ctx context.Context, repositoryURL string) ([]RepoEvent, error)
	FetchTeamDiscussionComments(ctx context.Context, org, teamName string) ([]DiscussionComment, error)
}

// PullComment a struct for local, simplified representation of a PullRequestComment
type PullComment struct {
	Handle             string
	ID                 int64
	Body               string
	ReactionTotalCount int
	ReactionPlusOne    int
	ReactionMinusOne   int
	ReactionLaugh      int
	ReactionConfused   int
	ReactionHeart      int
	ReactionHooray     int
	CreatedAt          string
}

// RepoEvent a struct for local, simplified representation of an RepoEvent for a repository
type RepoEvent struct {
	Handle    string
	Repo      string
	Type      string
	Payload   string
	CreatedAt string
}

// DiscussionComment a struct for local, simplified representation of a DiscussionComment
type DiscussionComment struct {
	Handle             string
	ID                 int64
	Title              string
	Body               string
	ReactionTotalCount int
	ReactionPlusOne    int
	ReactionMinusOne   int
	ReactionLaugh      int
	ReactionConfused   int
	ReactionHeart      int
	ReactionHooray     int
	CreatedAt          string
}
