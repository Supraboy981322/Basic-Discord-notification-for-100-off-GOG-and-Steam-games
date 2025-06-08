#!/bin/bash
discord_url="YOUR_DISCORD_WEBHOOK_URL_GOES_HERE"

generate_post_data() {
  cat <<EOF
{
  "embeds": [{
    "title": "GOG has $(cat 'gog-number.txt') games are 100% off!",
    "description": "[You can see it here]($(cat 'gog-url.txt'))",
    "color": "10181046"
  }]
}
EOF
}
# POST request to Discord Webhook
curl -H "Content-Type: application/json" -X POST -d "$(generate_post_data)" $discord_url