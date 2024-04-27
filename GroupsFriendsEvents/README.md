# Groups-Friends-Events

Welcome to the Unischedule Backend API documentation. This API provides endpoints to interact with user data, group data, and event data for the Unischedule application. Below you will find information on how to install, build, and use the API, as well as troubleshooting tips.

## Installation Guide

To install and run the Unischedule Backend API locally or in your environment, follow these steps:

1. Clone the Repository: `git clone https://github.com/ISIS3510-202410-Team-13/Backend.git`
2. Navigate to the directory: `cd GroupsFriendsEvents`
3. Install Dependencies: `npm install`
4. Set up Firebase: Follow the instructions in the Firebase documentation to set up Firebase Admin SDK and obtain your service account JSON file.
5. Run the Server: `npm start`
6. Verify Installation: Open a web browser and navigate to `http://localhost:3000` to ensure the server is running correctly.

## Usage Guide

### User Endpoints

#### Get User Groups

- URL: `/user/:id/groups`
- Method: GET
- Description: Retrieves the groups that a user is a member of.
- Parameters:
  - `id`: The ID of the user.
- Response:
  - Status Code: 200 OK
  - Content: JSON array of group objects containing group information.

#### Get User Friends

- URL: `/user/:id/friends`
- Method: GET
- Description: Retrieves the friends of a user.
- Parameters:
  - `id`: The ID of the user.
- Response:
  - Status Code: 200 OK
  - Content: JSON array of user objects containing friend information.

#### Get User Events

- URL: `/user/:id/events`
- Method: GET
- Description: Retrieves the events associated with a user.
- Parameters:
  - `id`: The ID of the user.
- Response:
  - Status Code: 200 OK
  - Content: JSON array of event objects containing event information.

#### Create User Event

- URL: `/user/:id/events`
- Method: POST
- Description: Adds a new event for a user.
- Parameters:
  - `id`: The ID of the user.
- Request Body:
  - JSON object containing event data.
- Response:
  - Status Code: 201 Created
  - Content: JSON object with message and event ID.

### Group Endpoints

#### Get Group Members

- URL: `/group/:id/members`
- Method: GET
- Description: Retrieves the members of a group.
- Parameters:
  - `id`: The ID of the group.
- Response:
  - Status Code: 200 OK
  - Content: JSON array of user objects containing member information.

#### Get Group Events

- URL: `/group/:id/events`
- Method: GET
- Description: Retrieves the events associated with a group.
- Parameters:
  - `id`: The ID of the group.
- Response:
  - Status Code: 200 OK
  - Content: JSON array of event objects containing event information.

## Troubleshooting

If you encounter any issues while using the Unischedule Groups-Friends-Events API, refer to the following troubleshooting guide to identify and resolve common problems:

### User Not Found

- Error Message: "User not found"
- Cause: The specified user ID does not exist in the database.
- Solution: Verify that the user ID is correct and exists in the database.

### Group Not Found

- Error Message: "Group not found"
- Cause: The specified group ID does not exist in the database.
- Solution: Verify that the group ID is correct and exists in the database.

### Firebase Setup

- Error Message: "Firebase setup error"
- Cause: Firebase Admin SDK is not properly set up or initialized.
- Solution: Follow the Firebase documentation to ensure proper setup and initialization of the Firebase Admin SDK.

### Invalid Request Parameters

- Error Message: "Invalid request parameters"
- Cause: The parameters or request body format are invalid.
- Solution: Check the API documentation for correct parameter and request body formats.

### Server Error

- Error Message: "Internal server error"
- Cause: An unexpected error occurred on the server.
- Solution: Check the server logs for more information and try again.
