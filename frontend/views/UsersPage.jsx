'use client';

import { UserList } from '@/features/users';

export function UsersPage() {
  return (
    <div className="page">
      <h1 className="page-title">People</h1>
      <p className="small muted" style={{ marginTop: '-0.5rem', marginBottom: '1rem' }}>
        Open someone&apos;s profile to send them a message.
      </p>
      <UserList />
    </div>
  );
}
