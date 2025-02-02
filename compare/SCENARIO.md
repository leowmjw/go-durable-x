# SCENARIO

## Objective

Describes a Saga pattern where it coordinates with eventual consistency the following 3 independent systems:

- Hotel Booking System
- Flight Booking System
- Car Rental System

It will simulate a distributed system where the user will be travelling to a foreign country for a holiday.

For a start; we will ignore Payment Systems; might add it for more complex demo


## User Acceptance Cases

### Happy Flow

- User books a hotel, books a flight and books a car successfully

### Sad Flows

- User books a hotel, tries to book a flight but fails; cancels hotel
- User books a hotel; is unsuccessful.  The system will retry twice the first day, once per day for 1 week;
      if it continues to fail.  If succeeds, sends email to user to signal to continue then merge back to main flow.
- User books a hotel, books a flight; but unsuccessfully books car; cancels flight, hotel
- User books a hotel, books a flight; but unsuccessfully books car.  Awaits signal from user to accept it.  If not accepted; cancels flight, hotel
- User books a hotel, books a flight and books a car successfully; after 1 day receives signal notification 
      that hotel booking canceled; cancel flight, car
- User books a hotel, books a flight and books a car successfully; 2 days before flight receives signal notification
  that flight booking canceled; cancel hotel, car
- User books a hotel, books a flight and books a car successfully; on day of flight receives signal email
  that car booking canceled.  Awaits signal from user to accept it.  If not accepted; cancels flight, hotel capturing fatal errors

