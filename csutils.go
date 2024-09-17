package csutils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// CloudSecureInfo struct holds the API credentials and tenant ID
type CloudSecureInfo struct {
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
	TenantID  string `json:"tenantID"`
}

// CloudSecureConfig struct holds the configuration
type CloudSecureConfig struct {
	CloudSecures     map[string]CloudSecureInfo `json:"cloudsecures"`
	DefaultCloudName string                     `json:"default_cloud_name"`
}

// LoadOrCreateCloudSecureConfig loads or creates the configuration
func LoadOrCreateCloudSecureConfig(configFile string) (CloudSecureConfig, error) {
	var config CloudSecureConfig

	// Check if the specified config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("Configuration file not found, please provide the details for the first CloudSecure:")
		config.CloudSecures = make(map[string]CloudSecureInfo)
		cloudSecureInfo := CreateNewCloudSecureInfo()

		fmt.Print("CloudSecure Name: ")
		reader := bufio.NewReader(os.Stdin)
		cloudSecureName, _ := reader.ReadString('\n')
		cloudSecureName = strings.TrimSpace(cloudSecureName)

		config.CloudSecures[cloudSecureName] = cloudSecureInfo
		config.DefaultCloudName = cloudSecureName

		SaveCloudSecureConfig(configFile, config)
		fmt.Println("Configuration saved to", configFile)
	} else {
		configData, err := os.ReadFile(configFile)
		if err != nil {
			return config, err
		}
		json.Unmarshal(configData, &config)
	}

	// Check if any CloudSecure is missing or if default CloudSecure name is not set
	if len(config.CloudSecures) == 0 || config.DefaultCloudName == "" {
		fmt.Println("Invalid configuration. Adding a new CloudSecure:")
		cloudSecureInfo := CreateNewCloudSecureInfo()

		fmt.Print("CloudSecure Name: ")
		reader := bufio.NewReader(os.Stdin)
		cloudSecureName, _ := reader.ReadString('\n')
		cloudSecureName = strings.TrimSpace(cloudSecureName)

		if config.CloudSecures == nil {
			config.CloudSecures = make(map[string]CloudSecureInfo)
		}
		config.CloudSecures[cloudSecureName] = cloudSecureInfo
		if config.DefaultCloudName == "" {
			config.DefaultCloudName = cloudSecureName
		}

		SaveCloudSecureConfig(configFile, config)
		fmt.Println("Updated and saved configuration to", configFile)
	}

	return config, nil
}

// CreateNewCloudSecureInfo prompts the user for CloudSecure information
func CreateNewCloudSecureInfo() CloudSecureInfo {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("API Key: ")
	apiKey, _ := reader.ReadString('\n')
	fmt.Print("API Secret: ")
	apiSecret, _ := reader.ReadString('\n')
	fmt.Print("Tenant ID: ")
	tenantID, _ := reader.ReadString('\n')

	return CloudSecureInfo{
		APIKey:    strings.TrimSpace(apiKey),
		APISecret: strings.TrimSpace(apiSecret),
		TenantID:  strings.TrimSpace(tenantID),
	}
}

// SaveCloudSecureConfig saves the configuration to a JSON file
func SaveCloudSecureConfig(configFile string, config CloudSecureConfig) error {
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, configData, 0644)
}
