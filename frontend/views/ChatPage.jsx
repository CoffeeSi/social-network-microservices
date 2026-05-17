'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { useAuth } from '@/app/providers/AuthProvider';
import { getChats } from '@/features/chat/api/chatApi';
import { fetchPeerDisplayNames, formatChatTitle } from '@/features/chat/lib/chatTitles';
import { ChatSidebar, ChatThread } from '@/features/chat';
import { Card } from '@/shared/ui/Card';

export function ChatPage() {
  const params = useParams();
  const raw = params?.chatId;
  const chatId = Array.isArray(raw) ? raw[0] : raw;
  const { userId, ready } = useAuth();
  const [mainTitle, setMainTitle] = useState('Messages');

  useEffect(() => {
    if (!ready || !chatId) {
      setMainTitle('Messages');
      return;
    }
    let cancelled = false;
    (async () => {
      try {
        const data = await getChats();
        const chat = data.chats?.find((c) => c.id === chatId);
        if (!chat || cancelled) return;
        const names = await fetchPeerDisplayNames([chat], userId);
        if (cancelled) return;
        setMainTitle(formatChatTitle(chat, names, userId));
      } catch {
        if (!cancelled) setMainTitle(`Chat ${String(chatId).slice(0, 8)}…`);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [chatId, userId, ready]);

  return (
    <div className="page page--chat">
      <Card title="Conversations" className="card--chat-list">
        <ChatSidebar />
      </Card>
      <Card title={mainTitle} className="card--chat-main">
        <ChatThread chatId={chatId} />
      </Card>
    </div>
  );
}
