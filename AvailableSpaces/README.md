# AvailableSpaces Service

The AvailableSpaces service, a critical component of our application, facilitates the retrieval of information regarding available classrooms within specified time frames and days of the week at Universidad de los Andes' campus. Developed in TypeScript using Express.js and Node.js, this service serves as an API endpoint to support scheduling and resource allocation for university activities. It is deployed on Google Cloud Platform's App Engine.

## Installation Guide

To install and run the AvailableSpaces service locally or in your environment, follow these steps:

1. Clone the Repository: `git clone https://github.com/ISIS3510-202410-Team-13/Backend.git`
2. Navigate to the Backend Directory: `cd Backend/AvailableSpaces`
3. Build the Docker Image: `docker build -t available_spaces .`
4. Run the Docker Container: `docker run -d -p 3000:3000 available_spaces`
5. Verify Installation: `docker ps`

## Usage Guide

### Health Checkpoint

Make a GET request to the following endpoint:

```
GET https://available-spaces-dot-unischedule-5ee93.uc.r.appspot.com/health
```

If the service is running correctly, you will receive the message `Available Spaces Server is running`.

### Request Schema

To query the AvailableSpaces service, make a POST request to the following endpoint:

```
POST https://available-spaces-dot-unischedule-5ee93.uc.r.appspot.com/spaces
```

Include the following JSON body in your request:

```json
{
  "dayOfWeek": "l|m|i|j|v|s|d",
  "startTime": "hhmm",
  "endTime": "hhmm"
}
```


- `dayOfWeek`: A single letter string representing the day of the week. Valid values are 'l' (Monday), 'm' (Tuesday), 'i' (Wednesday), 'j' (Thursday), 'v' (Friday), 's' (Saturday), and 'd' (Sunday).
- `startTime`: A four-digit string representing the start time in 24-hour format (e.g., "0800" for 8:00 AM).
- `endTime`: A four-digit string representing the end time in 24-hour format (e.g., "1700" for 5:00 PM).


### Response Schema

The response will be a JSON array containing objects with the following structure:

```json
[
  {
    "building": "string",
    "room": "string",
    "availableFrom": "hhmm",
    "availableUntil": "hhmm",
    "minutesAvailable": "integer"
  },
]
```


- `building`: The building code where the space is located (e.g., "ML" for Mario Laserna building).
- `room`: The room number or identifier.
- `availableFrom`: The time the space is available from in 24-hour format.
- `availableUntil`: The time the space is available until in 24-hour format.
- `minutesAvailable`: The duration of availability in minutes.

> Notice that the endpoint returns a list of these objects (an empty list is possible)

### Example Call

```bash
curl -X POST https://available-spaces-dot-unischedule-5ee93.uc.r.appspot.com/spaces \
-H "Content-Type: application/json" \
-d '{"dayOfWeek": "m", "startTime": "0800", "endTime": "1200"}'
```

## Troubleshooting

If you encounter any issues while using the AvailableSpaces service, refer to the following troubleshooting guide to identify and resolve common problems:

### Google's App Engine Bad Request:

You will get the following error if you send any data over the body of a GET request. If you're querying `/health`, remove the body completely; if you're querying `/spaces`, change the method type to POST.

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

You will get the following errors if the parameters of your request are invalid. You can check `src/validate.ts` for more information.

1. Request Body Missing:

- Error Message: `{ "message": "Request body is missing" }`
- Cause: The request sent to the service is missing the JSON body.
- Solution: Ensure that your request includes a valid JSON body with the required parameters.

2. Missing Required Fields:

- Error Message: `{ "message": "Missing required fields in request body" }`
- Cause: The request body is missing one or more required fields (dayOfWeek, startTime, endTime).
- Solution: Include all required fields in the request body before sending the request.

3. Invalid Data Types:

- Error Message: `{ "message": "Invalid data types in request body" }`
- Cause: One or more fields in the request body have invalid data types.
- Solution: Ensure that all fields in the request body have the correct data types (dayOfWeek and time fields should be strings).

4. Invalid Day of Week:

- Error Message: `{ "message": "Invalid day of week" }`
- Cause: The dayOfWeek parameter in the request body is not a valid single-letter string representing a day of the week.
- Solution: Check that the dayOfWeek parameter is a single-letter string ('l', 'm', 'i', 'j', 'v', 's', or 'd').

5. Invalid Time Format:

- Error Message: `{ "message": "Invalid time format" }`
- Cause: The startTime or endTime parameter in the request body does not match the expected four-digit string format (hhmm).
- Solution: Ensure that the time parameters follow the four-digit string format (e.g., "0800" for 8:00 AM).

6. Start Time Must be Before End Time:

- Error Message: `{ "message": "Start time must be before end time" }`
- Cause: The startTime parameter is equal to or later than the endTime parameter in the request body.
- Solution: Adjust the startTime and endTime parameters to ensure that the start time comes before the end time.

7. Invalid Time Range:

- Error Message: `{ "message": "Invalid time range" }`
- Cause: The startTime or endTime parameter in the request body is outside the valid time range (00:00 - 23:59).
- Solution: Verify that the time parameters fall within the valid range (00:00 - 23:59).

8. Invalid Minutes:

- Error Message: `{ "message": "Invalid minutes" }`
- Cause: The minutes portion of the startTime or endTime parameter in the request body is greater than 59.
- Solution: Ensure that the minutes portion of the time parameters does not exceed 59.
