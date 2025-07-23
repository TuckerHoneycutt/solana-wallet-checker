package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"solana-wallet-checker/models"
	"strconv"
	"strings"
	"time"
)

const (
	SOLANA_RPC_URL = "https://api.mainnet-beta.solana.com"
	COINGECKO_API  = "https://api.coingecko.com/api/v3"
)

type SolanaService struct {
	httpClient    *http.Client
	configService *ConfigService
}

func NewSolanaService(configService *ConfigService) *SolanaService {
	return &SolanaService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		configService: configService,
	}
}

// GetWalletBalance fetches the complete wallet balance information
func (s *SolanaService) GetWalletBalance(walletAddress string) (*models.WalletBalance, error) {
	// Validate wallet address format (basic validation)
	if len(walletAddress) < 32 || len(walletAddress) > 44 {
		return nil, fmt.Errorf("invalid wallet address format")
	}

	walletBalance := &models.WalletBalance{
		WalletAddress: walletAddress,
		Tokens:        []models.TokenBalance{},
	}

	// Get SOL balance
	solBalance, err := s.getSOLBalance(walletAddress)
	if err != nil {
		fmt.Printf("Warning: failed to get SOL balance: %v\n", err)
		solBalance = 0
	}
	walletBalance.SOLBalance = solBalance

	// Get SOL price in USD
	solPrice, err := s.getSOLPrice()
	if err != nil {
		fmt.Printf("Warning: failed to get SOL price: %v\n", err)
		solPrice = 0
	} else {
		fmt.Printf("SOL price: $%.4f\n", solPrice)
	}
	walletBalance.SOLUSDBalance = solBalance * solPrice
	fmt.Printf("SOL balance: %.6f SOL = $%.2f\n", solBalance, walletBalance.SOLUSDBalance)

	// Get token accounts (only bluechip tokens)
	tokens, err := s.getTokenBalances(walletAddress)
	if err != nil {
		fmt.Printf("Warning: failed to get token balances: %v\n", err)
		tokens = []models.TokenBalance{}
	}
	walletBalance.Tokens = tokens

	// Calculate total USD value
	totalUSD := walletBalance.SOLUSDBalance
	for _, token := range walletBalance.Tokens {
		totalUSD += token.USDBalance
	}
	walletBalance.TotalUSDValue = totalUSD

	fmt.Printf("Total portfolio value: $%.2f\n", totalUSD)

	return walletBalance, nil
}

// getSOLBalance fetches the SOL balance for a wallet
func (s *SolanaService) getSOLBalance(walletAddress string) (float64, error) {
	payload := fmt.Sprintf(`{
        "jsonrpc": "2.0",
        "id": 1,
        "method": "getBalance",
        "params": ["%s"]
    }`, walletAddress)

	resp, err := s.httpClient.Post(SOLANA_RPC_URL, "application/json", strings.NewReader(payload))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Result struct {
			Value int64 `json:"value"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// Convert lamports to SOL (1 SOL = 1e9 lamports)
	return float64(result.Result.Value) / 1e9, nil
}

// getSOLPrice fetches the current SOL price in USD
func (s *SolanaService) getSOLPrice() (float64, error) {
	resp, err := s.httpClient.Get(COINGECKO_API + "/simple/price?ids=solana&vs_currencies=usd")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Solana struct {
			USD float64 `json:"usd"`
		} `json:"solana"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Solana.USD, nil
}

// getTokenBalances fetches bluechip token balances for a wallet
func (s *SolanaService) getTokenBalances(walletAddress string) ([]models.TokenBalance, error) {
	payload := fmt.Sprintf(`{
        "jsonrpc": "2.0",
        "id": 1,
        "method": "getTokenAccountsByOwner",
        "params": [
            "%s",
            {
                "programId": "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
            },
            {
                "encoding": "jsonParsed"
            }
        ]
    }`, walletAddress)

	resp, err := s.httpClient.Post(SOLANA_RPC_URL, "application/json", strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Result struct {
			Value []struct {
				Account struct {
					Data struct {
						Parsed struct {
							Info struct {
								TokenAmount struct {
									Amount         string `json:"amount"`
									Decimals       int    `json:"decimals"`
									UIAmountString string `json:"uiAmountString"`
								} `json:"tokenAmount"`
								Mint string `json:"mint"`
							} `json:"info"`
						} `json:"parsed"`
					} `json:"data"`
				} `json:"account"`
			} `json:"value"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var tokens []models.TokenBalance
	var coingeckoIDs []string

	// First pass: collect all bluechip tokens and their coingecko IDs
	for _, account := range result.Result.Value {
		tokenInfo := account.Account.Data.Parsed.Info

		// Skip tokens with zero balance
		if tokenInfo.TokenAmount.UIAmountString == "0" || tokenInfo.TokenAmount.UIAmountString == "" {
			continue
		}

		// Only process bluechip tokens
		if !s.configService.IsBluechipToken(tokenInfo.Mint) {
			continue
		}

		balance, err := strconv.ParseFloat(tokenInfo.TokenAmount.UIAmountString, 64)
		if err != nil || balance <= 0 {
			continue
		}

		// Get token info from config
		configTokenInfo, exists := s.configService.GetTokenInfo(tokenInfo.Mint)
		if !exists {
			continue
		}

		tokenBalance := models.TokenBalance{
			TokenAddress: tokenInfo.Mint,
			TokenName:    configTokenInfo.Name,
			TokenSymbol:  configTokenInfo.Symbol,
			Balance:      balance,
			USDBalance:   0, // Will be set after getting price
			Decimals:     configTokenInfo.Decimals,
			LogoURI:      configTokenInfo.LogoURI,
		}

		tokens = append(tokens, tokenBalance)
		coingeckoIDs = append(coingeckoIDs, configTokenInfo.CoingeckoID)
	}

	// Get prices for all bluechip tokens found
	if len(coingeckoIDs) > 0 {
		fmt.Printf("Fetching prices for: %v\n", coingeckoIDs)
		prices, err := s.getBluechipPrices(coingeckoIDs)
		if err != nil {
			fmt.Printf("Error fetching prices: %v\n", err)
		} else {
			fmt.Printf("Received prices: %v\n", prices)
			// Update USD balances
			for i := range tokens {
				configTokenInfo, _ := s.configService.GetTokenInfo(tokens[i].TokenAddress)
				if price, exists := prices[configTokenInfo.CoingeckoID]; exists && price > 0 {
					tokens[i].USDBalance = tokens[i].Balance * price
					fmt.Printf("Token %s: balance=%.6f, price=%.4f, usd=%.2f\n",
						tokens[i].TokenSymbol, tokens[i].Balance, price, tokens[i].USDBalance)
				} else {
					fmt.Printf("No price found for %s (coingecko_id: %s)\n",
						tokens[i].TokenSymbol, configTokenInfo.CoingeckoID)
				}
			}
		}
	}

	return tokens, nil
}

// getBluechipPrices fetches prices for multiple tokens from CoinGecko
func (s *SolanaService) getBluechipPrices(coingeckoIDs []string) (map[string]float64, error) {
	// Remove duplicates from the slice
	uniqueIDs := make(map[string]bool)
	var cleanIDs []string
	for _, id := range coingeckoIDs {
		if !uniqueIDs[id] {
			uniqueIDs[id] = true
			cleanIDs = append(cleanIDs, id)
		}
	}

	ids := strings.Join(cleanIDs, ",")
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", COINGECKO_API, ids)

	fmt.Printf("Fetching prices from: %s\n", url)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch prices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("CoinGecko API returned status %d", resp.StatusCode)
	}

	var result map[string]struct {
		USD float64 `json:"usd"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode price response: %w", err)
	}

	prices := make(map[string]float64)
	for id, priceInfo := range result {
		prices[id] = priceInfo.USD
		fmt.Printf("Price for %s: $%.4f\n", id, priceInfo.USD)
	}

	return prices, nil
}
