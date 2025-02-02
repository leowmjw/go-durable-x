package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"

	"github.com/leowmjw/go-durable-x/temporal/types"
)

func Test_TravelBookingWorkflow_HappyPath(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activities
	env.OnActivity(BookHotelActivity, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(BookFlightActivity, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(BookCarActivity, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(SendEmailActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	booking := types.TravelBooking{
		BookingID: "TEST-123",
		UserID:    "user-1",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour * 7),
		HotelBooking: &types.HotelBooking{
			HotelID:  "hotel-1",
			RoomType: "deluxe",
			Price:    200.0,
		},
		FlightBooking: &types.FlightBooking{
			FlightNumber: "FL123",
			SeatClass:    "economy",
			Price:        500.0,
		},
		CarBooking: &types.CarBooking{
			CarType: "SUV",
			Price:   100.0,
		},
	}

	env.ExecuteWorkflow(TravelBookingWorkflow, booking)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func Test_TravelBookingWorkflow_FlightFailure(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activities
	env.OnActivity(BookHotelActivity, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(BookFlightActivity, mock.Anything, mock.Anything).Return(
		fmt.Errorf("flight booking failed"))
	env.OnActivity(CancelHotelActivity, mock.Anything, mock.Anything).Return(nil)

	booking := types.TravelBooking{
		BookingID: "TEST-124",
		UserID:    "user-1",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour * 7),
		HotelBooking: &types.HotelBooking{
			HotelID:  "hotel-1",
			RoomType: "deluxe",
			Price:    200.0,
		},
		FlightBooking: &types.FlightBooking{
			FlightNumber: "FL123",
			SeatClass:    "economy",
			Price:        500.0,
		},
		CarBooking: &types.CarBooking{
			CarType: "SUV",
			Price:   100.0,
		},
	}

	env.ExecuteWorkflow(TravelBookingWorkflow, booking)

	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}

func Test_TravelBookingWorkflow_CarFailure(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activities
	env.OnActivity(BookHotelActivity, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(BookFlightActivity, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(BookCarActivity, mock.Anything, mock.Anything).Return(
		fmt.Errorf("car booking failed"))
	env.OnActivity(CancelHotelActivity, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(CancelFlightActivity, mock.Anything, mock.Anything).Return(nil)

	booking := types.TravelBooking{
		BookingID: "TEST-125",
		UserID:    "user-1",
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour * 7),
		HotelBooking: &types.HotelBooking{
			HotelID:  "hotel-1",
			RoomType: "deluxe",
			Price:    200.0,
		},
		FlightBooking: &types.FlightBooking{
			FlightNumber: "FL123",
			SeatClass:    "economy",
			Price:        500.0,
		},
		CarBooking: &types.CarBooking{
			CarType: "SUV",
			Price:   100.0,
		},
	}

	env.ExecuteWorkflow(TravelBookingWorkflow, booking)

	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}
