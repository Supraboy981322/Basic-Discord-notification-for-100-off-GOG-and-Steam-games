#!/bin/bash

cd ~/
echo "NOTE: this script should only be run without sudo"
echo "downloading package"
echo "The package will be placed in your /home/$USER directory"
# download package
mkdir free-games-checker
cd free-games-checker
wget "https://github.com/Supraboy981322/Basic-Discord-notification-for-100-off-GOG-and-Steam-games/raw/refs/heads/main/free-games-checker.zip" &> /dev/null
unzip free-games-checker.zip > /dev/null
rm free-games-checker.zip

# ask user for discord api webhook url
echo "In order to communicate with Discord, this script needs a Discord API webhook URL, please create one if not already done."
echo "If you do not know how to do this, Google is your friend, you can use it."
read DiscordURL
echo "Setting the value of 'discord_url' in 'send-gog.sh' and 'send-steam.sh'"
# set the value of "discord_url" in send-gog.sh and send-steam.sh
DEFAULT_DISCORD_URL="YOUR_DISCORD_WEBHOOK_URL_GOES_HERE"
sed -i "s|${DEFAULT_DISCORD_URL}|${DiscordURL}|g" send-steam.sh
sed -i "s|${DEFAULT_DISCORD_URL}|${DiscordURL}|g" send-gog.sh

# ask user how frequently to check for sales and what time
echo "The script needs to know how frequently to check for game sales"
echo "Keep in mind that it will notify you if there is a 100% off game everytime it finds one, even if it has already notified you of it."
echo "How often should it run?"
echo "Daily (d), weekly (w), monthly (m)"
read -p "daily (d), weekly (w), monthly (m): " FREQUENCY
echo "At what hour of the day should it run?"
echo "Use 24 hour format, ex: '23' (which is 11 PM)"
read HOUR
echo "What minute of the hour?"
echo "ex: '45' (for 45 minutes after the hour)"
read MINUTE

# create cronjob

if [[ "$FREQUENCY" == "d" ]]; then 
   DAY="*"
   WEEK="*"
   MONTH="*"
elif [[ "$FREQUENCY" == "w" ]]; then
   DAY="*"
   MONTH="*"
   WEEK="0"
elif [[ "$FREQUENCY" == "m" ]]; then
   DAY="1"
   MONTH="*"
   WEEK="*"
fi

crontab -l > mycron
#echo new cron into cron file
echo "$MINUTE $HOUR $DAY $MONTH $WEEK $HOME/free-games-checker/./master.sh" >> mycron
#install new cron file
crontab mycron
rm mycron

echo "Thank you for installing my Steam and GOG free game sale Discord bot."