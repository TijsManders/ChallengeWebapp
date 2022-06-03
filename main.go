package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Pagina struct {
	Tafel1 bool
	Tafel2 bool
}

type TafelData struct {
	Tafel1JSON bool `json:"Tafel1Status"`
	Tafel2JSON bool `json:"Tafel2Status"`
}

var (
	Tafel1Value bool
	Tafel2Value bool
	data        Pagina
)

func main() {
	http.HandleFunc("/", RadioButtons)
	http.ListenAndServe("localhost:80", nil)
}

func VraagAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		response, err := http.Get("http://localhost:4000/get")
		if err != nil {
			fmt.Println(err)
		}
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(responseData, &data)
		// fmt.Fprintf(w, "")
		// fmt.Println(data.Tafel1, data.Tafel2)
	}
}

func StuurAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		TafelDataAPI := TafelData{
			Tafel1JSON: Tafel1Value,
			Tafel2JSON: Tafel2Value,
		}
		payloadBuf := new(bytes.Buffer)
		json.NewEncoder(payloadBuf).Encode(TafelDataAPI)
		resp, err := http.Post("http://localhost:4000/", "application/json", payloadBuf)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		// Hiermee kan de data die verstuurd wordt gelezen worden
		// body, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// sb := string(body)
		// log.Printf(sb)

	}
}

func RadioButtons(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}

		Tafel1V, err := strconv.ParseBool(r.Form.Get("Tafel1I"))
		if err != nil {
			fmt.Println(err)
		}
		Tafel1Value = Tafel1V

		Tafel2V, err := strconv.ParseBool(r.Form.Get("Tafel2I"))
		if err != nil {
			fmt.Println(err)
		}
		Tafel2Value = Tafel2V
	}
	data = Pagina{
		Tafel1: Tafel1Value,
		Tafel2: Tafel2Value,
	}
	StuurAPI(w, r)
	VraagAPI(w, r)
	tmpl.Execute(w, data)

}
