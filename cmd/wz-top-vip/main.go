// Command wz-top-vip calcule un classement pondéré des VIP d'une chaîne Twitch
// à partir des tops uptime et messages de l'API WizeBot.
//
// Utilisation :
//
//	wz-top-vip [options] [fichier-vip]
//
// Le fichier VIP doit contenir un pseudo Twitch par ligne.
// Si absent, le programme cherche vips.txt dans le même dossier que le binaire.
// La clé API WizeBot en lecture est fournie via -apikey ou WIZEBOT_API_READ.
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
		fmt.Fprintf(os.Stderr, "erreur : %v\n\n", err)
		fmt.Fprintf(os.Stderr, "Lancez le programme avec -help pour afficher l'aide.\n")
		return err
	}

	vips, err := vip.LoadFile(cfg.VIPFile)
	if err != nil {
		return err
	}

	client := wizebot.New(cfg.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Fprintf(os.Stderr, "Récupération du top uptime (%s, limite=%d)...\n", cfg.Period, cfg.FetchLimit)
	uptimeEntries, err := client.FetchTop(ctx, wizebot.TopTypeUptime, string(cfg.Period), cfg.FetchLimit)
	if err != nil {
		return fmt.Errorf("récupération du top uptime : %w", err)
	}

	fmt.Fprintf(os.Stderr, "Récupération du top messages (%s, limite=%d)...\n", cfg.Period, cfg.FetchLimit)
	messageEntries, err := client.FetchTop(ctx, wizebot.TopTypeMessage, string(cfg.Period), cfg.FetchLimit)
	if err != nil {
		return fmt.Errorf("récupération du top messages : %w", err)
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

// pauseIfInteractive attend que l'utilisateur appuie sur Entrée avant de quitter,
// uniquement lorsque stdin est un terminal interactif (pas un pipe ni une
// redirection). Cela évite que la fenêtre console se ferme immédiatement
// lors d'un double-clic depuis l'explorateur Windows.
func pauseIfInteractive() {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return
	}
	fmt.Fprint(os.Stderr, "\nAppuyez sur Entrée pour quitter...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}
