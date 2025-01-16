package main

import (
	"fmt"
	"math"
)

var graph map[string]map[string]float64
var costs map[string]float64
var parents map[string]string
var processed []string

func main() {

	var width, height int
	fmt.Scan(&width, &height)

	// Инициализация графа
	graph = make(map[string]map[string]float64)
	costs = make(map[string]float64)
	parents = make(map[string]string)

	// Заполнение графа и стоимости из ввода
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			var weight float64
			fmt.Scan(&weight)

			if weight > 0 { // Если вес больше 0, добавляем узел в граф
				node := fmt.Sprintf("%d,%d", i, j)
				if _, ok := graph[node]; !ok {
					graph[node] = make(map[string]float64)
				}

				// Добавляем соседние узлы (вверх, вниз, влево, вправо)
				addNeighbors(i, j, width, height, weight)
			}
			// Инициализация стоимости достижения узла
			costs[fmt.Sprintf("%d,%d", i, j)] = math.Inf(1)
		}
	}

	// Ввод координат старта и финиша
	var startRow, startCol, endRow, endCol int
	fmt.Scan(&startRow, &startCol)
	fmt.Scan(&endRow, &endCol) // Считываем endRow и endCol здесь

	startNode := fmt.Sprintf("%d,%d", startRow, startCol) // Определяем стартовый узел
	endNode := fmt.Sprintf("%d,%d", endRow, endCol)       // Определяем конечный узел

	// Установка начальной стоимости
	costs[startNode] = 0
	parents[startNode] = ""

	// Поиск кратчайшего пути
	node := findLowestCostNode(costs)

	// Поиск кратчайшего пути
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

	fmt.Printf("Cost from %s to all nodes: %v\n", startNode, costs)
	fmt.Printf("Shortest path to finish: ")

	// Восстанавливает путь от конечного узла к стартовому
	path := []string{}
	currentNode := endNode
	for currentNode != "" {
		path = append([]string{currentNode}, path...)
		currentNode = parents[currentNode]
	}

	fmt.Println(path)
}

// addNeighbors добавляет соседей для текущего узла в зависимости от их расположения в сетке
func addNeighbors(row, col, width, height int, weight float64) {
	node := fmt.Sprintf("%d,%d", row, col)
	if row > 0 { // вверх
		addEdge(node, fmt.Sprintf("%d,%d", row-1, col), weight)
	}
	if row < height-1 { // вниз
		addEdge(node, fmt.Sprintf("%d,%d", row+1, col), weight)
	}
	if col > 0 { // влево
		addEdge(node, fmt.Sprintf("%d,%d", row, col-1), weight)
	}
	if col < width-1 { // вправо
		addEdge(node, fmt.Sprintf("%d,%d", row, col+1), weight)
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
