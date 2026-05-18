'use client';

import Link from 'next/link';
import { useAuth } from '@/app/providers/AuthProvider';
import { routes } from '@/app/router/routes';
import { Card } from '@/shared/ui/Card';

export function HomePage() {
  const { isAuthenticated } = useAuth();

  return (
    <div className="page page--center">
      <Card title="MAXat">
        <p className="lede">
          Feature-based Next.js UI for your microservices. It calls a REST-shaped API gateway — map these routes in
          your gateway to the existing gRPC services (auth, users, content, chat).
        </p>
        {isAuthenticated ? (
          <Link className="text-link" href={routes.feed}>
            Go to feed →
          </Link>
        ) : (
          <p>
            <Link className="text-link" href={routes.login}>
              Sign in
            </Link>{' '}
            or{' '}
            <Link className="text-link" href={routes.register}>
              create an account
            </Link>
            .
          </p>
        )}
      </Card>
    </div>
  );
}
