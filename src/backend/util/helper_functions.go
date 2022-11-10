package util

import "fmt"

//checks if a point, given coordinates, is on land (false) or in water (true)
func IsPointInWater(g graph, p point) bool {
	//edges := g.edges
	//do point in polygon test for all way-polygons
	return false
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
