package usecase

import "errors"

const (
	alreadyOnboarded = "already onboarded"
	notFound         = "not found"
	unknown          = "unknown"
)

type OnboardingError struct {
	Reason  string
	Message string
}

type OnboardingStatus interface {
	Status() string
}

func ReasonForError(err error) string {
	if reason := OnboardingStatus(nil); errors.As(err, &reason) {
		return reason.Status()
	}
	return unknown
}

func (e *OnboardingError) Error() string { return e.Message }

func (e *OnboardingError) Status() string { return e.Reason }

func IsAlreadyOnboarded(err error) bool {
	return ReasonForError(err) == alreadyOnboarded
}

func IsNotFound(err error) bool {
	return ReasonForError(err) == notFound
}
