import { Suspense } from 'react';
import { LoginPage } from '@/views/LoginPage';
import { Spinner } from '@/shared/ui/Spinner';

export default function Page() {
  return (
    <Suspense
      fallback={
        <div className="page page--center">
          <Spinner />
        </div>
      }
    >
      <LoginPage />
    </Suspense>
  );
}
