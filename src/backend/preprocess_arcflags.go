package backend

import (
	"container/heap"
	"fmt"
	"math"
	"sync"
	"time"
)

var wg sync.WaitGroup
var processTimer time.Time

func PreprocessArcFlags(graph Graph, numOfRows int, numOfPoleRowPartitions int) ArcData {
	println("Generating arc flags...")
	maxThreads := 16
	fmt.Printf("Preprocessing on %d threads...\n", maxThreads)
	arcFlags := [][]bool{}
	nodePartitionMatrix, numOfPartitions := CreatePartitions(graph, numOfRows, numOfPoleRowPartitions)

	//fill empty arc flags array
	for i := 0; i < len(graph.Sources); i++ {
		tmp := []bool{}
		for j := 0; j < numOfPartitions; j++ {
			tmp = append(tmp, false)
		}
		arcFlags = append(arcFlags, tmp)
	}
	boundaryNodeIDs := getBoundaryNodeIDS(graph, nodePartitionMatrix)
	ctr := 0
	processTimer = time.Now()
	totalTime := time.Now()
	println("Starting multithreading")
	for {
		ctr++
		start := maxThreads * (ctr - 1)
		end := maxThreads * ctr
		if ctr%125 == 1 {
			printPreProcessProgress(start, len(boundaryNodeIDs))
		}

		if len(boundaryNodeIDs) < end {
			end = len(boundaryNodeIDs)
		}
		for i := start; i < end; i++ {
			wg.Add(1)
			go singleSourceArcFlagPreprocess(graph, boundaryNodeIDs[i], arcFlags, numOfPartitions, nodePartitionMatrix, true)
		}
		wg.Wait()
		if end >= len(boundaryNodeIDs) {
			break
		}
	}
	fmt.Printf("Time to generate arc flags: %.3fs\n", time.Since(totalTime).Seconds())
	return ArcData{arcFlags, nodePartitionMatrix, numOfPartitions}
}

func CreatePartitions(graph Graph, numOfRows int, numOfPoleRowPartitions int) ([][]int, int) {
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
	for id := range partitionRows {
		//calc how many partitions per row
		partitionsPerRow = append(partitionsPerRow, int(math.Ceil(float64(nodesPerRow[id])/(float64(nodesPerRow[0])/float64(numOfPoleRowPartitions)))))
	}

	// calc column split of graph matrix rows
	//same as nodematrix but values are IDs of partitions
	numberOfPartitions := 0
	nodePartitionMatrix := [][]int{}
	partitionCounterStart := 0
	maxPart := 0
	for partRowID, partRow := range partitionRows {
		numOfColumns := partitionsPerRow[partRowID]
		//for every row in the original node matrix determine the cutoffs between the partition columns
		for _, graphRowID := range partRow {
			nodePartitionMatrix = append(nodePartitionMatrix, []int{})
			//nodes in row / columns rounded up
			nodesInThisGraphRow := len(graph.NodeMatrix[graphRowID])
			//start id of this clumn
			startID := 0
			//end id of this column (element of this id not included)
			cutOffID := 0
			//current column we are adding
			colCtr := 0
			for colCtr < numOfColumns {
				cutOffID = int(math.Round(float64((colCtr+1.0)*nodesInThisGraphRow) / float64(numOfColumns)))
				startID = int(math.Round(float64(colCtr*nodesInThisGraphRow) / float64(numOfColumns)))
				// 3 columns for 10 nodes should go -> 4, 4, 2
				for i := startID; i < cutOffID && i < nodesInThisGraphRow; i++ {
					//only add partition of node if node exists
					nodePartitionMatrix[graphRowID] = append(nodePartitionMatrix[graphRowID], partitionCounterStart+colCtr)
					if partitionCounterStart+colCtr > maxPart {
						maxPart = partitionCounterStart + colCtr
					}
				}
				colCtr++
				//determine each cut off
			}
			//fill rows (to the right with highest partition in this row) if necessary (e.g. if round down for nodespercolumn)
			//this guarantees that every node on the rightmost side has a valid partition
			for len(nodePartitionMatrix[graphRowID]) < len(graph.NodeMatrix[graphRowID]) {
				nodePartitionMatrix[graphRowID] = append(nodePartitionMatrix[graphRowID], partitionCounterStart+colCtr-1)
				if partitionCounterStart+colCtr-1 > maxPart {
					maxPart = partitionCounterStart + colCtr - 1
				}
			}
			numberOfPartitions = partitionCounterStart + colCtr
		}
		partitionCounterStart += numOfColumns
	}
	// DIVIDE PARTITION ROWS INTO PARTITIONS

	fmt.Printf("Rough nodes per partition: %d\n---\n", nodesPerRow[numOfRows/2]/numOfPoleRowPartitions)
	fmt.Printf("node matrix rows: %d\npartition matrix rows: %d\n", len(graph.NodeMatrix), len(nodePartitionMatrix))
	fmt.Printf("number of partitions total: %d\n", numberOfPartitions)
	fmt.Printf("rough height of partitions: %d\n", rowCount/numOfRows)
	numOfMiddleColumns := nodePartitionMatrix[len(nodePartitionMatrix)/2][len(nodePartitionMatrix[len(nodePartitionMatrix)/2])-1] - nodePartitionMatrix[len(nodePartitionMatrix)/2][0]

	fmt.Printf("rough width of partitions: %d\n", (nodesPerRow[numOfRows/2]/numOfMiddleColumns)/(rowCount/numOfRows))
	return nodePartitionMatrix, numberOfPartitions
}

// returns all boundary nodes of all partitions
func getBoundaryNodeIDS(graph Graph, nodePartitionMatrix [][]int) []int {
	boundaryNodeIDS := []int{}
	for rowID, row := range graph.NodeMatrix {
		for colID, nodeID := range row {
			currentNodePartition := nodePartitionMatrix[rowID][colID]
			neighList := GetGraphNeighbors(graph.Targets, graph.Offsets, graph.Weights, nodeID)
			shouldAdd := false
			for _, idAndDistance := range neighList {
				row, col := GetNodeMatrixIndex(idAndDistance[0], graph)
				if nodePartitionMatrix[row][col] != currentNodePartition {
					shouldAdd = true
				}
			}
			if shouldAdd {
				boundaryNodeIDS = append(boundaryNodeIDS, nodeID)
			}
		}
	}
	return boundaryNodeIDS
}

// // returns all neighbro node IDs connected to the input node
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

		neighborIDList = append(neighborIDList, []int{destinations[i], weights[i]})

	}
	return neighborIDList
}

// dijkstra for a single source node (a boundary node) to all other nodes - adds all arc flags to this source node
func singleSourceArcFlagPreprocess(graph Graph, sourceID int, arcFlags [][]bool, numOfPartitions int, nodePartitionMatrix [][]int, test bool) {
	hasBeenPopped := make([]bool, len(graph.Nodes))
	var distance []int
	//here prevEdges means previous edgeID
	var prevEdges []int
	// for each node, we save all partitions that were "touched" >before< opening this node -> these are the flags of the backwards edge, as soon as target node of edge (id of the prev node) is opened
	var partitionsVisited [][]bool
	nodesPoppedCounter := 0
	//priority queue datastructure (see priority_queue.go)
	var prioQ = make(PriorityQueue, 1)

	for i := 0; i < len(graph.Nodes); i++ {
		prevEdges = append(prevEdges, -1)
		distance = append(distance, 50000000)
		partitionsVisited = append(partitionsVisited, make([]bool, numOfPartitions))
	}

	row, col := GetNodeMatrixIndex(sourceID, graph)
	sourceNodesPartition := nodePartitionMatrix[row][col]
	partitionsVisited[sourceID][sourceNodesPartition] = true

	distance[sourceID] = 0
	prioQ[0] = &Item{value: sourceID, priority: distance[sourceID]}
	heap.Init(&prioQ)
	for {
		//gets "best" next node
		if prioQ.Len() < 1 {
			break
		}
		node := heap.Pop(&prioQ).(*Item)
		thisNodeID := node.value
		if hasBeenPopped[thisNodeID] {
			continue
		}
		if distance[node.value] >= 50000000 {
			break
		}
		hasBeenPopped[thisNodeID] = true
		if prevEdges[thisNodeID] >= 0 {
			row, col := GetNodeMatrixIndex(thisNodeID, graph)
			thisNodesPartition := nodePartitionMatrix[row][col]
			// if prev edge is not -1
			// node popped -> prev is final -> save visited notes of prev + current one
			//add all previously visited flags to this node
			prevNodeID := graph.Sources[prevEdges[thisNodeID]]
			// -> add (stored) arcflags for reverseEdgeOf(prev)
			reverseEdgeID := getReverseEdgeID(graph, prevEdges[thisNodeID])
			//partitionsVisited[thisNodeID] = partitionsVisited[prevNodeID]  <-- bad!
			copy(partitionsVisited[thisNodeID], partitionsVisited[prevNodeID]) // <-- good!
			for pID, f := range partitionsVisited[thisNodeID] {
				arcFlags[reverseEdgeID][pID] = (f || arcFlags[reverseEdgeID][pID])
			}
			//below version has far less false "true"s but is also very inefficient and for some reason not exactly correct
			partitionsVisited[thisNodeID][thisNodesPartition] = true
		}

		nodesPoppedCounter++
		// if we are at the destination then we break!

		// gets all neighbor/connected nodes
		neighbors := getArcFlagPreProcessNeighbors(graph.Targets, graph.Offsets, graph.Weights, node.value, hasBeenPopped)
		for _, neighbor := range neighbors {
			alt := distance[node.value] + neighbor[1]
			// neighbor [0] is target node ID, neighbor[2] is ID of edge which goes to this node
			if alt < distance[neighbor[0]] {
				distance[neighbor[0]] = alt
				prevEdges[neighbor[0]] = neighbor[2]
				//just re-queue items with better value instead of updating it
				heap.Push(&prioQ, &Item{value: neighbor[0], priority: alt})
			}
		}
	}
	if test {
		defer wg.Done()
	}
}

// returns the neighbors of a node by their [nodeID, distance, edgeID which leads to new node]
func getArcFlagPreProcessNeighbors(destinations []int, offsets []int, weights []int, nodeID int, hasBeenPopped []bool) [][]int {
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
		if !hasBeenPopped[destinations[i]] {
			neighborIDList = append(neighborIDList, []int{destinations[i], weights[i], i})
		}
	}
	return neighborIDList
}

func printPreProcessProgress(current int, max int) {
	progress := float64(current) / float64(max)
	fmt.Printf("Arc flag pre-processing | Progress: %2.2f%s%d%s%d boundary nodes completed (%.3fs)\n\r", 100*progress, "% - ", current, " out of ", max, time.Since(processTimer).Seconds())

	processTimer = time.Now()
}
