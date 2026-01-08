// Pulsar Documentation - Main JavaScript

(function() {
  'use strict';

  // Theme Toggle
  function initTheme() {
    const themeToggle = document.querySelector('.theme-toggle');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const savedTheme = localStorage.getItem('theme');

    // Set initial theme
    if (savedTheme) {
      document.documentElement.setAttribute('data-theme', savedTheme);
    } else if (!prefersDark) {
      document.documentElement.setAttribute('data-theme', 'light');
    }

    // Update icon
    updateThemeIcon();

    if (themeToggle) {
      themeToggle.addEventListener('click', toggleTheme);
    }
  }

  function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';

    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
    updateThemeIcon();
  }

  function updateThemeIcon() {
    const themeToggle = document.querySelector('.theme-toggle');
    if (!themeToggle) return;

    const isDark = document.documentElement.getAttribute('data-theme') !== 'light';
    themeToggle.innerHTML = isDark
      ? '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>'
      : '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>';
  }

  // Mobile Menu Toggle
  function initMobileMenu() {
    const menuBtn = document.querySelector('.mobile-menu-btn');
    const navMenu = document.querySelector('.navbar-nav');

    if (menuBtn && navMenu) {
      menuBtn.addEventListener('click', function() {
        navMenu.classList.toggle('open');
        const isOpen = navMenu.classList.contains('open');
        menuBtn.setAttribute('aria-expanded', isOpen);
        menuBtn.innerHTML = isOpen
          ? '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>'
          : '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 12h18M3 6h18M3 18h18"/></svg>';
      });
    }
  }

  // Docs Sidebar Toggle (Mobile)
  function initDocsSidebar() {
    const sidebarToggle = document.querySelector('.sidebar-toggle');
    const sidebar = document.querySelector('.docs-sidebar');

    if (sidebarToggle && sidebar) {
      sidebarToggle.addEventListener('click', function() {
        sidebar.classList.toggle('open');
      });

      // Close sidebar when clicking outside
      document.addEventListener('click', function(e) {
        if (sidebar.classList.contains('open') &&
            !sidebar.contains(e.target) &&
            !sidebarToggle.contains(e.target)) {
          sidebar.classList.remove('open');
        }
      });
    }

    // Mark active link
    const currentPath = window.location.pathname;
    const sidebarLinks = document.querySelectorAll('.docs-sidebar-nav a');
    sidebarLinks.forEach(link => {
      if (link.getAttribute('href') === currentPath ||
          currentPath.endsWith(link.getAttribute('href'))) {
        link.classList.add('active');
      }
    });
  }

  // Copy to Clipboard
  function initCopyButtons() {
    const copyButtons = document.querySelectorAll('.copy-btn');

    copyButtons.forEach(btn => {
      btn.addEventListener('click', async function() {
        const codeBlock = this.closest('.code-block');
        const code = codeBlock.querySelector('code').textContent;

        try {
          await navigator.clipboard.writeText(code);
          const originalText = this.textContent;
          this.textContent = 'Copied!';
          this.style.background = 'var(--success)';

          setTimeout(() => {
            this.textContent = originalText;
            this.style.background = '';
          }, 2000);
        } catch (err) {
          console.error('Failed to copy:', err);
        }
      });
    });
  }

  // Smooth Scroll for Anchor Links
  function initSmoothScroll() {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
      anchor.addEventListener('click', function(e) {
        const targetId = this.getAttribute('href');
        if (targetId === '#') return;

        const target = document.querySelector(targetId);
        if (target) {
          e.preventDefault();
          const navHeight = document.querySelector('.navbar')?.offsetHeight || 0;
          const targetPosition = target.getBoundingClientRect().top + window.pageYOffset - navHeight - 20;

          window.scrollTo({
            top: targetPosition,
            behavior: 'smooth'
          });
        }
      });
    });
  }

  // Animate elements on scroll
  function initScrollAnimations() {
    const animatedElements = document.querySelectorAll('.feature-card, .step, .tech-item, .docs-card');

    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          entry.target.classList.add('animate-fade-in');
          observer.unobserve(entry.target);
        }
      });
    }, {
      threshold: 0.1,
      rootMargin: '0px 0px -50px 0px'
    });

    animatedElements.forEach(el => {
      el.style.opacity = '0';
      observer.observe(el);
    });
  }

  // Table of Contents Generator
  function generateTOC() {
    const tocContainer = document.querySelector('.toc');
    const content = document.querySelector('.docs-content');

    if (!tocContainer || !content) return;

    const headings = content.querySelectorAll('h2, h3');
    if (headings.length === 0) return;

    const toc = document.createElement('ul');
    toc.className = 'toc-list';

    headings.forEach((heading, index) => {
      // Add ID if not present
      if (!heading.id) {
        heading.id = `heading-${index}`;
      }

      const li = document.createElement('li');
      li.className = heading.tagName === 'H3' ? 'toc-item toc-h3' : 'toc-item';

      const a = document.createElement('a');
      a.href = `#${heading.id}`;
      a.textContent = heading.textContent;

      li.appendChild(a);
      toc.appendChild(li);
    });

    tocContainer.appendChild(toc);
  }

  // Initialize on DOM Ready
  document.addEventListener('DOMContentLoaded', function() {
    initTheme();
    initMobileMenu();
    initDocsSidebar();
    initCopyButtons();
    initSmoothScroll();
    initScrollAnimations();
    generateTOC();
  });
})();
