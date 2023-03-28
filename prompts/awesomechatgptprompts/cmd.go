package awesomechatgptprompts

import (
	"fmt"
	"net/http"
	"encoding/csv"
)

func Exec() map[string]string {
	list := make(map[string]string)

	//// Set up a timer to run the API request every hour
	//ticker := time.NewTicker(time.Hour)
	//defer ticker.Stop()
	//
	//for ; true; <-ticker.C {
	// Make a GET request to the GitHub API to retrieve the file contents
	url := fmt.Sprintf("https://raw.githubusercontent.com/f/awesome-chatgpt-prompts/main/prompts.csv")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		//continue
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		//continue
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.LazyQuotes = true
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error parsing CSV data:", err)
		return map[string]string{}
	}

	for _, row := range data {
		list[row[0]] = row[1]
	}
	return list
}
