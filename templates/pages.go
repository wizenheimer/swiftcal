// templates/pages.go
package templates

// GetWelcomePageHTML returns the HTML for the welcome page
func GetWelcomePageHTML(emailDomain string) string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>Welcome to swiftcal</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
        .container { max-width: 600px; margin: 0 auto; }
        h1 { color: #2c3e50; }
        p { color: #7f8c8d; line-height: 1.6; }
        .highlight { color: #3498db; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎉 Welcome to swiftcal!</h1>
        <p>We're delighted to have you on board! Your account has been successfully created and is ready to help you manage your calendar more efficiently.</p>
        <p>Simply forward any email containing event details to <span class="highlight">swiftcal@` + emailDomain + `</span>, and our AI will automatically add it to your Google Calendar for you.</p>
        <p>We're excited to see how swiftcal can help streamline your scheduling!</p>
    </div>
</body>
</html>`
}

// Get404PageHTML returns the HTML for the 404 page
func Get404PageHTML(emailDomain string) string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>Page Not Found - swiftcal</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
        .container { max-width: 600px; margin: 0 auto; }
        h1 { color: #e74c3c; }
        p { color: #7f8c8d; line-height: 1.6; }
    </style>
</head>
<body>
    <div class="container">
        <h1>404 - Page Not Found</h1>
        <p>We apologize, but the page you're looking for doesn't exist. If you need assistance, please don't hesitate to reach out to us at <a href="mailto:hey@` + emailDomain + `">hey@` + emailDomain + `</a>. We're here to help!</p>
    </div>
</body>
</html>`
}

// GetInvitedPageHTML returns the HTML for the invited page
func GetInvitedPageHTML() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>Guests Invited - swiftcal</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
        .container { max-width: 600px; margin: 0 auto; }
        h1 { color: #27ae60; }
        p { color: #7f8c8d; line-height: 1.6; }
    </style>
</head>
<body>
    <div class="container">
        <h1>✅ Wonderful! Guests Have Been Invited</h1>
        <p>Perfect! We've successfully sent calendar invitations to all the additional attendees for your event. They should receive their invitations shortly.</p>
        <p>Thank you for using swiftcal to help manage your event!</p>
    </div>
</body>
</html>`
}
