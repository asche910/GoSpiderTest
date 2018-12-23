package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	HuabanUrl  = "http://huaban.com/partner/uc/aimeinv/pins/"
	BaseImgUrl = "http://hbimg.b0.upaiyun.com/"

	// 下载目录
	BasePath = "C:\\Users\\Asche\\go\\src\\GoSpiderTest\\imgs\\"

	// 需要翻页的次数
	PageNum = 10
)

var count = 0

func main() {
	fmt.Println("Spider Start: ------------------------------>")

	html := getHtml(HuabanUrl)

	for i := 0; i < PageNum + 1; i++ {
		id := parsePage(html)
		html = getHtml(nextPage(id))
	}
	fmt.Println("任务完成，共计下载：", count)
}

func parsePage(strHtml string) string{
	// 利用正则爬取页面中图片的key
	reg := regexp.MustCompile("\"key\":\"(.*?)\"")
	keys := reg.FindAllStringSubmatch(strHtml, -1)

	// 利用正则爬取页面中图片的pin_id
	rege := regexp.MustCompile(`"pin_id":(\d+),`)
	ids := rege.FindAllStringSubmatch(strHtml, -1)


	for i := 0; i < len(keys) ; i++ {
		key := keys[i][1]
		// 过滤掉非图片类型的key
		if len(key) < 46{
			continue
		}
		fmt.Println(parseKey(key))
		downImg(parseKey(key))
		count++
	}
	// 得到用于翻页的id
	fmt.Println(ids[len(ids)-1][1])
	return ids[len(ids)-1][1]
}

// 由id获取下一页的请求地址
func nextPage(id string) string{
	url := HuabanUrl + "?max=" + id + "&limit=8&wfl=1"
	return url
}

// 由key得到图片地址
func parseKey(key string)string{
	return BaseImgUrl + key
}

func getHtml(url string) string{
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("A error occurred!")
	}
	defer resp.Body.Close()

	htmlBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(htmlBytes)
}

// 图片下载函数
func downImg(url string) {
	// 处理前缀为//的url
	if string([]rune(url)[:2]) == "//" {
		url = "http:" + url
		fmt.Println(url)
	}

	// 解决文件无图片格式后缀的问题
	fileName := path.Base(url)
	if !strings.Contains(fileName, ".png") {
		fileName += ".png"
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("A error occurred!")
		return
	}
	defer resp.Body.Close()
	reader := bufio.NewReaderSize(resp.Body, 32*1024)

	file, _ := os.Create(BasePath + fileName)
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d\n", written)
}
