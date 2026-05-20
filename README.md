# wz-top-vip

Classement pondéré des VIP d'une chaîne Twitch à partir des tops WizeBot
(temps de visionnage et messages).

## Installation

```bash
make build
```

Cross-compilation pour les trois plateformes :

```bash
make build-all
# → dist/wz-top-vip-linux-amd64
# → dist/wz-top-vip-windows-amd64.exe
# → dist/wz-top-vip-darwin-arm64
```

## Utilisation

```
wz-top-vip [flags] <vip-file>
```

| Flag              | Défaut  | Description |
|-------------------|---------|-------------|
| `-apikey`         | —       | Clé API WizeBot READ (sinon `WIZEBOT_API_READ`) |
| `-period`         | `month` | Période : `week` ou `month` |
| `-message-weight` | `50`    | Poids des messages dans le score (0-100, en %) |
| `-top`            | `3`     | Nombre de VIP à afficher |
| `-fetch-limit`    | `100`   | Entrées récupérées par top API (max 100) |

**Le fichier VIP** doit contenir un pseudo Twitch par ligne.  
Les lignes vides et les lignes commençant par `#` sont ignorées.

### Exemples

```bash
# Top 3 du mois, poids 50/50
WIZEBOT_API_READ=xxxxx wz-top-vip vips.txt

# Top 5 de la semaine, 70% messages
wz-top-vip -apikey xxxxx -period week -message-weight 70 -top 5 vips.txt

# 100% uptime
wz-top-vip -apikey xxxxx -message-weight 0 vips.txt
```

### Exemple de sortie

```
════════════════════════════════════════════
        wz-top-vip — Classement VIP
════════════════════════════════════════════
Période       : month
Poids uptime  : 30 %
Poids messages: 70 %

── Méthode de calcul ───────────────────────
  Chaque métrique est normalisée sur le maximum
  observé dans les résultats de l'API :

  uptime_norm   = uptime_brut   / 45230  (max uptime)
  messages_norm = messages_brut / 1240   (max messages)

  score = uptime_norm × 30 % + messages_norm × 70 %
────────────────────────────────────────────

Top 3 VIP (sur 12 VIP présents dans les tops)
────────────────────────────────────────────
 1. pseudo_a                    score= 84.2%  uptime=  41200 (norm= 91.1%)  messages= 1100 (norm= 88.7%)
 2. pseudo_b                    score= 67.5%  uptime=  38000 (norm= 84.1%)  messages=  720 (norm= 58.1%)
 3. pseudo_c                    score= 55.3%  uptime=  22400 (norm= 49.5%)  messages=  900 (norm= 72.6%)
════════════════════════════════════════════
```

## Méthode de scoring

Pour chaque VIP présent dans au moins un des deux tops WizeBot :

1. `uptime_norm   = uptime_brut   / max(uptime_brut)`
2. `messages_norm = messages_brut / max(messages_brut)`
3. `score = uptime_norm × poids_uptime + messages_norm × poids_messages`

Les VIP absents des tops (score API nul) apparaissent avec une valeur 0.
