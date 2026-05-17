import { useEffect, useState } from 'react';
import { listPosts, toggleLike } from '@/features/content/api/contentApi';
import { Spinner } from '@/shared/ui/Spinner';
import { Button } from '@/shared/ui/Button';

function formatPostTime(ts) {
  if (!ts) return '';
  try {
    return new Date(ts).toLocaleString();
  } catch {
    return String(ts);
  }
}

export function PostList({ currentUserId, refreshKey = 0 }) {
  const [state, setState] = useState({ loading: true, error: '', posts: [] });

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const data = await listPosts({ page_size: 30, page: 1 });
        if (!cancelled) setState({ loading: false, error: '', posts: data.posts ?? [] });
      } catch (e) {
        if (!cancelled)
          setState({
            loading: false,
            error: e instanceof Error ? e.message : 'Failed to load posts',
            posts: [],
          });
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [refreshKey]);

  async function onLike(post) {
    if (!currentUserId) return;
    try {
      const res = await toggleLike(post.id, { user_id: currentUserId });
      setState((s) => ({
        ...s,
        posts: s.posts.map((p) =>
          p.id === post.id ? { ...p, like_count: res.new_like_count ?? p.like_count } : p,
        ),
      }));
    } catch {
      /* noop */
    }
  }

  if (state.loading) return <Spinner />;
  if (state.error) return <p className="muted">{state.error}</p>;
  if (state.posts.length === 0) return <p className="muted">No posts yet. Be the first to share something.</p>;

  return (
    <ul className="post-list">
      {state.posts.map((post) => (
        <li key={post.id} className="post-card">
          <header className="post-card__head">
            <span className="mono small">author {post.author_id?.slice(0, 8) ?? '—'}…</span>
            <time className="small muted">{formatPostTime(post.created_at)}</time>
          </header>
          <p className="post-card__body">{post.content}</p>
          <footer className="post-card__foot">
            <span className="small muted">{post.like_count ?? 0} likes · {post.comment_count ?? 0} comments</span>
            {currentUserId ? (
              <Button type="button" variant="secondary" className="small-btn" onClick={() => onLike(post)}>
                Like
              </Button>
            ) : null}
          </footer>
        </li>
      ))}
    </ul>
  );
}
