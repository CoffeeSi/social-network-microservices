import { apiJson } from '@/shared/api/client';
import { endpoints } from '@/shared/api/endpoints';

export function getChats() {
  return apiJson(endpoints.chat.chats);
}

/** @param {{ target_user_id: string }} payload */
export function createDirectChat(payload) {
  return apiJson(endpoints.chat.direct, { method: 'POST', body: payload });
}

/** @param {{ name: string, participant_ids: string[] }} payload */
export function createGroupChat(payload) {
  return apiJson(endpoints.chat.group, { method: 'POST', body: payload });
}

/** @param {string} chatId @param {{ name: string }} payload */
export function editGroupChat(chatId, payload) {
  return apiJson(endpoints.chat.chat(chatId), { method: 'PATCH', body: payload });
}

/** @param {string} chatId */
export function deleteGroupChat(chatId) {
  return apiJson(endpoints.chat.chat(chatId), { method: 'DELETE' });
}

/** @param {string} chatId */
export function leaveGroupChat(chatId) {
  return apiJson(endpoints.chat.leave(chatId), { method: 'POST', body: {} });
}

/** @param {string} chatId @param {{ user_id: string }} payload */
export function addParticipant(chatId, payload) {
  return apiJson(endpoints.chat.participants(chatId), { method: 'POST', body: payload });
}

/** @param {string} chatId @param {string} userId */
export function removeParticipant(chatId, userId) {
  return apiJson(endpoints.chat.participant(chatId, userId), { method: 'DELETE' });
}

/** @param {string} chatId @param {{ page?: number, page_size?: number }} [params] */
export function getMessages(chatId, params = {}) {
  const qs = new URLSearchParams();
  if (params.page != null) qs.set('page', String(params.page));
  if (params.page_size != null) qs.set('page_size', String(params.page_size));
  const q = qs.toString();
  return apiJson(`${endpoints.chat.messages(chatId)}${q ? `?${q}` : ''}`);
}

/** @param {string} chatId @param {{ content: string }} payload */
export function sendMessage(chatId, payload) {
  return apiJson(endpoints.chat.messages(chatId), { method: 'POST', body: payload });
}

/** @param {string} chatId @param {string} messageId @param {{ new_content: string }} payload */
export function editMessage(chatId, messageId, payload) {
  return apiJson(endpoints.chat.message(chatId, messageId), {
    method: 'PATCH',
    body: payload,
  });
}

/** @param {string} chatId @param {string} messageId */
export function deleteMessage(chatId, messageId) {
  return apiJson(endpoints.chat.message(chatId, messageId), { method: 'DELETE' });
}

/** @param {string} chatId @param {string} messageId */
export function sendReadReceipt(chatId, messageId) {
  return apiJson(endpoints.chat.readReceipt(chatId, messageId), { method: 'POST', body: {} });
}

/** @param {string} chatId @param {{ is_typing: boolean }} payload */
export function sendTypingStatus(chatId, payload) {
  return apiJson(endpoints.chat.typing(chatId), { method: 'POST', body: payload });
}
