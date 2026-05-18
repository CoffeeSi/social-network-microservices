'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/app/providers/AuthProvider';
import { routes } from '@/app/router/routes';
import { Spinner } from '@/shared/ui/Spinner';

export default function ProfileRedirectPage() {
    const { ready, userId } = useAuth();
    const router = useRouter();

    useEffect(() => {
        if (!ready) return;
        if (userId) router.replace(routes.userProfile(userId));
        else router.replace(routes.login);
    }, [ready, userId, router]);

    return (
        <div className="page page--center">
            <Spinner />
        </div>
    );
}
