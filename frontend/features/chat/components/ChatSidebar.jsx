'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/app/providers/AuthProvider';
import { getChats } from '@/features/chat/api/chatApi';
import { fetchPeerDisplayNames, formatChatTitle } from '@/features/chat/lib/chatTitles';
import { routes } from '@/app/router/routes';
import { clsx } from '@/shared/lib/clsx';
import { Spinner } from '@/shared/ui/Spinner';

export function ChatSidebar() {
  const pathname = usePathname();
  const { userId, ready } = useAuth();
  const [state, setState] = useState({
    loading: true,
    error: '',
    chats: [],
    nameByUserId: {},
  });

  useEffect(() => {
    if (!ready) return;

    let cancelled = false;
    (async () => {
      try {
        const data = await getChats();
        const chats = data.chats ?? [];
        const nameByUserId = await fetchPeerDisplayNames(chats, userId);
        if (!cancelled) {
          setState({ loading: false, error: '', chats, nameByUserId });
        }
      } catch (e) {
        if (!cancelled) {
          setState({
            loading: false,
            error: e instanceof Error ? e.message : 'Failed to load chats',
            chats: [],
            nameByUserId: {},
          });
        }
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [userId, ready]);

  if (!ready || state.loading) return <Spinner />;
  if (state.error) return <p className="muted small">{state.error}</p>;

  if (state.chats.length === 0) {
    return <p className="muted small">No conversations yet. Start one from another screen once the API is wired.</p>;
  }

  return (
    <nav className="chat-sidebar">
      <ul>
        {state.chats.map((c) => {
          const href = routes.chatThread(c.id);
          const active = pathname === href;
          const title = formatChatTitle(c, state.nameByUserId, userId);
          return (
            <li key={c.id}>
              <Link href={href} className={clsx('chat-sidebar__link', active && 'active')}>
                <span className="chat-sidebar__title">{title}</span>
                <span className="chat-sidebar__preview small muted">
                  {c.last_message?.content?.slice(0, 80) ?? '—'}
                </span>
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
