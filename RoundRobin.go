package main

import (
	"fmt"
	"net/http"
)

type Server struct {
	Ip    string
	URL   string
	Alive bool
}

func main() {
	var startingServerList []Server

	startingServerList = append(startingServerList, Server{"127.0.0.1", "http://localhost:5001/", false})
	startingServerList = append(startingServerList, Server{"127.0.0.1", "http://localhost:5002/", false})
	startingServerList = append(startingServerList, Server{"127.0.0.1", "http://localhost:5003/", false})
	startingServerList = append(startingServerList, Server{"127.0.0.1", "http://localhost:5004/", false})

	test(startingServerList)
}

func test(startingServerList []Server) {
	var list = startingServerList
	fmt.Println(list)
	for i := 0; i < 100000; i++ {
		list = loadBalancer(list)
		fmt.Println(list)
	}
}

func loadBalancer(startingServerList []Server) []Server {
	returningServerList := startingServerList
	for i := 0; i < len(startingServerList); i++ {
		startingServerList[i].Alive = checkIsAlive(startingServerList[i].URL)
		for j := 0; j < len(returningServerList); j++ {
			if !startingServerList[i].Alive && returningServerList[j].URL == startingServerList[i].URL {
				returningServerList = removeFromList(returningServerList, j)
			}
		}
	}
	returningServerList = balancingСyclicQueue(returningServerList)
	return returningServerList
}

//checkIsAlive проверка отвечает ли сервер
func checkIsAlive(url string) bool {
	//fmt.Println("Checking URL ", url)
	resp, err := http.Get(url)
	if err != nil {
		//fmt.Printf("Connection error. %s\n", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		//fmt.Printf("Error. http-status: %s\n", resp.StatusCode)
		return false
	}
	//fmt.Printf("Online. http-status: %d\n", resp.StatusCode)
	return true
}

func balancingСyclicQueue(list []Server) []Server {
	list = append(list, list[0])
	list = removeFromList(list, 0)
	return list
}

func removeFromList(list []Server, index int) []Server {
	list = list[:index+copy(list[index:], list[index+1:])]
	return list
}
