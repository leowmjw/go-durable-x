# OUTPUT

## Show auto-retry for non-terminal actions

```shell
$ curl -v  http://localhost:8080/TravelBookingService/BookTravel -H "Content-Type: application/json"  
  -d '{ "userID": "user-123", "startDate": "2024-02-08T00:00:00Z", "endDate": "2024-02-12T00:00:00Z", 
  "hotelBooking": {"hotelID": "hotel-123", "roomType": "deluxe", "price": 200.0}
  }'

2025/02/11 20:18:10 INFO Handling invocation method=TravelBookingService/BookTravel invocationID=inv_1jIZ17fELUaZ3FBgMPlZMA0NVNwJCAp8Vr
time=2025-02-11T20:18:10.176+08:00 level=INFO msg="starting travel booking workflow" bookingId=test-123
2025/02/11 20:18:10 INFO Handling invocation method=TravelBookingService/BookHotel invocationID=inv_1jtfSPjN5c277o7ODbSubu7WqpbIez8UyR
time=2025-02-11T20:18:10.183+08:00 level=INFO msg="booking hotel" hotelId=HTL-123456
2025/02/11 20:18:11 ERROR Invocation returned a non-terminal failure method=TravelBookingService/BookHotel invocationID=inv_1jtfSPjN5c277o7ODbSubu7WqpbIez8UyR err="funcBookHotel: hotel booking failed: service unavailable"
2025/02/11 20:18:11 INFO Handling invocation method=TravelBookingService/BookHotel invocationID=inv_1jtfSPjN5c277o7ODbSubu7WqpbIez8UyR
time=2025-02-11T20:18:11.250+08:00 level=INFO msg="booking hotel" hotelId=HTL-123456
2025/02/11 20:18:12 ERROR Invocation returned a non-terminal failure method=TravelBookingService/BookHotel invocationID=inv_1jtfSPjN5c277o7ODbSubu7WqpbIez8UyR err="funcBookHotel: hotel booking failed: service unavailable"
2025/02/11 20:18:12 INFO Handling invocation method=TravelBookingService/BookHotel invocationID=inv_1jtfSPjN5c277o7ODbSubu7WqpbIez8UyR
time=2025-02-11T20:18:12.368+08:00 level=INFO msg="booking hotel" hotelId=HTL-123456
2025/02/11 20:18:13 ERROR Invocation returned a non-terminal failure method=TravelBookingService/BookHotel invocationID=inv_1jtfSPjN5c277o7ODbSubu7WqpbIez8UyR err="funcBookHotel: hotel booking failed: service unavailable"
2025/02/11 20:18:13 INFO Handling invocation method=TravelBookingService/BookHotel invocationID=inv_1jtfSPjN5c277o7ODbSubu7WqpbIez8UyR
time=2025-02-11T20:18:13.615+08:00 level=INFO msg="booking hotel" hotelId=HTL-123456
time=2025-02-11T20:18:14.615+08:00 level=INFO msg="Hotel booked successfully" booking_ref=HTL-705789830 hotel_id=HTL-123456
(types.TravelBooking) {
 BookingID: (string) (len=36) "a47fc16a-b601-4c23-8950-0f556df67c13",
...
```