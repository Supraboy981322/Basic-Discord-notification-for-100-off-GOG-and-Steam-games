<p align="center">
  <img src="https://github.com/Supraboy981322/Basic-Discord-notification-for-100-off-GOG-and-Steam-games/raw/refs/heads/main/logo.png">
</p>

> [!Warning]
> The install script is currently being re-written, please do not use the install script.

# free-games-checker

### A Basic Discord bot for games that are 100% off on GOG and Steam

A basic Discord bot that searches Steam and GOG (with more stores planned) for games that are 100% off.

This was originally written in Bash, but has been rewritten in Go. Due to the fact that some stores, like Epic (which is planned), have measures in place to block bots, it might take a while to rewrite my original Bash scripts for scraping them in Go, as I need to research libraries I can use (and those stores aren't my highest priority at the moment). For Itch, there is no filtering or sorting built-in to their search, so, originally, I used a headless browser with a userscript that injects filtering and sorting into the webpage, but I want to do this properly for my rewrite. So, any stores beyond what I have now (GOG and Steam) might take a while for me to get around to writting new scrapers. I still have the original Bash scripts, but I am not open-sourcing them because they are super messy and very janky (there's also a LOT of dependencies, which I neglected to make note of, so you'd have to manually find and install all the dependencies).

---

# Instructions:
- Manually download and compile. Better instructions will be created after the install script is written.
