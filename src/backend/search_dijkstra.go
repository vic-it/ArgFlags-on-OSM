package backend

import (
	"container/heap"
	"time"
)

var DistCopy []int
var PrevCopy []int

// calculates the shortest path between two nodes (on a graph) via dijkstras algorithm
func CalculateDijkstra(graph Graph, sourceID int, destID int) (int, []int, float64, float64, int) {
	initTime := time.Now()
	numOfNodes := len(graph.Nodes)
	visited := make([]bool, numOfNodes)
	dist := make([]int, numOfNodes)
	prev := make([]int, numOfNodes)

	// copies pre-filled lists to save time
	copy(prev, PrevCopy)
	copy(dist, DistCopy)
	nodesPoppedCounter := 0

	//priority queue datastructure (see priority_queue.go)
	prioQ := &PriorityQueue{{priority: 0, value: sourceID}}

	dist[sourceID] = 0
	initTimeDiff := time.Since(initTime).Seconds()
	searchTime := time.Now()
	for prioQ.Len() > 0 {
		//gets "best" next node
		node := heap.Pop(prioQ).(*Item)
		currentNodeID := node.value
		// skip previously popped nodes (because we dont update PQ, we re-push)
		if visited[currentNodeID] {
			continue
		}
		// if we are at the destination we break!
		if currentNodeID == destID {
			break
		}
		visited[currentNodeID] = true
		nodesPoppedCounter++

		// gets all neighbor/connected nodes
		startIndex := graph.Offsets[currentNodeID]
		endIndex := graph.Offsets[currentNodeID+1]

		for i := startIndex; i < endIndex; i++ {
			neighbor := graph.Targets[i]
			// skip nodes we already popped
			if visited[neighbor] {
				continue
			}
			alt := dist[currentNodeID] + graph.Weights[i]
			if alt < dist[neighbor] {
				dist[neighbor] = alt
				prev[neighbor] = currentNodeID
				//just re-queue items with better value instead of updating it
				heap.Push(prioQ, &Item{value: neighbor, priority: alt})
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
	if dist[destID] >= 500000000 {
		dist[destID] = -1
	}
	//if distance is "-1" -> no path found,
	searchTimeDiff := time.Since(searchTime).Seconds()
	return dist[destID], path, initTimeDiff, searchTimeDiff, nodesPoppedCounter
}

func PrepArrays(n int) {
	for i := 0; i < n; i++ {
		// prepare these values beforehand because copying them to a new slice is more efficient than preparing them for each run
		DistCopy = append(DistCopy, 500000000)
		PrevCopy = append(PrevCopy, -1)
	}
}
