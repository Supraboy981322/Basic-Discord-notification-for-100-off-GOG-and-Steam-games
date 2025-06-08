# A Basic Discord bot for games that are 100% off on GOG and Steam
I wrote this because I just felt like writing my own Shell script to recieve notifications on Discord for games that are 100% off on GOG and Steam

When the cronjob starts it downloads the Steam and GOG web pages, then looks for keywords on the page, and sends a Discord notification to a Discord webhook address if any are found.


---
# Instructions:
- Download the setup shell script and DO NOT execute it as su.
    
        wget "https://github.com/Supraboy981322/Basic-Discord-notification-for-100-off-GOG-and-Steam-games/raw/refs/heads/main/setup.sh" && sudo chmod +x setup.sh && sudo setup.sh

- The setup script will tell you what to do from here.