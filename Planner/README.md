# Planner Service

The Planner service, a critical component of our application, facilitates event scheduling for users within specific time frames and days of the week at Universidad de los Andes. Developed in Go, this service serves as an API endpoint to calculate time slot availability of users through the week. It is deployed on Google Cloud Platform's App Engine.

## Installation Guide

To install and run the Planner service locally or in your environment, follow these steps:

1. Clone the Repository: `git clone https://github.com/ISIS3510-202410-Team-13/Backend.git`
2. Navigate to the Planner Directory: `cd Backend/Planner`
3. Build the Docker Image: `docker build -t planner .`
4. Run the Docker Container: `docker run -d -p 8080:8080 planner`
5. Verify Installation: `docker ps`

## Usage Guide

### Request Schema

To query the Planner service, make a GET request to the following endpoint:

```
GET http://localhost:8080/
``` 

Include the following JSON body in your request:

```json
{
  "dayOfWeek": "l|m|i|j|v|s|d",
  "startTime": "hhmm",
  "endTime": "hhmm"
}
```

- dayOfWeek: A single letter string representing the day of the week. Valid values are 'l' (Monday), 'm' (Tuesday), 'i' (Wednesday), 'j' (Thursday), 'v' (Friday), 's' (Saturday), and 'd' (Sunday).
- startTime: A four-digit string representing the start time in 24-hour format (e.g., "0800" for 8:00 AM).
- endTime: A four-digit string representing the end time in 24-hour format (e.g., "1700" for 5:00 PM).


### Response Schema

The response will be a JSON array containing objects with the following structure:

```json
[
  {
    "DayOfWeek": "string",
    "StartTime": "string",
    "EndTime": "string",
    "UsersAvailable": ["string"],
    "AmountAvailable": "integer"
  },
  ...
]
```

- DayOfWeek: The day of the week.
- StartTime: The start time of the event.
- EndTime: The end time of the event.
- UsersAvailable: The users available during this time interval.
- AmountAvailable: The number of users available during this time interval.

Note that the endpoint returns a list of these objects (an empty list is possible).

### Example Call

```bash
curl -X GET http://localhost:8080/ \
-H "Content-Type: application/json" \
-d '{"dayOfWeek": "m", "startTime": "0800", "endTime": "1200"}'
```

## Troubleshooting

If you encounter any issues while using the Planner service, refer to the following troubleshooting guide to identify and resolve common problems:


- **Request Body Missing:**
  - Error Message: `{ "message": "Request body is missing" }`
  - Cause: The request sent to the service is missing the JSON body.
  - Solution: Ensure that your request includes a valid JSON body with the required parameters.

- **Missing Required Fields:**
  - Error Message: `{ "message": "Missing required fields in request body" }`
  - Cause: The request body is missing one or more required fields (dayOfWeek, startTime, endTime).
  - Solution: Include all required fields in the request body before sending the request.

- **Invalid Data Types:**
  - Error Message: `{ "message": "Invalid data types in request body" }`
  - Cause: One or more fields in the request body have invalid data types.
  - Solution: Ensure that all fields in the request body have the correct data types (dayOfWeek and time fields should be strings).

- **Invalid Day of Week:**
  - Error Message: `{ "message": "Invalid day of week" }`
  - Cause: The dayOfWeek parameter in the request body is not a valid single-letter string representing a day of the week.
  - Solution: Check that the dayOfWeek parameter is a single-letter string ('l', 'm', 'i', 'j', 'v', 's', or 'd').

- **Invalid Time Format:**
  - Error Message: `{ "message": "Invalid time format" }`
  - Cause: The startTime or endTime parameter in the request body does not match the expected four-digit string format (hhmm).
  - Solution: Ensure that the time parameters follow the four-digit string format (e.g., "0800" for 8:00 AM).

- **Start Time Must be Before End Time:**
  - Error Message: `{ "message": "Start time must be before end time" }`
  - Cause: The startTime parameter is equal to or later than the endTime parameter in the request body.
  - Solution: Adjust the startTime and endTime parameters to ensure that the start time comes before the end time.

- **Invalid Time Range:**
  - Error Message: `{ "message": "Invalid time range" }`
  - Cause: The startTime or endTime parameter in the request body is outside the valid time range (00:00 - 23:59).
  - Solution: Verify that the time parameters fall within the valid range (00:00 - 23:59).

- **Invalid Minutes:**
  - Error Message: `{ "message": "Invalid minutes" }`
  - Cause: The minutes portion of the startTime or endTime parameter in the request body is greater than 59.
  - Solution: Ensure that the minutes portion of the time parameters does not exceed 59.