/**
 * Central route paths for the App Router and navigation.
 * Extend this as you add gateway-backed screens.
 */
export const routes = {
  home: '/',
  login: '/login',
  register: '/register',
  feed: '/app/feed',
  users: '/app/users',
  /** Current user — redirects to /app/users/{id from JWT} */
  profile: '/app/profile',
  userProfile: (userId) => `/app/users/${encodeURIComponent(userId)}`,
  chat: '/app/chat',
  chatThread: (chatId) => `/app/chat/${encodeURIComponent(chatId)}`,
};
