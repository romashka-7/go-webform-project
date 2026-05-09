// DOM готов
document.addEventListener("DOMContentLoaded", function () {
  // ========== Мобильное меню ==========
  const hamburger = document.getElementById("hamburger");
  const mobileMenu = document.getElementById("mobileMenu");
  const mobileDropdowns = document.querySelectorAll(".mobile-dropdown");

  if (hamburger) {
    hamburger.addEventListener("click", function () {
      mobileMenu.classList.toggle("active");
      hamburger.classList.toggle("active");

      // Анимация гамбургера в крестик
      const bars = document.querySelectorAll(".bar");
      if (mobileMenu.classList.contains("active")) {
        bars[0].style.transform = "rotate(45deg) translate(5px, 5px)";
        bars[1].style.opacity = "0";
        bars[2].style.transform = "rotate(-45deg) translate(7px, -6px)";
      } else {
        bars[0].style.transform = "none";
        bars[1].style.opacity = "1";
        bars[2].style.transform = "none";
      }
    });
    // ========== FAQ аккордеон ==========
    const faqItems = document.querySelectorAll(".faq-item");

    faqItems.forEach((item) => {
      const question = item.querySelector(".faq-question");

      if (question) {
        question.addEventListener("click", () => {
          // Закрываем все остальные открытые вопросы
          faqItems.forEach((otherItem) => {
            if (otherItem !== item && otherItem.classList.contains("active")) {
              otherItem.classList.remove("active");
            }
          });

          // Переключаем текущий вопрос
          item.classList.toggle("active");
        });
      }
    });
  }

  // Раскрытие подменю в мобильной версии
  mobileDropdowns.forEach((dropdown) => {
    const link = dropdown.querySelector(".mobile-link");
    const submenu = dropdown.querySelector(".mobile-submenu");

    if (link && submenu) {
      link.addEventListener("click", function (e) {
        e.preventDefault();
        submenu.classList.toggle("active");

        // Поворот стрелочки
        const icon = link.querySelector("i");
        if (icon) {
          icon.style.transform = submenu.classList.contains("active")
            ? "rotate(180deg)"
            : "rotate(0deg)";
        }
      });
    }
  });

  // ========== Слайдер ==========
  const slider = document.querySelector(".slider");
  const prevBtn = document.querySelector(".prev-btn");
  const nextBtn = document.querySelector(".next-btn");
  const dots = document.querySelectorAll(".dot");
  const slides = document.querySelectorAll(".slider-item");

  let currentSlide = 0;
  const totalSlides = slides.length;

  // Функция обновления слайдера
  function updateSlider() {
    if (slider) {
      slider.style.transform = `translateX(-${currentSlide * 100}%)`;

      // Обновление точек
      dots.forEach((dot, index) => {
        dot.classList.toggle("active", index === currentSlide);
      });

      // Обновление активного слайда
      slides.forEach((slide, index) => {
        slide.classList.toggle("active", index === currentSlide);
      });
    }
  }

  // Следующий слайд
  if (nextBtn) {
    nextBtn.addEventListener("click", () => {
      currentSlide = (currentSlide + 1) % totalSlides;
      updateSlider();
    });
  }

  // Предыдущий слайд
  if (prevBtn) {
    prevBtn.addEventListener("click", () => {
      currentSlide = (currentSlide - 1 + totalSlides) % totalSlides;
      updateSlider();
    });
  }

  // Навигация по точкам
  dots.forEach((dot) => {
    dot.addEventListener("click", function () {
      currentSlide = parseInt(this.getAttribute("data-slide"));
      updateSlider();
    });
  });

  // Автопрокрутка слайдера
  let slideInterval = setInterval(() => {
    if (slider) {
      currentSlide = (currentSlide + 1) % totalSlides;
      updateSlider();
    }
  }, 5000);

  // Остановка автопрокрутки при наведении
  if (slider) {
    slider.addEventListener("mouseenter", () => {
      clearInterval(slideInterval);
    });

    slider.addEventListener("mouseleave", () => {
      slideInterval = setInterval(() => {
        currentSlide = (currentSlide + 1) % totalSlides;
        updateSlider();
      }, 5000);
    });
  }

  // ========== Форма ==========
  const contactForm = document.getElementById("contactForm");
  const submitBtn = document.getElementById("submitBtn");
  const spinner = document.getElementById("spinner");
  const formMessage = document.getElementById("formMessage");

  // Валидация формы
  function validateForm() {
    let isValid = true;
    document.querySelectorAll(".error-message").forEach((el) => {
      el.textContent = "";
    });

    const name = document.getElementById("name");
    if (!name.value.trim()) {
      document.getElementById("nameError").textContent = "Введите имя";
      isValid = false;
    }

    const email = document.getElementById("email");
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!email.value.trim()) {
      document.getElementById("emailError").textContent = "Введите email";
      isValid = false;
    } else if (!emailPattern.test(email.value)) {
      document.getElementById("emailError").textContent =
        "Введите корректный email";
      isValid = false;
    }

    const message = document.getElementById("message");
    if (!message.value.trim()) {
      document.getElementById("messageError").textContent = "Введите сообщение";
      isValid = false;
    }

    const agreement = document.getElementById("agreement");
    if (!agreement.checked) {
      document.getElementById("agreementError").textContent =
        "Необходимо согласие";
      isValid = false;
    }

    return isValid;
  }

  // Обработка отправки формы - REAL FORM CARRY
 if (contactForm) {
  contactForm.addEventListener("submit", async function (event) {

    event.preventDefault();

    if (!validateForm()) {
      return;
    }

    submitBtn.disabled = true;
    spinner.classList.remove("hidden");

    const data = {
  name: document.getElementById("name").value,
  phone: document.getElementById("phone").value,
  email: document.getElementById("email").value,
  birth_date: document.getElementById("birth_date").value,
  gender: document.querySelector('input[name="gender"]:checked')?.value,
  biography: document.getElementById("biography").value,
  agreement: document.getElementById("agreement").checked,
  languages: Array.from(
    document.getElementById("languages").selectedOptions
  ).map((option) => Number(option.value)),
};

    try {

      const response = await fetch("/api/applications/", {
        method: "POST",

        headers: {
          "Content-Type": "application/json",
        },

        body: JSON.stringify(data),
      });

      const result = await response.json();

      if (!response.ok) {
        throw new Error(result.message || "Ошибка отправки");
      }

      formMessage.textContent =
        "Заявка успешно отправлена";

      formMessage.className =
        "form-message success";

      console.log(result);

      contactForm.reset();

    } catch (error) {

      formMessage.textContent =
        error.message;

      formMessage.className =
        "form-message error";

      console.error(error);

    } finally {

      submitBtn.disabled = false;

      spinner.classList.add("hidden");
    }
  });
}

  // ========== Плавная прокрутка для якорных ссылок ==========
  document.querySelectorAll('a[href^="#"]').forEach((anchor) => {
    anchor.addEventListener("click", function (e) {
      const href = this.getAttribute("href");

      // Игнорируем пустые ссылки и ссылки только на #
      if (href === "#" || href === "") return;

      e.preventDefault();
      const targetElement = document.querySelector(href);

      if (targetElement) {
        // Закрываем мобильное меню при клике на ссылку
        if (mobileMenu.classList.contains("active")) {
          mobileMenu.classList.remove("active");
          hamburger.classList.remove("active");

          const bars = document.querySelectorAll(".bar");
          bars[0].style.transform = "none";
          bars[1].style.opacity = "1";
          bars[2].style.transform = "none";
        }

        // Показываем скрытые секции (например, услуги)
        if (href === "#services") {
          targetElement.style.display = "block";
        }

        // Плавная прокрутка
        window.scrollTo({
          top: targetElement.offsetTop - 80,
          behavior: "smooth",
        });

        // Добавляем небольшой хак для браузеров, которые не поддерживают behavior: smooth
        if (!("scrollBehavior" in document.documentElement.style)) {
          smoothScrollPolyfill(targetElement);
        }
      }
    });
  });

  // Полифилл для плавной прокрутки (для старых браузеров)
  function smoothScrollPolyfill(targetElement) {
    const targetPosition = targetElement.offsetTop - 80;
    const startPosition = window.pageYOffset;
    const distance = targetPosition - startPosition;
    const duration = 500;
    let start = null;

    function step(timestamp) {
      if (!start) start = timestamp;
      const progress = timestamp - start;
      window.scrollTo(
        0,
        easeInOutCubic(progress, startPosition, distance, duration)
      );
      if (progress < duration) {
        window.requestAnimationFrame(step);
      }
    }

    function easeInOutCubic(t, b, c, d) {
      t /= d / 2;
      if (t < 1) return (c / 2) * t * t * t + b;
      t -= 2;
      return (c / 2) * (t * t * t + 2) + b;
    }

    window.requestAnimationFrame(step);
  }

  // ========== Фиксированная навигация при скролле ==========
  let lastScrollTop = 0;
  const navbar = document.querySelector(".navbar");
  let isProgrammaticScroll = false;
  window.addEventListener("scroll", function () {
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;

    if (scrollTop > 100) {
      navbar.style.backgroundColor = "rgba(255, 255, 255, 0.98)";
      navbar.style.boxShadow = "0 5px 15px rgba(0, 0, 0, 0.1)";
    } else {
      navbar.style.backgroundColor = "rgba(255, 255, 255, 0.95)";
      navbar.style.boxShadow = "0 2px 10px rgba(0, 0, 0, 0.1)";
    }

    // Скрываем/показываем навигацию при скролле
    if (scrollTop > lastScrollTop && scrollTop > 200) {
      // Скролл вниз
      navbar.style.transform = "translateY(-100%)";
    } else {
      // Скролл вверх
      navbar.style.transform = "translateY(0)";
    }

    lastScrollTop = scrollTop;
  });
});
