package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func getData(input string, chatId string, configDir string, isInteractive bool) (serverChatId string) {
	// proxyUrl, _ := url.Parse("http://127.0.0.1:8080")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// Proxy:           http.ProxyURL(proxyUrl),
	}
	client := &http.Client{Transport: tr}
	var data = strings.NewReader(fmt.Sprintf(`{"prompt":"%v","options":{"parentMessageId":"%v"}}`, input, chatId))
	req, err := http.NewRequest("POST", "https://chatbot.theb.ai/api/chat-process", data)
	if err != nil {
		fmt.Println("\nSome error has occured.")
		fmt.Println("Error:", err)
		os.Exit(0)
	}
	// Setting all the required headers
	req.Header.Set("Host", "chatbot.theb.ai")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/112.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://chatbot.theb.ai")
	req.Header.Set("Referer", "https://chatbot.theb.ai/")
	resp, err := client.Do(req)
	if err != nil {
		stopSpin = true
		bold.Println("\rSome error has occured. Check your internet connection.")
		fmt.Println("\nError:", err)
		os.Exit(0)
	}
	code := resp.StatusCode

	defer resp.Body.Close()

	stopSpin = true
	fmt.Print("\r")

	scanner := bufio.NewScanner(resp.Body)

	// Variables
	var oldLine = ""
	var newLine = ""
	count := 0
	isCode := false
	isGreen := false
	tickCount := 0
	previousWasTick := false
	isTick := false

	// Print the Question
	if !isInteractive {
		fmt.Print("\r         ")
		bold.Printf("\r%v\n\n", input)
	} else {
		fmt.Println()
	}

	gotId := false
	id := ""
	// Handling each json
	for scanner.Scan() {
		var jsonObj map[string]interface{}
		line := scanner.Text()
		err := json.Unmarshal([]byte(line), &jsonObj)
		if err != nil {
			bold.Println("\rError. Your IP is being blocked by the server.")
			fmt.Println("Status Code:", code)
			os.Exit(0)
		}

		mainText := fmt.Sprintf("%s", jsonObj["text"])
		id = fmt.Sprintf("%s", jsonObj["id"])

		if !gotId {
			gotId = true
		}

		if count <= 0 {
			oldLine = mainText
			splitLine := strings.Split(oldLine, "")
			// Iterating through each word
			for _, word := range splitLine {
				// If its a backtick
				if word == "`" {
					tickCount++
					isTick = true

					if tickCount == 2 && !previousWasTick {
						tickCount = 0
					} else if tickCount == 6 {
						tickCount = 0
					}
					previousWasTick = true
					isGreen = false
					isCode = false

				} else {
					isTick = false
					// If its a normal word

					if tickCount == 1 {
						isGreen = true
					} else if tickCount == 3 {
						isCode = true
					}
					previousWasTick = false
				}

				if isCode {
					codeText.Print(word)
				} else if isGreen {
					boldBlue.Print(word)
				} else if !isTick {
					fmt.Print(word)
				}
			}
		} else {
			newLine = mainText
			result := strings.Replace(newLine, oldLine, "", -1)
			splitLine := strings.Split(result, "")
			whiteSpaceFound := false

			for _, word := range splitLine {
				// If its a backtick
				if word == "`" {
					tickCount++
					isTick = true

					if tickCount == 2 && !previousWasTick {
						tickCount = 0
					} else if tickCount == 6 {
						tickCount = 0
					}
					previousWasTick = true
					isGreen = false
					isCode = false

				} else {
					isTick = false
					if word == "\n" {
						whiteSpaceFound = true

					} else {
						whiteSpaceFound = false
					}
					// If its a normal word
					if tickCount == 1 {
						isGreen = true
					} else if tickCount == 3 {
						if previousWasTick {
						} else {
							if whiteSpaceFound {
								isCode = true
							}
						}
					}
					previousWasTick = false
				}

				if isCode {
					codeText.Print(word)
				} else if isGreen {
					boldBlue.Print(word)
				} else if !isTick {
					fmt.Print(word)
				}

			}
			oldLine = newLine
		}

		count++
	}
	fmt.Println("")
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	createConfig(configDir, id)
	return id
}

func createConfig(dir string, chatId string) {
	err := os.MkdirAll(dir, 0755)
	configTxt := "id:" + chatId
	if err != nil {
		fmt.Println(err)
	} else {
		os.WriteFile(dir+"/config.txt", []byte(configTxt), 0755)
	}
}

func loading(stop *bool) {
	spinChars := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		if *stop {
			break
		}
		fmt.Printf("\r%s Loading", spinChars[i])
		i = (i + 1) % len(spinChars)
		time.Sleep(100 * time.Millisecond)
	}
}