#!/bin/bash
mkdir gog-temp
cd gog-temp
curl -o raw.txt "$(cat 'gog-url.txt')" &> /dev/null
file="raw.txt"
search_string="Showing 0 games"
if grep -q "$search_string" "$file"; then
  echo "false"
else
  echo "true"
fi
echo "" > gog-temp.txt
grep -Eo 'Showing [0-9]+(\.[0-9]+)? games' raw.txt | grep -oE '[0-9]+(\.[0-9]+)?' > ../gog-number.txt
cd ..
rm -rf gog-temp