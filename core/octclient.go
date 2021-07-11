package core

import (
	"context"
	"github.com/shurcooL/githubv4"
	"sort"
	"time"
)

func NewOctclient(user string, token string) Octclient {
	return Octclient{
		user: user,
		client: githubv4.NewClient(
			newOauthClient(token),
		),
		oldestUpdatedAt: time.Now(),
		latestUpdatedAt: time.Time{},
	}
}

type Octclient struct {
	user            string
	client          *githubv4.Client
	oldestUpdatedAt time.Time
	latestUpdatedAt time.Time
}

func (o *Octclient) CallQueryRepositoriesContributedTo(ctx context.Context) (*UserRepositoriesContributedTo, error) {
	variables := map[string]interface{}{
		"login": githubv4.String(o.user),
	}
	query := UserRepositoriesContributedTo{}
	err := o.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	return &query, nil
}

func (o *Octclient) updateUpdatedTime(t time.Time) {
	if o.oldestUpdatedAt.After(t) {
		o.oldestUpdatedAt = t
	}
	if o.latestUpdatedAt.Before(t) {
		o.latestUpdatedAt = t
	}
}

func (o *Octclient) GetRepositoriesContributedTo(ctx context.Context, isSortBySize bool, reverse bool) (*Results, error) {
	result, err := o.CallQueryRepositoriesContributedTo(ctx)
	if err != nil {
		return nil, err
	}

	languageMap := map[string]int{}
	for _, node := range result.User.RepositoriesContributedTo.Nodes {
		o.updateUpdatedTime(node.UpdatedAt.Time)

		edges := node.Languages.Edges
		for _, language := range edges {
			languageMap[string(language.Node.Name)] += int(language.Size)
		}
	}

	var languages []LanguageSize
	for k, v := range languageMap {
		languages = append(languages, LanguageSize{
			Name: k,
			Size: v,
		})
	}

	if isSortBySize {
		sort.Slice(languages, func(i, j int) bool { return languages[i].Size > languages[j].Size })
	}

	if reverse {
		reverseSlice(languages)
	}

	return &Results{
		UpdatedRange: UpdatedRange{
			Oldest: o.oldestUpdatedAt,
			Latest: o.latestUpdatedAt,
		},
		LanguageSizes: languages,
	}, nil
}
