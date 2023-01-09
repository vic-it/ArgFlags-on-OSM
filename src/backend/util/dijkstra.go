package util

import (
	"container/heap"
	"fmt"
	"time"
)

// calculates the shortest path between two nodes (on a graph) via dijkstras algorithm
func CalculateDijkstra(graph Graph, sourceID int, destID int) (int, []int) {

	totalTime := time.Now()

	initTime := time.Now()
	nodesPoppedCounter := 0
	dist := make(map[int]int)
	prev := make(map[int]int)
	visited := make(map[int]bool)
	//priority queue datastructure (see priority_queue.go)
	var prioQ = make(PriorityQueue, graph.countOfWaterNodes)
	dist[sourceID] = 0
	i := 0

	for rowID, row := range graph.NodeInWaterMatrix {
		for columnID, isInWater := range row {
			nodeID := graph.NodeMatrix[rowID][columnID]
			if isInWater {
				if nodeID != sourceID {
					//adds super high default distance
					dist[nodeID] = 50000000
				}
				visited[nodeID] = false
				prev[nodeID] = -1
				//adds all nodes to the priority queue (heap)
				prioQ[i] = &Item{value: nodeID, priority: dist[nodeID], index: i}
				i++
			}
		}
	}
	heap.Init(&prioQ)

	fmt.Printf("Time to initialize search: %.3fs\n", time.Since(initTime).Seconds())
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
		neighbors := getGraphNeighbors(graph.Targets, graph.Offsets, graph.Weights, node.value, visited)
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
		if prioQ.Len() < 1 || dist[node.value] >= 50000000 {
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
	fmt.Printf("Time to search route: %.3fs\n", time.Since(searchTime).Seconds())
	fmt.Printf("Time total to calculate route: %.3fs\n", time.Since(totalTime).Seconds())
	fmt.Printf("distance: %dm\n", dist[destID])
	fmt.Printf("nodes in path: %d\n", len(path))
	fmt.Printf("Nodes popped: %d\n--\n", nodesPoppedCounter)
	return dist[destID], path
}

// returns all neighbro node IDs connected to the input node
func getGraphNeighbors(destinations []int, offsets []int, weights []int, nodeID int, visited map[int]bool) [][]int {
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
