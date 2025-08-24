package data_structure

const SkiplistMaxLevel = 32

/*
	/level 2: span=2 | forward\ --------------------------------------> /span=0 | forward\ ----> NULL
	|level 1: span=1 | forward| --------> /span=1 | forward\ ---------> |span=0 | forward| ----> NULL
	|ele                      |           |ele             |            |ele             |
	|score                    |           |score           |            |score           |
	|backward                 | <-------- |backward        | <--------- |backward        |
	\node1                    /           \node2           /            \node3           /
*/
