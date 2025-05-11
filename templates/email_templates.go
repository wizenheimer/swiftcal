// templates/email_templates.go
package templates

import (
	"fmt"
	"strings"
)

type EmailTemplate struct {
	HTML    string
	Subject string
}

func GetNoUserFoundTemplate(fromEmail, baseDomain, appDomain, emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`Welcome to swiftcal!<br><br>
We're excited to help you manage your calendar more efficiently. To get started, please click the link below to sign up.<br><br>
<a href="https://www.%s/signup-consent"><img src="https://%s/signup-with-google.png" alt="Sign Up with Google" width="182" height="42" style="display: block;"></a><br><br>
If you forwarded an email to have it added to your calendar, you'll need to forward it again after completing your signup.<br><br>

Note: If you've already signed up for swiftcal and would like to forward events from this email address, please send an email from your main Gmail account to <a href="mailto:swiftcal@%s?subject=add %s">swiftcal@%s</a> with the subject "add %s".<br>
`, baseDomain, appDomain, emailDomain, fromEmail, emailDomain, fromEmail)

	return EmailTemplate{HTML: html}
}

func GetUnverifiedEmailTemplate(fromEmail, emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`We're having trouble verifying your email configuration for %s. This might be due to email settings that prevent us from confirming your identity.

<br><br>If you need assistance, please don't hesitate to reach out: <a href="mailto:hey@%s">hey@%s</a>
`, fromEmail, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}

func GetOAuthFailedTemplate(appDomain, emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`We encountered an issue while connecting to your Google account. Please <a href="https://%s/signup">click here to authorize Google again</a>, and then forward your email thread once more.

Please ensure you complete the checkbox to allow swiftcal to access your calendar.
<img src="https://%s/swiftcalPermissions.png" alt="Google Permissions" width="394" height="170" style="display: block;">

<br><br>If you need any assistance, we're here to help: <a href="mailto:hey@%s">hey@%s</a><br>`, appDomain, appDomain, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}

func GetUnableToParseTemplate(emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`We weren't able to identify a date in your email. Please forward the thread again and include some additional context to help us understand the event details better.

<br><br>If you need assistance, please don't hesitate to reach out: <a href="mailto:hey@%s">hey@%s</a><br>`, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}

func GetAIParseErrorTemplate(description, emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`We weren't able to identify a date in your email. %s
Please forward the thread again and include some additional context to help us understand the event details better.

<br><br>If you need assistance, please don't hesitate to reach out: <a href="mailto:hey@%s">hey@%s</a><br>`, description, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}

func GetEventAddedTemplate(eventLink, eventDate, eventAttendees, emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`Great news! Your event has been successfully added to your calendar.
<br>Date: %s
<br>Attendees: %s
<br><a href="%s" style="display:inline-block; padding:10px 20px; margin:5px 0; background-color:#3498db; color:white; text-align:center; text-decoration:none; font-weight:bold; border-radius:5px; border:none; cursor:pointer;">View Event</a>

<br><br>If you need any assistance, we're here to help: <a href="mailto:hey@%s">hey@%s</a><br>`, eventDate, eventAttendees, eventLink, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}

func GetEventAddedAttendeesTemplate(eventLink, eventDate, inviteLink, eventAttendees, emailDomain string) EmailTemplate {
	attendeesList := strings.ReplaceAll(eventAttendees, ",", "<br>-")

	html := fmt.Sprintf(`Great news! Your event has been successfully added to your calendar.
<br>Date: %s
<br> <a href="%s" style="display:inline-block; padding:10px 20px; margin:5px 0; background-color:#3498db; color:white; text-align:center; text-decoration:none; font-weight:bold; border-radius:5px; border:none; cursor:pointer;">View Event</a>
<br> You may want to invite these attendees:
<br>- %s
<br><a href="%s" style="display:inline-block; padding:10px 20px; margin:5px 0; background-color:#3498db; color:white; text-align:center; text-decoration:none; font-weight:bold; border-radius:5px; border:none; cursor:pointer;">Invite Guests</a>

<br><br>If you need any assistance, we're here to help: <a href="mailto:hey@%s">hey@%s</a><br>`, eventDate, eventLink, attendeesList, inviteLink, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}

func GetICSEventTemplate(eventLink, emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`Perfect! We found an ICS file in your forwarded email and have successfully added this event to your calendar:
<br><a href="%s" style="display:inline-block; padding:10px 20px; margin:5px 0; background-color:#3498db; color:white; text-align:center; text-decoration:none; font-weight:bold; border-radius:5px; border:none; cursor:pointer;">View Event</a>

<br><br>If you need any assistance, we're here to help: <a href="mailto:hey@%s">hey@%s</a><br>`, eventLink, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}

func GetICSErrorTemplate(emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`We found an ICS file in your forwarded email, but encountered an issue while processing it. Please try forwarding the email again.

<br><br>If you need assistance, please don't hesitate to reach out: <a href="mailto:hey@%s">hey@%s</a><br>`, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}

func GetAddAdditionalEmailTemplate(verificationCode, originatorEmail, appDomain, emailDomain string) EmailTemplate {
	subject := fmt.Sprintf("%s would like you to add events to their calendar", originatorEmail)
	html := fmt.Sprintf(`%s has invited you to join their swiftcal account. To accept this invitation, please <a href="https://%s/auth/verifyAdditionalEmail?uuid=%s">click here</a>.<br><br>

Once you approve, you'll be able to forward any email to swiftcal@%s, and it will automatically be converted into an event in %s's calendar using our AI technology.<br><br>

If you prefer not to accept this invitation, you can simply ignore this email.

<br><br>If you need any assistance, we're here to help: <a href="mailto:hey@%s">hey@%s</a>
`, originatorEmail, appDomain, verificationCode, emailDomain, originatorEmail, emailDomain, emailDomain)

	return EmailTemplate{HTML: html, Subject: subject}
}

func GetAdditionalEmailInUseTemplate(emailToAdd, emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`We noticed that %s is already associated with another swiftcal account. If you'd like to add it to this account, the current account holder will need to send an email to swiftcal@%s with the subject "remove %s"

<br><br>If you need assistance, please don't hesitate to reach out: <a href="mailto:hey@%s">hey@%s</a><br>`, emailToAdd, emailDomain, emailToAdd, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}

func GetRemovalEmailInUseTemplate(emailToRemove, emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`We noticed that %s is associated with another swiftcal account, not your current account.

<br><br>If you need assistance, please don't hesitate to reach out: <a href="mailto:hey@%s">hey@%s</a><br>`, emailToRemove, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}

func GetEmailAddressRemovedTemplate(emailToRemove, emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`%s has been successfully removed from your account.

<br><br>If you need any assistance, we're here to help: <a href="mailto:hey@%s">hey@%s</a><br>`, emailToRemove, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}

func GetUserDeletedTemplate(emailDomain string) EmailTemplate {
	html := fmt.Sprintf(`Thank you for using swiftcal! Your account has been successfully deleted. You're always welcome to sign up again by sending an email to <a href="mailto:swiftcal@%s?subject=signup">swiftcal@%s</a>.<br><br>

You may want to disconnect swiftcal from your Google account <a href="https://myaccount.google.com/connections">here</a>.

<br><br>If you need any assistance, we're here to help: <a href="mailto:hey@%s">hey@%s</a><br>`, emailDomain, emailDomain, emailDomain, emailDomain)

	return EmailTemplate{HTML: html}
}
