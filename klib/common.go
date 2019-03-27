package klib

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nfnt/resize"
	"github.com/tealeg/xlsx"
)

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var AppDebug string

var mutexLog sync.Mutex
var mutexHttpLog sync.Mutex
var mutexSqlLog sync.Mutex

var coder = base64.NewEncoding(base64Table)
var tenToAny map[int]string = map[int]string{0: "0", 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "a", 11: "b", 12: "c", 13: "d", 14: "e", 15: "f", 16: "g", 17: "h", 18: "i", 19: "j", 20: "k", 21: "l", 22: "m", 23: "n", 24: "o", 25: "p", 26: "q", 27: "r", 28: "s", 29: "t", 30: "u", 31: "v", 32: "w", 33: "x", 34: "y", 35: "z", 36: ":", 37: ";", 38: "<", 39: "=", 40: ">", 41: "?", 42: "@", 43: "[", 44: "]", 45: "^", 46: "_", 47: "{", 48: "|", 49: "}", 50: "A", 51: "B", 52: "C", 53: "D", 54: "E", 55: "F", 56: "G", 57: "H", 58: "I", 59: "J", 60: "K", 61: "L", 62: "M", 63: "N", 64: "O", 65: "P", 66: "Q", 67: "R", 68: "S", 69: "T", 70: "U", 71: "V", 72: "W", 73: "X", 74: "Y", 75: "Z"}

/*
#string到int
int,err:=strconv.Atoi(string)
#string到int64
int64, err := strconv.ParseInt(string, 10, 64)
#int到string
string:=strconv.Itoa(int)
#int64到string
string:=strconv.FormatInt(int64,10)
*/

func base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}

func Mcrypt(data string, mode string) string {

	var passcrypt string
	key := "justware2013"
	if mode == "decode" {
		if len(data) < 3 {
			return ""
		}
		if data[0:2] != getOrd(data[2:]+key) {
			return ""
		}
		str := strings.Replace(data[2:], " ", "+", -1)
		str = strings.Replace(str, "-_", "+/", -1)

		enbyte, err := base64Decode([]byte(str))
		if err != nil {
			passcrypt = ""
		} else {
			passcrypt = string(enbyte)
		}
	}

	if mode == "encode" {
		enbyte := base64Encode([]byte(data))

		passcrypt = string(enbyte)
		passcrypt = strings.Replace(passcrypt, "+/", "-_", -1)
		passcrypt = getOrd(passcrypt+key) + passcrypt
	}
	return passcrypt
}

func getOrd(str string) string {
	sum := 0
	if len(str) == 0 {
		return "00"
	} else {
		sum = int(str[0])
		var i64 int64
		i64 = int64(sum % 265)
		ret := strconv.FormatInt(i64, 16)
		if len(ret) == 1 {
			return "0" + ret
		} else {
			return ret
		}
	}
}

func ReturnJson(ret map[string]interface{}) string {
	b, _ := json.Marshal(ret)
	return Mcrypt(string(b), "encode")
}

//校验字符串是否整数
func IsInt(str string) bool {

	if str == "" {
		return false
	}

	reg := regexp.MustCompile(`[^\d]`)

	if len(reg.FindAllString(str, -1)) > 0 {
		return false
	} else {
		return true
	}
}

/**
 * 单引号替换
 */
func FormatSql(str string) string {
	return strings.Replace(str, "'", "''", -1)
}

func FormatTimeToStr(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000000000")
}

func printT(t time.Time) time.Time {
	return t.Add(100)
}

//土方法判断进程是否启动
func AppInit() {
	iManPid := fmt.Sprint(os.Getpid())
	tmpDir := os.TempDir()
	fmt.Println(tmpDir)
	if err := ProcExsit(tmpDir); err == nil {
		pidFile, _ := os.Create(tmpDir + "/eposPack.pid")
		defer pidFile.Close()
		pidFile.WriteString(iManPid)
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 判断进程是否启动
func ProcExsit(tmpDir string) (err error) {
	iManPidFile, err := os.Open(tmpDir + "/eposPack.pid")
	defer iManPidFile.Close()
	if err == nil {
		filePid, err := ioutil.ReadAll(iManPidFile)
		if err == nil {
			pidStr := fmt.Sprintf("%s", filePid)
			pid, _ := strconv.Atoi(pidStr)
			_, err := os.FindProcess(pid)
			if err == nil {
				return errors.New("[ERROR] epos已启动.")
			}
		}
	}
	return nil
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

/**
 * 判断文件夹是否存在  存在返回 true 不存在返回false
 */
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//获取程序路径
func GetAppPath() string {
	//os.Getwd()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		return dir
	}
}

//http log
func WriteHttpLog(r *http.Request) {

	mutexHttpLog.Lock()
	defer mutexHttpLog.Unlock()

	dir := GetAppPath()
	if dir == "" {
		return
	}

	logPath := dir + "/log/"
	if ok, err := PathExists(logPath); !ok {
		err = os.Mkdir(logPath, os.ModePerm) //在当前目录下生成md目录
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	now := time.Now()
	logFile := "http_" + now.Format("200601") + ".log"

	fd, err := os.OpenFile(logPath+logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	str := `[` + now.Format("2006-01-02 15:04:05.0000") + `] ` + r.RemoteAddr + ` - "` + r.Method + ` ` + r.RequestURI + ` ` + r.Proto + `"` + "\n"
	buf := []byte(str)
	fd.Write(buf)
	fd.Close()

	fmt.Print(str)
	/*fmt.Println(r.Host)        //121.41.48.197:9090
	fmt.Println(r.Method)      //POST
	fmt.Println(r.RemoteAddr)  //59.172.170.173:19625
	fmt.Println(r.RequestURI)  // /
	fmt.Println(r.URL)         // /
	fmt.Println(r.UserAgent()) //
	fmt.Println(r.Referer())   //
	fmt.Println(r.Proto)       // HTTP/1.1
	*/
}

//log
func WriteLog(s string) {

	mutexLog.Lock()
	defer mutexLog.Unlock()

	dir := GetAppPath()
	if dir == "" {
		return
	}

	logPath := dir + "/log/"
	if ok, err := PathExists(logPath); !ok {
		err = os.Mkdir(logPath, os.ModePerm) //在当前目录下生成md目录
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	now := time.Now()
	logFile := now.Format("200601") + ".log"

	fd, err := os.OpenFile(logPath+logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	str := `[` + now.Format("2006-01-02 15:04:05.0000") + `] ` + s + "\n"

	buf := []byte(str)
	fd.Write(buf)
	fd.Close()

	fmt.Print(str)
}

//sql log
func WriteSqlLog(s string) {

	mutexSqlLog.Lock()
	defer mutexSqlLog.Unlock()

	dir := GetAppPath()
	if dir == "" {
		return
	}

	logPath := dir + "/log/"
	if ok, err := PathExists(logPath); !ok {
		err = os.Mkdir(logPath, os.ModePerm) //在当前目录下生成md目录
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	now := time.Now()
	logFile := "sql_" + now.Format("200601") + ".log"

	fd, err := os.OpenFile(logPath+logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Print(err)
		return
	}

	str := `[` + now.Format("2006-01-02 15:04:05.0000") + `] ` + s + "\n"
	buf := []byte(str)
	fd.Write(buf)
	fd.Close()
}

//err := SendToMail("yang**@yun*.com", "***", "smtp.exmail.qq.com:25", "397685131@qq.com", subject, body, "html")
func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

//获取上一级函数名称等信息
func GetCurrFuncInfo() string {

	pc, _, line, ok := runtime.Caller(1)
	if ok {
		f := runtime.FuncForPC(pc)
		return `(Func: ` + f.Name() + `;line: ` + strconv.Itoa(line) + `;)`
	} else {
		return ""
	}
}

func DownloadImg(url string, filename string) int64 {

	if len(url) <= 0 {
		return 0
	}

	var resultErr error

	for i := 0; i < 20; i++ {

		if i > 0 {
			WriteLog(`请求(` + url + `)失败, 重试中...(` + strconv.Itoa(i) + `)`)
			time.Sleep(200 * time.Millisecond)
		}
		client := &http.Client{}
		reqest, err := http.NewRequest("GET", url, nil)

		if err != nil {
			resultErr = err
			continue
		}

		reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		reqest.Header.Add("Accept-Encoding", "gzip, deflate")
		reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
		reqest.Header.Add("Cache-Control", "max-age=0")
		reqest.Header.Add("Connection", "keep-alive")
		//reqest.Header.Add("Host", "www.dianping.com")
		//reqest.Header.Add("Referer", "http://www.dianping.com/wuhan")
		reqest.Header.Add("Upgrade-Insecure-Requests", "1")
		reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
		response, err := client.Do(reqest)
		if err != nil {
			resultErr = err
			continue
		}

		if response.StatusCode == 200 || response.StatusCode == 304 {

			defer response.Body.Close()

			var bodyByte []byte
			ext := `jpg`

			if response.Header.Get("Content-Type") == "image/jpeg" {
				ext = `jpg`
				bodyByte, _ = ioutil.ReadAll(response.Body)
			} else if response.Header.Get("Content-Type") == "image/png" {
				ext = `png`
				bodyByte, _ = ioutil.ReadAll(response.Body)
			}

			out, _ := os.Create(filename + `.` + ext)
			defer out.Close()

			n, _ := io.Copy(out, bytes.NewReader(bodyByte))

			return n
		} else {
			response.Body.Close()
		}
	}

	/*suffixes := "avi|mpeg|3gp|mp3|mp4|wav|jpeg|gif|jpg|png|apk|exe|pdf|rar|zip|docx|doc"

	reg, _ := regexp.Compile(`(\w|\d|_)*.(` + suffixes + `)`)
	name := reg.FindStringSubmatch(url)[0]
	if name == "" {
		return 0
	}

	ext := strings.Split(name, ".")[1]

	index := strings.Index(url, name)

	//通过http请求获取图片的流文件
	resp, _ := http.Get(url[0:index] + name)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	out, _ := os.Create(filename + `.` + ext)
	defer out.Close()

	n, _ := io.Copy(out, bytes.NewReader(body))*/

	if resultErr != nil {
		WriteLog(resultErr.Error())
	}

	return 0
}

func GetPefixUrl(url string) string {

	suffixes := "avi|mpeg|3gp|mp3|mp4|wav|jpeg|gif|jpg|png|apk|exe|pdf|rar|zip|docx|doc"

	reg, _ := regexp.Compile(`(\w|\d|_)*.(` + suffixes + `)`)
	list := reg.FindStringSubmatch(url)
	if len(list) <= 0 {
		return GetPefixUrl2(url)
	}

	name := list[0]
	if name == "" {
		return GetPefixUrl2(url)
	}

	index := strings.Index(url, name)
	return url[0:index] + name
}

func GetPefixUrl2(url string) string {
	index := strings.Index(url, "%40")
	if index >= 0 {
		return url[0:index]
	} else {
		return url
	}
}

// 写超时警告日志 通用方法
func TimeoutWarning(tag, detailed string, start time.Time, timeLimit float64) {
	dis := time.Now().Sub(start).Seconds()
	if dis > timeLimit {
		//log.Warning(log.CENTER_COMMON_WARNING, tag, " detailed:", detailed, "TimeoutWarning using", dis, "s")
		//pubstr := fmt.Sprintf("%s count %v, using %f seconds", tag, count, dis)
		//stats.Publish(tag, pubstr)
		pubstr := fmt.Sprintf("%s(%s), using %f seconds", tag, detailed, dis)
		WriteLog(pubstr)
	}
}

func GetAddressLocation(region string, address string) (float64, float64, float64, error) {

	str := `http://apis.map.qq.com/ws/place/v1/search?`

	parameters := url.Values{}
	parameters.Add("keyword", address)
	parameters.Add("boundary", `region(`+region+`,0)`)
	parameters.Add("key", "LERBZ-VHXH6-TD5SB-M5P6M-YOG2Q-CFFVT")

	str += parameters.Encode()

	resp, err := http.Get(str)
	if err != nil {
		return 0, 0, -1, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, -1, errors.New("http返回statusCode:" + strconv.Itoa(resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)

	var dataMap map[string]interface{}
	err = json.Unmarshal(body, &dataMap)
	if err != nil {
		return 0, 0, -1, errors.New("返回数据JSON解析失败!")
	}

	var status float64 = -1
	if value, ok := dataMap["status"]; ok {
		status = value.(float64)
	}

	var message string = ``
	if value, ok := dataMap["message"]; ok {
		message = value.(string)
	}

	if status != 0 {
		return 0, 0, status, errors.New(message)
	}

	var data []interface{}
	if value, ok := dataMap["data"]; ok {
		data = value.([]interface{})
	} else {
		return 0, 0, status, errors.New("无法查询该地址(" + address + ")")
	}

	if len(data) == 0 {
		return 0, 0, status, errors.New("无法查询该地址(" + address + ")")
	}

	item := data[0].(map[string]interface{})

	var location map[string]interface{}
	if value, ok := item["location"]; ok {
		location = value.(map[string]interface{})
	} else {
		return 0, 0, status, errors.New("该地址没有经纬度坐标(" + address + ")")
	}

	var lat float64 = 0
	if value, ok := location["lat"]; ok {
		lat = value.(float64)
	} else {
		return 0, 0, status, errors.New("无法查询该地址lat(" + address + ")")
	}

	var lng float64 = 0
	if value, ok := location["lng"]; ok {
		lng = value.(float64)
	} else {
		return 0, 0, status, errors.New("无法查询该地址lng(" + address + ")")
	}

	return lat, lng, status, nil
}

func ReplaceSpecialChar(src string) string {

	isHave := false
	specialChar := `\/:*?"<>|` //特殊字符 windows文件名
	for i := 0; i < len(src); i++ {
		if strings.Contains(specialChar, src[i:i+1]) {
			isHave = true
			break
		}
	}

	newSrc := src
	if isHave {
		for i := 0; i < len(specialChar); i++ {
			newSrc = strings.Replace(newSrc, specialChar[i:i+1], `_`, -1)
		}
	}

	return newSrc
}

//js代码混淆加密
func JSEncode(src string) string {

	a := `62`

	c := strings.Replace(src, `\r\n`, ``, -1)
	c = strings.Replace(c, `'`, `\'`, -1)

	reg := regexp.MustCompile(`\b(\w+)\b`)
	result := reg.FindAllStringSubmatch(c, -1)
	if len(result) <= 0 {
		return src
	}

	var r []string
	for _, item := range result {
		r = append(r, item[1])
	}
	sort.Sort(sort.StringSlice(r))

	t := ``

	var p []string
	for _, item := range r {
		if item != t {
			p = append(p, item)
			t = item
		}
	}

	var ch string
	l := len(p)
	for i := 0; i < l; i++ {
		ch = num(i)
		re, _ := regexp.Compile(`\b` + regexp.QuoteMeta(p[i]) + `\b`)
		c = re.ReplaceAllString(c, ch)
		if ch == p[i] {
			p[i] = ``
		}
	}

	return "eval(function(p,a,c,k,e,d){e=function(c){return(c<a?'':e(parseInt(c/a)))+((c=c%a)>35?String.fromCharCode(c+29):c.toString(36))};if(!''.replace(/^/,String)){while(c--)d[e(c)]=k[c]||e(c);k=[function(e){return d[e]}];e=function(){return'\\\\w+'};c=1};while(c--)if(k[c])p=p.replace(new RegExp('\\\\b'+e(c)+'\\\\b','g'),k[c]);return p}(" + "'" + c + "'," + a + "," + strconv.Itoa(l) + ",'" + strings.Join(p, `|`) + "'.split('|'),0,{}))"
}

func num(c int) string {
	a := 62
	s := ``

	if c >= a {
		s = num(c / a)
	}

	c = c % a
	if c > 35 {
		s += string(c + 29)
	} else {
		s += decimalToAny(c, 36)
	}
	return s
}

// 10进制转任意进制
func decimalToAny(num, n int) string {
	new_num_str := ""
	var remainder int
	var remainder_string string
	for num != 0 {
		remainder = num % n
		if 76 > remainder && remainder > 9 {
			remainder_string = tenToAny[remainder]
		} else {
			remainder_string = strconv.Itoa(remainder)
		}
		new_num_str = remainder_string + new_num_str
		num = num / n
	}
	if len(new_num_str) == 0 {
		return `0`
	} else {
		return new_num_str
	}
}

// map根据value找key
func findkey(in string) int {
	result := -1
	for k, v := range tenToAny {
		if in == v {
			result = k
		}
	}
	return result
}

// 任意进制转10进制
func anyToDecimal(num string, n int) int {
	var new_num float64
	new_num = 0.0
	nNum := len(strings.Split(num, "")) - 1
	for _, value := range strings.Split(num, "") {
		tmp := float64(findkey(value))
		if tmp != -1 {
			new_num = new_num + tmp*math.Pow(float64(n), float64(nNum))
			nNum = nNum - 1
		} else {
			break
		}
	}
	return int(new_num)
}

//左边填充字符串
func PadLeft(s string, l int, c string) string {
	return PadString(s, l, true, c)
}

//右边填充字符串
func PadRight(s string, l int, c string) string {
	return PadString(s, l, false, c)
}

func PadString(s string, l int, f bool, c string) string {

	var result string

	for i := 0; i < l-len(s); i++ {
		result = result + c
	}

	if f {
		return result + s
	} else {
		return s + result
	}
}

func LoadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

/** 是否是图片 */
func IsPictureFormat(path string) (string, string, string) {
	temp := strings.Split(path, ".")
	if len(temp) <= 1 {
		return "", "", ""
	}
	mapRule := make(map[string]int64)
	mapRule["jpg"] = 1
	mapRule["png"] = 1
	mapRule["jpeg"] = 1
	// fmt.Println(temp[1]+"---")
	/** 添加其他格式 */
	if mapRule[temp[1]] == 1 {
		return path, temp[1], temp[0]
	} else {
		return "", "", ""
	}
}

func ImageCompress(
	getReadSizeFile func() (io.Reader, error),
	getDecodeFile func() (*os.File, error),
	to string,
	Quality,
	base int,
	format string) bool {

	/** 读取文件 */
	file_origin, err := getDecodeFile()
	defer file_origin.Close()
	if err != nil {
		WriteLog(`ImageCompress:` + err.Error())
		return false
	}
	var origin image.Image
	var config image.Config
	var temp io.Reader
	/** 读取尺寸 */
	temp, err = getReadSizeFile()
	if err != nil {
		WriteLog(`ImageCompress:` + err.Error())
		return false
	}
	var typeImage int64
	format = strings.ToLower(format)
	/** jpg 格式 */
	if format == "jpg" || format == "jpeg" {
		typeImage = 1
		origin, err = jpeg.Decode(file_origin)
		if err != nil {
			WriteLog(`ImageCompress:` + err.Error())
			return false
		}
		temp, err = getReadSizeFile()
		if err != nil {
			WriteLog(`ImageCompress:` + err.Error())
			return false
		}
		config, err = jpeg.DecodeConfig(temp)
		if err != nil {
			WriteLog(`ImageCompress:` + err.Error())
			return false
		}
	} else if format == "png" {
		typeImage = 0
		origin, err = png.Decode(file_origin)
		if err != nil {
			WriteLog(`ImageCompress:` + err.Error())
			return false
		}
		temp, err = getReadSizeFile()
		if err != nil {
			WriteLog(`ImageCompress:` + err.Error())
			return false
		}
		config, err = png.DecodeConfig(temp)
		if err != nil {
			WriteLog(`ImageCompress:` + err.Error())
			return false
		}
	}
	/** 做等比缩放 */
	//width := uint(base) /** 基准 */
	//height := uint(base * config.Height / config.Width)
	height := uint(base) /** 基准 */
	width := uint(base * config.Width / config.Height)

	canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)
	file_out, err := os.Create(to)
	defer file_out.Close()
	if err != nil {
		WriteLog(`ImageCompress:` + err.Error())
		return false
	}
	if typeImage == 0 {
		err = png.Encode(file_out, canvas)
		if err != nil {
			WriteLog(`ImageCompress:` + err.Error())
			return false
		}
	} else {
		err = jpeg.Encode(file_out, canvas, &jpeg.Options{Quality})
		if err != nil {
			WriteLog(`ImageCompress:` + err.Error())
			return false
		}
	}

	return true
}

func StringToTime(s string, f string) (t time.Time, e error) {
	//f "2006-01-02 15:04:05"//转化所需模板
	loc, _ := time.LoadLocation("Local")
	t, e = time.ParseInLocation(f, s, loc)

	return t, e
}

func MapToJson(ret map[string]interface{}) string {
	b, err := json.Marshal(ret)
	if err == nil {
		return string(b)
	} else {
		return ""
	}
}

func MD5ForStr(s string) string {

	h := md5.New()
	h.Write([]byte(s))
	b := h.Sum(nil)

	return strings.ToUpper(hex.EncodeToString(b))
}

//n 数组数量，s起始页码
func MakePageNoArray(pageNo, pageCount int) []int {

	leftRightNum := 2
	s := 1
	n := leftRightNum*2 + 1
	if pageCount > n {
		if pageNo <= leftRightNum {
			s = 1
		} else {
			s = pageNo - leftRightNum
		}
		if pageCount-s < n {
			s -= n - (pageCount - s) - 1
		}
		if s < 1 {
			s = 1
		}
	} else {
		n = pageCount
	}

	array := make([]int, n)
	for i := 0; i < n; i++ {
		array[i] = i + s
	}
	return array
}

func MapForKey(m map[string]interface{}, key string) string {

	s := ``
	if v, ok := m[key]; ok {
		s = v.(string)
	}
	return s
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
		//t.Field(i).Tag.Get(`type`);
	}

	return data
}

func SqlRow2Map(columns []string, rows *sql.Rows) map[string]interface{} {

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for i, _ := range columns {
		valuePtrs[i] = &values[i]
	}

	rows.Scan(valuePtrs...)

	var row map[string]interface{}
	row = make(map[string]interface{})

	for i, col := range columns {

		var v interface{}

		val := values[i]

		b, ok := val.([]byte)

		if ok {
			v = string(b)
		} else {
			v = val
		}

		if v == nil {
			v = ""
		}
		row[col] = v
	}
	return row
}

func SqlRows2Array(rows *sql.Rows) []map[string]interface{} {

	result := make([]map[string]interface{}, 0)

	//add row
	columns, _ := rows.Columns()
	for rows.Next() {
		result = append(result, SqlRow2Map(columns, rows))
	}

	return result
}

//读取excel的sheet到slice
func ReadXlsxMap(filepath string, isSkipHead bool) (bool, [][]string) {

	list := make([][]string, 0)

	xlsxFile, err := xlsx.OpenFile(filepath)
	if err != nil {
		WriteLog(`打开Excel文件失败！(` + filepath + `)`)

		return false, list
	}

	if len(xlsxFile.Sheets) <= 0 {
		WriteLog(`文件格式不是Excel！(` + filepath + `)`)

		return false, list
	}

	if isSkipHead {
		list = make([][]string, len(xlsxFile.Sheets[0].Rows)-1)
	} else {
		list = make([][]string, len(xlsxFile.Sheets[0].Rows))
	}

	for i, row := range xlsxFile.Sheets[0].Rows {
		data := make([]string, len(row.Cells))
		for j, cell := range row.Cells {
			data[j] = cell.String()
		}
		if isSkipHead {
			if i > 0 {
				list[i-1] = data
			}
		} else {
			list[i] = data
		}
	}
	return true, list
}

func FormatCamelCase(res []map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	s := ``

	for i := 0; i < len(res); i++ {
		row := make(map[string]interface{})
		for key, item := range res[i] {
			s = ``
			temp := strings.Split(key, `_`)
			for j := 0; j < len(temp); j++ {
				v := []rune(temp[j])
				for k := 0; k < len(v); k++ {
					if k == 0 {
						v[k] -= 32
					}
					s += string(v[k])
				}
			}
			row[s] = item
		}
		result = append(result, row)
	}
	return result
}
