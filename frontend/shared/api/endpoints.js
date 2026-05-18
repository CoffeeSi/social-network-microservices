/**
 * REST paths aligned with api-gateway/internal/transport/http/router.go
 */
export const endpoints = {
  auth: {
    register: '/api/v1/auth/register',
    login: '/api/v1/auth/login',
    verify: '/api/v1/auth/verify',
    refresh: '/api/v1/auth/refresh',
  },
  users: {
    list: '/api/v1/users',
    one: (id) => `/api/v1/users/${encodeURIComponent(id)}`,
    create: '/api/v1/users',
    changePassword: (id) => `/api/v1/users/${encodeURIComponent(id)}/password`,
  },
  content: {
    posts: '/api/v1/posts',
    myPosts: '/api/v1/posts/me',
    post: (id) => `/api/v1/posts/${encodeURIComponent(id)}`,
    postStats: (id) => `/api/v1/posts/${encodeURIComponent(id)}/stats`,
    comments: (postId) => `/api/v1/posts/${encodeURIComponent(postId)}/comments`,
    comment: (id) => `/api/v1/comments/${encodeURIComponent(id)}`,
    like: (postId) => `/api/v1/posts/${encodeURIComponent(postId)}/like`,
  },
  chat: {
    chats: '/api/v1/chats',
    direct: '/api/v1/chats/direct',
    group: '/api/v1/chats/group',
    chat: (id) => `/api/v1/chats/${encodeURIComponent(id)}`,
    leave: (id) => `/api/v1/chats/${encodeURIComponent(id)}/leave`,
    participants: (id) => `/api/v1/chats/${encodeURIComponent(id)}/participants`,
    participant: (id, userId) =>
      `/api/v1/chats/${encodeURIComponent(id)}/participants/${encodeURIComponent(userId)}`,
    messages: (chatId) => `/api/v1/chats/${encodeURIComponent(chatId)}/messages`,
    message: (chatId, messageId) =>
      `/api/v1/chats/${encodeURIComponent(chatId)}/messages/${encodeURIComponent(messageId)}`,
    readReceipt: (chatId, messageId) =>
      `/api/v1/chats/${encodeURIComponent(chatId)}/messages/${encodeURIComponent(messageId)}/read`,
    typing: (chatId) => `/api/v1/chats/${encodeURIComponent(chatId)}/typing`,
  },
};
