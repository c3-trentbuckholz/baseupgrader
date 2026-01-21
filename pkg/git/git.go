package git

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

// Diff represents a mapping of filenames to their corresponding patch data.
type Diff = map[string]string

const GITHUB_API_URL = "https://api.github.com"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type VcsClient interface {
	GetDiff(oldCommit string, newCommit string) (Diff, error)
	GetFiles(pattern string) ([]string, error)
}

type GitClient struct {
	httpClient HttpClient
	repoUrl    string
	owner      string
	repo       string
	authToken  string
}

type treeResponse struct {
	Tree []getFileResponse `json:"tree"`
}
type getFileResponse struct {
	Path string `json:"path"`
}

func NewGitClient(repoUrl string, authToken string, httpClient HttpClient) *GitClient {
	parts := strings.Split(repoUrl, "/")
	owner := parts[len(parts)-2]
	repo := strings.TrimSuffix(parts[len(parts)-1], ".git")
	return &GitClient{
		httpClient: httpClient,
		repoUrl:    repoUrl,
		owner:      owner,
		repo:       repo,
		authToken:  authToken,
	}
}

// GetFiles retrieves a list of files names in the Git repository, for the given directory path.
func (g *GitClient) GetFiles(dirPath string) ([]string, error) {
	apiUrl := fmt.Sprintf("%s/repos/%s/%s/git/trees/develop?recursive=1",
		GITHUB_API_URL, g.owner, g.repo)

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3.diff")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.authToken))

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var treeResp treeResponse
	if err := json.Unmarshal(output, &treeResp); err != nil {
		return nil, err
	}

	matchedFiles := make([]string, 0)
	for _, fileResp := range treeResp.Tree {
		if strings.HasPrefix(fileResp.Path, dirPath) {
			matchedFiles = append(matchedFiles, filepath.Base(fileResp.Path))
		}
	}
	return matchedFiles, nil
}

// GetDiff retrieves the diff between two commits in the Git repository.
func (g *GitClient) GetDiff(oldCommit string, newCommit string) (Diff, error) {
	apiUrl := fmt.Sprintf("%s/repos/%s/%s/compare/%s...%s",
		GITHUB_API_URL, g.owner, g.repo, oldCommit, newCommit)

	req, err := http.NewRequest("GET", apiUrl, nil)
	req.Header.Set("Accept", "application/vnd.github.v3.diff")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.authToken))
	if err != nil {
		return nil, err
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseDiff(string(output)), nil
}

func parseDiff(diffText string) Diff {
	diff := make(Diff)

	var (
		currentFile string
		patchLines  strings.Builder
	)
	for line := range strings.SplitSeq(diffText, "\n") {
		// Search for the start of a new file diff.
		if strings.HasPrefix(line, "diff --git") {
			// If not the first file, save the previous file's patch.
			if currentFile != "" {
				diff[currentFile] = patchLines.String()
				patchLines.Reset()
			}
			// Extract the filename from the diff header.
			parts := strings.Split(line, " ")
			if len(parts) >= 4 {
				parts = strings.Split(parts[len(parts)-1], "/")
				currentFile = parts[len(parts)-1]
			}
			continue
		}
		// Accumulate patch lines for the current file.
		if currentFile != "" {
			patchLines.WriteString(line + "\n")
		}
	}
	// Save the last file's patch if exists.
	if currentFile != "" {
		diff[currentFile] = strings.TrimSuffix(patchLines.String(), "\n")
	}
	return diff
}
