package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DraconDev/go-templ-htmx-ex/internal/middleware"
	"github.com/DraconDev/go-templ-htmx-ex/internal/utils/config"
	"github.com/DraconDev/go-templ-htmx-ex/templates/layouts"
	"github.com/DraconDev/go-templ-htmx-ex/templates/pages"
	"github.com/a-h/templ"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	Config *config.Config
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(config *config.Config) *PaymentHandler {
	return &PaymentHandler{
		Config: config,
	}
}

// PaymentPageHandler handles the payment page display
func (h *PaymentHandler) PaymentPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	// Get user info from middleware context
	userInfo := middleware.GetUserFromContext(r)
	if !userInfo.LoggedIn {
		// Redirect to home if not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Create payment page content with user data
	navigation := layouts.NavigationLoggedIn(userInfo)
	component := layouts.Layout("Payment", "Subscribe to access premium features and content.", navigation, pages.PaymentContent(userInfo))
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render payment page", http.StatusInternalServerError)
		return
	}
}

// CheckoutHandler handles payment checkout initiation
func (h *PaymentHandler) CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get user info from middleware context
	userInfo := middleware.GetUserFromContext(r)
	if !userInfo.LoggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Authentication required",
		})
		return
	}

	// Parse request body
	var req struct {
		PriceID string `json:"price_id"`
		ProductID string `json:"product_id"`
		SuccessURL string `json:"success_url"`
		CancelURL string `json:"cancel_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	// Validate required fields
	if req.PriceID == "" || req.ProductID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Missing required fields: price_id, product_id",
		})
		return
	}

	// Set default URLs if not provided
	baseURL := "http://localhost:8081"
	if req.SuccessURL == "" {
		req.SuccessURL = baseURL + "/payment/success"
	}
	if req.CancelURL == "" {
		req.CancelURL = baseURL + "/payment/cancel"
	}

	// Call payment microservice
	checkoutResp, err := h.createCheckoutSession(userInfo, req.PriceID, req.ProductID, req.SuccessURL, req.CancelURL)
	if err != nil {
		fmt.Printf("❌ PAYMENT: Failed to create checkout session: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to create checkout session",
		})
		return
	}

	// Return checkout URL to frontend
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"checkout_url": checkoutResp.CheckoutURL,
		"checkout_session_id": checkoutResp.CheckoutSessionID,
	})
}

// createCheckoutSession calls the payment microservice to create a Stripe checkout session
func (h *PaymentHandler) createCheckoutSession(userInfo layouts.UserInfo, priceID, productID, successURL, cancelURL string) (*CheckoutResponse, error) {
	// Create payment microservice request
	paymentReq := map[string]interface{}{
		"user_id": userInfo.Email, // Using email as user ID for now
		"email": userInfo.Email,
		"product_id": productID,
		"price_id": priceID,
		"success_url": successURL,
		"cancel_url": cancelURL,
	}

	// Make HTTP request to payment microservice
	client := &http.Client{Timeout: 10 * time.Second}
	jsonData, err := json.Marshal(paymentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Get API key from environment or config
	apiKey := h.Config.GetEnv("PAYMENT_MS_API_KEY", "")
	if apiKey == "" {
		return nil, fmt.Errorf("payment microservice API key not configured")
	}

	req, err := http.NewRequest("POST", "http://localhost:9000/api/v1/checkout/subscription", 
		http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment microservice: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var checkoutResp CheckoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&checkoutResp); err != nil {
		return nil, fmt.Errorf("failed to parse checkout response: %w", err)
	}

	return &checkoutResp, nil
}

// CheckoutResponse represents the response from payment microservice
type CheckoutResponse struct {
	CheckoutSessionID string `json:"checkout_session_id"`
	CheckoutURL string `json:"checkout_url"`
}

// SuccessHandler handles successful payment redirects
func (h *PaymentHandler) SuccessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	userInfo := middleware.GetUserFromContext(r)
	navigation := layouts.NavigationLoggedIn(userInfo)
	component := layouts.Layout("Payment Success", "Thank you for your purchase!", navigation, pages.PaymentSuccessContent())
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render success page", http.StatusInternalServerError)
		return
	}
}

// CancelHandler handles cancelled payment redirects
func (h *PaymentHandler) CancelHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	
	userInfo := middleware.GetUserFromContext(r)
	navigation := layouts.NavigationLoggedIn(userInfo)
	component := layouts.Layout("Payment Cancelled", "Payment was cancelled. You can try again.", navigation, pages.PaymentCancelContent())
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render cancel page", http.StatusInternalServerError)
		return
	}
}