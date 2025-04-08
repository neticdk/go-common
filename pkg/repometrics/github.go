package repometrics

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v71/github"
	"golang.org/x/mod/semver"
	"golang.org/x/oauth2"
)

// UpdateGitHub updates the metrics using data from a GitHub repository
func (m *Metrics) UpdateGitHub(ctx context.Context, client *github.Client, owner, repository string) error {
	if client == nil {
		token := os.Getenv("GITHUB_TOKEN")
		if token != "" {
			client = github.NewClient(oauth2.NewClient(ctx,
				oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})))
		} else {
			client = github.NewClient(nil)
		}
	}

	gh := &ghRepo{ctx: ctx, client: client, owner: owner, name: repository}

	repo, err := gh.fetchRepoInfo(m)
	if err != nil {
		return fmt.Errorf("fetching repository info: %w", err)
	}
	gh.repo = repo

	err = gh.fetchIssuesAndPRs(m)
	if err != nil {
		return fmt.Errorf("fetching issues and PRs: %w", err)
	}

	err = gh.fetchContributors(m)
	if err != nil {
		return fmt.Errorf("fetching contributors: %w", err)
	}

	err = gh.fetchCommits(m)
	if err != nil {
		return fmt.Errorf("fetching commits: %w", err)
	}

	err = gh.fetchReleases(m)
	if err != nil {
		return fmt.Errorf("fetching releases: %w", err)
	}

	m.CurrentVersion = gh.guessVersion()

	return nil
}

type ghRepo struct {
	ctx    context.Context
	client *github.Client
	repo   *github.Repository
	owner  string
	name   string
}

func (gr *ghRepo) fetchRepoInfo(m *Metrics) (*github.Repository, error) {
	r, _, err := gr.client.Repositories.Get(gr.ctx, gr.owner, gr.name)
	if err != nil {
		return nil, fmt.Errorf("getting repository info: %w", err)
	}
	gr.repo = r

	m.Name = r.GetName()
	m.URL = r.GetHTMLURL()
	m.IsCNCF = false
	m.IsKubernetesSIG = gr.owner == "kubernetes-sigs"
	m.License = r.GetLicense().GetName()
	if r.CreatedAt != nil {
		m.CreatedAt = &r.CreatedAt.Time
	}
	m.Stats.OpenIssuesNow = r.GetOpenIssuesCount()
	m.Stats.Likes = r.GetStargazersCount()
	m.Stats.Forks = r.GetForksCount()

	return r, nil
}

func (gr *ghRepo) fetchIssuesAndPRs(m *Metrics) error {
	counts := &issuePRCounts{}

	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	issuesOpts := &github.IssueListByRepoOptions{
		State:       "all",
		Since:       sixMonthsAgo,
		ListOptions: github.ListOptions{PerPage: perPage},
	}

	err := gr.paginate(func(page int) (int, error) {
		issuesOpts.Page = page
		pageIssues, resp, err := gr.client.Issues.ListByRepo(gr.ctx, gr.owner, gr.name, issuesOpts)
		if err != nil {
			return 0, fmt.Errorf("listing issues: %w", err)
		}
		for _, issue := range pageIssues {
			gr.updateIssuePRCounts(issue, counts)
		}
		return resp.NextPage, nil
	})
	if err != nil {
		return err
	}

	if m.Stats == nil {
		m.Stats = NewStats()
	}

	m.Stats.OpenedIssues6M = counts.openedIssues6M
	m.Stats.ClosedIssues6M = counts.closedIssues6M
	m.Stats.OpenedPRs6M = counts.openedPulls6M
	m.Stats.ClosedPRs6M = counts.closedPulls6M

	return nil
}

func (gr *ghRepo) updateIssuePRCounts(issue *github.Issue, counts *issuePRCounts) {
	if issue.IsPullRequest() {
		if issue.GetState() == "open" {
			counts.openedPulls6M++
		} else if issue.GetState() == "closed" {
			counts.closedPulls6M++
		}
	} else {
		if issue.GetState() == "open" {
			counts.openedIssues6M++
		} else if issue.GetState() == "closed" {
			counts.closedIssues6M++
		}
	}
}

func (gr *ghRepo) fetchContributors(m *Metrics) error {
	contributorStats, _, err := gr.client.Repositories.ListContributorsStats(gr.ctx, gr.owner, gr.name)
	if err != nil {
		return fmt.Errorf("listing contributor stats: %w", err)
	}

	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	topCommitters := make([]Contributor, 0)
	topCommitters1Y := make([]Contributor, 0)

	for _, stat := range contributorStats {
		if isBot(stat.GetAuthor().GetLogin()) {
			continue
		}

		contributor := Contributor{
			Name:    stat.GetAuthor().GetLogin(),
			Commits: stat.GetTotal(),
		}
		if m.Type == "github" {
			contributor.URL = fmt.Sprintf("https://github.com/%s", stat.GetAuthor().GetLogin())
		}
		topCommitters = append(topCommitters, contributor)

		contributor1Y := Contributor{
			Name: stat.GetAuthor().GetLogin(),
		}
		if m.Type == "github" {
			contributor1Y.URL = fmt.Sprintf("https://github.com/%s", stat.GetAuthor().GetLogin())
		}
		for _, week := range stat.Weeks {
			if week.Week.Unix() > oneYearAgo.Unix() {
				contributor1Y.Commits += week.GetCommits()
			}
		}
		topCommitters1Y = append(topCommitters1Y, contributor1Y)
	}

	sortContributors(topCommitters)
	sortContributors(topCommitters1Y)

	if m.Stats == nil {
		m.Stats = NewStats()
	}
	m.Stats.Contributors1Y = len(topCommitters1Y)
	m.Stats.TopCommitters = getTopContributors(topCommitters, 10)
	m.Stats.TopCommitters1Y = getTopContributors(topCommitters1Y, 10)

	return nil
}

func (gr *ghRepo) fetchCommits(m *Metrics) error {
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)
	commitsOpts := &github.CommitsListOptions{Since: sixMonthsAgo, ListOptions: github.ListOptions{PerPage: perPage}}

	var commits []*github.RepositoryCommit
	err := gr.paginate(func(page int) (int, error) {
		commitsOpts.Page = page
		pageCommits, resp, err := gr.client.Repositories.ListCommits(gr.ctx, gr.owner, gr.name, commitsOpts)
		if err != nil {
			return 0, fmt.Errorf("listing commits: %w", err)
		}
		commits = append(commits, pageCommits...)
		return resp.NextPage, nil
	})
	if err != nil {
		return err
	}

	lastCommit := gr.getLastCommit(commits)
	commitsPerMonth6M := gr.calculateCommitsPerMonth(commits, 6)

	if m.Stats == nil {
		m.Stats = NewStats()
	}
	m.Stats.LastCommit = lastCommit
	m.Stats.CommitsPerMonth6M = commitsPerMonth6M

	return nil
}

func (gr *ghRepo) getLastCommit(commits []*github.RepositoryCommit) *time.Time {
	if len(commits) > 0 {
		commit := commits[0].GetCommit().GetAuthor().GetDate().Time
		return &commit
	}
	return nil
}

func (gr *ghRepo) calculateCommitsPerMonth(commits []*github.RepositoryCommit, months int) int {
	if months <= 0 {
		return 0
	}
	return len(commits) / months
}

func (gr *ghRepo) fetchReleases(m *Metrics) error {
	releasesOpts := &github.ListOptions{PerPage: perPage}

	var releases []*github.RepositoryRelease
	err := gr.paginate(func(page int) (int, error) {
		releasesOpts.Page = page
		pageReleases, resp, err := gr.client.Repositories.ListReleases(gr.ctx, gr.owner, gr.name, releasesOpts)
		if err != nil {
			return 0, fmt.Errorf("listing releases: %w", err)
		}
		releases = append(releases, pageReleases...)
		return resp.NextPage, nil
	})
	if err != nil {
		return err
	}

	var firstReleaseTime, lastReleaseTime *time.Time
	firstRelease, lastRelease := gr.getFirstAndLastRelease(releases)

	if firstRelease != nil {
		t := firstRelease.GetPublishedAt().Time
		firstReleaseTime = &t
	}
	if lastRelease != nil {
		t := lastRelease.GetPublishedAt().Time
		lastReleaseTime = &t
	}
	releaseMetrics := CalculateReleaseMetrics(firstReleaseTime, lastReleaseTime, len(releases))

	if m.Stats == nil {
		m.Stats = NewStats()
	}

	m.Stats.FirstRelease = firstReleaseTime
	m.Stats.LastRelease = lastReleaseTime
	m.Stats.Releases = len(releases)
	m.Stats.ReleasesPerDay = releaseMetrics.PerDay
	m.Stats.ReleasesPerWeek = releaseMetrics.PerWeek
	m.Stats.ReleasesPerMonth = releaseMetrics.PerMonth
	m.Stats.ReleasesPerYear = releaseMetrics.PerYear

	return nil
}

func (gr *ghRepo) getFirstAndLastRelease(releases []*github.RepositoryRelease) (first *github.RepositoryRelease, last *github.RepositoryRelease) {
	if len(releases) == 0 {
		return nil, nil
	}
	first = releases[len(releases)-1]
	last = releases[0]
	return first, last
}

func (gr *ghRepo) guessVersion() string {
	if gr.repo == nil || gr.client == nil {
		return ""
	}

	// 1. Check for releases.
	releases, _, err := gr.client.Repositories.ListReleases(gr.ctx, gr.repo.GetOwner().GetLogin(), gr.repo.GetName(), nil)
	if err == nil && len(releases) > 0 {
		sort.Slice(releases, func(i, j int) bool {
			return semver.Compare(releases[i].GetTagName(), releases[j].GetTagName()) < 0
		})
		if releases[len(releases)-1].TagName != nil {
			return releases[len(releases)-1].GetTagName()
		}
	}

	// 2. Check for tags.
	tags, _, err := gr.client.Repositories.ListTags(gr.ctx, gr.repo.GetOwner().GetLogin(), gr.repo.GetName(), nil)
	if err == nil && len(tags) > 0 {
		sort.Slice(tags, func(i, j int) bool {
			return semver.Compare(tags[i].GetName(), tags[j].GetName()) < 0
		})
		return tags[len(tags)-1].GetName()
	}

	// 3. Check the default branch name.
	defaultBranch := gr.repo.GetDefaultBranch()
	if strings.HasPrefix(defaultBranch, "v") {
		return strings.TrimPrefix(defaultBranch, "v")
	}

	return "unknown"
}

func (gr *ghRepo) paginate(fetchPage func(page int) (int, error)) error {
	page := 1
	for {
		nextPage, err := fetchPage(page)
		if err != nil {
			return err
		}
		if nextPage == 0 {
			break
		}
		page = nextPage
	}
	return nil
}
