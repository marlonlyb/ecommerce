# API Contract - Tienda de Ropa MVP

## 1. Propósito

Contrato técnico final del MVP de tienda de ropa, alineado con el naming oficial `orders`, el esquema DB inicial y el stack real del proyecto: frontend React 19 + Vite + TypeScript y backend Go + Echo + PostgreSQL.

## 2. Decisiones técnicas cerradas

### 2.1 JWT

**Decisión**: usar un único access token JWT firmado con **HS256**, enviado por `Authorization: Bearer <token>`, con expiración corta/intermedia (**12 horas**) y almacenamiento en **sessionStorage** del frontend.

**Por qué esta estrategia**
- encaja con el backend actual basado en Bearer token
- evita agregar refresh tokens, cookies HttpOnly y tablas extra en el MVP
- reduce exposición frente a `localStorage` persistente
- mantiene operación simple para un despliegue pequeño

**Tradeoff explícito**
- más simple de implementar y operar
- menos robusto que access token corto + refresh token HttpOnly

**Cierre**: para este MVP se acepta el tradeoff y **no se implementa refresh token**.

### 2.2 PayPal

**Decisión**: usar flujo **backend-driven create + backend-driven capture**, dejando el webhook de PayPal como mecanismo de reconciliación secundaria, no como camino principal del checkout.

**Resumen del flujo**
1. frontend autenticado envía checkout a backend
2. backend valida stock, recalcula totales y crea `orders` en estado `pending_payment`
3. backend crea la order en PayPal y guarda `paypal_order_id`
4. frontend aprueba el pago con PayPal JS SDK
5. frontend notifica aprobación al backend
6. backend captura con PayPal, valida monto/estado y actualiza `orders` a `paid`
7. backend descuenta stock dentro de la confirmación final

### 2.3 Formato de errores API

**Decisión**: usar envelope único:

```json
{
  "error": {
    "code": "validation_error",
    "message": "Los datos enviados no son válidos",
    "details": [
      {
        "field": "email",
        "issue": "required"
      }
    ],
    "request_id": "req_123"
  }
}
```

`details` es opcional. `request_id` es opcional pero recomendado para trazabilidad.

## 3. Convenciones generales

- **Base path**: `/api/v1`
- **Formato**: JSON
- **Auth**: `Authorization: Bearer <jwt>`
- **Timestamps API**: ISO 8601 UTC
- **Moneda MVP**: `USD`
- **Naming oficial**: el recurso de compra/pedido se llama `orders`
- **Checkout**: el backend recalcula total siempre; el frontend nunca define el total final confiable

### 3.1 Respuesta exitosa

```json
{
  "data": {}
}
```

### 3.2 Respuesta paginada

```json
{
  "data": {
    "items": [],
    "pagination": {
      "page": 1,
      "limit": 12,
      "total": 0,
      "total_pages": 0
    }
  }
}
```

### 3.3 Respuesta de error

```json
{
  "error": {
    "code": "stock_insufficient",
    "message": "Una o más variantes ya no tienen stock suficiente",
    "details": [
      {
        "field": "items[0].variant_id",
        "issue": "available_stock=1"
      }
    ]
  }
}
```

## 4. Auth y autorización

### 4.1 Tipos de acceso

- **Público**: catálogo, registro, login
- **Cliente autenticado**: `me`, checkout y `orders` propias
- **Admin**: gestión de catálogo y `orders`

### 4.2 Claims mínimas del JWT

- `sub`: user id
- `email`
- `is_admin`
- `iat`
- `exp`

### 4.3 Reglas

- toda ruta `/private/*` requiere JWT válido
- toda ruta `/admin/*` requiere JWT válido e `is_admin = true`
- un cliente solo puede leer sus propias `orders`
- el JWT se invalida por expiración; logout en MVP es client-side

## 5. Modelos de respuesta

### 5.1 User

```json
{
  "id": "usr_1",
  "email": "cliente@email.com",
  "is_admin": false,
  "created_at": "2026-04-10T10:00:00Z"
}
```

### 5.2 ProductSummary

```json
{
  "id": "prod_1",
  "name": "Chaqueta Denim Azul",
  "slug": "chaqueta-denim-azul",
  "category": "chaquetas",
  "brand": "Marca X",
  "images": ["https://cdn.example.com/p1.jpg"],
  "active": true,
  "price_from": 49.99,
  "available_colors": ["azul", "negro"],
  "available_sizes": ["S", "M", "L"]
}
```

### 5.3 ProductDetail

```json
{
  "id": "prod_1",
  "name": "Chaqueta Denim Azul",
  "slug": "chaqueta-denim-azul",
  "description": "Chaqueta clásica de denim.",
  "category": "chaquetas",
  "brand": "Marca X",
  "images": ["https://cdn.example.com/p1.jpg"],
  "active": true,
  "variants": [
    {
      "id": "var_1",
      "product_id": "prod_1",
      "sku": "CHA-DEN-AZU-M",
      "color": "azul",
      "size": "M",
      "price": 49.99,
      "stock": 8,
      "image_url": "https://cdn.example.com/p1.jpg"
    }
  ]
}
```

### 5.4 Order

```json
{
  "id": "ord_1",
  "user_id": "usr_1",
  "status": "paid",
  "payment_provider": "paypal",
  "payment_status": "captured",
  "currency": "USD",
  "subtotal": 99.98,
  "total": 99.98,
  "paypal_order_id": "2GG279541U471931P",
  "paypal_capture_id": "3GG27954AB1234567",
  "paid_at": "2026-04-10T10:03:00Z",
  "created_at": "2026-04-10T10:00:00Z",
  "items": [
    {
      "id": "item_1",
      "product_id": "prod_1",
      "product_name": "Chaqueta Denim Azul",
      "variant_id": "var_1",
      "variant_sku": "CHA-DEN-AZU-M",
      "color": "azul",
      "size": "M",
      "unit_price": 49.99,
      "quantity": 2,
      "line_total": 99.98
    }
  ]
}
```

## 6. Endpoints públicos

### GET `/api/v1/public/products`

Lista productos activos con filtros.

**Query params**
- `search?: string`
- `category?: string`
- `color?: string`
- `size?: string`
- `min_price?: number`
- `max_price?: number`
- `page?: number`
- `limit?: number`
- `sort?: price_asc | price_desc | newest`

### GET `/api/v1/public/products/:id`

Devuelve detalle completo del producto activo.

### POST `/api/v1/public/register`

**Request**

```json
{
  "email": "cliente@email.com",
  "password": "secret123",
  "confirm_password": "secret123"
}
```

**Response 201**

```json
{
  "data": {
    "user": {
      "id": "usr_1",
      "email": "cliente@email.com",
      "is_admin": false,
      "created_at": "2026-04-10T10:00:00Z"
    }
  }
}
```

### POST `/api/v1/public/login`

**Request**

```json
{
  "email": "cliente@email.com",
  "password": "secret123"
}
```

**Response 200**

```json
{
  "data": {
    "user": {
      "id": "usr_1",
      "email": "cliente@email.com",
      "is_admin": false,
      "created_at": "2026-04-10T10:00:00Z"
    },
    "token": "jwt-token",
    "expires_in": 43200
  }
}
```

## 7. Endpoints privados cliente

### GET `/api/v1/private/me`

Devuelve identidad mínima del usuario autenticado.

### POST `/api/v1/private/orders/checkout/paypal`

Valida checkout, recalcula total, crea la `order` en estado `pending_payment`, crea la order en PayPal y devuelve el identificador de PayPal.

**Request**

```json
{
  "items": [
    {
      "variant_id": "var_1",
      "quantity": 2
    }
  ]
}
```

**Response 201**

```json
{
  "data": {
    "order": {
      "id": "ord_1",
      "user_id": "usr_1",
      "status": "pending_payment",
      "payment_provider": "paypal",
      "payment_status": "pending",
      "currency": "USD",
      "subtotal": 99.98,
      "total": 99.98,
      "paypal_order_id": "2GG279541U471931P",
      "created_at": "2026-04-10T10:00:00Z",
      "items": [
        {
          "id": "item_1",
          "product_id": "prod_1",
          "product_name": "Chaqueta Denim Azul",
          "variant_id": "var_1",
          "variant_sku": "CHA-DEN-AZU-M",
          "color": "azul",
          "size": "M",
          "unit_price": 49.99,
          "quantity": 2,
          "line_total": 99.98
        }
      ]
    },
    "paypal": {
      "order_id": "2GG279541U471931P"
    }
  }
}
```

### POST `/api/v1/private/orders/:id/paypal/capture`

Captura la order de PayPal desde backend y finaliza la `order` local.

**Request**

```json
{
  "paypal_order_id": "2GG279541U471931P"
}
```

**Response 200**

```json
{
  "data": {
    "order": {
      "id": "ord_1",
      "user_id": "usr_1",
      "status": "paid",
      "payment_provider": "paypal",
      "payment_status": "captured",
      "currency": "USD",
      "subtotal": 99.98,
      "total": 99.98,
      "paypal_order_id": "2GG279541U471931P",
      "paypal_capture_id": "3GG27954AB1234567",
      "paid_at": "2026-04-10T10:03:00Z",
      "created_at": "2026-04-10T10:00:00Z",
      "items": []
    }
  }
}
```

### GET `/api/v1/private/orders`

Lista las `orders` del usuario autenticado.

### GET `/api/v1/private/orders/:id`

Devuelve detalle de una `order` propia.

## 8. Endpoints admin

### GET `/api/v1/admin/products`
### POST `/api/v1/admin/products`
### PUT `/api/v1/admin/products/:id`
### PATCH `/api/v1/admin/products/:id/status`

### GET `/api/v1/admin/orders`

Query params opcionales:
- `status`
- `payment_status`
- `page`
- `limit`

### GET `/api/v1/admin/orders/:id`

### PATCH `/api/v1/admin/orders/:id/status`

**Request**

```json
{
  "status": "cancelled"
}
```

## 9. Estados y validaciones de negocio

### 9.1 Estados de `orders`

- `pending_payment`
- `paid`
- `payment_failed`
- `cancelled`
- `refunded`

### 9.2 Estados de pago

- `pending`
- `approved`
- `captured`
- `failed`
- `refunded`

### 9.3 Reglas críticas

- no crear `orders` sin items
- `quantity` debe ser mayor a 0
- el backend recalcula subtotal y total desde `product_variants.price`
- la captura debe validar monto, currency y ownership de la `order`
- el stock solo se descuenta cuando la captura queda confirmada
- el usuario solo puede consultar sus propias `orders`

## 10. Códigos mínimos de error

- `validation_error` → 400
- `authentication_required` → 401
- `invalid_credentials` → 401
- `forbidden` → 403
- `not_found` → 404
- `product_inactive` → 409
- `stock_insufficient` → 409
- `order_state_invalid` → 409
- `paypal_capture_failed` → 422
- `unexpected_error` → 500

## 11. Alineación contrato ↔ DB

- `users.id` ↔ `sub` del JWT
- `products.id` ↔ `product_variants.product_id`
- `product_variants.id` ↔ `order_items.variant_id`
- `orders.user_id` ↔ owner de la compra
- `orders.paypal_order_id` y `orders.paypal_capture_id` soportan reconciliación del pago
- `order_items` guarda snapshot mínimo de nombre, sku, color, size y precio para preservar historial

La estructura detallada de tablas, constraints e índices queda cerrada en `DB-Schema-TiendaRopa.md`.
