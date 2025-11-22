package dashboard

import (
	"net/http"

	"github.com/DraconDev/go-templ-htmx-ex/internal/clients/paymentms"
	"github.com/DraconDev/go-templ-htmx-ex/internal/handlers/auth/session"
	"github.com/DraconDev/go-templ-htmx-ex/internal/utils/config"
	"github.com/DraconDev/go-templ-htmx-ex/templates/pages"
)

type DashboardHandler struct {
	config        *config.Config
	paymentClient *paymentms.Client
}

func NewDashboardHandler(cfg *config.Config, paymentClient *paymentms.Client) *DashboardHandler {
	return &DashboardHandler{
		config:        cfg,
		paymentClient: paymentClient,
	}
}

func (h *DashboardHandler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get user info from session
	user, err := session.GetUserFromSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get subscription status
	// Note: In a real app, you'd get the user ID from the session.
	// For now, we'll use the email or a placeholder if ID is missing.
	userID := user.Email // Fallback since we might not have ID in session yet

	// Fetch subscription status
	// We use the product ID from config
	subStatus, err := h.paymentClient.GetSubscriptionStatus(r.Context(), userID, h.config.StripeProductID)

	// Prepare view model
	isPro := false
	status := "Free Plan"
	periodEnd := ""

	if err == nil && subStatus != nil {
		if subStatus.Status == "active" {
			isPro = true
			status = "Pro Plan"
			periodEnd = subStatus.CurrentPeriodEnd.Format("Jan 02, 2006")
		}
	}

	// Render template
	component := pages.Dashboard(user.Name, user.Email, user.Picture, status, isPro, periodEnd)
	component.Render(r.Context(), w)
}
