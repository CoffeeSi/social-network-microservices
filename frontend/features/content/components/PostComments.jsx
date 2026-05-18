'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { routes } from '@/app/router/routes';
import {
  createComment,
  deleteComment,
  listComments,
  updateComment,
} from '@/features/content/api/contentApi';
import { fetchDisplayNames } from '@/features/content/lib/authorNames';
import { Button } from '@/shared/ui/Button';
import { Spinner } from '@/shared/ui/Spinner';
import { formatDateTime, pickTimestampField } from '@/shared/lib/formatDateTime';

export function PostComments({ postId, currentUserId, open, onCountChange }) {
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [text, setText] = useState('');
  const [pending, setPending] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [editText, setEditText] = useState('');
  const [nameByUserId, setNameByUserId] = useState({});

  useEffect(() => {
    if (!open || !postId) return;
    let cancelled = false;
    setLoading(true);
    (async () => {
      try {
        const data = await listComments(postId, { page_size: 50, page: 1 });
        const list = data.comments ?? [];
        const names = await fetchDisplayNames(list.map((c) => c.user_id));
        if (!cancelled) {
          setComments(list);
          setNameByUserId(names);
          setError('');
          onCountChange?.(list.length);
        }
      } catch (e) {
        if (!cancelled) setError(e instanceof Error ? e.message : 'Failed to load comments');
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [open, postId, onCountChange]);

  async function onAdd(e) {
    e.preventDefault();
    if (!text.trim()) return;
    setPending(true);
    setError('');
    try {
      const res = await createComment(postId, { text: text.trim() });
      const comment = res?.comment ?? res;
      if (comment?.id) {
        if (comment.user_id && currentUserId === comment.user_id) {
          const selfName = nameByUserId[currentUserId];
          if (selfName) {
            setNameByUserId((m) => ({ ...m, [comment.user_id]: selfName }));
          } else if (currentUserId) {
            const names = await fetchDisplayNames([comment.user_id]);
            setNameByUserId((m) => ({ ...m, ...names }));
          }
        }
        setComments((c) => {
          const next = [...c, comment];
          onCountChange?.(next.length);
          return next;
        });
      }
      setText('');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not post comment');
    } finally {
      setPending(false);
    }
  }

  async function onSaveEdit(commentId) {
    if (!editText.trim()) return;
    setPending(true);
    try {
      const res = await updateComment(commentId, { text: editText.trim() });
      const updated = res?.comment ?? res;
      setComments((list) =>
        list.map((c) =>
          c.id === commentId ? { ...c, ...updated, text: updated?.text ?? editText.trim() } : c,
        ),
      );
      setEditingId(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not update comment');
    } finally {
      setPending(false);
    }
  }

  async function onDelete(commentId) {
    if (typeof window !== 'undefined' && !window.confirm('Delete this comment?')) return;
    try {
      await deleteComment(commentId);
      setComments((list) => {
        const next = list.filter((c) => c.id !== commentId);
        onCountChange?.(next.length);
        return next;
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not delete comment');
    }
  }

  if (!open) return null;

  return (
    <div className="post-comments">
      {error ? <p className="form__error small">{error}</p> : null}
      {loading ? (
        <Spinner />
      ) : (
        <ul className="comment-list">
          {comments.map((c) => {
            const isOwn = Boolean(currentUserId && c.user_id === currentUserId);
            const isEditing = editingId === c.id;
            return (
              <li key={c.id} className="comment-list__item">
                {isEditing ? (
                  <div className="comment-list__edit">
                    <input
                      className="input"
                      value={editText}
                      onChange={(ev) => setEditText(ev.target.value)}
                    />
                    <div className="comment-list__edit-actions">
                      <Button type="button" disabled={pending} onClick={() => onSaveEdit(c.id)}>
                        Save
                      </Button>
                      <Button type="button" variant="secondary" disabled={pending} onClick={() => setEditingId(null)}>
                        Cancel
                      </Button>
                    </div>
                  </div>
                ) : (
                  <>
                    <p className="comment-list__text">{c.text}</p>
                    <div className="comment-list__meta">
                      {c.user_id ? (
                        <Link href={routes.userProfile(c.user_id)} className="text-link small">
                          {nameByUserId[c.user_id] ?? '…'}
                        </Link>
                      ) : null}
                      <time className="small muted">
                        {formatDateTime(pickTimestampField(c, 'created_at', 'createdAt'))}
                      </time>
                      {isOwn ? (
                        <span className="comment-list__actions">
                          <button
                            type="button"
                            className="msg-action"
                            onClick={() => {
                              setEditingId(c.id);
                              setEditText(c.text ?? '');
                            }}
                          >
                            Edit
                          </button>
                          <button type="button" className="msg-action msg-action--danger" onClick={() => onDelete(c.id)}>
                            Delete
                          </button>
                        </span>
                      ) : null}
                    </div>
                  </>
                )}
              </li>
            );
          })}
          {comments.length === 0 ? <li className="muted small">No comments yet.</li> : null}
        </ul>
      )}
      {currentUserId ? (
        <form className="comment-form" onSubmit={onAdd}>
          <input
            className="input"
            value={text}
            onChange={(ev) => setText(ev.target.value)}
            placeholder="Write a comment…"
          />
          <Button type="submit" disabled={pending || !text.trim()} className="small-btn">
            Reply
          </Button>
        </form>
      ) : null}
    </div>
  );
}
