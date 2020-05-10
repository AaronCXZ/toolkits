package logger

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	// TypeTimeBaseRotate is time base logrotate tag
	TypeTimeBaseRotate = "time"
	// TypeSizeBaseRotate is size base logrotate tag
	TypeSizeBaseRotate = "size"
)

var (
	// ErrConfigFiltersNotFound not found filters
	ErrConfigFiltersNotFound = errors.New("Please define at least one filter")
	// ErrConfigBadAttributes wrong attribute
	ErrConfigBadAttributes = errors.New("Bad attributes setting")
	// ErrConfigLevelsNotFound not found levels
	ErrConfigLevelsNotFound = errors.New("Please define levels attribution")
	// ErrConfigFilePathNotFound not found file path
	ErrConfigFilePathNotFound = errors.New("Please define the file path")
	// ErrConfigFileRotateTypeNotFound not found rotate type
	ErrConfigFileRotateTypeNotFound = errors.New("Please define the file rotate type")
	// ErrConfigSocketAddressNotFound not found socket address
	ErrConfigSocketAddressNotFound = errors.New("Please define a socket address")
	// ErrConfigSocketNetworkNotFound not found socket port
	ErrConfigSocketNetworkNotFound = errors.New("Please define a socket network type")
)

// Config struct define the config struct used for file wirter
type Config struct {
	Filters  []filter `xml:"filter" yaml:"filter"`
	MinLevel string   `xml:"minlevel,attr" yaml:"minlevel"`
}

// log filter
type filter struct {
	Levels     string     `xml:"levels,attr" yaml:"levels"`
	Colored    bool       `xml:"colored,attr" yaml:"colored"`
	File       file       `xml:"file" yaml:"file"`
	RotateFile rotateFile `xml:"rotatefile" yaml:"rotatefile"`
	Console    console    `xml:"console" yaml:"console"`
	Socket     socket     `xml:"socket" yaml:"socket"`
}

type file struct {
	Path string `xml:"path,attr" yaml:"path"`
}

type rotateFile struct {
	Path        string `xml:"path,attr" yaml:"path"`
	Type        string `xml:"type,attr" yaml:"type"`
	RotateLines int    `xml:"rotateLines,attr" yaml:"rotateLines"`
	RotateSize  int64  `xml:"rotateSize,attr" yaml:"rotateSize"`
	Retentions  int64  `xml:"retentions,attr" yaml:"retentions"`
}

type console struct {
	// redirect stderr to stdout
	Redirect bool `xml:"redirect" yaml:"redirect"`
}

type socket struct {
	Network string `xml:"network,attr" yaml:"network"`
	Address string `xml:"address,attr" yaml:"address"`
}

// check if config is valid
func (config *Config) valid() error {
	// check minlevel validation
	if "" != config.MinLevel && !LevelFromString(config.MinLevel).valid() {
		return ErrConfigBadAttributes
	}

	// check filters len
	if len(config.Filters) < 1 {
		return ErrConfigFiltersNotFound
	}

	// check filter one by one
	for _, filter := range config.Filters {
		if "" == filter.Levels {
			return ErrConfigLevelsNotFound
		}

		if (file{}) != filter.File {
			// seem not needed now
			//if "" == filter.File.Path {
			//return ErrConfigFilePathNotFound
			//}
		} else if (rotateFile{}) != filter.RotateFile {
			if "" == filter.RotateFile.Path {
				return ErrConfigFilePathNotFound
			}

			if "" == filter.RotateFile.Type {
				return ErrConfigFileRotateTypeNotFound
			}
		} else if (socket{}) != filter.Socket {
			if "" == filter.Socket.Address {
				return ErrConfigSocketAddressNotFound
			}

			if "" == filter.Socket.Network {
				return ErrConfigSocketNetworkNotFound
			}
		}
	}

	return nil
}

// read config from a xml file
func readConfig(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if nil != err {
		return nil, err
	}
	defer file.Close()

	in, err := ioutil.ReadAll(file)
	if nil != err {
		return nil, err
	}

	config := new(Config)
	fileType := ConfigFileType(fileName)
	switch fileType {
	case "xml":
		err = xml.Unmarshal(in, config)
	case "json":
		err = json.Unmarshal(in, config)
	case "yaml", "yml":
		err = yaml.Unmarshal(in, config)
	}
	if nil != err {
		return nil, err
	}

	return config, err
}

func ConfigFileType(fileName string) string {
	names := strings.Split(fileName, ".")
	size := len(names)
	return names[size-1]
}

func NewConfig(minlevel string, filters ...filter) *Config {
	return &Config{Filters: filters, MinLevel: minlevel}
}

func NewFilter(level string, color bool, arg interface{}) filter {
	var f filter
	switch t := arg.(type) {
	case file:
		f = newFileFilter(level, color, t)
	case rotateFile:
		f = newRotateFileFilter(level, color, t)
	case console:
		f = newConsoleFilter(level, color, t)
	case socket:
		f = newSocketFilter(level, color, t)
	}
	return f
}

func newFileFilter(level string, color bool, arg file) filter {
	return filter{
		Levels:  level,
		Colored: color,
		File:    arg,
	}
}

func newRotateFileFilter(level string, color bool, arg rotateFile) filter {
	return filter{
		Levels:     level,
		Colored:    color,
		RotateFile: arg,
	}
}

func newConsoleFilter(level string, color bool, arg console) filter {
	return filter{
		Levels:  level,
		Colored: color,
		Console: arg,
	}
}

func newSocketFilter(level string, color bool, arg socket) filter {
	return filter{
		Levels:  level,
		Colored: color,
		Socket:  arg,
	}
}

func NewFile(f string) file {
	return file{
		Path: f,
	}
}

func NewTimeRotateFile(path string, retentions int64) rotateFile {
	return rotateFile{
		Path:       path,
		Type:       "time",
		Retentions: retentions,
	}
}

func NewSizeRotateFile(path string, rotateSize int64, rotateLines int) rotateFile {
	return rotateFile{
		Path:        path,
		Type:        "size",
		RotateSize:  rotateSize,
		RotateLines: rotateLines,
	}
}

func NewConsole(redirect bool) console {
	return console{
		Redirect: redirect,
	}
}

func NewSocket(network, address string) socket {
	return socket{
		Network: network,
		Address: address,
	}
}
