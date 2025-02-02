package types

import "time"

type BookingStatus string

const (
	StatusUnknown   BookingStatus = "UNKNOWN"
	StatusPending   BookingStatus = "PENDING"
	StatusConfirmed BookingStatus = "CONFIRMED"
	StatusFailed    BookingStatus = "FAILED"
	StatusCancelled BookingStatus = "CANCELLED"
)

type TravelBooking struct {
	BookingID   string
	UserID      string
	StartDate   time.Time
	EndDate     time.Time
	TotalAmount float64
	Status      BookingStatus

	// Individual bookings
	HotelBooking  *HotelBooking
	FlightBooking *FlightBooking
	CarBooking    *CarBooking
}

type HotelBooking struct {
	HotelID    string
	RoomType   string
	Price      float64
	Status     BookingStatus
	BookingRef string
}

type FlightBooking struct {
	FlightNumber string
	SeatClass    string
	Price        float64
	Status       BookingStatus
	BookingRef   string
}

type CarBooking struct {
	CarType    string
	Price      float64
	Status     BookingStatus
	BookingRef string
}

type BookingError struct {
	Component  string
	Message    string
	RetryCount int
}
