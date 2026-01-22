package report

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/c3-trentbuckholz/baseupgrader/pkg/git"
)

type HtmlReport struct {
	files    []string
	baseDiff git.Diff
}

func NewHtmlReport(files []string, baseDiff git.Diff) *HtmlReport {
	return &HtmlReport{
		files:    files,
		baseDiff: baseDiff,
	}
}

func (r *HtmlReport) Create() (string, error) {
	output := strings.Builder{}

	output.WriteString("<style>")
	output.WriteString(".file-header { font-weight: 600; font-size: 18px; padding: 12px 16px; border-bottom: 1px solid #d0d7de; color: #1f2328; }")
	output.WriteString(".diff-add { background-color: #e6ffed; color: #22863a; }")
	output.WriteString(".diff-del { background-color: #ffeef0; color: #cb2431; }")
	output.WriteString(".diff-hunk { color: #005cc5; background-color: #f1f8ff; }")
	output.WriteString("pre { font-family: monospace; white-space: pre-wrap; }")
	output.WriteString("</style>")

	output.WriteString("<pre>")
	for _, fileName := range r.files {
		if diff, ok := r.baseDiff[fileName]; ok {
			fmt.Fprintf(&output, "<div class='file-header'>%s</div>", html.EscapeString(fileName))

			for line := range strings.SplitSeq(diff, "\n") {
				escaped := html.EscapeString(line)

				switch {
				case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
					fmt.Fprintf(&output, "<div class='diff-add'>%s</div>", escaped)
				case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
					fmt.Fprintf(&output, "<div class='diff-del'>%s</div>", escaped)
				case strings.HasPrefix(line, "@@"):
					fmt.Fprintf(&output, "<div class='diff-hunk'>%s</div>", escaped)
				default:
					fmt.Fprintf(&output, "<div>%s</div>", escaped)
				}
			}
		}
	}
	output.WriteString("</pre>")
	return output.String(), nil
}

func (r *HtmlReport) Write(contents string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		_, err := fmt.Fprint(w, contents)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	fmt.Println("Server starting at http://localhost:8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		return err
	}
	return nil
}
