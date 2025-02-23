package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/leowmjw/go-durable-x/restate/activities"

	"github.com/leowmjw/go-durable-x/temporal/types"
	restate "github.com/restatedev/sdk-go"
	"github.com/restatedev/sdk-go/server"
)

const (
	ServiceName       = "TravelBookingService"
	RetryMaxAttempts  = 3
	RetryInitialDelay = time.Second
	RetryMaxDelay     = time.Hour * 24
)

// TravelBookingService handles the travel booking workflow
type TravelBookingService struct {
	logger           *slog.Logger
	bookingDetails   types.TravelBooking
	funcBookHotel    func(ctx context.Context, booking *types.HotelBooking) error
	funcCancelHotel  func(ctx restate.Context, bookingRef string) error
	mockBookFlight   func(ctx restate.Context, booking *types.FlightBooking) error
	mockCancelFlight func(ctx restate.Context, bookingRef string) error
	mockBookCar      func(ctx restate.Context, booking *types.CarBooking) error
	mockCancelCar    func(ctx restate.Context, bookingRef string) error
	mockSendEmail    func(ctx restate.Context, to, subject, body string) error
}

// BookHotel handles hotel booking
func (s *TravelBookingService) BookHotel(ctx restate.Context, booking *types.HotelBooking) error {
	if s.funcBookHotel == nil {
		return restate.TerminalError(fmt.Errorf("funcBookHotel is nil"), 4404)
	}
	//restate.Run()
	s.logger.Info("booking hotel", "hotelId", booking.HotelID)
	if err := s.funcBookHotel(ctx, booking); err != nil {
		return fmt.Errorf("funcBookHotel: %w", err)
	}
	// Attach the data ..
	s.bookingDetails.HotelBooking = booking
	spew.Dump(s.bookingDetails)
	return nil
}

// CancelHotel handles hotel cancellation
func (s *TravelBookingService) CancelHotel(ctx restate.Context, bookingRef string) error {
	if s.funcCancelHotel == nil {
		return fmt.Errorf("funcCancelHotel is nil")
	}
	s.logger.Info("cancelling hotel", "bookingRef", bookingRef)
	if err := s.funcCancelHotel(ctx, bookingRef); err != nil {
		return fmt.Errorf("funcCancelHotel: %w", err)
	}
	// DONE!!
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

func (s TravelBookingService) Greet(ctx restate.Context, fullName string) (string, error) {

	greetingID := restate.Rand(ctx).UUID()
	spew.Dump(greetingID)
	s.logger.Info("GID:", "SEQ:", greetingID.ClockSequence(), "NID:", greetingID.NodeID())

	return fmt.Sprintf("Goodbye cruel world!! ID: %s", greetingID.String()), nil
}

func (s TravelBookingService) Goodbye(ctx restate.Context) string {
	s.logger.Info("IN Goodbye")
	spew.Dump(ctx.Request())
	/*
		// Example data sent in via Request after route through the Ingress ...
		2025/02/09 21:58:16 INFO Handling invocation method=TravelBookingService/Goodbye invocationID=inv_13fsJWTcpGng41VPxq9Rfu0inTyQOfTKq5
				time=2025-02-09T21:58:16.746+08:00 level=INFO msg="IN Goodbye"
				(*state.Request)(0x1400029a058)({
				 ID: ([]uint8) (len=24 cap=24) {
				  00000000  62 de d3 ea fe 10 d9 25  01 94 eb 00 f7 a4 db 55  |b......%.......U|
				  00000010  cc 2e 01 03 df a1 73 84                           |......s.|
				 },
				 Headers: (map[string]string) (len=3) {
				  (string) (len=22) "x-restate-ingress-path": (string) (len=29) "/TravelBookingService/Goodbye",
				  (string) (len=10) "user-agent": (string) (len=10) "curl/8.1.2",
				  (string) (len=6) "accept": (string) (len=3) ""
			},
			AttemptHeaders: (map[string][]string) (len=3) {
			(string) (len=12) "Content-Type": ([]string) (len=1 cap=1) {
			(string) (len=37) "application/vnd.restate.invocation.v1"
			},
			(string) (len=6) "Accept": ([]string) (len=1 cap=1) {
			(string) (len=37) "application/vnd.restate.invocation.v1"
			},
			(string) (len=23) "X-Restate-Invocation-Id": ([]string) (len=1 cap=1) {
			(string) (len=38) "inv_13fsJWTcpGng41VPxq9Rfu0inTyQOfTKq5"
			}
			},
			Body: ([]uint8) <nil>
			})
	*/
	return "Adios!!"
}

// BookTravel is the main workflow handler
func (s *TravelBookingService) BookTravel(ctx restate.Context, booking types.TravelBooking) error {
	// Attach initial data that came from UI; like userid  + travel date
	// Generate a durable BookingID to be used
	s.bookingDetails.BookingID = restate.Rand(ctx).UUID().String()
	s.logger.Info("starting travel booking workflow", "bookingId", booking.BookingID)

	// Put some fake data
	booking.HotelBooking = &types.HotelBooking{
		HotelID:  "HTL-123456",
		RoomType: "Deluxe",
		Price:    100,
	}
	// Book hotel
	if _, err := restate.Service[restate.Void](ctx, ServiceName, "BookHotel").Request(booking.HotelBooking); err != nil {
		s.logger.Error("failed to book hotel", "error", err)
		return err
	}
	// Book flight
	if _, err := restate.Service[restate.Void](ctx, ServiceName, "BookFlight").Request(booking.FlightBooking); err != nil {
		s.logger.Error("failed to book flight", "error", err)
		// Compensate hotel booking
		if _, cerr := restate.Service[restate.Void](ctx, ServiceName, "CancelHotel").Request(booking.HotelBooking.BookingRef); cerr != nil {
			s.logger.Error("failed to cancel hotel", "error", cerr)
		}
		return err
	}
	/*

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
	*/

	booking.Status = types.StatusConfirmed
	return nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	a := activities.NewActivities(logger)
	svc := &TravelBookingService{
		logger:         logger,
		bookingDetails: types.TravelBooking{},
		funcBookHotel:  a.BookHotel,
	}

	// Create and register the service
	if err := server.NewRestate().
		Bind(restate.Reflect(svc)).
		Start(context.Background(), "0.0.0.0:9080"); err != nil {
		logger.Error("service failed", "error", err)
		os.Exit(1)
	}
}
