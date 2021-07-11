package core

import (
	"github.com/shurcooL/githubv4"
	"time"
)

type UpdatedRange struct {
	Oldest time.Time `json:"oldest"`
	Latest time.Time `json:"latest"`
}

type LanguageSize struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

type Results struct {
	UpdatedRange  UpdatedRange   `json:"updated_range"`
	LanguageSizes []LanguageSize `json:"language_sizes"`
}

type UserRepositoriesContributedTo struct {
	User struct {
		RepositoriesContributedTo struct {
			TotalCount githubv4.Int
			Nodes      []struct {
				UpdatedAt githubv4.DateTime
				Languages struct {
					Edges []struct {
						Size githubv4.Int
						Node struct {
							Name githubv4.String
						}
					}
				} `graphql:"languages(first: 100)"`
			}
		} `graphql:"repositoriesContributedTo(first: 100, includeUserRepositories: true, contributionTypes: [COMMIT], orderBy: {field: UPDATED_AT, direction: DESC})"`
	} `graphql:"user(login: $login)"`
}
