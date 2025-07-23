# Solana Wallet Checker

A web application built with Go and Echo framework that allows users to check the balance of Solana wallets, specifically focusing on bluechip tokens and SOL. The application provides real-time balance information with USD valuations using CoinGecko API.

## Features

- **Wallet Balance Checking**: Enter any Solana wallet address to view its balance
- **SOL Balance**: Shows native SOL balance with real-time USD conversion
- **Bluechip Token Support**: Tracks major tokens including:
  - SOL (Solana)
  - ETH (Ethereum via Wormhole)
  - USDC (USD Coin)
  - USDT (Tether)
  - WBTC (Wrapped Bitcoin via Wormhole)
  - cbBTC (Coinbase Wrapped BTC)
- **Real-time Pricing**: Uses CoinGecko API for current token prices
- **Modern UI**: Clean, responsive interface built with Tailwind CSS and Templ
- **RESTful API**: Built with Echo framework for fast, lightweight performance

## Screenshots

The application features a clean, modern interface with:
- A simple form to enter wallet addresses
- Real-time balance display with USD valuations
- Token logos and detailed balance information
- Responsive design that works on desktop and mobile

## Prerequisites

- Go 1.24.4 or higher
- Internet connection (for Solana RPC and CoinGecko API access)

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/solana-wallet-checker.git
   cd solana-wallet-checker
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## Usage

1. Open your web browser and navigate to `http://localhost:8080`
2. Enter a valid Solana wallet address in the input field
3. Click "Check Balance" to view the wallet's holdings
4. The application will display:
   - SOL balance with USD value
   - Bluechip token balances with USD values
   - Total portfolio value in USD

### Example Wallet Address
```
9WzDXwBbmkg8ZTbNMqUxvQRAyrZzDsGYdLVL9zYtAWWM
```

## API Endpoints

- `GET /` - Home page with wallet input form
- `GET /balance?wallet=<address>` - Display balance for specified wallet address

## Configuration

### Bluechip Tokens
The application tracks specific bluechip tokens defined in `config/bluechip_tokens.json`. Each token includes:
- Token mint address
- Symbol and name
- Decimal places
- Logo URI
- CoinGecko ID for price fetching

### Adding New Tokens
To add support for additional tokens:

1. Edit `config/bluechip_tokens.json`
2. Add a new entry with the token's mint address
3. Include the required fields (symbol, name, decimals, logoURI, coingecko_id)
4. Restart the application

## Architecture

```
solana-wallet-checker/
├── main.go                 # Application entry point
├── handlers/               # HTTP request handlers
│   └── handlers.go
├── services/               # Business logic services
│   ├── config.go          # Configuration management
│   └── solana.go          # Solana blockchain integration
├── models/                 # Data models
│   └── token.go
├── templates/              # UI templates (Templ)
│   ├── layout.templ       # Base layout
│   ├── home.templ         # Home page
│   └── balance.templ      # Balance display page
└── config/                 # Configuration files
    └── bluechip_tokens.json
```

## Dependencies

- **Echo v4**: High-performance HTTP framework
- **Templ**: Type-safe HTML templating
- **Go 1.24.4**: Programming language runtime

## External APIs

- **Solana RPC**: `https://api.mainnet-beta.solana.com` - For blockchain data
- **CoinGecko API**: `https://api.coingecko.com/api/v3` - For token prices

## Development

### Running in Development Mode
```bash
go run main.go
```

### Building for Production
```bash
go build -o solana-wallet-checker main.go
./solana-wallet-checker
```

### Project Structure
- **Handlers**: Handle HTTP requests and responses
- **Services**: Contain business logic for Solana integration and configuration
- **Models**: Define data structures for wallet balances and tokens
- **Templates**: UI components using Templ framework

## Error Handling

The application includes robust error handling for:
- Invalid wallet addresses
- Network connectivity issues
- API rate limiting
- Malformed responses

## Security Considerations

- Input validation for wallet addresses
- CORS middleware enabled
- No sensitive data storage
- Read-only blockchain access

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This application is for educational and informational purposes only. Always verify wallet addresses and double-check balance information. The application relies on external APIs and may not always reflect the most current data.

## Support

If you encounter any issues or have questions:
1. Check the existing issues in the repository
2. Create a new issue with detailed information about your problem
3. Include your Go version and operating system

---

Built with ❤️ using Go and the Solana blockchain
