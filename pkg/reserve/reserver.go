package reserve

type Reserver interface {
	Reserve(requestName string) error
	DeleteReservation() error

	CheckIn() error
	CheckOut() error
}
