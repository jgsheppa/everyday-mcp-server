package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type PullRequest struct {
	Number       int    `json:"number"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	HTMLURL      string `json:"html_url"`
	ChangedFiles int    `json:"changed_files"`
	Additions    int    `json:"additions"`
	Deletions    int    `json:"deletions"`
	User         struct {
		Login string `json:"login"`
	} `json:"user"`
	Draft bool `json:"draft"`
}

type GithubPullRequest struct {
	githubToken string
}

func NewGithubPullRequest(githubToken string) *GithubPullRequest {
	return &GithubPullRequest{githubToken}
}

func (gh *GithubPullRequest) ParseGitHubURL(url string) (owner, repo string, prNumber int, err error) {
	re := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)/pull/(\d+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) != 4 {
		return "", "", 0, fmt.Errorf("invalid GitHub PR URL format")
	}

	prNumber, err = strconv.Atoi(matches[3])
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid PR number: %v", err)
	}

	return matches[1], matches[2], prNumber, nil
}

func (s *GithubPullRequest) FetchPRDetails(owner, repo string, prNumber int) (*PullRequest, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, prNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.githubToken))
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}

	return &pr, nil
}
