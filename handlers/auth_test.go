package handlers

import (
	"log"
	"net/http"
)

// =============================================================================
// TEST HANDLERS
// =============================================================================
// These handlers are used for testing authentication flows and debugging.
// They should be disabled or removed in production environments.
// =============================================================================

// TestTokenRefreshHandler serves a test page with token refresh button
func (h *AuthHandler) TestTokenRefreshHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	testHTML := `
<!DOCTYPE html>
<html>
<head>
    <title>Auth Test - Token Refresh</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        // Helper function to check current user status
        function checkUserStatus() {
            console.log('👤 USER STATUS CHECK: === STARTED ===');
            const resultDiv = document.getElementById('user-result');
            
            fetch('/api/auth/user')
            .then(response => response.json())
            .then(data => {
                console.log('👤 USER STATUS CHECK: Response:', data);
                
                if (data.logged_in) {
                    resultDiv.innerHTML = '<p class="text-green-600">✅ Logged in as: ' + data.name + '</p>' +
                        '<p class="text-sm text-gray-600">Email: ' + data.email + '</p>';
                } else {
                    resultDiv.innerHTML = '<p class="text-red-600">❌ Not logged in</p>';
                }
            })
            .catch(error => {
                console.log('👤 USER STATUS CHECK: Error:', error);
                resultDiv.innerHTML = '<p class="text-red-600">❌ Error: ' + error.message + '</p>';
            });
        }
        
        // Log when page loads
        document.addEventListener('DOMContentLoaded', function() {
            console.log('🧪 AUTH TEST PAGE: Loaded - Check browser console for detailed testing logs');
        });
    </script>
</head>
<body class="bg-blue-300 min-h-screen">
    <div class="container mx-auto py-8 px-4">
        <h1 class="text-3xl font-bold text-center mb-8">Authentication Test Page</h1>
        
        <div class="max-w-2xl mx-auto space-y-6">
            <!-- Check Current User -->
            <div class="bg-white rounded-lg shadow p-6">
                <h2 class="text-xl font-semibold mb-4">Check Current User</h2>
                <p class="text-gray-600 mb-4">Check if user is currently logged in.</p>
                <button
                    onclick="checkUserStatus()"
                    class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                >
                    Check User Status
                </button>
                <div id="user-result" class="mt-4 p-3 bg-gray-100 rounded"></div>
            </div>
            
            <!-- OAuth Login Buttons -->
            <div class="bg-white rounded-lg shadow p-6">
                <h2 class="text-xl font-semibold mb-4">OAuth Login</h2>
                <p class="text-gray-600 mb-4">Use these to test the full OAuth flow.</p>
                <div class="space-x-4">
                    <a href="/auth/google" class="bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white font-bold py-3 px-6 rounded-lg transition-all duration-300 shadow-lg hover:shadow-red-500/25 border border-red-400/30">Login with Google</a>
                    <a href="/auth/github" class="bg-gradient-to-r from-gray-600 to-gray-700 hover:from-gray-700 hover:to-gray-800 text-white font-bold py-3 px-6 rounded-lg transition-all duration-300 shadow-lg hover:shadow-gray-500/25 border border-gray-400/30">Login with GitHub</a>
                    <a href="/auth/discord" class="bg-gradient-to-r from-indigo-500 to-indigo-600 hover:from-indigo-600 hover:to-indigo-700 text-white font-bold py-3 px-6 rounded-lg transition-all duration-300 shadow-lg hover:shadow-indigo-500/25 border border-indigo-400/30">Login with Discord</a>
                    <a href="/auth/microsoft" class="bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white font-bold py-3 px-6 rounded-lg transition-all duration-300 shadow-lg hover:shadow-blue-500/25 border border-blue-400/30">Login with Microsoft</a>
                </div>
            </div>
            
            <!-- Instructions -->
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                <h3 class="text-lg font-semibold text-yellow-800 mb-2">🔍 Testing Instructions</h3>
                <ol class="text-sm text-sm text-yellow-700 space-y-1">
                    <li>1. Open browser console (F12) to see detailed logs</li>
                    <li>2. Login with Google, GitHub, Discord, or Microsoft</li>
                    <li>3. Check console logs for every step of the process</li>
                </ol>
            </div>
        </div>
    </div>
</body>
</html>
`
	if _, err := w.Write([]byte(testHTML)); err != nil {
		log.Printf("Error writing test HTML: %v", err)
	}
}