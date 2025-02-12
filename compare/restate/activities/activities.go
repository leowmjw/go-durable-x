package activities

import (
	"context"
	"fmt"
	"github.com/leowmjw/go-durable-x/temporal/types"
	"log/slog"
	"math/rand"
	"time"
)

// Activities implementation
type Activities struct {
	logger *slog.Logger
}

func NewActivities(logger *slog.Logger) *Activities {
	return &Activities{
		logger: logger,
	}
}

// Hotel Activities
func (a *Activities) BookHotel(ctx context.Context, booking *types.HotelBooking) error {
	//spew.Dump(booking)
	// Simulate external API call
	time.Sleep(time.Second)

	// Simulate random failure
	if rand.Float32() < 0.8 { // 20% chance of failure
		return fmt.Errorf("hotel booking failed: service unavailable")
	}

	booking.BookingRef = fmt.Sprintf("HTL-%d", rand.Int31())
	booking.Status = types.StatusConfirmed

	a.logger.Info("Hotel booked successfully",
		slog.String("booking_ref", booking.BookingRef),
		slog.String("hotel_id", booking.HotelID))

	return nil
}

func (a *Activities) CancelHotel(ctx context.Context, bookingRef string) error {
	// Simulate external API call
	time.Sleep(time.Second)

	a.logger.Info("Hotel booking cancelled",
		slog.String("booking_ref", bookingRef))

	return nil
}

// Flight Activities
func (a *Activities) BookFlight(ctx context.Context, booking *types.FlightBooking) error {
	// Simulate external API call
	time.Sleep(time.Second)

	// Simulate random failure
	if rand.Float32() < 0.2 { // 20% chance of failure
		return fmt.Errorf("flight booking failed: no seats available")
	}

	booking.BookingRef = fmt.Sprintf("FLT-%d", rand.Int31())
	booking.Status = types.StatusConfirmed

	a.logger.Info("Flight booked successfully",
		slog.String("booking_ref", booking.BookingRef),
		slog.String("flight_number", booking.FlightNumber))

	return nil
}

func (a *Activities) CancelFlight(ctx context.Context, bookingRef string) error {
	// Simulate external API call
	time.Sleep(time.Second)

	a.logger.Info("Flight booking cancelled",
		slog.String("booking_ref", bookingRef))

	return nil
}

// Car Activities
func (a *Activities) BookCar(ctx context.Context, booking *types.CarBooking) error {
	// Simulate external API call
	time.Sleep(time.Second)

	// Simulate random failure
	if rand.Float32() < 0.2 { // 20% chance of failure
		return fmt.Errorf("car booking failed: no cars available")
	}

	booking.BookingRef = fmt.Sprintf("CAR-%d", rand.Int31())
	booking.Status = types.StatusConfirmed

	a.logger.Info("Car booked successfully",
		slog.String("booking_ref", booking.BookingRef),
		slog.String("car_type", booking.CarType))

	return nil
}

func (a *Activities) CancelCar(ctx context.Context, bookingRef string) error {
	// Simulate external API call
	time.Sleep(time.Second)

	a.logger.Info("Car booking cancelled",
		slog.String("booking_ref", bookingRef))

	return nil
}

// Notification Activities
func (a *Activities) SendEmail(ctx context.Context, to string, subject string, body string) error {
	// Simulate sending email
	time.Sleep(time.Second)

	a.logger.Info("Email sent",
		slog.String("to", to),
		slog.String("subject", subject))

	return nil
}
