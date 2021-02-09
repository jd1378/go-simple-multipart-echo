// with help from https://github.com/abvarun226/blog-source-code/blob/master/multipart-requests-in-go/multipart-related/server/main.go

package main

import (
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

// from stackoverflow
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/formdata", func(w http.ResponseWriter, r *http.Request) {
		contentType, params, parseErr := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if parseErr != nil || !strings.HasPrefix(contentType, "multipart/") {
			http.Error(w, "expecting a multipart message", http.StatusBadRequest)
			return
		}

		multipartReader := multipart.NewReader(r.Body, params["boundary"])
		defer r.Body.Close()
		log.Print("Received request, parts: \n")

		textResponse := ""
		for {
			part, err := multipartReader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, "unexpected error when retrieving a part of the message", http.StatusInternalServerError)
				return
			}
			defer part.Close()

			partBytes, err := ioutil.ReadAll(part)
			if err != nil {
				http.Error(w, "failed to read content of the part", http.StatusInternalServerError)
				return
			}
			partNameAndValue := part.FormName() + "=" + string(partBytes)
			textResponse += partNameAndValue + "&"
			log.Print(partNameAndValue + "\n")
		}

		w.Write([]byte(TrimSuffix(textResponse, "&")))
	})

	log.Print("serving on localhost:8080 \n")
	http.ListenAndServe(":8080", mux)

}
