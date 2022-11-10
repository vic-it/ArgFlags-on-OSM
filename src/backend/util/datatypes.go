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
	nodes []node
	ways  []way
	//more stuff here
}
