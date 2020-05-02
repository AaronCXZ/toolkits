package file

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// 文件的绝对路径
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// 文件的绝对路径
func RealPath(fp string) (string, error) {
	if path.IsAbs(fp) {
		return fp, nil
	}
	wd, err := os.Getwd()
	return path.Join(wd, fp), err
}

// 获取文件所在目录的绝对路径
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// 返回路径的最后一个元素
func Basename(fp string) string {
	return path.Base(fp)
}

// 返回路径去掉最后一个元素后的剩余部分
func Dir(fp string) string {
	return path.Dir(fp)
}

// 判断目录是否存在，不存在则创建
func InsureDir(fp string) error {
	if IsExist(fp) {
		return nil
	}
	return os.MkdirAll(fp, os.ModePerm)
}

// 创建目录
func EnsureDir(fp string) error {
	return os.MkdirAll(fp, os.ModePerm)
}

// 确保目录存在并可写
func EnsureDirRW(dataDir string) error {
	err := EnsureDir(dataDir)
	if err != nil {
		return err
	}

	checkFile := fmt.Sprintf("%s/rw.%d", dataDir, time.Now().UnixNano())
	fd, err := Create(checkFile)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("open %s: rw permission denied", dataDir)
		}
		return err
	}
	Close(fd)
	Remove(checkFile)

	return nil
}

// 判断文件或目录是否存在
func IsExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

// 创建并打开文件
func Create(name string) (*os.File, error) {
	return os.Create(name)
}

//删除文件
func Remove(name string) error {
	return os.Remove(name)
}

// 关闭文件
func Close(fd *os.File) error {
	return fd.Close()
}

// 是不是文件
func IsFile(fp string) bool {
	f, e := os.Stat(fp)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// 取消软连接删除文件
func Unlink(fp string) error {
	if !IsExist(fp) {
		return nil
	}
	return os.Remove(fp)
}

// 获取文件的mtime
func FileMTime(fp string) (int64, error) {
	f, e := os.Stat(fp)
	if e != nil {
		return 0, e
	}
	return f.ModTime().Unix(), nil
}

// 获取文件的大小
func FileSize(fp string) (int64, error) {
	f, e := os.Stat(fp)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

// 目录下的所有目录
func DirsUnder(dirPath string) ([]string, error) {
	if !IsExist(dirPath) {
		return []string{}, nil
	}

	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	ret := make([]string, 0, sz)
	for i := 0; i < sz; i++ {
		if fs[i].IsDir() {
			name := fs[i].Name()
			if name != "." && name != ".." {
				ret = append(ret, name)
			}
		}
	}

	return ret, nil
}

// 目录下的所有文件
func FilesUnder(dirPath string) ([]string, error) {
	if !IsExist(dirPath) {
		return []string{}, nil
	}

	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	ret := make([]string, 0, sz)
	for i := 0; i < sz; i++ {
		if !fs[i].IsDir() {
			ret = append(ret, fs[i].Name())
		}
	}

	return ret, nil
}

// 读取文件内容，[]byte类型
func ReadBytes(cpath string) ([]byte, error) {
	if !IsExist(cpath) {
		return nil, fmt.Errorf("%s not exists", cpath)
	}

	if !IsFile(cpath) {
		return nil, fmt.Errorf("%s not file", cpath)
	}

	return ioutil.ReadFile(cpath)
}

// 读取文件内容，string类型
func ReadString(cpath string) (string, error) {
	bs, err := ReadBytes(cpath)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

// 读取文件内容，删除首尾的空白
func ReadStringTrim(cpath string) (string, error) {
	out, err := ReadString(cpath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

// 读取yaml文件
func ReadYaml(cpath string, cptr interface{}) error {
	bs, err := ReadBytes(cpath)
	if err != nil {
		return fmt.Errorf("cannot read %s: %s", cpath, err.Error())
	}

	err = yaml.Unmarshal(bs, cptr)
	if err != nil {
		return fmt.Errorf("cannot parse %s: %s", cpath, err.Error())
	}

	return nil
}

// 读取json文件
func ReadJson(cpath string, cptr interface{}) error {
	bs, err := ReadBytes(cpath)
	if err != nil {
		return fmt.Errorf("cannot read %s: %s", cpath, err.Error())
	}

	err = json.Unmarshal(bs, cptr)
	if err != nil {
		return fmt.Errorf("cannot parse %s: %s", cpath, err.Error())
	}

	return nil
}

// 将byte写入到文件
func WriteBytes(filePath string, b []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}

// 将string写入文件
func WriteString(filePath string, s string) (int, error) {
	return WriteBytes(filePath, []byte(s))
}

// 下载文件
func Download(toFile, url string) error {
	out, err := os.Create(toFile)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.Body == nil {
		return fmt.Errorf("%s response body is nil", url)
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// 文件内容做md5
func MD5(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// 打开文件，追加模式
func OpenLogFile(fp string) (*os.File, error) {
	return os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

// 读取一行
func ReadLineBytes(r *bufio.Reader) ([]byte, error) {
	return r.ReadBytes('\n')
}

func ReadLineString(r *bufio.Reader) (string, error) {
	return r.ReadString('\n')
}
