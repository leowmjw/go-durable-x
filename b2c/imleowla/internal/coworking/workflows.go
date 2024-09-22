package coworking

import (
	"errors"
	"time"

	"go.temporal.io/sdk/workflow"
)

// BookingInput represents the input parameters for the booking workflow
type BookingInput struct {
	Room     string
	Date     string
	TimeSlot string
	UserID   string
	Price    float64
}

// BookingResult represents the output of the booking workflow
type BookingResult struct {
	Confirmed             bool
	PaymentConfirmationID string
	ErrorMessage          string
}

// BookingWorkflow handles the process of booking a coworking space
type RescheduleInput struct {
	NewRoom       string
	NewDate       string
	NewTimeSlot   string
	RescheduleFee float64
}

func BookingWorkflow(ctx workflow.Context, input BookingInput) (BookingResult, error) {
	var result BookingResult
	currentBooking := input
	rescheduleAttempts := 0
	maxRescheduleAttempts := 6

	err := workflow.SetQueryHandler(ctx, "getBookingStatus", func() (BookingInput, error) {
		return currentBooking, nil
	})
	if err != nil {
		return result, err
	}

	// Initial booking logic
	var isAvailable bool
	err = workflow.ExecuteActivity(ctx, CheckAvailability, input.Room, input.Date, input.TimeSlot).Get(ctx, &isAvailable)
	if err != nil {
		result.ErrorMessage = "Error checking availability"
		return result, err
	}
	if !isAvailable {
		result.ErrorMessage = "Room not available"
		return result, nil
	}

	var paymentConfirmationID string
	err = workflow.ExecuteActivity(ctx, ProcessPayment, input.UserID, input.Price).Get(ctx, &paymentConfirmationID)
	if err != nil {
		result.ErrorMessage = "Payment failed"
		return result, err
	}

	result.Confirmed = true
	result.PaymentConfirmationID = paymentConfirmationID

	// Handle rescheduling
	rescheduleChannel := workflow.GetSignalChannel(ctx, "reschedule")
	for {
		var rescheduleInput RescheduleInput
		if !rescheduleChannel.ReceiveAsync(&rescheduleInput) {
			break
		}

		if rescheduleAttempts >= maxRescheduleAttempts {
			result.ErrorMessage = "Maximum rescheduling attempts exceeded"
			break
		}

		bookingDate, _ := time.Parse("2006-01-02", currentBooking.Date)
		if time.Now().After(bookingDate) {
			result.ErrorMessage = "Rescheduling not allowed on the day of booking"
			break
		}

		err := handleReschedule(ctx, &currentBooking, rescheduleInput, &result)
		if err != nil {
			result.ErrorMessage = "Rescheduling failed: " + err.Error()
			break
		}

		rescheduleAttempts++
	}

	return result, nil
}

func handleReschedule(ctx workflow.Context, currentBooking *BookingInput, rescheduleInput RescheduleInput, result *BookingResult) error {
	var isAvailable bool
	err := workflow.ExecuteActivity(ctx, CheckAvailability, rescheduleInput.NewRoom, rescheduleInput.NewDate, rescheduleInput.NewTimeSlot).Get(ctx, &isAvailable)
	if err != nil {
		return errors.New("Error checking availability for new booking")
	}
	if !isAvailable {
		return errors.New("New room not available")
	}

	var paymentConfirmationID string
	err = workflow.ExecuteActivity(ctx, ProcessPayment, currentBooking.UserID, rescheduleInput.RescheduleFee).Get(ctx, &paymentConfirmationID)
	if err != nil {
		return errors.New("Rescheduling payment failed")
	}

	currentBooking.Room = rescheduleInput.NewRoom
	currentBooking.Date = rescheduleInput.NewDate
	currentBooking.TimeSlot = rescheduleInput.NewTimeSlot
	result.PaymentConfirmationID = paymentConfirmationID

	return nil
}
