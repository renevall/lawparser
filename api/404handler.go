package api

import(
	"net/http"
	"io/ioutil"
)

type hookedResponseWriter struct{
	http.ResponseWriter
	ignore bool
}

func (hrw *hookedResponseWriter) WriteHeader(status int){
	
	if status == 404{
		hrw.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		hrw.ResponseWriter.WriteHeader(200)
		hrw.ignore = true
		
		// hrw.ResponseWriter.Write(file)
		file, err := ioutil.ReadFile("app/index.html")
		if err != nil{
			 http.Error(hrw.ResponseWriter, err.Error(), http.StatusInternalServerError)
			 return
		}
		hrw.ResponseWriter.Write(file)
	}
	hrw.ResponseWriter.WriteHeader(status)
}

func (hrw *hookedResponseWriter) Write(p []byte) (int, error) {
    if hrw.ignore {
        return len(p), nil
    }
    return hrw.ResponseWriter.Write(p)
}

type NotFoundHook struct {
    h http.Handler
}

func (nfh NotFoundHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    nfh.h.ServeHTTP(&hookedResponseWriter{ResponseWriter: w}, r)
}