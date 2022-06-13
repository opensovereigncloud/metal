package entity

type Order struct {
	Name      string `json:"requestName"`
	Namespace string `json:"requestNamespace"`
	Ordered   bool   `json:"ordered"`
}

type ReservationStatus string

const (
	ReservationStatusAvailable ReservationStatus = "Available"
	ReservationStatusReserved  ReservationStatus = "Reserved"
	ReservationStatusPending   ReservationStatus = "Pending"
	ReservationStatusError     ReservationStatus = "Error"
	ReservationStatusRunning   ReservationStatus = "Running"
)

type Reservation struct {
	OrderName        string            `json:"orderName"`
	OrderNamespace   string            `json:"orderNamespace"`
	RequestName      string            `json:"requestName"`
	RequestNamespace string            `json:"requestNamespace"`
	Status           ReservationStatus `json:"reservationStatus"`
}

func (o *Order) IsOrdered() bool {
	return o.Ordered
}

func (r *Reservation) IsReserved() bool {
	if r.Status == ReservationStatusAvailable {
		return false
	}
	return r.Status == ReservationStatusPending || r.Status != ReservationStatusError
}
