'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { routes } from '@/app/router/routes';
import { listUsers } from '@/features/users/api/usersApi';
import { Spinner } from '@/shared/ui/Spinner';
import { clsx } from '@/shared/lib/clsx';
import { formatDateTime, pickTimestampField } from '@/shared/lib/formatDateTime';

export function UserList() {
  const [state, setState] = useState({ loading: true, error: '', data: null });

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const data = await listUsers({ page_size: 50, page: 1 });
        if (!cancelled) setState({ loading: false, error: '', data });
      } catch (e) {
        if (!cancelled)
          setState({ loading: false, error: e instanceof Error ? e.message : 'Failed to load', data: null });
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  if (state.loading) return <Spinner />;
  if (state.error) return <p className="muted">{state.error}</p>;

  const users = state.data?.users ?? [];
  if (users.length === 0) return <p className="muted">No users yet.</p>;

  return (
    <ul className="user-list">
      {users.map((u) => (
        <li key={u.id}>
          <Link href={routes.userProfile(u.id)} className={clsx('user-list__item', 'user-list__item--link')}>
            <span className="user-list__name">
              {u.first_name} {u.last_name}
            </span>
            <span className="user-list__meta">{u.email}</span>
            <span className="user-list__meta">
              Joined {formatDateTime(pickTimestampField(u, 'created_at', 'createdAt'))}
            </span>
          </Link>
        </li>
      ))}
    </ul>
  );
}
