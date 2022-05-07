package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	client := &http.Client{}
	var data = strings.NewReader(`from=en&to=zh&query=bad&transtype=realtime&simple_means_flag=3&sign=424354.170643&token=8db8c0f0e089e3dbf510ea52f8a88d4b&domain=common`)
	req, err := http.NewRequest("POST", "https://fanyi.baidu.com/v2transapi?from=en&to=zh", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "__yjs_duid=1_023dad79550dce1c8a501f57e678af6b1645977456887; BDUSS=mk2eTFLcHdGTlpPdGk5Y3FpRlNoUFd2clo4WGJ3blF2SjRiV1A2RHlqVGlOVTFpRVFBQUFBJCQAAAAAAAAAAAEAAABT5x6btPK59rXEvfDLv8i4AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOKoJWLiqCViT; BDUSS_BFESS=mk2eTFLcHdGTlpPdGk5Y3FpRlNoUFd2clo4WGJ3blF2SjRiV1A2RHlqVGlOVTFpRVFBQUFBJCQAAAAAAAAAAAEAAABT5x6btPK59rXEvfDLv8i4AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOKoJWLiqCViT; BAIDUID=55C11F053336F2B438EBF386C6B96DD2:FG=1; BIDUPSID=55C11F053336F2B438EBF386C6B96DD2; PSTM=1648645090; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; BAIDUID_BFESS=55C11F053336F2B438EBF386C6B96DD2:FG=1; BA_HECTOR=8ha0058k0l010h05t71h7cad50r; delPer=0; PSINO=6; H_PS_PSSID=36309_36367_36165_34584_35978_35802_36232_26350_36311_36061; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1651915350; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; Hm_lpvt_64ecd82404c51e03dc91cb9e8c025574=1651915488; ab_sr=1.0.1_NjI3ZWE4MjA4ZDRkZTYwNzY1MmI0NTQ5MjZkZDFhMmFkMjVmYjg2ZDYxOWFlNzI0YTViOThkMGI2MDllNjQzMzliMWY4Mzg2Mzg1YThhZGQ3MTcwNDdkNzBmOTVjZjgxMzM0YTMyYWI1M2NmYTVjNGMzMjk4ZTZkYmI5ODljZGZkN2M1YTc2NDhmNzZmMzEwNTUzOGZjZDFhYjU0NjA5ZTM1ZDAyNTA2ODk2Zjc4ZmRlZmM1MDdjYjgwNmQ3NDc4")
	req.Header.Set("Origin", "https://fanyi.baidu.com")
	req.Header.Set("Referer", "https://fanyi.baidu.com/translate?aldtype=16047&query=&keyfrom=baidu&smartresult=dict&lang=auto2zh")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
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
	fmt.Printf("%s\n", bodyText)
}
