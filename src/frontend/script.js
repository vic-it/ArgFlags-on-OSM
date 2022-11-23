
function clickHandler(event){
    console.log(event)
    let url = new URL("http://localhost:8080/getpoint");
    url.searchParams.append("lon", {lon: event.latlng.lon});
    url.searchParams.append("lat", {lat: event.latlng.lat});
    fetch(url).then((response) => {
        console.log(response.text())
    })
}

function addMarker(lon, lat, type){

}

function getClosestGridPoint({lon, lat}){

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