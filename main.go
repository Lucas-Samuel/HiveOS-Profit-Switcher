package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

type coinsProfit struct {
	tag   string
	value float64
}

type Configs struct {
	APIKey        string  `json:"api_key"`
	FarmID        string  `json:"farm_id"`
	CoinDiference string  `json:"coin_diference"`
	Workers       Workers `json:"workers"`
}

type Workers []struct {
	Name        string `json:"name"`
	WtmEndpoint string `json:"wtm_endpoint"`
	Coins       Coins  `json:"coins"`
}

type Coins []struct {
	Tag string `json:"tag"`
	Fs  string `json:"fs"`
}

var configs Configs

func (configs *Configs) load() {
	content, err := ioutil.ReadFile("configs.json")

	if err != nil {
		fmt.Println("Error when opening config file: ", err)
		os.Exit(1)
	}

	json.Unmarshal([]byte(content), &configs)

	if configs.APIKey == "" {
		fmt.Println("API Key not set")
		os.Exit(1)
	}

	if configs.FarmID == "" {
		fmt.Println("Farm id not set")
		os.Exit(1)
	}
}

func main() {
	configs.load()

	var algoMap = map[string]string{
		"AUTOLYKOS": "al_p", "BEAMHASHIII": "eqb_p", "CORTEX": "cx_p", "CRYPTONIGHTFASTV2": "cnf_p", "CRYPTONIGHTGPU": "cng_p", "CRYPTONIGHTHAVEN": "cnh_p", "CUCKAROO29S": "cr29_p", "CUCKATOO31": "ct31_p", "CUCKATOO32": "ct32_p", "CUCKOOCYCLE": "cc_p", "EQUIHASH (210,9)": "eqa_p", "EQUIHASHZERO": "eqz_p", "ETCHASH": "e4g_p", "ETHASH": "eth_p", "ETHASH4": "e4g_p", "FIROPOW": "fpw_p", "KAWPOW": "kpw_p", "NEOSCRYPT": "ns_p", "OCTOPUS": "ops_p", "PROGPOW": "ppw_p", "PROGPOWZ": "ppw_p", "RANDOMX": "rmx_p", "UBQHASH": "e4g_p", "VERTHASH": "vh_p", "X25X": "x25x_p", "ZELHASH": "zlh_p", "ZHASH": "zh_p",
	}

	result := requestHive("GET", "/farms/FARM_ID/workers", nil)
	data := result["data"].([]interface{})

	for _, value := range data {
		worker := value.(map[string]interface{})

		var workerKey int = -1
		for key, cw := range configs.Workers {
			if cw.Name == worker["name"] {
				workerKey = key
				break
			}
		}

		if workerKey < 0 {
			fmt.Println("Worker \"" + worker["name"].(string) + "\" not found on config file")
			os.Exit(1)
		}

		config := configs.Workers[workerKey]
		currentFs := worker["flight_sheet"].(map[string]interface{})

		var fsKey int = -1
		for key, coin := range config.Coins {
			if coin.Fs == currentFs["name"] {
				fsKey = key
				break
			}
		}

		if fsKey < 0 {
			fmt.Println("Flight sheet \"" + currentFs["name"].(string) + "\" not found on config file")
			os.Exit(1)
		}

		currentCoin := config.Coins[fsKey].Tag

		resultBtc := request("https://api.coindesk.com/v1/bpi/currentprice.json")
		btcPrice := resultBtc["bpi"].(map[string]interface{})["USD"].(map[string]interface{})["rate_float"].(float64)

		WtmEndpoint := config.WtmEndpoint
		params, _ := url.ParseQuery(WtmEndpoint)

		powerCost, _ := strconv.ParseFloat(params["factor[cost]"][0], 32)

		resultWTM := request(WtmEndpoint)
		coins := resultWTM["coins"].(map[string]interface{})

		var profits []coinsProfit
		var currentCoinPrice float64

		for _, coin := range coins {
			coin := coin.(map[string]interface{})

			btcRevenue, _ := strconv.ParseFloat(coin["btc_revenue24"].(string), 64)
			algo := strings.ToUpper(coin["algorithm"].(string))

			if algo == "ETHASH" && coin["tag"] != "ETH" && coin["tag"] != "NICEHASH" {
				algo = "ETHASH4"
			}

			consumption, _ := strconv.ParseFloat(params["factor["+algoMap[algo]+"]"][0], 64)

			dailyPowerCost := 24 * (consumption / 1000) * powerCost
			dailyRevenue := btcRevenue * btcPrice
			dailyProfit := dailyRevenue - dailyPowerCost

			var key = coin["tag"]

			if key == "NICEHASH" {
				key = coin["tag"].(string) + "-" + algo
			}

			profit := coinsProfit{tag: key.(string), value: dailyProfit}
			profits = append(profits, profit)

			if profit.tag == currentCoin {
				currentCoinPrice = profit.value
			}
		}

		sort.SliceStable(profits, func(i, j int) bool {
			return profits[i].value > profits[j].value
		})

		bestCoin := profits[0].tag
		bestCoinPrice := profits[0].value

		coinDiference, _ := strconv.ParseFloat(configs.CoinDiference, 32)

		currentCoinPrice += currentCoinPrice * (coinDiference / 100)

		if bestCoin == currentCoin || bestCoinPrice < currentCoinPrice {
			fmt.Println("Already in best coin")
			os.Exit(0)
		}

		var newFsId float64 = -1
		var newFsName string

		for _, coin := range config.Coins {
			if coin.Tag == bestCoin {
				newFsName = coin.Fs
				break
			}
		}

		if newFsName == "" {
			fmt.Println("flight Sheet not found for coin \"" + bestCoin + "\"")
			os.Exit(1)
		}

		result := requestHive("GET", "/farms/FARM_ID/fs", nil)
		data := result["data"].([]interface{})

		for _, value := range data {
			sheet := value.(map[string]interface{})

			if sheet["name"] == newFsName {
				newFsId = sheet["id"].(float64)
				break
			}
		}

		if newFsId < 0 {
			fmt.Println("flight Sheet \"" + newFsName + "\" not found on HiveOS")
			os.Exit(1)
		}

		payload, _ := json.Marshal(map[string]interface{}{
			"fs_id": newFsId,
		})

		requestHive("PATCH", "/farms/FARM_ID/workers/"+strconv.Itoa(int(worker["id"].(float64))), bytes.NewBuffer(payload))

		fmt.Println("Worker \""+worker["name"].(string)+"\" flight Sheet updated to \""+newFsName+"\". Estimated profit in 24h: $", bestCoinPrice)
	}
}

func requestHive(method string, url string, payload io.Reader) map[string]interface{} {
	client := http.Client{}
	req, err := http.NewRequest(method, "https://api2.hiveos.farm/api/v2"+strings.Replace(url, "FARM_ID", configs.FarmID, 1), payload)

	if err != nil {
		fmt.Println("Error in communication with HiveOS: ", err)
		os.Exit(1)
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + configs.APIKey},
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	return result
}

func request(url string) map[string]interface{} {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Error in communication with "+url+": ", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	return result
}
