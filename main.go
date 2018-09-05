package main

import (
	"fmt"
	"net/http"
	//"os"
	//"path"
	"bufio"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Imgdata struct {
	uncompleted int
	taskArray   []string
	taskMask    []bool
	resultArray []string
}

func getRandomImg(dataset string, task string, dd Imgdata) (string, int) {

	if dd.uncompleted == 0 {
		fmt.Println("success")
		return "", -1
	} else {
		fmt.Println("uncompleted:", dd.uncompleted)
	}
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(dd.uncompleted)

	k := 0
	//?????????????????
	for i := -1; i < x; k++ {
		if dd.taskMask[k] {
			i++
		}
	}
	return dd.taskArray[k], k
}

func setResult(result string, index int, dd Imgdata) Imgdata {
	if index < len(dd.taskArray) && index >= 0 {
		dd.resultArray[index] = result
		dd.taskMask[index] = false
		dd.uncompleted--
	}
	return dd
}

func readFile2DataSet(fileName string) Imgdata {
	var ret Imgdata
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return ret
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	var size = stat.Size()
	fmt.Println("File size =", size)

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		ret.taskArray = append(ret.taskArray, line)
		ret.taskMask = append(ret.taskMask, true)
		ret.resultArray = append(ret.resultArray, "")
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return ret
			}
		}
	}
	ret.uncompleted = len(ret.taskMask)
	return ret
	///////////////////////////////////////////////////
}
func route(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	dataset := strings.Join(r.Form["dataset"], "")
	task := strings.Join(r.Form["task"], "")
	if len(dataset) == 0 || len(task) == 0 {
		fmt.Fprintln(w, "parameter error")
		r.Body.Close()
		return
	}
	r.Body.Close()
	img, index := getRandomImg("1111", "2222", allData)
	index = index
	img = img
	send := `{"img":"` + dataPath + img + `","index":"` + strconv.Itoa(index) + `"}`
	fmt.Fprintln(w, send)
}

func submit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	index := strings.Join(r.Form["index"], "")
	result := strings.Join(r.Form["result"], "")
	if len(result) == 0 || len(index) == 0 {
		fmt.Fprintln(w, "parameter error")
		r.Body.Close()
		return
	}
	idx, _ := strconv.Atoi(index)
	allData = setResult(result, idx, allData)

	fmt.Fprintln(w, "submit_success")
}
func process(w http.ResponseWriter, r *http.Request) {
	ret := ""
	cc := 0
	for _, v := range allData.taskMask {
		if v {
			ret += "0"
		} else {
			ret += "1"
			cc += 1
		}
	}
	fmt.Fprintln(w, strconv.Itoa(cc)+"<br/>"+ret)
}
func result(w http.ResponseWriter, r *http.Request) {
	ret := "["
	for k := 0; k < len(allData.resultArray); k++ {
		if len(allData.resultArray[k]) != 0 {
			ret += (`{"name":"` + allData.taskArray[k] + `","result":"` + allData.resultArray[k] + `"},`)
		}
	}
	ret += "]"
	fmt.Fprintln(w, ret)
}

var allData Imgdata
var dataPath = "http://10.231.56.131/bmp/"

func main() {
	allData = readFile2DataSet("imglist.txt")
	//
	/*
		for i := 0; i < 3000; i++ {
			img, index := getRandomImg("1111", "2222", allData)
			fmt.Println(img, index, 1110-i)
			allData = setResult(img, index, allData)

		}
	*/
	http.HandleFunc("/getImg/", route)
	http.HandleFunc("/submit/", submit)
	http.HandleFunc("/process/", process)
	http.HandleFunc("/result/", result)
	http.ListenAndServe(":8090", nil)
}

//http://127.0.0.1:8090/submit/?index=118&result=23232323
//http://127.0.0.1:8090/getImg/?dataset=111&task=555
//http://127.0.0.1:8090/process/
