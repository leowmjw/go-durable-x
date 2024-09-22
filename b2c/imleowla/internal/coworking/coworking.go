package coworking

import "go.temporal.io/sdk/workflow"

// RoomEntryInput represents the input for the room entry workflow
type RoomEntryInput struct {
	BookingID string
	UserID    string
}

// RoomEntryResult represents the output of the room entry workflow
type RoomEntryResult struct {
	Granted bool
	Message string
}

// RoomEntryWorkflow handles the process of validating and granting room access
func RoomEntryWorkflow(ctx workflow.Context, input RoomEntryInput) (RoomEntryResult, error) {
	var result RoomEntryResult

	// Execute the ValidateRoomAccess activity
	var accessGranted bool
	err := workflow.ExecuteActivity(ctx, ValidateRoomAccess, input.BookingID, input.UserID).Get(ctx, &accessGranted)
	if err != nil {
		return result, err
	}

	if accessGranted {
		result.Granted = true
		result.Message = "Access granted"
	} else {
		result.Granted = false
		result.Message = "Access denied"
	}

	return result, nil
}
