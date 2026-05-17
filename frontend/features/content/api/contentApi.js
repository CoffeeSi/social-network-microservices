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

/** @param {string} id */
export function getPost(id) {
  return apiJson(endpoints.content.post(id));
}

/** @param {{ author_id: string, content: string, media_urls?: string[] }} payload */
export function createPost(payload) {
  return apiJson(endpoints.content.posts, { method: 'POST', body: payload });
}

/** @param {string} id @param {{ user_id: string, content?: string, media_urls?: string[] }} payload */
export function updatePost(id, payload) {
  return apiJson(endpoints.content.post(id), { method: 'PATCH', body: payload });
}

/** @param {string} id @param {{ user_id: string }} [query] */
export function deletePost(id, query) {
  const qs = query?.user_id ? `?user_id=${encodeURIComponent(query.user_id)}` : '';
  return apiJson(`${endpoints.content.post(id)}${qs}`, { method: 'DELETE' });
}

/** @param {string} postId @param {{ page_size?: number, page?: number }} [params] */
export function listComments(postId, params = {}) {
  const qs = new URLSearchParams();
  if (params.page_size != null) qs.set('page_size', String(params.page_size));
  if (params.page != null) qs.set('page', String(params.page));
  const q = qs.toString();
  return apiJson(`${endpoints.content.comments(postId)}${q ? `?${q}` : ''}`);
}

/** @param {string} postId @param {{ user_id: string, text: string }} payload */
export function createComment(postId, payload) {
  return apiJson(endpoints.content.comments(postId), { method: 'POST', body: payload });
}

/** @param {string} postId @param {{ user_id: string }} payload */
export function toggleLike(postId, payload) {
  return apiJson(endpoints.content.like(postId), { method: 'POST', body: payload });
}
