// Package output renders scoring results as human-readable console text.
package output

import (
	"fmt"
	"io"
	"math"

	"github.com/vignemail1/wz-top-vip/internal/scoring"
)

// Print writes a full explanation of the computation and the top-N ranking
// to w. topN caps the number of displayed results.
func Print(w io.Writer, period string, result scoring.Result, topN int) {
	p := result.Params

	fmt.Fprintln(w, "════════════════════════════════════════════")
	fmt.Fprintln(w, "        wz-top-vip — Classement VIP")
	fmt.Fprintln(w, "════════════════════════════════════════════")
	fmt.Fprintf(w, "Période       : %s\n", period)
	fmt.Fprintf(w, "Poids uptime  : %d %%\n", p.UptimeWeightPct)
	fmt.Fprintf(w, "Poids messages: %d %%\n", p.MessageWeightPct)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "── Méthode de calcul ───────────────────────")
	fmt.Fprintln(w, "  Chaque métrique est normalisée sur le maximum")
	fmt.Fprintln(w, "  observé dans les résultats de l'API :")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  uptime_norm   = uptime_brut   / %.0f  (max uptime)\n", p.MaxUptimeRaw)
	fmt.Fprintf(w, "  messages_norm = messages_brut / %.0f  (max messages)\n", p.MaxMessagesRaw)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  score = uptime_norm × %d %% + messages_norm × %d %%\n",
		p.UptimeWeightPct, p.MessageWeightPct)
	fmt.Fprintln(w, "────────────────────────────────────────────")
	fmt.Fprintln(w)

	total := len(result.Viewers)
	display := topN
	if total == 0 {
		fmt.Fprintln(w, "⚠  Aucun VIP trouvé dans les tops WizeBot pour cette période.")
		return
	}
	if display > total {
		display = total
	}

	fmt.Fprintf(w, "Top %d VIP (sur %d VIP présents dans les tops)\n", display, total)
	fmt.Fprintln(w, "────────────────────────────────────────────")

	for i := 0; i < display; i++ {
		v := result.Viewers[i]
		fmt.Fprintf(w, "%2d. %-25s  score=%s  uptime=%7.0f (norm=%s)  messages=%5.0f (norm=%s)\n",
			i+1,
			v.Name,
			pct(v.Score),
			v.UptimeRaw,
			pct(v.UptimeNorm),
			v.MessagesRaw,
			pct(v.MessagesNorm),
		)
	}

	fmt.Fprintln(w, "════════════════════════════════════════════")
}

// pct formats a normalised value (0..1) as a human-readable percentage string.
func pct(v float64) string {
	return fmt.Sprintf("%5.1f%%", math.Round(v*1000)/10)
}
