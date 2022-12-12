1.  Export this ZIP file to any directory but keep the following structure within the project:	 
	|
- /OSM
	|	- /data
	|		- graph.graph
	|		- possibly other files...
	|			
	|	- /OSM-FMI
	|		- /src
	|		- go.sum
	|		- all other files...
						

2. Download and install Golang (GO)
	- the program was tested on a "Ubuntu 20.04 64bit" VM
	- with go version "go 1.19.4"
	
3. Open a terminal in the OSM/OSM-FMI/src directory and type "go run main.go"
	- In the current version the program will import a pre-processed (by the same program) graph with roughly 1 million nodes
	- It will then open a web interface on port 8080

4. Go to localhost:8080 in a browser of your choice (tested on firefox and chrome)
	- Source and Destination nodes are chosen with clicks on the map
		- It will snap to the (roughly) nearest node on water
		- the clicks will alternate between source and destination nodes
	- On the top you can see the nodes you chose with their IDs, longitudes and latitudes
	- To calculate the shortest path between the nodes press "Calculate Route"
		- If there is a valid path the distance will be shown next to the button and the path will be drawn onto the map
		- If there is no valid path the same text field will say so!
	- In case anything goes wrong, there is a backup graph in /OSM-FMI
		