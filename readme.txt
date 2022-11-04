1. get test data
	- go to https://download.geofabrik.de/index.html
	- download antarctica.osm.pbf

2. extract coast lines
	- download Osmosis https://github.com/openstreetmap/osmosis/releases/latest
	- set path variable https://learnosm.org/en/osm-data/osmosis/
	- open cli (cmd) and navigate (cd) to folder with test data
	- write "osmosis --read-pbf ANTARCTICA_FILE_NAME.osm.pbf --way-key-value keyValueList="natural.coastline" --write-xml output.osm"
		- documentation here: https://wiki.openstreetmap.org/wiki/Osmosis#Example_usage
		- and here: https://wiki.openstreetmap.org/wiki/Osmosis/Detailed_Usage_0.48
	-> you now have the extracted coast lines as OSM XML file

3. map OSM XML file to GeoJSON file
	- download and install "Anaconda"/"conda-forge" https://www.anaconda.com/products/distribution
	- download gdal over conda (which includes ogr2ogr) https://anaconda.org/conda-forge/gdal
	- type "ogr2ogr" and if you get a "missing ...dll file" try "conda install krb5" to somehow (?) install the missing dll file
	- download and install libsqlite3 and libexpat (somehow? TODO) 
	- see build dependencies https://gdal.org/drivers/vector/osm.html
	- hopefully transform the OSM XML file created in [2] to a json (or geojson) file with "ogr2ogr -f GeoJSON coastlines.geojson output.osm"
	- upload to geojson.io and check results

4. TODO...

X. Extra
	- download and install go (Golang)
		- I installed go itself as well as the language support for VSCode (language support can be installed from within vscode as an extension)
	- DO NOT UPLOAD THE BIG .XML .OSM .OSM.PBF OR SIMILAR FILES TO THE GITHUB REPOSITORY (preferably keep them in a seperate folder/directory, outside of the repository)
	
		