package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

const serverURL = "http://localhost:8082"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "protocli",
		Short: "Protokey CLI client",
		Long:  "A command line client for interacting with the Protokey key-value store",
	}

	var setCmd = &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a key-value pair",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			setKey(args[0], args[1])
		},
	}

	var getCmd = &cobra.Command{
		Use:   "get <key>",
		Short: "Get a value by key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			getKey(args[0])
		},
	}

	var keysCmd = &cobra.Command{
		Use:   "keys <prefix>",
		Short: "Get keys matching prefix",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			getKeys(args[0])
		},
	}

	rootCmd.AddCommand(setCmd, getCmd, keysCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setKey(key, value string) {
	v, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("Value must be a valid 32-bit integer")
		os.Exit(1)
	}

	payload := map[string]interface{}{
		"key":   key,
		"value": v,
	}
	body, status, err := sendRequest("POST", "/set", payload)
	if err != nil {
		fmt.Println("Request failed:", err)
		os.Exit(1)
	}
	if status != http.StatusOK {
		fmt.Printf("Error (%d): %s\n", status, string(body))
		os.Exit(1)
	}

	fmt.Println("OK")
}

func getKey(key string) {
	body, status, err := sendRequest("GET", "/get?key="+key, nil)
	if err != nil {
		fmt.Println("Request failed:", err)
		os.Exit(1)
	}
	if status != http.StatusOK {
		fmt.Printf("Error (%d): %s\n", status, string(body))
		os.Exit(1)
	}

	var resp struct {
		Value int `json:"value"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println("Invalid response format:", err)
		os.Exit(1)
	}

	fmt.Println(resp.Value)
}

func getKeys(prefix string) {
	body, status, err := sendRequest("GET", "/keys?prefix="+prefix, nil)
	if err != nil {
		fmt.Println("Request failed:", err)
		os.Exit(1)
	}
	if status != http.StatusOK {
		fmt.Printf("Error (%d): %s\n", status, string(body))
		os.Exit(1)
	}

	var resp struct {
		Keys []string `json:"keys"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println("Invalid response format:", err)
		os.Exit(1)
	}

	for _, k := range resp.Keys {
		fmt.Println(k)
	}
}

func sendRequest(method, endpoint string, payload interface{}) ([]byte, int, error) {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, 0, fmt.Errorf("json.Marshal: %w", err)
		}
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, serverURL+endpoint, body)
	if err != nil {
		return nil, 0, fmt.Errorf("http.NewRequest: %w", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("Do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}
