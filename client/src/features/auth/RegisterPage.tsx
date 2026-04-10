import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { register as apiRegister } from './api';
import { AppError, API_ERROR_CODES } from '../../shared/api/errors';

export function RegisterPage() {
  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);

    if (password !== confirmPassword) {
      setError('Passwords do not match.');
      return;
    }

    setSubmitting(true);

    try {
      await apiRegister({ email, password, confirm_password: confirmPassword });
      // After successful registration, redirect to login with a hint
      navigate('/login', { replace: true, state: { registered: true } });
    } catch (err) {
      if (err instanceof AppError) {
        if (err.code === API_ERROR_CODES.VALIDATION_ERROR) {
          // Try to extract field-specific errors
          const fieldErrors = err.details
            .map((d) => {
              const field = d.field ?? '';
              if (field.includes('email')) return 'Invalid email address or already in use.';
              if (field.includes('password')) return 'Password does not meet requirements.';
              return d.issue;
            })
            .join(' ');
          setError(fieldErrors || 'Please check your input and try again.');
        } else {
          setError(err.message);
        }
      } else {
        setError('An unexpected error occurred. Please try again.');
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <section className="auth-page">
      <article className="card">
        <p className="eyebrow">Create account</p>
        <h2>Register</h2>

        {error ? (
          <div className="auth-page__error" role="alert">
            {error}
          </div>
        ) : null}

        <form className="auth-page__form" onSubmit={handleSubmit}>
          <label className="auth-page__label">
            Email
            <input
              type="email"
              className="auth-page__input"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoComplete="email"
              disabled={submitting}
            />
          </label>

          <label className="auth-page__label">
            Password
            <input
              type="password"
              className="auth-page__input"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              autoComplete="new-password"
              minLength={6}
              disabled={submitting}
            />
          </label>

          <label className="auth-page__label">
            Confirm password
            <input
              type="password"
              className="auth-page__input"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              autoComplete="new-password"
              minLength={6}
              disabled={submitting}
            />
          </label>

          <button
            type="submit"
            className="btn btn--primary"
            disabled={submitting}
          >
            {submitting ? 'Creating account…' : 'Create account'}
          </button>
        </form>

        <p className="auth-page__alt">
          Already have an account?{' '}
          <Link to="/login">Sign in</Link>
        </p>
      </article>
    </section>
  );
}
