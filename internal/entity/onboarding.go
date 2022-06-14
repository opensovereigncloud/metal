package entity

type Onboarding struct {
	RequestName                   string
	RequestNamespace              string
	InitializationObjectName      string
	InitializationObjectNamespace string
}

type Initialization struct {
	Require bool
	Error   error
}
