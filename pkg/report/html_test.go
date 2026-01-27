package report

import (
	"testing"

	"github.com/c3-trentbuckholz/baseupgrader/pkg/git"
	"github.com/stretchr/testify/assert"
)

func Test_CreateHtmlReport(t *testing.T) {
	output, err := NewHtmlReport([]string{"file1.txt", "file2.txt"}, git.Diff{
		"file1.txt": "diff --git a/file1.txt b/file1.txt\n+Added line in file1",
		"file2.txt": "diff --git a/file2.txt b/file2.txt\n-Removed line in file2",
	}).Create()
	assert.NoError(t, err)

	expectedOutput := `<style>.file-header { font-weight: 600; font-size: 18px;` +
		` padding: 12px 16px; border-bottom: 1px solid #d0d7de; color: #1f2328; }` +
		`.diff-add { background-color: #e6ffed; color: #22863a; }.diff-del ` +
		`{ background-color: #ffeef0; color: #cb2431; }.diff-hunk ` +
		`{ color: #005cc5; background-color: #f1f8ff; }pre { font-family: monospace;` +
		` white-space: pre-wrap; }</style><pre><div class='file-header'>file1.txt</div>` +
		`<div>diff --git a/file1.txt b/file1.txt</div><div class='diff-add'>+Added line in ` +
		`file1</div><div class='file-header'>file2.txt</div><div>diff --git a/file2.txt ` +
		`b/file2.txt</div><div class='diff-del'>-Removed line in file2</div></pre>`
	assert.Equal(t, expectedOutput, output)
}
