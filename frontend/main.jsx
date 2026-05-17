'use client';

import { AuthProvider } from '@/app/providers/AuthProvider';

export function Main({ children }) {
  return <AuthProvider>{children}</AuthProvider>;
}
