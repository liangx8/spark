package view
import (
	"net/http"
)
type View struct{
	Render func(http.ResponseWriter)
}

func HtmlWithStatus(status int,name string,data interface{}) *View{
	return &View{
		Render:func(w http.ResponseWriter){
			if status > 0 {
				w.WriteHeader(status)
			}
		},
	}
}
