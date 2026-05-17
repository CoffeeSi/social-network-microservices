'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/app/providers/AuthProvider';
import { routes } from '@/app/router/routes';
import { getUser } from '@/features/users/api/usersApi';
import { createDirectChat } from '@/features/chat/api/chatApi';
import { Button } from '@/shared/ui/Button';
import { Card } from '@/shared/ui/Card';
import { Spinner } from '@/shared/ui/Spinner';

function formatTs(ts) {
  if (ts == null) return '—';
  if (typeof ts === 'string') {
    try {
      return new Date(ts).toLocaleString();
    } catch {
      return ts;
    }
  }
  if (typeof ts === 'object' && ts.seconds != null) {
    try {
      return new Date(Number(ts.seconds) * 1000).toLocaleString();
    } catch {
      return '—';
    }
  }
  return '—';
}

export function UserProfilePage({ userId }) {
  const { userId: meId, ready } = useAuth();
  const router = useRouter();
  const [user, setUser] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);
  const [msgPending, setMsgPending] = useState(false);

  useEffect(() => {
    if (!userId) {
      setLoading(false);
      setError('User not specified');
      return;
    }
    let cancelled = false;
    setLoading(true);
    (async () => {
      try {
        const u = await getUser(userId);
        if (!cancelled) {
          setUser(u);
          setError('');
        }
      } catch (e) {
        if (!cancelled) setError(e instanceof Error ? e.message : 'Failed to load profile');
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [userId]);

  async function onMessage() {
    if (!meId || !user?.id || meId === user.id) return;
    setMsgPending(true);
    setError('');
    try {
      const res = await createDirectChat({ target_user_id: user.id });
      const chat = res?.chat ?? res;
      const id = chat?.id;
      if (id) router.push(routes.chatThread(id));
      else setError('Chat created but response had no id — check the API');
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Could not open chat');
    } finally {
      setMsgPending(false);
    }
  }

  if (loading) return <Spinner />;
  if (error && !user) return <p className="muted">{error}</p>;
  if (!user) return <p className="muted">User not found.</p>;

  const isSelf = Boolean(meId && user.id === meId);
  const canMessage = ready && Boolean(meId) && !isSelf;

  return (
    <div className="page page--narrow" style={{ maxWidth: 520 }}>
      <Card title={`${user.first_name} ${user.last_name}`}>
        {error ? <p className="form__error">{error}</p> : null}
        <dl className="detail-list">
          <div>
            <dt>Email</dt>
            <dd>{user.email}</dd>
          </div>
          <div>
            <dt>Active</dt>
            <dd>{user.is_active ? 'Yes' : 'No'}</dd>
          </div>
          <div>
            <dt>Joined</dt>
            <dd>{formatTs(user.created_at)}</dd>
          </div>
        </dl>
        <div className="profile-actions">
          {canMessage ? (
            <Button type="button" disabled={msgPending} onClick={onMessage}>
              {msgPending ? 'Opening…' : 'Message'}
            </Button>
          ) : null}
          <Link href={routes.users} className="btn btn--secondary profile-actions__link">
            Back to people
          </Link>
        </div>
      </Card>
    </div>
  );
}
