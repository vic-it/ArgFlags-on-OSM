package backend

import (
	"container/heap"
	"time"
)

// calculates the shortest path between two nodes (on a graph) via dijkstras algorithm
func CalculateDijkstra(graph Graph, sourceID int, destID int) (int, []int, float64, float64, int) {
	initTime := time.Now()
	visited := make([]bool, len(graph.Nodes))
	dist := make([]int, len(graph.Nodes))
	prev := make([]int, len(graph.Nodes))
	nodesPoppedCounter := 0
	//priority queue datastructure (see priority_queue.go)
	prioQ := &PriorityQueue{{priority: 0, value: sourceID}}

	//simply iterating over every single node
	for nodeID := range graph.Nodes {
		dist[nodeID] = -1
		prev[nodeID] = -1
	}

	dist[sourceID] = 0
	initTimeDiff := time.Since(initTime).Seconds()
	searchTime := time.Now()
	for prioQ.Len() > 0 {
		//gets "best" next node
		node := heap.Pop(prioQ).(*Item)
		if visited[node.value] {
			continue
		}
		if node.value == destID {
			break
		}
		nodeID := node.value
		visited[nodeID] = true
		nodesPoppedCounter++
		// if we are at the destination then we break!

		// gets all neighbor/connected nodes
		startIndex := graph.Offsets[nodeID]
		endIndex := 0
		var neighbors [][]int
		if nodeID == len(graph.Offsets)-1 {
			endIndex = len(graph.Targets)
		} else {
			endIndex = graph.Offsets[nodeID+1]
		}
		for i := startIndex; i < endIndex; i++ {
			if visited[graph.Targets[i]] {
				continue
			}
			neighbors = append(neighbors, []int{graph.Targets[i], graph.Weights[i]})

		}
		//for each neighbor...
		for _, neighbor := range neighbors {
			alt := dist[node.value] + neighbor[1]
			// neighbor [0] is target node ID
			if dist[neighbor[0]] < 0 || alt < dist[neighbor[0]] {
				dist[neighbor[0]] = alt
				prev[neighbor[0]] = node.value
				//just re-queue items with better value instead of updating it
				heap.Push(prioQ, &Item{value: neighbor[0], priority: alt})
			}
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
	searchTimeDiff := time.Since(searchTime).Seconds()
	return dist[destID], path, initTimeDiff, searchTimeDiff, nodesPoppedCounter
}
