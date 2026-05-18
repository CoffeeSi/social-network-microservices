import { apiJson } from '@/shared/api/client';
import { endpoints } from '@/shared/api/endpoints';

/** @param {{ page_size?: number, page?: number, author_id?: string }} [params] */
export function listPosts(params = {}) {
  const qs = new URLSearchParams();
  if (params.page_size != null) qs.set('page_size', String(params.page_size));
  if (params.page != null) qs.set('page', String(params.page));
  if (params.author_id) qs.set('author_id', params.author_id);
  const q = qs.toString();
  return apiJson(`${endpoints.content.posts}${q ? `?${q}` : ''}`);
}

/** @param {{ page_size?: number, page?: number }} [params] */
export function getMyPosts(params = {}) {
  const qs = new URLSearchParams();
  if (params.page_size != null) qs.set('page_size', String(params.page_size));
  if (params.page != null) qs.set('page', String(params.page));
  const q = qs.toString();
  return apiJson(`${endpoints.content.myPosts}${q ? `?${q}` : ''}`);
}

/** @param {string} id */
export function getPost(id) {
  return apiJson(endpoints.content.post(id));
}

/** @param {string} id */
export function getPostStats(id) {
  return apiJson(endpoints.content.postStats(id));
}

/** @param {{ content: string, media_urls?: string[] }} payload */
export function createPost(payload) {
  return apiJson(endpoints.content.posts, { method: 'POST', body: payload });
}

/** @param {string} id @param {{ content?: string, media_urls?: string[] }} payload */
export function updatePost(id, payload) {
  return apiJson(endpoints.content.post(id), { method: 'PATCH', body: payload });
}

/** @param {string} id */
export function deletePost(id) {
  return apiJson(endpoints.content.post(id), { method: 'DELETE' });
}

/** @param {string} postId @param {{ page_size?: number, page?: number }} [params] */
export function listComments(postId, params = {}) {
  const qs = new URLSearchParams();
  if (params.page_size != null) qs.set('page_size', String(params.page_size));
  if (params.page != null) qs.set('page', String(params.page));
  const q = qs.toString();
  return apiJson(`${endpoints.content.comments(postId)}${q ? `?${q}` : ''}`);
}

/** @param {string} postId @param {{ text: string }} payload */
export function createComment(postId, payload) {
  return apiJson(endpoints.content.comments(postId), { method: 'POST', body: payload });
}

/** @param {string} commentId @param {{ text: string }} payload */
export function updateComment(commentId, payload) {
  return apiJson(endpoints.content.comment(commentId), { method: 'PATCH', body: payload });
}

/** @param {string} commentId */
export function deleteComment(commentId) {
  return apiJson(endpoints.content.comment(commentId), { method: 'DELETE' });
}

/** @param {string} postId */
export function toggleLike(postId) {
  return apiJson(endpoints.content.like(postId), { method: 'POST', body: {} });
}
