package util

import (
	"fmt"
	"math"
)

func CreatePartitions(graph Graph, numOfRows int, numOfPoleRowPartitions int) [][]int {
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
	fmt.Printf("Total rows in graph matrix: %d - totalRows in partition row matrix: %d \n", rowCount, currentGRowID)
	for id, row := range partitionRows {
		//calc how many partitions per row
		partitionsPerRow = append(partitionsPerRow, int(math.Ceil(float64(nodesPerRow[id])/(float64(nodesPerRow[0])/float64(numOfPoleRowPartitions)))))
		fmt.Printf("partition row %d: %d rows - %d nodes - %d partitions\n", id, len(row), nodesPerRow[id], partitionsPerRow[id])
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
				}

				colCtr++
				//determine each cut off
			}
			//fill rows if necessary (e.g. if round down for nodespercolumn)
			for len(nodePartitionMatrix[graphRowID]) < len(graph.NodeMatrix[graphRowID]) {
				nodePartitionMatrix[graphRowID] = append(nodePartitionMatrix[graphRowID], partitionCounterStart+colCtr-1)
			}
		}
		partitionCounterStart += numOfColumns
	}
	// DIVIDE PARTITION ROWS INTO PARTITIONS

	fmt.Printf("node matrix rows: %d\npartition matrix rows: %d\n", len(graph.NodeMatrix), len(nodePartitionMatrix))
	for rowID, row := range graph.NodeMatrix {
		fmt.Printf("diff: %d\n", len(row)-len(nodePartitionMatrix[rowID]))
		fmt.Println(nodePartitionMatrix[rowID])
	}
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
	return nodePartitionMatrix
}
