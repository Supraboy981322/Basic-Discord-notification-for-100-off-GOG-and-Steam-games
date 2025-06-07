discord_url="YOUR_DISCORD_WEBHOOK_URL_GOES_HERE"

generate_post_data() {
  cat <<EOF
{
  "embeds": [{
    "title": "$(cat 'steam-results.txt')",
    "description": "$(cat 'steam-number.txt') games are 100% off on Steam.",
    "color": "3447003"
  }]
}
EOF
}
# POST request to Discord Webhook
curl -H "Content-Type: application/json" -X POST -d "$(generate_post_data)" $discord_url