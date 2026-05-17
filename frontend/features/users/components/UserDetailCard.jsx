import { Card } from '@/shared/ui/Card';

export function UserDetailCard({ user }) {
  if (!user) {
    return (
      <Card title="Profile">
        <p className="muted">Select a user to view details.</p>
      </Card>
    );
  }

  return (
    <Card title={`${user.first_name} ${user.last_name}`}>
      <dl className="detail-list">
        <div>
          <dt>Email</dt>
          <dd>{user.email}</dd>
        </div>
        <div>
          <dt>Active</dt>
          <dd>{user.is_active ? 'Yes' : 'No'}</dd>
        </div>
      </dl>
    </Card>
  );
}
