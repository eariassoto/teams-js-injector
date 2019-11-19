package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/net/websocket"
)

const (
	TeamsPath = "C:\\Users\\%s\\AppData\\Local\\Microsoft\\Teams\\current\\Teams.exe"
)

func killAllTeamsProcesses() {
	cmnd := exec.Command("taskkill", "/f", "/im", "Teams.exe")
	cmnd.Start()
	cmnd.Wait()
}

func launchTeamProcess(path string, debugPort int) {
	cmnd := exec.Command(path, fmt.Sprintf("--remote-debugging-port=%d", debugPort))
	cmnd.Start()
}

func getChatWsDebuggerURL(debugPort int) string {
	log.Println("Make sure you have a chat windows open")
	urlFound := false
	var url string
	for !urlFound {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/json", debugPort))
		if err != nil {
			continue
		}
		var result []map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&result)
		for i := 0; i < len(result); i++ {
			if result[i]["title"] == "Chat | Microsoft Teams" {
				url = result[i]["webSocketDebuggerUrl"].(string)
				urlFound = true
			}
		}
	}
	return url
}

func makeRequestMessage(payload string) map[string]interface{} {
	request := map[string]interface{}{}
	request["id"] = 1337
	request["method"] = "Runtime.evaluate"

	payloadStruct := map[string]interface{}{}
	payloadStruct["expression"] = payload

	request["params"] = payloadStruct

	return request
}

func sendMessage(wsURL string, message map[string]interface{}) {
	b, err := json.Marshal(message)
	log.Println(string(b))
	if err != nil {
		log.Fatalln(err)
	}

	origin := "http://localhost/"
	conn, err := websocket.Dial(wsURL, "", origin)
	if err != nil {
		log.Panicln(err)
	}
	defer conn.Close()

	if err = websocket.JSON.Send(conn, message); err != nil {
		log.Panicln(err)
	}
}

func readFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	fileStr := string(b)

	return fileStr
}

func main() {
	debugPort := flag.Int("debug-port", 9222, "Port number for Chromium remote debugging")

	user := os.Getenv("USERNAME")
	defaultPath := fmt.Sprintf(TeamsPath, user)
	teamsExePath := flag.String("teams-path", defaultPath, "Location of Teams executable")

	payloadFile := flag.String("payload-file", "payload.js", "Javascript file to inject")
	flag.Parse()

	killAllTeamsProcesses()

	launchTeamProcess(*teamsExePath, *debugPort)

	wsURL := getChatWsDebuggerURL(*debugPort)

	payload := readFile(*payloadFile)
	log.Println(payload)
	request := makeRequestMessage(payload)

	sendMessage(wsURL, request)

	log.Println("Program finished, you can now close this windows")
}
