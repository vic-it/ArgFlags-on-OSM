This is a programming project for a university course. The goal was to implement a much faster, more sophisticated alternative to the common dijkstra/A* approaches for shortest path computations on a graph.
The Arc-Flag approach provided in this project allows for a 20 to 100 times faster computation of the shortest path in a roughly 1,000,000 node graph.
The backend was mainly written in GO while the GUI was made in leaflet and javascript.

How to use:
1.  Export this ZIP file to any directory but keep the following structure within the project:	 
	|
	|- /OSM
	|	- /data
	|		- graph.graph
	|		- arc.flags
	|			
	|	- /OSMFMI
	|		- /src
	|		- go.sum
	|		- all other files...
						

2. Download and install Golang (GO)
	- the program was tested on a "Ubuntu 20.04 64bit" VM
	- with go version "go 1.19.4"
	
3. Open a terminal in the OSM/OSMFMI/src directory and type "go run main.go"
	- In the current version the program will import the pre-processed (by the same program) graph and arc flags respectively
	- It will then open a web interface on port 8080

4. Go to localhost:8080 in a browser of your choice (tested on firefox and chrome)
	- Source and Destination nodes are chosen with clicks on the map
		- It will snap to the (roughly) nearest node on water
		- the clicks will alternate between source and destination nodes
	- On the top you can see the nodes you chose with their IDs, longitudes and latitudes
	- To calculate the shortest path between the nodes press "Calculate Route"
		- If there is a valid path the distance will be shown next to the button and the path will be drawn onto the map
		- If there is no valid path the same text field will say so!
	- To run tests (anywhere from 100 to 10000 runs), at the bottom of the page there is an intuitive interface for it
	- Do note that because the nodes are spread equidistantly on the graph, there are usually multiple shortest paths
		- because of this the dijkstra and the arc flags algorithm might produce different paths, although both are optimal and identical in length
		
