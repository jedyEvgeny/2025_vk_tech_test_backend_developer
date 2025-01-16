package main

import (
	"fmt"
	"math"
)

var (
	graph     map[string]map[string]float64
	costs     map[string]float64
	parents   map[string]string
	processed []string
)

func main() {
	initGraph()
	fillGraph()

	startNode, endNode := findCoordinatesStartAndFinishNodes()
	fillInitialCost(startNode)

	node := findLowestCostNode(costs)
	findShortestPath(node)
	path := retrievePath(endNode)

	printPath(path)
}

// initGraph инициализирует граф
func initGraph() {
	graph = make(map[string]map[string]float64)
	costs = make(map[string]float64)
	parents = make(map[string]string)
}

// fillGraph заполняет граф и стоимость из ввода
func fillGraph() {
	var width, height int
	fmt.Scan(&width, &height)

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			var weight float64
			fmt.Scan(&weight)

			if weight > 0 { // Если вес больше 0, добавляем узел в граф
				node := fmt.Sprintf("%d %d", i, j)
				if _, ok := graph[node]; !ok {
					graph[node] = make(map[string]float64)
				}

				// Добавляем соседние узлы (вверх, вниз, влево, вправо)
				addNeighbors(i, j, width, height, weight)
			}
			// Инициализация стоимости достижения узла
			costs[fmt.Sprintf("%d %d", i, j)] = math.Inf(1)
		}
	}
}

// findCoordinatesStartAndFinishNodes формирует координаты начального
// и конечного узлов графа
func findCoordinatesStartAndFinishNodes() (string, string) {
	var startRow, startCol, endRow, endCol int

	fmt.Scan(&startRow, &startCol)
	fmt.Scan(&endRow, &endCol)

	startNode := fmt.Sprintf("%d %d", startRow, startCol)
	endNode := fmt.Sprintf("%d %d", endRow, endCol)

	return startNode, endNode
}

// fillInitialCost устанавливает начальный вес
func fillInitialCost(startNode string) {
	costs[startNode] = 0
	parents[startNode] = ""
}

// addNeighbors добавляет соседей для текущего узла в зависимости от их расположения в сетке
func addNeighbors(row, col, width, height int, weight float64) {
	node := fmt.Sprintf("%d %d", row, col)
	if row > 0 { // вверх
		addEdge(node, fmt.Sprintf("%d %d", row-1, col), weight)
	}
	if row < height-1 { // вниз
		addEdge(node, fmt.Sprintf("%d %d", row+1, col), weight)
	}
	if col > 0 { // влево
		addEdge(node, fmt.Sprintf("%d %d", row, col-1), weight)
	}
	if col < width-1 { // вправо
		addEdge(node, fmt.Sprintf("%d %d", row, col+1), weight)
	}
}

// addEdge добавляет ребро в граф
func addEdge(from, to string, weight float64) {
	if _, ok := graph[from]; !ok {
		graph[from] = make(map[string]float64)
	}
	graph[from][to] = weight
}

// findLowestCostNode находит узел с наименьшей стоимостью
func findLowestCostNode(costs map[string]float64) string {
	lowestCost := math.Inf(1)
	lowestCostNode := ""

	for node := range costs {
		cost := costs[node]
		if cost < lowestCost && !isContains(processed, node) {
			lowestCost = cost
			lowestCostNode = node
		}
	}

	return lowestCostNode
}

// isContains проверяет, есть ли элемент в срезе
func isContains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}

	return false
}

// findShortestPath находит кратчайший путь
func findShortestPath(node string) {
	for node != "" {
		cost := costs[node]
		neighbors := graph[node]

		// Обработка соседних узлов
		for neighbor := range neighbors {
			newCost := cost + neighbors[neighbor]
			if costs[neighbor] > newCost {
				costs[neighbor] = newCost
				parents[neighbor] = node
			}
		}

		// Помечаем узел как обработанный
		processed = append(processed, node)

		node = findLowestCostNode(costs)
	}
}

// retrievePath восстанавливает путь от конечного узла к стартовому
func retrievePath(endNode string) []string {
	path := []string{}
	currentNode := endNode
	for currentNode != "" {
		path = append([]string{currentNode}, path...)
		currentNode = parents[currentNode]
	}

	return path
}

func printPath(path []string) {
	for idx, coordinates := range path {
		fmt.Println(coordinates)
		if idx == len(path)-1 {
			fmt.Println(".")
		}
	}
}
