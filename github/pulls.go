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
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

func (s *fetcher) FetchPullRequestComments(ctx context.Context, repositoryURL string) ([]PullComment, error) {

	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	url, err := url.Parse(repositoryURL)
	if err != nil {
		return nil, err
	}
	values := strings.Split(url.Path, "/")
	if len(values) < 3 {
		return nil, errors.New("invalid repository url")
	}
	owner, repo := values[1], values[2]

	listOpts := github.PullRequestListCommentsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var pullRequestComments []*github.PullRequestComment
	var pullComments []PullComment
	var resp *github.Response
	for {
		pullRequestComments, resp, err = s.client.PullRequests.ListComments(ctx, owner, repo, 0, &listOpts)
		if err != nil {
			return nil, err
		}

		var pullComment PullComment
		var id int64
		var body string
		var handle string
		var reactionTotalCount int
		var reactionPlusOne int
		var reactionMinusOne int
		var reactionLaugh int
		var reactionConfused int
		var reactionHeart int
		var reactionHooray int
		var commentCreatedAt string
		var time time.Time
		for _, prc := range pullRequestComments {
			id = *prc.ID
			body = *prc.Body
			handle = *prc.User.Login
			reactionTotalCount = *prc.Reactions.TotalCount
			reactionPlusOne = *prc.Reactions.PlusOne
			reactionMinusOne = *prc.Reactions.MinusOne
			reactionLaugh = *prc.Reactions.Laugh
			reactionConfused = *prc.Reactions.Confused
			reactionHeart = *prc.Reactions.Heart
			reactionHooray = *prc.Reactions.Hooray
			time = *prc.CreatedAt
			commentCreatedAt = time.Format("2006-01-02")
			pullComment = PullComment{ID: id, Body: body, Handle: handle, CreatedAt: commentCreatedAt,
				ReactionTotalCount: reactionTotalCount, ReactionPlusOne: reactionPlusOne,
				ReactionMinusOne: reactionMinusOne, ReactionLaugh: reactionLaugh,
				ReactionConfused: reactionConfused, ReactionHeart: reactionHeart, ReactionHooray: reactionHooray}
			pullComments = append(pullComments, pullComment)
		}

		if resp.NextPage == 0 {
			break
		}
		listOpts.Page = resp.NextPage
	}

	return pullComments, nil
}
