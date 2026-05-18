/**
 * REST paths expected from an API gateway that fronts the gRPC microservices.
 * Align your gateway routes with these, or change the paths here to match.
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
    post: (id) => `/api/v1/posts/${encodeURIComponent(id)}`,
    comments: (postId) => `/api/v1/posts/${encodeURIComponent(postId)}/comments`,
    comment: (postId, commentId) =>
      `/api/v1/posts/${encodeURIComponent(postId)}/comments/${encodeURIComponent(commentId)}`,
    like: (postId) => `/api/v1/posts/${encodeURIComponent(postId)}/like`,
  },
  chat: {
    chats: '/api/v1/chats',
    direct: '/api/v1/chats/direct',
    group: '/api/v1/chats/group',
    messages: (chatId) => `/api/v1/chats/${encodeURIComponent(chatId)}/messages`,
    message: (chatId, messageId) =>
      `/api/v1/chats/${encodeURIComponent(chatId)}/messages/${encodeURIComponent(messageId)}`,
  },
};
