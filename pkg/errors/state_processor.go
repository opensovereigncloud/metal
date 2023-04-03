package errors

type StateProcErrorReason string

func (s StateProcErrorReason) String() string {
	return string(s)
}

type StateProcError interface {
	Error() string
	Reason() StateProcErrorReason
	Message() string
}
