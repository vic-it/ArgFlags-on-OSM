package util

type basic struct {
	nodes map[int64]node
	ways  map[int64]way
}

type tag struct {
	key   string
	value string
}

type node struct {
	lat  float64
	lon  float64
	tags []tag
}

type way struct {
	nodes []int64
	tags  []tag
}

type graph struct {
	nodes []node
	ways  []way
	//more stuff here
}
