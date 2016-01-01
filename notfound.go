package spark

import (
	"net/http"
)
const (
	noFoundHtml = `<html>
<head><title>Not Found Error</title>
<style type="text/css">
html, body {
font-family: "Roboto", sans-serif;
color: #333333;
background-color: #ea5343;
margin: 0px;
}
h1 {
color: #d04526;
background-color: #ffffff;
padding: 20px;
border-bottom: 1px dashed #2b3848;
}
pre {
margin: 20px;
padding: 20px;
border: 2px solid #2b3848;
background-color: #ffffff;
}
</style>
</head><body>
<h1>Not Found Error</h1>
<pre style="font-weight: bold;">The resources you specified is not available on this server anymore</pre>

</body>
</html>`
)

func NotFound(w http.ResponseWriter,r *http.Request){
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(noFoundHtml))
}
