package main

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
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
	funcCancelHotel  func(ctx context.Context, bookingRef string) error
	funcBookFlight   func(ctx context.Context, booking *types.FlightBooking) error
	funcCancelFlight func(ctx context.Context, bookingRef string) error
	funcBookCar      func(ctx context.Context, booking *types.CarBooking) error
	funcCancelCar    func(ctx context.Context, bookingRef string) error
	funcSendEmail    func(ctx context.Context, to, subject, body string) error
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
	// spew.Dump(s.bookingDetails)
	return nil
}

// CancelHotel handles hotel cancellation
func (s *TravelBookingService) CancelHotel(ctx restate.Context, bookingRef string) error {
	if s.funcCancelHotel == nil {
		return restate.TerminalError(fmt.Errorf("funcCancelHotel is nil"), 4404)
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
	if s.funcBookFlight == nil {
		return restate.TerminalError(fmt.Errorf("funcBookFlight is nil"), 4404)
	}
	s.logger.Info("booking flight", "flightNumber", booking.FlightNumber)
	booking.Status = types.StatusConfirmed
	booking.BookingRef = "F123" // Simulated booking reference
	return nil
}

// CancelFlight handles flight cancellation
func (s *TravelBookingService) CancelFlight(ctx restate.Context, bookingRef string) error {
	if s.funcCancelFlight == nil {
		return restate.TerminalError(fmt.Errorf("funcCancelFlight is nil"), 4404)
	}
	s.logger.Info("cancelling flight", "bookingRef", bookingRef)
	// Call activity with run to simulate external API call
	// Wrap it in restate.Run so it remembers the result
	if err := s.funcCancelFlight(ctx, bookingRef); err != nil {
		return fmt.Errorf("funcCancelFlight: %w", err)
	}
	return nil
}

// BookCar handles car booking
func (s *TravelBookingService) BookCar(ctx restate.Context, booking *types.CarBooking) error {
	if s.funcBookCar == nil {
		return restate.TerminalError(fmt.Errorf("funcBookCar is nil"), 4404)
	}
	s.logger.Info("booking car", "carType", booking.CarType)
	booking.Status = types.StatusConfirmed
	booking.BookingRef = "C123" // Simulated booking reference
	return nil
}

// CancelCar handles car cancellation
func (s *TravelBookingService) CancelCar(ctx restate.Context, bookingRef string) error {
	if s.funcCancelCar == nil {
		return restate.TerminalError(fmt.Errorf("funcCancelCar is nil"), 4404)
	}
	s.logger.Info("cancelling car", "bookingRef", bookingRef)
	if err := s.funcCancelCar(ctx, bookingRef); err != nil {
		return fmt.Errorf("funcCancelCar: %w", err)
	}
	return nil
}

// SendEmail handles email notifications
func (s *TravelBookingService) SendEmail(ctx restate.Context, to string, subject string, body string) error {
	if s.funcSendEmail == nil {
		return restate.TerminalError(fmt.Errorf("funcSendEmail is nil"), 4404)
	}

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
	// fmt.Println("BOOK FOR:", booking.UserID) // booking.UserID
	s.bookingDetails.UserID = booking.UserID
	s.bookingDetails.TotalAmount = booking.TotalAmount
	s.bookingDetails.EndDate = booking.EndDate
	// Generate a durable BookingID to be used
	s.bookingDetails.BookingID = restate.Rand(ctx).UUID().String()
	s.logger.Info("starting travel booking workflow", "bookingId", booking.BookingID)

	// Put some fake data; if needed
	// booking.HotelBooking = &types.HotelBooking{
	// 	HotelID:  "HTL-123456",
	// 	RoomType: "Deluxe",
	// 	Price:    100,
	// }
	// Attach teh data from the example calls

	// Book hotel
	if _, err := restate.Service[restate.Void](ctx, ServiceName, "BookHotel").Request(booking.HotelBooking); err != nil {
		s.logger.Error("failed to book hotel", "error", err)
		s.bookingDetails.Status = types.StatusFailed

		return err
	}
	// Book flight
	if _, err := restate.Service[restate.Void](ctx, ServiceName, "BookFlight").Request(booking.FlightBooking); err != nil {
		s.logger.Error("failed to book flight", "error", err)
		s.bookingDetails.Status = types.StatusFailed

		// Compensate hotel booking
		if _, cerr := restate.Service[restate.Void](ctx, ServiceName, "CancelHotel").Request(booking.HotelBooking.BookingRef); cerr != nil {
			s.logger.Error("failed to cancel hotel", "error", cerr)
			return cerr
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
	fmt.Println("================ FINAL BOOKING DETAILS ======================")
	s.bookingDetails.Status = types.StatusConfirmed
	spew.Dump(s.bookingDetails)
	fmt.Println("********************* FINAL BOOKING DETAILS *********************")

	return nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	a := activities.NewActivities(logger)
	svc := &TravelBookingService{
		logger:         logger,
		bookingDetails: types.TravelBooking{},
		funcBookHotel:  a.BookHotel,
		// funcCancelHotel:  a.CancelHotel,
		funcBookFlight: a.BookFlight,
		// funcCancelFlight: a.CancelFlight,
		// funcBookCar:      a.BookCar,
		// funcCancelCar:    a.CancelCar,
		// funcSendEmail:    a.SendEmail,
	}

	// Start Restate server in a goroutine
	go func() {
		// Create and register the Restate service
		if err := server.NewRestate().
			Bind(restate.Reflect(svc)).
			Start(context.Background(), "0.0.0.0:9080"); err != nil {
			logger.Error("restate service failed", "error", err)
			os.Exit(1)
		}
	}()

	// Create a separate HTTP server for the demo interface
	webServer := &http.Server{
		Addr: "0.0.0.0:8888",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/demo" {
				templatePath := filepath.Join("templates", "demo.html")
				tmpl, err := template.ParseFiles(templatePath)
				if err != nil {
					logger.Error("error parsing template", "error", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				if err := tmpl.Execute(w, nil); err != nil {
					logger.Error("error executing template", "error", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			} else {
				http.NotFound(w, r)
			}
		}),
	}

	// Log startup
	logger.Info("web server started", "address", "0.0.0.0:8888")
	logger.Info("restate server started", "address", "0.0.0.0:9080")
	logger.Info("restate ingress available at", "address", "0.0.0.0:8080")

	// Start the web server
	if err := webServer.ListenAndServe(); err != nil {
		logger.Error("web server failed", "error", err)
		os.Exit(1)
	}
}
