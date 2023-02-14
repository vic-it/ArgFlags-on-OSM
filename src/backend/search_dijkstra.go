package backend

import (
	"container/heap"
	"time"
)

var dijkstraVisited [1000000]bool
var dijkstraDistance [1000000]int
var dijkstraPrev [1000000]int

// calculates the shortest path between two nodes (on a graph) via dijkstras algorithm
func CalculateDijkstra(graph Graph, sourceID int, destID int) (int, []int, float64, float64, int) {

	initTime := time.Now()
	nodesPoppedCounter := 0
	//priority queue datastructure (see priority_queue.go)
	var prioQ = make(PriorityQueue, 1)

	for rowID, row := range graph.NodeInWaterMatrix {
		for columnID, isInWater := range row {
			nodeID := graph.NodeMatrix[rowID][columnID]
			if isInWater {
				dijkstraVisited[nodeID] = false
				dijkstraDistance[nodeID] = 50000000
				dijkstraPrev[nodeID] = -1
			}
		}
	}

	dijkstraDistance[sourceID] = 0
	prioQ[0] = &Item{value: sourceID, priority: dijkstraDistance[sourceID], index: 0}
	heap.Init(&prioQ)
	initTimeDiff := time.Since(initTime).Seconds()
	searchTime := time.Now()
	for {
		//gets "best" next node
		node := heap.Pop(&prioQ).(*Item)
		if node.value == destID {
			break
		}
		dijkstraVisited[node.value] = true
		nodesPoppedCounter++
		// if we are at the destination then we break!

		// gets all neighbor/connected nodes
		neighbors := GetGraphNeighbors(graph.Targets, graph.Offsets, graph.Weights, node.value)
		for _, neighbor := range neighbors {
			alt := dijkstraDistance[node.value] + neighbor[1]
			// neighbor [0] is target node ID
			if alt < dijkstraDistance[neighbor[0]] {
				dijkstraDistance[neighbor[0]] = alt
				dijkstraPrev[neighbor[0]] = node.value
				//just re-queue items with better value instead of updating it
				heap.Push(&prioQ, &Item{value: neighbor[0], priority: alt, index: neighbor[0]})
			}
		}
		if prioQ.Len() < 1 || dijkstraDistance[node.value] >= 50000000 {
			dijkstraDistance[destID] = -1
			break
		}
	}

	var path []int
	currentNode := destID
	path = append(path, currentNode)
	// starts from the destination node and iterates backwards to source node, creating the path
	for dijkstraPrev[currentNode] >= 0 {
		path = append(path, dijkstraPrev[currentNode])
		currentNode = dijkstraPrev[currentNode]
	}
	//if distance is "-1" -> no path found,
	searchTimeDiff := time.Since(searchTime).Seconds()
	return dijkstraDistance[destID], path, initTimeDiff, searchTimeDiff, nodesPoppedCounter
}

// returns all neighbro node IDs connected to the input node
func GetGraphNeighbors(destinations []int, offsets []int, weights []int, nodeID int) [][]int {
	// start index of edges determined by offset list
	startIndex := offsets[nodeID]
	endIndex := 0
	var neighborIDList [][]int
	if nodeID == len(offsets)-1 {
		endIndex = len(destinations)
	} else {
		endIndex = offsets[nodeID+1]
	}
	for i := startIndex; i < endIndex; i++ {
		if !dijkstraVisited[destinations[i]] {
			neighborIDList = append(neighborIDList, []int{destinations[i], weights[i]})
		}
	}
	return neighborIDList
}
