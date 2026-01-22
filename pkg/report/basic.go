package report

import (
	"fmt"
	"strings"

	"github.com/c3-trentbuckholz/baseupgrader/pkg/git"
)

type BasicReport struct {
	files    []string
	baseDiff git.Diff
}

func NewBasicReport(files []string, baseDiff git.Diff) *BasicReport {
	return &BasicReport{
		files:    files,
		baseDiff: baseDiff,
	}
}

func (r *BasicReport) Create() (string, error) {
	output := strings.Builder{}
	for _, fileName := range r.files {
		if _, ok := r.baseDiff[fileName]; ok {
			output.WriteString(fileName + "\n")
			output.WriteString("--------------------------------------------\n")
			output.WriteString(r.baseDiff[fileName] + "\n")
			output.WriteString("============================================\n")
		}
	}
	return output.String(), nil
}

func (r *BasicReport) Write(contents string) error {
	fmt.Println(contents)
	return nil
}
