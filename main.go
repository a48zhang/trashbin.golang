package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const MAXGRT = 50

var (
	cookie = "" //
	url    = "http://kjyy.ccnu.edu.cn/ClientWeb/pro/ajax/data/searchAccount.aspx?type=logonname&ReservaApply=ReservaApply&term="
	cli    http.Client
	task   chan string
	wg     sync.WaitGroup
	data   cmap
)

type cmap struct {
	l    sync.RWMutex
	data map[string]string
}

type Info struct {
	Id          string `json:"id"`
	Pid         string `json:"Pid"`
	Name        string `json:"name"`
	Label       string `json:"label"`
	SzLogonName string `json:"szLogonName"`
	SzHandPhone string `json:"szHandPhone"`
	SzTel       string `json:"szTel"`
	SzEmail     string `json:"szEmail"`
}

func main() {
	data.data = make(map[string]string)
	var begin, end int
	cli = http.Client{}
	
	fmt.Println("from? to?")
	fmt.Scanf("%d %d", &begin, &end)

	task = make(chan string, 1000)
	go func(){for begin <= end {
		task <- strconv.Itoa(begin)
		begin++
	}}()

	wg.Add(MAXGRT)
	for i := 1; i <= MAXGRT; i++ {
		worm()
	}
	close(task)
	wg.Wait()
	tmp, _ := json.Marshal(data.data)
	out, _ := os.OpenFile("data.json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	out.Write(tmp)
	out.Close()
}

func worm() {
	for id := range task {
		req, _ := http.NewRequest("GET", url+id, nil)
		req.Header = map[string][]string{
			"Cookie": {cookie},
		}
		res, _ := cli.Do(req)

		var tmp []Info
		content, _ := io.ReadAll(res.Body)
		json.Unmarshal(content, &tmp)
		if len(tmp) > 0 {
			tmp[0].Pid = id
			tmp[0].SzLogonName = id
			data.l.Lock()
			data.data[id] = tmp[0].Name
			data.l.Unlock()
			fmt.Println(tmp[0])
		} else {
			fmt.Println(id + " nil")
		}
	}
	wg.Done()
}
