# Roadmap de Ejecución - Tienda de Ropa MVP

## 1. Enfoque

Roadmap incremental para un equipo pequeño. La prioridad sigue siendo habilitar primero el flujo de venta completo, pero ahora con cuatro decisiones ya cerradas: JWT, PayPal, formato de error y esquema DB inicial.

## 2. Readiness actual

### Cerrado
- naming oficial `orders`
- estrategia JWT del MVP
- flujo técnico recomendado de PayPal
- formato uniforme de errores API
- esquema DB inicial del MVP

### Iniciado
- base de implementación backend alineada a `orders`
- tipos Go del contrato nuevo
- migración SQL inicial de soporte al MVP

### Pendiente
- handlers/services/repositories completos de catálogo con variantes
- checkout PayPal end-to-end
- historial completo de `orders`
- admin básico operativo

## 3. Roadmap por etapas

### Etapa 0 - Cierre técnico y base de arranque

**Objetivo**: eliminar ambigüedades antes del desarrollo funcional.

**Salida**
- contrato API actualizado
- specs frontend/backend actualizadas
- esquema DB inicial definido
- implementación base mínima arrancada

### Etapa 1 - Catálogo y variantes

**Objetivo**: exponer catálogo real con `products` + `product_variants`.

**Backend**
- consolidar modelo `products`
- incorporar `product_variants`
- `GET /api/v1/public/products`
- `GET /api/v1/public/products/:id`

**Frontend**
- layout público
- catálogo navegable
- detalle de producto con variantes

### Etapa 2 - Auth y sesión

**Objetivo**: permitir sesión funcional para checkout y área privada.

**Backend**
- `POST /api/v1/public/register`
- `POST /api/v1/public/login`
- `GET /api/v1/private/me`
- middleware auth/admin alineado al contrato

**Frontend**
- login
- register
- restauración de sesión desde `sessionStorage`
- guards de navegación

### Etapa 3 - Checkout PayPal y creación de `orders`

**Objetivo**: cerrar compra real de punta a punta.

**Backend**
- `POST /api/v1/private/orders/checkout/paypal`
- `POST /api/v1/private/orders/:id/paypal/capture`
- persistencia de `orders` y `order_items`
- validación de stock y confirmación final

**Frontend**
- checkout autenticado
- integración PayPal JS SDK
- limpieza de carrito al éxito

### Etapa 4 - Historial del cliente

**Objetivo**: permitir seguimiento post-compra.

**Backend**
- `GET /api/v1/private/orders`
- `GET /api/v1/private/orders/:id`

**Frontend**
- listado de `orders`
- detalle de `order`

### Etapa 5 - Admin básico

**Objetivo**: operación mínima real.

**Backend**
- CRUD básico de productos
- listado y detalle admin de `orders`
- cambio de estado validado

**Frontend**
- admin products
- admin orders

## 4. Riesgos principales

1. **Coexistencia con esquema legado**
   - hoy existen `purchase_orders` e `invoices`
   - se requiere transición controlada hacia `orders`

2. **Cambio a variantes reales**
   - el repositorio actual todavía modela precio a nivel producto
   - el MVP final necesita precio y stock por variante

3. **Error envelope aún no cableado en todo el backend**
   - ya quedó definido
   - falta adopción progresiva en handlers/responses existentes

## 5. Orden recomendado inmediato

1. implementar repositorio/migraciones de `product_variants`, `orders` y `order_items`
2. ajustar auth al contrato final (`sub`, expiración y `me`)
3. construir checkout backend-driven de PayPal
4. cerrar historial de `orders`
5. completar admin básico
