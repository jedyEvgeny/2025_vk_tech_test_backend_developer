package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type coordinates struct {
	startRow int
	startCol int
	endRow   int
	endCol   int
}

var (
	graph     map[string]map[string]uint
	costs     map[string]uint
	parents   map[string]string
	processed []string
)

const (
	errInput = "недопустимые входящие данные: %v"
)

func main() {
	initGraph()
	width, height, err := fillGraph()
	if err != nil {
		log.Fatalf(errInput, err)
	}

	c := readCoordinates()
	startNode, endNode, err := c.findStartAndFinishNodes(width, height)
	if err != nil {
		log.Fatalf(errInput, err)
	}

	fillInitialCost(startNode)

	node := findLowestCostNode(costs)
	findShortestPath(node)
	path := retrievePath(endNode)

	printPath(path)
}

// initGraph инициализирует структуры данных
func initGraph() {
	graph = make(map[string]map[string]uint)
	costs = make(map[string]uint)
	parents = make(map[string]string)
}

// fillGraph заполняет граф из стандартного ввода
func fillGraph() (int, int, error) {
	var width, height int
	fmt.Scan(&width, &height)
	if width <= 0 || height <= 0 {
		return 0, 0, fmt.Errorf("размер матрицы должен быть положительным. Имеется: ширина=%d, высота=%d", width, height)
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			var weight uint
			fmt.Scan(&weight)

			if weight < 0 || weight > 9 {
				return 0, 0, fmt.Errorf("вес должен быть в диапазоне 0...9. Имеется: %d", weight)
			}

			if weight > 0 { // Если вес больше 0, добавляем узел в граф
				node := fmt.Sprintf("%d %d", i, j)
				if _, ok := graph[node]; !ok {
					graph[node] = make(map[string]uint)
				}

				// Добавляем соседние узлы (вверх, вниз, влево, вправо)
				addNeighbors(i, j, width, height, weight)
			}
			// Инициализация стоимости достижения узла
			costs[fmt.Sprintf("%d %d", i, j)] = math.MaxUint
		}
	}

	return width, height, nil
}

func readCoordinates() *coordinates {
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		log.Fatal(scanner.Err())
	}

	input := scanner.Text()

	elems := strings.Split(input, " ")
	if len(elems) != 4 {
		log.Fatalf("Ожидалось 4 координаты, имеется: %d. Также проверьте кол-во элементов в последней строке матрицы",
			len(elems))
	}

	c := &coordinates{}
	for idx, el := range elems {
		num, err := strconv.Atoi(el)
		if err != nil {
			log.Fatalf("%d-я координата не число, а '%s': %v",
				idx+1, el, err)
		}
		switch idx {
		case 0:
			c.startRow = num
		case 1:
			c.startCol = num
		case 2:
			c.endRow = num
		case 3:
			c.endCol = num
		}
	}

	return c
}

// findStartAndFinishNodes формирует координаты начального
// и конечного узлов графа
func (c *coordinates) findStartAndFinishNodes(width, height int) (string, string, error) {
	if c.startCol < 0 || c.startCol >= width {
		return "", "", fmt.Errorf("начальная координата столбца '%d' вне ширины матрицы: '%d'", c.startCol, width)
	}
	if c.endCol < 0 || c.endCol >= width {
		return "", "", fmt.Errorf("конечная координата столбца '%d' вне ширины матрицы: '%d'", c.endCol, width)
	}
	if c.startRow < 0 || c.startRow >= height {
		return "", "", fmt.Errorf("начальная координата строки '%d' вне высоты матрицы: '%d'", c.startRow, height)
	}
	if c.endRow < 0 || c.endRow >= height {
		return "", "", fmt.Errorf("конечная координата строки '%d' вне высоты матрицы: '%d'", c.endRow, height)
	}

	startNode := fmt.Sprintf("%d %d", c.startRow, c.startCol)
	endNode := fmt.Sprintf("%d %d", c.endRow, c.endCol)

	return startNode, endNode, nil
}

// fillInitialCost устанавливает начальный вес
func fillInitialCost(startNode string) {
	costs[startNode] = 0
	parents[startNode] = ""
}

// addNeighbors добавляет соседей для текущего узла в зависимости от их расположения в сетке
func addNeighbors(row, col, width, height int, weight uint) {
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
func addEdge(from, to string, weight uint) {
	if _, ok := graph[from]; !ok {
		graph[from] = make(map[string]uint)
	}
	graph[from][to] = weight
}

// findLowestCostNode находит узел с наименьшей стоимостью
func findLowestCostNode(costs map[string]uint) string {
	var lowestCost uint = math.MaxUint
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

// printPath выводит в терминал кратчайший маршрут между узлами
func printPath(path []string) {
	for idx, coordinates := range path {
		fmt.Println(coordinates)
		if idx == len(path)-1 {
			fmt.Println(".")
		}
	}
}
