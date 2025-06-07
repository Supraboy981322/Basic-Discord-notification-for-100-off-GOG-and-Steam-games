mkdir gog-temp
cd gog-temp
curl -o raw.txt "https://www.gog.com/en/games?priceRange=0,0&discounted=true&hideDLCs=true" &> /dev/null
file="raw.txt"
search_string="Showing 0 games"
if grep -q "$search_string" "$file"; then
  echo "GOG does not have any 100% off games"
else
  echo "GOG has games that are 100% off!"
fi
echo "" > gog-temp.txt
grep -Eo 'Showing [0-9]+(\.[0-9]+)? games' raw.txt | grep -oE '[0-9]+(\.[0-9]+)?' > ../gog-number.txt
cd ..
rm -rf gog-temp