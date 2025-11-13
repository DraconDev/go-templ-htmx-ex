# Token Refresh Testing Plan

## 🎯 Testing Objective
Verify that the token refresh mechanism works correctly when access tokens expire.

## 🔧 Test Implementation Steps

### 1. **Token Expiration Simulation**
Since access tokens expire after 1 hour, we need a way to test this:
- Option A: Modify JWT generation to create short-lived tokens for testing
- Option B: Manually expire the session_token cookie
- Option C: Create a test endpoint to force token expiration

### 2. **Refresh Flow Verification**
Test the complete refresh cycle:
1. User is logged in with valid token
2. Token expires 
3. User tries to access protected page
4. Frontend detects expired token
5. Frontend calls `/api/auth/refresh`
6. Server uses refresh_token cookie to get new tokens
7. Server sets new session_token cookie
8. User continues browsing seamlessly

### 3. **Test Checklist**
- [ ] **Refresh endpoint exists and responds correctly**
- [ ] **Refresh token cookie is sent with request** (automatic browser behavior)
- [ ] **Auth service returns new tokens** when given valid refresh_token
- [ ] **New session_token is set in response**
- [ ] **Old session_token is replaced** (not duplicated)
- [ ] **User experience is seamless** (no logout/redirect needed)

## 🧪 Test Implementation

Let me implement a test endpoint to force token expiration and test the refresh flow.