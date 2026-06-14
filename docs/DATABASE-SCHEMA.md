# Database Schema & Persistence Documentation

## Overview

This document describes how all objects in the Delícias da Lú API are persisted to **Google Cloud Firestore**, a NoSQL cloud database. Firestore uses a hierarchical document structure with collections containing documents.

**Key Technologies:**
- **Database**: Google Cloud Firestore
- **Client Library**: `cloud.google.com/go/firestore`
- **Project ID**: Retrieved from `GCP_PROJECT_ID` environment variable (default: `project-4419255d-5de2-41f6-82b`)

---

## Architecture Overview

### Connection & Initialization

```go
// From cmd/main.go
projectID := os.Getenv("GCP_PROJECT_ID")
client, err := firestore.NewClient(ctx, projectID)
```

All repositories receive the Firestore client and use it to perform CRUD operations. The connection is established once during application bootstrap and reused across all requests.

---

## Collection Structure

### 1. **Collection: `users`**

Stores admin user accounts for authentication and authorization.

**Firestore Path:** `/users/{userId}`

**Document Schema:**

```json
{
  "id": "string",
  "username": "string",
  "email": "string",
  "password": "string (hashed)",
  "role": "string (admin | manager)",
  "lastLogin": "timestamp",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

**Example Document:**

```json
{
  "id": "user-001",
  "username": "juliosshoji",
  "email": "julio@delicias-da-lu.com.br",
  "password": "$2a$10$...", // bcrypt hash
  "role": "admin",
  "lastLogin": "2026-06-14T10:30:00Z",
  "createdAt": "2026-01-01T00:00:00Z",
  "updatedAt": "2026-06-14T10:30:00Z"
}
```

**Persistence Operations:**

| Operation | Method | Implementation |
|-----------|--------|-----------------|
| Get by username | `Query` | `Collection("users").Where("username", "==", username).Documents()` |
| Get by ID | `Document` | `Collection("users").Doc(id).Get()` |
| Create | `Set` | `Collection("users").Doc(userId).Set(userObject)` |
| Update last login | `Update` | `Collection("users").Doc(id).Update({"lastLogin": now})` |

**Access Pattern:** Users are queried by username during login, then cached in JWT tokens

---

### 2. **Collection: `menu`**

Stores all available menu items (bolos, doces, etc.).

**Firestore Path:** `/menu/{itemId}`

**Document Schema:**

```json
{
  "id": "string",
  "name": "string",
  "category": "string",
  "price": "float64",
  "unit": "string",
  "image": "string (URL path)",
  "description": "string",
  "active": "boolean",
  "order": "integer (display order)",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

**Example Document:**

```json
{
  "id": "frutas-abacaxi-creme",
  "name": "Abacaxi com Creme",
  "category": "Sabores de Frutas",
  "price": 65.00,
  "unit": "kg",
  "image": "/images/bolos/frutas-abacaxi-creme.jpg",
  "description": "Bolo macio com camadas de abacaxi fresco e creme",
  "active": true,
  "order": 1,
  "createdAt": "2026-01-15T09:00:00Z",
  "updatedAt": "2026-06-10T14:30:00Z"
}
```

**Persistence Operations:**

| Operation | Method | Implementation |
|-----------|--------|-----------------|
| Get all | `OrderBy` | `Collection("menu").OrderBy("order", Asc).Documents()` |
| Get by ID | `Document` | `Collection("menu").Doc(id).Get()` |
| Create | `Set` | `Collection("menu").Doc(itemId).Set(menuItem)` |
| Update | `Set` | `Collection("menu").Doc(id).Set(updatedItem)` |
| Delete | `Delete` | `Collection("menu").Doc(id).Delete()` |

**Filtering:** Performed in-memory after retrieval by active status and category
**Sorting:** By `order` field ascending for display

---

### 3. **Collection: `cakebuilder`**

Stores custom cake builder components organized by type. Each document represents a single component option.

**Firestore Path:** `/cakebuilder/{componentId}`

**Document Schema:**

```json
{
  "id": "string",
  "name": "string",
  "type": "string (massa | recheio | cobertura | decoracao)",
  "price": "float64",
  "image": "string (URL path)",
  "active": "boolean",
  "order": "integer (display order within type)",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

**Example Documents:**

```json
{
  "id": "massa-chocolate",
  "name": "Chocolate",
  "type": "massa",
  "price": 15.50,
  "image": "/images/cake-builder/massas/chocolate.jpg",
  "active": true,
  "order": 1,
  "createdAt": "2026-01-20T10:00:00Z",
  "updatedAt": "2026-06-10T14:30:00Z"
}
```

**Persistence Operations:**

| Operation | Method | Implementation |
|-----------|--------|-----------------|
| Get all | `Documents` | `Collection("cakebuilder").Documents()` |
| Get by type | `Where` | `Collection("cakebuilder").Where("type", "==", type).Documents()` |
| Get by ID | `Document` | `Collection("cakebuilder").Doc(id).Get()` |
| Create | `Set` | `Collection("cakebuilder").Doc(componentId).Set(component)` |
| Update | `Set` | `Collection("cakebuilder").Doc(id).Set(updatedComponent)` |
| Delete | `Delete` | `Collection("cakebuilder").Doc(id).Delete()` |

**Types:** masa, recheio, cobertura, decoracao
**Display:** Grouped by type on frontend, sorted by `order` field

---

### 4. **Collection: `orders`**

Stores all customer orders with full order details including items and customer information.

**Firestore Path:** `/orders/{orderId}`

**Document Schema:**

```json
{
  "id": "string",
  "items": [
    {
      "type": "string (menu | cakeBuilder)",
      "menuItemId": "string (if type=menu)",
      "cakeCustomization": {
        "massa": "string",
        "recheio": "string",
        "cobertura": "string",
        "decoracao": "string"
      },
      "quantity": "integer",
      "unitPrice": "float64",
      "subtotal": "float64"
    }
  ],
  "customerInfo": {
    "name": "string",
    "phone": "string",
    "email": "string",
    "deliveryDate": "string (YYYY-MM-DD)",
    "notes": "string (allergies, special instructions)"
  },
  "status": "string (pending | confirmed | preparing | ready | delivered | cancelled)",
  "totalPrice": "float64",
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

**Example Document:**

```json
{
  "id": "order-20260614-001",
  "items": [
    {
      "type": "menu",
      "menuItemId": "frutas-abacaxi-creme",
      "quantity": 1,
      "unitPrice": 65.00,
      "subtotal": 65.00
    },
    {
      "type": "cakeBuilder",
      "cakeCustomization": {
        "massa": "massa-chocolate",
        "recheio": "recheio-morango",
        "cobertura": "cobertura-brigadeiro",
        "decoracao": "decoracao-confeitos"
      },
      "quantity": 2,
      "unitPrice": 45.99,
      "subtotal": 91.98
    }
  ],
  "customerInfo": {
    "name": "João Silva",
    "phone": "11999999999",
    "email": "joao@example.com",
    "deliveryDate": "2026-06-20",
    "notes": "Sem amendoim - alérgico"
  },
  "status": "pending",
  "totalPrice": 156.98,
  "createdAt": "2026-06-14T14:30:00Z",
  "updatedAt": "2026-06-14T14:30:00Z"
}
```

**Persistence Operations:**

| Operation | Method | Implementation |
|-----------|--------|-----------------|
| Get all (paginated) | `OrderBy` | `Collection("orders").OrderBy("createdAt", Desc).Documents()` |
| Get by ID | `Document` | `Collection("orders").Doc(id).Get()` |
| Create | `Set` | `Collection("orders").Doc(orderId).Set(order)` |
| Update status | `Update` | `Collection("orders").Doc(id).Update({"status": newStatus})` |

**Sorting:** By creation date descending (newest first)
**Filtering:** By status and pagination parameters applied in-memory
**Nested Data:** Items and customer info stored as sub-documents within the order

---

### 5. **Collection: `contacts`**

Stores business contact information. This is typically a singleton collection (single document).

**Firestore Path:** `/contacts/main` (or single document in collection)

**Document Schema:**

```json
{
  "whatsapp": {
    "number": "string (country code + number)",
    "link": "string (wa.me URL)"
  },
  "email": "string",
  "instagram": "string",
  "address": "string",
  "phone": "string"
}
```

**Example Document:**

```json
{
  "whatsapp": {
    "number": "5511987654321",
    "link": "https://wa.me/5511987654321"
  },
  "email": "contato@delicias-da-lu.com.br",
  "instagram": "delicias_da_lu",
  "address": "Rua das Flores, 123 - Centro",
  "phone": "1133334444"
}
```

**Persistence Operations:**

| Operation | Method | Implementation |
|-----------|--------|-----------------|
| Get | `Document` | `Collection("contacts").Doc("main").Get()` |
| Update | `Set` | `Collection("contacts").Doc("main").Set(updatedContacts)` |

**Pattern:** Treated as a singleton for simplicity

---

### 6. **Collection: `config`**

Stores aggregated site configuration combining menu, cake builder, and contact information.

**Firestore Path:** `/config/main` (or single document in collection)

**Document Schema:**

```json
{
  "menu": {
    "items": [ /* array of MenuItem objects */ ],
    "sectionLabels": {
      "bolos": "🍰 Bolos",
      "doces": "🍬 Doces"
    },
    "customSections": []
  },
  "cakeBuilder": {
    "massas": [ /* CakeBuilderComponent[] */ ],
    "recheios": [ /* CakeBuilderComponent[] */ ],
    "coberturas": [ /* CakeBuilderComponent[] */ ],
    "decoracoes": [ /* CakeBuilderComponent[] */ ]
  },
  "contacts": { /* Contact object */ }
}
```

**Persistence Pattern:** 
- Not directly written as a single aggregated document
- Instead, publicly retrieved as computed response combining data from `menu`, `cakebuilder`, and `contacts` collections
- For admin, full config can include additional metadata

---

### 7. **Collection: `error-types`**

Stores documentation and HTML descriptions for each error type. Used for error tracking and user-facing error documentation.

**Firestore Path:** `/error-types/{errorIdentifier}`

**Document Schema:**

```json
{
  "html": "string (HTML documentation)",
  "updatedAt": "timestamp"
}
```

**Example Document:**

```json
{
  "html": "<h1>Not Found Error</h1><p>The requested resource was not found...</p>",
  "updatedAt": "2026-06-01T00:00:00Z"
}
```

**Persistence Operations:**

| Operation | Method | Implementation |
|-----------|--------|-----------------|
| Get by identifier | `Document` | `Collection("error-types").Doc(identifier).Get()` |
| Create/Update | `Set` | `Collection("error-types").Doc(identifier).Set(errorType)` |
| Delete | `Delete` | `Collection("error-types").Doc(identifier).Delete()` |

---

### 8. **Collection: `error-instances`**

Stores detailed error occurrence records for debugging and analysis. Each request that encounters an error creates a document here with trace ID as the identifier.

**Firestore Path:** `/error-instances/{traceId}`

**Document Schema:**

```json
{
  "trace_id": "string (unique request trace)",
  "type": "string (error type identifier)",
  "title": "string",
  "status": "integer (HTTP status code)",
  "request_url": "string (full request path)",
  "request_date": "timestamp",
  "request_body": "string (request body, if any)",
  "user_agent": "string",
  "html": "string (formatted error HTML)"
}
```

**Example Document:**

```json
{
  "trace_id": "a1b2c3d4e5f6g7h8",
  "type": "not-found",
  "title": "Menu Item Not Found",
  "status": 404,
  "request_url": "GET /v1/menu/items/invalid-id-99999",
  "request_date": "2026-06-14T14:30:00Z",
  "request_body": "",
  "user_agent": "Mozilla/5.0...",
  "html": "<h1>Not Found</h1><p>No menu item found with ID: invalid-id-99999</p>"
}
```

**Persistence Operations:**

| Operation | Method | Implementation |
|-----------|--------|-----------------|
| Get by trace ID | `Document` | `Collection("error-instances").Doc(traceId).Get()` |
| Create | `Set` | `Collection("error-instances").Doc(traceId).Set(errorInstance)` |

**Automatic Recording:** Error middleware automatically records errors when they occur

---

## Data Relationships & Normalization

### Foreign Key Patterns

Firestore uses **denormalization and reference patterns** instead of traditional foreign keys:

#### 1. **Menu Items in Orders**
```
Order.items[].menuItemId → menu collection
```
- Order stores menu item ID, not full object
- UI fetches menu details separately or caches locally
- Allows menu prices to change without updating orders

#### 2. **Cake Builder Components in Orders**
```
Order.items[].cakeCustomization.{massa,recheio,cobertura,decoracao} → cakebuilder collection
```
- Each component ID references a document in cakebuilder
- Allows custom cakes to be reconstructed with current prices

#### 3. **Aggregated Data in Config**
```
config → menu + cakebuilder + contacts
```
- Config retrieval fetches all three collections
- Computed response combines data (not a stored join)

---

## Common Persistence Patterns

### Pattern 1: Document Writes (Create/Update)

```go
// Set entire document (create or overwrite)
_, err := r.client.Collection("menu").Doc(item.ID).Set(ctx, item)

// Update specific fields without overwriting
_, err := r.client.Collection("orders").Doc(id).Update(ctx, []firestore.Update{
  {Path: "status", Value: newStatus},
  {Path: "updatedAt", Value: time.Now()},
})
```

### Pattern 2: Document Reads

```go
// Get single document
doc, err := r.client.Collection("users").Doc(id).Get(ctx)
var user User
doc.DataTo(&user)

// Get multiple documents with query
docs, err := r.client.Collection("menu").
  OrderBy("order", firestore.Asc).
  Documents(ctx).GetAll()
```

### Pattern 3: Filtering & Querying

```go
// Query with condition
docs, err := r.client.Collection("users").
  Where("username", "==", username).
  Documents(ctx).GetAll()

// In-memory filtering (after retrieval)
for _, doc := range docs {
  var item MenuItem
  doc.DataTo(&item)
  if (active == nil || item.Active == *active) && 
     (category == "" || item.Category == category) {
    items = append(items, item)
  }
}
```

### Pattern 4: Ordering & Pagination

```go
// Order by field
query := r.client.Collection("orders").OrderBy("createdAt", firestore.Desc)

// Pagination in-memory
var orders []Order
for i, doc := range docs {
  if i < offset || i >= offset+limit {
    continue
  }
  var ord Order
  doc.DataTo(&ord)
  orders = append(orders, ord)
}
```

---

## Indexes

Firestore automatically creates indexes for:
- Single field queries (e.g., `WHERE username == "x"`)
- Document ID lookups

**Composite indexes** may be needed for:
- `menu.OrderBy("order").Desc` + `WHERE active == true`
- `orders.OrderBy("createdAt").Desc` + `WHERE status == "pending"`

These are typically auto-created on first query or can be manually created in Firestore console.

---

## Data Consistency & Transactions

**Current Implementation:** Each repository method performs individual document operations without transactions.

**Implications:**
- Creating an order reads menu/cake items but doesn't lock them
- Concurrent updates don't use pessimistic locking
- Eventually consistent

**For Multi-Document Atomicity:** Use Firestore transactions:

```go
err := client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
  // Multiple operations, all succeed or all fail
  return nil
})
```

---

## Performance Considerations

### Indexing Strategy

| Collection | Key Indexes | Composite Indexes |
|-----------|-------------|-------------------|
| users | username | - |
| menu | order, active | order + active |
| cakebuilder | type, order | type + order |
| orders | createdAt, status | createdAt + status |

### Query Optimization

1. **Firestore prefers indexed equality queries** - place equality conditions first
2. **Avoid full table scans** - use WHERE clauses whenever possible
3. **Pagination** - retrieve only needed documents with limits
4. **In-memory filtering** - filter results after retrieval for complex logic

### Cost Optimization

- Read cost: 1 document read ≈ $0.06 per 100K reads (approximate)
- Write cost: 1 document write ≈ $0.18 per 100K writes
- Strategies:
  - Batch related data to reduce queries
  - Cache config/menu data on client side
  - Archive old error instances periodically

---

## Backup & Recovery

**Automatic Backups:**
- Firestore provides built-in snapshots
- Enable from Cloud Console for production
- Can restore to specific point-in-time

**Manual Export:**
```bash
gcloud firestore export gs://bucket-name/export-2026-06-14
```

**Recovery:**
```bash
gcloud firestore import gs://bucket-name/export-2026-06-14
```

---

## Environment Configuration

**Required Environment Variables:**

```bash
# Google Cloud Project
GCP_PROJECT_ID=project-4419255d-5de2-41f6-82b

# Google Cloud Credentials (one of):
# - Service account JSON file path
# - Application Default Credentials (if running in GCP)
# - GOOGLE_APPLICATION_CREDENTIALS=path/to/service-account.json
```

---

## Troubleshooting Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| "Collection not found" | Collection doesn't exist yet | Create first document in collection |
| "Missing index" | Complex query without index | Create composite index in Firestore console |
| "Permission denied" | IAM/Security Rules | Check Firestore security rules |
| "Timeout" | Large data transfer or network latency | Add pagination, reduce document size |
| "Document decode error" | Schema mismatch | Verify struct tags match field names |

---

## Future Considerations

1. **Sub-collections:** Consider moving order items to sub-collections:
   ```
   /orders/{orderId}/items/{itemId}
   ```

2. **Denormalization:** Cache computed values to reduce queries

3. **Archiving:** Move old orders to archive collection for performance

4. **Real-time Updates:** Use Firestore listeners for live order status updates

5. **Audit Trail:** Add collection for tracking all changes for compliance
