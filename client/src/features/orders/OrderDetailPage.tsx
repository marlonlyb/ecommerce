import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';

import { fetchOrderById } from './api';
import type { Order } from '../../shared/types/order';
import { ORDER_STATUSES, PAYMENT_STATUSES } from '../../shared/types/order';
import { AppError } from '../../shared/api/errors';

const STATUS_LABELS: Record<string, string> = {
  [ORDER_STATUSES.PENDING_PAYMENT]: 'Pending Payment',
  [ORDER_STATUSES.PAID]: 'Paid',
  [ORDER_STATUSES.PAYMENT_FAILED]: 'Payment Failed',
  [ORDER_STATUSES.CANCELLED]: 'Cancelled',
  [ORDER_STATUSES.REFUNDED]: 'Refunded',
};

const PAYMENT_STATUS_LABELS: Record<string, string> = {
  [PAYMENT_STATUSES.PENDING]: 'Pending',
  [PAYMENT_STATUSES.APPROVED]: 'Approved',
  [PAYMENT_STATUSES.CAPTURED]: 'Captured',
  [PAYMENT_STATUSES.FAILED]: 'Failed',
  [PAYMENT_STATUSES.REFUNDED]: 'Refunded',
};

export function OrderDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    let cancelled = false;

    fetchOrderById(id)
      .then((data) => {
        if (!cancelled) {
          setOrder(data);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (!cancelled) {
          setError(
            err instanceof AppError ? err.message : 'Failed to load order.',
          );
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [id]);

  if (loading) {
    return (
      <section className="card-stack">
        <p className="orders__loading">Loading order…</p>
      </section>
    );
  }

  if (error || !order) {
    return (
      <section className="card-stack">
        <article className="card">
          <div className="orders__error" role="alert">{error ?? 'Order not found.'}</div>
          <Link className="btn btn--ghost" to="/profile/orders">
            Back to orders
          </Link>
        </article>
      </section>
    );
  }

  return (
    <section className="card-stack">
      <article className="card">
        <Link className="detail__back" to="/profile/orders">
          ← Back to orders
        </Link>

        <p className="eyebrow">Order detail</p>
        <h2>Order #{order.id}</h2>

        <dl className="order-detail__fields">
          <dt>Status</dt>
          <dd>
            <span className={`orders__status orders__status--${order.status}`}>
              {STATUS_LABELS[order.status] ?? order.status}
            </span>
          </dd>

          <dt>Payment</dt>
          <dd>
            {PAYMENT_STATUS_LABELS[order.payment_status] ?? order.payment_status}
            {order.payment_provider !== '' ? ` via ${order.payment_provider}` : ''}
          </dd>

          <dt>Created</dt>
          <dd>{new Date(order.created_at).toLocaleString()}</dd>

          {order.paid_at ? (
            <>
              <dt>Paid at</dt>
              <dd>{new Date(order.paid_at).toLocaleString()}</dd>
            </>
          ) : null}

          <dt>Currency</dt>
          <dd>{order.currency}</dd>
        </dl>
      </article>

      <article className="card">
        <h3>Items</h3>

        <div className="order-detail__items">
          {order.items.map((item) => (
            <div key={item.id} className="order-detail__item">
              <div className="order-detail__item-info">
                <strong>{item.product_name}</strong>
                <span className="order-detail__item-variant">
                  {item.color} / {item.size} — {item.variant_sku}
                </span>
              </div>
              <div className="order-detail__item-pricing">
                <span>${item.unit_price.toFixed(2)} × {item.quantity}</span>
                <strong>${item.line_total.toFixed(2)}</strong>
              </div>
            </div>
          ))}
        </div>

        <div className="order-detail__totals">
          <div className="order-detail__totals-row">
            <span>Subtotal</span>
            <span>${order.subtotal.toFixed(2)}</span>
          </div>
          <div className="order-detail__totals-row order-detail__totals-row--total">
            <span>Total</span>
            <span>${order.total.toFixed(2)}</span>
          </div>
        </div>
      </article>
    </section>
  );
}
