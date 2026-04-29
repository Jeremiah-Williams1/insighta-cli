# Insighta Labs+ — CLI

Command line interface for the Insighta Labs+ platform.

## Installation

```bash
go install github.com/jeremiah-williams1/insighta-cli@latest
```

Or build locally:
```bash
git clone <repo>
cd insighta-cli
go build -o insighta .
sudo mv insighta /usr/local/bin/
```

## Environment Variables

```env
API_BASE_URL=http://localhost:8080  # or your deployed backend URL
GITHUB_CLIENT_ID=                   # same as backend
```

## Authentication

```bash
insighta login     # Opens GitHub OAuth in browser, saves tokens locally
insighta logout    # Clears saved tokens
insighta whoami    # Shows your user ID and role
```

Tokens are stored at `~/.insighta/credentials.json`.
The CLI handles token refresh automatically.

## Profile Commands

```bash
# List profiles
insighta profiles list
insighta profiles list --gender male
insighta profiles list --country NG --age-group adult
insighta profiles list --min-age 25 --max-age 40
insighta profiles list --sort-by age --order desc
insighta profiles list --page 2 --limit 20

# Get a profile by ID
insighta profiles get <id>

# Search using natural language
insighta profiles search "young males from nigeria"

# Create a profile (admin only)
insighta profiles create --name "Harriet Tubman"

# Export to CSV
insighta profiles export --format csv
insighta profiles export --format csv --gender male --country NG
```

## Token Handling

Access tokens expire every 3 minutes. The CLI automatically attempts a refresh using the stored refresh token. If the refresh token has also expired, you will be prompted to run `insighta login` again.