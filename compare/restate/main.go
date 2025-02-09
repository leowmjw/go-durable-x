package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/leowmjw/go-durable-x/temporal/types"
	"github.com/restatedev/sdk-go"
	"github.com/restatedev/sdk-go/server"
)

const (
	ServiceName         = "travel-booking"
	RetryMaxAttempts    = 3
	RetryInitialDelay   = time.Second
	RetryMaxDelay       = time.Hour * 24
)

// TravelBookingService handles the travel booking workflow
type TravelBookingService struct {
	logger     *slog.Logger
	mockBookTravel func(ctx restate.Context, booking types.TravelBooking) error
	mockBookHotel  func(ctx restate.Context, booking *types.HotelBooking) error
	mockCancelHotel func(ctx restate.Context, bookingRef string) error
	mockBookFlight  func(ctx restate.Context, booking *types.FlightBooking) error
	mockCancelFlight func(ctx restate.Context, bookingRef string) error
	mockBookCar     func(ctx restate.Context, booking *types.CarBooking) error
	mockCancelCar   func(ctx restate.Context, bookingRef string) error
	mockSendEmail   func(ctx restate.Context, to, subject, body string) error
}

// BookHotel handles hotel booking
func (s *TravelBookingService) BookHotel(ctx restate.Context, booking *types.HotelBooking) error {
	if s.mockBookHotel != nil {
		return s.mockBookHotel(ctx, booking)
	}
	s.logger.Info("booking hotel", "hotelId", booking.HotelID)
	booking.Status = types.StatusConfirmed
	booking.BookingRef = "H123" // Simulated booking reference
	return nil
}

// CancelHotel handles hotel cancellation
func (s *TravelBookingService) CancelHotel(ctx restate.Context, bookingRef string) error {
	if s.mockCancelHotel != nil {
		return s.mockCancelHotel(ctx, bookingRef)
	}
	s.logger.Info("cancelling hotel", "bookingRef", bookingRef)
	return nil
}

// BookFlight handles flight booking
func (s *TravelBookingService) BookFlight(ctx restate.Context, booking *types.FlightBooking) error {
	if s.mockBookFlight != nil {
		return s.mockBookFlight(ctx, booking)
	}
	s.logger.Info("booking flight", "flightNumber", booking.FlightNumber)
	booking.Status = types.StatusConfirmed
	booking.BookingRef = "F123" // Simulated booking reference
	return nil
}

// CancelFlight handles flight cancellation
func (s *TravelBookingService) CancelFlight(ctx restate.Context, bookingRef string) error {
	if s.mockCancelFlight != nil {
		return s.mockCancelFlight(ctx, bookingRef)
	}
	s.logger.Info("cancelling flight", "bookingRef", bookingRef)
	return nil
}

// BookCar handles car booking
func (s *TravelBookingService) BookCar(ctx restate.Context, booking *types.CarBooking) error {
	if s.mockBookCar != nil {
		return s.mockBookCar(ctx, booking)
	}
	s.logger.Info("booking car", "carType", booking.CarType)
	booking.Status = types.StatusConfirmed
	booking.BookingRef = "C123" // Simulated booking reference
	return nil
}

// CancelCar handles car cancellation
func (s *TravelBookingService) CancelCar(ctx restate.Context, bookingRef string) error {
	if s.mockCancelCar != nil {
		return s.mockCancelCar(ctx, bookingRef)
	}
	s.logger.Info("cancelling car", "bookingRef", bookingRef)
	return nil
}

// SendEmail handles email notifications
func (s *TravelBookingService) SendEmail(ctx restate.Context, to string, subject string, body string) error {
	s.logger.Info("sending email", "to", to, "subject", subject)
	return nil
}

// BookTravel is the main workflow handler
func (s *TravelBookingService) BookTravel(ctx restate.Context, booking types.TravelBooking) error {
	if s.mockBookTravel != nil {
		return s.mockBookTravel(ctx, booking)
	}
	s.logger.Info("starting travel booking workflow", "bookingId", booking.BookingID)

	// Book hotel
	if _, err := restate.Service[error](ctx, ServiceName, "BookHotel").Request(booking.HotelBooking); err != nil {
		s.logger.Error("failed to book hotel", "error", err)
		return err
	}

	// Book flight
	if _, err := restate.Service[error](ctx, ServiceName, "BookFlight").Request(booking.FlightBooking); err != nil {
		s.logger.Error("failed to book flight", "error", err)
		// Compensate hotel booking
		if _, cerr := restate.Service[error](ctx, ServiceName, "CancelHotel").Request(booking.HotelBooking.BookingRef); cerr != nil {
			s.logger.Error("failed to cancel hotel", "error", cerr)
		}
		return err
	}

	// Book car
	if _, err := restate.Service[error](ctx, ServiceName, "BookCar").Request(booking.CarBooking); err != nil {
		s.logger.Error("failed to book car", "error", err)
		// Compensate flight and hotel bookings
		if _, cerr := restate.Service[error](ctx, ServiceName, "CancelFlight").Request(booking.FlightBooking.BookingRef); cerr != nil {
			s.logger.Error("failed to cancel flight", "error", cerr)
		}
		if _, cerr := restate.Service[error](ctx, ServiceName, "CancelHotel").Request(booking.HotelBooking.BookingRef); cerr != nil {
			s.logger.Error("failed to cancel hotel", "error", cerr)
		}
		return err
	}

	booking.Status = types.StatusConfirmed
	return nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := &TravelBookingService{logger: logger}

	// Create and register the service
	if err := server.NewRestate().
		Bind(restate.Reflect(svc)).
		Start(context.Background(), "0.0.0.0:9080"); err != nil {
		logger.Error("service failed", "error", err)
		os.Exit(1)
	}
}
