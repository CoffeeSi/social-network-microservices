'use client';

import { useState } from 'react';
import { useAuth } from '@/app/providers/AuthProvider';
import { PostComposer, PostList } from '@/features/content';
import { Card } from '@/shared/ui/Card';

export function FeedPage() {
  const { userId } = useAuth();
  const [refreshKey, setRefreshKey] = useState(0);

  return (
    <div className="page">
      <div>
        <Card title="New post">
          <PostComposer authorId={userId} onCreated={() => setRefreshKey((k) => k + 1)} />
        </Card>
        <Card title="Feed">
          <PostList currentUserId={userId} refreshKey={refreshKey} />
        </Card>
      </div>
    </div>
  );
}
