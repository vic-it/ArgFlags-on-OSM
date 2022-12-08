package util

import (
	"container/heap"
)

// NODE ID
// distances between nodes
// connections between nodes (edges) src and target
func CalculateDijkstra(graph Graph, sourceID int, destID int) (int, []int) {
	dist := make(map[int]int)
	prev := make(map[int]int)
	var prioQ = make(PriorityQueue, len(graph.Nodes))
	dist[sourceID] = 0
	i := 0
	for nodeID, _ := range graph.Nodes {
		// a bit bigger than the circumference of the earth -> in meters
		if nodeID != sourceID {
			dist[nodeID] = 50000000
		}
		prev[nodeID] = -1
		prioQ[i] = &Item{value: nodeID, priority: dist[nodeID], index: i}
		i++
	}
	heap.Init(&prioQ)

	for {
		node := heap.Pop(&prioQ).(*Item)
		if node.value == destID {
			break
		}
		neighbors := getGraphNeighbors(graph.Targets, graph.Offsets, graph.Weights, node.value)
		for _, neighbor := range neighbors {
			alt := dist[node.value] + neighbor[1]
			if alt < dist[neighbor[0]] {
				dist[neighbor[0]] = alt
				prev[neighbor[0]] = node.value
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
	for prev[currentNode] >= 0 {
		path = append(path, prev[currentNode])
		currentNode = prev[currentNode]
	}
	return dist[destID], path
}

func getGraphNeighbors(destinations []int, offsets []int, weights []int, nodeID int) [][]int {
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
