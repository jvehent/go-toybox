package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getWeather(location string) (map[string]interface{}, error) {
	url := fmt.Sprintf("http://wttr.in/%s?format=j1", location)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %w", err)
	}
	// Ensure body is closed even if errors occur later in the function.
	// Must be after the err check for http.Get, as resp might be nil.
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := ioutil.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("failed to get weather data: status code %d", resp.StatusCode)
		if readErr == nil && len(bodyBytes) > 0 {
			// Attempt to include the API's error message if available
			errorMsg = fmt.Sprintf("%s - %s", errorMsg, string(bodyBytes))
		}
		return nil, fmt.Errorf(errorMsg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// Make sure this error is specific about JSON parsing
		return nil, fmt.Errorf("failed to parse JSON response from weather service: %w. Response body: %s", err, string(body))
	}
	return data, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: weather <location>")
		os.Exit(1)
	}
	location := os.Args[1]

	weatherData, err := getWeather(location)
	if err != nil {
		// User-friendly message for errors from getWeather
		fmt.Fprintf(os.Stderr, "Error fetching weather for '%s': %v\n", location, err)
		os.Exit(1)
	}

	formattedOutput := formatWeatherOutput(weatherData)
	fmt.Print(formattedOutput) // Use fmt.Print as formatWeatherOutput will include newlines
}

// formatWeatherOutput takes weather data and returns a user-friendly formatted string.
func formatWeatherOutput(weatherData map[string]interface{}) string {
	// Helper function to safely extract nested string values (getPath)
	// This function is defined inside formatWeatherOutput to keep it self-contained
	// as it's specific to formatting this particular data structure.
	var getPath func(data map[string]interface{}, path ...string) string
	getPath = func(data map[string]interface{}, path ...string) string {
		current := interface{}(data)
		for _, key := range path {
			if currentMap, ok := current.(map[string]interface{}); ok {
				if val, found := currentMap[key]; found {
					current = val
				} else {
					return "N/A"
				}
			} else {
				// Handle cases where path expects a map, but it's an array (e.g. "areaName": [{"value": "City"}])
				// This specific handling for "0" might need to be generalized if other paths have similar structures.
				if key == "0" {
					if currentArray, ok := current.([]interface{}); ok {
						if len(currentArray) > 0 {
							current = currentArray[0]
							continue // Continue to next key in path with the element from array
						} else {
							return "N/A" // Array is empty
						}
					}
				}
				return "N/A" // Path element not a map where expected, and not handled array case
			}
		}
		if result, ok := current.(string); ok {
			return result
		}
		if result, ok := current.(float64); ok {
			return fmt.Sprintf("%.0f", result)
		}
		if currentMap, ok := current.(map[string]interface{}); ok {
			if val, found := currentMap["value"]; found {
				if strVal, ok := val.(string); ok {
					return strVal
				}
			}
		}
		return "N/A"
	}

	// Helper to get from an array of maps, then path (getFromMapArray)
	// Also defined inside formatWeatherOutput.
	var getFromMapArray func(data map[string]interface{}, arrayKey string, index int, path ...string) string
	getFromMapArray = func(data map[string]interface{}, arrayKey string, index int, path ...string) string {
		arrayInterface, arrayFound := data[arrayKey]
		if !arrayFound {
			return "N/A"
		}
		array, ok := arrayInterface.([]interface{})
		if !ok {
			return "N/A"
		}
		if index < 0 || index >= len(array) {
			return "N/A"
		}
		itemMap, ok := array[index].(map[string]interface{})
		if !ok {
			return "N/A"
		}
		return getPath(itemMap, path...)
	}

	// Extracting weather information
	areaName := getFromMapArray(weatherData, "nearest_area", 0, "areaName", "0", "value")
	region := getFromMapArray(weatherData, "nearest_area", 0, "region", "0", "value")
	country := getFromMapArray(weatherData, "nearest_area", 0, "country", "0", "value")

	currentCondition := getFromMapArray(weatherData, "current_condition", 0, "weatherDesc", "0", "value")
	tempC := getFromMapArray(weatherData, "current_condition", 0, "temp_C")
	feelsLikeC := getFromMapArray(weatherData, "current_condition", 0, "FeelsLikeC")
	humidity := getFromMapArray(weatherData, "current_condition", 0, "humidity")
	windspeedKmph := getFromMapArray(weatherData, "current_condition", 0, "windspeedKmph")

	// Using fmt.Sprintf to build the output string
	output := fmt.Sprintf("Weather for %s, %s, %s:\n", areaName, region, country)
	output += fmt.Sprintf("Current condition: %s\n", currentCondition)
	output += fmt.Sprintf("Temperature: %s°C\n", tempC)
	output += fmt.Sprintf("Feels like: %s°C\n", feelsLikeC)
	output += fmt.Sprintf("Humidity: %s%%\n", humidity)
	output += fmt.Sprintf("Wind: %s km/h\n", windspeedKmph)
	
	return output
}
