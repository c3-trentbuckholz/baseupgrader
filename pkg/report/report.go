package report

type Reporter interface {
	Create() (string, error)
	Write(contents string) error
}
