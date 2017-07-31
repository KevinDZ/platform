package main

import (
	"bytes"
	_ "errors"
	"flag"
	"fmt"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/zalando/gin-glog"
	_ "io"
	"io/ioutil"
	_ "log"
	"net"
	"net/http"
	"net/http/httputil"
	_ "net/url"
	"runtime"
	"strings"
	"time"
)

const __LOG_HEADER_FILED__ = "__logheader__"

func LogE(c *gin.Context, str ...string) {
	header, _ := c.Get(__LOG_HEADER_FILED__)
	requestJson := c.PostForm("request")
	glog.Error(header, "request:[request="+requestJson+"]")
	glog.Error(header, str)
}

func LogI(c *gin.Context, str ...string) {
	header, _ := c.Get(__LOG_HEADER_FILED__)
	glog.Info(header, str)
}

func LogW(c *gin.Context, str ...string) {
	header, _ := c.Get(__LOG_HEADER_FILED__)
	requestJson := c.PostForm("request")
	glog.Warning(header, "request:[request="+requestJson+"]")
	glog.Warning(header, str)
}

type ResponseWriter struct {
	gin.ResponseWriter
	header string
}

func (r ResponseWriter) WriteString(response string) (int, error) {
	glog.Info(r.header + "response:[response=" + response + "]")
	return r.ResponseWriter.WriteString(response)
}

func ginLogMiddleware(c *gin.Context) {
	// header
	now := time.Now()
	timestamp := now.Format("2006-01-02|15:04:05")
	s := now.Format("20060102150405")
	ip, _ := getClientIPByRequestRemoteAddr(c.Request)
	id := s + strings.Replace(strings.Replace(ip, ".", "_", -1), ":", "_", -1)
	header := timestamp + " " + ip + " " + id + " " + c.Request.URL.Path + "\n"

	// log requset
	requestJson := c.PostForm("request")
	response := &ResponseWriter{header: header, ResponseWriter: c.Writer}
	c.Writer = response
	glog.Info(header + "request:[requesst=" + requestJson + "]")

	// for router log
	c.Set(__LOG_HEADER_FILED__, header)
	c.Next()
	/*
		  status := c.Writer.Status()
		  if status != 200 {
			glog.Error(header + "request:[requesst=" + requestJson + "]")
			glog.Error(header + "request:[requesst=" + response.response + "]")
			//glog.Error(header + "response:[response=" + response + "]")
		  }
		  //*/
}

func ginUseLogger(engine *gin.Engine) {
	flag.Parse()
	engine.Use(ginglog.Logger(3 * time.Second))
	engine.Use(ginLogMiddleware)
	engine.Use(Recovery())
	/*
	  glog.Warning("warning")
	  glog.Error("err")
	  glog.Info("info")
	  glog.V(2).Infoln("This line will be printed if you use -v=N with N >= 2.")
	*/
}

/* https://github.com/gin-gonic/gin/issues/604 */
func getClientIPByRequestRemoteAddr(req *http.Request) (ip string, err error) {

	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", err
	} else {
		return ip + ":" + port, nil
	}
}

/* Recovery */
var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return RecoveryWithWriter( /*DefaultErrorWriter*/ )
}

func RecoveryWithWriter( /*out io.Writer*/ ) gin.HandlerFunc {
	/*
		var logger *log.Logger
		if out != nil {
			logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
		}
	*/
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//if logger != nil {
				reset := ""
				stack := stack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				LogE(c, fmt.Sprintf("[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httprequest), err, stack, reset))
				//}
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//seeruntime/debug.*T·ptrmethod
	// and want
	//want*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
