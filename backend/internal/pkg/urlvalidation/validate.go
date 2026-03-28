package urlvalidation

import (
	"fmt"
	"net"
	"net/url"
)

var privateRanges []*net.IPNet

func init() {
	cidrs := []string{
		"127.0.0.0/8",    // loopback
		"10.0.0.0/8",     // RFC 1918
		"172.16.0.0/12",  // RFC 1918
		"192.168.0.0/16", // RFC 1918
		"169.254.0.0/16", // link-local
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 unique local
		"fe80::/10",      // IPv6 link-local
	}
	for _, cidr := range cidrs {
		_, network, _ := net.ParseCIDR(cidr)
		privateRanges = append(privateRanges, network)
	}
}

func isPrivateIP(ip net.IP) bool {
	for _, r := range privateRanges {
		if r.Contains(ip) {
			return true
		}
	}
	return false
}

// ValidateWebhookURL checks that a URL is safe for outgoing webhook delivery.
// It rejects non-HTTP(S) schemes and URLs that resolve to private/reserved IPs.
func ValidateWebhookURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("unsupported scheme %q: only http and https are allowed", parsed.Scheme)
	}

	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("URL must have a hostname")
	}

	// Check if host is a literal IP
	if ip := net.ParseIP(host); ip != nil {
		if isPrivateIP(ip) {
			return fmt.Errorf("webhook URL must not target private/reserved IP address")
		}
		return nil
	}

	// Resolve hostname and check all IPs
	ips, err := net.LookupHost(host)
	if err != nil {
		return fmt.Errorf("failed to resolve hostname %q: %w", host, err)
	}

	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ip != nil && isPrivateIP(ip) {
			return fmt.Errorf("webhook URL must not resolve to private/reserved IP address")
		}
	}

	return nil
}
