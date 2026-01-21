package git

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockHttpClient struct {
	mock.Mock
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func Test_GetFilesSuccess(t *testing.T) {
	var httpClient mockHttpClient
	httpClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body: io.NopCloser(strings.NewReader(
			`{"tree":[{"path":"path/to/dir/file1.txt"},{"path":"path/to/dir/file2.txt"},{"path":"path/to/other/file3.txt"}]}`)),
	}, nil)

	gitClient := NewGitClient("https://github.com/owner/repo.git", "testtoken", &httpClient)

	files, err := gitClient.GetFiles("path/to/dir")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	httpClient.AssertExpectations(t)

	expectedFiles := []string{"file1.txt", "file2.txt"}
	assert.Equal(t, len(expectedFiles), len(files), "file count mismatch")
	for i, file := range files {
		assert.Equal(t, expectedFiles[i], file, "file mismatch")
	}
}

func Test_GetFilesHttpError(t *testing.T) {
	var httpClient mockHttpClient
	httpClient.On("Do", mock.Anything).Return(&http.Response{}, fmt.Errorf("network error"))

	gitClient := NewGitClient("https://github.com/owner/repo.git", "testtoken", &httpClient)

	_, err := gitClient.GetFiles("path/to/dir")
	assert.NotNil(t, err, "expected error, got nil")

	httpClient.AssertExpectations(t)
}

func Test_GetFilesMalformedHttpResponse(t *testing.T) {
	var httpClient mockHttpClient
	httpClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body: io.NopCloser(strings.NewReader(
			`{"tree":[{"path":"path/to/dir/file1.txt"},{"path""path/to/dir/file2.txt"},{"path":"path/to/other/file3.txt}]}`)),
	}, nil)

	gitClient := NewGitClient("https://github.com/owner/repo.git", "testtoken", &httpClient)

	_, err := gitClient.GetFiles("path/to/dir")

	httpClient.AssertExpectations(t)

	assert.NotNil(t, err, "expected error, got nil")
}

func Test_GetFilesNoResponseBody(t *testing.T) {
	var httpClient mockHttpClient
	httpClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil)

	gitClient := NewGitClient("https://github.com/owner/repo.git", "testtoken", &httpClient)

	_, err := gitClient.GetFiles("path/to/dir")

	httpClient.AssertExpectations(t)

	assert.NotNil(t, err, "expected error, got nil")
}

func Test_GetDiffSuccess(t *testing.T) {
	var httpClient mockHttpClient
	httpClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body: io.NopCloser(strings.NewReader(
			`diff --git a/path/to/file1.txt b/path/to/file1.txt
index e69de29..4b825dc 100644
--- a/path/to/file1.txt
+++ b/path/to/file1.txt
@@ -0,0 +1 @@
+Hello World
diff --git a/path/to/file2.txt b/path/to/file2.txt
index e69de29..4b825dc 100644
--- a/path/to/file2.txt
+++ b/path/to/file2.txt
@@ -0,0 +1 @@
+Goodbye World
`)),
	}, nil)

	gitClient := NewGitClient("https://github.com/owner/repo.git", "testtoken", &httpClient)

	diff, err := gitClient.GetDiff("oldCommitHash", "newCommitHash")
	assert.Nil(t, err, "expected no error, got %v", err)

	httpClient.AssertExpectations(t)

	expectedDiff := Diff{
		"file1.txt": `index e69de29..4b825dc 100644
--- a/path/to/file1.txt
+++ b/path/to/file1.txt
@@ -0,0 +1 @@
+Hello World
`,
		"file2.txt": `index e69de29..4b825dc 100644
--- a/path/to/file2.txt
+++ b/path/to/file2.txt
@@ -0,0 +1 @@
+Goodbye World
`,
	}
	assert.Equal(t, expectedDiff, diff, "diff mismatch")
}

func Test_GetDiffHttpError(t *testing.T) {
	var httpClient mockHttpClient
	httpClient.On("Do", mock.Anything).Return(&http.Response{}, fmt.Errorf("network error"))

	gitClient := NewGitClient("https://github.com/owner/repo.git", "testtoken", &httpClient)

	_, err := gitClient.GetDiff("oldCommitHash", "newCommitHash")
	assert.NotNil(t, err, "expected error, got nil")

	httpClient.AssertExpectations(t)
}

func Test_GetDiffEmptyResponseBody(t *testing.T) {
	var httpClient mockHttpClient
	httpClient.On("Do", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body: io.NopCloser(strings.NewReader(
			"")),
	}, nil)

	gitClient := NewGitClient("https://github.com/owner/repo.git", "testtoken", &httpClient)

	diff, err := gitClient.GetDiff("oldCommitHash", "newCommitHash")
	assert.Nil(t, err, "expected no error, got %v", err)

	httpClient.AssertExpectations(t)

	assert.Equal(t, Diff{}, diff, "expected empty diff, got %v", diff)
}
