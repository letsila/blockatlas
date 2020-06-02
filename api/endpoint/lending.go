package endpoint

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
)

// @Summary Get Lending Providers Info.
// @ID lending_providers
// @Description Get lending providers, their info, and supported assets.
// @Produce json
// @Tags Lending
// @Success 200 {object} blockatlas.DocsResponse Docs: []blockatlas.LendingProvider
// @Router /v1/lending/providers [get]
func ServeProviders(c *gin.Context, apis map[string]blockatlas.LendingAPI) {
	p, err := getProviders(apis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, blockatlas.DocsResponse{Docs: &p})
}

// @Summary Get Lending Rates.
// @ID lending_providers
// @Description Get lending rates, for one or more assets, of a provider.
// @Accept json
// @Produce json
// @Tags Lending
// @Param provider path string true "Lending provider name"
// @Param request body blockatlas.RatesRequest true "Request, containing one or more assets (token symbols)"
// @Success 200 {object} blockatlas.DocsResponse Docs: []blockatlas.LendingAssetRates
// @Router /v1/lending/rates/:provider [post]
func ServeRates(c *gin.Context, apis map[string]blockatlas.LendingAPI) {
	provider, ok := c.Params.Get("provider")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fatal: missing provider"})
		return
	}
	var req blockatlas.RatesRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	p, err := getRates(provider, req, apis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, blockatlas.DocsResponse{Docs: p})
}

// @Summary Get Account Contracts.
// @ID lending_providers
// @Description Get lending constracts, for one or more adresses, one or more assets (at one provider).
// @Accept json
// @Produce json
// @Tags Lending
// @Param provider path string true "Lending provider name"
// @Param request body blockatlas.RatesRequest true "Request, containing one or more assets (token symbols)"
// @Success 200 {object} blockatlas.DocsResponse Docs: []blockatlas.AccountLendingContracts
// @Router /v1/lending/account/:provider [post]
func ServeAccount(c *gin.Context, apis map[string]blockatlas.LendingAPI) {
	provider, ok := c.Params.Get("provider")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Fatal: missing provider"})
		return
	}
	var req blockatlas.AccountRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	p, err := getAccounts(provider, req, apis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, blockatlas.DocsResponse{Docs: &p})
}

// GetProviders return provider info list
func getProviders(apis map[string]blockatlas.LendingAPI) ([]blockatlas.LendingProvider, error) {
	ret := []blockatlas.LendingProvider{}
	for _, api := range apis {
		prov, err := api.GetProviderInfo()
		if err != nil {
			continue
		}
		ret = append(ret, prov)
	}
	return ret, nil
}

// GetRates return rates info
func getRates(provider string, req blockatlas.RatesRequest, apis map[string]blockatlas.LendingAPI) ([]blockatlas.LendingAssetRates, error) {
	api, ok := apis[provider]
	if !ok {
		return nil, fmt.Errorf("Unknown provider %v", provider)
	}
	return api.GetCurrentLendingRates(req.Assets)
}

// GetAccounts return account contract
func getAccounts(provider string, req blockatlas.AccountRequest, apis map[string]blockatlas.LendingAPI) ([]blockatlas.AccountLendingContracts, error) {
	api, ok := apis[provider]
	if !ok {
		return nil, fmt.Errorf("Unknown provider %v", provider)
	}
	return api.GetAccountLendingContracts(req)
}
