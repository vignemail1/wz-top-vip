# wz-top-vip

Classement pondéré des VIP d'une chaîne Twitch à partir des tops WizeBot
(temps de visionnage et nombre de messages).

Le programme interroge l'API WizeBot pour récupérer les classements mensuels
ou hebdomadaires, filtre les résultats sur une liste de VIPs fournie, normalise
les deux métriques et calcule un score pondéré pour chaque VIP.

---

## Prérequis

- Un compte [WizeBot](https://wizebot.tv) associé à votre chaîne Twitch
- Une clé API WizeBot en lecture (`[R]`) — voir la section [Clé API](#clé-api)
- La liste des VIPs de votre chaîne — voir la section [Liste des VIPs](#liste-des-vips)

---

## Clé API

La clé API WizeBot en lecture est nécessaire pour interroger les classements.

1. Connectez-vous à votre espace WizeBot
2. Rendez-vous sur la page de gestion des APIs :
   **https://panel.wizebot.tv/development_api_management#**
3. Copiez la clé **[R]** (lecture seule)

La clé peut être fournie de trois façons (par ordre de priorité) :

| Méthode | Exemple |
|---|---|
| Flag `-apikey` | `wz-top-vip -apikey VOTRE_CLE vips.txt` |
| Variable d'environnement | `export WIZEBOT_API_READ=VOTRE_CLE` |
| Saisie interactive | Le programme la demande au lancement si absente |

Lors de la saisie interactive, les caractères sont masqués et une confirmation
du nombre de caractères saisis est affichée.

---

## Liste des VIPs

### Récupérer la liste depuis Twitch

1. Ouvrez votre navigateur et allez sur votre chaîne Twitch :
   `https://www.twitch.tv/VOTRE_CHAINE`
2. Ouvrez le chat
3. Tapez la commande `/vips` dans le chat et envoyez-la

Twitch affiche une liste comme :

```
Les VIPs actuels de cette chaîne sont : pseudo1, pseudo2, pseudo3, pseudo4.
```

### Créer le fichier `vips.txt`

Créez un fichier texte nommé `vips.txt` (ou tout autre nom de votre choix)
avec **un pseudo par ligne**, en reprenant exactement les pseudos affichés par Twitch :

```
pseudo1
pseudo2
pseudo3
pseudo4
```

Règles du fichier :
- Un pseudo Twitch par ligne
- Les lignes vides sont ignorées
- Les lignes commençant par `#` sont ignorées (commentaires)
- La casse n'a pas d'importance (`MonPseudo` = `monpseudo`)

Exemple de fichier avec commentaires :

```
# VIPs de la chaîne — mis à jour mai 2026
pseudo1
pseudo2
# pseudo3  ← temporairement exclu
pseudo4
```

---

## Installation

### Télécharger la dernière release

Rendez-vous sur la page des releases GitHub :
**https://github.com/vignemail1/wz-top-vip/releases/latest**

Téléchargez l'archive correspondant à votre système d'exploitation :

| Système | Fichier à télécharger |
|---|---|
| Linux (Intel/AMD 64 bits) | `wz-top-vip-linux-amd64.tar.gz` |
| Windows (Intel/AMD 64 bits) | `wz-top-vip-windows-amd64.zip` |
| macOS (Apple Silicon M1/M2/M3) | `wz-top-vip-darwin-arm64.tar.gz` |

### Décompresser l'archive

**Linux / macOS — Terminal :**

```bash
tar -xzf wz-top-vip-linux-amd64.tar.gz    # Linux
tar -xzf wz-top-vip-darwin-arm64.tar.gz   # macOS
```

**Windows — Explorateur de fichiers :**

Faites un clic droit sur le fichier `.zip` → **Extraire tout...**

Ou depuis PowerShell :

```powershell
Expand-Archive wz-top-vip-windows-amd64.zip .
```

### Rendre le binaire exécutable (Linux / macOS)

```bash
chmod +x wz-top-vip-linux-amd64    # Linux
chmod +x wz-top-vip-darwin-arm64   # macOS
```

---

## Utilisation

### Ouvrir un terminal

- **Linux** : `Ctrl+Alt+T` ou recherchez « Terminal » dans vos applications
- **macOS** : `Cmd+Espace` → tapez `Terminal` → Entrée
- **Windows** : `Win+R` → tapez `powershell` → Entrée  
  *(ou cherchez « PowerShell » dans le menu Démarrer)*

Placez-vous dans le dossier contenant le binaire et votre fichier `vips.txt` :

```bash
cd ~/Téléchargements    # adaptez selon votre dossier
```

### Afficher l'aide

```bash
# Linux
./wz-top-vip-linux-amd64 -help

# macOS
./wz-top-vip-darwin-arm64 -help

# Windows (PowerShell)
.\wz-top-vip-windows-amd64.exe -help
```

Sortie :

```
Usage: wz-top-vip [flags] <vip-file>

  <vip-file>  Path to a text file containing one Twitch username per line (VIP list).

Flags:
  -apikey string
        WizeBot read API key (overrides WIZEBOT_API_READ env var)
  -fetch-limit int
        Number of entries to fetch from WizeBot API per ranking (1-100) (default 100)
  -message-weight int
        Weight of messages in score (0-100, as percentage) (default 50)
  -period string
        Time period: week or month (default "month")
  -top int
        Number of top VIPs to display (default 3)

Environment variables:
  WIZEBOT_API_READ  WizeBot read API key (used if -apikey is not set)

API key location:
  https://panel.wizebot.tv/development_api_management#
```

### Lancer le programme

**Exemple minimal** (top 3 du mois, poids 50 % / 50 %) :

```bash
# Linux
./wz-top-vip-linux-amd64 vips.txt

# macOS
./wz-top-vip-darwin-arm64 vips.txt

# Windows (PowerShell)
.\wz-top-vip-windows-amd64.exe vips.txt
```

Si la clé API n'est pas configurée, le programme la demande :

```
Aucune clé API WizeBot trouvée (-apikey / WIZEBOT_API_READ).
Vous pouvez obtenir votre clé API [R] (lecture) ici :
  https://panel.wizebot.tv/development_api_management#

Clé API WizeBot [R] : ****************
Clé saisie : **************** (16 caractères)
```

**Exemple avancé** — top 5 de la semaine, 70 % messages / 30 % uptime :

```bash
./wz-top-vip-darwin-arm64 -period week -message-weight 70 -top 5 vips.txt
```

**Avec la clé en variable d'environnement (recommandé) :**

```bash
# Linux / macOS
export WIZEBOT_API_READ=VOTRE_CLE
./wz-top-vip-darwin-arm64 vips.txt

# Windows (PowerShell)
$env:WIZEBOT_API_READ = "VOTRE_CLE"
.\wz-top-vip-windows-amd64.exe vips.txt
```

---

## Exemple de sortie

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

---

## Flags disponibles

| Flag | Défaut | Description |
|---|---|---|
| `-apikey` | — | Clé API WizeBot [R] (sinon `WIZEBOT_API_READ`) |
| `-period` | `month` | Période : `week` ou `month` |
| `-message-weight` | `50` | Poids des messages dans le score (0–100, en %) |
| `-top` | `3` | Nombre de VIP à afficher |
| `-fetch-limit` | `100` | Entrées récupérées par top API (max 100) |

---

## Méthode de scoring

Pour chaque VIP présent dans au moins un des deux tops WizeBot :

1. `uptime_norm   = uptime_brut   / max(uptime_brut)`
2. `messages_norm = messages_brut / max(messages_brut)`
3. `score = uptime_norm × poids_uptime + messages_norm × poids_messages`

Les VIP absents des tops (hors du top 100 de la période) apparaissent avec un score de 0.

---

## Compilation depuis les sources

Requiert Go 1.24+.

```bash
git clone https://github.com/vignemail1/wz-top-vip.git
cd wz-top-vip
make build          # binaire local
make build-all      # Linux + Windows + macOS dans dist/
make vet            # go vet ./...
```

---

## Licence

MIT
