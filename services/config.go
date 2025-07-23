package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"solana-wallet-checker/models"
)

type ConfigService struct {
	bluechipTokens map[string]models.TokenInfo
}

func NewConfigService(configPath string) (*ConfigService, error) {
	config, err := loadBluechipConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &ConfigService{
		bluechipTokens: config.BluechipTokens,
	}, nil
}

func loadBluechipConfig(configPath string) (*models.BluechipConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config models.BluechipConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *ConfigService) IsBluechipToken(mintAddress string) bool {
	_, exists := c.bluechipTokens[mintAddress]
	return exists
}

func (c *ConfigService) GetTokenInfo(mintAddress string) (models.TokenInfo, bool) {
	tokenInfo, exists := c.bluechipTokens[mintAddress]
	return tokenInfo, exists
}

func (c *ConfigService) GetAllBluechipTokens() map[string]models.TokenInfo {
	return c.bluechipTokens
}
