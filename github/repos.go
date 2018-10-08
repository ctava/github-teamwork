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

func (s *fetcher) FetchRepoEvents(ctx context.Context, repositoryURL string) ([]RepoEvent, error) {

	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	if repositoryURL == "" {
		return nil, errors.New("repositoryURL is nil")
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

	listOpts := github.ListOptions{PerPage: 30}

	var githubEvents []*github.Event
	var events []RepoEvent
	var resp *github.Response
	for {
		githubEvents, resp, err = s.client.Activity.ListRepositoryEvents(ctx, owner, repo, &listOpts)
		if err != nil {
			return nil, err
		}
		var event RepoEvent
		var repo string
		var actor string
		var eventType string
		var createdDateAt time.Time
		var createdAt string
		for _, e := range githubEvents {
			actor = *e.Actor.Login
			repo = *e.Repo.Name
			eventType = *e.Type
			createdDateAt = *e.CreatedAt
			createdAt = createdDateAt.Format("2006-01-02")
			event = RepoEvent{Handle: actor, Type: eventType, CreatedAt: createdAt, Repo: repo}
			events = append(events, event)
		}
		if resp.NextPage == 0 {
			break
		}

		listOpts.Page = resp.NextPage

	}

	return events, nil
}
