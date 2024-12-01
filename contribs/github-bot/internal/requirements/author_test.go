package requirements

import (
	"context"
	"fmt"
	"testing"

	"github-bot/internal/client"
	"github-bot/internal/logger"
	"github-bot/internal/utils"
	"github.com/stretchr/testify/assert"

	"github.com/google/go-github/v64/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/xlab/treeprint"
)

func TestAuthor(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		name        string
		user        string
		author      string
		isSatisfied bool
	}{
		{"author match", "user", "user", true},
		{"author doesn't match", "user", "author", false},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			pr := &github.PullRequest{
				User: &github.User{Login: github.String(testCase.author)},
			}
			details := treeprint.New()
			requirement := Author(testCase.user)

			assert.Equal(t, requirement.IsSatisfied(pr, details), testCase.isSatisfied, fmt.Sprintf("requirement should have a satisfied status: %t", testCase.isSatisfied))
			assert.True(t, utils.TestLastNodeStatus(t, testCase.isSatisfied, details), fmt.Sprintf("requirement details should have a status: %t", testCase.isSatisfied))
		})
	}
}

func TestAuthorInTeam(t *testing.T) {
	t.Parallel()

	members := []*github.User{
		{Login: github.String("notTheRightOne")},
		{Login: github.String("user")},
		{Login: github.String("anotherOne")},
	}

	for _, testCase := range []struct {
		name        string
		user        string
		members     []*github.User
		isSatisfied bool
	}{
		{"empty member list", "user", []*github.User{}, false},
		{"member list contains user", "user", members, true},
		{"member list doesn't contain user", "user2", members, false},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mockedHTTPClient := mock.NewMockedHTTPClient(
				mock.WithRequestMatchPages(
					mock.EndpointPattern{
						Pattern: "/orgs/teams/team/members",
						Method:  "GET",
					},
					testCase.members,
				),
			)

			gh := &client.GitHub{
				Client: github.NewClient(mockedHTTPClient),
				Ctx:    context.Background(),
				Logger: logger.NewNoopLogger(),
			}

			pr := &github.PullRequest{
				User: &github.User{Login: github.String(testCase.user)},
			}
			details := treeprint.New()
			requirement := AuthorInTeam(gh, "team")

			assert.Equal(t, requirement.IsSatisfied(pr, details), testCase.isSatisfied, fmt.Sprintf("requirement should have a satisfied status: %t", testCase.isSatisfied))
			assert.True(t, utils.TestLastNodeStatus(t, testCase.isSatisfied, details), fmt.Sprintf("requirement details should have a status: %t", testCase.isSatisfied))
		})
	}
}