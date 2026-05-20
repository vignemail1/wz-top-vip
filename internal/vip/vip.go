// Package vip loads and normalises the list of VIP usernames from a text file.
package vip

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadFile reads a newline-separated list of Twitch usernames from path.
// Empty lines and lines starting with '#' are ignored.
// All usernames are lowercased and trimmed for consistent comparison.
func LoadFile(path string) (map[string]struct{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening VIP file %q: %w", path, err)
	}
	defer f.Close()

	vips := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		vips[strings.ToLower(line)] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading VIP file %q: %w", path, err)
	}
	if len(vips) == 0 {
		return nil, fmt.Errorf("VIP file %q contains no valid usernames", path)
	}
	return vips, nil
}
