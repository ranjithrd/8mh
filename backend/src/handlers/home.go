package handlers

import (
	"backend/src/db"
	"backend/src/repos"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	statsRepo = repos.StatsRepo{}
	loanRepo  = repos.LoanRepo{}
)

var (
	blockchainVerificationCache struct {
		sync.RWMutex
		isValid       bool
		lastChecked   time.Time
		cacheDuration time.Duration
	}
)

func init() {
	blockchainVerificationCache.cacheDuration = 10 * time.Second
	blockchainVerificationCache.isValid = true
	blockchainVerificationCache.lastChecked = time.Time{}
}

func getBlockchainIntegrity() bool {
	blockchainVerificationCache.RLock()
	if time.Since(blockchainVerificationCache.lastChecked) < blockchainVerificationCache.cacheDuration {
		cached := blockchainVerificationCache.isValid
		blockchainVerificationCache.RUnlock()
		return cached
	}
	blockchainVerificationCache.RUnlock()

	blockchainVerificationCache.Lock()
	defer blockchainVerificationCache.Unlock()

	if time.Since(blockchainVerificationCache.lastChecked) < blockchainVerificationCache.cacheDuration {
		return blockchainVerificationCache.isValid
	}

	valid, _ := db.VerifyEntireChain()
	blockchainVerificationCache.isValid = valid
	blockchainVerificationCache.lastChecked = time.Now()

	return valid
}

type HomeResponse struct {
	TotalAssets         int64  `json:"total_assets" example:"1000000"`
	TotalLoans          int64  `json:"total_loans" example:"500000"`
	TotalProfit         int64  `json:"total_profit" example:"50000"`
	DividendExpected    *int64 `json:"dividend_expected,omitempty" example:"5000"`
	Role                string `json:"role" example:"member"`
	BlockchainIntegrity bool   `json:"blockchain_integrity" example:"true"`
}

// Home godoc
// @Summary Get home dashboard statistics
// @Description Returns financial statistics based on user role (member, manager, auditor)
// @Tags home
// @Produce json
// @Security SessionAuth
// @Success 200 {object} HomeResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/home [get]
func Home(c echo.Context) error {
	user := c.Get("user").(*repos.UserWithSession)

	totalAssets, err := statsRepo.GetTotalAssets()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get total assets"})
	}

	totalLoans, err := loanRepo.GetTotalLoansAmount()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get total loans"})
	}

	totalProfit, err := loanRepo.GetTotalProfit()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get total profit"})
	}

	response := HomeResponse{
		TotalAssets:         totalAssets,
		TotalLoans:          totalLoans,
		TotalProfit:         totalProfit,
		Role:                user.Role,
		BlockchainIntegrity: getBlockchainIntegrity(),
	}

	if user.Role == "member" {
		totalShares, err := statsRepo.GetTotalSharesBalance()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to calculate dividend"})
		}

		if totalShares > 0 {
			dividend := (totalProfit * int64(user.SharesBalance)) / totalShares
			response.DividendExpected = &dividend
		} else {
			zero := int64(0)
			response.DividendExpected = &zero
		}
	}

	return c.JSON(http.StatusOK, response)
}
