package backend

import (
	"container/heap"
	"time"
)

func CalculateArcFlagDijkstra(graph Graph, sourceID int, destID int, arcData ArcData, nodePartitionList []int) (int, []int, float64, float64, int) {
	nodePartitionMatrix := arcData.NodePartitionMatrix
	arcFlags := arcData.ArcFlags
	//totalTime := time.Now()
	destPartition := 0
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
				//prioQ[i] = &Item{value: nodeID, priority: dist[nodeID], index: i}
			}
			if nodeID == destID {
				destPartition = nodePartitionMatrix[rowID][columnID]
			}
		}
	}
	// for nodeID, _ := range graph.Nodes {
	// 	dist[nodeID] = 50000000
	// 	prev[nodeID] = -1
	// }

	dijkstraDistance[sourceID] = 0
	prioQ[0] = &Item{value: sourceID, priority: dijkstraDistance[sourceID], index: 0}
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
		dijkstraVisited[node.value] = true
		nodesPoppedCounter++
		// if we are at the destination then we break!

		// gets all neighbor/connected nodes

		currentNodePartition := nodePartitionList[node.value] //
		neighbors := getArcFlagRouteNeighbors(graph.Targets, graph.Offsets, graph.Weights, node.value, arcFlags, destPartition, currentNodePartition)
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
	//println(x / 1000)
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
	// fmt.Printf("Time to search route: %.3fs\n", searchTimeDiff)
	// fmt.Printf("Time total to calculate route: %.3fs\n", time.Since(totalTime).Seconds())
	// fmt.Printf("distance: %dm\n", dist[destID])
	// fmt.Printf("nodes in path: %d\n", len(path))
	// fmt.Printf("Nodes popped: %d\n--\n", nodesPoppedCounter)
	return dijkstraDistance[destID], path, initTimeDiff, searchTimeDiff, nodesPoppedCounter
}

// returns the neighbors by their [nodeID, distance]
func getArcFlagRouteNeighbors(destinations []int, offsets []int, weights []int, nodeID int, arcFlags [][]bool, destPartition int, currentNodePartition int) [][]int {
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
		if (arcFlags[i][destPartition] || destPartition == currentNodePartition) && !dijkstraVisited[destinations[i]] {
			neighborIDList = append(neighborIDList, []int{destinations[i], weights[i]})
		}
	}
	return neighborIDList
}
