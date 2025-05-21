package snaptradeclient

import (
	"fmt"

	snaptrade "github.com/passiv/snaptrade-sdks/sdks/go"
)

type SnapTradeClient struct {
	api        *snaptrade.APIClient
	userID     string
	userSecret string
}

func New(clientId, clientSecret, userId, userSecret string) *SnapTradeClient {
	configuration := snaptrade.NewConfiguration()
	configuration.SetPartnerClientId(clientId)
	configuration.SetConsumerKey(clientSecret)
	client := snaptrade.NewAPIClient(configuration)

	if client == nil {
		panic("Failed to create API client")
	}
	return &SnapTradeClient{
		api:        client,
		userID:     userId,
		userSecret: userSecret,
	}
}

func (c *SnapTradeClient) RegisterUser() (string, error) {
	response, _, err := c.api.AuthenticationApi.RegisterSnapTradeUser(snaptrade.SnapTradeRegisterUserRequestBody{
		UserId: c.userID,
	}).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}
	if response.UserSecret == nil {
		return "", fmt.Errorf("user secret is nil")
	}
	return *response.UserSecret, nil
}

func (c *SnapTradeClient) DeleteUser() error {
	_, _, err := c.api.AuthenticationApi.DeleteSnapTradeUser(c.userID).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (c *SnapTradeClient) LoginUserAndGetRedirectURI(brokerageSlug string) (string, error) {
	response, _, err := c.api.AuthenticationApi.LoginSnapTradeUser(c.userID, c.userSecret).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to login to SnapTrade: %w", err)
	}
	redirectURI := response.LoginRedirectURI.GetRedirectURI() + "&connectionType=trade&broker=" + brokerageSlug
	return redirectURI, nil
}

func (c *SnapTradeClient) ListUserAccounts() ([]snaptrade.Account, error) {
	response, _, err := c.api.AccountInformationApi.ListUserAccounts(c.userID, c.userSecret).Execute()
	if err != nil {
		return nil, fmt.Errorf("error retrieving accounts: %w", err)
	}
	return response, nil
}

func (c *SnapTradeClient) GetUserAccountPositions(accountID string) ([]snaptrade.Position, error) {
	positions, _, err := c.api.AccountInformationApi.GetUserAccountPositions(c.userID, c.userSecret, accountID).Execute()
	if err != nil {
		return nil, fmt.Errorf("error retrieving positions for account %s: %w", accountID, err)
	}
	return positions, nil
}

func (c *SnapTradeClient) GetUserAccountOrders(accountID string) ([]snaptrade.AccountOrderRecord, error) {
	orders, _, err := c.api.AccountInformationApi.GetUserAccountOrders(c.userID, c.userSecret, accountID).Execute()
	if err != nil {
		return nil, fmt.Errorf("error retrieving orders for account %s: %w", accountID, err)
	}
	return orders, nil
}

func (c *SnapTradeClient) PlaceForceOrder(accountID, action, ticker string, quantity float32) (*snaptrade.AccountOrderRecord, error) {
	orderRecord, _, err := c.api.TradingApi.PlaceForceOrder(c.userID, c.userSecret, snaptrade.ManualTradeFormWithOptions{
		AccountId:   accountID,
		Action:      snaptrade.ActionStrictWithOptions(action),
		Symbol:      *snaptrade.NewNullableString(&ticker),
		OrderType:   snaptrade.OrderTypeStrict("Market"),
		TimeInForce: snaptrade.TimeInForceStrict("GTC"),
		Units:       *snaptrade.NewNullableFloat32(&quantity),
	}).Execute()
	if err != nil {
		return nil, fmt.Errorf("error placing order: %w", err)
	}
	return orderRecord, nil
}
