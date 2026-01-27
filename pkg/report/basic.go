package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/c3-trentbuckholz/baseupgrader/pkg/git"
)

type BasicReport struct {
	files    []string
	baseDiff git.Diff
	writer   io.Writer
}

func NewBasicReport(files []string, baseDiff git.Diff, writer io.Writer) *BasicReport {
	return &BasicReport{
		files:    files,
		baseDiff: baseDiff,
		writer:   writer,
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
	_, err := fmt.Fprint(r.writer, contents)
	return err
}
