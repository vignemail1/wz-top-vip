// Package config handles CLI flag parsing, environment variable resolution
// and validation of all program parameters.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"
)

const (
	defaultPeriod        = "month"
	defaultMessageWeight = 50
	defaultTop           = 3
	defaultFetchLimit    = 100
	maxFetchLimit        = 100

	defaultVIPFile = "vips.txt"
	apiKeyURL      = "https://panel.wizebot.tv/development_api_management#"
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
// prompts interactively if still missing, and validates all parameters.
func Parse() (*Config, error) {
	var (
		apiKey        string
		period        string
		messageWeight int
		topN          int
		fetchLimit    int
	)

	flag.StringVar(&apiKey, "apikey", "", "Clé API WizeBot en lecture (remplace la variable WIZEBOT_API_READ)")
	flag.StringVar(&period, "period", defaultPeriod, "Période : week (semaine) ou month (mois)")
	flag.IntVar(&messageWeight, "message-weight", defaultMessageWeight, "Poids des messages dans le score (0-100, en pourcentage)")
	flag.IntVar(&topN, "top", defaultTop, "Nombre de VIP à afficher")
	flag.IntVar(&fetchLimit, "fetch-limit", defaultFetchLimit, "Nombre d'entrées à récupérer par classement WizeBot (1-100)")

	flag.Usage = usage
	flag.Parse()

	// Ordre de résolution : flag > variable d'environnement > saisie interactive.
	if apiKey == "" {
		apiKey = os.Getenv("WIZEBOT_API_READ")
	}
	if apiKey == "" {
		var err error
		apiKey, err = promptAPIKey()
		if err != nil {
			return nil, fmt.Errorf("lecture de la clé API : %w", err)
		}
	}

	// Résolution du fichier VIP :
	// 1. Argument positionnel fourni en ligne de commande.
	// 2. vips.txt dans le même dossier que le binaire (silencieux).
	// 3. Erreur explicite.
	var vipFile string
	if flag.NArg() >= 1 {
		vipFile = flag.Arg(0)
	} else {
		vipFile = vipFileNextToExecutable()
		if vipFile == "" {
			return nil, fmt.Errorf(
				"fichier VIP manquant : passez-le en argument ou placez un fichier %q à côté du programme",
				defaultVIPFile,
			)
		}
	}

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

// vipFileNextToExecutable returns the path to vips.txt located in the same
// directory as the running executable, or an empty string if not found.
func vipFileNextToExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	candidate := filepath.Join(filepath.Dir(exePath), defaultVIPFile)
	if _, err := os.Stat(candidate); err == nil {
		return candidate
	}
	return ""
}

// promptAPIKey displays instructions and reads the API key with masked echo.
func promptAPIKey() (string, error) {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Aucune clé API WizeBot trouvée (-apikey / WIZEBOT_API_READ).")
	fmt.Fprintln(os.Stderr, "Vous pouvez obtenir votre clé API [R] (lecture) ici :")
	fmt.Fprintln(os.Stderr, " "+apiKeyURL)
	fmt.Fprintln(os.Stderr)
	fmt.Fprint(os.Stderr, "Clé API WizeBot [R] : ")

	raw, err := term.ReadPassword(int(os.Stderr.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}

	key := strings.TrimSpace(string(raw))
	if key == "" {
		return "", errors.New("la clé API ne peut pas être vide")
	}

	mask := strings.Repeat("*", len(key))
	fmt.Fprintf(os.Stderr, "Clé saisie : %s (%d caractères)\n", mask, len(key))
	return key, nil
}

func validate(c *Config) (*Config, error) {
	if c.APIKey == "" {
		return nil, errors.New("clé API requise : utilisez -apikey ou définissez WIZEBOT_API_READ")
	}
	if !c.Period.Valid() {
		return nil, fmt.Errorf("période invalide %q : doit être \"week\" ou \"month\"", c.Period)
	}
	if c.MessageWeightPct < 0 || c.MessageWeightPct > 100 {
		return nil, fmt.Errorf("-message-weight %d invalide : doit être entre 0 et 100", c.MessageWeightPct)
	}
	if c.TopN < 1 {
		return nil, fmt.Errorf("-top %d invalide : doit être >= 1", c.TopN)
	}
	if c.FetchLimit < 1 || c.FetchLimit > maxFetchLimit {
		return nil, fmt.Errorf("-fetch-limit %d invalide : doit être entre 1 et %d", c.FetchLimit, maxFetchLimit)
	}
	return c, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `Utilisation : wz-top-vip [options] [fichier-vip]

  [fichier-vip]  Chemin vers un fichier texte contenant un pseudo Twitch par ligne
                 (liste des VIPs). Optionnel : si absent, le programme cherche
                 automatiquement un fichier vips.txt à côté du programme.

Options :
`)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
Variables d'environnement :
  WIZEBOT_API_READ  Clé API WizeBot en lecture (utilisée si -apikey est absent)

Obtenir la clé API :
  %s

Exemples :
  wz-top-vip vips.txt
  wz-top-vip -apikey xxxxx -period week -message-weight 70 -top 5 vips.txt
`, apiKeyURL)
}
