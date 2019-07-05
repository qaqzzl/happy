/**
默认加载项目下config.conf文件
自定义设置配置文件请使用InitConfig方法
*/
package happy

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const middle = "."
var CONFIG = make(map[string]string)

func init() {
	initConf("config.conf")
}

func InitConfig(path string) {
	initConf(path)
}

func ConfigGet(node, key string) string {
	key = node + middle + key
	v, found := CONFIG[key]
	if !found {
		return ""
	}
	return v
}

func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

func initConf(path string) {
	//c.Mymap = make(map[string]string)
	var strcet string
	f, err := os.Open(path)
	if err != nil { //抛出错误信息
		//panic("配置文件读取失败")
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := strcet + middle + frist
		CONFIG[key] = strings.TrimSpace(second)
	}
}
