// Package config handles CLI flag parsing, environment variable resolution
// and validation of all program parameters.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const (
	defaultPeriod        = "month"
	defaultMessageWeight = 50
	defaultTop           = 3
	defaultFetchLimit    = 100
	maxFetchLimit        = 100
)

// Period represents a WizeBot ranking time window.
type Period string

const (
	PeriodWeek  Period = "week"
	PeriodMonth Period = "month"
)

// Valid reports whether p is a supported Period value.
func (p Period) Valid() bool {
	return p == PeriodWeek || p == PeriodMonth
}

// Config holds all resolved program parameters.
type Config struct {
	APIKey           string
	VIPFile          string
	Period           Period
	MessageWeightPct int
	UptimeWeightPct  int
	TopN             int
	FetchLimit       int
}

// MessageWeight returns the message weight as a float64 coefficient.
func (c *Config) MessageWeight() float64 {
	return float64(c.MessageWeightPct) / 100.0
}

// UptimeWeight returns the uptime weight as a float64 coefficient.
func (c *Config) UptimeWeight() float64 {
	return float64(c.UptimeWeightPct) / 100.0
}

// Parse parses os.Args, resolves the API key from flag or environment,
// and validates all parameters. It returns a ready-to-use Config or an error.
func Parse() (*Config, error) {
	var (
		apiKey        string
		period        string
		messageWeight int
		topN          int
		fetchLimit    int
	)

	flag.StringVar(&apiKey, "apikey", "", "WizeBot read API key (overrides WIZEBOT_API_READ env var)")
	flag.StringVar(&period, "period", defaultPeriod, "Time period: week or month")
	flag.IntVar(&messageWeight, "message-weight", defaultMessageWeight, "Weight of messages in score (0-100, as percentage)")
	flag.IntVar(&topN, "top", defaultTop, "Number of top VIPs to display")
	flag.IntVar(&fetchLimit, "fetch-limit", defaultFetchLimit, "Number of entries to fetch from WizeBot API per ranking (1-100)")

	flag.Usage = usage
	flag.Parse()

	if apiKey == "" {
		apiKey = os.Getenv("WIZEBOT_API_READ")
	}

	if flag.NArg() < 1 {
		return nil, errors.New("missing required argument: path to VIP file")
	}
	vipFile := flag.Arg(0)

	return validate(&Config{
		APIKey:           apiKey,
		VIPFile:          vipFile,
		Period:           Period(period),
		MessageWeightPct: messageWeight,
		UptimeWeightPct:  100 - messageWeight,
		TopN:             topN,
		FetchLimit:       fetchLimit,
	})
}

func validate(c *Config) (*Config, error) {
	if c.APIKey == "" {
		return nil, errors.New("API key is required: use -apikey flag or set WIZEBOT_API_READ environment variable")
	}
	if !c.Period.Valid() {
		return nil, fmt.Errorf("invalid period %q: must be \"week\" or \"month\"", c.Period)
	}
	if c.MessageWeightPct < 0 || c.MessageWeightPct > 100 {
		return nil, fmt.Errorf("invalid -message-weight %d: must be between 0 and 100", c.MessageWeightPct)
	}
	if c.TopN < 1 {
		return nil, fmt.Errorf("invalid -top %d: must be >= 1", c.TopN)
	}
	if c.FetchLimit < 1 || c.FetchLimit > maxFetchLimit {
		return nil, fmt.Errorf("invalid -fetch-limit %d: must be between 1 and %d", c.FetchLimit, maxFetchLimit)
	}
	return c, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: wz-top-vip [flags] <vip-file>

  <vip-file>  Path to a text file containing one Twitch username per line (VIP list).

Flags:
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
Environment variables:
  WIZEBOT_API_READ  WizeBot read API key (used if -apikey is not set)

Examples:
  wz-top-vip vips.txt
  wz-top-vip -apikey xxxxx -period week -message-weight 70 -top 5 vips.txt
`)
}
