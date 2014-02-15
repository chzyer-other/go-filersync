package filersync

import (
	"runtime"
	"path"
	"strconv"
	"os"
	"fmt"
	go_log "log"
)

var loglevel = 1
var log = Logger{}
type Logger struct {}

func init() {
	for _, i := range os.Args {
		if i == "-v" { loglevel = 0 }
	}
}

func (l Logger) Infof(f string, info ...interface{}) { l.output("INFO", []interface{}{fmt.Sprintf(f, info...)}) }
func (l Logger) Info(info ...interface{}) { l.output("INFO", info) }
func (l Logger) Debug(info ...interface{}) {
	if loglevel > 0 { return }
	l.output("DEBUG", info)
}
func (l Logger) Error(info ...interface{}) { l.output("ERROR", info) }
func (l Logger) Warn(info ...interface{}) { l.output("WARN", info) }
func (l Logger) Stack() {
	a := make([]byte, 1024*1024)
	runtime.Stack(a, true)
	println(string(a))
}

func (l Logger) output(tag string, info []interface{}) {
	pc, f, line, _ := runtime.Caller(2)
	name := runtime.FuncForPC(pc).Name()
	f = path.Base(f)
	info = append([]interface{}{"["+tag+"][" + f + ":" + strconv.Itoa(line) + "]["+name+"]"}, info...)
	go_log.Println(info...)
}

