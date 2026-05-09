'use strict';

const API = {
  applications: '/api/applications/',
  login: '/api/login',
  me: '/api/me',
  logout: '/api/logout',
};

const state = {
  isEditMode: false,
  applicationId: null,
};

document.addEventListener('DOMContentLoaded', () => {
  initLoginForm();
  initApplicationForm();
});

function $(id) {
  return document.getElementById(id);
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
  console.log('HTTP', response.status, url, text);

  let result = null;

  if (text.trim() !== '') {
    try {
      result = JSON.parse(text);
    } catch {
      throw new Error(text || `Ошибка HTTP ${response.status}`);
    }
  }

  if (!response.ok) {
    throw new Error(result?.message || text || `Ошибка HTTP ${response.status}`);
  }

  return result;
}

function initLoginForm() {
  const loginForm = $('loginForm');
  if (!loginForm) return;

  loginForm.addEventListener('submit', async (event) => {
    event.preventDefault();

    const login = $('login')?.value.trim();
    const password = $('password')?.value.trim();

    if (!login || !password) {
      showMessage('loginMessage', 'Введите логин и пароль', 'error');
      return;
    }

    try {
      const result = await requestJSON(API.login, {
        method: 'POST',
        body: JSON.stringify({ login, password }),
      });

      showMessage('loginMessage', result.message || 'Авторизация успешна', 'success');

      const applicationId = Number(result?.data?.application_id);

      if (applicationId) {
        state.isEditMode = true;
        state.applicationId = applicationId;

        await loadApplicationForEdit(applicationId);

        showMessage(
          'loginMessage',
          'Авторизация успешна. Данные загружены в форму, теперь их можно изменить.',
          'success'
        );
      
        document.getElementById('contact')?.scrollIntoView({ behavior: 'smooth' });
      }

      console.log('LOGIN RESULT:', result);
    } catch (error) {
      showMessage('loginMessage', error.message, 'error');
    }
  });
}

function initApplicationForm() {
  const contactForm = $('contactForm');
  if (!contactForm) return;

  contactForm.addEventListener('submit', async (event) => {
    event.preventDefault();

    clearErrors();

    const application = collectApplicationData();
    const validationError = validateApplication(application);

    if (validationError) {
      showMessage('formMessage', validationError, 'error');
      return;
    }

    setLoading(true);

    try {
      let result;

      if (state.isEditMode && state.applicationId) {
        result = await requestJSON(`${API.applications}${state.applicationId}`, {
          method: 'PUT',
          body: JSON.stringify(application),
        });

        showMessage('formMessage', result.message || 'Данные успешно обновлены', 'success');
      } else {
        result = await requestJSON(API.applications, {
          method: 'POST',
          body: JSON.stringify(application),
        });

        const login = result?.data?.login || '';
        const password = result?.data?.password || '';
        const profileUrl = result?.data?.profile_url || '/';

        showMessage(
          'formMessage',
          `<strong>${escapeHTML(result.message || 'Заявка успешно отправлена')}</strong><br><br>` +
            `Сохраните данные для входа:<br>` +
            `<b>Логин:</b> ${escapeHTML(login)}<br>` +
            `<b>Пароль:</b> ${escapeHTML(password)}<br>` +
            `<b>Вход:</b> ${escapeHTML(profileUrl)}`,
          'success'
        );
      }

      console.log('FORM RESULT:', result);
    } catch (error) {
      showMessage('formMessage', error.message, 'error');
    } finally {
      setLoading(false);
    }
  });
}

function collectApplicationData() {
  return {
    name: $('name')?.value.trim() || '',
    phone: $('phone')?.value.trim() || '',
    email: $('email')?.value.trim() || '',
    birth_date: $('birth_date')?.value || '',
    gender: document.querySelector('input[name="gender"]:checked')?.value || '',
    biography: $('biography')?.value.trim() || '',
    agreement: Boolean($('agreement')?.checked),
    languages: Array.from($('languages')?.selectedOptions || []).map((option) =>
      Number(option.value)
    ),
  };
}

function validateApplication(application) {
  let hasError = false;

  if (!application.name) {
    setError('nameError', 'Введите ФИО');
    hasError = true;
  }

  if (!application.email || !application.email.includes('@')) {
    setError('emailError', 'Введите корректный email');
    hasError = true;
  }

  if (!application.phone) {
    setError('phoneError', 'Введите телефон');
    hasError = true;
  }

  if (!application.birth_date) {
    setError('birthDateError', 'Введите дату рождения');
    hasError = true;
  }

  if (!application.gender) {
    setError('genderError', 'Выберите пол');
    hasError = true;
  }

  if (application.languages.length === 0) {
    setError('languagesError', 'Выберите хотя бы один язык');
    hasError = true;
  }

  if (!application.biography) {
    setError('biographyError', 'Введите биографию');
    hasError = true;
  }

  if (!application.agreement) {
    setError('agreementError', 'Необходимо согласие');
    hasError = true;
  }

  return hasError ? 'Проверьте поля формы' : null;
}

function clearErrors() {
  document.querySelectorAll('.error-message').forEach((element) => {
    element.textContent = '';
  });
}

function setError(id, message) {
  const element = $(id);
  if (element) element.textContent = message;
}

function showMessage(id, message, type) {
  const element = $(id);
  if (!element) return;

  element.innerHTML = message;
  element.className = `form-message ${type}`;
}

function setLoading(isLoading) {
  const submitBtn = $('submitBtn');
  const spinner = $('spinner');

  if (submitBtn) submitBtn.disabled = isLoading;
  if (spinner) spinner.classList.toggle('hidden', !isLoading);
}

function escapeHTML(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}

async function loadApplicationForEdit(applicationId) {
  const result = await requestJSON(`${API.applications}${applicationId}`, {
    method: 'GET',
  });

  const app = result.data || result;

  fillApplicationForm(app);
}

function fillApplicationForm(app) {
  $('name').value = app.name || '';
  $('phone').value = app.phone || '';
  $('email').value = app.email || '';
  $('birth_date').value = app.birth_date || '';
  $('biography').value = app.biography || '';
  $('agreement').checked = Boolean(app.agreement);

  const genderInput = document.querySelector(
    `input[name="gender"][value="${app.gender}"]`
  );
  if (genderInput) {
    genderInput.checked = true;
  }

  const languages = app.languages || [];

  Array.from($('languages').options).forEach((option) => {
    option.selected = languages.includes(Number(option.value));
  });

  const submitText = document.querySelector('#submitBtn span');
  if (submitText) {
    submitText.textContent = 'Обновить данные';
  }
}