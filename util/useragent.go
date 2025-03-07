package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	browserTemplates = `
Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36
Mozilla/5.0 (Macintosh; Intel Mac OS X 12_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.7 Safari/605.1.15
Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36
Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1
`

	osTemplates = `
Win 10
Mac OS X 12.6
Ubuntu 22.04 LTS
iOS 16.0
`

	versionTemplates = `
Chrome 103.0.5060.114
Safari 16.7
Chrome 103.0.5060.114
Mobile/15E148
`

	deviceTemplates = `
Desktop
Tablet
Mobile
`
)

func getRandomTemplate(template string) string {
	tokens := strings.Split(template, "\n")
	// Use crypto/rand for secure random number generation
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(tokens))))
	if err != nil {
		// Fallback to a simple index if crypto/rand fails
		return tokens[0]
	}
	return tokens[n.Int64()]
}

// GetRandomUserAgent returns a random User-Agent string that is a combination of
// a browser type, OS, version and device type. The browser types are all from
// the 91.0.4472 series, the OSes are Windows 10, Mac OS X 10.15.7, Ubuntu 20.04
// LTS, and iOS 14.0, the version numbers are 91.0.4472.124, 16.1.1,
// 91.0.4472.101, and 15E148, and the device types are Desktop, Tablet, and
// Mobile.
func GetRandomUserAgent() string {
	browser := getRandomTemplate(browserTemplates)
	os := getRandomTemplate(osTemplates)
	version := getRandomTemplate(versionTemplates)
	device := getRandomTemplate(deviceTemplates)

	return fmt.Sprintf("%s %s %s %s", browser, os, version, device)
}
