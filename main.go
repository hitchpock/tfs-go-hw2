package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Структура подписчиков в User.
type Subscriber struct {
	Email     string `json:"Email"`
	CreatedAt string `json:"Created_at"`
}

func (s *Subscriber) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	}{
		Email:     s.Email,
		CreatedAt: s.CreatedAt,
	})
}

// Стуктура User входного файла.
type User struct {
	Nick        string       `json:"Nick"`
	Email       string       `json:"Email"`
	CreatedAt   string       `json:"Created_at"`
	Subscribers []Subscriber `json:"Subscribers"`
}

// Структура для выходного файла.
type ResultStruct struct {
	ID   int          `json:"id"`
	From string       `json:"from"`
	To   string       `json:"to"`
	Path []Subscriber `json:"path,omitempty"`
}

func main() {
	dataCSV, err := readCSVFile("input.csv")
	if err != nil {
		log.Fatal("func readCSVFile crash: ", err)
	}

	userJSON, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatal("unable to read file 'users.json'")
	}

	users := []User{}

	err = json.Unmarshal(userJSON, &users)
	if err != nil {
		log.Fatal("unable to unmarshal userJSON: ", err)
	}

	finalList := CreateResult(users, dataCSV)

	err = WriteToJSON(finalList)
	if err != nil {
		log.Fatal("WriteToJSON crash: ", err)
	}
}

// Создание "социального графа".
func GraphFromSliceSub(users []User) map[string][]Subscriber {
	graph := make(map[string][]Subscriber)

	for _, user := range users {
		for _, sub := range user.Subscribers {
			graph[sub.Email] = append(graph[sub.Email], Subscriber{user.Email, user.CreatedAt})
		}
	}

	return graph
}

// Нахождение связи между двумя пользователями.
// Основан на алгоритме https://github.com/TheAlgorithms/Python/blob/master/graphs/bfs_shortest_path.py
func BFSSub(graph map[string][]Subscriber, start Subscriber, goal string) []Subscriber {
	if start.Email == goal {
		return nil
	}

	var explored []Subscriber

	queue := [][]Subscriber{{start}}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		node := path[len(path)-1]

		if !inAlt(explored, node) {
			neighbors := graph[node.Email]

			for _, neighbor := range neighbors {
				newPath := make([]Subscriber, len(path))
				copy(newPath, path)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)

				if neighbor.Email == goal {
					return newPath[1 : len(newPath)-1]
				}
			}

			explored = append(explored, node)
		}
	}

	return nil
}

// Встпомогательная функция bfs.
func inAlt(s []Subscriber, str Subscriber) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == str {
			return true
		}
	}

	return false
}

// Извлечение информации из csv файла.
func readCSVFile(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %s", err)
	}
	defer file.Close()

	// Парсинг csv файлов.
	fileReader := csv.NewReader(file)

	dataCSV, err := fileReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %s", err)
	}

	return dataCSV, err
}

// Нахождение связей между всеми парами пользователей из списка.
func CreateResult(users []User, dataCSV [][]string) []ResultStruct {
	gr := GraphFromSliceSub(users)
	finalList := make([]ResultStruct, 0, len(dataCSV))

	for _, line := range dataCSV {
		start := line[0]
		goal := line[1]
		count := 1

		for _, user := range users {
			if user.Email == start {
				sub := Subscriber{user.Email, user.CreatedAt}
				path := BFSSub(gr, sub, goal)
				note := ResultStruct{ID: count, From: start, To: goal, Path: path}
				finalList = append(finalList, note)
				count++
			}
		}
	}

	return finalList
}

// Запись итогового списка в json файл.
func WriteToJSON(finalList []ResultStruct) error {
	file, err := os.Create("result.json")
	if err != nil {
		return fmt.Errorf("unable to create file: %s", err)
	}
	defer file.Close()

	resultEncoder := json.NewEncoder(file)
	resultEncoder.SetIndent("   ", "\t")

	err = resultEncoder.Encode(finalList)
	if err != nil {
		return fmt.Errorf("unable to Encode: %s", err)
	}

	return nil
}
