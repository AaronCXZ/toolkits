// 模板方法模式
/*
模版方法模式使用继承机制，把通用步骤和通用方法放到父类中，把具体实现延迟到子类中实现。使得实现符合开闭原则。
如实例代码中通用步骤在父类中实现（准备、下载、保存、收尾）下载和保存的具体实现留到子类中，并且提供 保存方法的默认实现。
因为Golang不提供继承机制，需要使用匿名组合模拟实现继承。
此处需要注意：因为父类需要调用子类方法，所以子类需要匿名组合父类的同时，父类需要持有子类的引用。
*/

package templatemethod

import "fmt"

type Downloader interface {
	Download(uri string)
}

type template struct {
	implement
	uri string
}

type implement interface {
	download()
	save()
}

func NewTemplate(impl implement) *template {
	return &template{implement: impl}
}

func (t *template) Download(uri string) {
	t.uri = uri
	fmt.Println("prepare downloading")
	t.implement.download()
	t.implement.save()
	fmt.Println("finish downloading")
}

func (t *template) save() {
	fmt.Println("default save")
}

type HTTPDownloader struct {
	*template
}

func NewHTTPDownloader() Downloader {
	downloader := &HTTPDownloader{}
	template := NewTemplate(downloader)
	downloader.template = template
	return downloader
}

func (d *HTTPDownloader) download() {
	fmt.Printf("download %s via http\n", d.uri)
}

func (d *HTTPDownloader) save() {
	fmt.Println("http save")
}

type FTPDownloader struct {
	*template
}

func NewFTPDownloader() Downloader {
	downloader := &FTPDownloader{}
	template := NewTemplate(downloader)
	downloader.template = template
	return downloader
}

func (d *FTPDownloader) download() {
	fmt.Printf("download %s via ftp\n", d.uri)
}
