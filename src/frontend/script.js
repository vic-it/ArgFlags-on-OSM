var isStart = true;
//maybe store node IDs too in start dest objects
var start = {marker: null, lon: 0, lat: 0, id: -1}
var dest = {marker: null, lon: 0, lat: 0, id: -1}
var path1 = null
var path2 = null
var path3 = null
var paths = []
var estTimePerTest = 0.2
var slider = document.getElementById("numOfTests");
var output = document.getElementById("numDisplay");
var queryIntervall

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
function abortHandler(){
    let url = new URL("http://localhost:8080/abort");  
    document.querySelector('#abortButton').disabled = true;
    fetch(url)
    clearInterval(queryIntervall)
    setTimeout(()=>{
        document.querySelector('#testButton').disabled = false;
        document.querySelector('#routeButton').disabled = false;
        document.getElementById("dijkstraBar").style.width = "0%"
        document.getElementById("dijkstraBar").innerHTML ='';
        document.getElementById("arcflagBar").style.width = "0%"
        document.getElementById("arcflagBar").innerHTML ='';},500)
}
function testsHandler(){
    isAborted = false
    document.querySelector('#testButton').disabled = true;
    document.querySelector('#routeButton').disabled = true;
    document.querySelector('#abortButton').disabled = false;
    document.getElementById("dijkstraBar").style.width = "0%"
    document.getElementById("dijkstraBar").innerHTML ='';
    document.getElementById("arcflagBar").style.width = "0%"
    document.getElementById("arcflagBar").innerHTML ='';
    var numOfTests = slider.value
    let url = new URL("http://localhost:8080/testalgorithms");
    url.searchParams.append("num", numOfTests);    
    fetch(url)
    queryIntervall = setInterval(()=>{
        let url = new URL("http://localhost:8080/querytestprogress");
        fetch(url).then((response) => {
            //extracts the response
            answer = response.text()
            return answer
        }).then((result) => {
            //takes body of response and splits it into the coordinates lon then lat
            answer = result.split("-")
            dProgress = +answer[0]
            aProgress = +answer[1]
            updateBars(dProgress, aProgress)
            //0 lon, 1 lat, the + before the coordinates casts the string valued coordinates into number values
            //i hate javascript
            if(dProgress== 100.0){
                setTimeout(()=>{
                    document.getElementById("dijkstraResult").innerHTML = answer[2]},600)            }            
            if(aProgress== 100.0){                                
                setTimeout(()=>{
                    document.getElementById("arcflagResult").innerHTML = answer[3]},600)
            }
            if(dProgress == 100.0 && aProgress == 100.0){
                document.querySelector('#testButton').disabled = false;
                document.querySelector('#routeButton').disabled = false;
                document.querySelector('#abortButton').disabled = true;
                clearInterval(queryIntervall)
            }
        })
    },500) 
    //document.querySelector('#testButton').disabled = false;
}

function updateBars(dProg, aProg){
    var dBar = document.getElementById("dijkstraBar");
    var aBar = document.getElementById("arcflagBar");
    var dWidth = +dBar.style.width.replace("%","");
    var aWidth = +aBar.style.width.replace("%","");
    var dStepSize = (dProg - dWidth) / 25.0
    var aStepSize = (aProg - aWidth) / 25.0
    var dIntervall = setInterval(dFrame, 20);
    function dFrame() {
        if (dWidth >= dProg) {
            clearInterval(dIntervall);
        } else {
            dBar.style.width = dWidth + '%';
            dBar.innerHTML = dWidth.toFixed(1) + '%';
            dWidth += dStepSize;
        }
        // if(dProg == 100){            
        //     dBar.style.width = '100%';
        //     dBar.innerHTML = '100%';
        // }
    }
    var aIntervall = setInterval(aFrame, 20);
    function aFrame() {
        if (aWidth >= aProg) {
            clearInterval(aIntervall);
        } else {
            aBar.style.width = aWidth + '%';
            aBar.innerHTML = aWidth.toFixed(1)  + '%';
            aWidth += aStepSize;
        }
        // if(aProg == 100){            
        //     aBar.style.width = '100%';
        //     aBar.innerHTML = '100%';
        // }
    }
      
}

function routeHandler(){
    srcID = start.id
    destID = dest.id
    mode = document.querySelector('input[name="alg"]:checked').value;
    if(start.id <0 || dest.id<0){
        return
    }
    
    let url = new URL("http://localhost:8080/getroute");
    url.searchParams.append("src", srcID);
    url.searchParams.append("dest", destID);
    url.searchParams.append("mode", mode);

    fetch(url).then((response) => {
        //extracts the response
        answer = response.text()
        return answer
    }).then((result) => {
        //takes body of response and splits it into the coordinates lon then lat
        for(var path of paths){
            map.removeLayer(path)
        }
        distCoords = result.split("y")
        if (+distCoords[0] <0 || +distCoords[0]>45000000){
            document.getElementById("distance").value = "NO PATH FOUND"   
        } else{
            distance = +distCoords[0]
            rawCoordinates = distCoords[4].split("x")
            var coordinates = [];
            for (i = 0; i < rawCoordinates.length-1; i++){
                c = rawCoordinates[i].split("z")
                lat = +c[1]
                lon = +c[0]
                coordinates.push([lat, lon])
                if(lon>90 && +rawCoordinates[i+1].split("z")[0]<-90){
                    coordinates.push([lat, 180])
                    paths.push(L.polyline(coordinates, {color: 'blue'}).addTo(map))
                    coordinates = []
                    coordinates.push([lat, -179.5])
                } else if(lon<-90 && +rawCoordinates[i+1].split("z")[0]>90){
                    coordinates.push([lat, -179.5])
                    paths.push(L.polyline(coordinates, {color: 'blue'}).addTo(map))
                    coordinates = []
                    coordinates.push([lat, 180])
                    }
            }
                c = rawCoordinates[rawCoordinates.length-1].split("z")
                lat = +c[1]
                lon = +c[0]
                coordinates.push([lat, lon])
                paths.push(L.polyline(coordinates, {color: 'blue'}).addTo(map))
            document.getElementById("distance").value = ""+distance+"km"
        }
        nodesPopped = distCoords[1]
        initTime = +distCoords[2]
        searchTime = +distCoords[3]
        //total time in ms, above times in s
        totalTime = Math. round((initTime+searchTime)*1000)
        document.getElementById("nodes").value = ""+nodesPopped
        document.getElementById("time").value = ""+totalTime+"ms"
    })
}

//coordinates in lat-longs
function drawRoute(coordinates){      
    
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

slider.oninput = function() {
    output.innerHTML = this.value;    
    minutes = Math.floor(estTimePerTest*this.value / 60);
    seconds = Math.floor(estTimePerTest*this.value - minutes * 60);
    document.querySelector('#estTime').innerHTML = minutes+"m "+seconds+"s"
}
slider.oninput()

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