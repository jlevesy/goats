package testing

type Status int

const (
	StatusUnknown Status = iota
	StatusSuccess
	StatusFailure
	StatusFatal
)

func (s Status) String() string {
	switch s {
	case StatusSuccess:
		return "success"
	case StatusFailure:
		return "failure"
	case StatusFatal:
		return "fatal"
	default:
		return "unknown"
	}
}
