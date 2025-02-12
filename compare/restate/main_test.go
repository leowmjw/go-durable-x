package main

// Too complicated; must for now use full integration tests
// Does not seem to have any testsuite capabilities like Temporal for now ..
/*

func TestTravelBookingService(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := &TravelBookingService{logger: logger}

	// Create test server
	srv := server.NewRestate()
	// Register service
	srv.Bind(restate.Reflect(svc))
	// Start server
	go func() {
		if err := srv.Start(context.Background(), "0.0.0.0:9081"); err != nil {
			logger.Error("test server failed", "error", err)
			os.Exit(1)
		}
	}()

	t.Run("happy path - all bookings successful", func(t *testing.T) {
		booking := types.TravelBooking{
			BookingID: "test-123",
			UserID:    "user-123",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			HotelBooking: &types.HotelBooking{
				HotelID:  "hotel-123",
				RoomType: "deluxe",
				Price:    200.0,
			},
			FlightBooking: &types.FlightBooking{
				FlightNumber: "FL123",
				SeatClass:   "economy",
				Price:       300.0,
			},
			CarBooking: &types.CarBooking{
				CarType: "SUV",
				Price:   100.0,
			},
		}

		// Call the service
		err := svc.BookTravel(context.Background(), booking)
		assert.NoError(t, err)
		assert.Equal(t, types.StatusConfirmed, booking.Status)
	})

	t.Run("sad path - flight booking fails", func(t *testing.T) {
		booking := types.TravelBooking{
			BookingID: "test-456",
			UserID:    "user-456",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			HotelBooking: &types.HotelBooking{
				HotelID:  "hotel-456",
				RoomType: "standard",
				Price:    150.0,
			},
			FlightBooking: &types.FlightBooking{
				FlightNumber: "FL456",
				SeatClass:   "business",
				Price:       500.0,
			},
			CarBooking: &types.CarBooking{
				CarType: "compact",
				Price:   80.0,
			},
		}

		// Mock flight booking to fail
		srv.MockService(ServiceName, "BookFlight").Return(nil, assert.AnError)

		// Call the service
		err := svc.BookTravel(context.Background(), booking)
		assert.Error(t, err)
		assert.Equal(t, types.StatusFailed, booking.Status)
	})

	t.Run("sad path - car booking fails", func(t *testing.T) {
		booking := types.TravelBooking{
			BookingID: "test-789",
			UserID:    "user-789",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			HotelBooking: &types.HotelBooking{
				HotelID:  "hotel-789",
				RoomType: "suite",
				Price:    300.0,
			},
			FlightBooking: &types.FlightBooking{
				FlightNumber: "FL789",
				SeatClass:   "first",
				Price:       1000.0,
			},
			CarBooking: &types.CarBooking{
				CarType: "luxury",
				Price:   200.0,
			},
		}

		// Mock car booking to fail
		srv.MockService(ServiceName, "BookCar").Return(nil, assert.AnError)

		// Call the service
		err := svc.BookTravel(context.Background(), booking)
		assert.Error(t, err)
		assert.Equal(t, types.StatusFailed, booking.Status)
	})
}
*/
