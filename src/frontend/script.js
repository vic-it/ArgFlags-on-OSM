var isStart = true;
//maybe store node IDs too in start dest objects
var start = {marker: null, lon: 0, lat: 0, id: -1}
var dest = {marker: null, lon: 0, lat: 0, id: -1}
var path = null


function clickHandler(event){
    let url = new URL("http://localhost:8080/getpoint");
    url.searchParams.append("lon", event.latlng.lng);
    url.searchParams.append("lat", event.latlng.lat);
    console.log("clicked point (lon/lat): \n["+event.latlng.lng.toFixed(6)+", "+event.latlng.lat.toFixed(6)+"]")
    
    //fetches result
    fetch(url).then((response) => {
        //extracts the response
        answer = response.text()
        return answer
    }).then((result) => {
        //takes body of response and splits it into the coordinates lon then lat
        coordinates = result.split("x")
        //0 lon, 1 lat, the + before the coordinates casts the string valued coordinates into number values
        //i hate javascript
        
        console.log("closest grid point (lon/lat): \n["+coordinates[0]+", "+coordinates[1]+"] id = "+coordinates[2])
        console.log("-----------------")
        addMarker(+coordinates[0],+coordinates[1], +coordinates[2])
    })
}

function routeHandler(){
    srcID = start.id
    destID = dest.id
    if(start.id <0 || dest.id<0){
        return
    }
    
    let url = new URL("http://localhost:8080/getroute");
    url.searchParams.append("src", srcID);
    url.searchParams.append("dest", destID);

    fetch(url).then((response) => {
        //extracts the response
        answer = response.text()
        return answer
    }).then((result) => {
        //takes body of response and splits it into the coordinates lon then lat
        distCoords = result.split("y")
        if (+distCoords[0] <0){
            console.log("NO PATH FOUND!")
        } else{
            distance = +distCoords[0]
            rawCoordinates = distCoords[1].split("x")
            var coordinates = [];
            for (const coord in rawCoordinates) {
                lat = +coord.split(z)[1]
                lon = +coord.split(z)[0]
                coordinates.append([lat, lon])
            }
            
            document.getElementById("distance").value = ""+distance+"m"
            addRoute(coordinates)
            console.log("Path found with distance: "+distance+"m")
        }
        //0 lon, 1 lat, the + before the coordinates casts the string valued coordinates into number values
        //i hate javascript
        
        console.log("closest grid point (lon/lat): \n["+coordinates[0]+", "+coordinates[1]+"] id = "+coordinates[2])
        console.log("-----------------")
        addMarker(+coordinates[0],+coordinates[1], +coordinates[2])
    })
}

//coordinates in lat-longs
function addRoute(coordinates){
    if(path !=null){
        map.removeLayer(path)
    }    
    path = L.polyline(coordinates, {color: 'blue'}).addTo(map)
}

// type 0 -> starting node, type 1 -> destination node
function addMarker(lon, lat, id){
    if(isStart){
        if(start.marker != null){
        map.removeLayer(start.marker)
        }
        start.marker = L.marker([lat, lon]).bindPopup("Source-ID: "+id+" | Lon: "+lon+" | Lat: "+lat)
        start.lon = lon
        start.lat = lat
        start.id = id
        start.marker.addTo(map)
        document.getElementById("idField1").value = id
        document.getElementById("lonField1").value = lon
        document.getElementById("latField1").value = lat
    } else {
        if(dest.marker !=null){
        map.removeLayer(dest.marker)
        }
        dest.marker = L.marker([lat, lon]).bindPopup("Dest.-ID: "+id+" | Lon: "+lon+" | Lat: "+lat)
        dest.lon = lon
        dest.lat = lat
        dest.id = id
        dest.marker.addTo(map)
        document.getElementById("idField2").value = id
        document.getElementById("lonField2").value = lon
        document.getElementById("latField2").value = lat
    }
    isStart =! isStart
}




// https://www.youtube.com/watch?v=Fk-P5l7DJjo
let mapOpions = {
    center: [0, 0],
    zoom: 2
}
var map = L.map('map', mapOpions);
let layer = L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png')
layer.addTo(map)
map.on("click", clickHandler)
//add "usable" box
var bounds = [[-90, -180], [90, 180]];
L.rectangle(bounds, {color: "#aaffaa", weight: 3, fillOpacity: 0}).addTo(map)
map.fitBounds(bounds)