'use strict';

document.addEventListener('DOMContentLoaded', async () => {
  await loadApplications();
  await loadStats();
});

async function requestJSON(url, options = {}) {
  const response = await fetch(url, {
    credentials: 'same-origin',
    ...options,
    headers: {
      ...(options.body ? { 'Content-Type': 'application/json' } : {}),
      ...(options.headers || {}),
    },
  });

  const text = await response.text();

  let result;
  try {
    result = JSON.parse(text);
  } catch {
    throw new Error(text);
  }

  if (!response.ok) {
    throw new Error(result.message || 'Ошибка запроса');
  }

  return result;
}

async function loadApplications() {
  const result = await requestJSON('/admin/applications');
  const applications = result.data || [];

  const container = document.getElementById('applicationsTable');

  if (applications.length === 0) {
    container.innerHTML = '<p>Заявок нет</p>';
    return;
  }

  container.innerHTML = `
    <table border="1" cellpadding="8" cellspacing="0" style="width:100%; border-collapse:collapse;">
      <thead>
        <tr>
          <th>ID</th>
          <th>ФИО</th>
          <th>Телефон</th>
          <th>Email</th>
          <th>Дата рождения</th>
          <th>Пол</th>
          <th>Биография</th>
          <th>Языки</th>
          <th>Действия</th>
        </tr>
      </thead>
      <tbody>
        ${applications.map(app => renderRow(app)).join('')}
      </tbody>
    </table>
  `;
}

function renderRow(app) {
  return `
    <tr id="row-${app.id}">
      <td>${app.id}</td>
      <td><input value="${escapeHTML(app.name || '')}" id="name-${app.id}"></td>
      <td><input value="${escapeHTML(app.phone || '')}" id="phone-${app.id}"></td>
      <td><input value="${escapeHTML(app.email || '')}" id="email-${app.id}"></td>
      <td><input type="date" value="${escapeHTML(app.birth_date || '')}" id="birth-${app.id}"></td>
      <td>
        <select id="gender-${app.id}">
          <option value="male" ${app.gender === 'male' ? 'selected' : ''}>male</option>
          <option value="female" ${app.gender === 'female' ? 'selected' : ''}>female</option>
        </select>
      </td>
      <td><textarea id="bio-${app.id}">${escapeHTML(app.biography || '')}</textarea></td>
      <td>
        <select multiple id="languages-${app.id}">
          ${renderLanguages(app.languages || [])}
        </select>
      </td>
      <td>
        <button onclick="updateApplication(${app.id})">Сохранить</button>
        <button onclick="deleteApplication(${app.id})">Удалить</button>
      </td>
    </tr>
  `;
}

function renderLanguages(selected = []) {
  const languages = [
    { id: 1, name: 'Pascal' },
    { id: 2, name: 'C' },
    { id: 3, name: 'C++' },
    { id: 4, name: 'JavaScript' },
    { id: 5, name: 'PHP' },
    { id: 6, name: 'Python' },
    { id: 7, name: 'Java' },
    { id: 8, name: 'Haskell' },
    { id: 9, name: 'Clojure' },
    { id: 10, name: 'Prolog' },
    { id: 11, name: 'Scala' },
    { id: 12, name: 'Go' },
  ];

  return languages.map(lang => `
    <option
      value="${lang.id}"
      ${selected.includes(lang.id) ? 'selected' : ''}
    >
      ${lang.name}
    </option>
  `).join('');
}

async function updateApplication(id) {
  const data = {
    name: document.getElementById(`name-${id}`).value.trim(),
    phone: document.getElementById(`phone-${id}`).value.trim(),
    email: document.getElementById(`email-${id}`).value.trim(),
    birth_date: document.getElementById(`birth-${id}`).value,
    gender: document.getElementById(`gender-${id}`).value,
    biography: document.getElementById(`bio-${id}`).value.trim(),
    agreement: true,
    languages: Array.from(
  document.getElementById(`languages-${id}`).selectedOptions
).map(option => Number(option.value)),
  };

  const result = await requestJSON(`/admin/applications/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });

  showMessage(result.message || 'Заявка обновлена');
  await loadApplications();
  await loadStats();
}

async function deleteApplication(id) {
  if (!confirm('Удалить заявку?')) return;

  const result = await requestJSON(`/admin/applications/${id}`, {
    method: 'DELETE',
  });

  showMessage(result.message || 'Заявка удалена');
  await loadApplications();
  await loadStats();
}

async function loadStats() {
  const result = await requestJSON('/admin/stats');
  const stats = result.data || {};
  const languages = stats.languages || [];

  document.getElementById('statsTable').innerHTML = `
    <p>Всего заявок: <b>${stats.total_applications}</b></p>
    <p>Всего пользователей: <b>${stats.total_users}</b></p>
    <p>Всего сессий: <b>${stats.total_sessions}</b></p>

    <table border="1" cellpadding="8" cellspacing="0" style="border-collapse:collapse;">
      <thead>
        <tr>
          <th>Язык</th>
          <th>Количество</th>
        </tr>
      </thead>
      <tbody>
        ${languages.map(item => `
          <tr>
            <td>${escapeHTML(item.language)}</td>
            <td>${item.count}</td>
          </tr>
        `).join('')}
      </tbody>
    </table>
  `;
}

function showMessage(message) {
  const block = document.getElementById('adminMessage');
  block.textContent = message;
}

function escapeHTML(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}