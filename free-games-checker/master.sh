cd ~/free-games-checker/
./gog.sh > gog-results.txt
./steam.sh > steam-results.txt

if ! grep -q "GOG does not have any 100% off games" "gog-results.txt"; then
        ./send-gog.sh
fi

if ! grep -q "Steam does not have any 100% off games" "steam-results.txt"; then
        ./send-steam.sh
fi

./cleanup.sh