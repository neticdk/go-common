package repometrics

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v69/github"
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

	m.LatestVersion = gh.guessVersion()

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
	m.IsApache = gr.owner == "apache"
	m.License = r.GetLicense().GetName()
	if r.CreatedAt != nil {
		m.CreatedAt = &r.CreatedAt.Time
	}
	m.Stats.OpenedIssuesNow = r.GetOpenIssuesCount()
	m.Stats.Likes = r.GetStargazersCount()
	m.Stats.Forks = r.GetForksCount()

	return r, nil
}

func (gr *ghRepo) fetchIssuesAndPRs(m *Metrics) error {
	if m.Stats == nil {
		m.Stats = NewStats()
	}

	counts1Y := &issueCounts{}
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	err := gr.updateIssues(oneYearAgo, counts1Y)
	if err == nil {
		m.Stats.OpenedIssues1Y = counts1Y.openedIssues
		m.Stats.ClosedIssues1Y = counts1Y.closedIssues
		m.Stats.OpenedPRs1Y = counts1Y.openedPulls
		m.Stats.ClosedPRs1Y = counts1Y.closedPulls
		m.Stats.OpenedFeatures1Y = counts1Y.openedFeatures
		m.Stats.ClosedFeatures1Y = counts1Y.closedFeatures
		m.Stats.OpenedBugs1Y = counts1Y.openedBugs
		m.Stats.ClosedBugs1Y = counts1Y.closedBugs

		counts9M := &issueCounts{}
		nineMonthsAgo := time.Now().AddDate(0, -9, 0)
		err := gr.updateIssues(nineMonthsAgo, counts9M)
		if err == nil {
			m.Stats.OpenedIssues9M = counts9M.openedIssues
			m.Stats.ClosedIssues9M = counts9M.closedIssues
			m.Stats.OpenedPRs9M = counts9M.openedPulls
			m.Stats.ClosedPRs9M = counts9M.closedPulls
			m.Stats.OpenedFeatures9M = counts9M.openedFeatures
			m.Stats.ClosedFeatures9M = counts9M.closedFeatures
			m.Stats.OpenedBugs9M = counts9M.openedBugs
			m.Stats.ClosedBugs9M = counts9M.closedBugs

			counts6M := &issueCounts{}
			sixMonthsAgo := time.Now().AddDate(0, -6, 0)
			err := gr.updateIssues(sixMonthsAgo, counts6M)
			if err == nil {
				m.Stats.OpenedIssues6M = counts6M.openedIssues
				m.Stats.ClosedIssues6M = counts6M.closedIssues
				m.Stats.OpenedPRs6M = counts6M.openedPulls
				m.Stats.ClosedPRs6M = counts6M.closedPulls
				m.Stats.OpenedFeatures6M = counts6M.openedFeatures
				m.Stats.ClosedFeatures6M = counts6M.closedFeatures
				m.Stats.OpenedBugs6M = counts6M.openedBugs
				m.Stats.ClosedBugs6M = counts6M.closedBugs

				counts3M := &issueCounts{}
				threeMonthsAgo := time.Now().AddDate(0, -3, 0)
				err := gr.updateIssues(threeMonthsAgo, counts3M)
				if err == nil {
					m.Stats.OpenedIssues3M = counts3M.openedIssues
					m.Stats.ClosedIssues3M = counts3M.closedIssues
					m.Stats.OpenedPRs3M = counts3M.openedPulls
					m.Stats.ClosedPRs3M = counts3M.closedPulls
					m.Stats.OpenedFeatures3M = counts3M.openedFeatures
					m.Stats.ClosedFeatures3M = counts3M.closedFeatures
					m.Stats.OpenedBugs3M = counts3M.openedBugs
					m.Stats.ClosedBugs3M = counts3M.closedBugs

					counts1M := &issueCounts{}
					oneMonthsAgo := time.Now().AddDate(0, -1, 0)
					err := gr.updateIssues(oneMonthsAgo, counts1M)
					if err == nil {
						m.Stats.OpenedIssues1M = counts1M.openedIssues
						m.Stats.ClosedIssues1M = counts1M.closedIssues
						m.Stats.OpenedPRs1M = counts1M.openedPulls
						m.Stats.ClosedPRs1M = counts1M.closedPulls
						m.Stats.OpenedFeatures1M = counts1M.openedFeatures
						m.Stats.ClosedFeatures1M = counts1M.closedFeatures
						m.Stats.OpenedBugs1M = counts1M.openedBugs
						m.Stats.ClosedBugs1M = counts1M.closedBugs

						countsNow := &issueCounts{}
						rightNow := time.Now().AddDate(0, 0, -1)
						err := gr.updateIssues(rightNow, countsNow)
						if err == nil {
							m.Stats.OpenedIssuesNow = countsNow.openedIssues
							m.Stats.ClosedIssuesNow = countsNow.closedIssues
							m.Stats.OpenedPRsNow = countsNow.openedPulls
							m.Stats.ClosedPRsNow = countsNow.closedPulls
							m.Stats.OpenedFeaturesNow = countsNow.openedFeatures
							m.Stats.ClosedFeaturesNow = countsNow.closedFeatures
							m.Stats.OpenedBugsNow = countsNow.openedBugs
							m.Stats.ClosedBugsNow = countsNow.closedBugs
						}
					}
				}
			}
		}
	}
	return err
}

func (gr *ghRepo) updateIssues(pit time.Time,
	observations *issueCounts) error {
	issuesOpts := &github.IssueListByRepoOptions{
		State:       "all",
		Since:       pit,
		ListOptions: github.ListOptions{PerPage: perPage},
	}
	err := gr.paginate(func(page int) (int, error) {
		issuesOpts.Page = page
		pageIssues, resp, err := gr.client.Issues.ListByRepo(gr.ctx, gr.owner, gr.name, issuesOpts)
		if err != nil {
			return 0, fmt.Errorf("listing issues: %w", err)
		}
		for _, issue := range pageIssues {
			gr.updateIssuePRCounts(issue, observations)
			gr.updateFeatureCounts(issue, observations)
			gr.updateBugsCounts(issue, observations)
		}
		return resp.NextPage, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (gr *ghRepo) updateIssuePRCounts(issue *github.Issue,
	counts *issueCounts) {
	if issue.IsPullRequest() {
		if issue.GetState() == "open" {
			counts.openedPulls++
		} else if issue.GetState() == "closed" {
			counts.closedPulls++
		}
	} else {
		if issue.GetState() == "open" {
			counts.openedIssues++
		} else if issue.GetState() == "closed" {
			counts.closedIssues++
		}
	}
}

func (gr *ghRepo) updateFeatureCounts(issue *github.Issue,
	counts *issueCounts) {
	if issue.Labels != nil {
		for i := 0; i < len(issue.Labels); i++ {
			if strings.Contains(*issue.Labels[i].Name, "enhancement") ||
				strings.Contains(*issue.Labels[i].Name, "feat") {
				if issue.GetState() == "open" {
					counts.openedFeatures++
				} else if issue.GetState() == "closed" {
					counts.closedFeatures++
				}
			}
		}
	}
}

func (gr *ghRepo) updateBugsCounts(issue *github.Issue,
	counts *issueCounts) {
	if issue.Labels != nil {
		for i := 0; i < len(issue.Labels); i++ {
			if strings.Contains(*issue.Labels[i].Name, "bug") {
				if issue.GetState() == "open" {
					counts.openedBugs++
				} else if issue.GetState() == "closed" {
					counts.closedBugs++
				}
			}
		}
	}
}

func (gr *ghRepo) fetchContributors(m *Metrics) error {
	contributorStats, _, err := gr.client.Repositories.ListContributorsStats(gr.ctx, gr.owner, gr.name)
	if err != nil {
		return fmt.Errorf("listing contributor stats: %w", err)
	}

	topCommitters := make([]Contributor, 0)
	topCommitters1Y := make([]Contributor, 0)
	topCommitters9M := make([]Contributor, 0)
	topCommitters6M := make([]Contributor, 0)
	topCommitters3M := make([]Contributor, 0)
	topCommitters1M := make([]Contributor, 0)

	for _, stat := range contributorStats {
		if isBot(stat.GetAuthor().GetLogin()) {
			continue
		}

		contributor := Contributor{
			Name:    stat.GetAuthor().GetLogin(),
			Commits: stat.GetTotal(),
		}
		ecoType := m.Type

		if ecoType == "github" {
			contributor.URL = fmt.Sprintf("https://github.com/%s", stat.GetAuthor().GetLogin())
		}
		topCommitters = append(topCommitters, contributor)

		oneYearAgo := time.Now().AddDate(-1, 0, 0)
		topCommitters1Y = handleCommittersInPeriod(oneYearAgo, topCommitters1Y, stat, ecoType)

		nineMonthsAgo := time.Now().AddDate(0, -9, 0)
		topCommitters9M = handleCommittersInPeriod(nineMonthsAgo, topCommitters9M, stat, ecoType)

		sixMonthsAgo := time.Now().AddDate(0, -6, 0)
		topCommitters6M = handleCommittersInPeriod(sixMonthsAgo, topCommitters6M, stat, ecoType)

		threeMonthsAgo := time.Now().AddDate(0, -3, 0)
		topCommitters3M = handleCommittersInPeriod(threeMonthsAgo, topCommitters3M, stat, ecoType)

		oneMonthAgo := time.Now().AddDate(0, -1, 0)
		topCommitters1M = handleCommittersInPeriod(oneMonthAgo, topCommitters1M, stat, ecoType)
	}

	sortContributors(topCommitters)
	sortContributors(topCommitters1Y)
	sortContributors(topCommitters9M)
	sortContributors(topCommitters6M)
	sortContributors(topCommitters3M)
	sortContributors(topCommitters1M)

	if m.Stats == nil {
		m.Stats = NewStats()
	}
	m.Stats.Contributors1Y = len(topCommitters1Y)
	m.Stats.TopCommitters = getTopContributors(topCommitters, 10)
	m.Stats.TopCommitters1Y = getTopContributors(topCommitters1Y, 10)
	m.Stats.TopCommitters9M = getTopContributors(topCommitters9M, 10)
	m.Stats.TopCommitters6M = getTopContributors(topCommitters6M, 10)
	m.Stats.TopCommitters3M = getTopContributors(topCommitters3M, 10)
	m.Stats.TopCommitters1M = getTopContributors(topCommitters1M, 10)
	return nil
}

func handleCommittersInPeriod(period time.Time, committersInPeriod []Contributor, stat *github.ContributorStats, ecoType RepoType) []Contributor {
	contributorPeriod := Contributor{
		Name: stat.GetAuthor().GetLogin(),
	}
	if ecoType == "github" {
		contributorPeriod.URL = fmt.Sprintf("https://github.com/%s", stat.GetAuthor().GetLogin())
	}
	for _, week := range stat.Weeks {
		if week.Week.Unix() > period.Unix() {
			contributorPeriod.Commits += week.GetCommits()
		}
	}
	committersInPeriod = append(committersInPeriod, contributorPeriod)
	return committersInPeriod
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
	verifiedCommitsPerMonth6M := gr.calculateVerifiedCommitsPerMonth(commits, 6)

	if m.Stats == nil {
		m.Stats = NewStats()
	}
	m.Stats.LastCommit = lastCommit
	m.Stats.CommitsPerMonth6M = commitsPerMonth6M
	m.Stats.VerifiedCommitsPerMonth6M = verifiedCommitsPerMonth6M
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

func (gr *ghRepo) calculateVerifiedCommitsPerMonth(commits []*github.RepositoryCommit, months int) int {
	if months <= 0 {
		return 0
	}
	verifiedCommits := 0
	for _, commit := range commits {
		if commit.GetCommit().GetVerification().GetVerified() {
			verifiedCommits++
		}
	}
	return verifiedCommits / months
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
	m.Stats.NoOfReleases = len(releases)
	m.Stats.Releases = statistifyReleases(releases)
	m.Stats.ReleasesPerDay = releaseMetrics.PerDay
	m.Stats.ReleasesPerWeek = releaseMetrics.PerWeek
	m.Stats.ReleasesPerMonth = releaseMetrics.PerMonth
	m.Stats.ReleasesPerYear = releaseMetrics.PerYear

	return nil
}

func statistifyReleases(releases []*github.RepositoryRelease) []Release {
	releaseStats := make([]Release, 0)
	for _, release := range releases {
		releaseStats = append(releaseStats, Release{
			Name:       release.GetTagName(),
			Date:       release.GetPublishedAt().Time,
			ReleaseURL: release.GetHTMLURL(),
			AssetsURL:  release.GetAssetsURL(),
			UploadURL:  release.GetUploadURL(),
			TarballURL: release.GetTarballURL(),
		})
	}
	return releaseStats
}

func (gr *ghRepo) getFirstAndLastRelease(releases []*github.RepositoryRelease) (*github.RepositoryRelease, *github.RepositoryRelease) {
	if len(releases) == 0 {
		return nil, nil
	}
	firstRelease := releases[len(releases)-1]
	lastRelease := releases[0]
	return firstRelease, lastRelease
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
