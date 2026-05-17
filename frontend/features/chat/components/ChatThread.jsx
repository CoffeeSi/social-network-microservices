'use client';

import Link from 'next/link';
import { useAuth } from '@/app/providers/AuthProvider';
import { routes } from '@/app/router/routes';
import {
  deleteMessage,
  editMessage,
  getMessages,
  sendMessage,
} from '@/features/chat/api/chatApi';
import { Button } from '@/shared/ui/Button';
import { Spinner } from '@/shared/ui/Spinner';
import { clsx } from '@/shared/lib/clsx';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';

function messageTimeMs(m) {
  const t = m?.created_at;
  if (t == null) return 0;
  if (typeof t === 'string') {
    const ms = Date.parse(t);
    return Number.isFinite(ms) ? ms : 0;
  }
  if (typeof t === 'object' && t.seconds != null) return Number(t.seconds) * 1000;
  return 0;
}

function sortMessagesAsc(list) {
  return [...(list ?? [])].sort((a, b) => messageTimeMs(a) - messageTimeMs(b));
}

function formatMessageClock(m) {
  const ms = messageTimeMs(m);
  if (!ms) return '';
  return new Date(ms).toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' });
}

function messageIso(m) {
  const ms = messageTimeMs(m);
  return ms ? new Date(ms).toISOString() : undefined;
}

export function ChatThread({ chatId }) {
  const { userId } = useAuth();
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [text, setText] = useState('');
  const [sending, setSending] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [editText, setEditText] = useState('');
  const [editPending, setEditPending] = useState(false);
  const scrollRef = useRef(null);

  const sorted = useMemo(() => sortMessagesAsc(messages), [messages]);

  const scrollToBottom = useCallback(() => {
    const el = scrollRef.current;
    if (!el) return;
    el.scrollTop = el.scrollHeight;
  }, []);

  useEffect(() => {
    if (!chatId) {
      setMessages([]);
      setLoading(false);
      return;
    }
    let cancelled = false;
    setLoading(true);
    setEditingId(null);
    (async () => {
      try {
        const data = await getMessages(chatId, { page: 1, page_size: 100 });
        if (!cancelled) {
          setMessages(data.messages ?? []);
          setError('');
        }
      } catch (e) {
        if (!cancelled) setError(e instanceof Error ? e.message : 'Failed to load messages');
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [chatId]);

  useEffect(() => {
    if (!loading) scrollToBottom();
  }, [sorted, loading, scrollToBottom]);

  async function onSend(e) {
    e.preventDefault();
    if (!chatId || !text.trim()) return;
    setSending(true);
    setError('');
    try {
      const res = await sendMessage(chatId, { content: text.trim() });
      if (res?.message) setMessages((m) => [...m, res.message]);
      setText('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Send failed');
    } finally {
      setSending(false);
    }
  }

  function startEdit(msg) {
    setEditingId(msg.id);
    setEditText(msg.content ?? '');
  }

  function cancelEdit() {
    setEditingId(null);
    setEditText('');
  }

  async function saveEdit(messageId) {
    if (!chatId || !editText.trim()) return;
    setEditPending(true);
    setError('');
    try {
      const res = await editMessage(chatId, messageId, { new_content: editText.trim() });
      const updated = res?.message;
      if (updated) {
        setMessages((list) => list.map((m) => (m.id === messageId ? { ...m, ...updated } : m)));
      }
      cancelEdit();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Edit failed');
    } finally {
      setEditPending(false);
    }
  }

  async function onDelete(messageId) {
    if (!chatId) return;
    if (typeof window !== 'undefined' && !window.confirm('Delete this message?')) return;
    setError('');
    try {
      await deleteMessage(chatId, messageId);
      setMessages((list) => list.filter((m) => m.id !== messageId));
      if (editingId === messageId) cancelEdit();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Delete failed');
    }
  }

  if (!chatId) {
    return <p className="muted">Pick a chat from the list.</p>;
  }
  if (loading) return <Spinner />;

  return (
    <div className="chat-thread">
      {error ? <p className="form__error">{error}</p> : null}
      <div className="chat-thread__scroll" ref={scrollRef}>
        <ul className="msg-list">
          {sorted.map((msg) => {
            const isOwn = Boolean(userId && msg.sender_id === userId);
            const isEditing = editingId === msg.id;
            const clock = formatMessageClock(msg);
            const iso = messageIso(msg);

            return (
              <li key={msg.id} className={clsx('msg-list__row', isOwn && 'msg-list__row--own')}>
                {isOwn && !isEditing ? (
                  <div className="msg-list__side-actions">
                    <button type="button" className="msg-action" onClick={() => startEdit(msg)}>
                      Edit
                    </button>
                    <button type="button" className="msg-action msg-action--danger" onClick={() => onDelete(msg.id)}>
                      Delete
                    </button>
                  </div>
                ) : null}
                <div className={clsx('msg-list__bubble', isOwn && 'msg-list__bubble--own')}>
                  {isEditing ? (
                    <div className="msg-list__edit">
                      <textarea
                        className="input input--area msg-list__edit-input"
                        rows={3}
                        value={editText}
                        onChange={(ev) => setEditText(ev.target.value)}
                      />
                      <div className="msg-list__edit-actions">
                        <Button
                          type="button"
                          disabled={editPending || !editText.trim()}
                          onClick={() => saveEdit(msg.id)}
                        >
                          {editPending ? '…' : 'Save'}
                        </Button>
                        <Button type="button" variant="secondary" disabled={editPending} onClick={cancelEdit}>
                          Cancel
                        </Button>
                      </div>
                    </div>
                  ) : (
                    <>
                      <p className="msg-list__text">{msg.content}</p>
                      <div className="msg-list__meta-line">
                        {msg.sender_id ? (
                          <Link href={routes.userProfile(msg.sender_id)} className="msg-list__profile-link">
                            {isOwn ? 'You' : 'Profile'}
                          </Link>
                        ) : null}
                        {clock ? (
                          <time className="msg-list__time" dateTime={iso}>
                            {clock}
                          </time>
                        ) : null}
                      </div>
                    </>
                  )}
                </div>
              </li>
            );
          })}
        </ul>
      </div>
      <div className="chat-thread__input-wrap">
        <form className="chat-thread__input" onSubmit={onSend}>
          <input
            className="input"
            value={text}
            onChange={(ev) => setText(ev.target.value)}
            placeholder="Message…"
          />
          <Button type="submit" disabled={sending || !text.trim()}>
            Send
          </Button>
        </form>
      </div>
    </div>
  );
}
