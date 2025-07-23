package models

// TokenBalance represents a token balance in a wallet
type TokenBalance struct {
	TokenAddress string  `json:"tokenAddress"`
	TokenName    string  `json:"tokenName"`
	TokenSymbol  string  `json:"tokenSymbol"`
	Balance      float64 `json:"balance"`
	USDBalance   float64 `json:"usdBalance"`
	Decimals     int     `json:"decimals"`
	LogoURI      string  `json:"logoURI"`
}

// WalletBalance represents the complete wallet balance information
type WalletBalance struct {
	WalletAddress string         `json:"walletAddress"`
	SOLBalance    float64        `json:"solBalance"`
	SOLUSDBalance float64        `json:"solUsdBalance"`
	Tokens        []TokenBalance `json:"tokens"`
	TotalUSDValue float64        `json:"totalUsdValue"`
}

// SolanaAPIResponse represents the response from Solana API
type SolanaAPIResponse struct {
	Result struct {
		Value []struct {
			Account struct {
				Data     []string `json:"data"`
				Owner    string   `json:"owner"`
				Lamports int64    `json:"lamports"`
			} `json:"account"`
			Pubkey string `json:"pubkey"`
		} `json:"value"`
	} `json:"result"`
}

// TokenInfo represents token metadata
type TokenInfo struct {
	Symbol      string  `json:"symbol"`
	Name        string  `json:"name"`
	Decimals    int     `json:"decimals"`
	LogoURI     string  `json:"logoURI"`
	CoingeckoID string  `json:"coingecko_id"`
	Price       float64 `json:"price"`
}

// BluechipConfig represents the configuration for bluechip tokens
type BluechipConfig struct {
	BluechipTokens map[string]TokenInfo `json:"bluechip_tokens"`
}
