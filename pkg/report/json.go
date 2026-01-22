package report

import (
	"encoding/json"
	"fmt"

	"github.com/c3-trentbuckholz/baseupgrader/pkg/git"
)

type JsonReport struct {
	files    []string
	baseDiff git.Diff
}

func NewJsonReport(files []string, baseDiff git.Diff) *JsonReport {
	return &JsonReport{
		files:    files,
		baseDiff: baseDiff,
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
	fmt.Println(contents)
	return nil
}
