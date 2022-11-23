var isStart = true;
//maybe store node IDs too in start dest objects
var start = {marker: null, lon: 0, lat: 0}
var dest = {marker: null, lon: 0, lat: 0}


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
        
        console.log("closest grid point (lon/lat): \n["+coordinates[0]+", "+coordinates[1]+"]")
        console.log("-----------------")
        addMarker(+coordinates[0],+coordinates[1])
    })
}


// type 0 -> starting node, type 1 -> destination node
function addMarker(lon, lat){
    if(isStart){
        if(start.marker != null){
        map.removeLayer(start.marker)
        }
        start.marker = L.marker([lat, lon]).bindPopup("start")
        start.lon = lon
        start.lat = lat
        start.marker.addTo(map)
    } else {
        if(dest.marker !=null){
        map.removeLayer(dest.marker)
        }
        dest.marker = L.marker([lat, lon]).bindPopup("destination")
        dest.lon = lon
        dest.lat = lat
        dest.marker.addTo(map)
    }
    isStart =! isStart
}
function calculateRoute(){

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