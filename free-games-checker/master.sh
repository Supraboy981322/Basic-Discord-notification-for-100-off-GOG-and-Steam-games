#!/bin/bash
cd ~/free-games-checker/
./gog.sh > gog-results.txt
./steam.sh > steam-results.txt

if $(cat 'gog-results.txt'); then
        ./send-gog.sh
fi

if $(cat 'steam-results.txt'); then
        ./send-steam.sh
fi

./cleanup.sh