'use client';

import { RegisterForm } from '@/features/auth';
import { Card } from '@/shared/ui/Card';

export function RegisterPage() {
  return (
    <div className="page page--narrow">
      <Card title="Create account">
        <RegisterForm />
      </Card>
    </div>
  );
}
