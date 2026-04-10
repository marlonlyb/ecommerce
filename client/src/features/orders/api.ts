import { httpGet } from '../../shared/api/http';
import type { Order, OrderListResponse } from '../../shared/types/order';

/**
 * Fetch the authenticated user's orders.
 * GET /api/v1/private/orders → { data: { items: Order[] } }
 */
export function fetchOrders(): Promise<OrderListResponse> {
  return httpGet<OrderListResponse>('/api/v1/private/orders');
}

/**
 * Fetch a single order by ID (must belong to the authenticated user).
 * GET /api/v1/private/orders/:id → { data: Order }
 */
export function fetchOrderById(id: string): Promise<Order> {
  return httpGet<Order>(`/api/v1/private/orders/${id}`);
}
