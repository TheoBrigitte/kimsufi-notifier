#!/bin/bash
#
# Display hardware and country code for Kimsufi server models.

URL=${URL:-https://www.kimsufi.com/fr/serveurs.xml}

echo "> fetching data from $URL" 1>&2
DATA=$(curl -qSs "$URL")
echo "> fetched data"

# count number of tables
table_count=$(echo $DATA | pup 'div#main table.homepage-table' -n)
tables=$(echo $DATA | pup 'div#main table.homepage-table')

# process each table
for (( i=1; i<=$table_count; i++ )); do
	# count number of rows
	table=$(echo $tables | pup 'table:nth-child('$i')')
	row_count=$(echo $table | pup 'tr' -n)

	# get country code from header row
	country=$(echo $table | pup 'tr:nth-child(1) th:nth-child(10) span attr{class}' | tr -d '[:space:]')

	# get model name and hardware code from each row
	# skip first header row
	for (( j=2; j<=$row_count; j++)); do
		name=$(echo $table | pup 'tr:nth-child('$j') td:first-child text{}' | tr -d '[:space:]')
		hardware=$(echo $table | pup 'tr:nth-child('$j') td:nth-child(11) attr{data-ref}' | tr -d '[:space:]')
		cpu=$(echo $table | pup 'tr:nth-child('$j') td:nth-child(2) span:nth-child(2) text{}' | tr -d '[:space:]')
		threads=$(echo $table | pup 'tr:nth-child('$j') td:nth-child(3) span text{}' | tr -d '[:space:]')
		freq=$(echo $table | pup 'tr:nth-child('$j') td:nth-child(4) span text{}' | tr -d '[:space:]')
		ram=$(echo $table | pup 'tr:nth-child('$j') td:nth-child(5) text{}')
		disk=$(echo $table | pup 'tr:nth-child('$j') td:nth-child(6) text{}' | awk '{$1=$1};1' | tr -d '\n')
		network=$(echo $table | pup 'tr:nth-child('$j') td:nth-child(7) text{}' | tr -d '[:space:]' )
		price=$(echo $table | pup 'tr:nth-child('$j') td:nth-child(9) span span text{}' | tr -d '[:space:]' )
		availability=$(echo $table | pup 'tr:nth-child('$j') td:nth-child(11) style text{}' | tr -d '[:space:]')
		echo -e "model=$name\thardware=$hardware\tcpu=$cpu\t$freq ($threads)\tram=$ram \tdisk=$disk \tnetwork=$network\tprice=$price\tcountry=$country\tavailable=$availability"
	done
done
exit
