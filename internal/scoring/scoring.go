// Package scoring merges uptime and message rankings, filters VIPs,
// normalises values and computes a weighted score for each viewer.
package scoring

import (
	"sort"
	"strings"

	"github.com/vignemail1/wz-top-vip/internal/wizebot"
)

// ViewerScore holds raw, normalised and computed values for one viewer.
type ViewerScore struct {
	Name         string
	UptimeRaw    float64
	MessagesRaw  float64
	UptimeNorm   float64
	MessagesNorm float64
	Score        float64
	IsVIP        bool
}

// Params holds the weights and normalization maxima used during computation.
type Params struct {
	UptimeWeightPct  int
	MessageWeightPct int
	MaxUptimeRaw     float64
	MaxMessagesRaw   float64
}

// Result is the full output of the scoring computation.
type Result struct {
	Params  Params
	Viewers []ViewerScore // all VIP viewers found, sorted by Score desc
}

// Compute merges the two API ranking lists, keeps only VIP entries,
// normalises each metric against its observed maximum, and calculates
// the weighted score.
func Compute(
	uptimeEntries []wizebot.RankingEntry,
	messageEntries []wizebot.RankingEntry,
	vips map[string]struct{},
	uptimeWeightPct int,
	messageWeightPct int,
) Result {
	type raw struct {
		uptime  float64
		message float64
	}

	merged := make(map[string]*raw)

	for _, e := range uptimeEntries {
		key := strings.ToLower(strings.TrimSpace(e.UserName))
		if _, ok := vips[key]; !ok {
			continue
		}
		if _, exists := merged[key]; !exists {
			merged[key] = &raw{}
		}
		merged[key].uptime = e.Value
	}

	for _, e := range messageEntries {
		key := strings.ToLower(strings.TrimSpace(e.UserName))
		if _, ok := vips[key]; !ok {
			continue
		}
		if _, exists := merged[key]; !exists {
			merged[key] = &raw{}
		}
		merged[key].message = e.Value
	}

	// Ensure every VIP appears, even those absent from both tops (score = 0).
	for name := range vips {
		if _, exists := merged[name]; !exists {
			merged[name] = &raw{}
		}
	}

	var maxUptime, maxMessages float64
	for _, r := range merged {
		if r.uptime > maxUptime {
			maxUptime = r.uptime
		}
		if r.message > maxMessages {
			maxMessages = r.message
		}
	}

	uw := float64(uptimeWeightPct) / 100.0
	mw := float64(messageWeightPct) / 100.0

	viewers := make([]ViewerScore, 0, len(merged))
	for name, r := range merged {
		vs := ViewerScore{
			Name:        name,
			UptimeRaw:   r.uptime,
			MessagesRaw: r.message,
			IsVIP:       true,
		}
		if maxUptime > 0 {
			vs.UptimeNorm = r.uptime / maxUptime
		}
		if maxMessages > 0 {
			vs.MessagesNorm = r.message / maxMessages
		}
		vs.Score = vs.UptimeNorm*uw + vs.MessagesNorm*mw
		viewers = append(viewers, vs)
	}

	sort.Slice(viewers, func(i, j int) bool {
		if viewers[i].Score != viewers[j].Score {
			return viewers[i].Score > viewers[j].Score
		}
		return viewers[i].Name < viewers[j].Name
	})

	return Result{
		Params: Params{
			UptimeWeightPct:  uptimeWeightPct,
			MessageWeightPct: messageWeightPct,
			MaxUptimeRaw:     maxUptime,
			MaxMessagesRaw:   maxMessages,
		},
		Viewers: viewers,
	}
}
