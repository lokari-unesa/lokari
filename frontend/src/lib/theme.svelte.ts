function createTheme() {
  let isDark = $state(false);

  function applyTheme(dark: boolean) {
    if (typeof window !== 'undefined') {
      const root = document.documentElement;
      if (dark) {
        root.classList.add('dark');
      } else {
        root.classList.remove('dark');
      }
    }
  }

  // Initialize theme on client side
  if (typeof window !== 'undefined') {
    const saved = localStorage.getItem('lokari-theme');
    if (saved) {
      isDark = saved === 'dark';
    } else {
      isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    applyTheme(isDark);
  }

  return {
    get isDark() {
      return isDark;
    },
    toggle: () => {
      isDark = !isDark;
      if (typeof window !== 'undefined') {
        localStorage.setItem('lokari-theme', isDark ? 'dark' : 'light');
        applyTheme(isDark);
      }
    }
  };
}

export const theme = createTheme();
