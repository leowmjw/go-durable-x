package coworking

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

var testSuite = &testsuite.WorkflowTestSuite{}

// TestBookingWorkflow_Success ensures the booking workflow completes successfully.
func TestBookingWorkflow_Success(t *testing.T) {
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
