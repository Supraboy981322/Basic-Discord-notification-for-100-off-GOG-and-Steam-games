<p align="center">
  <img src="https://github.com/Supraboy981322/Basic-Discord-notification-for-100-off-GOG-and-Steam-games/raw/refs/heads/main/logo.png">
</p>

# A Basic Discord bot for games that are 100% off on GOG and Steam
I wrote this because I just felt like writing my own Shell script to recieve notifications on Discord for games that are 100% off on GOG and Steam

> [!NOTE]
> I am doing a rewrite of this in Go (in my spare time), but I will only update this repository each time a new store is completely integrated.

> [!INFO]
> I am still waiting for my friend to test the setup process, so it's not fully polished yet. Also, I'm not completely finished either, it works perfectly as is, but I just want to make some `.txt` files for long variable values that are defined in more than one file.

When the cronjob starts it downloads the Steam and GOG web pages, then looks for keywords on the page, and sends a Discord notification to a Discord webhook address if any are found.


---
# Instructions:
- Download the setup shell script and DO NOT execute it as `su`.
    
        wget "https://github.com/Supraboy981322/Basic-Discord-notification-for-100-off-GOG-and-Steam-games/raw/refs/heads/main/setup.sh" && sudo chmod +x setup.sh && ./setup.sh

- The setup script will tell you what to do from here.
