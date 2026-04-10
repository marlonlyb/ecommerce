import { useSession } from '../../app/providers/SessionProvider';
import { Link } from 'react-router-dom';

export function ProfilePage() {
  const { user, loading } = useSession();

  if (loading) {
    return (
      <section className="card-stack">
        <p className="profile__loading">Loading profile…</p>
      </section>
    );
  }

  if (!user) {
    // This shouldn't happen since the route is behind RequireAuth,
    // but handle gracefully.
    return (
      <section className="card-stack">
        <article className="card">
          <p>Session data not available. Please sign in again.</p>
        </article>
      </section>
    );
  }

  return (
    <section className="card-stack">
      <article className="card">
        <p className="eyebrow">Your profile</p>
        <h2>{user.email}</h2>

        <dl className="profile__fields">
          <dt>ID</dt>
          <dd>{user.id}</dd>

          <dt>Email</dt>
          <dd>{user.email}</dd>

          <dt>Role</dt>
          <dd>{user.is_admin ? 'Admin' : 'Customer'}</dd>

          <dt>Member since</dt>
          <dd>{new Date(user.created_at).toLocaleDateString()}</dd>
        </dl>
      </article>

      <article className="card card--muted">
        <h3>Quick links</h3>
        <div className="profile__links">
          <Link className="btn btn--ghost" to="/profile/orders">
            View order history
          </Link>
        </div>
      </article>
    </section>
  );
}
