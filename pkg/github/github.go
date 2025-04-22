package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v71/github"
)

// CreateGitHubRepository creates a new GitHub repository
func CreateGitHubRepository(ctx context.Context, client *github.Client, owner, name string) (*github.Repository, *github.Response, error) {
	repo := &github.Repository{
		Name:    github.Ptr(name),
		Private: github.Ptr(true),
	}
	repo, res, err := client.Repositories.Create(ctx, owner, repo)
	if err != nil {
		return nil, nil, err
	}
	return repo, res, nil
}

func GetReleaseByTag(ctx context.Context, client *github.Client, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error) {
	release, res, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	if err != nil {
		return nil, nil, fmt.Errorf(`getting release tagged %q: %w`, tag, err)
	}
	return release, res, nil
}

func GetLatestRelease(ctx context.Context, client *github.Client, owner, repo string) (*github.RepositoryRelease, *github.Response, error) {
	release, res, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return nil, nil, fmt.Errorf("getting latest release: %w", err)
	}
	return release, res, nil
}

func GetReleaseByTagOrLatest(ctx context.Context, client *github.Client, owner, repo, tag string) (*github.RepositoryRelease, *github.Response, error) {
	if tag == "" {
		return GetLatestRelease(ctx, client, owner, repo)
	}
	return GetReleaseByTag(ctx, client, owner, repo, tag)
}

func GetOrganizationTeams(ctx context.Context, client *github.Client, org string) ([]*github.Team, *github.Response, error) {
	teams, res, err := client.Teams.ListTeams(ctx, org, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("getting teams: %w", err)
	}
	return teams, res, nil
}
