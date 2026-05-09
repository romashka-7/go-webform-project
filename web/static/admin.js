'use strict';

const LANGUAGES = [
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

document.addEventListener('DOMContentLoaded', async () => {
  document.getElementById('refreshBtn')?.addEventListener('click', loadAdminData);

  await loadAdminData();
});

async function loadAdminData() {
  try {
    await Promise.all([
      loadApplications(),
      loadStats(),
    ]);
  } catch (error) {
    showMessage(error.message, 'error');
  }
}

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
    throw new Error(text || 'Сервер вернул не JSON');
  }

  if (!response.ok) {
    throw new Error(result.message || 'Ошибка запроса');
  }

  return result;
}

async function loadApplications() {
  const result = await requestJSON('/admin/applications');
  const applications = result.data || [];

  const tbody = document.getElementById('applicationsBody');

  if (!applications.length) {
    tbody.innerHTML = `
      <tr>
        <td colspan="9" class="empty-cell">Заявок пока нет</td>
      </tr>
    `;
    return;
  }

  tbody.innerHTML = applications.map(renderApplicationRow).join('');
}

function renderApplicationRow(app) {
  return `
    <tr id="row-${app.id}">
      <td><span class="badge">#${app.id}</span></td>

      <td>
        <input id="name-${app.id}" value="${escapeHTML(app.name || '')}" />
      </td>

      <td>
        <input id="phone-${app.id}" value="${escapeHTML(app.phone || '')}" />
      </td>

      <td>
        <input id="email-${app.id}" type="email" value="${escapeHTML(app.email || '')}" />
      </td>

      <td>
        <input id="birth-${app.id}" type="date" value="${escapeHTML(app.birth_date || '')}" />
      </td>

      <td>
        <select id="gender-${app.id}">
          <option value="male" ${app.gender === 'male' ? 'selected' : ''}>Мужской</option>
          <option value="female" ${app.gender === 'female' ? 'selected' : ''}>Женский</option>
        </select>
      </td>

      <td>
        <textarea id="bio-${app.id}">${escapeHTML(app.biography || '')}</textarea>
      </td>

      <td>
        <select id="languages-${app.id}" multiple>
          ${renderLanguageOptions(app.languages || [])}
        </select>
      </td>

      <td>
        <div class="actions">
          <button class="action-btn save-btn" onclick="updateApplication(${app.id})">
            Сохранить
          </button>

          <button class="action-btn delete-btn" onclick="deleteApplication(${app.id})">
            Удалить
          </button>
        </div>
      </td>
    </tr>
  `;
}

function renderLanguageOptions(selectedLanguages) {
  return LANGUAGES.map((language) => {
    const selected = selectedLanguages.includes(language.id) ? 'selected' : '';

    return `
      <option value="${language.id}" ${selected}>
        ${language.name}
      </option>
    `;
  }).join('');
}

async function updateApplication(id) {
  try {
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
      ).map((option) => Number(option.value)),
    };

    const result = await requestJSON(`/admin/applications/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });

    showMessage(result.message || 'Заявка обновлена', 'success');

    await loadAdminData();
  } catch (error) {
    showMessage(error.message, 'error');
  }
}

async function deleteApplication(id) {
  const confirmed = confirm(`Удалить заявку #${id}?`);

  if (!confirmed) return;

  try {
    const result = await requestJSON(`/admin/applications/${id}`, {
      method: 'DELETE',
    });

    showMessage(result.message || 'Заявка удалена', 'success');

    await loadAdminData();
  } catch (error) {
    showMessage(error.message, 'error');
  }
}

async function loadStats() {
  const result = await requestJSON('/admin/stats');
  const stats = result.data || {};
  const languages = stats.languages || [];

  renderStatsCards(stats);
  renderStatsTable(languages);
}

function renderStatsCards(stats) {
  const container = document.getElementById('statsCards');

  container.innerHTML = `
    <div class="stat-card">
      <span>Всего заявок</span>
      <strong>${stats.total_applications ?? 0}</strong>
    </div>

    <div class="stat-card">
      <span>Пользователей</span>
      <strong>${stats.total_users ?? 0}</strong>
    </div>

    <div class="stat-card">
      <span>Активных сессий</span>
      <strong>${stats.total_sessions ?? 0}</strong>
    </div>
  `;
}

function renderStatsTable(languages) {
  const container = document.getElementById('statsTable');

  container.innerHTML = `
    <table class="stats-table">
      <thead>
        <tr>
          <th>Язык</th>
          <th>Количество</th>
        </tr>
      </thead>
      <tbody>
        ${languages.map((item) => `
          <tr>
            <td>${escapeHTML(item.language)}</td>
            <td><span class="badge">${item.count}</span></td>
          </tr>
        `).join('')}
      </tbody>
    </table>
  `;
}

function showMessage(message, type) {
  const messageBlock = document.getElementById('adminMessage');

  messageBlock.textContent = message;
  messageBlock.className = `admin-message ${type}`;

  setTimeout(() => {
    messageBlock.className = 'admin-message hidden';
  }, 3500);
}

function escapeHTML(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}