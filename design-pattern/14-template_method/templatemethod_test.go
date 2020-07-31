package templatemethod

import "testing"

func ExampleHTTPDownloader() {
	var downloader Downloader = NewHTTPDownloader()
	downloader.Download("http://example.com/abc.zip")
}

func ExampleFTPDownloader() {
	var downloader Downloader = NewFTPDownloader()
	downloader.Download("ftp://example.com/abc.zip")
}

func TestDownloader(t *testing.T) {
	ExampleFTPDownloader()
	ExampleHTTPDownloader()
}
