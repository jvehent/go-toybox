package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatWeatherOutput(t *testing.T) {
	mockJSONResponse := `{
		"current_condition": [
			{
				"FeelsLikeC": "10",
				"humidity": "80",
				"temp_C": "12",
				"weatherDesc": [
					{
						"value": "Sunny"
					}
				],
				"windspeedKmph": "15"
			}
		],
		"nearest_area": [
			{
				"areaName": [
					{
						"value": "TestCity"
					}
				],
				"country": [
					{
						"value": "TestCountry"
					}
				],
				"region": [
					{
						"value": "TestRegion"
					}
				]
			}
		]
	}`

	var weatherData map[string]interface{}
	err := json.Unmarshal([]byte(mockJSONResponse), &weatherData)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock JSON: %v", err)
	}

	actualOutput := formatWeatherOutput(weatherData)

	expectedSubstrings := []string{
		"Weather for TestCity, TestRegion, TestCountry:",
		"Current condition: Sunny",
		"Temperature: 12°C",
		"Feels like: 10°C",
		"Humidity: 80%",
		"Wind: 15 km/h",
	}

	for _, expected := range expectedSubstrings {
		if !strings.Contains(actualOutput, expected) {
			t.Errorf("Output does not contain expected string '%s'.\nActual output:\n%s", expected, actualOutput)
		}
	}

	// Test with missing data to ensure "N/A" is handled
	mockJSONWithMissingData := `{
		"current_condition": [
			{
				"FeelsLikeC": "10",
				"temp_C": "12"
				// Missing humidity, weatherDesc, windspeedKmph
			}
		],
		"nearest_area": [
			{
				"areaName": [
					{
						"value": "AnotherCity"
					}
				]
				// Missing country, region
			}
		]
	}`
	var weatherDataMissing map[string]interface{}
	err = json.Unmarshal([]byte(mockJSONWithMissingData), &weatherDataMissing)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock JSON with missing data: %v", err)
	}

	actualOutputMissing := formatWeatherOutput(weatherDataMissing)
	expectedSubstringsMissing := []string{
		"Weather for AnotherCity, N/A, N/A:",
		"Current condition: N/A",
		"Temperature: 12°C",
		"Feels like: 10°C",
		"Humidity: N/A%", // Expect N/A for humidity
		"Wind: N/A km/h",  // Expect N/A for wind
	}

	for _, expected := range expectedSubstringsMissing {
		if !strings.Contains(actualOutputMissing, expected) {
			t.Errorf("Output with missing data does not contain expected string '%s'.\nActual output:\n%s", expected, actualOutputMissing)
		}
	}
}
