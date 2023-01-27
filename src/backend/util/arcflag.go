package util

import (
	"container/heap"
	"fmt"
	"math"
	"sync"
	"time"
)

var wg sync.WaitGroup
var processTimer time.Time

func PreprocessArcFlags(graph Graph, numOfRows int, numOfPoleRowPartitions int) ([][]bool, [][]int) {
	println("Generating arc flags...")
	maxThreads := 10
	fmt.Printf("Preprocessing on %d threads...\n", maxThreads)
	arcFlags := [][]bool{}
	nodePartitionMatrix, numOfPartitions := createPartitions(graph, numOfRows, numOfPoleRowPartitions)

	//fill empty arc flags array
	for i := 0; i < len(graph.Sources); i++ {
		tmp := []bool{}
		for j := 0; j < numOfPartitions; j++ {
			tmp = append(tmp, false)
		}
		arcFlags = append(arcFlags, tmp)
	}

	boundaryNodeIDs := getBoundaryNodeIDS(graph, nodePartitionMatrix)

	processTimer = time.Now()
	ctr := 0
	altTime := time.Now()

	//calculate arc flags in batches of "maxThreads" at once until all boundary nodes went through
	// calcArcDijkstraForNode(graph, 51230, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 510, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 54320, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 521360, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 51230, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 12350, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 150, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 56340, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 530, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 23250, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 12350, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 152150, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 58560, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 55320, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 98750, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 251230, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 424350, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 64350, arcFlags, numOfPartitions, nodePartitionMatrix, false)
	// calcArcDijkstraForNode(graph, 50, arcFlags, numOfPartitions, nodePartitionMatrix, false)

	fmt.Printf("alt vers: %.3fs\n", time.Since(altTime).Seconds())
	totalTime := time.Now()

	println("Starting multithreading")
	for {
		ctr++
		start := maxThreads * (ctr - 1)
		end := maxThreads * ctr
		if ctr%10 == 1 {
			printPreProcessProgress(start, len(boundaryNodeIDs))
		}

		if len(boundaryNodeIDs) < end {
			end = len(boundaryNodeIDs)
		}
		for i := start; i < end; i++ {
			wg.Add(1)
			go calcArcDijkstraForNode(graph, boundaryNodeIDs[i], arcFlags, numOfPartitions, nodePartitionMatrix, true)
		}
		wg.Wait()
		if end >= len(boundaryNodeIDs) {
			break
		}
	}

	ensureBidirectionality(graph, arcFlags)
	fmt.Printf("Time to generate arc flags: %.3fs\n", time.Since(totalTime).Seconds())
	return arcFlags, nodePartitionMatrix
}
func createPartitions(graph Graph, numOfRows int, numOfPoleRowPartitions int) ([][]int, int) {
	// nodeCount := len(graph.Nodes)
	rowCount := len(graph.NodeMatrix)
	graphRowsPerPartitionRow := int(math.Ceil(float64(len(graph.NodeMatrix)) / float64(numOfRows)))
	partitionRows := [][]int{}
	nodesPerRow := []int{}
	partitionsPerRow := []int{}
	currentGRowID := 0

	for i := 0; i < numOfRows; i++ {
		rowOfIDs := []int{}
		numOfNodesInThisRow := 0
		for len(rowOfIDs) < graphRowsPerPartitionRow && currentGRowID < rowCount {
			rowOfIDs = append(rowOfIDs, currentGRowID)
			numOfNodesInThisRow += len(graph.NodeMatrix[currentGRowID])
			currentGRowID++
		}
		nodesPerRow = append(nodesPerRow, numOfNodesInThisRow)
		partitionRows = append(partitionRows, rowOfIDs)
	}
	for id, _ := range partitionRows {
		//calc how many partitions per row
		partitionsPerRow = append(partitionsPerRow, int(math.Ceil(float64(nodesPerRow[id])/(float64(nodesPerRow[0])/float64(numOfPoleRowPartitions)))))
		// fmt.Printf("partition row %d: %d rows - %d nodes - %d partitions\n", id, len(row), nodesPerRow[id], partitionsPerRow[id])
	}
	fmt.Printf("Rough nodes per partition: %d\n---\n", nodesPerRow[0]/numOfPoleRowPartitions)
	//---------------------------------------------------------------------------------
	//---------------------------------------------------------------------------------
	//---------------------------------------------------------------------------------
	//---------------------------------------------------------------------------------
	//---------------------------------------------------------------------------------
	//---------------------------------------------------------------------------------
	//---------------------------------------------------------------------------------
	//---------------------------------------------------------------------------------
	//---------------------------------------------------------------------------------
	// calc column split of graph matrix rows
	//same as nodematrix but values are IDs of partitions
	numberOfPartitions := 0
	nodePartitionMatrix := [][]int{}
	partitionCounterStart := 0
	for partRowID, partRow := range partitionRows {
		numOfColumns := partitionsPerRow[partRowID]
		//for every row in the original node matrix determine the cutoffs between the partition columns
		for _, graphRowID := range partRow {
			nodePartitionMatrix = append(nodePartitionMatrix, []int{})
			//nodes in row / columns rounded up
			nodesInThisGraphRow := len(graph.NodeMatrix[graphRowID])
			nodesPerColumn := int(math.Round(float64(nodesInThisGraphRow) / float64(numOfColumns)))
			//start id of this clumn
			startID := 0
			//end id of this column (element of this id not included)
			cutOffID := 0
			//current column we are adding
			colCtr := 0
			for colCtr < numOfColumns {
				cutOffID = (colCtr + 1) * nodesPerColumn
				startID = colCtr * nodesPerColumn
				// 3 columns for 10 nodes should go -> 4, 4, 2
				for i := startID; i < cutOffID && i < nodesInThisGraphRow; i++ {
					//only add partition of node if node exists
					nodePartitionMatrix[graphRowID] = append(nodePartitionMatrix[graphRowID], partitionCounterStart+colCtr)
					numberOfPartitions = partitionCounterStart + colCtr + 1
				}

				colCtr++
				//determine each cut off
			}
			//fill rows if necessary (e.g. if round down for nodespercolumn)
			for len(nodePartitionMatrix[graphRowID]) < len(graph.NodeMatrix[graphRowID]) {
				nodePartitionMatrix[graphRowID] = append(nodePartitionMatrix[graphRowID], partitionCounterStart+colCtr-1)
				numberOfPartitions = partitionCounterStart + colCtr - 1 + 1
			}
		}
		partitionCounterStart += numOfColumns
	}
	// DIVIDE PARTITION ROWS INTO PARTITIONS

	fmt.Printf("node matrix rows: %d\npartition matrix rows: %d\n", len(graph.NodeMatrix), len(nodePartitionMatrix))
	fmt.Printf("number of partitions total: %d\n", numberOfPartitions)
	// for rowID, row := range graph.NodeMatrix {
	// 	fmt.Printf("diff: %d\n", len(row)-len(nodePartitionMatrix[rowID]))
	// 	fmt.Println(nodePartitionMatrix[rowID])
	// }
	//fmt.Println(nodePartitionMatrix[0])
	// fmt.Println(nodePartitionMatrix[1])
	// fmt.Println(nodePartitionMatrix[2])
	// fmt.Println(nodePartitionMatrix[4])
	// fmt.Println(nodePartitionMatrix[5])
	// fmt.Println(nodePartitionMatrix[6])
	// println("---------")
	// fmt.Println(nodePartitionMatrix[len(nodePartitionMatrix)-7])
	// fmt.Println(nodePartitionMatrix[len(nodePartitionMatrix)-6])
	// fmt.Println(nodePartitionMatrix[len(nodePartitionMatrix)-5])
	// fmt.Println(nodePartitionMatrix[len(nodePartitionMatrix)-4])
	// fmt.Println(nodePartitionMatrix[len(nodePartitionMatrix)-3])
	// fmt.Println(nodePartitionMatrix[len(nodePartitionMatrix)-2])
	// fmt.Println(nodePartitionMatrix[len(nodePartitionMatrix)-1])
	return nodePartitionMatrix, numberOfPartitions
}

func getBoundaryNodeIDS(graph Graph, nodePartitionMatrix [][]int) []int {
	boundaryNodeIDS := []int{}
	for rowID, row := range graph.NodeMatrix {
		for colID, nodeID := range row {
			currentNodePartition := nodePartitionMatrix[rowID][colID]
			neighList := GetGraphNeighbors(graph.Targets, graph.Offsets, graph.Weights, nodeID)
			for _, idAndDistance := range neighList {
				row, col := GetNodeMatrixIndex(idAndDistance[0], graph)
				if nodePartitionMatrix[row][col] != currentNodePartition {
					boundaryNodeIDS = append(boundaryNodeIDS, nodeID)

				}
			}
		}
	}
	return boundaryNodeIDS
}

func calcArcDijkstraForNode(graph Graph, sourceID int, arcFlags [][]bool, numOfPartitions int, nodePartitionMatrix [][]int, test bool) {

	// totalTime := time.Now()
	var distance []int
	//here prev means previous edgeID
	var prev []int
	//sorts all popped nodes by their distance to source node, starting at the source, going to the furthest (actually reachable) node
	var ascendingPoppedList []int
	//list of true/false for nodeID - true if node was in sub-tree of previously checked path -> path for this node also checked
	var checked []bool
	nodesPoppedCounter := 0
	//priority queue datastructure (see priority_queue.go)
	var prioQ = make(PriorityQueue, 1)

	for i := 0; i < len(graph.Nodes); i++ {
		prev = append(prev, -1)
		distance = append(distance, 50000000)
		checked = append(checked, false)
	}

	distance[sourceID] = 0
	prioQ[0] = &Item{value: sourceID, priority: distance[sourceID], index: 0}
	heap.Init(&prioQ)
	for {
		//gets "best" next node
		node := heap.Pop(&prioQ).(*Item)
		if distance[node.value] >= 50000000 {
			break
		}
		ascendingPoppedList = append(ascendingPoppedList, node.value)

		nodesPoppedCounter++
		// if we are at the destination then we break!

		// gets all neighbor/connected nodes
		neighbors := getArcFlagPreProcessNeighbors(graph.Targets, graph.Offsets, graph.Weights, node.value)
		for _, neighbor := range neighbors {
			alt := distance[node.value] + neighbor[1]
			// neighbor [0] is target node ID, neighbor[2] is ID of edge which goes to this node
			if alt < distance[neighbor[0]] {
				distance[neighbor[0]] = alt
				prev[neighbor[0]] = neighbor[2]
				//just re-queue items with better value instead of updating it
				heap.Push(&prioQ, &Item{value: neighbor[0], priority: alt, index: neighbor[0]})
			}
		}
		if prioQ.Len() < 1 {
			break
		}
	}
	//maybe not over all edges multiple times????????? array for each node of flags for this node-- maybe als only over boundary nodes and then do nonboundary nodes afterwards alltogether
	// for id, val := range prev {
	// 	if val >= 0 {
	// 		addArcFlags(graph, id, arcFlags, prev, numOfPartitions, nodePartitionMatrix)
	// 	}
	// }
	// fmt.Printf("time to calc full dijkstra: %.3fs\n", time.Since(totalTime).Seconds())
	x := 0
	// for _, nodeID := range ascendingPoppedList {
	// 	if prev[nodeID] >= 0 && !checked[nodeID] {
	// 		addArcFlags(graph, nodeID, arcFlags, prev, numOfPartitions, nodePartitionMatrix, checked)
	// 		x++
	// 	}
	// }
	for i := len(ascendingPoppedList) - 1; i >= 0; i-- {
		nodeID := ascendingPoppedList[i]
		if prev[nodeID] >= 0 && !checked[nodeID] {
			addArcFlags(graph, nodeID, arcFlags, prev, numOfPartitions, nodePartitionMatrix, checked)
			x++
		}
	}

	// fmt.Printf("nodes backwards iterated through: %d\n", x)
	if test {
		defer wg.Done()
	}
}

func addArcFlags(graph Graph, nodeID int, arcFlags [][]bool, prev []int, numOfPartitions int, nodePartitionMatrix [][]int, checkList []bool) {
	currNode := nodeID
	partitionFlags := []bool{}
	edgeIDList := []int{}
	//fill empty partition flags
	for i := 0; i <= numOfPartitions; i++ {
		partitionFlags = append(partitionFlags, false)
	}
	//collect all edges we go through to get to source node, as well as all partitions we move through (all flags we need to set to 1)
	for prev[currNode] >= 0 {
		checkList[currNode] = true
		row, col := GetNodeMatrixIndex(currNode, graph)
		partitionFlags[nodePartitionMatrix[row][col]] = true
		edgeIDList = append(edgeIDList, prev[currNode])
		currNode = graph.Sources[prev[currNode]]
	}
	for _, edge := range edgeIDList {
		for i, flag := range partitionFlags {
			if flag {
				arcFlags[edge][i] = flag
			}
		}
	}
}

func getArcFlagPreProcessNeighbors(destinations []int, offsets []int, weights []int, nodeID int) [][]int {
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
		neighborIDList = append(neighborIDList, []int{destinations[i], weights[i], i})
	}
	return neighborIDList
}

func printPreProcessProgress(current int, max int) {
	progress := float64(current) / float64(max)
	fmt.Printf("Arc flag pre-processing | Progress: %2.2f%s%d%s%d boundary nodes completed (%.3fs)\n\r", 100*progress, "% - ", current, " out of ", max, time.Since(processTimer).Seconds())

	processTimer = time.Now()
}
func CalculateArcFlagDijkstra(graph Graph, sourceID int, destID int, arcFlags [][]bool, nodePartitionMatrix [][]int) (int, []int, float64, float64, int) {

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
				visited[nodeID] = false
				distance[nodeID] = 50000000
				prev[nodeID] = -1
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
		neighbors := getArcFlagRouteNeighbors(graph.Targets, graph.Offsets, graph.Weights, node.value, arcFlags, destPartition)
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

func getArcFlagRouteNeighbors(destinations []int, offsets []int, weights []int, nodeID int, arcFlags [][]bool, destPartition int) [][]int {
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
		if !visited[destinations[i]] && arcFlags[i][destPartition] {
			neighborIDList = append(neighborIDList, []int{destinations[i], weights[i]})
		}
	}
	return neighborIDList
}

func ensureBidirectionality(graph Graph, arcFlags [][]bool) {
	println("Ensuring bi-directionality of arc flags...")
	timer := time.Now()
	for edgeID, flags := range arcFlags {
		reverseEdgeID := getReverseEdgeID(graph, edgeID)
		if reverseEdgeID >= 0 {
			for partitionID, flag := range flags {
				if flag {
					arcFlags[reverseEdgeID][partitionID] = true
				}
			}
		}
	}
	getReverseEdgeID(graph, 200)
	fmt.Printf("Time to ensure bi-directionality of arc flags: %.3fs\n", time.Since(timer).Seconds())
}
