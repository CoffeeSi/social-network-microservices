'use client';

import { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { routes } from '@/app/router/routes';
import {
  deletePost,
  getMyPosts,
  getPostStats,
  listPosts,
  toggleLike,
  updatePost,
} from '@/features/content/api/contentApi';
import { fetchDisplayNames } from '@/features/content/lib/authorNames';
import { PostComments } from '@/features/content/components/PostComments';
import { Button } from '@/shared/ui/Button';
import { Spinner } from '@/shared/ui/Spinner';
import { formatDateTime, pickTimestampField } from '@/shared/lib/formatDateTime';

function likeCountFromResponse(res) {
  return res?.new_like_count ?? res?.newLikeCount;
}

function PostCard({ post, currentUserId, authorName, onRemoved, onUpdated }) {
  const [commentsOpen, setCommentsOpen] = useState(false);
  const [commentCount, setCommentCount] = useState(post.comment_count ?? 0);
  const [likeCount, setLikeCount] = useState(post.like_count ?? 0);
  const [editing, setEditing] = useState(false);
  const [editContent, setEditContent] = useState(post.content ?? '');
  const [pending, setPending] = useState(false);
  const [error, setError] = useState('');

  const isOwn = Boolean(currentUserId && post.author_id === currentUserId);
  const authorLabel = authorName || 'Unknown';

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const stats = await getPostStats(post.id);
        if (!cancelled) {
          if (stats.like_count != null) setLikeCount(stats.like_count);
          if (stats.comment_count != null) setCommentCount(stats.comment_count);
        }
      } catch {
        /* stats optional */
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [post.id]);

  async function onLike() {
    if (!currentUserId) return;
    setError('');
    try {
      const res = await toggleLike(post.id);
      const next = likeCountFromResponse(res);
      if (next != null) setLikeCount(next);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not like post');
    }
  }

  async function onSaveEdit() {
    if (!editContent.trim()) return;
    setPending(true);
    setError('');
    try {
      const res = await updatePost(post.id, { content: editContent.trim() });
      const updated = res?.post ?? res;
      onUpdated?.({ ...post, ...updated, content: updated?.content ?? editContent.trim() });
      setEditing(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not update post');
    } finally {
      setPending(false);
    }
  }

  async function onDelete() {
    if (typeof window !== 'undefined' && !window.confirm('Delete this post?')) return;
    setError('');
    try {
      await deletePost(post.id);
      onRemoved?.(post.id);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Could not delete post');
    }
  }

  return (
    <li className="post-card">
      <header className="post-card__head">
        {post.author_id ? (
          <Link href={routes.userProfile(post.author_id)} className="text-link small post-card__author">
            {authorLabel}
          </Link>
        ) : (
          <span className="small">—</span>
        )}
        <time className="small muted">
          {formatDateTime(pickTimestampField(post, 'created_at', 'createdAt'))}
        </time>
      </header>
      {error ? <p className="form__error small">{error}</p> : null}
      {editing ? (
        <div className="post-card__edit">
          <textarea
            className="input input--area"
            rows={3}
            value={editContent}
            onChange={(ev) => setEditContent(ev.target.value)}
          />
          <div className="post-card__edit-actions">
            <Button type="button" disabled={pending} onClick={onSaveEdit}>
              Save
            </Button>
            <Button type="button" variant="secondary" disabled={pending} onClick={() => setEditing(false)}>
              Cancel
            </Button>
          </div>
        </div>
      ) : (
        <p className="post-card__body">{post.content}</p>
      )}
      <footer className="post-card__foot">
        <span className="small muted">
          {likeCount} likes · {commentCount} comments
        </span>
        <div className="post-card__actions">
          {currentUserId ? (
            <Button type="button" variant="secondary" className="small-btn" onClick={onLike}>
              Like
            </Button>
          ) : null}
          <Button
            type="button"
            variant="secondary"
            className="small-btn"
            onClick={() => setCommentsOpen((o) => !o)}
          >
            {commentsOpen ? 'Hide comments' : 'Comments'}
          </Button>
          {isOwn && !editing ? (
            <>
              <button type="button" className="msg-action" onClick={() => setEditing(true)}>
                Edit
              </button>
              <button type="button" className="msg-action msg-action--danger" onClick={onDelete}>
                Delete
              </button>
            </>
          ) : null}
        </div>
      </footer>
      <PostComments
        postId={post.id}
        currentUserId={currentUserId}
        open={commentsOpen}
        onCountChange={setCommentCount}
      />
    </li>
  );
}

/** @param {{ currentUserId?: string, refreshKey?: number, mode?: 'feed' | 'mine' }} props */
export function PostList({ currentUserId, refreshKey = 0, mode = 'feed' }) {
  const [state, setState] = useState({ loading: true, error: '', posts: [], nameByUserId: {} });

  const load = useCallback(async () => {
    setState((s) => ({ ...s, loading: true, error: '' }));
    try {
      const data =
        mode === 'mine' ? await getMyPosts({ page_size: 30, page: 1 }) : await listPosts({ page_size: 30, page: 1 });
      const posts = data.posts ?? [];
      const authorIds = posts.map((p) => p.author_id).filter(Boolean);
      const nameByUserId = await fetchDisplayNames(authorIds);
      setState({ loading: false, error: '', posts, nameByUserId });
    } catch (e) {
      setState({
        loading: false,
        error: e instanceof Error ? e.message : 'Failed to load posts',
        posts: [],
        nameByUserId: {},
      });
    }
  }, [mode]);

  useEffect(() => {
    load();
  }, [load, refreshKey]);

  function onRemoved(id) {
    setState((s) => ({ ...s, posts: s.posts.filter((p) => p.id !== id) }));
  }

  function onUpdated(updated) {
    setState((s) => ({
      ...s,
      posts: s.posts.map((p) => (p.id === updated.id ? { ...p, ...updated } : p)),
    }));
  }

  if (state.loading) return <Spinner />;
  if (state.error) return <p className="muted">{state.error}</p>;
  if (state.posts.length === 0) {
    return (
      <p className="muted">
        {mode === 'mine' ? 'You have not posted yet.' : 'No posts yet. Be the first to share something.'}
      </p>
    );
  }

  return (
    <ul className="post-list">
      {state.posts.map((post) => (
        <PostCard
          key={post.id}
          post={post}
          currentUserId={currentUserId}
          authorName={state.nameByUserId[post.author_id]}
          onRemoved={onRemoved}
          onUpdated={onUpdated}
        />
      ))}
    </ul>
  );
}
