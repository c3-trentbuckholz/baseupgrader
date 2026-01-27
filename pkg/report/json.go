package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/c3-trentbuckholz/baseupgrader/pkg/git"
)

type JsonReport struct {
	files    []string
	baseDiff git.Diff
	writer   io.Writer
}

func NewJsonReport(files []string, baseDiff git.Diff, writer io.Writer) *JsonReport {
	return &JsonReport{
		files:    files,
		baseDiff: baseDiff,
		writer:   writer,
	}
}

func (r *JsonReport) Create() (string, error) {
	outputDiff := git.Diff{}
	for _, fileName := range r.files {
		if _, ok := r.baseDiff[fileName]; ok {
			outputDiff[fileName] = r.baseDiff[fileName]
		}
	}
	output, err := json.Marshal(outputDiff)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (r *JsonReport) Write(contents string) error {
	_, err := fmt.Fprint(r.writer, contents)
	return err
}
