package templates

// Helper functions for styling and status management

// getStatusColor returns the CSS class for status indicator colors
func getStatusColor(status string) string {
	switch status {
	case "success":
		return "w-4 h-4 rounded-full bg-green-500"
	case "error":
		return "w-4 h-4 rounded-full bg-red-500"
	case "warning":
		return "w-4 h-4 rounded-full bg-yellow-500"
	default:
		return "w-4 h-4 rounded-full bg-gray-500"
	}
}

// getStatusTextColor returns the CSS class for status text colors
func getStatusTextColor(status string) string {
	switch status {
	case "success":
		return "text-green-700"
	case "error":
		return "text-red-700"
	case "warning":
		return "text-yellow-700"
	default:
		return "text-gray-700"
	}
}

// getStatusIcon returns the appropriate icon for a status
func getStatusIcon(status string) string {
	switch status {
	case "success":
		return "✓"
	case "error":
		return "✗"
	case "warning":
		return "⚠"
	default:
		return "?"
	}
}
