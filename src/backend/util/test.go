package util

import (
	"fmt"
	"math/rand"
	"time"
)

func TestAlgorithms(graph Graph, nodePartitionMatrix [][]int, arcFlags [][]bool, numberOfPaths int) {
	paths := generateTestPaths(graph, numberOfPaths)
	testDijkstra(graph, numberOfPaths, paths)
	testArcFlagDijkstra(graph, numberOfPaths, paths, nodePartitionMatrix, arcFlags)
}
func testDijkstra(graph Graph, n int, paths [][]int) {
	fmt.Printf("Starting Dijkstra test for %d routes...\n", n)
	totalTime := time.Now()
	fails := 0
	totalDistance := 0
	totalInitTime := 0.0
	totalSearchTime := 0.0
	totalNodesPopped := 0
	totalPathLength := 0
	maxDistance := 0
	maxInitTime := 0.0
	maxSearchTime := 0.0
	maxNodesPopped := 0
	maxPathLength := 0
	for i := 0; i < n; i++ {
		dist, path, initTime, searchTime, nodesPopped := CalculateDijkstra(graph, paths[i][0], paths[i][1])
		if dist >= 0 {
			totalDistance += dist
			totalPathLength += len(path)
			if len(path) > maxPathLength {
				maxPathLength = len(path)
			}
		} else {
			fails++
		}
		totalInitTime += initTime
		totalSearchTime += searchTime
		totalNodesPopped += nodesPopped
		if dist > maxDistance {
			maxDistance = dist
		}
		if initTime > maxInitTime {
			maxInitTime = initTime
		}
		if searchTime > maxSearchTime {
			maxSearchTime = searchTime
		}
		if nodesPopped > maxNodesPopped {
			maxNodesPopped = nodesPopped
		}

		if i%1000 == 0 {
			printTestProgress("Dijkstra", i, n)
		}
	}
	//TEST: 1000 runs with average initialization time of: xxx and average search time of: xxx and longest total search of: xxx
	//print results
	fmt.Printf("-------\nTotal test time for %d routes: %.3fs (avg: %.3fs)\n", n, time.Since(totalTime).Seconds(), time.Since(totalTime).Seconds()/float64(n))
	fmt.Printf("Average initialization time: %.3fs (max: %.3fs)\n", totalInitTime/float64(n), maxInitTime)
	fmt.Printf("Average search time: %.3fs (max: %.3fs)\n", totalSearchTime/float64(n), maxSearchTime)
	fmt.Printf("Average heap pops: %d pops (max: %d pops)\n", totalNodesPopped/n, maxNodesPopped)
	fmt.Printf("Average distance: %dkm (max: %dkm)\n", totalDistance/(n-fails), maxDistance)
	fmt.Printf("Average path length: %d nodes (max: %d nodes)\n", totalPathLength/(n-fails), maxPathLength)
	fmt.Printf("Number of routes with no viable path: %d\n-------\n", fails)
}

func testArcFlagDijkstra(graph Graph, n int, paths [][]int, nodePartitionMatrix [][]int, arcFlags [][]bool) {
	fmt.Printf("Starting Arc Flags test for %d routes...\n", n)
	totalTime := time.Now()
	fails := 0
	totalDistance := 0
	totalInitTime := 0.0
	totalSearchTime := 0.0
	totalNodesPopped := 0
	totalPathLength := 0
	maxDistance := 0
	maxInitTime := 0.0
	maxSearchTime := 0.0
	maxNodesPopped := 0
	maxPathLength := 0
	for i := 0; i < n; i++ {
		dist, path, initTime, searchTime, nodesPopped := CalculateArcFlagDijkstra(graph, paths[i][0], paths[i][1], arcFlags, nodePartitionMatrix)
		if dist >= 0 {
			totalDistance += dist
			totalPathLength += len(path)
			if len(path) > maxPathLength {
				maxPathLength = len(path)
			}
		} else {
			fails++
		}
		totalInitTime += initTime
		totalSearchTime += searchTime
		totalNodesPopped += nodesPopped
		if dist > maxDistance {
			maxDistance = dist
		}
		if initTime > maxInitTime {
			maxInitTime = initTime
		}
		if searchTime > maxSearchTime {
			maxSearchTime = searchTime
		}
		if nodesPopped > maxNodesPopped {
			maxNodesPopped = nodesPopped
		}

		if i%100 == 0 {
			printTestProgress("Arc Flags", i, n)
		}
	}
	//TEST: 1000 runs with average initialization time of: xxx and average search time of: xxx and longest total search of: xxx
	//print results
	fmt.Printf("-------\nTotal test time for %d routes: %.3fs (avg: %.3fs)\n", n, time.Since(totalTime).Seconds(), time.Since(totalTime).Seconds()/float64(n))
	fmt.Printf("Average initialization time: %.3fs (max: %.3fs)\n", totalInitTime/float64(n), maxInitTime)
	fmt.Printf("Average search time: %.3fs (max: %.3fs)\n", totalSearchTime/float64(n), maxSearchTime)
	fmt.Printf("Average heap pops: %d pops (max: %d pops)\n", totalNodesPopped/n, maxNodesPopped)
	fmt.Printf("Average distance: %dkm (max: %dkm)\n", totalDistance/(n-fails), maxDistance)
	fmt.Printf("Average path length: %d nodes (max: %d nodes)\n", totalPathLength/(n-fails), maxPathLength)
	fmt.Printf("Number of routes with no viable path: %d\n-------\n", fails)
}

func generateTestPaths(graph Graph, numberOfPaths int) [][]int {
	rand.Seed(time.Now().UnixNano())
	output := [][]int{}
	numOfNodes := len(graph.Nodes)
	for len(output) < numberOfPaths {
		id1 := rand.Intn(numOfNodes)
		id2 := rand.Intn(numOfNodes)
		row1, col1 := GetNodeMatrixIndex(id1, graph)
		row2, col2 := GetNodeMatrixIndex(id2, graph)
		if id1 != id2 {
			if graph.NodeInWaterMatrix[row1][col1] && graph.NodeInWaterMatrix[row2][col2] {
				output = append(output, []int{id1, id2})
			}
		}
	}
	return output
}

func GetNodeMatrixIndex(n int, graph Graph) (int, int) {
	rowIdx := 0
	in := n
	matrix := graph.NodeMatrix
	for n >= len(matrix[rowIdx]) {
		n -= len(matrix[rowIdx])
		rowIdx++
	}
	if matrix[rowIdx][n] != in {
		fmt.Printf("FAULTY NODE ID IN MATRIX")
	}
	// fmt.Printf("Searched for ID: %d - found ID: %d at [%d][%d]\n", in, matrix[rowIdx][n], rowIdx, n)
	return rowIdx, n
}

func printTestProgress(algName string, current int, max int) {
	progress := float64(current) / float64(max)
	fmt.Printf("Running test for %s | Progress: %2.2f%s%d%s%d runs completed\n\r", algName, 100*progress, "% - ", current, " out of ", max)
}
