import { getUser } from '@/features/users/api/usersApi';

/**
 * @param {string} userId
 * @returns {Promise<string>}
 */
export async function fetchUserDisplayName(userId) {
  if (!userId) return 'Unknown';
  try {
    const u = await getUser(userId);
    const label = `${u.first_name ?? ''} ${u.last_name ?? ''}`.trim();
    return label || userId.slice(0, 8);
  } catch {
    return userId.slice(0, 8);
  }
}

/**
 * @param {string[]} userIds
 * @returns {Promise<Record<string, string>>}
 */
export async function fetchDisplayNames(userIds) {
  /** @type {Record<string, string>} */
  const nameByUserId = {};
  const unique = [...new Set(userIds.filter(Boolean))];
  await Promise.all(
    unique.map(async (id) => {
      nameByUserId[id] = await fetchUserDisplayName(id);
    }),
  );
  return nameByUserId;
}
