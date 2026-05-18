'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { routes } from '@/app/router/routes';
import {
  addParticipant,
  deleteGroupChat,
  editGroupChat,
  leaveGroupChat,
  removeParticipant,
} from '@/features/chat/api/chatApi';
import { fetchDisplayNames } from '@/features/content/lib/authorNames';
import { Button } from '@/shared/ui/Button';

export function ChatGroupPanel({ chat, onChanged }) {
  const router = useRouter();
  const [name, setName] = useState(chat?.name ?? '');
  const [newUserId, setNewUserId] = useState('');
  const [error, setError] = useState('');
  const [pending, setPending] = useState(false);
  const [nameByUserId, setNameByUserId] = useState({});

  const participants = chat?.participant_ids ?? [];

  useEffect(() => {
    if (!chat?.is_group || participants.length === 0) {
      setNameByUserId({});
      return;
    }
    let cancelled = false;
    (async () => {
      const names = await fetchDisplayNames(participants);
      if (!cancelled) setNameByUserId(names);
    })();
    return () => {
      cancelled = true;
    };
  }, [chat?.id, chat?.is_group, participants.join(',')]);

  if (!chat?.is_group) return null;

  async function onRename(e) {
    e.preventDefault();
    if (!name.trim()) return;
    setPending(true);
    setError('');
    try {
      await editGroupChat(chat.id, { name: name.trim() });
      onChanged?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not rename group');
    } finally {
      setPending(false);
    }
  }

  async function onAdd(e) {
    e.preventDefault();
    if (!newUserId.trim()) return;
    setPending(true);
    setError('');
    try {
      await addParticipant(chat.id, { user_id: newUserId.trim() });
      setNewUserId('');
      onChanged?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not add participant');
    } finally {
      setPending(false);
    }
  }

  async function onRemove(userId) {
    setPending(true);
    setError('');
    try {
      await removeParticipant(chat.id, userId);
      onChanged?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not remove participant');
    } finally {
      setPending(false);
    }
  }

  async function onLeave() {
    if (typeof window !== 'undefined' && !window.confirm('Leave this group?')) return;
    setPending(true);
    try {
      await leaveGroupChat(chat.id);
      router.push(routes.chat);
      onChanged?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not leave group');
    } finally {
      setPending(false);
    }
  }

  async function onDelete() {
    if (typeof window !== 'undefined' && !window.confirm('Delete this group for everyone?')) return;
    setPending(true);
    try {
      await deleteGroupChat(chat.id);
      router.push(routes.chat);
      onChanged?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not delete group');
    } finally {
      setPending(false);
    }
  }

  return (
    <details className="chat-group-panel">
      <summary className="chat-group-panel__summary small">Group settings</summary>
      <div className="chat-group-panel__body">
        {error ? <p className="form__error small">{error}</p> : null}
        <form className="chat-group-panel__form" onSubmit={onRename}>
          <label className="field">
            <span className="field__label">Name</span>
            <input className="input" value={name} onChange={(ev) => setName(ev.target.value)} />
          </label>
          <Button type="submit" disabled={pending} className="small-btn">
            Save name
          </Button>
        </form>
        <form className="chat-group-panel__form" onSubmit={onAdd}>
          <label className="field">
            <span className="field__label">Add participant (user id)</span>
            <input className="input" value={newUserId} onChange={(ev) => setNewUserId(ev.target.value)} />
          </label>
          <Button type="submit" disabled={pending} className="small-btn">
            Add
          </Button>
        </form>
        <ul className="chat-group-panel__members">
          {participants.map((id) => (
            <li key={id} className="chat-group-panel__member">
              <Link href={routes.userProfile(id)} className="text-link small chat-group-panel__member-name">
                {nameByUserId[id] ?? '…'}
              </Link>
              <button
                type="button"
                className="msg-action msg-action--danger"
                disabled={pending}
                onClick={() => onRemove(id)}
              >
                Remove
              </button>
            </li>
          ))}
        </ul>
        <div className="chat-group-panel__danger">
          <Button type="button" variant="secondary" disabled={pending} onClick={onLeave}>
            Leave group
          </Button>
          <Button type="button" variant="secondary" disabled={pending} onClick={onDelete}>
            Delete group
          </Button>
        </div>
      </div>
    </details>
  );
}
