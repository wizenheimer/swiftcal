// templates/prompts.go
package templates

// GetEventExtractionPrompt returns the OpenAI prompt for extracting events from emails
func GetEventExtractionPrompt() string {
	return `
Hello! I'd like to help you extract all the events mentioned in this email thread. Please review the following email content and identify every event that's been discussed.

I'll need you to return the information in this format:
{
  "events": [
    {
      "summary": "the title of the event",
      "location": "a location of the event if one has been given",
      "description": "a description of the event if one has been given",
      "conference_call": true or false, if the event is a conference call or virtual,
      "date": "DD MMMM YYYY - the date of the event",
      "start_time": "HH:mm - the start time of the event in 24 hour format",
      "end_time": "HH:mm - the end time of the event in 24 hour format",
      "attendees": ["list of attendees email addresses"]
    }
  ]
}

Please make sure to capture all distinct events mentioned in the email. Each event with a different date, time, or purpose should be listed separately in the events array.

Here's what to look for:
- The text begins with a Date - that's when the email was sent
- The next line shows the subject of the email thread
- For attendees, consider everyone in the thread, but also think about the email content and subject
- If it's a transactional email (like a receipt or automated message), only include the sender as an attendee
- If it's an email thread, focus on the most recent email but keep the context from the entire thread
- Relative dates like "next tuesday" are perfectly fine - just calculate the actual date based on when the email was sent
- If there aren't enough details for the summary or description, simply use "Event" as a placeholder

To create an event, you'll need at least a date. If you can't find a date for any event, please let me know with this response:
{
  "error": "No date provided",
  "description": "A brief explanation of what information was missing from the email"
}

Here are some examples to help guide you:

---EXAMPLE 1 START---
email_text:
Date: Tue, 26 Mar 2024 12:38:21 +0000
Subject: Fwd: Get ready for the Genius Bar
From: Timmy Jimmy <timmy@gmail.com>
---------- Forwarded message ---------
From: Fitness Clinic <noreply@email.apple.com>
Date: Tue, Mar 26, 2024 at 11:04 AM
Subject: Get ready for the Genius Bar
To: <timmy@gmail.com>

Your upcoming Genius Bar appointment.

Review steps below and check in with a Specialist when you arrive.

For convenience and a quicker check-in, add your appointment to Apple Wallet in your iOS device. Or show this code to a Specialist.

Wednesday, April 3, 2024
10:20

Add to Calendar
Apple Covent Garden

No. 1-7 The Piazza
London

Get directions, view store details, and read store-specific health and safety information
iPhone

Case ID: 102258148113

events_json:
{
  "events": [
    {
      "summary": "Genius Bar",
      "location": "Apple Covent Garden",
      "description": "Case ID: 102258148113",
      "conference_call": false,
      "date": "3 April 2024",
      "start_time": "10:20",
      "end_time": null,
      "attendees": ["timmy@gmail.com"]
    }
  ]
}
--- EXAMPLE 1 END ---

---EXAMPLE 2 START---
email_text:
Date: Thu, 21 Mar 2024 11:38:21 +0000
Subject: Fwd: Investing Holdings Strategic Initiative
From: jeff harry <jeff@investing.com>
---------- Forwarded message ---------
From: Richard Soom <rsoom@toom.com>
Sent: Thursday, March 21, 2024 11:19 AM
To: jeff harry <jeff@investing.com>
Cc: Joe Doe <Joe@investing.com>
Subject: RE: Investing Holdings Strategic Initiative

Thanks jeff, and hello Joe.

May I suggest 3/26 at 3:00 pm ET?

Richard

From: jeff harry <jeff@investing.com>
Sent: Thursday, March 21, 2024 11:18 AM
To: Richard Soom <rsoom@toom.com>
Cc: Joe Doe <Joe@investing.com>
Subject: Investing Holdings Strategic Initiative

Hi Richard,

Updating our previous correspondence, Joe Doe, Investings CIO (cc'd) and I would like a call with at your earliest convenience to discuss:

Introduction to Soom Toom
Investing's progress on sourcing deals to date
potential opportunities to work together
M&A mandate
Fairness opinion
merging with company where Soom Toom is the advisor to the go forward operating company

You had proposed March 28th at 9:30, 10:30 or 11am, do you have any availability prior to that time?

Best regards

events_json:
{
  "events": [
    {
      "summary": "Investing Holdings Strategic Initiative",
      "location": null,
      "description": "Introduction to Soom Toom, Investing's progress on sourcing deals to date, potential opportunities to work together",
      "conference_call": true,
      "date": "26 March 2024",
      "start_time": "15:00",
      "end_time": null,
      "attendees": ["rsoom@toom.com", "jeff@investing.com", "Joe@investing.com"]
    }
  ]
}
--- EXAMPLE 2 END ---

---EXAMPLE 3 START---
email_text:
Date: Fri, 5 Apr 2024 01:08:21 +0000
Subject: find a new suit
From: jeff john <jeff@john.com>
go to h&m next saturday at 2pm

events_json:
{
  "events": [
    {
      "summary": "Find new suit",
      "location": "H&M",
      "description": "find new suit from H&M",
      "conference_call": false,
      "date": "13 April 2024",
      "start_time": "14:00",
      "end_time": null,
      "attendees": ["jeff@john.com"]
    }
  ]
}
--- EXAMPLE 3 END ---

Please respond with JSON only.
`
}

// GetTimezoneExtractionPrompt returns the OpenAI prompt for extracting timezone from emails
func GetTimezoneExtractionPrompt() string {
	return `
Hello! I'd like to help you determine the timezone for the events mentioned in this email thread. Please review the following email content and identify the appropriate IANA Time Zone.

I'll need you to return the information in this format:
{
  "reason": "Brief reasoning of why the timezone was chosen",
  "timezone": "IANA Time Zone Database formatted string"
}

Here's what to look for:
- If a timezone is explicitly mentioned in the email, use that
- If no timezone is given but there's a location, you can infer the timezone based on the location
- Do your best to determine the timezone from the available details
- If no location or timezone information is available, you can leave it as null

Here are some examples to help guide you:

---EXAMPLE 1 START---
email_text:
Date: Tue, 26 Mar 2024 12:38:21 +0000
Subject: Fwd: Get ready for the Genius Bar
From: Timmy Jimmy <timmy@gmail.com>
---------- Forwarded message ---------
From: Apple Support <noreply@email.apple.com>
Date: Tue, Mar 26, 2024 at 11:04 AM
Subject: Get ready for the Genius Bar
To: <timmy@gmail.com>

Your upcoming Genius Bar appointment.

Review steps below and check in with a Specialist when you arrive.

For convenience and a quicker check-in, add your appointment to Apple Wallet in your iOS device. Or show this code to a Specialist.

Wednesday, April 3, 2024
10:20

Add to Calendar
Apple Covent Garden

No. 1-7 The Piazza
London

Get directions, view store details, and read store-specific health and safety information
iPhone

Case ID: 102258148113

timezone_json:
{
  "reason": "The timezone was not explicitly given, but it can be inferred from the location: London",
  "timezone": "Europe/London"
}
--- EXAMPLE 1 END ---

---EXAMPLE 2 START---
email_text:
Date: Thu, 21 Mar 2024 11:38:21 +0000
Subject: Fwd: Investing Holdings Strategic Initiative
From: jeff harry <jeff@investing.com>
---------- Forwarded message ---------
From: Richard Soom <rsoom@toom.com>
Sent: Thursday, March 21, 2024 11:19 AM
To: jeff harry <jeff@investing.com>
Cc: Joe Doe <Joe@investing.com>
Subject: RE: Investing Holdings Strategic Initiative

Thanks jeff, and hello Joe.

May I suggest 3/26 at 3:00 pm ET?

Richard

timezone_json:
{
  "reason": "The timezone ET is referenced, which is eastern time in America",
  "timezone": "America/New_York"
}
---EXAMPLE 2 END---

---EXAMPLE 3 START---
email_text:
Date: Fri, 5 Apr 2024 01:08:21 +0000
Subject: find a new suit
From: jeff john <jeff@john.com>
go to h&m next saturday at 2pm

timezone_json:
{
  "reason": "No timezone or location is given, it is not possible to determine a timezone",
  "timezone": null
}
--- EXAMPLE 3 END ---

Please respond with JSON only.
`
}
