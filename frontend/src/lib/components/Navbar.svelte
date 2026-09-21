<script lang="ts">
  import { Search, Menu, X, Route, Sun, Moon } from "lucide-svelte";
  import { page } from "$app/stores";
  import { cn } from "$lib/utils";
  import { i18n } from "$lib/i18n.svelte";
  import { theme } from "$lib/theme.svelte";

  let open = $state(false);

  // Dynamic links to support instant translation
  const getNavLinks = () => [
    { label: i18n.t("nav.home"), path: "/" },
    { label: i18n.t("nav.riskMap"), path: "/danger-map" },
    { label: i18n.t("nav.evacuation"), path: "/safe-routes" },
    { label: i18n.t("nav.info"), path: "/news" },
    { label: i18n.t("nav.about"), path: "/about" },
  ];
</script>

<header class="sticky top-0 z-50 h-18 bg-card border-b border-border">
  <div class="mx-auto h-full max-w-360 px-3 sm:px-4 lg:px-8 flex items-center justify-between gap-2 sm:gap-6 relative z-50 bg-card">
    <!-- Logo -->
    <a href="/" class="flex items-center gap-1.5 sm:gap-2.5 shrink-0 min-w-0">
      <img src="/logo.webp" alt="LOKARI Logo" class="h-9 w-auto object-contain" />
      <div class="flex flex-col leading-none">
        <span class="font-display font-extrabold text-[1.0625rem] text-foreground tracking-tight">LOKARI</span>
        <span class="hidden sm:block text-[0.6875rem] font-medium text-muted-foreground mt-0.5">Desa Jarak, Kediri</span>
      </div>
    </a>

    <!-- Center nav -->
    <nav class="hidden lg:flex items-center gap-1">
      {#each getNavLinks() as link}
        {@const active = $page.url.pathname === link.path}
        <a
          href={link.path}
          class={cn(
            "relative px-3.5 py-2 rounded-lg text-[0.9375rem] font-medium transition-colors",
            active ? "text-primary" : "text-muted-foreground hover:text-foreground hover:bg-muted"
          )}
        >
          {link.label}
          {#if active}
            <span class="absolute left-3.5 right-3.5 -bottom-px h-0.5 rounded-full bg-primary"></span>
          {/if}
        </a>
      {/each}
    </nav>

    <!-- Right -->
    <div class="flex items-center gap-1 sm:gap-2.5 shrink-0">
      <!-- Language Switcher -->
      <button 
        type="button"
        onclick={() => i18n.toggle()}
        class="hidden sm:grid place-items-center h-10 px-3 rounded-xl border border-border text-sm font-semibold text-muted-foreground hover:bg-muted hover:text-foreground transition-colors uppercase"
        title="Ubah Bahasa / Change Language"
      >
        {i18n.locale}
      </button>

      <!-- Theme Switcher -->
      <button 
        type="button"
        onclick={() => theme.toggle()}
        class="hidden sm:grid place-items-center w-10 h-10 rounded-xl border border-border text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
        aria-label="Toggle theme"
      >
        {#if theme.isDark}
          <Sun class="w-5 h-5" />
        {:else}
          <Moon class="w-5 h-5" />
        {/if}
      </button>

      <!-- Search button (Icon on mobile, Text on desktop) -->
      <a
        href="/search"
        class="grid place-items-center md:flex md:items-center w-10 md:w-auto h-10 md:px-4 rounded-xl bg-primary text-primary-foreground text-[0.875rem] font-semibold hover:bg-primary-dark transition-colors"
        aria-label={i18n.t('ui.search.info')}
      >
        <Search class="w-4 h-4 md:mr-2" />
        <span class="hidden md:inline">{i18n.t('ui.search.info')}</span>
      </a>

      <!-- Mobile Menu Toggle -->
      <button
        type="button"
        onclick={() => open = !open}
        class="lg:hidden grid place-items-center w-10 h-10 rounded-xl border border-border text-foreground"
        aria-label="Menu"
      >
        {#if open}
          <X class="w-5 h-5" />
        {:else}
          <Menu class="w-5 h-5" />
        {/if}
      </button>
    </div>
  </div>

  <!-- Mobile menu dropdown -->
  {#if open}
    <!-- Overlay for click-outside -->
    <div 
      class="lg:hidden fixed inset-0 z-40 bg-black/20" 
      onclick={() => open = false}
      onkeydown={(e) => e.key === 'Escape' && (open = false)}
      role="button"
      tabindex="0"
      aria-label="Tutup menu"
    ></div>

    <div class="lg:hidden absolute left-0 right-0 top-18 bg-card border-b border-border shadow-lg z-50">
      <nav class="px-4 py-3 flex flex-col">
        {#each getNavLinks() as link}
          {@const active = $page.url.pathname === link.path}
          <a
            href={link.path}
            onclick={() => open = false}
            class={cn(
              "px-3 py-3 rounded-lg text-[0.9375rem] font-medium",
              active ? "bg-accent text-primary" : "text-foreground hover:bg-muted"
            )}
          >
            {link.label}
          </a>
        {/each}
        
        <div class="mt-4 pt-4 border-t border-border flex items-center justify-between">
          <button 
            onclick={() => i18n.toggle()}
            class="flex-1 flex items-center justify-center gap-2 h-11 rounded-xl border border-border text-sm font-semibold uppercase text-foreground"
          >
            {i18n.t('ui.language')} {i18n.locale}
          </button>
          <div class="w-2"></div>
          <button 
            onclick={() => theme.toggle()}
            class="flex-1 flex items-center justify-center gap-2 h-11 rounded-xl border border-border text-sm font-semibold text-foreground"
          >
            {#if theme.isDark}
              <Sun class="w-4 h-4" /> {i18n.t('ui.lightMode')}
            {:else}
              <Moon class="w-4 h-4" /> {i18n.t('ui.darkMode')}
            {/if}
          </button>
        </div>

        <div class="mt-4">
          <a
            href="/safe-routes"
            onclick={() => open = false}
            class="flex items-center justify-center gap-1.5 w-full h-11 px-3 rounded-xl bg-red-600 text-white text-[0.875rem] font-semibold hover:bg-red-700 transition-colors shadow-sm"
          >
            <Route class="w-4 h-4" /> {i18n.t('nav.evacuation')}
          </a>
        </div>
      </nav>
    </div>
  {/if}
</header>
