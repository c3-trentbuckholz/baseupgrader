package report

import (
	"strings"
	"testing"

	"github.com/c3-trentbuckholz/baseupgrader/pkg/git"
	"github.com/stretchr/testify/assert"
)

func Test_CreateBasicReport(t *testing.T) {
	output, err := NewBasicReport([]string{"file1.txt", "file2.txt"}, git.Diff{
		"file1.txt": "diff --git a/file1.txt b/file1.txt\n+Added line in file1",
		"file2.txt": "diff --git a/file2.txt b/file2.txt\n-Removed line in file2",
	}, nil).Create()
	assert.NoError(t, err)

	expectedOutput := `file1.txt
--------------------------------------------
diff --git a/file1.txt b/file1.txt
+Added line in file1
============================================
file2.txt
--------------------------------------------
diff --git a/file2.txt b/file2.txt
-Removed line in file2
============================================
`
	assert.Equal(t, expectedOutput, output)
}

func Test_WriteBasicReport(t *testing.T) {
	var outputBuilder strings.Builder
	err := NewBasicReport(nil, nil, &outputBuilder).Write("Test content")

	assert.NoError(t, err)
	assert.Equal(t, "Test content", outputBuilder.String())
}
