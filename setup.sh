echo "If you did not, please be sure to run this script under sudo"
echo "downloading package"
echo "The package will be placed in your /home/$USER directory"
# download package

# ask user for discord api webhook url
echo "In order to communicate with Discord, this script needs a Discord API webhook URL, please create one if not already done."
echo "If you do not know how to do this, Google is your friend, you can use it."
# make temporary variable for that webhook url

echo "Setting the value of 'discord_url' in 'send-gog.sh' and 'send-steam.sh'"
# set the value of "discord_url" in send-gog.sh and send-steam.sh

# ask user how frequently to check for sales and what time
echo "The script needs to know how frequently to check for game sales"
echo "Keep in mind that it will notify you if there is a 100% off game everytime it finds one, even if it has already notified you of it."

# create cronjob

echo "Thank you for installing my basic Steam and GOG free game sale Discord bot."