package domain

type BookingDetails struct {
	UserId    int32
	BookingId int32
	ShowId    int32
	Seats     []string

	BookingStatus         string
	SeatReservationStatus string
	PaymentStatus         string
}
