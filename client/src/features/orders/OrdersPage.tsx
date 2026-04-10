import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

import { fetchOrders } from './api';
import type { Order } from '../../shared/types/order';
import { ORDER_STATUSES } from '../../shared/types/order';
import { AppError } from '../../shared/api/errors';

const STATUS_LABELS: Record<string, string> = {
  [ORDER_STATUSES.PENDING_PAYMENT]: 'Pending Payment',
  [ORDER_STATUSES.PAID]: 'Paid',
  [ORDER_STATUSES.PAYMENT_FAILED]: 'Payment Failed',
  [ORDER_STATUSES.CANCELLED]: 'Cancelled',
  [ORDER_STATUSES.REFUNDED]: 'Refunded',
};

export function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    fetchOrders()
      .then((res) => {
        if (!cancelled) {
          setOrders(res.items);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (!cancelled) {
          setError(
            err instanceof AppError ? err.message : 'Failed to load orders.',
          );
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  if (loading) {
    return (
      <section className="card-stack">
        <p className="orders__loading">Loading orders…</p>
      </section>
    );
  }

  if (error) {
    return (
      <section className="card-stack">
        <article className="card">
          <div className="orders__error" role="alert">{error}</div>
        </article>
      </section>
    );
  }

  if (orders.length === 0) {
    return (
      <section className="card-stack">
        <article className="card">
          <p className="eyebrow">No orders yet</p>
          <h2>Order history</h2>
          <p>You haven't placed any orders yet.</p>
          <Link className="btn btn--primary" to="/products">
            Browse products
          </Link>
        </article>
      </section>
    );
  }

  return (
    <section className="card-stack">
      <article className="card">
        <p className="eyebrow">Your orders</p>
        <h2>Order history</h2>
      </article>

      <div className="orders__list">
        {orders.map((order) => (
          <Link
            key={order.id}
            to={`/profile/orders/${order.id}`}
            className="orders__item card"
          >
            <div className="orders__item-header">
              <span className="orders__item-id">#{order.id}</span>
              <span className={`orders__status orders__status--${order.status}`}>
                {STATUS_LABELS[order.status] ?? order.status}
              </span>
            </div>

            <div className="orders__item-meta">
              <span>{new Date(order.created_at).toLocaleDateString()}</span>
              <span>{order.items.length} item{order.items.length !== 1 ? 's' : ''}</span>
              <span className="orders__item-total">${order.total.toFixed(2)}</span>
            </div>
          </Link>
        ))}
      </div>
    </section>
  );
}
