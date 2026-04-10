# Spec Técnica Backend - Tienda de Ropa MVP

## 1. Objetivo

Definir la base técnica de backend para el MVP de tienda de ropa sobre Go + Echo + PostgreSQL, cerrando JWT, PayPal, formato de error y esquema DB inicial con naming oficial `orders`.

## 2. Principios

- simpleza operativa antes que sofisticación
- validación real en backend, nunca solo en frontend
- compatibilidad incremental con el backend existente
- stock, pago y persistencia final bajo control transaccional
- naming consistente: `orders` en API, documentación y nuevas piezas técnicas

## 3. Estado actual del repositorio

El repositorio ya tiene backend Go + Echo y PostgreSQL, pero hoy mezcla un modelo legado basado en:

- `purchase_orders`
- `invoices`
- productos sin tabla de variantes
- respuestas HTTP con envelope distinto al contrato deseado

### Implicación

La implementación inicial del MVP debe ser **incremental y aditiva**: se define el modelo final de `orders`, pero sin intentar reescribir todo el backend legado en un solo paso.

## 4. Decisiones técnicas cerradas

### 4.1 JWT

**Decisión**: access token JWT único, firmado con HS256, expiración de 12 horas, enviado por `Authorization: Bearer <token>`.

**Frontend**: guardar en `sessionStorage` para restauración de sesión en la pestaña activa.

**No incluido en MVP**
- refresh token
- cookie HttpOnly
- blacklist de tokens

**Tradeoff**
- a favor: muy simple, compatible con el backend actual
- en contra: menor robustez frente a sesiones largas y rotación segura

### 4.2 PayPal

**Decisión**: create order y capture desde backend; webhook solo como reconciliación.

**Flujo backend recomendado**
1. recibir items del checkout autenticado
2. cargar variantes reales desde DB
3. validar stock y estado de producto
4. recalcular subtotal y total
5. crear `orders` + `order_items` con `status=pending_payment`
6. crear order en PayPal y guardar `paypal_order_id`
7. al aprobar frontend, capturar desde backend
8. confirmar estado `captured`, persistir `paypal_capture_id`, `paid_at` y `status=paid`
9. descontar stock en la misma transacción de confirmación final

### 4.3 Error format

**Decisión**: estandarizar a un único envelope:

```json
{
  "error": {
    "code": "validation_error",
    "message": "Los datos enviados no son válidos",
    "details": [
      { "field": "email", "issue": "required" }
    ],
    "request_id": "req_123"
  }
}
```

## 5. Arquitectura recomendada

### 5.1 Capas

#### Handlers
- parsean request
- validan shape básico
- invocan services
- transforman respuesta HTTP

#### Services
- auth y autorización contextual
- validación de stock
- recálculo de totales
- transiciones de `orders`
- coordinación con PayPal

#### Repositories
- acceso SQL
- queries por tabla agregada
- transacciones para capture/finalización de order

## 6. Esquema DB base del MVP

El esquema final queda especificado en `DB-Schema-TiendaRopa.md`.

### 6.1 Tablas mínimas obligatorias

- `users`
- `products`
- `product_variants`
- `orders`
- `order_items`

### 6.2 Decisiones de modelado

- `products` representa catálogo publicable
- `product_variants` concentra precio y stock real vendible
- `orders` guarda estado de negocio y estado de pago
- `order_items` guarda snapshot mínimo para historial
- no se agregan tablas extra de carrito ni refresh tokens en MVP

## 7. Responsabilidades por módulo

### 7.1 Catálogo público

**Handlers**
- `GET /api/v1/public/products`
- `GET /api/v1/public/products/:id`

**Services**
- aplicar filtros
- exponer solo productos activos
- resolver `price_from`, colores y talles disponibles

### 7.2 Auth

**Handlers**
- `POST /api/v1/public/register`
- `POST /api/v1/public/login`
- `GET /api/v1/private/me`

**Services**
- crear usuario
- comparar password hash
- emitir JWT con `sub`, `email`, `is_admin`, `iat`, `exp`

### 7.3 Checkout y `orders`

**Handlers**
- `POST /api/v1/private/orders/checkout/paypal`
- `POST /api/v1/private/orders/:id/paypal/capture`
- `GET /api/v1/private/orders`
- `GET /api/v1/private/orders/:id`

**Services**
- validar items
- cargar variantes reales
- validar stock suficiente
- crear `orders` pendientes
- crear/capturar en PayPal
- confirmar stock y pago

### 7.4 Admin

**Handlers**
- `GET /api/v1/admin/products`
- `POST /api/v1/admin/products`
- `PUT /api/v1/admin/products/:id`
- `PATCH /api/v1/admin/products/:id/status`
- `GET /api/v1/admin/orders`
- `GET /api/v1/admin/orders/:id`
- `PATCH /api/v1/admin/orders/:id/status`

## 8. Reglas de seguridad mínimas

- password solo como hash
- nunca confiar en precio o total enviados por cliente
- nunca confiar en `is_admin` del request
- validar ownership de `orders`
- validar capture PayPal contra order local
- no exponer productos inactivos en catálogo público

## 9. Estados y transiciones

### 9.1 `orders.status`

- `pending_payment`
- `paid`
- `payment_failed`
- `cancelled`
- `refunded`

### 9.2 `orders.payment_status`

- `pending`
- `approved`
- `captured`
- `failed`
- `refunded`

### 9.3 Transiciones válidas mínimas

- `pending_payment` → `paid`
- `pending_payment` → `payment_failed`
- `pending_payment` → `cancelled`
- `paid` → `refunded`

## 10. Error mapping mínimo

| Código | HTTP | Caso |
|---|---:|---|
| `validation_error` | 400 | payload inválido |
| `authentication_required` | 401 | token ausente o inválido |
| `invalid_credentials` | 401 | login fallido |
| `forbidden` | 403 | acceso sin permisos |
| `not_found` | 404 | recurso inexistente |
| `stock_insufficient` | 409 | sin stock suficiente |
| `order_state_invalid` | 409 | transición inválida |
| `paypal_capture_failed` | 422 | PayPal no confirmó la captura |
| `unexpected_error` | 500 | error no controlado |

## 11. Ajuste contrato ↔ esquema DB

- `product_variants.price` es la fuente única para cálculo de checkout
- `orders` reemplaza conceptualmente al agregado legado `purchase_orders + invoices`
- `orders.paypal_order_id` y `orders.paypal_capture_id` soportan auditoría mínima
- `order_items.line_total` se persiste para historial y trazabilidad

## 12. Punto mínimo y seguro de arranque

El repositorio todavía no tiene el módulo final de `orders`, pero sí un backend funcional y middleware JWT básico. Por eso, el arranque más seguro es:

1. documentar decisiones cerradas
2. introducir el esquema SQL base del MVP
3. introducir tipos Go alineados al contrato nuevo
4. abrir alias/ruta inicial con naming `orders`

Esto permite avanzar sin romper de golpe el flujo legado.
