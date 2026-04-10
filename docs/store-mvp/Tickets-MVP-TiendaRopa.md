# Tickets MVP - Tienda de Ropa

## Objetivo

Traducir el PRD del MVP a tickets accionables para un equipo pequeño, manteniendo foco en compra simple, operación admin básica y naming oficial `orders`.

## Convenciones

- Tipos: `feature`, `task`, `bug`, `spike`
- Prioridad: `P0`, `P1`, `P2`
- Área: `Frontend`, `Backend`, `Shared`
- Dependencias: usar IDs de tickets de este documento

---

## Épica E1 - Base pública y navegación

### FE-001 - Layout público y routing base
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Frontend
- **Descripción**: Implementar layout público base con navegación hacia home, catálogo, carrito, login, registro y perfil, dejando la estructura inicial lista para el resto del MVP.
- **Criterio de aceptación**:
  - Existen rutas públicas para `/`, `/products`, `/products/:id`, `/cart`, `/login`, `/register`, `/profile`.
  - El header permite navegar a catálogo, carrito y sesión.
  - La estructura funciona en móvil y escritorio sin romper navegación.
- **Dependencias**: ninguna

### FE-002 - Home con acceso claro al catálogo
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Frontend
- **Descripción**: Construir una home simple con CTA principal hacia catálogo y sección breve de productos destacados o recientes.
- **Criterio de aceptación**:
  - La home muestra CTA visible hacia `/products`.
  - Se muestra un bloque de productos destacados o recientes.
  - Hay estados de loading, error y vacío.
- **Dependencias**: FE-001, BE-001

### SH-001 - Definir convenciones de respuesta API y errores
- **Tipo**: task
- **Prioridad**: P0
- **Área**: Shared
- **Descripción**: Acordar estructura base de `data`, `meta` y `error` para evitar divergencias entre frontend y backend.
- **Criterio de aceptación**:
  - Existe un formato documentado para respuestas exitosas y de error.
  - Frontend y backend usan los mismos códigos de error principales.
  - El contrato distingue errores de validación, auth, negocio y sistema.
- **Dependencias**: ninguna

---

## Épica E2 - Catálogo y detalle de producto

### BE-001 - Endpoint público de listado de productos
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Backend
- **Descripción**: Implementar `GET /api/v1/public/products` con productos activos, búsqueda, filtros y paginación simple.
- **Criterio de aceptación**:
  - Solo devuelve productos activos.
  - Soporta `search`, `category`, `color`, `size`, `min_price`, `max_price`, `page`, `limit`, `sort`.
  - Devuelve `items` y `pagination` en estructura consistente.
- **Dependencias**: SH-001, BE-002

### BE-002 - Modelo de catálogo para producto y variantes
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Backend
- **Descripción**: Consolidar entidades y relaciones mínimas para `products` y `product_variants`, con stock por variante y control de activo/inactivo.
- **Criterio de aceptación**:
  - El dominio soporta talla, color, SKU, precio y stock por variante.
  - El producto puede activarse o desactivarse sin borrado físico.
  - El catálogo público puede filtrar usando datos reales del modelo.
- **Dependencias**: ninguna

### FE-003 - Página de catálogo con búsqueda y filtros
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Frontend
- **Descripción**: Consumir el listado público y construir la página de catálogo con buscador, filtros y estados UX básicos.
- **Criterio de aceptación**:
  - La página lista productos usando backend real.
  - Los filtros y búsqueda se reflejan en query params.
  - Existe acción de limpiar filtros.
  - El catálogo muestra loading, error y vacío de forma clara.
- **Dependencias**: FE-001, BE-001

### BE-003 - Endpoint público de detalle de producto
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Backend
- **Descripción**: Implementar `GET /api/v1/public/products/:id` con detalle completo, galería y variantes disponibles.
- **Criterio de aceptación**:
  - Devuelve nombre, descripción, imágenes, categoría y variantes.
  - Incluye `stock` y `price` por variante.
  - No expone productos inactivos por API pública.
- **Dependencias**: BE-002, SH-001

### FE-004 - Detalle de producto y selección de variante
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Frontend
- **Descripción**: Construir vista de detalle con galería, datos principales, selección de color/talla y validación visual de stock.
- **Criterio de aceptación**:
  - El usuario puede seleccionar variante antes de agregar al carrito.
  - Si no hay stock, el CTA queda bloqueado y el estado es visible.
  - El detalle usa datos reales del endpoint público.
- **Dependencias**: BE-003

---

## Épica E3 - Carrito y preparación de compra

### FE-005 - Store local de carrito y persistencia
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Frontend
- **Descripción**: Implementar el manejo local del carrito con persistencia en `localStorage` y validaciones básicas contra stock conocido.
- **Criterio de aceptación**:
  - El carrito persiste entre recargas.
  - Cada item guarda `product_id`, `variant_id`, talla, color, precio y cantidad.
  - No permite cantidades menores a 1 ni mayores al stock conocido.
- **Dependencias**: FE-004

### FE-006 - Vista de carrito con edición y resumen
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Frontend
- **Descripción**: Implementar `/cart` con edición de cantidades, eliminación de items y resumen de compra.
- **Criterio de aceptación**:
  - Se puede sumar, restar y remover items.
  - El total se recalcula correctamente.
  - Hay CTA claro para iniciar checkout.
- **Dependencias**: FE-005

---

## Épica E4 - Autenticación y sesión

### BE-004 - Registro, login y endpoint `me`
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Backend
- **Descripción**: Implementar `POST /api/v1/public/register`, `POST /api/v1/public/login` y `GET /api/v1/private/me` usando JWT y rol admin básico.
- **Criterio de aceptación**:
  - Registro crea usuario con email y password válidos.
  - Login devuelve token y datos mínimos del usuario.
  - `me` responde solo con token válido.
- **Dependencias**: SH-001

### FE-007 - Formularios de registro y login
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Frontend
- **Descripción**: Construir formularios de auth con validación básica, feedback de error y persistencia de sesión según el contrato acordado.
- **Criterio de aceptación**:
  - Login y registro consumen endpoints reales.
  - Los errores de validación y credenciales se muestran claramente.
  - La sesión queda disponible para rutas privadas.
- **Dependencias**: BE-004

### FE-008 - Protección de rutas cliente y admin
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Frontend
- **Descripción**: Proteger checkout, perfil y admin con guardas de navegación y redirect al destino original tras login.
- **Criterio de aceptación**:
  - Un visitante no autenticado no entra a checkout o perfil.
  - Un usuario no admin no accede al módulo admin.
  - Tras login exitoso el usuario vuelve al destino intentado.
- **Dependencias**: FE-007, BE-004

### SH-002 - Definir estrategia de almacenamiento del JWT
- **Tipo**: task
- **Prioridad**: P0
- **Área**: Shared
- **Descripción**: Acordar dónde vive el token y cómo se restaura la sesión para evitar ambigüedad de implementación.
- **Criterio de aceptación**:
  - La estrategia queda documentada y aplicada igual en frontend y backend.
  - Se define manejo de expiración y logout.
  - No se almacena información sensible adicional en `localStorage`.
- **Dependencias**: SH-001

---

## Épica E5 - Checkout y orders del cliente

### FE-009 - Página de checkout autenticado
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Frontend
- **Descripción**: Construir checkout simple con resumen de items, total, gating por auth y continuidad desde carrito.
- **Criterio de aceptación**:
  - Solo usuarios autenticados pueden confirmar compra.
  - Si el usuario llega sin sesión, es redirigido a login y luego vuelve a checkout.
  - Se usa el carrito real como fuente de datos.
- **Dependencias**: FE-006, FE-008

### BE-005 - Crear `orders` del cliente
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Backend
- **Descripción**: Implementar `POST /api/v1/private/orders` recalculando totales, validando stock y guardando referencia de pago.
- **Criterio de aceptación**:
  - El backend rechaza requests sin items válidos.
  - Valida stock actualizado por variante antes de confirmar.
  - Persiste orden, items, total y `payment_reference`.
  - Usa naming oficial `orders` en rutas, dominio y documentación.
- **Dependencias**: BE-002, BE-004, SH-001

### FE-010 - Integración PayPal y confirmación de order
- **Tipo**: feature
- **Prioridad**: P0
- **Área**: Frontend
- **Descripción**: Integrar PayPal en el checkout y enviar el payload final al backend para crear la `order`.
- **Criterio de aceptación**:
  - Se manejan estados de éxito, cancelación y error.
  - El frontend envía `payment_method`, `payment_reference` e items.
  - El carrito se limpia tras compra exitosa.
- **Dependencias**: FE-009, BE-005

### BE-006 - Listado y detalle de orders del cliente
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Backend
- **Descripción**: Implementar `GET /api/v1/private/orders` y `GET /api/v1/private/orders/:id` para historial y detalle del usuario autenticado.
- **Criterio de aceptación**:
  - Solo devuelve orders del usuario autenticado.
  - El detalle incluye items, total, estado y fecha.
  - Devuelve `404` si la order no pertenece al usuario.
- **Dependencias**: BE-005

### FE-011 - Historial y detalle de orders del cliente
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Frontend
- **Descripción**: Implementar `/profile/orders` y `/profile/orders/:id` con listado y detalle claros para seguimiento post-compra.
- **Criterio de aceptación**:
  - El usuario ve sus orders con estado, fecha y total.
  - Puede entrar al detalle de una order.
  - Se manejan estados vacíos y errores.
- **Dependencias**: FE-008, BE-006

### SH-003 - Cerrar contrato de PayPal y estados de payment
- **Tipo**: task
- **Prioridad**: P1
- **Área**: Shared
- **Descripción**: Definir el payload mínimo de integración PayPal y cómo se traduce a estados internos del MVP.
- **Criterio de aceptación**:
  - Existe un contrato claro para `payment_method` y `payment_reference`.
  - Se documenta cuándo una order queda `pending`, `paid` o `failed`.
  - Frontend y backend usan el mismo flujo nominal y de error.
- **Dependencias**: FE-009, BE-005

---

## Épica E6 - Administración de catálogo

### BE-007 - Listado admin de productos
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Backend
- **Descripción**: Implementar `GET /api/v1/admin/products` con filtros básicos para operación interna.
- **Criterio de aceptación**:
  - Permite buscar por nombre o slug.
  - Permite filtrar por estado activo/inactivo.
  - Requiere rol admin.
- **Dependencias**: BE-002, BE-004

### FE-012 - Listado admin de productos
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Frontend
- **Descripción**: Construir vista admin de productos con tabla simple, estado y accesos a crear/editar.
- **Criterio de aceptación**:
  - Lista productos usando endpoint admin real.
  - Permite navegar a crear y editar.
  - Muestra activo/inactivo de forma visible.
- **Dependencias**: FE-008, BE-007

### BE-008 - Crear y editar productos con variantes
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Backend
- **Descripción**: Implementar `POST /api/v1/admin/products` y `PUT /api/v1/admin/products/:id` para administrar producto y variantes del MVP.
- **Criterio de aceptación**:
  - Permite crear productos con variantes iniciales.
  - Permite editar datos principales y variantes.
  - Valida SKU único, precio válido y stock no negativo.
- **Dependencias**: BE-002, BE-004, SH-001

### FE-013 - Formulario admin de producto y variantes
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Frontend
- **Descripción**: Construir formulario de alta/edición de producto con gestión simple de variantes y stock.
- **Criterio de aceptación**:
  - Permite crear y editar producto.
  - Permite agregar y editar variantes con talla, color, SKU, precio y stock.
  - Muestra validaciones inline y errores de API.
- **Dependencias**: FE-012, BE-008

### BE-009 - Activar y desactivar productos
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Backend
- **Descripción**: Implementar `PATCH /api/v1/admin/products/:id/status` para controlar visibilidad del catálogo sin borrado físico.
- **Criterio de aceptación**:
  - Un admin puede cambiar `active` de forma explícita.
  - El cambio impacta catálogo público.
  - El endpoint requiere rol admin.
- **Dependencias**: BE-007

### FE-014 - Acción admin de activación y desactivación
- **Tipo**: feature
- **Prioridad**: P2
- **Área**: Frontend
- **Descripción**: Incorporar acción rápida para activar o desactivar productos desde el listado admin.
- **Criterio de aceptación**:
  - El cambio actualiza el estado visible sin recargar toda la app.
  - Hay feedback de éxito o error.
  - La acción está disponible solo para admins.
- **Dependencias**: FE-012, BE-009

---

## Épica E7 - Administración de orders

### BE-010 - Listado admin de orders
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Backend
- **Descripción**: Implementar `GET /api/v1/admin/orders` con filtros simples por estado y búsqueda operativa.
- **Criterio de aceptación**:
  - Devuelve orders con estado, total, fecha y datos básicos del cliente.
  - Soporta filtro por estado.
  - Requiere rol admin.
- **Dependencias**: BE-005, BE-004

### FE-015 - Listado admin de orders
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Frontend
- **Descripción**: Construir vista administrativa para revisar orders y entrar a su detalle.
- **Criterio de aceptación**:
  - Muestra listado con estado, total, fecha y cliente.
  - Permite filtrar por estado.
  - Permite navegar al detalle de la order.
- **Dependencias**: FE-008, BE-010

### BE-011 - Detalle admin y cambio de estado de orders
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Backend
- **Descripción**: Implementar `GET /api/v1/admin/orders/:id` y `PATCH /api/v1/admin/orders/:id/status` con validación de transiciones básicas.
- **Criterio de aceptación**:
  - El detalle devuelve items, totales, estado y referencia de pago.
  - El cambio de estado solo acepta estados válidos.
  - Rechaza transiciones inconsistentes según la regla definida.
- **Dependencias**: BE-010, SH-001

### FE-016 - Detalle admin y actualización de estado
- **Tipo**: feature
- **Prioridad**: P1
- **Área**: Frontend
- **Descripción**: Implementar detalle admin de order con selector de estado y feedback de actualización.
- **Criterio de aceptación**:
  - El admin ve el detalle completo de la order.
  - Puede actualizar el estado usando endpoint real.
  - Se reflejan errores de transición o autorización.
- **Dependencias**: FE-015, BE-011

---

## Épica E8 - Hardening MVP

### SH-004 - Revisar UX states y accesibilidad mínima
- **Tipo**: task
- **Prioridad**: P1
- **Área**: Shared
- **Descripción**: Recorrer el MVP y normalizar loading, error, vacío, validaciones visibles y accesibilidad básica.
- **Criterio de aceptación**:
  - Las vistas clave muestran estados consistentes.
  - Formularios y acciones principales son navegables por teclado.
  - Los mensajes críticos son comprensibles para cliente y admin.
- **Dependencias**: FE-003, FE-004, FE-006, FE-011, FE-013, FE-016

### BE-012 - Estandarizar validaciones, errores y logging mínimo
- **Tipo**: task
- **Prioridad**: P1
- **Área**: Backend
- **Descripción**: Consolidar validaciones de entrada, códigos de error y observabilidad mínima para soporte operativo del MVP.
- **Criterio de aceptación**:
  - Los endpoints usan estructura uniforme de error.
  - Los errores de negocio comunes tienen códigos estables.
  - Existen logs mínimos para auth, creación de orders y cambios admin.
- **Dependencias**: BE-005, BE-008, BE-011, SH-001

---

## Secuencia sugerida de ejecución

1. SH-001, BE-002, FE-001
2. BE-001, BE-003, FE-002, FE-003, FE-004
3. FE-005, FE-006
4. BE-004, SH-002, FE-007, FE-008
5. FE-009, BE-005, FE-010, BE-006, FE-011, SH-003
6. BE-007, BE-008, BE-009, FE-012, FE-013, FE-014
7. BE-010, BE-011, FE-015, FE-016
8. SH-004, BE-012
