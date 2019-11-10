package log

import (
	"fmt"
	"os"
	// "io"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	prefix     string
	fileWriter logFileWriter
	filePath   string = ""
	nameFormat string = "2006_01_02"
)

func Init() bool {
	s := strings.Split(os.Args[0], "/")
	prefix = s[len(s)-1]
	f := fileWriter.createFile(prefix, time.Now())
	if f == nil {
		return false
	}
	fileWriter.f = f
	fileWriter.stdout = true
	log.Logger = log.Output(&fileWriter)

	zerolog.TimeFieldFormat = time.RFC3339

	go fileWriter.check()

	log.Info().Msg("log init ok")

	return true
}
func InitWithPathAndFormat(path, nameFmt string) bool {

	filePath = path
	nameFormat = nameFmt

	s := strings.Split(os.Args[0], "/")
	prefix = s[len(s)-1]
	f := fileWriter.createFile(prefix, time.Now())
	if f == nil {
		return false
	}
	fileWriter.f = f
	fileWriter.stdout = true
	log.Logger = log.Output(&fileWriter)

	zerolog.TimeFieldFormat = time.RFC3339

	go fileWriter.check()

	log.Info().Msg("log init with path ok")

	return true
}

func EnableStdOut() {
	fileWriter.stdout = true
}
func DisableStdOut() {
	fileWriter.stdout = false
}
func SetLogLevel(lv int) {
	log.Logger = log.Level(zerolog.Level(lv))
}
func Close() {
	fileWriter.Close()
}

// log
type logFileWriter struct {
	// io.Writer
	mu     sync.Mutex
	f      *os.File
	stdout bool
}

// 实现io.Writer接口
func (l *logFileWriter) Write(p []byte) (n int, e error) {
	if l.stdout {
		os.Stderr.Write(p)
	}
	tf := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&l.f)))
	f := (*os.File)(tf)
	if f != nil {
		l.mu.Lock()
		n, e := f.Write(p)
		l.mu.Unlock()
		return n, e
	}
	return 0, nil
}
func (l *logFileWriter) Close() {
	tf := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&l.f)), nil)
	f := (*os.File)(tf)
	if f != nil {
		f.Sync()
		f.Close()
	}
}

func (l *logFileWriter) createFile(prefix string, now time.Time) *os.File {
	t := now.Format(nameFormat)
	n := fmt.Sprintf("%s%s_%s.log", filePath, prefix, t)
	f, e := os.OpenFile(n, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if e != nil {
		return nil
	}
	return f
}
func (l *logFileWriter) check() {
	for {
		nd := time.Now().Add(time.Hour * 24)
		nd = time.Date(nd.Year(), nd.Month(), nd.Day(), 0, 0, 0, 0, nd.Location())
		tm := time.NewTimer(nd.Sub(time.Now()))
		tms := time.NewTimer(time.Second)
		select {
		case <-tm.C:
			{
				f := l.createFile(prefix, nd)
				if f != nil {
					oldf := l.f
					atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&l.f)), unsafe.Pointer(f))
					time.Sleep(10 * time.Second)
					oldf.Sync()
					oldf.Close()
				}
			}
		case <-tms.C:
			{
				now := time.Now()
				if !isLogExist(prefix, now) {
					f := fileWriter.createFile(prefix, now)
					if f != nil {
						atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&l.f)), unsafe.Pointer(f))
					}
				}
			}
		}
	}
}

func isLogExist(prefix string, now time.Time) bool {
	t := now.Format("2006_01_02")
	n := fmt.Sprintf("%s_%s.log", prefix, t)
	return isExist(n)
}
func isExist(path string) bool {
	_, e := os.Stat(path)
	return e == nil || os.IsExist(e)
}

func Info() *zerolog.Event {
	return log.Info()
}
func Debug() *zerolog.Event {
	return log.Debug()
}
func Error() *zerolog.Event {
	return log.Error()
}
func Warn() *zerolog.Event {
	return log.Warn()
}
func Fatal() *zerolog.Event {
	return log.Fatal()
}
func Print(v ...interface{}) {
	log.Print(v...)
}
func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
