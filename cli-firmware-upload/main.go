package main

import (
	"crypto/md5"
	"encoding/hex"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	md5Hash := hex.EncodeToString(hashInBytes)

	return md5Hash, nil
}

func setHash(){
	
	
		filePath := "../.pio/build/nodemcuv2/firmware.bin"
		hash, err := calculateMD5(filePath)
		url := "http://192.168.0.6/ota/start?mode=fr&hash=" + hash
	
		// Create the request
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
	
		// Set the necessary headers
		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0")
		request.Header.Set("Accept", "*/*")
		request.Header.Set("Accept-Language", "en-US,en;q=0.5")
		request.Header.Set("Accept-Encoding", "gzip, deflate")
		request.Header.Set("Referer", "http://192.168.0.6/update")
		request.Header.Set("Connection", "keep-alive")
	
		// Make the request
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer response.Body.Close()
	
		// Print the response status and headers
		fmt.Println("Response Status:", response.Status)
		fmt.Println("Response Headers:", response.Header)
	
		// You can read and print the response body if needed
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		fmt.Println("Response Body:", string(body))

}

func uploadFirmware(){
	url := "http://192.168.0.6/ota/upload"
	filePath := "../.pio/build/nodemcuv2/firmware.bin"

	// Create a buffer to store the request body
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the file to the request body
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fileWriter, err := writer.CreateFormFile("file", "firmware.bin")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		fmt.Println("Error copying file content:", err)
		return
	}

	// Close the multipart writer to finalize the request body
	writer.Close()

	// Create the request
	request, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the necessary headers
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Accept-Language", "en-US,en;q=0.9")
	request.Header.Set("Referer", "http://192.168.0.6/update")
	request.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	request.Header.Set("Accept-Encoding", "gzip, deflate")


	// Make the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer response.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Response:", string(body))
}



func main() {
	setHash()
	uploadFirmware()
}
