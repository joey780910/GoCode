package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Ubike struct {
	RetCode int `json:"retCode"`
}

type UbikeSteps struct {
	steps []Step
}

type StepFromJson struct {
	// Sno   string `json:"sno"`
	Sna string `json:"sna"`
	Tot string `json:"tot"`
	Sbi string `json:"sbi"`
	// Sarea string `json:"sarea"`
	// Mday  string `json:"mday"`
	Lat string `json:"lat"`
	Lng string `json:"lng"`
	Ar  string `json:"ar"`
	// Sareaen string `json:"sareaen"`
	Snaen string `json:"snaen"`
	Aren  string `json:"aren"`
	Bemp  string `json:"bemp"`
	Act   string `json:"act"`
}

type Step struct {
	// Sno   string
	Sna string
	Tot string
	Sbi string
	// Sarea string
	// Mday  string
	Lat float64
	Lng float64
	Ar  string
	// Sareaen string
	Snaen string
	Aren  string
	Bemp  string
	Act   string
}

type UbikeInfo struct {
	RetCode int
	RetVal  []Step
}

func main() {

	url := "http://data.taipei/youbike"

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// fmt.Println(string(body))

	_Ubike := Ubike{}
	jsonErr := json.Unmarshal(body, &_Ubike)
	if jsonErr != nil {
		fmt.Println(jsonErr)
		return
	}

	ubike := UbikeRetCode(body)
	fmt.Println("retCode", ubike.RetCode)

	//接收的Json資料
	stepsFormJSON := UbikeRetValJSONSplit(string(body))

	//可用的資料(some string type change)
	steps := GetConvertStep(stepsFormJSON)
	fmt.Println("All Step Count=", len(steps))

	testStep := GetStepByCoordinate(steps, steps[5].Lat, steps[5].Lng)
	fmt.Println("GetStepByCoordinate Test", steps[5].Lat, testStep.Lat)

	nearbySteps := GetNearbySteps(steps, 25.051343, 121.582167, 25.047580, 121.574958)
	fmt.Println("nearbySteps count", len(nearbySteps))
	for i := 0; i < len(nearbySteps); i++ {
		fmt.Println(nearbySteps[i].Lat, nearbySteps[i].Lng)
	}

	ubikeInfo := UbikeInfoToJSON(_Ubike.RetCode, nearbySteps)
	fmt.Println(string(ubikeInfo))
}

// UbikeRetCode 接收retCode
func UbikeRetCode(body []byte) (_Ubike Ubike) {
	jsonErr := json.Unmarshal(body, &_Ubike)
	if jsonErr != nil {
		fmt.Println(jsonErr)
	}
	return
}

// UbikeRetValJSONSplit 取出retVal
func UbikeRetValJSONSplit(str string) (steps []StepFromJson) {
	s := strings.Split(string(str), ":{")
	for i := 0; i < len(s); i++ {
		_step := StepFromJson{}
		str := "{" + strings.Split(s[i], "}")[0] + "}"

		// fmt.Println(str)

		err := json.Unmarshal([]byte(str), &_step)
		if err != nil {
			// fmt.Println(err)
			continue
		}
		steps = append(steps, _step)
	}
	return
}

// GetConvertStep 將retVel中某些string轉換成其他type
func GetConvertStep(stepsJSON []StepFromJson) (newSteps []Step) {
	for i := 0; i < len(stepsJSON); i++ {
		convertStep := Step{}
		convertStep.Sna = stepsJSON[i].Sna
		convertStep.Tot = stepsJSON[i].Tot
		convertStep.Sbi = stepsJSON[i].Sbi
		convertStep.Lat = StringToFloat(stepsJSON[i].Lat)
		convertStep.Lng = StringToFloat(stepsJSON[i].Lng)
		convertStep.Snaen = stepsJSON[i].Snaen
		convertStep.Aren = stepsJSON[i].Aren
		convertStep.Bemp = stepsJSON[i].Bemp
		convertStep.Act = stepsJSON[i].Act

		newSteps = append(newSteps, convertStep)
	}
	return
}

// GetStepByCoordinate 找出lat,lng相同的租借站
func GetStepByCoordinate(steps []Step, lat float64, lng float64) (step Step) {
	for i := 0; i < len(steps); i++ {
		if lat == steps[i].Lat && lng == steps[i].Lng {
			step = steps[i]
			return
		}
	}
	return
}

// GetNearbySteps 找範圍內的所有租借站
func GetNearbySteps(steps []Step, maxLat float64, maxLng float64, mimLat float64, minLng float64) (nearbySteps []Step) {
	for i := 0; i < len(steps); i++ {
		if steps[i].Lat < maxLat &&
			steps[i].Lat > mimLat &&
			steps[i].Lng < maxLng &&
			steps[i].Lng > minLng {
			nearbySteps = append(nearbySteps, steps[i])
		}
	}
	return
}

// StringToFloat 將string轉乘float64
func StringToFloat(str string) (retVal float64) {
	retVal, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// UbikeInfoToJSON 轉為JSON
func UbikeInfoToJSON(_retCode int, _retVal []Step) (ubikeInfoJSON []byte) {
	ubikeInfo := UbikeInfo{RetCode: _retCode, RetVal: _retVal}
	// ubikeInfo.retCode = _retCode
	// ubikeInfo.retVal = _retVal

	out, err := json.Marshal(ubikeInfo)
	if err != nil {
		fmt.Println(err)
	}
	ubikeInfoJSON = out
	return
}

type response1 struct {
	RetCode int
	Fruits  []string
}

func aa() {
	res1D := &response1{
		RetCode: 1}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))
	fmt.Println("end")
	// res1D := &UbikeInfo{
	// 	retCode: 1,
	// }
	// res1B, _ := json.Marshal(res1D)
	// fmt.Println(string(res1B))
}
