/**
 * Central route paths for the App Router and navigation.
 * Extend this as you add gateway-backed screens.
 */
export const routes = {
  home: '/',
  login: '/login',
  register: '/register',
  confirmEmail: '/confirm-email',
  feed: '/feed',
  users: '/users',
  profile: '/profile',
  userProfile: (userId) => `/users/${encodeURIComponent(userId)}`,
  chat: '/chat',
  chatThread: (chatId) => `/chat/${encodeURIComponent(chatId)}`,
};
