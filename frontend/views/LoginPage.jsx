'use client';

import { useSearchParams } from 'next/navigation';
import { LoginForm } from '@/features/auth';
import { Card } from '@/shared/ui/Card';

export function LoginPage() {
  const searchParams = useSearchParams();
  const registered = searchParams.get('registered') === '1';

  return (
    <div className="page page--narrow">
      <Card title="Sign in">
        {registered ? <p className="banner-success">Account created. You can sign in now.</p> : null}
        <LoginForm />
      </Card>
    </div>
  );
}
