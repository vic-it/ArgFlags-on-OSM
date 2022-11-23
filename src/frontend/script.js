var isStart = true;
var start = {marker: null, lon: 0, lat: 0}
var dest = {marker: null, lon: 0, lat: 0}


function clickHandler(event){
    console.log(event)
    let url = new URL("http://localhost:8080/getpoint");
    url.searchParams.append("lon", event.latlng.lng);
    url.searchParams.append("lat", event.latlng.lat);
    
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
        addMarker(+coordinates[0],+coordinates[1])
    })
}


// type 0 -> starting node, type 1 -> destination node
function addMarker(lon, lat){
    if(isStart){
        if(start.marker != null){
        map.removeLayer(start.marker)
        }
        start.marker = L.marker([lat, lon])
        start.lon = lon
        start.lat = lat
        start.marker.addTo(map)
    } else {
        if(dest.marker !=null){
        map.removeLayer(dest.marker)
        }
        dest.marker = L.marker([lat, lon])
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
    zoom: 5
}
var map = L.map('map', mapOpions);
let layer = L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png')
layer.addTo(map)
map.on("click", clickHandler)