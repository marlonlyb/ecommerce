# Spec Técnica Frontend - Tienda de Ropa MVP

## 1. Objetivo

Definir la propuesta técnica de frontend para el MVP de tienda de ropa, cerrando la estrategia de sesión JWT, el flujo PayPal recomendado y la forma uniforme de consumir errores API.

## 2. Stack confirmado

- React 19
- Vite
- TypeScript
- React Router v7
- Fetch API
- `localStorage` solo para carrito
- `sessionStorage` para el JWT del MVP

## 3. Estructura recomendada

```txt
src/
  app/
    router/
    providers/
  features/
    catalog/
    product-detail/
    cart/
    auth/
    checkout/
    orders/
    admin-products/
    admin-orders/
  shared/
    api/
    components/
    guards/
    types/
    utils/
```

## 4. Decisiones técnicas cerradas

### 4.1 JWT y sesión

**Decisión**: el frontend guarda el access token JWT en `sessionStorage` y lo envía como Bearer token.

**Motivo**
- es compatible con el backend actual
- es más simple que implementar cookies HttpOnly + refresh token
- evita persistencia larga en `localStorage`

**Reglas frontend**
- restaurar sesión desde `sessionStorage` al boot
- resolver usuario actual con `GET /api/v1/private/me`
- si `401`, limpiar sesión y redirigir según contexto
- logout = limpiar token + user + estados derivados

### 4.2 Flujo PayPal recomendado

**Decisión**: el frontend no crea ni captura órdenes de PayPal directamente contra PayPal. Siempre orquesta al backend.

**Flujo UI/API**
1. usuario autenticado llega a checkout
2. frontend envía items a `POST /api/v1/private/orders/checkout/paypal`
3. backend devuelve `order.id` local y `paypal.order_id`
4. PayPal JS SDK usa ese `paypal.order_id`
5. en `onApprove`, frontend llama `POST /api/v1/private/orders/:id/paypal/capture`
6. si backend confirma `paid`, frontend limpia carrito y navega a detalle de `order`
7. si hay cancelación o error, la `order` queda en estado consistente para reintento o soporte

### 4.3 Manejo de errores

El frontend debe consumir un único formato:

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

**Regla de cliente**
- `shared/api/http.ts` parsea el envelope
- `shared/api/errors.ts` mapea a un `AppError`
- formularios usan `details` por campo cuando exista
- toasts o banners usan `message`

## 5. Rutas funcionales

### Públicas
- `/`
- `/products`
- `/products/:id`
- `/cart`
- `/login`
- `/register`

### Privadas cliente
- `/checkout`
- `/profile`
- `/profile/orders`
- `/profile/orders/:id`

### Privadas admin
- `/admin`
- `/admin/products`
- `/admin/products/new`
- `/admin/products/:id`
- `/admin/orders`
- `/admin/orders/:id`

## 6. Estado y persistencia

### 6.1 Estado local

- formularios
- loading de interacción
- selección de variante
- estados de PayPal checkout

### 6.2 Estado compartido

- sesión actual
- carrito
- filtros del catálogo en URL

### 6.3 Persistencia

- carrito → `localStorage`
- token JWT → `sessionStorage`
- perfil actual → memoria, derivado de `/private/me`

## 7. Tipos frontend mínimos

```ts
interface SessionUser {
  id: string;
  email: string;
  is_admin: boolean;
  created_at: string;
}

interface CartItem {
  product_id: string;
  product_name: string;
  product_image: string;
  variant_id: string;
  variant_sku: string;
  color: string;
  size: string;
  unit_price: number;
  quantity: number;
  available_stock: number;
}

interface ApiErrorDetail {
  field?: string;
  issue: string;
}

interface ApiErrorPayload {
  code: string;
  message: string;
  details?: ApiErrorDetail[];
  request_id?: string;
}
```

## 8. Capa API recomendada

```txt
shared/api/http.ts
shared/api/errors.ts
features/catalog/api.ts
features/auth/api.ts
features/checkout/api.ts
features/orders/api.ts
features/admin-products/api.ts
features/admin-orders/api.ts
```

### Reglas

- no usar `fetch()` directo dentro de componentes visuales
- adjuntar Authorization header solo si hay token
- mapear 401/403/409 de forma semántica
- centralizar parseo del nuevo error envelope

## 9. Checkout y UX

### 9.1 Reglas de checkout

- checkout requiere auth para iniciar pago
- el botón de pago se habilita solo si hay items válidos
- el frontend nunca calcula el total definitivo como fuente de verdad
- el carrito se limpia solo si la captura backend devuelve `paid`

### 9.2 Estados UX mínimos

- `idle`
- `creating_paypal_order`
- `awaiting_paypal_approval`
- `capturing_payment`
- `success`
- `error`

## 10. Guards

- `ProtectedRoute`: checkout, perfil y `orders`
- `AdminRoute`: todo `/admin`
- guardar `returnTo` para volver a checkout tras login

## 11. Contrato técnico derivado

Las definiciones cerradas se apoyan en:

- `API-Contract-TiendaRopa.md`
- `Spec-Backend-TiendaRopa.md`
- `DB-Schema-TiendaRopa.md`

El frontend debe alinearse con esos artefactos y usar `orders` como naming único.
