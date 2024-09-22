package coworking

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

var testSuite = &testsuite.WorkflowTestSuite{}

// TestBookingWorkflow_Success ensures the booking workflow completes successfully.
func TestBookingWorkflow_Success(t *testing.T) {
	t.Parallel()

	testEnv := testSuite.NewTestWorkflowEnvironment()
	testEnv.RegisterWorkflow(BookingWorkflow)

	// Mock availability checking activity
	testEnv.OnActivity(CheckAvailability, mock.Anything, "Room A", "2024-09-20", "10:00-14:00").Return(true, nil)

	// Mock payment processing
	testEnv.OnActivity(ProcessPayment, mock.Anything, "user1", 100.00).Return("payment-confirmation-id", nil)

	// Execute the booking workflow
	testEnv.ExecuteWorkflow(BookingWorkflow, BookingInput{
		Room:     "Room A",
		Date:     "2024-09-20",
		TimeSlot: "10:00-14:00",
		UserID:   "user1",
		Price:    100.00,
	})

	// Assert that the workflow completed successfully
	assert.True(t, testEnv.IsWorkflowCompleted())
	assert.NoError(t, testEnv.GetWorkflowError())

	var result BookingResult
	assert.NoError(t, testEnv.GetWorkflowResult(&result))
	assert.True(t, result.Confirmed)
	assert.Equal(t, "payment-confirmation-id", result.PaymentConfirmationID)
	assert.Empty(t, result.ErrorMessage)
}

// TestBookingWorkflow_RoomUnavailable ensures booking fails when the room is unavailable.
func TestBookingWorkflow_RoomUnavailable(t *testing.T) {
	t.Parallel()

	testEnv := testSuite.NewTestWorkflowEnvironment()
	testEnv.RegisterWorkflow(BookingWorkflow)

	// Mock availability check returning unavailable
	testEnv.OnActivity(CheckAvailability, mock.Anything, "Room A", "2024-09-20", "10:00-14:00").Return(false, nil)

	// Execute the booking workflow
	testEnv.ExecuteWorkflow(BookingWorkflow, BookingInput{
		Room:     "Room A",
		Date:     "2024-09-20",
		TimeSlot: "10:00-14:00",
		UserID:   "user1",
		Price:    100.00,
	})

	// Assert that the workflow completed (even though it failed to book)
	assert.True(t, testEnv.IsWorkflowCompleted())
	assert.NoError(t, testEnv.GetWorkflowError())

	var result BookingResult
	assert.NoError(t, testEnv.GetWorkflowResult(&result))
	assert.False(t, result.Confirmed)
	assert.NotEmpty(t, result.ErrorMessage)
}

// TestBookingWorkflow_InvalidTimeSlot tests invalid time format or booking in the past.
func TestBookingWorkflow_InvalidTimeSlot(t *testing.T) {
	t.Parallel()

	testEnv := testSuite.NewTestWorkflowEnvironment()
	testEnv.RegisterWorkflow(BookingWorkflow)

	// Simulate an invalid time slot (e.g., past or incorrect format)
	invalidDate := "invalid-date"
	testEnv.ExecuteWorkflow(BookingWorkflow, BookingInput{
		Room:     "Room A",
		Date:     invalidDate,
		TimeSlot: "10:00-14:00",
		UserID:   "user1",
		Price:    100.00,
	})

	assert.True(t, testEnv.IsWorkflowCompleted())
	assert.NoError(t, testEnv.GetWorkflowError())

	var result BookingResult
	assert.NoError(t, testEnv.GetWorkflowResult(&result))
	assert.False(t, result.Confirmed)
	assert.NotEmpty(t, result.ErrorMessage)
}

// Below better with an actual server; for full integration test ..
// TestBookingWorkflow_ConcurrentBooking tests concurrent booking requests for the same time slot.
func TestBookingWorkflow_ConcurrentBooking(t *testing.T) {
	t.Parallel()

	testEnv := testSuite.NewTestWorkflowEnvironment()
	testEnv.RegisterWorkflow(BookingWorkflow)
	// Concurrent env will run ..
	testEnvConc := testSuite.NewTestWorkflowEnvironment()
	testEnvConc.RegisterWorkflow(BookingWorkflow)

	// Simulate concurrent booking for the same time slot
	testEnv.OnActivity(CheckAvailability, mock.Anything, "Room A", "2024-09-20", "10:00-14:00").Return(true, nil).Once()
	testEnv.OnActivity(CheckAvailability, mock.Anything, "Room A", "2024-09-20", "10:00-14:00").Return(false, nil).Once()

	// Execute the booking workflow twice to simulate concurrency
	testEnv.ExecuteWorkflow(BookingWorkflow, BookingInput{
		Room:     "Room A",
		Date:     "2024-09-20",
		TimeSlot: "10:00-14:00",
		UserID:   "user1",
		Price:    100.00,
	})

	var result1 BookingResult
	assert.NoError(t, testEnv.GetWorkflowResult(&result1))
	assert.True(t, result1.Confirmed)

	testEnvConc.ExecuteWorkflow(BookingWorkflow, BookingInput{
		Room:     "Room A",
		Date:     "2024-09-20",
		TimeSlot: "10:00-14:00",
		UserID:   "user2",
		Price:    100.00,
	})

	var result2 BookingResult
	assert.NoError(t, testEnv.GetWorkflowResult(&result2))
	assert.False(t, result2.Confirmed)
	assert.NotEmpty(t, result2.ErrorMessage)
}

// TestBookingWorkflow_PaymentFailure tests failure during payment processing.
func TestBookingWorkflow_PaymentFailure(t *testing.T) {
	t.Parallel()

	testEnv := testSuite.NewTestWorkflowEnvironment()
	testEnv.RegisterWorkflow(BookingWorkflow)

	// Mock availability as true
	testEnv.OnActivity(CheckAvailability, mock.Anything, "Room A", "2024-09-20", "10:00-14:00").Return(true, nil)

	// Mock payment failure
	testEnv.OnActivity(ProcessPayment, mock.Anything, "user1", 100.00).Return("", errors.New("payment failure"))

	// Execute the booking workflow
	testEnv.ExecuteWorkflow(BookingWorkflow, BookingInput{
		Room:     "Room A",
		Date:     "2024-09-20",
		TimeSlot: "10:00-14:00",
		UserID:   "user1",
		Price:    100.00,
	})

	// Assert that the workflow completed (even though payment failed)
	assert.True(t, testEnv.IsWorkflowCompleted())
	assert.NoError(t, testEnv.GetWorkflowError())

	var result BookingResult
	assert.NoError(t, testEnv.GetWorkflowResult(&result))
	assert.False(t, result.Confirmed)
	assert.Empty(t, result.PaymentConfirmationID)
	assert.NotEmpty(t, result.ErrorMessage)
}

// TestBookingWorkflow_OffHoursBooking tests that bookings cannot be made during off-hours.
func TestBookingWorkflow_OffHoursBooking(t *testing.T) {
	t.Parallel()

	testEnv := testSuite.NewTestWorkflowEnvironment()
	testEnv.RegisterWorkflow(BookingWorkflow)

	// Simulate booking during off-hours (e.g., 2 AM)
	testEnv.OnActivity(CheckAvailability, mock.Anything, "Room A", "2024-09-20", "02:00-03:00").Return(false, nil)

	// Execute the booking workflow
	testEnv.ExecuteWorkflow(BookingWorkflow, BookingInput{
		Room:     "Room A",
		Date:     "2024-09-20",
		TimeSlot: "02:00-03:00",
		UserID:   "user1",
		Price:    100.00,
	})

	// Assert that the workflow failed due to off-hours booking
	assert.True(t, testEnv.IsWorkflowCompleted())
	assert.NoError(t, testEnv.GetWorkflowError())

	var result BookingResult
	assert.NoError(t, testEnv.GetWorkflowResult(&result))
	assert.False(t, result.Confirmed)
	assert.NotEmpty(t, result.ErrorMessage)
}

func TestBookingWorkflow_RescheduleSuccess(t *testing.T) {
	testEnv := testSuite.NewTestWorkflowEnvironment()
	testEnv.RegisterWorkflow(BookingWorkflow)

	// Mock initial booking
	testEnv.OnActivity(CheckAvailability, mock.Anything, "Room A", "2024-09-20", "10:00-14:00").Return(true, nil)
	testEnv.OnActivity(ProcessPayment, mock.Anything, "user1", 100.00).Return("payment-confirmation-id", nil)

	// Mock rescheduling
	testEnv.OnActivity(CheckAvailability, mock.Anything, "Room B", "2024-09-21", "14:00-18:00").Return(true, nil)
	testEnv.OnActivity(ProcessPayment, mock.Anything, "user1", 50.00).Return("reschedule-payment-id", nil)

	testEnv.RegisterDelayedCallback(func() {
		testEnv.SignalWorkflow("reschedule", RescheduleInput{
			NewRoom: "Room B", NewDate: "2024-09-21", NewTimeSlot: "14:00-18:00", RescheduleFee: 50.00,
		})
	}, time.Hour*2)

	testEnv.ExecuteWorkflow(BookingWorkflow, BookingInput{
		Room: "Room A", Date: "2024-09-20", TimeSlot: "10:00-14:00", UserID: "user1", Price: 100.00,
	})

	testEnv.AssertExpectations(t)

	var result BookingResult
	testEnv.GetWorkflowResult(&result)

	assert.True(t, result.Confirmed)
	assert.Equal(t, "reschedule-payment-id", result.PaymentConfirmationID)
	assert.Empty(t, result.ErrorMessage)
}

func TestBookingWorkflow_RescheduleTooLate(t *testing.T) {
	testEnv := testSuite.NewTestWorkflowEnvironment()
	testEnv.RegisterWorkflow(BookingWorkflow)

	// Mock initial booking
	testEnv.OnActivity(CheckAvailability, mock.Anything, "Room A", "2024-09-20", "10:00-14:00").Return(true, nil)
	testEnv.OnActivity(ProcessPayment, mock.Anything, "user1", 100.00).Return("payment-confirmation-id", nil)

	testEnv.RegisterDelayedCallback(func() {
		testEnv.SignalWorkflow("reschedule", RescheduleInput{
			NewRoom: "Room B", NewDate: "2024-09-21", NewTimeSlot: "14:00-18:00", RescheduleFee: 50.00,
		})
	}, time.Hour*24*7) // Simulate rescheduling attempt on the day of booking

	testEnv.ExecuteWorkflow(BookingWorkflow, BookingInput{
		Room: "Room A", Date: "2024-09-20", TimeSlot: "10:00-14:00", UserID: "user1", Price: 100.00,
	})

	testEnv.AssertExpectations(t)

	var result BookingResult
	testEnv.GetWorkflowResult(&result)

	assert.True(t, result.Confirmed)
	assert.Equal(t, "payment-confirmation-id", result.PaymentConfirmationID)
	assert.Contains(t, result.ErrorMessage, "Rescheduling not allowed on the day of booking")
}

func TestBookingWorkflow_RescheduleMaxAttemptsExceeded(t *testing.T) {
	testEnv := testSuite.NewTestWorkflowEnvironment()
	testEnv.RegisterWorkflow(BookingWorkflow)

	// Mock initial booking
	testEnv.OnActivity(CheckAvailability, mock.Anything, "Room A", "2024-09-20", "10:00-14:00").Return(true, nil)
	testEnv.OnActivity(ProcessPayment, mock.Anything, "user1", 100.00).Return("payment-confirmation-id", nil)

	// Mock successful rescheduling attempts
	for i := 0; i < 6; i++ {
		testEnv.OnActivity(CheckAvailability, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
		testEnv.OnActivity(ProcessPayment, mock.Anything, "user1", 50.00).Return(fmt.Sprintf("reschedule-payment-id-%d", i), nil)
	}

	// Register 6 successful reschedules
	for i := 0; i < 6; i++ {
		testEnv.RegisterDelayedCallback(func() {
			testEnv.SignalWorkflow("reschedule", RescheduleInput{
				NewRoom: fmt.Sprintf("Room %c", 'B'+i), NewDate: "2024-09-21", NewTimeSlot: "14:00-18:00", RescheduleFee: 50.00,
			})
		}, time.Hour*24*time.Duration(i+1))
	}

	// Register 7th reschedule attempt (should fail)
	testEnv.RegisterDelayedCallback(func() {
		testEnv.SignalWorkflow("reschedule", RescheduleInput{
			NewRoom: "Room H", NewDate: "2024-09-22", NewTimeSlot: "14:00-18:00", RescheduleFee: 50.00,
		})
	}, time.Hour*24*7)

	testEnv.ExecuteWorkflow(BookingWorkflow, BookingInput{
		Room: "Room A", Date: "2024-09-20", TimeSlot: "10:00-14:00", UserID: "user1", Price: 100.00,
	})

	testEnv.AssertExpectations(t)

	var result BookingResult
	testEnv.GetWorkflowResult(&result)

	assert.True(t, result.Confirmed)
	assert.Equal(t, "reschedule-payment-id-5", result.PaymentConfirmationID)
	assert.Contains(t, result.ErrorMessage, "Maximum rescheduling attempts exceeded")
}

// Test Cases for the day itself ..

func TestRoomEntryWorkflow_Success(t *testing.T) {
	t.Parallel()

	testEnv := testSuite.NewTestWorkflowEnvironment()

	// Register the workflow
	testEnv.RegisterWorkflow(RoomEntryWorkflow)

	// Mock access validation for the user
	testEnv.OnActivity(ValidateRoomAccess, mock.Anything, "booking-id", "user1").Return(true, nil)

	// In TestRoomEntryWorkflow_Success and TestRoomEntryWorkflow_AccessDenied,
	// replace the ExecuteWorkflow call with:

	testEnv.ExecuteWorkflow(RoomEntryWorkflow, RoomEntryInput{
		BookingID: "booking-id",
		UserID:    "user1",
	})

	// Also, add assertions for the result:

	var result RoomEntryResult
	err := testEnv.GetWorkflowResult(&result)
	assert.NoError(t, err)
	assert.True(t, result.Granted)
	assert.Equal(t, "Access granted", result.Message)

}

func TestRoomEntryWorkflow_AccessDenied(t *testing.T) {
	t.Parallel()

	testEnv := testSuite.NewTestWorkflowEnvironment()

	// Register the workflow
	testEnv.RegisterWorkflow(RoomEntryWorkflow)

	// Mock access denial
	testEnv.OnActivity(ValidateRoomAccess, mock.Anything, "booking-id", "user1").Return(false, errors.New("access denied"))

	// In TestRoomEntryWorkflow_Success and TestRoomEntryWorkflow_AccessDenied,
	// replace the ExecuteWorkflow call with:

	// Execute the room entry workflow
	testEnv.ExecuteWorkflow(RoomEntryWorkflow, RoomEntryInput{
		BookingID: "booking-id",
		UserID:    "user1",
	})

	// Also, add assertions for the result:

	var result RoomEntryResult
	err := testEnv.GetWorkflowResult(&result)
	assert.NoError(t, err)

	// For the AccessDenied test, adjust the last two assertions:
	assert.False(t, result.Granted)
	assert.Equal(t, "Access denied", result.Message)

}
