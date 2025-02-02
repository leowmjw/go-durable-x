package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	"github.com/leowmjw/go-durable-x/temporal/activities"
	"github.com/leowmjw/go-durable-x/temporal/types"
)

// Workflow constants
const (
	TaskQueueName        = "TravelBookingTaskQueue"
	RetryMaxAttempts     = 3
	RetryInitialInterval = time.Second
	RetryMaxInterval     = time.Hour * 24
)

// Activity functions
func BookHotelActivity(ctx context.Context, booking *types.HotelBooking) error {
	activities := activities.NewActivities(slog.Default())
	return activities.BookHotel(ctx, booking)
}

func CancelHotelActivity(ctx context.Context, bookingRef string) error {
	activities := activities.NewActivities(slog.Default())
	return activities.CancelHotel(ctx, bookingRef)
}

func BookFlightActivity(ctx context.Context, booking *types.FlightBooking) error {
	activities := activities.NewActivities(slog.Default())
	return activities.BookFlight(ctx, booking)
}

func CancelFlightActivity(ctx context.Context, bookingRef string) error {
	activities := activities.NewActivities(slog.Default())
	return activities.CancelFlight(ctx, bookingRef)
}

func BookCarActivity(ctx context.Context, booking *types.CarBooking) error {
	activities := activities.NewActivities(slog.Default())
	return activities.BookCar(ctx, booking)
}

func CancelCarActivity(ctx context.Context, bookingRef string) error {
	activities := activities.NewActivities(slog.Default())
	return activities.CancelCar(ctx, bookingRef)
}

func SendEmailActivity(ctx context.Context, to string, subject string, body string) error {
	activities := activities.NewActivities(slog.Default())
	return activities.SendEmail(ctx, to, subject, body)
}

// Activities interfaces for better testability
type (
	HotelBookingActivities interface {
		BookHotel(ctx context.Context, booking *types.HotelBooking) error
		CancelHotel(ctx context.Context, bookingRef string) error
	}

	FlightBookingActivities interface {
		BookFlight(ctx context.Context, booking *types.FlightBooking) error
		CancelFlight(ctx context.Context, bookingRef string) error
	}

	CarBookingActivities interface {
		BookCar(ctx context.Context, booking *types.CarBooking) error
		CancelCar(ctx context.Context, bookingRef string) error
	}

	NotificationActivities interface {
		SendEmail(ctx context.Context, to string, subject string, body string) error
	}
)

// TravelBookingWorkflow orchestrates the entire booking process
func TravelBookingWorkflow(ctx workflow.Context, booking types.TravelBooking) error {
	logger := workflow.GetLogger(ctx)

	// Setup retry policy for activities
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    RetryInitialInterval,
		BackoffCoefficient: 2.0,
		MaximumInterval:    RetryMaxInterval,
		MaximumAttempts:    RetryMaxAttempts,
	}

	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
		RetryPolicy:         retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	// Step 1: Book Hotel
	var hotelBookingErr error
	err := workflow.ExecuteActivity(ctx, BookHotelActivity, booking.HotelBooking).Get(ctx, &hotelBookingErr)
	if err != nil {
		logger.Error("Failed to book hotel", slog.String("error", err.Error()))
		return err
	}

	// Step 2: Book Flight
	var flightBookingErr error
	err = workflow.ExecuteActivity(ctx, BookFlightActivity, booking.FlightBooking).Get(ctx, &flightBookingErr)
	if err != nil {
		// Compensate: Cancel Hotel
		_ = workflow.ExecuteActivity(ctx, CancelHotelActivity, booking.HotelBooking.BookingRef).Get(ctx, nil)
		logger.Error("Failed to book flight", slog.String("error", err.Error()))
		return err
	}

	// Step 3: Book Car
	var carBookingErr error
	err = workflow.ExecuteActivity(ctx, BookCarActivity, booking.CarBooking).Get(ctx, &carBookingErr)
	if err != nil {
		// Compensate: Cancel Flight and Hotel
		_ = workflow.ExecuteActivity(ctx, CancelFlightActivity, booking.FlightBooking.BookingRef).Get(ctx, nil)
		_ = workflow.ExecuteActivity(ctx, CancelHotelActivity, booking.HotelBooking.BookingRef).Get(ctx, nil)
		logger.Error("Failed to book car", slog.String("error", err.Error()))
		return err
	}

	// All bookings successful
	booking.Status = types.StatusConfirmed

	// Send confirmation email
	err = workflow.ExecuteActivity(ctx, SendEmailActivity,
		"user@example.com",
		"Travel Booking Confirmed",
		fmt.Sprintf("Your travel booking %s has been confirmed", booking.BookingID)).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to send confirmation email", slog.String("error", err.Error()))
		// Non-critical error, don't fail the workflow
	}

	return nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create Temporal client
	c, err := client.NewClient(client.Options{})
	if err != nil {
		logger.Error("Failed to create Temporal client", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer c.Close()

	// Create worker
	w := worker.New(c, TaskQueueName, worker.Options{})

	// Register workflow and activities
	w.RegisterWorkflow(TravelBookingWorkflow)
	w.RegisterActivity(BookHotelActivity)
	w.RegisterActivity(CancelHotelActivity)
	w.RegisterActivity(BookFlightActivity)
	w.RegisterActivity(CancelFlightActivity)
	w.RegisterActivity(BookCarActivity)
	w.RegisterActivity(CancelCarActivity)
	w.RegisterActivity(SendEmailActivity)

	// Start worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		logger.Error("Worker failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
