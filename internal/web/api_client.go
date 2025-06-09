package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Order struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	DishID          string    `json:"dish_id"`
	Status          string    `json:"status"`
	DishName        string    `json:"dish_name"`
	DishDescription string    `json:"dish_description"`
	DishPrice       float64   `json:"dish_price"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type APIClientImpl struct {
	baseURL    string
	httpClient *http.Client
}

func NewAPIClient(baseURL string) *APIClientImpl {
	return &APIClientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *APIClientImpl) GetOrders() ([]Order, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/orders", c.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching orders: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var orders []Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return orders, nil
}

func (c *APIClientImpl) UpdateOrderStatus(orderID, status string) error {
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/orders/%s/status", c.baseURL, orderID), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	req.Header.Set("Content-Type", "application/json")

	body := map[string]string{"status": status}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling body: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error updating order status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// log response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		log.Printf("Response body: %s", string(body))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
