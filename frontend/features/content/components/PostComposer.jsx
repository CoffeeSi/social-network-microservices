'use client';

import { useMemo, useState } from 'react';
import { createPost } from '@/features/content/api/contentApi';
import { Button } from '@/shared/ui/Button';
import { clsx } from '@/shared/lib/clsx';

export function PostComposer({ authorId, onCreated }) {
  const [content, setContent] = useState('');
  const [error, setError] = useState('');
  const [pending, setPending] = useState(false);
  const canSubmit = useMemo(() => Boolean(authorId) && content.trim().length > 0, [authorId, content]);

  async function submit(e) {
    e.preventDefault();
    if (!authorId || !canSubmit) return;
    setError('');
    setPending(true);
    try {
      const res = await createPost({ content: content.trim() });
      const post = res?.post ?? res;
      setContent('');
      onCreated?.(post);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not publish');
    } finally {
      setPending(false);
    }
  }

  return (
    <form className="composer" onSubmit={submit}>
      {!authorId ? (
        <p className="muted small">
          Could not read user id from the access token. Sign in again — the token should include a <code>sub</code>{' '}
          claim.
        </p>
      ) : null}
      {error ? <p className="form__error">{error}</p> : null}
      <label className="field">
        <span className="field__label">What’s on your mind?</span>
        <textarea
          name="content"
          className={clsx('input', 'input--area')}
          rows={4}
          value={content}
          onChange={(ev) => setContent(ev.target.value)}
          required
        />
      </label>
      <Button type="submit" disabled={pending || !canSubmit}>
        {pending ? 'Publishing…' : 'Publish'}
      </Button>
    </form>
  );
}
