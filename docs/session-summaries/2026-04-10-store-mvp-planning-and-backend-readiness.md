# 2026-04-10 — Store MVP planning and backend readiness

## Goal

Definir y ordenar la base funcional y técnica del MVP de la tienda de ropa para dejar el backend listo y poder iniciar la implementación del frontend con menos ambigüedad.

## Scope covered

- revisión y reescritura del PRD del frontend de tienda
- definición del admin básico dentro del MVP
- backlog inicial, roadmap y tickets
- contrato API y specs técnicas de frontend/backend
- decisión oficial de naming: `orders`
- definición técnica de JWT, PayPal, errores API y esquema DB
- implementación base y verificación local del backend
- reorganización de la documentación del MVP dentro de `docs/store-mvp/`

## Key decisions

1. **Naming oficial:** usar `orders` en lugar de `purchaseorders`.
2. **Modelo comercial:** producto con variantes por talla/color y stock por variante.
3. **Flujo de compra:** carrito como visitante y autenticación al momento de checkout.
4. **Operación mínima:** el admin entra al MVP para gestionar catálogo, variantes, stock y pedidos.
5. **Auth:** JWT para web con estrategia documentada y contrato backend/frontend definido.
6. **PayPal:** flujo backend-driven para create/capture; la captura correcta depende de la aprobación previa del comprador en PayPal.
7. **Errores API:** envelope uniforme `{"error": {...}}` en las superficies nuevas o adaptadas.

## Documentation created or updated

### Store MVP docs
- `docs/store-mvp/PRD-Frontend-TiendaRopa.md`
- `docs/store-mvp/API-Contract-TiendaRopa.md`
- `docs/store-mvp/DB-Schema-TiendaRopa.md`
- `docs/store-mvp/Roadmap-Ejecucion-TiendaRopa.md`
- `docs/store-mvp/Tickets-MVP-TiendaRopa.md`
- `docs/store-mvp/Spec-Frontend-TiendaRopa.md`
- `docs/store-mvp/Spec-Backend-TiendaRopa.md`

### General docs
- `README.md`

## Backend work completed

### New or adapted backend surfaces
- `POST /api/v1/public/register`
- `POST /api/v1/public/login`
- `GET /api/v1/private/me`
- `GET /api/v1/public/products`
- `GET /api/v1/public/products/:id`
- `POST /api/v1/private/orders/checkout/paypal`
- `POST /api/v1/private/orders/:id/paypal/capture`
- `GET /api/v1/private/orders`
- `GET /api/v1/private/orders/:id`

### Supporting implementation
- contrato store en `model/store_contract.go`
- cliente PayPal de orders en `infrastructure/paypal/orders.go`
- handler y servicio de `orders`
- repositorio PostgreSQL para `orders`
- migración base `sqlmigrations/20260410_1200_store_mvp_base.sql`

## Local verification performed

- migraciones aplicadas
- `go build ./...` OK
- `go test ./...` OK
- backend levanta localmente
- register/login OK
- `/private/me` OK
- catálogo y detalle con variantes OK
- listado y detalle de `orders` OK
- checkout PayPal OK con credenciales sandbox válidas
- capture responde correctamente con error controlado si se intenta sin aprobación previa del comprador

## Important findings

- el backend ya está listo para que el frontend inicie auth, catálogo, detalle, carrito y checkout PayPal
- la captura de PayPal no puede considerarse exitosa solo desde un smoke test backend, porque requiere aprobación real del comprador desde la UI/SDK
- tener múltiples `.md` del MVP en la raíz ensuciaba el repo y dificultaba ubicar la documentación relevante para frontend

## Reorganization performed

- los documentos del MVP se movieron a `docs/store-mvp/`
- este historial cronológico se empezó a guardar en `docs/session-summaries/`
- el resumen legado `sdd-architecture-hardening-summary.md` fue movido y renombrado con fecha para que quede ordenado cronológicamente

## Pending

- iniciar implementación base del frontend
- validar el flujo completo de aprobación + captura PayPal desde la UI real
- completar más adelante las superficies admin si se quiere cerrar todo el backoffice del MVP

## Recommended next step

Empezar el frontend con este orden:

1. auth y sesión
2. catálogo
3. detalle de producto con variantes
4. carrito
5. checkout PayPal
6. captura post-aprobación
7. historial de `orders`

## Relevant files

- `cmd/server.go`
- `cmd/routes/routes_orders.go`
- `domain/services/order.go`
- `infrastructure/handlers/order.go`
- `infrastructure/paypal/orders.go`
- `infrastructure/postgres/order.go`
- `model/store_contract.go`
- `sqlmigrations/20260410_1200_store_mvp_base.sql`
