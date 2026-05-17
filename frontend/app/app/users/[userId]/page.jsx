'use client';

import { useParams } from 'next/navigation';
import { UserProfilePage } from '@/views/UserProfilePage';

export default function Page() {
  const params = useParams();
  const raw = params?.userId;
  const userId = Array.isArray(raw) ? raw[0] : raw;
  return <UserProfilePage userId={userId} />;
}
