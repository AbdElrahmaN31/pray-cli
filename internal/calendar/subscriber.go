// Package calendar provides calendar generation and ICS file handling
package calendar

import (
	"fmt"
	"io"
)

// Instructions contains subscription instructions for different platforms
type Instructions struct {
	URL string
}

// NewInstructions creates new subscription instructions
func NewInstructions(url string) *Instructions {
	return &Instructions{URL: url}
}

// Print writes the subscription instructions to the writer
func (i *Instructions) Print(w io.Writer) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "📅 Calendar Subscription Instructions")
	fmt.Fprintln(w, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Your calendar URL:")
	fmt.Fprintln(w, i.URL)
	fmt.Fprintln(w)

	// Google Calendar
	fmt.Fprintln(w, "📍 Google Calendar:")
	fmt.Fprintln(w, "   1. Go to calendar.google.com")
	fmt.Fprintln(w, "   2. Click '+' next to 'Other calendars'")
	fmt.Fprintln(w, "   3. Select 'From URL'")
	fmt.Fprintln(w, "   4. Paste the URL above")
	fmt.Fprintln(w, "   5. Click 'Add calendar'")
	fmt.Fprintln(w)

	// Apple Calendar
	fmt.Fprintln(w, "🍎 Apple Calendar (macOS/iOS):")
	fmt.Fprintln(w, "   macOS:")
	fmt.Fprintln(w, "   1. Open Calendar app")
	fmt.Fprintln(w, "   2. File → New Calendar Subscription")
	fmt.Fprintln(w, "   3. Paste the URL above")
	fmt.Fprintln(w, "   4. Click 'Subscribe'")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "   iOS:")
	fmt.Fprintln(w, "   1. Go to Settings → Calendar → Accounts")
	fmt.Fprintln(w, "   2. Add Account → Other")
	fmt.Fprintln(w, "   3. Add Subscribed Calendar")
	fmt.Fprintln(w, "   4. Paste the URL above")
	fmt.Fprintln(w)

	// Outlook
	fmt.Fprintln(w, "📧 Microsoft Outlook:")
	fmt.Fprintln(w, "   Desktop:")
	fmt.Fprintln(w, "   1. Go to Calendar view")
	fmt.Fprintln(w, "   2. File → Account Settings → Account Settings")
	fmt.Fprintln(w, "   3. Internet Calendars → New")
	fmt.Fprintln(w, "   4. Paste the URL above")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "   Web (outlook.com):")
	fmt.Fprintln(w, "   1. Go to Calendar")
	fmt.Fprintln(w, "   2. Add calendar → Subscribe from web")
	fmt.Fprintln(w, "   3. Paste the URL above")
	fmt.Fprintln(w)

	// Thunderbird
	fmt.Fprintln(w, "🦊 Mozilla Thunderbird:")
	fmt.Fprintln(w, "   1. Open Calendar tab")
	fmt.Fprintln(w, "   2. File → New → Calendar")
	fmt.Fprintln(w, "   3. Select 'On the Network'")
	fmt.Fprintln(w, "   4. Format: iCalendar (ICS)")
	fmt.Fprintln(w, "   5. Paste the URL above")
	fmt.Fprintln(w)

	// Note
	fmt.Fprintln(w, "💡 Note:")
	fmt.Fprintln(w, "   - Subscribed calendars auto-update (usually every 24h)")
	fmt.Fprintln(w, "   - To change prayer times, update your location and generate a new URL")
	fmt.Fprintln(w, "   - Events include reminders based on your alarm settings")
	fmt.Fprintln(w)
}

// PrintShort writes abbreviated subscription instructions
func (i *Instructions) PrintShort(w io.Writer) {
	fmt.Fprintln(w)
	fmt.Fprintln(w, "📅 Subscribe to this calendar:")
	fmt.Fprintln(w)
	fmt.Fprintln(w, i.URL)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Run 'pray calendar subscribe' for detailed instructions.")
	fmt.Fprintln(w)
}
