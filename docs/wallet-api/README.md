# Wallet API - Bruno Collection

This folder contains Bruno API collection for testing the Wallet API endpoints.

## Setup

1. Install Bruno from [usebruno.com](https://www.usebruno.com/)
2. Import this folder into Bruno
3. Update `Env-Local.bru` if needed:

```javascript
vars {
  base_url: http://127.0.0.1:8080
  auth_cookie: access_token
}
```

## Testing Guide

### 1. Authentication

**Register** → Create new account
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Login** → Get auth cookies
```json
{
  "username": "testuser",
  "password": "password123"
}
```

The Login request automatically saves cookies for use in other requests.

### 2. Account Operations

**Create Account** → Create a wallet
```json
{
  "account_name": "My Wallet",
  "currency": "IDR"
}
```

Save the returned `id` to use in subsequent requests.

### 3. Transaction Operations

**Deposit** → Add funds
```json
{
  "amount": 100000,
  "description": "Initial deposit"
}
```

**Withdraw** → Remove funds
```json
{
  "amount": 50000,
  "description": "ATM Withdrawal"
}
```

**Transfer** → Send to another account
```json
{
  "to_account_id": "destination-uuid",
  "amount": 25000,
  "description": "Transfer for food"
}
```

Note: You need to create 2 accounts to test transfers.

### 4. View History

**Get Transactions** → List all transactions with pagination
```
Query params: limit=10, offset=0
```

## Request Variables

Update these variables in Bruno for testing:

- `{{base_url}}` - API base URL (default: http://127.0.0.1:8080)
- `{{access_token}}` - Auto-populated from Login response
- `{{refresh_token}}` - Auto-populated from Login response
- `{{account_id}}` - Set manually from Create Account response
- `{{from_account_id}}` - Source account for transfers
- `{{to_account_id}}` - Destination account for transfers

## Race Condition Testing

To test the pessimistic locking (SELECT FOR UPDATE):

1. Use Bruno's "Runner" feature
2. Run multiple **Deposit** or **Transfer** requests simultaneously
3. All requests should complete without balance inconsistency

Example: Deposit 100000 three times simultaneously
- Expected final balance = initial + 300000 ✓
- Without locking: Could be initial + 100000 (race condition) ✗
