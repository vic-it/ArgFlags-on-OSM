1. get test data
	- go to https://download.geofabrik.de/index.html
	- download antarctica.osm.pbf

2. extract coast lines	
	- current go program extracts them automatically (main.go -> backend/util/coastlines.go)

3. map OSM XML file to GeoJSON file
	- TODO: implement "BasictoGEOJson" function in transform.go (backend/util/transform.go)

4. TODO...

X. Extra
	- download and install go (Golang)
		- I installed go itself as well as the language support for VSCode (language support can be installed from within vscode as an extension)
	- DO NOT UPLOAD THE BIG .XML .OSM .OSM.PBF OR SIMILAR LARGE FILES TO THE GITHUB REPOSITORY
		keep them in a seperate folder, e.g. in a structure like:
			-OSM
				-OSM-Project (this is the whole github project)
					-src
						-main.go
						-...
					-...
				-data
					-antarctica.osm.pbf ("light weight dummy data" for testing the whole program)
					-...
	- start program with "go run main.go"
	- make sure the filepath declared in main.go points to your antarctica.osm.pbf file
		