package handlers

import (
	"fmt"
	"net/http"
	"solana-wallet-checker/models"
	"solana-wallet-checker/services"
	"solana-wallet-checker/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	solanaService *services.SolanaService
}

func NewHandlers(configService *services.ConfigService) *Handlers {
	return &Handlers{
		solanaService: services.NewSolanaService(configService),
	}
}

// HomeHandler renders the home page
func (h *Handlers) HomeHandler(c echo.Context) error {
	component := templates.Home()
	return templ.Handler(component).Component.Render(c.Request().Context(), c.Response().Writer)
}

// BalanceHandler handles the balance page
func (h *Handlers) BalanceHandler(c echo.Context) error {
	walletAddress := c.QueryParam("wallet")

	if walletAddress == "" {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	fmt.Printf("Checking wallet: %s\n", walletAddress)

	// Get wallet balance from Solana service
	walletBalance, err := h.solanaService.GetWalletBalance(walletAddress)
	if err != nil {
		fmt.Printf("Error getting wallet balance: %v\n", err)
		// In a real app, you'd want better error handling
		// For now, create a mock response with error info
		walletBalance = &models.WalletBalance{
			WalletAddress: walletAddress,
			SOLBalance:    0,
			SOLUSDBalance: 0,
			Tokens:        []models.TokenBalance{},
			TotalUSDValue: 0,
		}
	}

	fmt.Printf("Final wallet balance - SOL: %.6f ($%.2f), Tokens: %d, Total: $%.2f\n",
		walletBalance.SOLBalance, walletBalance.SOLUSDBalance,
		len(walletBalance.Tokens), walletBalance.TotalUSDValue)

	component := templates.Balance(walletBalance)
	return templ.Handler(component).Component.Render(c.Request().Context(), c.Response().Writer)
}
