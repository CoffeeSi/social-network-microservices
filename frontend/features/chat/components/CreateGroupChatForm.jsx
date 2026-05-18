'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { routes } from '@/app/router/routes';
import { createGroupChat } from '@/features/chat/api/chatApi';
import { Button } from '@/shared/ui/Button';

export function CreateGroupChatForm({ onCreated }) {
  const router = useRouter();
  const [name, setName] = useState('');
  const [participants, setParticipants] = useState('');
  const [error, setError] = useState('');
  const [pending, setPending] = useState(false);
  const [open, setOpen] = useState(false);

  async function onSubmit(e) {
    e.preventDefault();
    const ids = participants
      .split(/[\s,;]+/)
      .map((s) => s.trim())
      .filter(Boolean);
    if (!name.trim() || ids.length === 0) {
      setError('Enter a group name and at least one participant user id.');
      return;
    }
    setPending(true);
    setError('');
    try {
      const res = await createGroupChat({ name: name.trim(), participant_ids: ids });
      const chat = res?.chat ?? res;
      setName('');
      setParticipants('');
      setOpen(false);
      onCreated?.();
      if (chat?.id) router.push(routes.chatThread(chat.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not create group');
    } finally {
      setPending(false);
    }
  }

  if (!open) {
    return (
      <Button type="button" variant="secondary" className="small-btn chat-sidebar__new" onClick={() => setOpen(true)}>
        New group
      </Button>
    );
  }

  return (
    <form className="group-chat-form" onSubmit={onSubmit}>
      {error ? <p className="form__error small">{error}</p> : null}
      <label className="field">
        <span className="field__label">Group name</span>
        <input className="input" value={name} onChange={(ev) => setName(ev.target.value)} required />
      </label>
      <label className="field">
        <span className="field__label">Participant user IDs (comma-separated)</span>
        <textarea
          className="input input--area"
          rows={2}
          value={participants}
          onChange={(ev) => setParticipants(ev.target.value)}
          placeholder="uuid-1, uuid-2"
          required
        />
      </label>
      <div className="group-chat-form__actions">
        <Button type="submit" disabled={pending}>
          {pending ? 'Creating…' : 'Create'}
        </Button>
        <Button type="button" variant="secondary" disabled={pending} onClick={() => setOpen(false)}>
          Cancel
        </Button>
      </div>
    </form>
  );
}
