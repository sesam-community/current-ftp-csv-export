package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/webguerilla/ftps"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	server string
	user string
	password string
	ftpPort int
)

type ProjectLine struct {
	Projectref string
	Projectname string
	Duration float64
	Amount float64
	DateUtc string
}

func main(){
	wsPort := os.Getenv("WS_PORT")
	//ftp setup
	server = os.Getenv("FTP_SERVER")
	user = os.Getenv("FTP_USER")
	password =os.Getenv("FTP_PASSWORD")
	ftpPort, _ = strconv.Atoi(os.Getenv("FTP_PORT"))


	if wsPort == "" {
		wsPort = "8080"
	}

	log.Printf("Starting service on port %s", wsPort)

	router := mux.NewRouter()
	router.HandleFunc("/{topic}", PublishMessage).Methods("GET")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", wsPort), router))
}

func PublishMessage(w http.ResponseWriter, r *http.Request) {
	ftps := new(ftps.FTPS)
	ftps.TLSConfig.ServerName = server
	ftps.Debug = true

	err := ftps.Connect(server, ftpPort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ftps.Login(user, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := ftps.RetrieveFileData("currentdata.csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ftps.Quit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reader := csv.NewReader(bytes.NewReader(data))
	reader.FieldsPerRecord = -1
	reader.Comma = ';'

	csvData, err := reader.ReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var line ProjectLine
	var lines []ProjectLine

	currentDate := time.Now().UTC().Format("2006-01-02T15:04:05-0700")
	for i, each := range csvData {
		if i == 0 {//header
			continue
		}
		line.Projectref = each[0]
		line.Projectname = DecodeAnsiToUtf(each[1])
		line.Duration, _ = strconv.ParseFloat(each[2], 64)
		line.Amount, _ = strconv.ParseFloat(each[3], 64)
		line.DateUtc = currentDate
		lines = append(lines, line)
	}

	jsonData, err := json.Marshal(lines)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)


}

func DecodeAnsiToUtf(s string) string{
	var byteArr = []byte(s)
	var buf bytes.Buffer
	var r rune
	for _, b := range byteArr {
		r = rune(b)
		buf.WriteRune(r)
	}
	return string(buf.Bytes())
}
