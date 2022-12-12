package util

import (
	"container/heap"
)

// calculates the shortest path between two nodes (on a graph) via dijkstras algorithm
func CalculateDijkstra(graph Graph, sourceID int, destID int) (int, []int) {
	dist := make(map[int]int)
	prev := make(map[int]int)
	//priority queue datastructure (see priority_queue.go)
	var prioQ = make(PriorityQueue, len(graph.Nodes))
	dist[sourceID] = 0
	i := 0
	for nodeID, _ := range graph.Nodes {
		// a bit bigger than the circumference of the earth -> in meters
		if nodeID != sourceID {
			//adds super high default distance
			dist[nodeID] = 50000000
		}
		//sets previous node ids no -1 since there is no path yet
		prev[nodeID] = -1
		//adds all nodes to the priority queue (heap)
		prioQ[i] = &Item{value: nodeID, priority: dist[nodeID], index: i}
		i++
	}
	heap.Init(&prioQ)

	for {
		//gets "best" next node
		node := heap.Pop(&prioQ).(*Item)
		// if we are at the destination then we break!
		if node.value == destID {
			break
		}
		// gets all neighbor/connected nodes
		neighbors := getGraphNeighbors(graph.Targets, graph.Offsets, graph.Weights, node.value)
		for _, neighbor := range neighbors {
			alt := dist[node.value] + neighbor[1]
			// neighbor [0] is target node ID
			if alt < dist[neighbor[0]] {
				dist[neighbor[0]] = alt
				prev[neighbor[0]] = node.value
				//just re-queue items with better value instead of updating it
				heap.Push(&prioQ, &Item{value: neighbor[0], priority: alt, index: neighbor[0]})
			}
		}
		if prioQ.Len() < 1 {
			dist[destID] = -1
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
	}
	//if distance is "-1" -> no path found,
	return dist[destID], path
}

// returns all neighbro node IDs connected to the input node
func getGraphNeighbors(destinations []int, offsets []int, weights []int, nodeID int) [][]int {
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
		neighborIDList = append(neighborIDList, []int{destinations[i], weights[i]})
	}
	return neighborIDList
}
