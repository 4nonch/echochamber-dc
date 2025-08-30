package patterns

import (
	"regexp"
)

var (
	MessageLink = regexp.MustCompile(`^https://discord\.com/channels/\d+/\d+/(\d+)$`)
)
