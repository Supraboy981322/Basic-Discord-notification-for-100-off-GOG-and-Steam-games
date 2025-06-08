#!/bin/bash
mkdir steam-temp
cd steam-temp
curl -o raw.txt "$(cat 'steam-url.txt')" &> /dev/null
file="raw.txt"
search_string="Showing 0 games"
if grep -q "$search_string" "$file"; then
  echo "false"
else
  echo "true"
fi
echo "" > steam-temp.txt
grep -Eo '[0-9]+(\.[0-9]+)? results match your search' raw.txt | grep -oE '[0-9]+(\.[0-9]+)?' > ../steam-number.txt
cd ..
rm -rf steam-temp