"use client";

import { Suspense } from 'react';
import ConfirmEmailClient from './ConfirmEmail.client';

export const dynamic = 'force-dynamic';

export default function ConfirmEmailPage() {
    return (
        <Suspense fallback={<p>Загрузка...</p>}>
            <ConfirmEmailClient />
        </Suspense>
    );
}
