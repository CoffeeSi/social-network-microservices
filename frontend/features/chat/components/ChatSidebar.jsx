'use client';

import { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/app/providers/AuthProvider';
import { getChats } from '@/features/chat/api/chatApi';
import { CreateGroupChatForm } from '@/features/chat/components/CreateGroupChatForm';
import { fetchPeerDisplayNames, formatChatTitle } from '@/features/chat/lib/chatTitles';
import { routes } from '@/app/router/routes';
import { clsx } from '@/shared/lib/clsx';
import { Spinner } from '@/shared/ui/Spinner';

export function ChatSidebar() {
  const pathname = usePathname();
  const { userId, ready } = useAuth();
  const [reloadKey, setReloadKey] = useState(0);
  const [state, setState] = useState({
    loading: true,
    error: '',
    chats: [],
    nameByUserId: {},
  });

  const reload = useCallback(() => setReloadKey((k) => k + 1), []);

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
  }, [userId, ready, reloadKey]);

  if (!ready || state.loading) return <Spinner />;
  if (state.error) return <p className="muted small">{state.error}</p>;

  return (
    <div className="chat-sidebar">
      <CreateGroupChatForm onCreated={reload} />
      {state.chats.length === 0 ? (
        <p className="muted small">No conversations yet. Message someone from their profile or create a group.</p>
      ) : (
        <nav>
          <ul>
            {state.chats.map((c) => {
              const href = routes.chatThread(c.id);
              const active = pathname === href;
              const title = formatChatTitle(c, state.nameByUserId, userId);
              return (
                <li key={c.id}>
                  <Link href={href} className={clsx('chat-sidebar__link', active && 'active')}>
                    <span className="chat-sidebar__title">
                      {title}
                      {c.is_group ? <span className="chat-sidebar__badge">group</span> : null}
                    </span>
                    <span className="chat-sidebar__preview small muted">
                      {c.last_message?.content?.slice(0, 80) ?? '—'}
                    </span>
                  </Link>
                </li>
              );
            })}
          </ul>
        </nav>
      )}
    </div>
  );
}
