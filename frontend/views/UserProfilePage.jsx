'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/app/providers/AuthProvider';
import { routes } from '@/app/router/routes';
import { createDirectChat } from '@/features/chat/api/chatApi';
import { changePassword, deleteUser, getUser, updateUser } from '@/features/users/api/usersApi';
import { Button } from '@/shared/ui/Button';
import { Card } from '@/shared/ui/Card';
import { Spinner } from '@/shared/ui/Spinner';
import { formatDateTime, parseTimestamp, pickTimestampField } from '@/shared/lib/formatDateTime';

function dobInputValue(user) {
  const raw = pickTimestampField(user, 'dob', 'date_of_birth');
  const d = parseTimestamp(raw);
  return d ? d.toISOString().slice(0, 10) : '';
}

export function UserProfilePage({ userId }) {
  const { userId: meId, ready, signOut } = useAuth();
  const router = useRouter();
  const [user, setUser] = useState(null);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(true);
  const [msgPending, setMsgPending] = useState(false);
  const [editMode, setEditMode] = useState(false);
  const [form, setForm] = useState({ first_name: '', last_name: '', email: '', dob: '' });
  const [pwd, setPwd] = useState({ old_password: '', new_password: '' });
  const [savePending, setSavePending] = useState(false);

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
          setForm({
            first_name: u.first_name ?? '',
            last_name: u.last_name ?? '',
            email: u.email ?? '',
            dob: dobInputValue(u),
          });
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

  async function onSaveProfile(e) {
    e.preventDefault();
    if (!user?.id) return;
    setSavePending(true);
    setError('');
    setSuccess('');
    try {
      const body = {
        first_name: form.first_name.trim(),
        last_name: form.last_name.trim(),
        email: form.email.trim(),
      };
      if (form.dob) body.dob = form.dob;
      const updated = await updateUser(user.id, body);
      setUser(updated?.user ?? updated);
      setEditMode(false);
      setSuccess('Profile updated.');
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Could not update profile');
    } finally {
      setSavePending(false);
    }
  }

  async function onChangePassword(e) {
    e.preventDefault();
    if (!user?.id) return;
    setSavePending(true);
    setError('');
    setSuccess('');
    try {
      await changePassword(user.id, {
        old_password: pwd.old_password,
        new_password: pwd.new_password,
      });
      setPwd({ old_password: '', new_password: '' });
      setSuccess('Password changed.');
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Could not change password');
    } finally {
      setSavePending(false);
    }
  }

  async function onDeleteAccount() {
    if (!user?.id) return;
    if (typeof window !== 'undefined' && !window.confirm('Delete your account permanently?')) return;
    setSavePending(true);
    setError('');
    try {
      await deleteUser(user.id);
      signOut();
      router.push(routes.home);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Could not delete account');
    } finally {
      setSavePending(false);
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
        {success ? <p className="banner-success">{success}</p> : null}
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
            <dd>{formatDateTime(pickTimestampField(user, 'created_at', 'createdAt'))}</dd>
          </div>
        </dl>

        {isSelf ? (
          <>
            {!editMode ? (
              <Button type="button" className="small-btn" onClick={() => setEditMode(true)}>
                Edit profile
              </Button>
            ) : (
              <form className="form profile-form" onSubmit={onSaveProfile}>
                <label className="field">
                  <span className="field__label">First name</span>
                  <input
                    className="input"
                    value={form.first_name}
                    onChange={(ev) => setForm((f) => ({ ...f, first_name: ev.target.value }))}
                    required
                  />
                </label>
                <label className="field">
                  <span className="field__label">Last name</span>
                  <input
                    className="input"
                    value={form.last_name}
                    onChange={(ev) => setForm((f) => ({ ...f, last_name: ev.target.value }))}
                    required
                  />
                </label>
                <label className="field">
                  <span className="field__label">Email</span>
                  <input
                    className="input"
                    type="email"
                    value={form.email}
                    onChange={(ev) => setForm((f) => ({ ...f, email: ev.target.value }))}
                    required
                  />
                </label>
                <label className="field">
                  <span className="field__label">Date of birth</span>
                  <input
                    className="input"
                    type="date"
                    value={form.dob}
                    onChange={(ev) => setForm((f) => ({ ...f, dob: ev.target.value }))}
                  />
                </label>
                <div className="profile-form__actions">
                  <Button type="submit" disabled={savePending}>
                    Save
                  </Button>
                  <Button type="button" variant="secondary" disabled={savePending} onClick={() => setEditMode(false)}>
                    Cancel
                  </Button>
                </div>
              </form>
            )}

            <form className="form profile-form" onSubmit={onChangePassword}>
              <h3 className="profile-form__heading">Change password</h3>
              <label className="field">
                <span className="field__label">Current password</span>
                <input
                  className="input"
                  type="password"
                  value={pwd.old_password}
                  onChange={(ev) => setPwd((p) => ({ ...p, old_password: ev.target.value }))}
                  required
                />
              </label>
              <label className="field">
                <span className="field__label">New password</span>
                <input
                  className="input"
                  type="password"
                  value={pwd.new_password}
                  onChange={(ev) => setPwd((p) => ({ ...p, new_password: ev.target.value }))}
                  required
                />
              </label>
              <Button type="submit" disabled={savePending}>
                Update password
              </Button>
            </form>

            <div className="profile-form__danger">
              <Button type="button" variant="secondary" disabled={savePending} onClick={onDeleteAccount}>
                Delete account
              </Button>
            </div>
          </>
        ) : null}

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
