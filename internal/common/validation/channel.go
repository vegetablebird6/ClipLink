package validation

import "regexp"

var channelIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{8,128}$`)

// IsValidChannelID keeps channel IDs usable in URLs/headers while preventing
// weak, ambiguous, or unexpectedly large bearer tokens.
func IsValidChannelID(channelID string) bool {
	return channelIDPattern.MatchString(channelID)
}
