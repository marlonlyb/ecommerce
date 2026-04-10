# DB Schema - Tienda de Ropa MVP

## 1. Objetivo

Definir el esquema inicial de base de datos del MVP para catálogo, auth, checkout PayPal, historial de `orders` y admin básico.

## 2. Decisiones de diseño

- naming oficial: `orders`
- PostgreSQL como fuente de verdad
- UUID como PK en todas las tablas principales
- timestamps en DB como `TIMESTAMPTZ`
- precio y stock viven en `product_variants`, no en `products`
- `order_items` preserva snapshot histórico mínimo
- no se agregan tablas extra de carrito, refresh tokens o auditoría avanzada en el MVP

## 3. Tablas

### 3.1 `users`

| Columna | Tipo | Null | Default | Notas |
|---|---|---|---|---|
| `id` | `UUID` | no | — | PK |
| `email` | `VARCHAR(254)` | no | — | único, guardar normalizado en lower-case |
| `password_hash` | `VARCHAR(72)` | no | — | bcrypt |
| `is_admin` | `BOOLEAN` | no | `FALSE` | rol mínimo MVP |
| `created_at` | `TIMESTAMPTZ` | no | `NOW()` | |
| `updated_at` | `TIMESTAMPTZ` | no | `NOW()` | |

**Constraints**
- `PRIMARY KEY (id)`
- `UNIQUE (email)`

**Índices**
- `ux_users_email`

### 3.2 `products`

| Columna | Tipo | Null | Default | Notas |
|---|---|---|---|---|
| `id` | `UUID` | no | — | PK |
| `name` | `VARCHAR(160)` | no | — | nombre público |
| `slug` | `VARCHAR(180)` | no | — | único para URL |
| `description` | `TEXT` | no | — | |
| `category` | `VARCHAR(80)` | no | — | filtro MVP |
| `brand` | `VARCHAR(80)` | sí | — | opcional |
| `images` | `JSONB` | no | `'[]'::jsonb` | array de URLs |
| `active` | `BOOLEAN` | no | `TRUE` | visibilidad pública |
| `created_at` | `TIMESTAMPTZ` | no | `NOW()` | |
| `updated_at` | `TIMESTAMPTZ` | no | `NOW()` | |

**Constraints**
- `PRIMARY KEY (id)`
- `UNIQUE (slug)`
- `CHECK (jsonb_typeof(images) = 'array')`

**Índices**
- `ux_products_slug`
- `ix_products_category`
- `ix_products_active`

### 3.3 `product_variants`

| Columna | Tipo | Null | Default | Notas |
|---|---|---|---|---|
| `id` | `UUID` | no | — | PK |
| `product_id` | `UUID` | no | — | FK a `products.id` |
| `sku` | `VARCHAR(120)` | no | — | único |
| `color` | `VARCHAR(60)` | no | — | |
| `size` | `VARCHAR(30)` | no | — | |
| `price` | `NUMERIC(10,2)` | no | — | precio vendible |
| `stock` | `INTEGER` | no | `0` | stock actual |
| `image_url` | `TEXT` | sí | — | opcional |
| `created_at` | `TIMESTAMPTZ` | no | `NOW()` | |
| `updated_at` | `TIMESTAMPTZ` | no | `NOW()` | |

**Constraints**
- `PRIMARY KEY (id)`
- `FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE RESTRICT ON DELETE CASCADE`
- `UNIQUE (sku)`
- `UNIQUE (product_id, color, size)`
- `CHECK (price >= 0)`
- `CHECK (stock >= 0)`

**Índices**
- `ux_product_variants_sku`
- `ux_product_variants_product_color_size`
- `ix_product_variants_product_id`

### 3.4 `orders`

| Columna | Tipo | Null | Default | Notas |
|---|---|---|---|---|
| `id` | `UUID` | no | — | PK |
| `user_id` | `UUID` | no | — | FK a `users.id` |
| `status` | `VARCHAR(32)` | no | `'pending_payment'` | estado de negocio |
| `payment_provider` | `VARCHAR(24)` | no | `'paypal'` | MVP solo PayPal |
| `payment_status` | `VARCHAR(24)` | no | `'pending'` | estado técnico del pago |
| `currency` | `CHAR(3)` | no | `'USD'` | MVP |
| `subtotal` | `NUMERIC(10,2)` | no | — | recalculado backend |
| `total` | `NUMERIC(10,2)` | no | — | total final cobrado |
| `paypal_order_id` | `VARCHAR(64)` | sí | — | id de order PayPal |
| `paypal_capture_id` | `VARCHAR(64)` | sí | — | id de captura PayPal |
| `paid_at` | `TIMESTAMPTZ` | sí | — | fecha de captura confirmada |
| `created_at` | `TIMESTAMPTZ` | no | `NOW()` | |
| `updated_at` | `TIMESTAMPTZ` | no | `NOW()` | |

**Constraints**
- `PRIMARY KEY (id)`
- `FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE RESTRICT ON DELETE RESTRICT`
- `CHECK (status IN ('pending_payment', 'paid', 'payment_failed', 'cancelled', 'refunded'))`
- `CHECK (payment_provider IN ('paypal'))`
- `CHECK (payment_status IN ('pending', 'approved', 'captured', 'failed', 'refunded'))`
- `CHECK (subtotal >= 0)`
- `CHECK (total >= 0)`
- `UNIQUE (paypal_order_id)`
- `UNIQUE (paypal_capture_id)`

**Índices**
- `ix_orders_user_id_created_at`
- `ix_orders_status`
- `ix_orders_payment_status`

### 3.5 `order_items`

| Columna | Tipo | Null | Default | Notas |
|---|---|---|---|---|
| `id` | `UUID` | no | — | PK |
| `order_id` | `UUID` | no | — | FK a `orders.id` |
| `product_id` | `UUID` | no | — | FK a `products.id` |
| `variant_id` | `UUID` | no | — | FK a `product_variants.id` |
| `product_name` | `VARCHAR(160)` | no | — | snapshot |
| `variant_sku` | `VARCHAR(120)` | no | — | snapshot |
| `color` | `VARCHAR(60)` | no | — | snapshot |
| `size` | `VARCHAR(30)` | no | — | snapshot |
| `unit_price` | `NUMERIC(10,2)` | no | — | snapshot |
| `quantity` | `INTEGER` | no | — | |
| `line_total` | `NUMERIC(10,2)` | no | — | snapshot |
| `created_at` | `TIMESTAMPTZ` | no | `NOW()` | |

**Constraints**
- `PRIMARY KEY (id)`
- `FOREIGN KEY (order_id) REFERENCES orders(id) ON UPDATE RESTRICT ON DELETE CASCADE`
- `FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE RESTRICT ON DELETE RESTRICT`
- `FOREIGN KEY (variant_id) REFERENCES product_variants(id) ON UPDATE RESTRICT ON DELETE RESTRICT`
- `CHECK (quantity > 0)`
- `CHECK (unit_price >= 0)`
- `CHECK (line_total >= 0)`

**Índices**
- `ix_order_items_order_id`
- `ix_order_items_variant_id`

## 4. Relaciones

- un `user` tiene muchas `orders`
- un `product` tiene muchas `product_variants`
- una `order` tiene muchos `order_items`
- cada `order_item` referencia una variante real y además guarda snapshot histórico

## 5. Decisiones relevantes

### 5.1 Por qué precio en `product_variants`

Porque ropa requiere vender por combinación de talla/color y el stock también vive ahí.

### 5.2 Por qué `orders` guarda estado de negocio y pago

Porque para el MVP evita tablas extra y mantiene simple el historial del cliente y la operación admin.

### 5.3 Por qué `order_items` duplica datos

Para preservar historial aunque cambie el catálogo después.

## 6. Compatibilidad con el repositorio actual

El repositorio hoy ya tiene tablas y migraciones legadas (`purchase_orders`, `invoices`, productos sin variantes). La implementación base del MVP debe ser aditiva:

- agregar columnas faltantes a `users` y `products`
- crear `product_variants`
- crear `orders`
- crear `order_items`
- mantener coexistencia temporal con el esquema legado hasta migrar servicios y handlers

## 7. Contrato API alineado

Este esquema es la fuente de verdad para:

- `API-Contract-TiendaRopa.md`
- `Spec-Backend-TiendaRopa.md`
- `Spec-Frontend-TiendaRopa.md`

Toda nueva implementación debe usar `orders` como naming definitivo.
