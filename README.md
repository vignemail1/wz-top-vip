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
| Option `-apikey` | `wz-top-vip -apikey VOTRE_CLE vips.txt` |
| Variable d'environnement | `WIZEBOT_API_READ=VOTRE_CLE` |
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

Créez un fichier texte nommé `vips.txt` avec **un pseudo par ligne**,
en reprenant exactement les pseudos affichés par Twitch :

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

**Windows — Explorateur de fichiers :**

Faites un clic droit sur le fichier `.zip` → **Extraire tout...**
Choisissez un dossier facile d'accès (par exemple `Bureau` ou `Téléchargements`).

**Linux / macOS — Terminal :**

```bash
tar -xzf wz-top-vip-linux-amd64.tar.gz    # Linux
tar -xzf wz-top-vip-darwin-arm64.tar.gz   # macOS
chmod +x wz-top-vip-linux-amd64            # Linux
chmod +x wz-top-vip-darwin-arm64           # macOS
```

---

## Utilisation sous Windows

### Méthode rapide — double-clic

Cette méthode ne nécessite aucune connaissance technique.

1. Extrayez le `.zip` dans un dossier (voir section [Installation](#installation))
2. Placez votre fichier `vips.txt` **dans le même dossier** que le fichier
   `wz-top-vip-windows-amd64.exe`
3. **Double-cliquez** sur `wz-top-vip-windows-amd64.exe`
4. Une fenêtre noire (console) s'ouvre
5. Si la clé API n'est pas encore configurée, le programme la demande :

```
Aucune clé API WizeBot trouvée (-apikey / WIZEBOT_API_READ).
Vous pouvez obtenir votre clé API [R] (lecture) ici :
  https://panel.wizebot.tv/development_api_management#

Clé API WizeBot [R] : ****************
Clé saisie : **************** (16 caractères)
```

6. Saisissez votre clé API (les caractères ne s'affichent pas) puis appuyez sur **Entrée**
7. Le classement s'affiche
8. Appuyez sur **Entrée** pour fermer la fenêtre

> **Astuce :** Pour ne plus avoir à saisir la clé à chaque lancement, configurez-la
> une fois en variable d'environnement Windows (voir la méthode avancée ci-dessous).

---

### Méthode avancée — PowerShell

Utilisez cette méthode pour passer des options (période, poids, top N...)
ou pour configurer la clé API de façon permanente.

#### Ouvrir PowerShell dans le bon dossier

Dans l'Explorateur de fichiers, naviguez jusqu'au dossier contenant le `.exe`,
puis dans la barre d'adresse tapez `powershell` et appuyez sur **Entrée**.

Ou ouvrez PowerShell depuis le menu Démarrer et naviguez :

```powershell
cd C:\Users\VotreNom\Downloads\wz-top-vip
```

#### Afficher l'aide

```powershell
.\wz-top-vip-windows-amd64.exe -help
```

#### Lancement simple (fichier `vips.txt` dans le même dossier)

```powershell
.\wz-top-vip-windows-amd64.exe
```

#### Lancement avec un fichier VIP spécifique

```powershell
.\wz-top-vip-windows-amd64.exe C:\Users\VotreNom\Documents\mes-vips.txt
```

#### Passer la clé API en argument (usage ponctuel)

```powershell
.\wz-top-vip-windows-amd64.exe -apikey VOTRE_CLE vips.txt
```

#### Configurer la clé API pour la session PowerShell en cours

```powershell
$env:WIZEBOT_API_READ = "VOTRE_CLE"
.\wz-top-vip-windows-amd64.exe vips.txt
```

#### Configurer la clé API de façon permanente (pour tous les lancements futurs)

```powershell
[System.Environment]::SetEnvironmentVariable("WIZEBOT_API_READ", "VOTRE_CLE", "User")
```

Fermez et réouvrez PowerShell (ou double-cliquez sur le `.exe`) pour que le
changement soit pris en compte.

#### Exemples avec options

```powershell
# Top 5 VIP de la semaine
.\wz-top-vip-windows-amd64.exe -period week -top 5 vips.txt

# Top 3 du mois, 70 % messages / 30 % uptime
.\wz-top-vip-windows-amd64.exe -message-weight 70 vips.txt

# Top 10 de la semaine, 100 % uptime (ignorer les messages)
.\wz-top-vip-windows-amd64.exe -period week -message-weight 0 -top 10 vips.txt
```

---

## Utilisation sous Linux / macOS

### Ouvrir un terminal

- **Linux** : `Ctrl+Alt+T` ou recherchez « Terminal » dans vos applications
- **macOS** : `Cmd+Espace` → tapez `Terminal` → Entrée

Placez-vous dans le dossier contenant le binaire :

```bash
cd ~/Téléchargements/wz-top-vip    # adaptez selon votre dossier
```

### Afficher l'aide

```bash
./wz-top-vip-linux-amd64 -help     # Linux
./wz-top-vip-darwin-arm64 -help    # macOS
```

### Lancement simple

```bash
# Avec vips.txt dans le même dossier
./wz-top-vip-linux-amd64

# Avec un fichier VIP spécifique
./wz-top-vip-darwin-arm64 ~/Documents/mes-vips.txt
```

### Clé API via variable d'environnement (recommandé)

```bash
export WIZEBOT_API_READ="VOTRE_CLE"
./wz-top-vip-linux-amd64 vips.txt
```

Pour la rendre permanente, ajoutez la ligne `export WIZEBOT_API_READ="VOTRE_CLE"`
dans votre `~/.bashrc`, `~/.zshrc` ou `~/.profile`.

### Exemples avec options

```bash
# Top 5 VIP de la semaine
./wz-top-vip-linux-amd64 -period week -top 5 vips.txt

# Top 3 du mois, 70 % messages / 30 % uptime
./wz-top-vip-darwin-arm64 -message-weight 70 vips.txt

# Top 10 de la semaine, 100 % uptime
./wz-top-vip-linux-amd64 -period week -message-weight 0 -top 10 vips.txt
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

## Options disponibles

| Option | Défaut | Description |
|---|---|---|
| `-apikey` | — | Clé API WizeBot [R] (sinon `WIZEBOT_API_READ`) |
| `-period` | `month` | Période : `week` (semaine) ou `month` (mois) |
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
