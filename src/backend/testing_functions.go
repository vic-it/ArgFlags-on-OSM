package backend

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var dijkstaProgress float64 = 0
var dijkstraResult string = ""
var arcFlagProgress float64 = 0
var arcFlagResult string = ""
var shouldAbort = false

func GetTestStatusAsString() string {
	output := fmt.Sprintf("%.1f-%.1f-%s-%s", dijkstaProgress, arcFlagProgress, dijkstraResult, arcFlagResult)
	return output
}

func TestAlgorithms(graph Graph, arcData ArcData, numberOfPaths int, nodePartitionList []int) {
	shouldAbort = false
	dijkstaProgress = 0.0
	arcFlagProgress = 0.0
	dijkstraResult = ""
	arcFlagResult = ""
	paths := generateTestPaths(graph, numberOfPaths)
	testDijkstra(graph, numberOfPaths, paths)
	testArcFlagDijkstra(graph, numberOfPaths, paths, arcData, nodePartitionList)
}

func AbortTests() {
	shouldAbort = true
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
		if shouldAbort {
			println("Dijkstra Aborted")
			return
		}
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
		dijkstaProgress = 100.0 * float64(i) / float64(n-1)
		if i%100 == 0 {
			printTestProgress("Dijkstra", i, n)
		}
	}
	//TEST: 1000 runs with average initialization time of: xxx and average search time of: xxx and longest total search of: xxx
	//print results
	stats := fmt.Sprintf("Total test time for %d routes: %.1fs (avg: %.1fms)\n", n, time.Since(totalTime).Seconds(), 1000*time.Since(totalTime).Seconds()/float64(n))
	stats += fmt.Sprintf("Average initialization time: %.1fms (max: %.1fms)\n", 1000*totalInitTime/float64(n), 1000*maxInitTime)
	stats += fmt.Sprintf("Average search time: %.1fms (max: %.1fms)\n", 1000*totalSearchTime/float64(n), 1000*maxSearchTime)
	stats += fmt.Sprintf("Average heap pops: %d pops (max: %d pops)\n", totalNodesPopped/n, maxNodesPopped)
	stats += fmt.Sprintf("Average distance: %dkm (max: %dkm)\n", totalDistance/(n-fails), maxDistance)
	//stats += fmt.Sprintf("Average path length: %d nodes (max: %d nodes)\n", totalPathLength/(n-fails), maxPathLength)
	stats += fmt.Sprintf("Number of routes with no viable path: %d\n", fails)
	println("-------")
	println("NORMAL DIJKSTRA RESULTS:")
	fmt.Printf(stats)
	dijkstraResult = strings.ReplaceAll(stats, "\n", "<br>")
	println("-------")
}

func testArcFlagDijkstra(graph Graph, n int, paths [][]int, arcData ArcData, nodePartitionList []int) {
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
		if shouldAbort {
			println("Arc flags Aborted")
			return
		}
		dist, path, initTime, searchTime, nodesPopped := CalculateArcFlagDijkstra(graph, paths[i][0], paths[i][1], arcData, nodePartitionList)
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

		arcFlagProgress = 100.0 * float64(i) / float64(n-1)
		if i%100 == 0 {
			printTestProgress("Arc Flags", i, n)
		}
	}
	//TEST: 1000 runs with average initialization time of: xxx and average search time of: xxx and longest total search of: xxx
	//print results
	stats := fmt.Sprintf("Total test time for %d routes: %.1fs (avg: %.1fms)\n", n, time.Since(totalTime).Seconds(), 1000*time.Since(totalTime).Seconds()/float64(n))
	stats += fmt.Sprintf("Average initialization time: %.1fms (max: %.1fms)\n", 1000*totalInitTime/float64(n), 1000*maxInitTime)
	stats += fmt.Sprintf("Average search time: %.1fms (max: %.1fms)\n", 1000*totalSearchTime/float64(n), 1000*maxSearchTime)
	stats += fmt.Sprintf("Average heap pops: %d pops (max: %d pops)\n", totalNodesPopped/n, maxNodesPopped)
	stats += fmt.Sprintf("Average distance: %dkm (max: %dkm)\n", totalDistance/(n-fails), maxDistance)
	//stats += fmt.Sprintf("Average path length: %d nodes (max: %d nodes)\n", totalPathLength/(n-fails), maxPathLength)
	stats += fmt.Sprintf("Number of routes with no viable path: %d\n", fails)
	println("-------")
	println("ARC FLAG RESULTS:")
	fmt.Printf(stats)

	arcFlagResult = strings.ReplaceAll(stats, "\n", "<br>")
	println("-------")
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
