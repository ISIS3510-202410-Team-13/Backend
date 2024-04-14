# Planner Service

The Planner service, a critical component of our application, facilitates event scheduling for users within specific time frames and days of the week at Universidad de los Andes. Developed in Go, this service serves as an API endpoint to calculate time slot availability of users through the week. It is deployed on Google Cloud Platform's App Engine.

## Installation Guide

To install and run the Planner service locally or in your environment, follow these steps:

1. Clone the Repository: `git clone https://github.com/ISIS3510-202410-Team-13/Backend.git`
2. Navigate to the Planner Directory: `cd Backend/Planner`
3. Build the Docker Image: `docker build -t planner .`
4. Run the Docker Container: `docker run -d -p 8080:8080 planner`
5. Verify Installation: `docker ps`

### Build Guide

To build the Planner service from source, follow these steps:

1. Navigate to the Planner Directory: `cd Backend/Planner`
2. Build the Go Binary: `go build -o planner.exe ./cmd/main.go`
3. Run the Binary: `./planner.exe`
4. Verify Installation: `curl http://localhost:8080/hello`

### Test Guide

This project implements unit tests for the Planner service to ensure the correctness of its functionality. To run tests for all the micro-service code, follow these steps:

1. Navigate to the Planner Directory: `cd Backend/Planner`
2. Run the Tests with Coverage: `go test -coverprofile=coverage.out ./...`
3. Generate the Coverage Report: `go tool cover -html=coverage.out -o coverage.html`
4. Open the Coverage Report: `open coverage.html`

## Usage Guide

### Health Checkpoint

Make a GET request to the following endpoint:

```
GET https://planner-dot-unischedule-5ee93.uc.r.appspot.com/hello
```

If the service is running correctly, you will receive the message `Hello from the Planner service!`.

### Request Schema

To query the Planner service, make a POST request to the following endpoint:

```
POST https://planner-dot-unischedule-5ee93.uc.r.appspot.com/planner
``` 

Include a JSON body with the following schema in your request:

```json
{
  "userId": [
    {
      "dayOfWeek": "l|m|i|j|v|s|d",
      "startTime": "hhmm",
      "endTime": "hhmm"
    },
  ]
}
```

This JSON body is an object with a different entry for each user. Each User ID key contains an array of objects representing the time slots where that user is busy (i.e. its events). Each entry has the following structure:

- `userId`: A string representing the user ID. It must be unique. Each user can have zero or more events with the following structure:

  - `dayOfWeek`: A single letter string representing the day of the week. Valid values are 'l' (Monday), 'm' (Tuesday), 'i' (Wednesday), 'j' (Thursday), 'v' (Friday), 's' (Saturday), and 'd' (Sunday).
  - `startTime`: A four-digit string representing the start time in 24-hour format (e.g., "0800" for 8:00 AM).
  - `endTime`: A four-digit string representing the end time in 24-hour format (e.g., "1700" for 5:00 PM).

> Refer to the troubleshooting section for common errors and solutions related to the request schema.

### Response Schema

The response will be a JSON array containing objects with the following structure:

```json
[
  {
    "dayOfWeek": "string",
    "startTime": "string",
    "endTime": "string",
    "usersAvailable": "integer",
    "attendees": ["string"],
    "duration": "integer"
  }
]
```

Each element of this collection represents a time interval where users are available to schedule events. For each day of the week, the service returns values that span over the entire day and which do not overlap (i.e. it is a set partition over the 24-hours of the day). The fields of each object are as follows:

- `dayOfWeek`: The day of the week. This field is a string with the same semantics as the dayOfWeek field in the request.
- `startTime`: The start time of the event. This field is a string with the same semantics as the startTime field in the request.
- `endTime`: The end time of the event. This field is a string with the same semantics as the endTime field in the request.
- `usersAvailable`: The number of users available during this time interval.
- `attendees`: The list of users available during this time interval (by their ID).
- `duration`: The duration of the time interval in minutes.

Note that the endpoint returns a list of these objects (an empty list is possible).

### Example Call

Here's an example of a valid request body:

```json
{
  "user1": [
    {
      "dayOfWeek": "v",
      "startTime": "0900",
      "endTime": "1100"
    },
    {
      "dayOfWeek": "v",
      "startTime": "1400",
      "endTime": "1600"
    }
  ],
  "user2": [
    {
      "dayOfWeek": "v",
      "startTime": "0800",
      "endTime": "1000"
    },
    {
      "dayOfWeek": "v",
      "startTime": "1200",
      "endTime": "1400"
    }
  ],
  "user3": [
    {
      "dayOfWeek": "v",
      "startTime": "1000",
      "endTime": "1200"
    },
    {
      "dayOfWeek": "v",
      "startTime": "1100",
      "endTime": "1300"
    }
  ]
}
```

You can make a POST request to the Planner service with the previous example using the following cURL command:

```bash
curl -X POST https://planner-dot-unischedule-5ee93.uc.r.appspot.com/planner \
-H "Content-Type: application/json" \
-d '{"user1":[{"dayOfWeek":"v","startTime":"0900","endTime":"1100"},{"dayOfWeek":"v","startTime":"1400","endTime":"1600"}],"user2":[{"dayOfWeek":"v","startTime":"0800","endTime":"1000"},{"dayOfWeek":"v","startTime":"1200","endTime":"1400"}],"user3":[{"dayOfWeek":"v","startTime":"1000","endTime":"1200"},{"dayOfWeek":"v","startTime":"1100","endTime":"1300"}]}'
```

This command will return a JSON response like the following: 

```json
[
  {
    "dayOfWeek": "m",
    "startTime": "00:00",
    "endTime": "23:59",
    "usersAvailable": 3,
    "attendees": [
      "user1",
      "user2",
      "user3"
    ],
    "duration": 1440
  },
  {
    "dayOfWeek": "d",
    "startTime": "00:00",
    "endTime": "23:59",
    "usersAvailable": 3,
    "attendees": [
      "user1",
      "user2",
      "user3"
    ],
    "duration": 1440
  },
  {
    "dayOfWeek": "j",
    "startTime": "00:00",
    "endTime": "23:59",
    "usersAvailable": 3,
    "attendees": [
      "user1",
      "user2",
      "user3"
    ],
    "duration": 1440
  },
  {
    "dayOfWeek": "v",
    "startTime": "00:00",
    "endTime": "07:59",
    "usersAvailable": 3,
    "attendees": [
      "user1",
      "user2",
      "user3"
    ],
    "duration": 480
  },
  {
    "dayOfWeek": "v",
    "startTime": "08:00",
    "endTime": "08:59",
    "usersAvailable": 2,
    "attendees": [
      "user1",
      "user3"
    ],
    "duration": 60
  },
  {
    "dayOfWeek": "v",
    "startTime": "09:00",
    "endTime": "09:59",
    "usersAvailable": 1,
    "attendees": [
      "user3"
    ],
    "duration": 60
  },
  {
    "dayOfWeek": "v",
    "startTime": "10:00",
    "endTime": "10:00",
    "usersAvailable": 0,
    "attendees": [],
    "duration": 1
  },
  {
    "dayOfWeek": "v",
    "startTime": "10:01",
    "endTime": "11:00",
    "usersAvailable": 1,
    "attendees": [
      "user2"
    ],
    "duration": 60
  },
  {
    "dayOfWeek": "v",
    "startTime": "11:01",
    "endTime": "11:59",
    "usersAvailable": 2,
    "attendees": [
      "user1",
      "user2"
    ],
    "duration": 59
  },
  {
    "dayOfWeek": "v",
    "startTime": "12:00",
    "endTime": "12:00",
    "usersAvailable": 1,
    "attendees": [
      "user1"
    ],
    "duration": 1
  },
  {
    "dayOfWeek": "v",
    "startTime": "12:01",
    "endTime": "13:59",
    "usersAvailable": 2,
    "attendees": [
      "user1",
      "user3"
    ],
    "duration": 119
  },
  {
    "dayOfWeek": "v",
    "startTime": "14:00",
    "endTime": "14:00",
    "usersAvailable": 1,
    "attendees": [
      "user3"
    ],
    "duration": 1
  },
  {
    "dayOfWeek": "v",
    "startTime": "14:01",
    "endTime": "16:00",
    "usersAvailable": 2,
    "attendees": [
      "user2",
      "user3"
    ],
    "duration": 120
  },
  {
    "dayOfWeek": "v",
    "startTime": "16:01",
    "endTime": "23:59",
    "usersAvailable": 3,
    "attendees": [
      "user1",
      "user2",
      "user3"
    ],
    "duration": 479
  },
  {
    "dayOfWeek": "i",
    "startTime": "00:00",
    "endTime": "23:59",
    "usersAvailable": 3,
    "attendees": [
      "user1",
      "user2",
      "user3"
    ],
    "duration": 1440
  },
  {
    "dayOfWeek": "s",
    "startTime": "00:00",
    "endTime": "23:59",
    "usersAvailable": 3,
    "attendees": [
      "user1",
      "user2",
      "user3"
    ],
    "duration": 1440
  },
  {
    "dayOfWeek": "l",
    "startTime": "00:00",
    "endTime": "23:59",
    "usersAvailable": 3,
    "attendees": [
      "user1",
      "user2",
      "user3"
    ],
    "duration": 1440
  }
]
```

## Troubleshooting

If you encounter any issues while using the Planner service, refer to the following troubleshooting guide to identify and resolve common problems:

### Google's App Engine Bad Request:

You will get the following error if you send any data over the body of a GET request. If you're querying `/hello`, remove the body completely; if you're querying `/planner`, change the method type to POST.

```html
<!DOCTYPE html>
<html lang=en>
  <meta charset=utf-8>
  <meta name=viewport
    content="initial-scale=1, minimum-scale=1, width=device-width">
  <title>Error 400 (Bad Request)!!1</title>
  <style>
    *{margin:0;padding:0}html,code{font:15px/22px arial,sans-serif}html{background:#fff;color:#222;padding:15px}body{margin:7% auto 0;max-width:390px;min-height:180px;padding:30px 0 15px}* > body{background:url(//www.google.com/images/errors/robot.png) 100% 5px no-repeat;padding-right:205px}p{margin:11px 0 22px;overflow:hidden}ins{color:#777;text-decoration:none}a img{border:0}@media screen and (max-width:772px){body{background:none;margin-top:0;max-width:none;padding-right:0}}#logo{background:url(//www.google.com/images/branding/googlelogo/1x/googlelogo_color_150x54dp.png) no-repeat;margin-left:-5px}@media only screen and (min-resolution:192dpi){#logo{background:url(//www.google.com/images/branding/googlelogo/2x/googlelogo_color_150x54dp.png) no-repeat 0% 0%/100% 100%;-moz-border-image:url(//www.google.com/images/branding/googlelogo/2x/googlelogo_color_150x54dp.png) 0}}@media only screen and (-webkit-min-device-pixel-ratio:2){#logo{background:url(//www.google.com/images/branding/googlelogo/2x/googlelogo_color_150x54dp.png) no-repeat;-webkit-background-size:100% 100%}}#logo{display:inline-block;height:54px;width:150px}
  </style>
  <a href=//www.google.com/><span id=logo aria-label=Google></span></a>
  <p><b>400.</b> <ins>That’s an error.</ins>
    <p>Your client has issued a malformed or illegal request. <ins>That’s all we
        know.</ins>
```

### Parameters

You will get the following errors if the parameters of your request are invalid (in particular, those associated with the events for each user). You can check `api/middlewares/planner_middleware.go` for more information. The API will return an error if you don't follow the required format.

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
