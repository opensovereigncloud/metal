package entity

type Onboarding struct {
	RequestName                   string
	RequestNamespace              string
	InitializationObjectName      string
	InitializationObjectNamespace string
}

func (o *Onboarding) IsOnboarded() bool {
	return true
}
