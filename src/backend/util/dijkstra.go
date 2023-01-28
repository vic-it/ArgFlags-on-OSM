package util

import (
	"container/heap"
	"time"
)

var visited [1000000]bool
var distance [1000000]int
var prev [1000000]int

// calculates the shortest path between two nodes (on a graph) via dijkstras algorithm
func CalculateDijkstra(graph Graph, sourceID int, destID int) (int, []int, float64, float64, int) {

	//totalTime := time.Now()

	initTime := time.Now()
	nodesPoppedCounter := 0
	//priority queue datastructure (see priority_queue.go)
	var prioQ = make(PriorityQueue, 1)

	for rowID, row := range graph.NodeInWaterMatrix {
		for columnID, isInWater := range row {
			nodeID := graph.NodeMatrix[rowID][columnID]
			if isInWater {
				visited[nodeID] = false
				distance[nodeID] = 50000000
				prev[nodeID] = -1
				//prioQ[i] = &Item{value: nodeID, priority: dist[nodeID], index: i}
			}
		}
	}
	// for nodeID, _ := range graph.Nodes {
	// 	dist[nodeID] = 50000000
	// 	prev[nodeID] = -1
	// }

	distance[sourceID] = 0
	prioQ[0] = &Item{value: sourceID, priority: distance[sourceID], index: 0}
	heap.Init(&prioQ)
	initTimeDiff := time.Since(initTime).Seconds()
	//fmt.Printf("Time to initialize search: %.3fs\n", initTimeDiff)
	searchTime := time.Now()
	for {
		//gets "best" next node
		node := heap.Pop(&prioQ).(*Item)
		if node.value == destID {
			break
		}
		visited[node.value] = true
		nodesPoppedCounter++
		// if we are at the destination then we break!

		// gets all neighbor/connected nodes
		neighbors := GetGraphNeighbors(graph.Targets, graph.Offsets, graph.Weights, node.value)
		for _, neighbor := range neighbors {
			alt := distance[node.value] + neighbor[1]
			// neighbor [0] is target node ID
			if alt < distance[neighbor[0]] {
				distance[neighbor[0]] = alt
				prev[neighbor[0]] = node.value
				//just re-queue items with better value instead of updating it
				heap.Push(&prioQ, &Item{value: neighbor[0], priority: alt, index: neighbor[0]})
			}
		}
		if prioQ.Len() < 1 || distance[node.value] >= 50000000 {
			distance[destID] = -1
			break
		}
	}

	var path []int
	currentNode := destID
	path = append(path, currentNode)
	// starts from the destination node and iterates backwards to source node, creating the path
	for prev[currentNode] >= 0 {
		path = append(path, prev[currentNode])
		currentNode = prev[currentNode]
		//maybe error here??
		if currentNode < 0 {
			distance[destID] = -1
			break
		}
	}
	//if distance is "-1" -> no path found,
	searchTimeDiff := time.Since(searchTime).Seconds()
	// fmt.Printf("Time to search route: %.3fs\n", searchTimeDiff)
	// fmt.Printf("Time total to calculate route: %.3fs\n", time.Since(totalTime).Seconds())
	// fmt.Printf("distance: %dm\n", dist[destID])
	// fmt.Printf("nodes in path: %d\n", len(path))
	// fmt.Printf("Nodes popped: %d\n--\n", nodesPoppedCounter)
	return distance[destID], path, initTimeDiff, searchTimeDiff, nodesPoppedCounter
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
		if !visited[destinations[i]] {
			neighborIDList = append(neighborIDList, []int{destinations[i], weights[i]})
		}
	}
	return neighborIDList
}
