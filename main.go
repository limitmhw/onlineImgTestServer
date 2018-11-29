package main

import (
	"fmt"
	"net/http"
	//"os"
	//"path"
	"bufio"
	"flag"
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

func writeFile(fileName string, content string) {
	fileName = fileName + `.txt`
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, _ = f.Write([]byte(content))
	}
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
	send := `{"img":"` + (*dataPath) + img + `","index":"` + strconv.Itoa(index) + `"}`
	fmt.Fprintln(w, send)
}

func gethtml(w http.ResponseWriter, r *http.Request) {
	aa := `<html>

<head>
<title>我的第一个 HTML 页面</title>
</head>

<body>
<center>
</br>
</br>
<div id="bb"  style="background-color:#ccc;width:500px;height:500px;"></div>
</br>
<p id='imgname'></p>
</br>
<button id="geti">获取一张图片</button>
<button class="p1" data-v="0">笑1</button>
<button class="p1" data-v="1">笑2</button>
<button class="p1" data-v="2">笑3</button>
<button class="p1" data-v="3">笑4</button>
</br></br>
<a href="/result/" target="view_window">查看结果</a>
</center>

</body>
<script src="http://libs.baidu.com/jquery/2.0.0/jquery.min.js"></script>
<script>
$(function(){
$("#geti").on("click",function(){
   $.ajax({
   	type: "POST", 
   	url: "/getImg/",
   	data: {
   		"dataset": "111",
		"task":"555"
   	}, 
   	success: function (dd) {
		console.log(dd);
	    let data = JSON.parse(dd);//eval('"'+dd+'"');
	    let url = data.img;
		$(".p1").attr('data-i',data.index)
	    $('#bb').css("background-image", 'url('+url+')');
		$('#imgname').text(data.img)
   	}
   });

})

$(".p1").on("click",function(){
	$(this).attr('data-i')
	$.ajax({
	   	type: "POST", 
	   	url: "/submit/",
	   	data: {
	   		"index": $(this).attr('data-i'),
			"result":$(this).attr('data-v')
	   	}, 
	   	success: function (dd) {
 			$("#geti").trigger("click");
	   	}
	   });

})

})

</script>

</html>`
	fmt.Fprintln(w, aa)
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
	writeFile(allData.taskArray[idx], allData.resultArray[idx])

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

var dataPath *string = flag.String("dataPath", "http://10.231.56.131/bmp/", "dataPath")
var imglist *string = flag.String("imglist", "imglist.txt", "imglist")

func main() {
	flag.Parse()
	allData = readFile2DataSet((*imglist))
	//
	/*
		for i := 0; i < 3000; i++ {
			img, index := getRandomImg("1111", "2222", allData)
			fmt.Println(img, index, 1110-i)
			allData = setResult(img, index, allData)

		}
	*/
	http.HandleFunc("/gethtml/", gethtml)
	http.HandleFunc("/getImg/", route)
	http.HandleFunc("/submit/", submit)
	http.HandleFunc("/process/", process)
	http.HandleFunc("/result/", result)
	http.ListenAndServe(":8090", nil)
}

//http://127.0.0.1:8090/submit/?index=118&result=23232323
//http://127.0.0.1:8090/getImg/?dataset=111&task=555
//http://127.0.0.1:8090/process/
