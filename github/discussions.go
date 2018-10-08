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

	"github.com/google/go-github/github"
)

func (s *fetcher) FetchTeamDiscussionComments(ctx context.Context, org, teamName string) ([]DiscussionComment, error) {

	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	listOpts := github.ListOptions{PerPage: 30}
	var teamID int64
	for {
		teams, resp, err := s.client.Teams.ListTeams(ctx, org, &listOpts)
		if err != nil {
			return nil, err
		}
		for _, t := range teams {
			if t.GetName() == teamName {
				teamID = t.GetID()
				break
			}
		}
		if resp.NextPage == 0 {
			break
		}
		listOpts.Page = resp.NextPage
	}
	if teamID == 0 {
		return nil, errors.New("TeamID is missing")
	}

	teamdiscussions, _, err := s.client.Teams.ListDiscussions(ctx, teamID, nil)
	if err != nil {
		return nil, err
	}
	var discussionComments []DiscussionComment
	var discussionComment DiscussionComment
	for _, td := range teamdiscussions {
		dcs, _, err := s.client.Teams.ListComments(ctx, teamID, *td.Number, nil)
		if err != nil {
			return nil, err
		}
		var body string
		var handle string
		var reactionTotalCount int
		var reactionPlusOne int
		var reactionMinusOne int
		var reactionLaugh int
		var reactionConfused int
		var reactionHeart int
		var reactionHooray int
		var time github.Timestamp
		var createdAt string
		var reactions github.Reactions
		for _, dc := range dcs {
			handle = *dc.Author.Login
			body = *dc.Body
			if reactions != (github.Reactions{}) {
				reactionTotalCount = *dc.Reactions.TotalCount
				reactionPlusOne = *dc.Reactions.PlusOne
				reactionMinusOne = *dc.Reactions.MinusOne
				reactionLaugh = *dc.Reactions.Laugh
				reactionConfused = *dc.Reactions.Confused
				reactionHeart = *dc.Reactions.Heart
				reactionHooray = *dc.Reactions.Hooray
			}
			time = *dc.CreatedAt
			createdAt = time.Format("2006-01-02")
			discussionComment = DiscussionComment{Body: body, Handle: handle, CreatedAt: createdAt, ReactionTotalCount: reactionTotalCount, ReactionPlusOne: reactionPlusOne, ReactionMinusOne: reactionMinusOne, ReactionLaugh: reactionLaugh, ReactionConfused: reactionConfused, ReactionHeart: reactionHeart, ReactionHooray: reactionHooray}
			discussionComments = append(discussionComments, discussionComment)
		}
	}

	return discussionComments, nil
}
