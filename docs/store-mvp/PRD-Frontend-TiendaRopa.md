# PRD - Tienda de Ropa (Frontend E-commerce + Admin MVP)

## 1. Visión del Producto

Tienda online de venta directa de ropa para público general, enfocada en una experiencia de compra simple, rápida y clara. El objetivo del MVP es permitir que cualquier visitante pueda explorar el catálogo, buscar, filtrar, seleccionar variantes de producto, agregar al carrito y completar la compra con la menor fricción posible.

Además del frontend público, el MVP debe incluir un módulo administrador básico para gestionar productos, variantes, stock y pedidos, de modo que la tienda pueda operar desde el primer día.

---

## 2. Objetivo del MVP

Construir una tienda virtual funcional con foco en:

- navegación simple
- selección de producto clara
- proceso de compra corto
- operación interna básica para administrar catálogo y pedidos

### Principios del MVP

1. **Simpleza antes que cantidad de funcionalidades**
2. **Comprar debe ser fácil desde móvil y escritorio**
3. **El administrador debe poder operar sin depender de desarrollos posteriores**
4. **Solo se incluyen funcionalidades esenciales para vender y gestionar pedidos**

---

## 3. Alcance MVP

| Módulo | Funcionalidades |
|--------|-----------------|
| **Home / Catálogo** | Listado de productos destacados o recientes, acceso rápido al catálogo |
| **Catálogo** | Listado de productos, búsqueda, filtros básicos, paginación o carga incremental simple |
| **Detalle de producto** | Galería, descripción, precio, variantes (talla/color), stock visible, agregar al carrito |
| **Carrito** | Agregar productos, editar cantidad, remover items, persistencia local |
| **Checkout** | Revisión de compra, autenticación requerida al finalizar, pago con PayPal |
| **Autenticación** | Registro, login, sesión persistente |
| **Usuario** | Historial de pedidos, detalle de pedido |
| **Administrador - Productos** | Crear, editar, activar/desactivar productos, gestionar variantes, stock e imágenes |
| **Administrador - Pedidos** | Listar pedidos, ver detalle, actualizar estado de pedido |

### Exclusiones para versiones posteriores

- wishlist
- reseñas
- newsletter
- cupones
- múltiples pasarelas de pago
- panel analítico avanzado
- gestión avanzada de promociones
- múltiples roles administrativos complejos

---

## 4. User Stories

### Catálogo y descubrimiento
1. Como visitante quiero ver el catálogo para explorar productos disponibles.
2. Como visitante quiero buscar productos por nombre para encontrar algo específico.
3. Como visitante quiero filtrar productos por categoría, precio, color y talla para reducir opciones.
4. Como visitante quiero ver el detalle de un producto con imágenes, precio, variantes y stock para decidir si comprarlo.

### Carrito
5. Como visitante quiero agregar productos al carrito para preparar mi compra sin iniciar sesión de inmediato.
6. Como visitante o cliente quiero editar la cantidad de productos del carrito.
7. Como visitante o cliente quiero remover productos del carrito.
8. Como visitante o cliente quiero ver subtotal y total estimado de mi compra.

### Checkout
9. Como cliente quiero iniciar sesión o registrarme al momento de comprar para completar el checkout.
10. Como cliente quiero pagar con PayPal para completar mi compra.
11. Como cliente quiero recibir confirmación clara cuando mi orden sea creada exitosamente.

### Autenticación y cuenta
12. Como visitante quiero registrarme fácilmente para poder comprar.
13. Como usuario registrado quiero iniciar sesión para acceder a mi cuenta y mis pedidos.
14. Como cliente quiero ver mi historial de pedidos para consultar mis compras previas.
15. Como cliente quiero ver el detalle de un pedido para revisar estado, fecha, total y productos.

### Administración
16. Como administrador quiero crear productos para publicar nuevo catálogo.
17. Como administrador quiero crear y editar variantes por talla y color para vender ropa correctamente.
18. Como administrador quiero actualizar stock por variante para evitar ventas de productos agotados.
19. Como administrador quiero activar o desactivar productos sin eliminarlos.
20. Como administrador quiero ver pedidos y actualizar su estado para dar soporte a la operación de la tienda.

---

## 5. Stack Tecnológico

- **Frontend**: React 19 + Vite + TypeScript
- **Routing**: React Router v7
- **HTTP Client**: Fetch API
- **Backend**: Go + Echo + PostgreSQL (existente)
- **Pagos**: PayPal SDK
- **Persistencia local**: localStorage para carrito y sesión de UI no sensible

### Artefactos técnicos derivados cerrados

Las decisiones técnicas del MVP quedaron cerradas y deben tomarse como fuente operativa en los siguientes artefactos:

- `API-Contract-TiendaRopa.md` → contrato final, errores API y endpoints
- `Spec-Frontend-TiendaRopa.md` → estrategia JWT frontend y flujo PayPal UI/API
- `Spec-Backend-TiendaRopa.md` → JWT backend, flujo PayPal backend-driven y reglas de `orders`
- `DB-Schema-TiendaRopa.md` → esquema inicial de `users`, `products`, `product_variants`, `orders` y `order_items`

En esos artefactos quedó cerrada la definición de JWT, PayPal, formato de error y esquema DB del MVP.

---

## 6. Requisitos No Funcionales

| Aspecto | Requisito |
|---------|-----------|
| **Rendimiento** | Carga inicial percibida < 3s en condiciones normales |
| **Responsive** | Mobile-first, compatible con móvil, tablet y escritorio |
| **Accesibilidad** | Navegación por teclado, labels correctos, contraste suficiente, estados visibles |
| **Seguridad** | JWT válido, validación de inputs, protección XSS, control de acceso por rol |
| **UX** | Feedback claro en loading, errores, éxito, estados vacíos y formularios |
| **Operatividad** | El admin debe poder crear productos y actualizar pedidos sin soporte técnico |
| **Consistencia** | Stock y variantes deben validarse tanto en catálogo como en checkout |

---

## 7. Navegación y Arquitectura de Información

La navegación debe ser clara y corta.

### Header público
- logo
- acceso a catálogo
- buscador
- carrito
- login / perfil

### Reglas UX
- el usuario debe poder llegar a un producto en pocos pasos
- el detalle debe mostrar primero información clave: imágenes, nombre, precio, variantes, stock, botón de compra
- el checkout no debe tener pasos innecesarios
- en móvil, los filtros deben abrirse en drawer o panel colapsable

---

## 8. Estructura de Rutas

### Público / Cliente

```txt
/                       → Home
/products               → Catálogo con búsqueda y filtros
/products/:id           → Detalle de producto
/cart                   → Carrito de compras
/checkout               → Checkout (requiere auth al confirmar compra)
/login                  → Login
/register               → Registro
/profile                → Perfil de usuario
/profile/orders         → Historial de pedidos
/profile/orders/:id     → Detalle de pedido
```

### Administrador

```txt
/admin                  → Inicio admin
/admin/products         → Listado de productos
/admin/products/new     → Crear producto
/admin/products/:id     → Editar producto
/admin/orders           → Listado de pedidos
/admin/orders/:id       → Detalle y gestión de pedido
```

---

## 9. Modelos de Datos (Frontend)

### Product
```typescript
interface Product {
  id: string;
  name: string;
  slug: string;
  description: string;
  category: string;
  brand?: string;
  images: string[];
  active: boolean;
  created_at: number;
  updated_at: number;
  variants: ProductVariant[];
}
```

### ProductVariant
```typescript
interface ProductVariant {
  id: string;
  product_id: string;
  color: string;
  size: string;
  sku: string;
  price: number;
  stock: number;
  image?: string;
}
```

### User
```typescript
interface User {
  id: string;
  email: string;
  is_admin: boolean;
}
```

### CartItem
```typescript
interface CartItem {
  product_id: string;
  product_name: string;
  product_image: string;
  variant_id: string;
  color: string;
  size: string;
  unit_price: number;
  quantity: number;
  available_stock: number;
}
```

### Order
```typescript
interface Order {
  id: string;
  user_id: string;
  total: number;
  status: 'pending' | 'paid' | 'failed' | 'cancelled' | 'shipped' | 'delivered';
  created_at: number;
  items: OrderItem[];
}
```

### OrderItem
```typescript
interface OrderItem {
  product_id: string;
  product_name: string;
  variant_id: string;
  color: string;
  size: string;
  unit_price: number;
  quantity: number;
}
```

---

## 10. Integración con Backend

### Público

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/api/v1/public/products` | GET | Listado de productos con búsqueda y filtros |
| `/api/v1/public/products/:id` | GET | Detalle de producto |
| `/api/v1/public/register` | POST | Registro de usuario |
| `/api/v1/public/login` | POST | Login de usuario |

### Privado - Cliente

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/api/v1/private/me` | GET | Datos del usuario autenticado |
| `/api/v1/private/orders` | POST | Crear orden de compra |
| `/api/v1/private/orders` | GET | Listar pedidos del usuario |
| `/api/v1/private/orders/:id` | GET | Ver detalle de pedido |

### Privado - Administrador

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `/api/v1/admin/products` | GET | Listar productos para administración |
| `/api/v1/admin/products` | POST | Crear producto |
| `/api/v1/admin/products/:id` | PUT | Editar producto |
| `/api/v1/admin/products/:id/status` | PATCH | Activar o desactivar producto |
| `/api/v1/admin/orders` | GET | Listar pedidos |
| `/api/v1/admin/orders/:id` | GET | Ver detalle de pedido |
| `/api/v1/admin/orders/:id/status` | PATCH | Actualizar estado del pedido |

---

## 11. Componentes UI

### Tienda pública

| Componente | Descripción |
|------------|-------------|
| `Navbar` | Logo, navegación principal, buscador, carrito, login/perfil |
| `ProductCard` | Imagen, nombre, precio base o desde, estado de stock |
| `ProductList` | Grid de productos |
| `ProductDetail` | Galería, descripción, variantes, stock, CTA agregar |
| `FilterSidebar` | Filtros por categoría, precio, color, talla |
| `SearchBar` | Búsqueda por nombre |
| `VariantSelector` | Selección de color y talla |
| `CartItem` | Producto agregado con variante, cantidad y subtotal |
| `CartSummary` | Resumen del carrito, subtotal y CTA checkout |
| `LoginForm` | Email y password |
| `RegisterForm` | Email, password y confirmación |
| `OrderHistory` | Lista de pedidos del usuario |
| `OrderDetail` | Detalle de pedido |
| `ProtectedRoute` | Protege rutas de usuario autenticado |
| `AdminRoute` | Protege rutas administrativas |

### Administración

| Componente | Descripción |
|------------|-------------|
| `AdminLayout` | Layout principal del módulo admin |
| `AdminProductList` | Lista de productos con acciones básicas |
| `AdminProductForm` | Formulario de producto |
| `VariantFormList` | Gestión de variantes de producto |
| `StockEditor` | Edición rápida de stock |
| `AdminOrderList` | Listado de pedidos |
| `AdminOrderDetail` | Vista detalle y cambio de estado |

---

## 12. Flujo de Usuario

### Flujo de compra

```txt
Visitante
    │
    ├─→ Home / Catálogo
    ├─→ Buscar / Filtrar productos
    ├─→ Ver detalle de producto
    ├─→ Seleccionar talla y color
    ├─→ Agregar al carrito
    ├─→ Revisar carrito
    └─→ Iniciar checkout
              │
              ├─→ Si no está autenticado → Login / Registro
              │                                │
              │                                └─→ Volver a checkout
              │
              └─→ Confirmar compra → PayPal → Orden creada → Perfil / detalle pedido
```

### Flujo administrativo

```txt
Administrador
    │
    ├─→ Login
    ├─→ /admin/products
    │      ├─→ Crear producto
    │      ├─→ Editar producto
    │      ├─→ Crear variantes
    │      └─→ Actualizar stock / activar / desactivar
    │
    └─→ /admin/orders
           ├─→ Ver pedidos
           ├─→ Ver detalle
           └─→ Cambiar estado del pedido
```

---

## 13. Reglas Funcionales Clave

1. Un producto de ropa puede tener múltiples variantes por talla y color.
2. El usuario no puede agregar al carrito un producto sin seleccionar variante cuando aplique.
3. No se puede agregar al carrito una cantidad mayor al stock disponible.
4. El carrito persiste en localStorage.
5. El usuario puede armar su carrito como visitante.
6. El login o registro se exige al momento de completar el checkout.
7. El stock debe validarse nuevamente antes de crear la orden.
8. Un administrador puede crear productos y variantes sin depender de otras herramientas externas.
9. Los productos pueden activarse o desactivarse sin eliminarse.
10. Los pedidos deben tener estados consistentes y visibles para cliente y administrador.

---

## 14. Criterios de Aceptación

### Catálogo
- [ ] Se muestran los productos activos del backend.
- [ ] La búsqueda por nombre retorna resultados relevantes.
- [ ] Los filtros por categoría, precio, color y talla reducen correctamente los resultados.
- [ ] El usuario puede limpiar filtros fácilmente.
- [ ] El catálogo funciona correctamente en móvil y escritorio.

### Detalle de producto
- [ ] Se muestran imágenes, nombre, descripción, precio y variantes del producto.
- [ ] El usuario puede seleccionar talla y color antes de agregar al carrito.
- [ ] Si no hay stock disponible, se informa claramente.
- [ ] No se puede agregar al carrito una variante agotada.

### Carrito
- [ ] Se pueden agregar múltiples productos distintos con variantes distintas.
- [ ] Se puede cambiar la cantidad de cada item.
- [ ] Se puede remover un item del carrito.
- [ ] Se muestra el total calculado correctamente.
- [ ] El carrito persiste en localStorage.
- [ ] El carrito conserva talla y color seleccionados.

### Checkout
- [ ] Un visitante puede iniciar el flujo de compra desde el carrito.
- [ ] Si el usuario no está autenticado, se redirige a login o registro antes de pagar.
- [ ] Luego de autenticarse, el usuario vuelve al checkout.
- [ ] El flujo de PayPal se ejecuta correctamente.
- [ ] Al completar el pago, se crea la orden en el backend.
- [ ] Si el stock cambió antes del pago, se informa al usuario y se evita una compra inconsistente.

### Autenticación
- [ ] El registro crea usuario en el backend.
- [ ] El login retorna un JWT válido.
- [ ] El token se almacena de forma segura según la estrategia definida por frontend/backend.
- [ ] Las rutas privadas y administrativas están protegidas correctamente.

### Usuario
- [ ] El usuario puede ver su historial de pedidos.
- [ ] Cada pedido muestra estado, total y fecha.
- [ ] El usuario puede ver el detalle de un pedido.

### Administración - Productos
- [ ] El administrador puede crear un producto.
- [ ] El administrador puede editar un producto existente.
- [ ] El administrador puede crear y editar variantes.
- [ ] El administrador puede actualizar stock por variante.
- [ ] El administrador puede activar o desactivar productos.

### Administración - Pedidos
- [ ] El administrador puede listar pedidos.
- [ ] El administrador puede ver el detalle de un pedido.
- [ ] El administrador puede cambiar el estado del pedido.

---

## 15. Prioridad de Implementación

### Fase 1 - Venta básica
1. Home y catálogo
2. Detalle de producto
3. Selección de variantes
4. Carrito
5. Login y registro
6. Checkout con PayPal
7. Historial y detalle de pedidos

### Fase 2 - Operación básica
8. Admin de productos
9. Admin de variantes y stock
10. Admin de pedidos

### Fase 3 - Mejoras futuras
11. Refinamientos UX
12. Métricas y reportes
13. Funcionalidades comerciales adicionales

---

## 16. Resumen Ejecutivo

Este MVP debe enfocarse en vender y operar con simplicidad.

La prioridad no es tener muchas funcionalidades, sino cubrir correctamente el recorrido esencial:

1. descubrir productos
2. seleccionar talla/color
3. agregar al carrito
4. autenticarse solo cuando haga falta
5. pagar sin fricción
6. gestionar productos y pedidos desde un admin básico

Si estas capacidades funcionan bien, la tienda tendrá una base sólida para crecer sin rehacer la arquitectura principal.

---

## 17. Backlog Inicial del MVP (Historias y Tareas)

El backlog se organiza por prioridad funcional. El objetivo es construir primero el flujo de venta y luego la operación administrativa mínima.

### Epic 1 - Base del frontend público

#### Historia 1.1 - Layout base y navegación pública
**Como** visitante  
**Quiero** navegar entre home, catálogo, carrito, login y perfil  
**Para** moverme por la tienda de forma simple.

**Tareas técnicas**
- crear layout público principal
- implementar `Navbar`
- configurar rutas públicas
- definir páginas base: home, catálogo, carrito, login, register, profile

#### Historia 1.2 - Home con acceso al catálogo
**Como** visitante  
**Quiero** ver una portada simple con acceso rápido a productos  
**Para** comenzar a explorar sin fricción.

**Tareas técnicas**
- crear página home
- mostrar productos destacados o recientes
- agregar CTA hacia catálogo

---

### Epic 2 - Catálogo y descubrimiento de productos

#### Historia 2.1 - Listado de productos
**Como** visitante  
**Quiero** ver el catálogo de productos activos  
**Para** explorar opciones disponibles.

**Tareas técnicas**
- consumir `GET /api/v1/public/products`
- renderizar grid de productos
- mostrar loading, error y estado vacío

#### Historia 2.2 - Búsqueda de productos
**Como** visitante  
**Quiero** buscar por nombre  
**Para** encontrar productos más rápido.

**Tareas técnicas**
- implementar `SearchBar`
- sincronizar búsqueda con query params
- reflejar búsqueda en la consulta al backend

#### Historia 2.3 - Filtros de catálogo
**Como** visitante  
**Quiero** filtrar por categoría, precio, color y talla  
**Para** reducir resultados.

**Tareas técnicas**
- implementar `FilterSidebar`
- persistir filtros en URL
- soportar filtros en móvil y escritorio
- agregar acción limpiar filtros

#### Historia 2.4 - Detalle de producto
**Como** visitante  
**Quiero** ver el detalle del producto  
**Para** tomar una decisión de compra.

**Tareas técnicas**
- consumir `GET /api/v1/public/products/:id`
- mostrar galería, descripción, variantes y stock
- implementar selección de variante

---

### Epic 3 - Carrito y selección de compra

#### Historia 3.1 - Agregar al carrito
**Como** visitante  
**Quiero** agregar una variante al carrito  
**Para** preparar mi compra.

**Tareas técnicas**
- validar variante seleccionada
- agregar item con variante, cantidad y precio
- guardar carrito en localStorage

#### Historia 3.2 - Gestión del carrito
**Como** visitante o cliente  
**Quiero** editar cantidades y remover productos  
**Para** ajustar mi compra antes de pagar.

**Tareas técnicas**
- crear vista `CartItem` y `CartSummary`
- permitir sumar, restar y remover
- recalcular subtotales y total
- bloquear cantidades mayores al stock conocido

---

### Epic 4 - Autenticación y sesión

#### Historia 4.1 - Registro
**Como** visitante  
**Quiero** registrarme  
**Para** poder comprar y consultar mis pedidos.

**Tareas técnicas**
- implementar formulario de registro
- consumir `POST /api/v1/public/register`
- mostrar feedback de error y éxito

#### Historia 4.2 - Login
**Como** usuario  
**Quiero** iniciar sesión  
**Para** acceder a checkout y perfil.

**Tareas técnicas**
- implementar formulario de login
- consumir `POST /api/v1/public/login`
- almacenar sesión según estrategia acordada
- consultar `GET /api/v1/private/me`

#### Historia 4.3 - Rutas protegidas
**Como** sistema  
**Quiero** proteger checkout, perfil y admin  
**Para** evitar accesos no autorizados.

**Tareas técnicas**
- implementar `ProtectedRoute`
- implementar `AdminRoute`
- soportar redirect al destino original después de login

---

### Epic 5 - Checkout y pedidos

#### Historia 5.1 - Checkout autenticado
**Como** cliente  
**Quiero** revisar mi compra antes de pagar  
**Para** confirmar que todo esté correcto.

**Tareas técnicas**
- crear página checkout
- redirigir a login si no hay sesión
- restaurar navegación al checkout luego del login

#### Historia 5.2 - Pago con PayPal
**Como** cliente  
**Quiero** pagar con PayPal  
**Para** finalizar mi compra.

**Tareas técnicas**
- integrar PayPal SDK
- manejar éxito, cancelación y error
- enviar payload final al backend

#### Historia 5.3 - Creación y consulta de pedidos
**Como** cliente  
**Quiero** que se cree mi pedido y luego poder consultarlo  
**Para** hacer seguimiento de la compra.

**Tareas técnicas**
- consumir `POST /api/v1/private/orders`
- consumir `GET /api/v1/private/orders`
- consumir `GET /api/v1/private/orders/:id`
- limpiar carrito luego de compra exitosa

---

### Epic 6 - Administración de catálogo

#### Historia 6.1 - Listado de productos admin
**Como** administrador  
**Quiero** ver el listado de productos  
**Para** gestionarlos rápidamente.

**Tareas técnicas**
- crear `AdminProductList`
- consumir `GET /api/v1/admin/products`
- mostrar estado activo/inactivo

#### Historia 6.2 - Crear y editar producto
**Como** administrador  
**Quiero** crear y editar productos  
**Para** mantener el catálogo actualizado.

**Tareas técnicas**
- crear `AdminProductForm`
- consumir `POST /api/v1/admin/products`
- consumir `PUT /api/v1/admin/products/:id`

#### Historia 6.3 - Gestionar variantes y stock
**Como** administrador  
**Quiero** administrar talla, color, SKU, precio y stock  
**Para** operar correctamente la tienda.

**Tareas técnicas**
- crear UI para variantes
- soportar altas, edición y baja lógica de variantes
- mostrar y editar stock por variante

#### Historia 6.4 - Activar y desactivar productos
**Como** administrador  
**Quiero** activar o desactivar productos  
**Para** controlar visibilidad sin borrar información.

**Tareas técnicas**
- consumir `PATCH /api/v1/admin/products/:id/status`
- actualizar estado en UI

---

### Epic 7 - Administración de pedidos

#### Historia 7.1 - Listado de pedidos
**Como** administrador  
**Quiero** ver los pedidos realizados  
**Para** gestionar la operación diaria.

**Tareas técnicas**
- crear `AdminOrderList`
- consumir `GET /api/v1/admin/orders`
- mostrar estado, fecha, total y cliente

#### Historia 7.2 - Detalle y actualización de pedido
**Como** administrador  
**Quiero** ver el detalle y cambiar el estado del pedido  
**Para** dar seguimiento al proceso de compra.

**Tareas técnicas**
- crear `AdminOrderDetail`
- consumir `GET /api/v1/admin/orders/:id`
- consumir `PATCH /api/v1/admin/orders/:id/status`

---

### Orden sugerido de implementación en backlog

1. Layout y routing base
2. Catálogo
3. Detalle de producto y variantes
4. Carrito
5. Registro y login
6. Checkout + PayPal
7. Pedidos del usuario
8. Admin de productos
9. Admin de variantes y stock
10. Admin de pedidos

---

## 18. Contrato API / Backend Propuesto

Este contrato busca dar soporte completo al frontend del MVP. Se prioriza consistencia, simplicidad y separación clara entre público, cliente autenticado y administrador.

### 18.1 Convenciones generales

- Base path: `/api/v1`
- Formato de respuesta: JSON
- Autenticación: JWT Bearer Token
- Fechas: ISO 8601 o timestamp consistente en toda la API
- Errores: estructura uniforme

### Respuesta exitosa sugerida

```json
{
  "data": {}
}
```

### Respuesta de error sugerida

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Los datos enviados no son válidos",
    "details": {}
  }
}
```

---

### 18.2 Catálogo público

#### GET `/api/v1/public/products`

Obtiene productos activos con filtros.

**Query params sugeridos**
- `search`: texto libre
- `category`: string
- `color`: string
- `size`: string
- `min_price`: number
- `max_price`: number
- `page`: number
- `limit`: number
- `sort`: `price_asc | price_desc | newest`

**Response**
```json
{
  "data": {
    "items": [
      {
        "id": "prod_1",
        "name": "Chaqueta Denim Azul",
        "slug": "chaqueta-denim-azul",
        "category": "chaquetas",
        "images": ["https://.../1.jpg"],
        "active": true,
        "price_from": 49.99,
        "available_colors": ["azul", "negro"],
        "available_sizes": ["S", "M", "L"]
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 12,
      "total": 120,
      "total_pages": 10
    }
  }
}
```

#### GET `/api/v1/public/products/:id`

Obtiene detalle completo de un producto.

**Response**
```json
{
  "data": {
    "id": "prod_1",
    "name": "Chaqueta Denim Azul",
    "slug": "chaqueta-denim-azul",
    "description": "Chaqueta clásica de denim.",
    "category": "chaquetas",
    "brand": "Marca X",
    "images": ["https://.../1.jpg", "https://.../2.jpg"],
    "active": true,
    "variants": [
      {
        "id": "var_1",
        "color": "azul",
        "size": "M",
        "sku": "CHA-DEN-AZU-M",
        "price": 49.99,
        "stock": 8,
        "image": "https://.../1.jpg"
      }
    ]
  }
}
```

---

### 18.3 Autenticación

#### POST `/api/v1/public/register`

**Request**
```json
{
  "email": "cliente@email.com",
  "password": "secret123",
  "confirm_password": "secret123"
}
```

**Response**
```json
{
  "data": {
    "user": {
      "id": "usr_1",
      "email": "cliente@email.com",
      "is_admin": false
    }
  }
}
```

#### POST `/api/v1/public/login`

**Request**
```json
{
  "email": "cliente@email.com",
  "password": "secret123"
}
```

**Response**
```json
{
  "data": {
    "token": "jwt-token",
    "user": {
      "id": "usr_1",
      "email": "cliente@email.com",
      "is_admin": false
    }
  }
}
```

#### GET `/api/v1/private/me`

**Response**
```json
{
  "data": {
    "id": "usr_1",
    "email": "cliente@email.com",
    "is_admin": false
  }
}
```

---

### 18.4 Pedidos del cliente

#### POST `/api/v1/private/orders`

Este endpoint debe validar stock nuevamente y registrar el resultado del pago.

**Request**
```json
{
  "payment_method": "paypal",
  "payment_reference": "PAYPAL-ORDER-ID",
  "items": [
    {
      "variant_id": "var_1",
      "quantity": 2
    }
  ]
}
```

**Response**
```json
{
  "data": {
    "id": "ord_1",
    "user_id": "usr_1",
    "status": "paid",
    "total": 99.98,
    "created_at": "2026-04-10T10:00:00Z",
    "items": [
      {
        "product_id": "prod_1",
        "product_name": "Chaqueta Denim Azul",
        "variant_id": "var_1",
        "color": "azul",
        "size": "M",
        "unit_price": 49.99,
        "quantity": 2
      }
    ]
  }
}
```

**Errores esperados**
- `VALIDATION_ERROR`
- `OUT_OF_STOCK`
- `PAYMENT_INVALID`
- `UNAUTHORIZED`

#### GET `/api/v1/private/orders`

Lista pedidos del usuario autenticado.

#### GET `/api/v1/private/orders/:id`

Obtiene detalle de un pedido del usuario autenticado.

---

### 18.5 Administración de productos

#### GET `/api/v1/admin/products`

Lista productos para administración.

**Query params sugeridos**
- `search`
- `status`: `active | inactive | all`
- `page`
- `limit`

#### POST `/api/v1/admin/products`

**Request**
```json
{
  "name": "Chaqueta Denim Azul",
  "slug": "chaqueta-denim-azul",
  "description": "Chaqueta clásica de denim.",
  "category": "chaquetas",
  "brand": "Marca X",
  "images": ["https://.../1.jpg"],
  "active": true,
  "variants": [
    {
      "color": "azul",
      "size": "M",
      "sku": "CHA-DEN-AZU-M",
      "price": 49.99,
      "stock": 8
    }
  ]
}
```

#### PUT `/api/v1/admin/products/:id`

Actualiza producto y variantes en una sola operación o en estructura compatible con el backend existente.

#### PATCH `/api/v1/admin/products/:id/status`

**Request**
```json
{
  "active": false
}
```

---

### 18.6 Administración de pedidos

#### GET `/api/v1/admin/orders`

Lista pedidos para operación.

**Query params sugeridos**
- `status`
- `search`
- `page`
- `limit`

#### GET `/api/v1/admin/orders/:id`

Obtiene detalle completo del pedido.

#### PATCH `/api/v1/admin/orders/:id/status`

**Request**
```json
{
  "status": "shipped"
}
```

**Estados permitidos**
- `pending`
- `paid`
- `failed`
- `cancelled`
- `shipped`
- `delivered`

---

### 18.7 Reglas de validación backend recomendadas

1. No permitir crear una orden sin items.
2. No permitir cantidades menores a 1.
3. Validar existencia y disponibilidad de cada variante.
4. Recalcular total en backend; no confiar solo en frontend.
5. No exponer productos inactivos en catálogo público.
6. Restringir endpoints admin a usuarios con `is_admin = true`.
7. Mantener consistencia de estados de pedido con transiciones válidas.

---

## 19. Próximo Paso Recomendado

Con este PRD actualizado, el siguiente paso lógico es convertir el backlog en:

1. **roadmap técnico por sprint**, o
2. **tareas de implementación frontend/backend**, o
3. **spec técnica más detallada para iniciar desarrollo**.

### Entregables derivados del PRD

Para separar producto de ejecución técnica, este PRD se complementa con los siguientes documentos derivados:

- `Tickets-MVP-TiendaRopa.md`
- `API-Contract-TiendaRopa.md`
- `Roadmap-Ejecucion-TiendaRopa.md`
- `Spec-Frontend-TiendaRopa.md`
- `Spec-Backend-TiendaRopa.md`

Regla de consistencia: el naming oficial del recurso es **`orders`** en frontend, backend y documentación.

---

## 20. Roadmap Sugerido por Sprints

Se propone una planificación incremental orientada a reducir riesgo funcional primero. El foco inicial debe ser habilitar el flujo de venta completo y luego la operación administrativa.

### Suposición de trabajo

- sprint de 1 a 2 semanas
- un equipo pequeño
- prioridad en funcionalidad estable antes que refinamiento visual avanzado

### Sprint 1 - Base pública y catálogo

**Objetivo**: dejar navegable la tienda y visible el catálogo.

**Alcance frontend**
- layout público principal
- navbar con navegación y acceso a carrito/login
- home simple con acceso al catálogo
- página de catálogo
- listado de productos
- búsqueda por nombre
- filtros básicos
- manejo de loading, error y estados vacíos

**Alcance backend**
- `GET /api/v1/public/products`
- `GET /api/v1/public/products/:id`
- soporte de filtros, búsqueda y paginación
- exponer solo productos activos

**Resultado esperado**
- el usuario puede navegar la tienda y descubrir productos

### Sprint 2 - Detalle de producto y carrito

**Objetivo**: permitir seleccionar productos y construir la compra.

**Alcance frontend**
- página de detalle de producto
- selector de variantes
- validación visual de stock
- agregar al carrito
- vista de carrito
- editar cantidad y remover items
- persistencia en localStorage

**Alcance backend**
- detalle de producto con variantes completas
- modelo estable de variantes con talla, color, precio y stock

**Resultado esperado**
- el usuario puede seleccionar talla/color y armar carrito correctamente

### Sprint 3 - Autenticación y checkout

**Objetivo**: cerrar el flujo de compra.

**Alcance frontend**
- registro
- login
- sesión persistente
- rutas protegidas
- redirección a checkout después de login
- página checkout
- integración PayPal

**Alcance backend**
- `POST /api/v1/public/register`
- `POST /api/v1/public/login`
- `GET /api/v1/private/me`
- `POST /api/v1/private/orders`
- validación de stock al crear orden
- registro de referencia de pago

**Resultado esperado**
- el usuario puede pasar de carrito a compra completa

### Sprint 4 - Pedidos del usuario

**Objetivo**: dar visibilidad post-compra al cliente.

**Alcance frontend**
- perfil de usuario
- historial de pedidos
- detalle de pedido

**Alcance backend**
- `GET /api/v1/private/orders`
- `GET /api/v1/private/orders/:id`

**Resultado esperado**
- el cliente puede consultar sus compras y estados

### Sprint 5 - Admin de productos

**Objetivo**: habilitar operación mínima del catálogo.

**Alcance frontend**
- layout admin
- listado de productos admin
- formulario de creación y edición
- gestión de variantes
- activación y desactivación de productos

**Alcance backend**
- `GET /api/v1/admin/products`
- `POST /api/v1/admin/products`
- `PUT /api/v1/admin/products/:id`
- `PATCH /api/v1/admin/products/:id/status`

**Resultado esperado**
- el administrador puede operar el catálogo sin soporte técnico

### Sprint 6 - Admin de pedidos y hardening MVP

**Objetivo**: completar operación y estabilizar el MVP.

**Alcance frontend**
- listado de pedidos admin
- detalle de pedido admin
- cambio de estado de pedido
- revisión de UX y mensajes de error
- ajustes responsive y accesibilidad básica

**Alcance backend**
- `GET /api/v1/admin/orders`
- `GET /api/v1/admin/orders/:id`
- `PATCH /api/v1/admin/orders/:id/status`
- validación de transiciones de estado

**Resultado esperado**
- tienda operable de punta a punta para cliente y administrador

---

## 21. Checklist Ejecutable por Equipo

Esta checklist separa responsabilidades para facilitar ejecución real.

### 21.1 Frontend

#### Base de aplicación
- [ ] configurar routing público
- [ ] configurar routing privado
- [ ] configurar routing admin
- [ ] crear layout público
- [ ] crear layout admin
- [ ] implementar navegación principal

#### Catálogo
- [ ] construir home simple
- [ ] construir listado de productos
- [ ] implementar buscador
- [ ] implementar filtros por categoría, precio, color y talla
- [ ] sincronizar filtros con URL
- [ ] mostrar loading, error y estado vacío

#### Detalle y carrito
- [ ] construir detalle de producto
- [ ] implementar selector de variantes
- [ ] bloquear agregar al carrito sin variante válida
- [ ] mostrar stock disponible o agotado
- [ ] construir carrito
- [ ] permitir editar cantidad
- [ ] permitir remover items
- [ ] persistir carrito en localStorage

#### Auth y sesión
- [ ] construir login
- [ ] construir registro
- [ ] guardar sesión según estrategia definida
- [ ] consultar usuario autenticado
- [ ] implementar `ProtectedRoute`
- [ ] implementar `AdminRoute`
- [ ] restaurar navegación después de login

#### Checkout y pedidos
- [ ] construir checkout
- [ ] integrar PayPal SDK
- [ ] enviar payload de compra al backend
- [ ] limpiar carrito tras compra exitosa
- [ ] construir historial de pedidos
- [ ] construir detalle de pedido

#### Admin
- [ ] construir listado de productos admin
- [ ] construir formulario de producto admin
- [ ] construir gestión de variantes
- [ ] construir editor de stock
- [ ] construir listado de pedidos admin
- [ ] construir detalle de pedido admin
- [ ] permitir cambio de estado de pedido

#### Calidad frontend
- [ ] validar responsive móvil y escritorio
- [ ] validar accesibilidad básica
- [ ] estandarizar estados de loading/error/success
- [ ] revisar consistencia de navegación y CTAs

### 21.2 Backend

#### Catálogo público
- [ ] implementar listado de productos activos
- [ ] implementar filtros por categoría, precio, color y talla
- [ ] implementar búsqueda por nombre
- [ ] implementar paginación
- [ ] implementar detalle de producto con variantes

#### Modelo de dominio
- [ ] definir entidad `products`
- [ ] definir entidad `product_variants`
- [ ] definir relación producto-variante
- [ ] asegurar stock por variante
- [ ] asegurar estado activo/inactivo de producto

#### Auth
- [ ] implementar registro de usuario
- [ ] implementar login
- [ ] emitir JWT
- [ ] exponer endpoint `me`
- [ ] proteger rutas privadas
- [ ] proteger rutas admin por rol

#### Pedidos y checkout
- [ ] implementar creación de orden
- [ ] validar stock antes de crear orden
- [ ] recalcular total en backend
- [ ] registrar referencia de pago PayPal
- [ ] definir estados válidos de pedido
- [ ] exponer historial de pedidos por usuario
- [ ] exponer detalle de pedido por usuario

#### Admin de productos
- [ ] implementar listado admin de productos
- [ ] implementar creación de producto
- [ ] implementar edición de producto
- [ ] implementar actualización de variantes
- [ ] implementar actualización de stock
- [ ] implementar activación/desactivación

#### Admin de pedidos
- [ ] implementar listado admin de pedidos
- [ ] implementar detalle admin de pedido
- [ ] implementar cambio de estado de pedido
- [ ] validar transiciones de estado

#### Calidad backend
- [ ] estandarizar respuestas de error
- [ ] validar payloads de entrada
- [ ] cubrir errores de stock insuficiente
- [ ] evitar exposición de productos inactivos
- [ ] asegurar autorización correcta en rutas privadas/admin

### 21.3 Dependencias compartidas

- [ ] acordar estrategia de almacenamiento del JWT
- [ ] acordar estructura final de errores de API
- [ ] acordar estructura final de paginación
- [ ] acordar contrato exacto de PayPal
- [ ] acordar convenciones de timestamps y estados
- [ ] confirmar naming final de recursos y propagarlo en backend/frontend

---

## 22. Recomendación de Ejecución

Para arrancar bien, conviene trabajar en este orden:

1. cerrar contrato backend
2. definir modelo de datos real de productos y variantes
3. construir frontend del catálogo contra contrato estable
4. cerrar checkout y órdenes
5. habilitar admin

La principal decisión pendiente antes de codificar fuerte es esta:

- **usar `orders` como naming oficial en frontend, backend y documentación**

Decisión recomendada: **usar `orders`** para mantener el sistema más claro y consistente.
