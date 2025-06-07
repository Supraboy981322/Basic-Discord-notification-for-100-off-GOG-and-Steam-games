mkdir steam-temp
cd steam-temp
curl -o raw.txt "https://store.steampowered.com/search/?ignore_preferences=1&maxprice=free&category1=998%2C993%2C997%2C996%2C990%2C994&specials=1&ndl=1" &> /dev/null
file="raw.txt"
search_string="Showing 0 games"
if grep -q "$search_string" "$file"; then
  echo "Steam does not have any 100% off games"
else
  echo "Steam has games that are 100% off!"
fi
echo "" > steam-temp.txt
grep -Eo '[0-9]+(\.[0-9]+)? results match your search' raw.txt | grep -oE '[0-9]+(\.[0-9]+)?' > ../steam-number.txt
cd ..
rm -rf steam-temp