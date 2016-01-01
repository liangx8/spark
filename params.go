package spark

import (
	"net/http"

)
type Params map[string]string
// when use this. *http.Request.Body would be touch. can read other content data anymore
func ParamsHandler(c Context,req *http.Request){
	vals := make(map[string]string)
	if err:=req.ParseForm();err != nil {
		panic(err)
	}

	for k,v := range req.Form{
		if len(v)>0 {
			vals[k]=v[0]
		}
	}
	c.Map(Params(vals))
}
