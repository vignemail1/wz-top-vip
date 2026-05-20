// Command wz-top-vip computes a weighted ranking of VIP viewers on a Twitch
// channel using WizeBot's uptime and message ranking APIs.
//
// Usage:
//
//	wz-top-vip [flags] <vip-file>
//
// The VIP file must contain one Twitch username per line.
// The WizeBot read API key is provided via -apikey or the WIZEBOT_API_READ
// environment variable.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vignemail1/wz-top-vip/internal/config"
	"github.com/vignemail1/wz-top-vip/internal/output"
	"github.com/vignemail1/wz-top-vip/internal/scoring"
	"github.com/vignemail1/wz-top-vip/internal/vip"
	"github.com/vignemail1/wz-top-vip/internal/wizebot"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "Run with -help for usage.\n")
		os.Exit(2)
	}

	vips, err := vip.LoadFile(cfg.VIPFile)
	if err != nil {
		return err
	}

	client := wizebot.New(cfg.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Fprintf(os.Stderr, "Récupération du top uptime (%s, limit=%d)...\n", cfg.Period, cfg.FetchLimit)
	uptimeEntries, err := client.FetchTop(ctx, wizebot.TopTypeUptime, string(cfg.Period), cfg.FetchLimit)
	if err != nil {
		return fmt.Errorf("fetching uptime top: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Récupération du top messages (%s, limit=%d)...\n", cfg.Period, cfg.FetchLimit)
	messageEntries, err := client.FetchTop(ctx, wizebot.TopTypeMessage, string(cfg.Period), cfg.FetchLimit)
	if err != nil {
		return fmt.Errorf("fetching message top: %w", err)
	}

	result := scoring.Compute(
		uptimeEntries,
		messageEntries,
		vips,
		cfg.UptimeWeightPct,
		cfg.MessageWeightPct,
	)

	output.Print(os.Stdout, string(cfg.Period), result, cfg.TopN)
	return nil
}
