package rollingwriter

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"gopkg.in/yaml.v2"
)

// 当WriterMode为none时使用的结构，无保护的writer: 不提供并发安全保障
type Writer struct {
	m             Manager
	file          *os.File // 当前的写入文件
	absPath       string
	fire          chan string
	cf            *Config
	rollingfilech chan string //
}

// 当WriterMode为lock时使用的结构，lock保护的writer: 提供由mutex保护的并发安全保障
type LockedWriter struct {
	Writer
	sync.Mutex
}

// 当WriterMode为async时使用的结构，同步writer，并发安全
type AsynchronousWriter struct {
	Writer
	ctx     chan int    // 有数据时退出写入
	queue   chan []byte // 缓存队列chan
	errChan chan error  // 数据写入错误chan
	closed  int32       // 默认为：0，当关闭时为：1
	wg      sync.WaitGroup
}

// 当WriterMode为buffer时使用的结构，异步write, 并发安全
type BufferWriter struct {
	Writer
	buf     *[]byte // 待写入数据
	swaping int32   // 缓存池中数据是否处理完的标志，默认为：0，没处理完为：1
}

// 同步并发是使用的临时缓存对象
var _asyncBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, BufferSize)
	},
}

// 根据配置生成RollingWriter，用于接收日志输入
func NewWriterFromConfig(c *Config) (RollingWriter, error) {
	// 判断配置
	if c.LogPath == "" || c.FileName == "" {
		return nil, ErrInvalidArgument
	}

	// 创建日志所在目录
	if err := os.MkdirAll(c.LogPath, 0700); err != nil {
		return nil, err
	}

	filepath := LogFilePath(c)
	// 打开日志文件
	file, err := os.OpenFile(filepath, DefualtFileFlag, DefualtFileMode)
	if err != nil {
		return nil, err
	}
	mng, err := NewManager(c)
	if err != nil {
		return nil, err
	}
	var rollingWriter RollingWriter
	writer := Writer{
		m:       mng,
		file:    file,
		absPath: filepath,
		fire:    mng.Fire(), // 最新的历史文件名称
		cf:      c,
	}

	if c.MaxRemain > 0 {
		// 保留历史日志文件名称的chan
		writer.rollingfilech = make(chan string, c.MaxRemain)
		// 读取日志目录下的所有文件信息
		dir, err := ioutil.ReadDir(c.LogPath)
		if err != nil {
			return nil, err
		}
		files := make([]string, 0, 10)
		// 查找日志目录中的历史日志文件
		for _, fi := range dir {
			if fi.IsDir() {
				continue
			}

			fileName := c.FileName + ".log"
			if strings.Contains(fi.Name(), fileName) {
				// 文件的后缀名
				fileSuffix := path.Ext(fi.Name())
				if len(fileSuffix) > 1 {
					// 将后缀字符串转化为时间，如果成功的为历史日志文件，添加到files
					_, err := time.Parse(c.TimeTagFormat, fileSuffix[1:])
					if err == nil {
						files = append(files, fi.Name())
					}
				}
			}
		}

		// 将files按照时间排序
		sort.Slice(files, func(i, j int) bool {
			fileSuffix1 := path.Ext(files[i])
			fileSuffix2 := path.Ext(files[j])
			t1, _ := time.Parse(c.TimeTagFormat, fileSuffix1[1:])
			t2, _ := time.Parse(c.TimeTagFormat, fileSuffix2[1:])
			return t1.Before(t2)
		})

		// 删除多余的历史日志文件
		for _, file := range files {
		retry:
			select {
			case writer.rollingfilech <- path.Join(c.LogPath, file):
			default:
				writer.DoRemove()
				goto retry
			}
		}
	}

	// 判断日志写入模式
	switch c.WriterMode {
	case "none":
		rollingWriter = &writer
	case "lock":
		rollingWriter = &LockedWriter{
			Writer: writer,
		}
	case "async":
		wr := &AsynchronousWriter{
			Writer:  writer,
			ctx:     make(chan int),
			queue:   make(chan []byte, QueueSize),
			errChan: make(chan error),
			closed:  0,
			wg:      sync.WaitGroup{},
		}
		wr.wg.Add(1)
		go wr.writer()
		wr.wg.Wait()
		rollingWriter = wr
	case "buffer":
		bf := make([]byte, 0, c.BufferWriterThreshold*2)
		rollingWriter = &BufferWriter{
			Writer:  writer,
			buf:     &bf,
			swaping: 0,
		}
	default:
		return nil, ErrInvalidArgument
	}
	return rollingWriter, nil
}

// 执行各个构造函数更新配置后生产RollingWriter
func NewWriter(ops ...Option) (RollingWriter, error) {
	cfg := NewDefaultConfig()
	for _, opt := range ops {
		opt(&cfg)
	}
	return NewWriterFromConfig(&cfg)
}

// 从配置文件读取配置,解析后生成RollingWriter,支持json和yaml类型
func NewWriterFromConfigFile(path string, typ string) (RollingWriter, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := NewDefaultConfig()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	switch typ {
	case "json":
		if err = json.Unmarshal(buf, &cfg); err != nil {
			return nil, err
		}
	case "yaml":
		if err = yaml.Unmarshal(buf, &cfg); err != nil {
			return nil, err
		}
	}

	return NewWriterFromConfig(&cfg)
}

// 删除过期的历史日志文件
func (w *Writer) DoRemove() {
	select {
	case file := <-w.rollingfilech:
		if err := os.Remove(file); err != nil {
			log.Println("error in remove log file", file, err)
		}
	}
}

// 压缩历史文件
func (w *Writer) CompressFile(oldfile *os.File, cmpname string) error {
	cmpfile, err := os.OpenFile(cmpname, DefualtFileFlag, DefualtFileMode)
	defer cmpfile.Close()
	if err != nil {
		return err
	}
	gw := gzip.NewWriter(cmpfile)
	defer gw.Close()

	// 设置下次读取oldfile文件时的偏移量，及从头开始读取oldfile到压缩文件
	if _, err := oldfile.Seek(0, 0); err != nil {
		return err
	}

	if _, err := io.Copy(gw, oldfile); err != nil {
		// 当压缩失败时删除压缩文件
		if errR := os.Remove(cmpname); errR != nil {
			return errR
		}
		return err
	}
	// 删除临时文件
	return os.Remove(cmpname + ".tmp")
}

// 执行日志滚动， file为生成的历史文件名称
func (w *Writer) Reopen(file string) error {
	// 重命名
	if err := os.Rename(w.absPath, file); err != nil {
		return err
	}
	// 打开新的日志文件
	newfile, err := os.OpenFile(w.absPath, DefualtFileFlag, DefualtFileMode)
	if err != nil {
		return err
	}

	// 原子性的将新打开的日志文件替换就日志文件，并返回就日志文件
	// 使用unsafe.Pointer直接操作了正在写入日志文件的指针
	// oldfile的指针指向最新生成的历史日志文件
	oldfile := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&w.file)), unsafe.Pointer(newfile))

	go func() {
		defer (*os.File)(oldfile).Close()
		// 执行历史日志文件压缩
		if w.cf.Compress {
			if err := os.Rename(file, file+".tmp"); err != nil {
				log.Println("error in compress rename tempfile", err)
				return
			}
			if err := w.CompressFile((*os.File)(oldfile), file); err != nil {
				log.Println("error in compress log file", err)
				return
			}
		}

		// 删除过期历史日志文件
		if w.cf.MaxRemain > 0 {
		retry:
			select {
			case w.rollingfilech <- file:
			default:
				w.DoRemove()
				goto retry
			}
		}
	}()
	return nil
}

// 没有lock的Write接口实现
func (w *Writer) Write(b []byte) (int, error) {
	select {
	// 触发日志滚动
	case filename := <-w.fire:
		if err := w.Reopen(filename); err != nil {
			return 0, err

		}
	default:

	}
	// 原子性的获取当前写入日志文件的指针
	fp := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&w.file)))
	file := (*os.File)(fp)
	return file.Write(b)
}

// 使用lock的Write接口实现
func (w *LockedWriter) Write(b []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	select {
	// 触发日志滚动
	case filename := <-w.fire:
		if err := w.Reopen(filename); err != nil {
			return 0, err
		}
	default:

	}
	n, err = w.file.Write(b)
	return n, err
}

// 同步并发的Write接口实现
func (w *AsynchronousWriter) Write(b []byte) (int, error) {
	if atomic.LoadInt32(&w.closed) == 0 {
		select {
		case err := <-w.errChan:
			return 0, err
		// 触发日志滚动
		case filename := <-w.fire:
			if err := w.Reopen(filename); err != nil {
				return 0, err
			}

			l := len(b)
			// 将数据写入_asyncBufferPool
			for len(b) > 0 {
				buf := _asyncBufferPool.Get().([]byte)
				n := copy(buf, b)
				w.queue <- buf[:n]
				b = b[n:]
			}
			return l, nil
		default:
			w.queue <- append(_asyncBufferPool.Get().([]byte)[0:0], b...)[:len(b)]
			return len(b), nil
		}
	}
	return 0, ErrClosed
}

// 异步并发的Write接口实现
func (w *BufferWriter) Write(b []byte) (int, error) {
	select {
	// 触发日志滚动
	case filename := <-w.fire:
		if err := w.Reopen(filename); err != nil {
			return 0, err
		}
	default:

	}
	// 读取所有待写入的数据
	buf := append(*w.buf, b...)
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&w.buf)), (unsafe.Pointer)(&buf))
	// 判断待写入数据大于缓存池，并且w.swaping==0，并且设置w.swaping=1
	if len(*w.buf) > w.cf.BufferWriterThreshold && atomic.CompareAndSwapInt32(&w.swaping, 0, 1) {
		// 创建新缓存池
		nb := make([]byte, 0, w.cf.BufferWriterThreshold*2)
		// 新缓存池代替旧缓存池，并返回就缓存池的指针
		ob := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&w.buf)), (unsafe.Pointer(&nb)))
		// 写入就缓存池中的数据
		w.file.Write(*(*[]byte)(ob))
		// 设置w.swaping=0
		atomic.StoreInt32(&w.swaping, 0)
	}
	return len(b), nil
}

// 没有lock的Close接口实现，借助atomic实现原子性操作
func (w *Writer) Close() error {
	return (*os.File)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&w.file)))).Close()
}

// 使用lock的Close接口实现
func (w *LockedWriter) Close() error {
	w.Lock()
	defer w.Unlock()
	return w.file.Close()
}

// 同步并发的Close接口实现
func (w *AsynchronousWriter) Close() error {
	// w.closed==0，并设置w.closed=1
	if atomic.CompareAndSwapInt32(&w.closed, 0, 1) {
		close(w.ctx)
		w.onClose()
		return w.file.Close()
	}
	return ErrClosed
}

// 将缓存队列中的数据处理完
func (w *AsynchronousWriter) onClose() {
	var err error
	for {
		select {
		case b := <-w.queue:
			if _, err = w.file.Write(b); err != nil {
				select {
				case w.errChan <- err:
				default:
					_asyncBufferPool.Put(b)
					return
				}
			}
			_asyncBufferPool.Put(b)
		default:
			return
		}
	}
}

// 同步并发是的数据写入
func (w *AsynchronousWriter) writer() {
	var err error
	w.wg.Done()
	for {
		select {
		case b := <-w.queue:
			if _, err = w.file.Write(b); err != nil {
				w.errChan <- err
			}
			_asyncBufferPool.Put(b)
		case <-w.ctx:
			return
		}
	}
}

// 异步并发的Close接口实现
func (w BufferWriter) Close() error {
	_, err := w.file.Write(*w.buf)
	if err != nil {
		return err
	}
	return w.file.Close()
}
