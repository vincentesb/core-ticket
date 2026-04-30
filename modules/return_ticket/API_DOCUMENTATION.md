# Return Ticket Module API Documentation

## Overview

The Return Ticket module provides endpoints to retrieve and analyze return tickets from the ticketing system. It allows filtering by product type and issue classification, with advanced analysis capabilities to identify patterns and common issues.

## Base URL

```
http://localhost:3005/return-ticket
```

## Authentication

All endpoints in this module support JWT authentication via the `Authorization` header:

```
Authorization: Bearer <JWT_TOKEN>
```

---

## Endpoints

### 1. Get Return Tickets by Issue Type

Retrieves return tickets filtered by product type and optionally by issue type.

**Endpoint:**
```
GET /return-ticket
```

**Method:** `GET`

**Request Body:**

```json
{
  "productTypeID": 1,
  "issueType": "Stock Opname"
}
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `productTypeID` | integer | Yes | Product type identifier (must be >= 1) |
| `issueType` | string | No | Filter by specific issue type (ASCII characters only) |

**Response:**

```json
[
  {
    "ticketNum": "TKT-001",
    "companyCode": "COMP001",
    "outletCode": "OUT001",
    "outletName": "Main Outlet",
    "issue": "Issue Description",
    "issueId": 1,
    "description": "Detailed description of the issue",
    "assignDeveloper": "developer@example.com",
    "checkNotes": "Initial check notes",
    "devCheckNotes": "Developer check notes",
    "solution": "Applied solution",
    "statusId": 1,
    "createdBy": "creator@example.com",
    "createdDate": "2024-01-15T10:30:00Z",
    "updatedBy": "updater@example.com",
    "updatedDate": "2024-01-16T14:45:00Z",
    "date": "2024-01-15T00:00:00Z",
    "ref": "REF123",
    "refFrom": "External System",
    "userId": "USER123",
    "productTypeId": 1
  }
]
```

**Response Field Descriptions:**

| Field | Type | Description |
|-------|------|-------------|
| `ticketNum` | string | Unique ticket number identifier |
| `companyCode` | string | Company code associated with the ticket |
| `outletCode` | string | Outlet/store code |
| `outletName` | string | Outlet/store name |
| `issue` | string \| null | Issue title or description |
| `issueId` | integer \| null | Issue classification ID |
| `description` | string | Full description of the ticket |
| `assignDeveloper` | string | Developer assigned to resolve the ticket |
| `checkNotes` | string | Initial check/inspection notes |
| `devCheckNotes` | string | Developer's technical notes and findings |
| `solution` | string \| null | Applied solution or resolution |
| `statusId` | integer | Current status of the ticket |
| `createdBy` | string | User who created the ticket |
| `createdDate` | datetime | Ticket creation timestamp |
| `updatedBy` | string | User who last updated the ticket |
| `updatedDate` | datetime | Last update timestamp |
| `date` | datetime | Ticket date reference |
| `ref` | string \| null | External reference number |
| `refFrom` | string \| null | Reference source system |
| `userId` | string | Associated user ID |
| `productTypeId` | integer | Product type classification |

**Example Request:**

```bash
curl -X GET http://localhost:3005/return-ticket \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "productTypeID": 1,
    "issueType": "Stock Opname"
  }'
```

**Status Codes:**

| Code | Description |
|------|-------------|
| 200 | Success - Returns array of tickets |
| 400 | Bad Request - Invalid input parameters |
| 401 | Unauthorized - Missing or invalid JWT token |
| 500 | Internal Server Error |

---

### 2. Get Return Tickets with Analysis

Retrieves all return tickets for a product type and performs comprehensive analysis, grouping tickets by issue types and identifying patterns.

**Endpoint:**
```
GET /return-ticket/analyze
```

**Method:** `GET`

**Request Body:**

```json
{
  "productTypeID": 1
}
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `productTypeID` | integer | Yes | Product type identifier (must be >= 1) |

**Response:**

```json
{
  "totalTickets": 25,
  "issueGroups": [
    {
      "issueType": "Stock Opname",
      "count": 8,
      "mainProblem": "Stock count discrepancies during physical inventory verification",
      "tickets": [
        {
          "ticketNum": "TKT-001",
          "companyCode": "COMP001",
          "outletCode": "OUT001",
          "outletName": "Main Outlet",
          "issue": "Stock mismatch",
          "issueId": 1,
          "description": "System stock doesn't match physical count",
          "assignDeveloper": "dev@example.com",
          "checkNotes": "Confirmed stock count discrepancy",
          "devCheckNotes": "Stock Opname issue - quantity variance",
          "solution": "Adjusted stock levels",
          "statusId": 2,
          "createdBy": "user1@example.com",
          "createdDate": "2024-01-10T09:00:00Z",
          "updatedBy": "user2@example.com",
          "updatedDate": "2024-01-11T16:30:00Z",
          "date": "2024-01-10T00:00:00Z",
          "ref": "REF100",
          "refFrom": "POS System",
          "userId": "USER001",
          "productTypeId": 1
        }
      ]
    },
    {
      "issueType": "POS Data Upload",
      "count": 6,
      "mainProblem": "Issues with data synchronization from POS terminals",
      "tickets": []
    },
    {
      "issueType": "Other",
      "count": 11,
      "mainProblem": "General maintenance and miscellaneous issues",
      "tickets": []
    }
  ],
  "mostCommonIssue": {
    "issueType": "Stock Opname",
    "count": 8,
    "mainProblem": "Stock count discrepancies during physical inventory verification",
    "tickets": []
  }
}
```

**Response Field Descriptions:**

| Field | Type | Description |
|-------|------|-------------|
| `totalTickets` | integer | Total number of return tickets for the product type |
| `issueGroups` | array | Array of grouped issues with statistics |
| `issueGroups[].issueType` | string | Category/classification of the issue |
| `issueGroups[].count` | integer | Number of tickets in this issue group |
| `issueGroups[].mainProblem` | string | Summary of the main problem in this group |
| `issueGroups[].tickets` | array | Full ticket details in this issue group |
| `mostCommonIssue` | object | Issue group with the highest ticket count |
| `mostCommonIssue.issueType` | string | Most frequently occurring issue type |
| `mostCommonIssue.count` | integer | Number of tickets with the most common issue |

**Issue Categories:**

The analysis automatically categorizes tickets into the following issue types:

| Category | Description |
|----------|-------------|
| **Stock Opname** | Stock count discrepancies and inventory verification issues |
| **POS Data Upload** | Data synchronization and upload issues from POS terminals |
| **Menu Consistency** | Menu item definition and data inconsistencies |
| **Order Processing** | Order creation, modification, and cancellation issues |
| **Price Synchronization** | Price update and synchronization problems |
| **Promo/Loyalty** | Promotional campaign and loyalty program issues |
| **Outlet Configuration** | Outlet setup and configuration problems |
| **Staff/Permissions** | User access and permission management issues |
| **System Integration** | Integration with external systems and APIs |
| **Database/Data** | Data integrity and database-related problems |
| **UI/UX** | User interface and user experience issues |
| **Performance** | System performance and optimization issues |
| **Security** | Security vulnerabilities and access control issues |
| **Other** | Miscellaneous issues not in other categories |

**Example Request:**

```bash
curl -X GET http://localhost:3005/return-ticket/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "productTypeID": 1
  }'
```

**Status Codes:**

| Code | Description |
|------|-------------|
| 200 | Success - Returns analysis result |
| 400 | Bad Request - Invalid input parameters |
| 401 | Unauthorized - Missing or invalid JWT token |
| 500 | Internal Server Error |

---

## Error Responses

All endpoints return error responses in the following format:

```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable error message",
  "details": "Additional error details (if available)"
}
```

**Common Error Codes:**

| Code | HTTP Status | Description |
|------|------------|-------------|
| `INVALID_INPUT` | 400 | Invalid or missing required parameters |
| `UNAUTHORIZED` | 401 | Authentication failed or JWT token expired |
| `PRODUCT_NOT_FOUND` | 404 | Specified product type ID not found |
| `DATABASE_ERROR` | 500 | Database query error |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

---

## Request/Response Examples

### Example 1: Get all return tickets for a product type

**Request:**
```bash
curl -X GET http://localhost:3005/return-ticket \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "productTypeID": 1
  }'
```

**Response:** (200 OK)
```json
[
  {
    "ticketNum": "TKT-0001",
    "companyCode": "MC",
    "outletCode": "OUT-001",
    "outletName": "McDonald's Central",
    "description": "Stock inventory mismatch in burger station",
    "statusId": 2,
    "createdDate": "2024-04-25T08:15:00Z",
    "productTypeId": 1
  }
]
```

### Example 2: Get tickets with issue type filter

**Request:**
```bash
curl -X GET http://localhost:3005/return-ticket \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "productTypeID": 1,
    "issueType": "Stock Opname"
  }'
```

**Response:** (200 OK)
```json
[
  {
    "ticketNum": "TKT-0005",
    "companyCode": "MC",
    "outletCode": "OUT-002",
    "outletName": "McDonald's North",
    "devCheckNotes": "Stock Opname issue",
    "statusId": 2,
    "createdDate": "2024-04-23T10:30:00Z",
    "productTypeId": 1
  },
  {
    "ticketNum": "TKT-0008",
    "companyCode": "MC",
    "outletCode": "OUT-003",
    "outletName": "McDonald's South",
    "devCheckNotes": "Stock Opname issue",
    "statusId": 3,
    "createdDate": "2024-04-22T14:45:00Z",
    "productTypeId": 1
  }
]
```

### Example 3: Get analysis report

**Request:**
```bash
curl -X GET http://localhost:3005/return-ticket/analyze \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  -d '{
    "productTypeID": 1
  }'
```

**Response:** (200 OK)
```json
{
  "totalTickets": 15,
  "issueGroups": [
    {
      "issueType": "Stock Opname",
      "count": 6,
      "mainProblem": "Stock count discrepancies during physical inventory",
      "tickets": [...]
    },
    {
      "issueType": "POS Data Upload",
      "count": 4,
      "mainProblem": "Data synchronization issues from POS terminals",
      "tickets": [...]
    },
    {
      "issueType": "Other",
      "count": 5,
      "mainProblem": "General maintenance and miscellaneous issues",
      "tickets": [...]
    }
  ],
  "mostCommonIssue": {
    "issueType": "Stock Opname",
    "count": 6,
    "mainProblem": "Stock count discrepancies during physical inventory",
    "tickets": [...]
  }
}
```

---

## Request Validation Rules

### ProductTypeID
- **Type:** Integer
- **Required:** Yes
- **Constraints:** Must be >= 1
- **Example:** `1`, `5`, `100`

### IssueType
- **Type:** String
- **Required:** No (only for `/return-ticket` endpoint)
- **Constraints:** ASCII characters only
- **Examples:** `Stock Opname`, `POS Data Upload`, `Menu Consistency`

---

## Rate Limiting

Currently, no rate limiting is enforced. This may be implemented in future versions.

---

## Versioning

This API uses URL versioning. The current version is `v1` (implicit, no prefix).

Future versions will use paths like:
- `/api/v2/return-ticket`
- `/api/v2/return-ticket/analyze`

---

## Support

For issues, questions, or feature requests related to this API, please contact the development team or create an issue in the project repository.
