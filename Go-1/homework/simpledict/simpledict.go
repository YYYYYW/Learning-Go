package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type AutoGenerated struct {
	Rc   int `json:"rc"`
	Wiki struct {
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

type HuoShanDiction struct {
	Words []struct {
		Source  int    `json:"source"`
		Text    string `json:"text"`
		PosList []struct {
			Type      int `json:"type"`
			Phonetics []struct {
				Type int    `json:"type"`
				Text string `json:"text"`
			} `json:"phonetics"`
			Explanations []struct {
				Text     string `json:"text"`
				Examples []struct {
					Type      int `json:"type"`
					Sentences []struct {
						Text      string `json:"text"`
						TransText string `json:"trans_text"`
					} `json:"sentences"`
				} `json:"examples"`
				Synonyms []interface{} `json:"synonyms"`
			} `json:"explanations"`
			Relevancys []interface{} `json:"relevancys"`
		} `json:"pos_list"`
	} `json:"words"`
	Phrases  []interface{} `json:"phrases"`
	BaseResp struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"base_resp"`
}

type HuoShanDictRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}

var l sync.Mutex

var isOutput = 0

func huoshanQuery(word string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	client := &http.Client{}
	request := HuoShanDictRequest{Text: word, Language: "en"}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	// var data = strings.NewReader(`{"text":"good\n","language":"en"}`)
	req, err := http.NewRequest("POST", "https://translate.volcengine.com/web/dict/match/v1/?msToken=&X-Bogus=DFSzswVLQDVBKiQrSWQR4Pt/pLv9&_signature=_02B4Z6wo00001yjd.7gAAIDCo5ZkWs8NEw8o3fsAAKhLNQOsQfcxG2TsBgxwl7OK8tsxiuNz4KNSUdOHH.FY7VnK2B-b7XGr1F-h2H9gFcyNMdkZ5K46sl8nevCpsimTy5CyHUvb73Mc5YpY95", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "translate.volcengine.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "x-jupiter-uuid=16519240543857514; i18next=zh-CN; ttcid=4dd1d6481a27405d8e088f7a6ab7d51a66; tt_scid=jlnel4RAjPdPt4fMI2z3xvfsjix73mk6PiqxOzXNcJE9sACpNxnSP4OjtY69fQQ833e2; s_v_web_id=verify_0b043b2145bfe9fab707ee95bf91ebc6; _tea_utm_cache_2018=undefined")
	req.Header.Set("origin", "https://translate.volcengine.com")
	req.Header.Set("referer", "https://translate.volcengine.com/translate?category=&home_language=zh&source_language=detect&target_language=zh&text=good%0A")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var huoshanDiction HuoShanDiction
	err = json.Unmarshal(bodyText, &huoshanDiction)
	if err != nil {
		log.Fatal(err)
	}

	l.Lock()
	if isOutput == 0 {
		fmt.Println(word, "\nUK:", huoshanDiction.Words[0].PosList[0].Phonetics[0].Text,
			"US:", huoshanDiction.Words[0].PosList[0].Phonetics[1].Text)
		for _, items := range huoshanDiction.Words[0].PosList {
			for _, exps := range items.Explanations {
				fmt.Println(exps.Text)
			}
		}
		isOutput = 1
		fmt.Println("-------------------------------Query by ????????????")
	}
	l.Unlock()
	// fmt.Printf("%#v\n", huoshanDiction)
}

func query(word string, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()
	client := &http.Client{}
	request := DictRequest{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("app-name", "xy")
	req.Header.Set("os-type", "web")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var autoGenerated AutoGenerated
	err = json.Unmarshal(bodyText, &autoGenerated)
	if err != nil {
		log.Fatal(err)
	}
	l.Lock()
	if isOutput == 0 {
		fmt.Println(word, "UK:", autoGenerated.Dictionary.Prons.En, "US:", autoGenerated.Dictionary.Prons.EnUs)
		for _, item := range autoGenerated.Dictionary.Explanations {
			fmt.Println(item)
		}
		isOutput = 1
		fmt.Println("-------------------------------Query by ????????????")
	}
	l.Unlock()
	//fmt.Printf("%s\n", bodyText)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	waitgroup := sync.WaitGroup{}
	waitgroup.Add(2)
	go huoshanQuery(word, &waitgroup)
	go query(word, &waitgroup)
	waitgroup.Wait()
}
