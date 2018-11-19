package main

import (
	"os"
        "fmt"
        "net/http"
        "log"
        "io/ioutil"
	"strings"
	"time"
	"strconv"
	"encoding/json"
	"bytes"
	"github.com/kylelemons/go-gypsy/yaml"
)

type MetricValue struct {
	Endpoint  string            `json:"endpoint"`
	Metric    string            `json:"metric"`
	Value     interface{}       `json:"value"`
	Step      int64             `json:"step"`
	Type      string            `json:"counterType"`
	Tags      string	    `json:"tags"`
	Timestamp int64             `json:"timestamp"`
}

func doPost(falcon_url string, buf []byte) error {
        resp, err := http.Post(falcon_url, "application/json", bytes.NewBuffer(buf))
        if err == nil {
            if resp.StatusCode != 200 {
                log.Printf("Post return status code %d\n", resp.StatusCode)
            }
        } else {
            log.Printf("Post warnning error:%s\n", err.Error())
        }
	resp.Body.Close()
	return nil
}

func getNginx(nginx_url string) (int, int, int, int, int, int, int){
        resp, err := http.Get(nginx_url)
        if err != nil {
                log.Println(err)
        }

        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
	s := strings.Split(string(body), " " )
        Active, err := strconv.Atoi(s[2])
        Reading, err := strconv.Atoi(s[11])
        Writing, err := strconv.Atoi(s[13])
        Waiting, err := strconv.Atoi(s[15])
        Accepts, err := strconv.Atoi(s[7])
        Handled, err := strconv.Atoi(s[8])
        Requests, err := strconv.Atoi(s[9])

        return Active, Reading, Writing, Waiting, Accepts, Handled, Requests

}

func main() {

	config, err := yaml.ReadFile("conf.yaml")
        if err != nil {
            fmt.Println(err)
        }

        host, err := os.Hostname()
        if err != nil {
            fmt.Println(err)
        }
        falcon_url, err := config.Get("falcon-url")
        nginx_url, err := config.Get("nginx-url")

	Active, Reading, Writing, Waiting, Accepts, Handled, Requests := getNginx(nginx_url)
	nginx_values := map[string]int{"active": Active, "reading": Reading,  "writing": Writing, "waiting": Waiting, "accepts": Accepts, "handled": Handled, "requests": Requests}

        if falcon_url != "" {
            for n_metric,n_value := range nginx_values {
                var Met = []MetricValue{{Endpoint: host, Metric: n_metric, Value: n_value, Step: 60, Type: "GAUGE", Tags: "nginx=status", Timestamp: time.Now().Unix()}}
                buf, _ := json.Marshal(Met)
	        err = doPost(falcon_url,buf)
	        if err != nil {
	    	log.Println("Cannot post the data:", err)
	        }
	    }
        }

}
