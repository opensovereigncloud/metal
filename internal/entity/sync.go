package entity

type Synchronization struct {
	SourceName      string
	SourceNamespace string
	TargetName      string
	TargetNamespace string
	SourceStatus    ReservationStatus
	TargetStatus    ReservationStatus
}

func (s *Synchronization) IsSyncNeeded() bool {
	return s.SourceStatus != s.TargetStatus
}
