package util

type EdgeCoordinate struct {
	edgeID     int
	coordinate float64
}

type Basic struct {
	Nodes map[int64]node
	//the ID of a way in this map is the ID of its first node -> for merge function
	ways map[int64]way
}

type tag struct {
	key   string
	value string
}

type node struct {
	//lat and lon are in degrees, not absolute position
	lat float64
	lon float64
}

type way struct {
	nodes      []int64
	lastNodeID int64
}

type graph struct {
	// list of edges with edgeID: [firstNodeID, secondNodeID]
	edges [][]int64
	// list of nodes with nodeID: [longitude, latitude]
	nodes [][]float64
	// list of source nodes
	sources []int
	targets []int
	weights []int
	offsets []int
}

// like nodes, except on a globe and in actualy position/distance instead of degrees
type point struct {
	x float64
	y float64
}

type point_threeD struct {
	x float64
	y float64
	z float64
}

// one edge of a multi-edge polygon (= way)
type edge struct {
	start  float64
	end    float64
	length float64
}

type Coastline struct {
	Nodes             map[int64][]float64
	Edges             [][]int64
	SortedLonEdgeList []EdgeCoordinate
	MaxLonDiff        float64
}
