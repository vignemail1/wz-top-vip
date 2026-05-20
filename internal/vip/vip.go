// Package vip loads and normalises the list of VIP usernames from a text file.
package vip

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadFile lit une liste de pseudos Twitch (un par ligne) depuis le fichier path.
// Les lignes vides et celles commençant par '#' sont ignorées.
// Tous les pseudos sont mis en minuscules pour une comparaison uniforme.
func LoadFile(path string) (map[string]struct{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ouverture du fichier VIP %q : %w", path, err)
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
		return nil, fmt.Errorf("lecture du fichier VIP %q : %w", path, err)
	}
	if len(vips) == 0 {
		return nil, fmt.Errorf("le fichier VIP %q ne contient aucun pseudo valide", path)
	}
	return vips, nil
}
