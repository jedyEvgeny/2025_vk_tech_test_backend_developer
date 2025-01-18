package main

import (
	"io"
	"os"
	"testing"
)

func TestFillGraph(t *testing.T) {
	tests := []struct {
		testName    string
		input       string
		expectError bool
	}{
		{
			testName:    "Корректный ввод",
			input:       "3 3\n1 2 3\n4 5 6\n7 8 9\n",
			expectError: false,
		},
		{
			testName:    "Ввод с символом не-цифрой",
			input:       "3 3\na 2 3\n4 5 6\n7 8 9\n",
			expectError: true,
		},
		{
			testName:    "Несоответствие матрицы в +",
			input:       "2 2\n1 2 3\n4 5\n",
			expectError: true,
		},
		{
			testName:    "Несоответствие матрицы в -",
			input:       "2 2\n1\n4 5\n",
			expectError: true,
		},
		{
			testName:    "Отрицательный вес",
			input:       "2 2\n-1 2\n3 4\n",
			expectError: true,
		},
		{
			testName:    "Только положительные веса",
			input:       "2 2\n1 2\n3 4\n",
			expectError: false,
		},
		{
			testName:    "Все веса отрицательные",
			input:       "2 2\n-1 -2\n-3 -4\n",
			expectError: true,
		},
		{
			testName:    "Отрицательный размер матрицы",
			input:       "-1 3\n1 2 3\n",
			expectError: true,
		},
		{
			testName:    "Только нули",
			input:       "2 2\n0 0\n0 0\n",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			restoreStdin := redirectStdin(tt.input)
			defer restoreStdin() // Восстанавливаем stdin после завершения теста

			initGraph()
			fillGraph()

			// Проверка на ошибки
			if tt.expectError && len(costs) > 0 {
				t.Errorf("Ожидалась ошибка, но ввод сохранён: %v", costs)
			}
		})
	}
}
func TestFindCoordinatesStartAndFinishNodes(t *testing.T) {
	tests := []struct {
		name          string
		coords        coordinates
		width         int
		height        int
		expectedStart string
		expectedEnd   string
		expectError   bool
	}{
		{
			name: "Корректные координаты",
			coords: coordinates{
				startRow: 0,
				startCol: 0,
				endRow:   2,
				endCol:   2,
			},
			width:         3,
			height:        3,
			expectedStart: "0 0",
			expectedEnd:   "2 2",
			expectError:   false,
		},
		{
			name: "Начальная координата столбца вне границ",
			coords: coordinates{
				startRow: 0,
				startCol: 3, // Некорректное значение
				endRow:   2,
				endCol:   2,
			},
			width:         3,
			height:        3,
			expectedStart: "",
			expectedEnd:   "",
			expectError:   true,
		},
		{
			name: "Конечная координата в +",
			coords: coordinates{
				startRow: 0,
				startCol: 0,
				endRow:   2,
				endCol:   3, // Некорректное значение
			},
			width:         3,
			height:        3,
			expectedStart: "0 0",
			expectedEnd:   "",
			expectError:   true,
		},
		{
			name: "Начальная координата в -",
			coords: coordinates{
				startRow: -1, // Некорректное значение
				startCol: 1,
				endRow:   2,
				endCol:   2,
			},
			width:         3,
			height:        3,
			expectedStart: "",
			expectedEnd:   "",
			expectError:   true,
		},
		{
			name: "Конечная координата строки вне границ (низ)",
			coords: coordinates{
				startRow: 0,
				startCol: 1,
				endRow:   3, // Некорректное значение
				endCol:   2,
			},
			width:         3,
			height:        3,
			expectedStart: "0 1",
			expectedEnd:   "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startNode, endNode, err := tt.coords.findStartAndFinishNodes(tt.width, tt.height)

			if tt.expectError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но её не было")
				}
			} else {
				if startNode != tt.expectedStart || endNode != tt.expectedEnd {
					t.Errorf("Ожидалось: %s %s, получено: %s %s",
						tt.expectedStart, tt.expectedEnd, startNode, endNode)
				}
			}
		})
	}
}

// redirectStdin перенаправляет стандартный ввод на входящие данные
func redirectStdin(input string) func() {
	reader, writer, _ := os.Pipe()
	os.Stdin = writer

	// Заполняем pipe входными данными
	go func() {
		defer writer.Close()
		io.WriteString(writer, input)
		writer.Close() // Закрываем writer после записи
	}()

	// Восстанавливаем стандартный ввод после завершения теста
	return func() {
		os.Stdin = reader
		reader.Close()
	}
}
