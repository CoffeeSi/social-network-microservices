'use client';

import { useState } from 'react';
import { useAuth } from '@/app/providers/AuthProvider';
import { PostComposer, PostList } from '@/features/content';
import { Card } from '@/shared/ui/Card';
import { clsx } from '@/shared/lib/clsx';

export function FeedPage() {
  const { userId } = useAuth();
  const [refreshKey, setRefreshKey] = useState(0);
  const [tab, setTab] = useState('feed');

  return (
    <div className="page">
      <div>
        <Card title="New post">
          <PostComposer authorId={userId} onCreated={() => setRefreshKey((k) => k + 1)} />
        </Card>
        <Card title={tab === 'feed' ? 'Feed' : 'My posts'}>
          <div className="feed-tabs">
            <button
              type="button"
              className={clsx('feed-tabs__btn', tab === 'feed' && 'feed-tabs__btn--active')}
              onClick={() => setTab('feed')}
            >
              Everyone
            </button>
            <button
              type="button"
              className={clsx('feed-tabs__btn', tab === 'mine' && 'feed-tabs__btn--active')}
              onClick={() => setTab('mine')}
            >
              My posts
            </button>
          </div>
          <PostList currentUserId={userId} refreshKey={refreshKey} mode={tab === 'mine' ? 'mine' : 'feed'} />
        </Card>
      </div>
    </div>
  );
}
