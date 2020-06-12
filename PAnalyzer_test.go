package pcrawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type ReqWSEntry struct {
	Key      string `json:"key"`
	Sentence string `json:"sentence"`
}

type ReqWSPayload struct {
	Payload []ReqWSEntry `json:"sentence_with_keys"`
}

func TestWSServer(t *testing.T) {
	//wsserver := "http://node1.rayer.idv.tw/"
	c := NewCrawler("gossiping")
	indexInfo, err := c.ParseIndex(100)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	collector := NewUserSentenceInfoCollector()

	for _, a := range indexInfo.Articles {
		doc, err := c.ParseDocument(a.Url)
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}
		// fmt.Printf("%+v\n", doc)
		//Parse and collect commit info
		collector.Collect(doc)
	}

	payload := ReqWSPayload{}
	//count := 0
	for k, v := range collector.GetMap() {
		for _, s := range v {
			payload.Payload = append(payload.Payload, ReqWSEntry{
				Key:      k,
				Sentence: s.Sentence,
			})
			//count++
		}
		//if count > 30 {
		//	break
		//}
	}
	body, _ := json.Marshal(&payload)
	fmt.Printf("payload : %v\n", string(body))
	resp, err := http.Post("http://node1.rayer.idv.tw:8000/", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	defer resp.Body.Close()
	res_json, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(res_json))
}