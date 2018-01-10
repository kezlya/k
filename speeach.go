package k

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func convert(message string) string {
	arrString := strings.Split(message, " ")
	var ret string
	for x := 0; x < len(arrString); x++ {
		ret += arrString[x] + "%20"
	}
	return ret
}

func SendWitMessage(message string) string {
	url := "https://api.wit.ai/message?v=20160225&q=" + convert(message)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer ")
	client := &http.Client{}
	resp, _ := client.Do(req)
	contents, _ := ioutil.ReadAll(resp.Body)
	return string(contents)
}

/**
*Sends an audio file to wit.ai, wit key must have been set prior to calling
*@param filename the full path to the file that is to be sent
*@return a string with the json data received
**/
func SendWitVoice(fileRef, key string) string {
	audio, err := ioutil.ReadFile(fileRef)
	if err != nil {
		log.Fatal("Error reading file:\n%v\n", err)

	}

	reader := bytes.NewReader(audio)

	url := "https://api.wit.ai/speech?v=20141022"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, reader)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "audio/wav")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
func SendWitBuff(buffer *bytes.Buffer) string {
	url := "https://api.wit.ai/speech?v=20141022"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, buffer)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer ")
	req.Header.Set("Content-Type", "audio/wav")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

/*
func ContinuousRecognition() {
	for {
		start()
	}
}
*/
func start() {
	/*
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	*/
	//cmd := "sox"
	cmd2 := "rec"
	arg2 := []string{
		//rec test.wav rate 32k silence 1 0.1 3% 1 3.0 3%
		"-t", "wav", "-",
		"rate", "32k",
		/*"rate", "16000", "channels", "1",*/
		"silence", "1", "0.1", "2%", "1", "3.0", "0.25%"}
	/*"silence", "1", "0.1", "0.1%", "1", "1.0", "0.1%"}*/
	/*args := []string{
	"-q",
	"-b", "16",
	"-d", "-t", "flac", "-",
	"rate", "16000", "channels", "1",
	//"silence", "1", "0.1", (ops.threshold || "0.1") + '%', "1", "1.0", (ops.threshold || "0.1") + '%'}
	"silence", "1", "0.1", "0.1" + "%", "1", "1.0", "0.1" + "%"}
	*/
	var byteArr []byte
	buf := bytes.NewBuffer(byteArr)
	cmdExec := exec.Command(cmd2, arg2...)
	stdout, err := cmdExec.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmdExec.Start()
	if err != nil {
		log.Fatal(err)
	}
	buf.ReadFrom(stdout)
	fmt.Println(SendWitBuff(buf))
}
