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

## Implementation

### Currently Implemented Features

1. Basic Workflow Structure
   - Temporal workflow setup with proper error handling
   - Activity definitions for Hotel, Flight, and Car bookings
   - Compensation logic for failed bookings
   - Basic retry policy with configurable attempts

2. Implemented Scenarios
   - Happy Flow: Complete successful booking of hotel, flight, and car
   - Flight Booking Failure: Hotel booking compensation
   - Car Booking Failure: Hotel and flight booking compensation

3. Testing Coverage
   - Unit tests for happy path
   - Unit tests for flight booking failure
   - Unit tests for car booking failure
   - Activity mocking and verification

### Pending Implementation

1. Complex Retry Scenarios
   - Hotel booking with weekly retry pattern
   - Custom retry policies per activity type
   - Progressive retry intervals (twice first day, once per day for a week)

2. User Signal Handling
   - Signal interfaces for user approvals
   - Timeout handling for user responses
   - Signal correlation with specific workflow instances
   - Email notification system integration

3. Time-Based Events
   - Delayed cancellation signals
   - Scheduled verification of bookings
   - Time-based triggers for status checks
   - Handling of booking expiration

4. Advanced Compensation Flows
   - Partial completion acceptance
   - User-approved partial bookings
   - Complex compensation chains
   - Fatal error capture and handling

5. State Management
   - Workflow state persistence
   - Recovery from partial completions
   - State transition logging
   - Audit trail maintenance

### Technical Debt and Future Improvements

1. Monitoring and Observability
   - Enhanced structured logging
   - Metrics collection
   - Performance tracking
   - Error rate monitoring

2. Resilience Features
   - Circuit breakers for external services
   - Fallback mechanisms
   - Rate limiting
   - Timeout configurations

3. Testing Improvements
   - Integration tests with external services
   - Load testing scenarios
   - Chaos testing for failure modes
   - End-to-end workflow testing
