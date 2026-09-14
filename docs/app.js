/**
 * md — Landing Page Interactive Features
 * Theme toggle, command copy, install tabs, terminal preview, and lightbox
 */

document.addEventListener('DOMContentLoaded', () => {
  initTheme();
  initInstallTabs();
  initCopyButtons();
  initTerminalTabs();
  initLightbox();
});

/* --------------------------------------------------------------------------
   Theme Toggle & Preference Sync
   -------------------------------------------------------------------------- */
function initTheme() {
  const themeToggleBtn = document.getElementById('theme-toggle');
  const sunIcon = document.getElementById('sun-icon');
  const moonIcon = document.getElementById('moon-icon');

  function updateIcons(isDark) {
    if (sunIcon && moonIcon) {
      if (isDark) {
        sunIcon.style.display = 'block';
        moonIcon.style.display = 'none';
      } else {
        sunIcon.style.display = 'none';
        moonIcon.style.display = 'block';
      }
    }
  }

  function getEffectiveTheme() {
    const saved = localStorage.getItem('color-scheme');
    if (saved === 'dark' || saved === 'light') return saved;
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }

  // Initial icon state
  updateIcons(getEffectiveTheme() === 'dark');

  // React to OS changes if user hasn't explicitly set a preference
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    if (!localStorage.getItem('color-scheme')) {
      updateIcons(e.matches);
    }
  });

  if (themeToggleBtn) {
    themeToggleBtn.addEventListener('click', () => {
      const current = getEffectiveTheme();
      const next = current === 'dark' ? 'light' : 'dark';
      document.documentElement.setAttribute('data-theme', next);
      localStorage.setItem('color-scheme', next);
      updateIcons(next === 'dark');
    });
  }
}

/* --------------------------------------------------------------------------
   Install Command Tabs
   -------------------------------------------------------------------------- */
function initInstallTabs() {
  const tabs = document.querySelectorAll('.install-tabs .tab-btn');
  const commandEl = document.getElementById('install-cmd-text');

  const commands = {
    homebrew: 'brew install smford/tap/md',
    go: 'go install github.com/smford/md/cmd/md@latest',
    source: 'git clone https://github.com/smford/md.git && cd md && make build',
    binary: 'curl -fsSL https://raw.githubusercontent.com/smford/md/main/install.sh | sh'
  };

  tabs.forEach((tab) => {
    tab.addEventListener('click', () => {
      tabs.forEach((t) => t.classList.remove('active'));
      tab.classList.add('active');

      const target = tab.getAttribute('data-target');
      if (commandEl && commands[target]) {
        commandEl.textContent = commands[target];
      }
    });
  });
}

/* --------------------------------------------------------------------------
   Copy to Clipboard
   -------------------------------------------------------------------------- */
function initCopyButtons() {
  document.querySelectorAll('.btn-copy').forEach((btn) => {
    btn.addEventListener('click', async () => {
      const targetId = btn.getAttribute('data-copy-target');
      let textToCopy = '';

      if (targetId) {
        const targetEl = document.getElementById(targetId);
        if (targetEl) {
          textToCopy = targetEl.textContent.trim();
        }
      } else {
        const codeEl = btn.closest('.code-box')?.querySelector('code') || btn.previousElementSibling;
        textToCopy = codeEl ? codeEl.textContent.trim() : '';
      }

      if (!textToCopy) return;

      try {
        await navigator.clipboard.writeText(textToCopy);
        const originalHtml = btn.innerHTML;
        btn.classList.add('copied');
        btn.innerHTML = `
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          <span>Copied!</span>
        `;

        setTimeout(() => {
          btn.classList.remove('copied');
          btn.innerHTML = originalHtml;
        }, 2000);
      } catch (err) {
        console.error('Failed to copy to clipboard', err);
      }
    });
  });
}

/* --------------------------------------------------------------------------
   Terminal Preview Window Tabs
   -------------------------------------------------------------------------- */
function initTerminalTabs() {
  const terminalTabs = document.querySelectorAll('.terminal-tab-btn');
  const views = document.querySelectorAll('.terminal-content-view');

  terminalTabs.forEach((tab) => {
    tab.addEventListener('click', () => {
      terminalTabs.forEach((t) => t.classList.remove('active'));
      views.forEach((v) => v.classList.remove('active'));

      tab.classList.add('active');
      const viewId = tab.getAttribute('data-view');
      const targetView = document.getElementById(viewId);
      if (targetView) {
        targetView.classList.add('active');
      }
    });
  });
}

/* --------------------------------------------------------------------------
   Screenshot Lightbox Modal
   -------------------------------------------------------------------------- */
function initLightbox() {
  const lightbox = document.getElementById('lightbox');
  const lightboxImg = document.getElementById('lightbox-img');
  const closeBtn = document.querySelector('.lightbox-close');

  if (!lightbox || !lightboxImg) return;

  document.querySelectorAll('.gallery-card').forEach((card) => {
    card.addEventListener('click', () => {
      const img = card.querySelector('img');
      if (img) {
        lightboxImg.src = img.src;
        lightboxImg.alt = img.alt || 'Screenshot Preview';
        lightbox.classList.add('active');
        document.body.style.overflow = 'hidden';
      }
    });
  });

  function closeLightbox() {
    lightbox.classList.remove('active');
    document.body.style.overflow = '';
  }

  if (closeBtn) {
    closeBtn.addEventListener('click', closeLightbox);
  }

  lightbox.addEventListener('click', (e) => {
    if (e.target === lightbox) {
      closeLightbox();
    }
  });

  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && lightbox.classList.contains('active')) {
      closeLightbox();
    }
  });
}
