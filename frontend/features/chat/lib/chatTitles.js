import { getUser } from '@/features/users/api/usersApi';

/**
 * Fetches display names for the "other" participant in each direct chat.
 * @param {Array<{ is_group?: boolean, participant_ids?: string[] }>} chats
 * @param {string | null} currentUserId
 * @returns {Promise<Record<string, string>>} user id → "First Last"
 */
export async function fetchPeerDisplayNames(chats, currentUserId) {
  /** @type {Record<string, string>} */
  const nameByUserId = {};
  if (!currentUserId) return nameByUserId;

  const ids = new Set();
  for (const c of chats) {
    if (c.is_group) continue;
    const other = c.participant_ids?.find((id) => id !== currentUserId);
    if (other) ids.add(other);
  }

  await Promise.all(
    [...ids].map(async (id) => {
      try {
        const u = await getUser(id);
        const label = `${u.first_name ?? ''} ${u.last_name ?? ''}`.trim();
        nameByUserId[id] = label || id.slice(0, 8);
      } catch {
        nameByUserId[id] = id.slice(0, 8);
      }
    }),
  );

  return nameByUserId;
}

/**
 * @param {{ is_group?: boolean, name?: string, participant_ids?: string[] }} chat
 * @param {Record<string, string>} nameByUserId
 * @param {string | null} currentUserId
 */
export function formatChatTitle(chat, nameByUserId, currentUserId) {
  if (!chat) return 'Messages';
  if (chat.is_group) return chat.name?.trim() || 'Group';

  if (currentUserId && chat.participant_ids?.length) {
    const other = chat.participant_ids.find((id) => id !== currentUserId);
    if (other && nameByUserId[other]) return nameByUserId[other];
  }

  return chat.name?.trim() || 'Direct message';
}
