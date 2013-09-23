package main

import (
  "fmt"
	"log"
	"net/http"
	"time"
)

type handler struct{}

const firstChunk = `<!DOCTYPE html>
<head><title>Test of HTML 12.2.5.4.8 script parser spin behavior</title>
<link href="" rel="stylesheet" type="text/css"></link>
</head>
<script>
var order = "";
function foo() { order += "foo "; }
function bar() { order += "bar "; }

function go() {
  order += "go ";

  var cssURL = "http://1.cuzillion.com/bin/resource.cgi?type=css&sleep=2&n=1"
  var sheet = document.createElement('link');
  sheet.rel = "stylesheet";
  sheet.type = "text/css";
  sheet.href = cssURL;
  document.head.appendChild(sheet);
  window.setTimeout(foo, 0);
  document.write("<script>order += 'inline '</scr"+"ipt> script closed");
  bar();
}

window.onload = function() {
  order += "load ";
  window.alert(order);
}
</script>
<body>
start of &lt;body&gt;
<img src="http://localhost:8000/404.gif" onerror="go()">`

const secondChunk = `end of &lt;body&gt;
</body>
`

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, firstChunk)
  if f, ok := w.(http.Flusher); ok {
    f.Flush()
  }
  time.Sleep(5 * time.Second)
  fmt.Fprint(w, secondChunk)
}

func main() {
	var myHandler handler
	s := &http.Server{
		Addr:           ":8080",
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 0,
	}
	log.Fatal(s.ListenAndServe())
}
