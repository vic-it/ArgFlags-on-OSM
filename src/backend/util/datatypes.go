package util

type basic struct {
	nodes map[int64]node
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
	edges  []edge
	points []point
	//more stuff here
}

// like nodes, except on a globe and in actualy position/distance instead of degrees
type point struct {
	x float64
	y float64
}

// one edge of a multi-edge polygon (= way)
type edge struct {
	start  float64
	end    float64
	length float64
}
