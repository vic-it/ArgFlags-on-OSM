package backend

type EdgeCoordinate struct {
	edgeID     int
	coordinate float64
}

type Basic struct {
	//the ID of a way in this map is the ID of its first node -> for merge function
	ways map[int64]way
}

type way struct {
	nodes      []int64
	lastNodeID int64
}

type Graph struct {
	// list of Nodes with nodeID: [longitude, latitude]
	Nodes [][]float64
	// list of source nodes
	Sources []int
	// list of edge destinations
	Targets []int
	// distances of edges
	Weights []int
	// offset for when edges for another node start
	Offsets []int
	// 2D matrix of nodes on the grid with the respective node IDs
	NodeMatrix [][]int
	// respective 2D matrix for the "PointMatrix" but instead of IDs it stores whether the node is in water or on land
	NodeInWaterMatrix [][]bool
	// number of nodes intended to create - usually close to len[nodes] but a bit higher
	intendedNodeQuantity int
	// number of nodes that are in water
	countOfWaterNodes int
}

type ArcData struct {
	ArcFlags            [][]bool
	NodePartitionMatrix [][]int
	NumberOfPartitions  int
}

type Coastline struct {
	Nodes             map[int64][]float64
	Edges             [][]int64
	SortedLonEdgeList []EdgeCoordinate
	MaxLonDiffs       []float64
	maxLonDiff        float64
}

type ArcFlagInfo struct {
	PartitionMatrix         [][]int
	NumberOfTotalPartitions int
	NodePartitions          map[int64][]int64
}
