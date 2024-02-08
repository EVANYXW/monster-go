package output

import (
	"fmt"
	"github.com/evanyxw/monster-go/pkg/async"
	"github.com/spf13/cast"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var Oput *Output

type Data struct {
	ConnNum int32
	GoCount int32
}

type Output struct {
	Chan       chan Data
	ServerName string
	Address    string
	RpcAddress string
	Url        string
}

// 定义颜色结构体
type Color struct {
	Foreground string
	Background string
}

// 定义颜色映射
var colors = map[string]int{
	"null":    0,
	"black":   30,
	"red":     31,
	"green":   32,
	"yellow":  33,
	"blue":    34,
	"magenta": 35,
	"cyan":    36,
	"white":   37,
}

var clear map[string]func() //create a map for storing clear func

var tamplate = `+--------------------------------------------------+
|   Server Name: {{Name}}
|   Server Addr: {{Addr}}  
|   Server Rpc Addr: {{RpcAddr}}
|
|   Player Num: {{PlayerNum}}
|   Go Num: {{GoNum}}
|   Pprof: {{Pprof}}
+--------------------------------------------------+`

func init() {
	clear = make(map[string]func()) //Initialize it

	linuxFun := clearCmd(func() *exec.Cmd {
		return exec.Command("clear") //Linux example, its tested
	})

	clear["linux"] = linuxFun
	clear["darwin"] = linuxFun
	clear["windows"] = clearCmd(func() *exec.Cmd {
		return exec.Command("cmd", "/c", "cls") //Windows example, its tested
	})
}

func NewOutput(name, addr, rpc, url string) {
	Oput = &Output{
		ServerName: name,
		Address:    addr,
		RpcAddress: rpc,
		Chan:       make(chan Data, 1),
		Url:        url,
	}

	Oput.Print(Data{0, async.GetGoCount()})
	Oput.Output()
}

func (s *Output) Update(data Data) {
	s.Chan <- data
}

func (s *Output) Output() {
	for {
		select {
		case data := <-s.Chan:
			fmt.Println(data)
			s.Clear()
			s.Print(data)
		}
	}
}

func clearCmd(fun func() *exec.Cmd) func() {
	return func() {
		cmd := fun()
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

}

// 清空终端屏幕
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func findMaxLen() {

}

func addLine(replacedText string) string {
	// 在每行末尾添加虚线
	lines := strings.Split(replacedText, "\n")
	maxLength := len(lines[0]) - 1
	for i, line := range lines {
		if i == 0 || i == len(lines)-1 {
			continue
		}
		spacesToAdd := maxLength - len(line)
		if spacesToAdd > 0 {
			lines[i] = line + strings.Repeat(" ", spacesToAdd) + "|"
		}
	}
	replacedText = strings.Join(lines, "\n")
	return replacedText
}

func (s *Output) Print(data Data) {
	str := tamplate
	str = strings.Replace(str, "{{Name}}", s.ServerName, -1)
	str = strings.Replace(str, "{{Addr}}", s.Address, -1)
	str = strings.Replace(str, "{{RpcAddr}}", s.RpcAddress, -1)
	str = strings.Replace(str, "{{PlayerNum}}", cast.ToString(data.ConnNum), -1)
	str = strings.Replace(str, "{{GoNum}}", cast.ToString(data.GoCount), -1)
	if s.Url != "" {
		str = strings.Replace(str, "{{Pprof}}", fmt.Sprintf("%s%s%s", "http://", s.Url, "/debug/pprof"), -1)
	}

	lineStr := addLine(str)
	ColorPrint("blue", "null", lineStr)

	//linkText := "Click here!"
	//linkURL := "www.example.com"

	//fmt.Printf("\033[4m%s\033[0m\n", linkText)
	//fmt.Printf("Visit %s for more information.\n", linkURL)
}

func (s *Output) Clear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func ColorPrint(foreground string, background string, text string) {
	foregroundCode := colors[foreground]      // 前景色代码
	backgroundCode := colors[background] + 10 // 背景色代码

	color := fmt.Sprintf("\033[%d;%dm", foregroundCode, backgroundCode) // 生成颜色代码
	reset := "\033[0m"                                                  // 重置颜色
	fmt.Printf("%s%s%s\n", color, text, reset)
}

func save() {
	//fmt.Println("\033[1;31mThis is red text\033[0m")     // 红色文本
	//fmt.Println("\033[1;32mThis is green text\033[0m")   // 绿色文本
	//fmt.Println("\033[1;33mThis is yellow text\033[0m")  // 黄色文本
	//fmt.Println("\033[1;34mThis is blue text\033[0m")    // 蓝色文本
	//fmt.Println("\033[1;35mThis is magenta text\033[0m") // 品红色文本
	//fmt.Println("\033[1;36mThis is cyan text\033[0m")    // 青色文本
	//
	//fmt.Println("\033[41;37mThis is red text with white background\033[0m")    // 红色文本，白色底色
	//fmt.Println("\033[42;30mThis is green text with black background\033[0m")  // 绿色文本，黑色底色
	//fmt.Println("\033[43;34mThis is yellow text with blue background\033[0m")  // 黄色文本，蓝色底色
	//fmt.Println("\033[44;33mThis is blue text with yellow background\033[0m")  // 蓝色文本，黄色底色
	//fmt.Println("\033[45;36mThis is magenta text with cyan background\033[0m") // 品红色文本，青色底色
	//fmt.Println("\033[46;31mThis is cyan text with red background\033[0m")
	fmt.Print("\x1b[31m") // 设置文本颜色为红色
	fmt.Print("This is red text")
	fmt.Print("\x1b[0m") // 重置文本颜色
	fmt.Print("\x1b[1m") // 设置文本为加粗
	fmt.Print("This is bold text")
	//fmt.Print("\x1b[0m") // 重置文本样式
	//fmt.Print("\x1b[3m") // 设置文本为斜体
	//fmt.Print("This is italic text")
	//fmt.Print("\x1b[0m") // 重置文本样式
	//fmt.Print("\x1b[4m") // 设置文本下划线
	//fmt.Print("This is underlined text")
	//fmt.Print("\x1b[0m") // 重置文本样式
}

// 设置文本颜色
func setFontColor(color string) {
	c := colors[color]
	fmt.Print(fmt.Sprintf("\x1b[%dm", c))
}

// 重置文本样式
func resetFont() {
	fmt.Print("\x1b[0m")
}

// 设置文本为加粗
func setFontWeight() {
	fmt.Print("\x1b[1m")
}

// 设置文本为斜体
func setFontXie() {
	fmt.Print("\x1b[3m")
}

// 设置文本下划线
func setFontUnderline() {
	fmt.Print("\x1b[4m")
}
