// Command wz-top-vip computes a weighted ranking of VIP viewers on a Twitch
// channel using WizeBot's uptime and message ranking APIs.
//
// Usage:
//
//	wz-top-vip [flags] [vip-file]
//
// The VIP file must contain one Twitch username per line.
// If omitted, the program looks for vips.txt next to the executable.
// The WizeBot read API key is provided via -apikey or the WIZEBOT_API_READ
// environment variable.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/vignemail1/wz-top-vip/internal/config"
	"github.com/vignemail1/wz-top-vip/internal/output"
	"github.com/vignemail1/wz-top-vip/internal/scoring"
	"github.com/vignemail1/wz-top-vip/internal/vip"
	"github.com/vignemail1/wz-top-vip/internal/wizebot"
)

func main() {
	err := run()
	pauseIfInteractive()
	if err != nil {
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "Run with -help for usage.\n")
		return err
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

// pauseIfInteractive waits for the user to press Enter before exiting, but
// only when stdin is an interactive terminal (i.e. not a pipe or redirect).
// This prevents the console window from closing immediately when the program
// is launched by double-clicking on Windows.
func pauseIfInteractive() {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return
	}
	fmt.Fprint(os.Stderr, "\nAppuyez sur Entrée pour quitter...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}
