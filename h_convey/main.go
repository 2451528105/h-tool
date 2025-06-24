package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"regexp"
)

const (
	MATCH_MP4_L = `(var mp4_l = "data:video/mp4;base64,)(.*?)(")`
	MATCH_MP4_P = `(var mp4_p = "data:video/mp4;base64,)(.*?)(")`
)

func main() {
	ConveyEndPanel("target_vedio/法语丛林尾卡_横.mp4", "target_html/德语森林尾卡.html", MATCH_MP4_L)
	ConveyEndPanel("target_vedio/法语冬季尾卡_横.mp4", "target_html/德语冬季尾卡.html", MATCH_MP4_L)
	ConveyEndPanel("target_vedio/法语魔法尾卡_横.mp4", "target_html/德语魔法尾卡.html", MATCH_MP4_L)
	ConveyEndPanel("target_vedio/法语丛林尾卡_竖.mp4", "target_html/德语森林尾卡.html", MATCH_MP4_P)
	ConveyEndPanel("target_vedio/法语冬季尾卡_竖.mp4", "target_html/德语冬季尾卡.html", MATCH_MP4_P)
	ConveyEndPanel("target_vedio/法语魔法尾卡_竖.mp4", "target_html/德语魔法尾卡.html", MATCH_MP4_P)
}

func ConveyEndPanel(targetFileName, originalFileName string, matchStr string) {
	file, err := os.ReadFile(targetFileName)
	if err != nil {
		fmt.Println(err)
	}
	base64Str := base64.StdEncoding.EncodeToString(file)
	if err != nil {
		fmt.Println(err)
	}

	file, _ = os.ReadFile(originalFileName)

	originalFile := string(file)
	// 编译正则表达式（含分组捕获）
	re := regexp.MustCompile(matchStr)

	// 执行替换（保留原字符串结构）
	replacedStr := re.ReplaceAllString(originalFile, "${1}"+base64Str+"${3}")
	//fmt.Println("替换后结果:", replacedStr)

	err = os.WriteFile(originalFileName, []byte(replacedStr), 0644)
	if err != nil {
		log.Fatal("保存失败:", err)
	}

}
