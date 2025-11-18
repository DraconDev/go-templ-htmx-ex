# Authentication Data Format - Working Reference vs Current

## 🔍 **Key Format Differences Found**

### **1. Auth Service Request Format**

**✅ WORKING (Reference):**
```go
// Reference uses "code" parameter
jsonData := map[string]string{"auth_code": code}
```

**❌ CURRENT (Your Implementation):**
```go
// Current uses "auth_code" parameter 
jsonData := map[string]string{"auth_code": auth_code}
```

### **2. Response Structure Expectations**

**✅ WORKING (Reference - Flexible Parsing):**
```go
// Reference first parses as generic map
var respData map[string]interface{}
json.Unmarshal(bodyBytes, &respData)

// Then extracts session_id specifically
if sessionInterface, exists := respData["session_id"]; exists {
    if sessionStr, ok := sessionInterface.(string); ok {
        sessionID = sessionStr
        hasSessionID = true
    }
}
```

**❌ CURRENT (Rigid Parsing):**
```go
// Current expects specific AuthResponse struct
var authResp models.AuthResponse
json.Unmarshal(bodyBytes, &authResp)
return &authResp, nil
```

### **3. Expected Response Field**

**✅ WORKING (Reference):**
```json
{
  "session_id": "actual-session-id-string"
}
```

**❌ CURRENT (Expects):**
```json
{
  "success": true,
  "user_id": "user-id", 
  "email": "user@example.com",
  "name": "User Name"
}
```

### **4. Session Token Format**

**✅ WORKING (Reference):**
```go
// Returns session_id as IdToken
return &models.TokenExchangeResponse{
    Success: true,
    IdToken: sessionID, // session_id from response
}, nil
```

**❌ CURRENT (Uses UserID):**
```go
// Uses UserID as session identifier
Value: authResp.UserID,
```

### **5. Session Validation**

**✅ WORKING (Reference):**
```go
// Uses session_id for validation
ValidateSession(sessionID string) (*models.AuthResponse, error)
{
    return s.CallAuthService(..., map[string]string{
        "session_id": sessionID,
    })
}
```

**❌ CURRENT (UserID-based):**
```go
// Same approach but expects different response format
ValidateSession(session_id string) (*models.AuthResponse, error)
{
    return s.callAuthService("/auth/session/refresh", map[string]string{
        "session_id": session_id,
    })
}
```

## 🔧 **Recommended Format Alignment**

To match your working reference, your system should expect:

### **Auth Service Request:**
```json
{
  "auth_code": "github_12345_cb67890"
}
```

### **Auth Service Response:**
```json
{
  "session_id": "actual-session-id-here"
}
```

### **API Response to Frontend:**
```json
{
  "success": true,
  "id_token": "actual-session-id-here"
}
```

The issue is your current implementation expects `AuthResponse` format but the auth service is returning `session_id` format like the working reference.