package util

import "fmt"

//checks if a point, given coordinates, is on land (false) or in water (true)
func IsPointInWater(p point) bool {
	//input should be: numberOfEdgesInTheWayToNorthpole, []edge
	//SQL command for relevant edges
	//number, questionableEdges := GetTestSQLResult(p)
	// for every edge e call isEdgeInTheWay(p, e)
	// if true -> number += 1

	return false
}

func isEdgeInTheWay(p point, e edge) bool {
	//check if edge is in the way
	//point : {lat, lon}
	//p.lat
	//edge : {point1, point2}
	//edge.point.lat
	return true
}

//takes a whole basic format map as input, checks if one way starts where another ends -> merges them
func MergeWays(inputMap basic) int {
	mergeCounter := 0
	ways := inputMap.ways
	for firstWayID, way := range ways {
		endWay, exists := ways[way.lastNodeID]
		if exists && firstWayID != way.lastNodeID {
			newWay := mergeTwoWays(way, endWay)
			inputMap.ways[firstWayID] = newWay
			delete(inputMap.ways, endWay.nodes[0])
			mergeCounter++
		}
	}
	fmt.Println("--------------attempt to merge touching ways--------------")
	fmt.Printf("ways merged: %d\n", mergeCounter)
	fmt.Printf("ways left: %d\n", len(inputMap.ways))
	return mergeCounter
}

// takes two connected ways and merges them into one
func mergeTwoWays(startWay way, endWay way) way {
	newNodes := append(startWay.nodes, endWay.nodes[1:]...)
	newWay := way{nodes: newNodes, lastNodeID: endWay.lastNodeID}
	return newWay
}
