<script lang="ts">
  import { Search, MapPin, Route, ShieldCheck, Info, ChevronRight, Navigation2, Loader2 } from "lucide-svelte";
  import { i18n } from "$lib/i18n.svelte";
  import { isQueryAllowed, ANOMALY_MESSAGE } from "$lib/queryGuard";

  let query = $state("");
  let answered = $state(false);
  let isLoading = $state(false);
  let errorMessage = $state("");
  let results = $state<any[]>([]);

  async function performSearch() {
    if (!query.trim()) return;

    // Filter Defensif: modul bersama $lib/queryGuard (validasi sebenarnya di backend)
    if (!isQueryAllowed(query)) {
      errorMessage = ANOMALY_MESSAGE;
      results = [];
      answered = true;
      return;
    }

    isLoading = true;
    answered = false;
    errorMessage = "";
    
    try {
      const res = await fetch("/api/search", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query })
      });
      const data = await res.json();
      
      if (res.status === 429 && data.error === "KUOTA_HABIS") {
        errorMessage = i18n.t(data.message as any);
        results = [];
      } else if (data.data) {
        results = data.data;
      } else {
        results = [];
      }
    } catch (e) {
      console.error(e);
      errorMessage = "Koneksi ke peladen terputus. Silakan coba lagi nanti.";
      results = [];
    } finally {
      isLoading = false;
      answered = true;
    }
  }
</script>

<div class="mx-auto max-w-5xl px-4 lg:px-8 py-10 overflow-x-hidden">
  <!-- Hero Area without wrapper -->
  <div class="flex flex-row items-center justify-center text-center mb-10 pt-4 gap-2.5 sm:gap-4">
    <span class="inline-flex items-center justify-center w-10 h-10 sm:w-16 sm:h-16 rounded-xl sm:rounded-2xl bg-primary/10 text-primary shadow-inner shrink-0">
      <Search class="w-5 h-5 sm:w-8 sm:h-8" />
    </span>
    <h1 class="text-[1.35rem] min-[400px]:text-[1.5rem] sm:text-[2.5rem] font-display font-extrabold text-foreground tracking-tight whitespace-nowrap">
      {i18n.t('page.search.title')}
    </h1>
  </div>
  
  <!-- Big Search Input -->
  <div class="relative max-w-3xl mx-auto flex flex-col sm:flex-row items-center bg-card rounded-2xl sm:rounded-full shadow-[0_8px_30px_rgb(0,0,0,0.04)] border border-border p-2 focus-within:ring-2 focus-within:ring-primary/50 transition-all mb-12">
    <div class="flex items-center w-full px-4 py-2 sm:py-0">
      <Search class="w-6 h-6 text-muted-foreground shrink-0" />
      <input
        bind:value={query}
        onkeydown={(e) => e.key === "Enter" && performSearch()}
        class="flex-1 bg-transparent px-4 py-3 text-[1.0625rem] outline-none placeholder:text-muted-foreground/60 w-full"
        placeholder={i18n.t('page.search.placeholder')}
      />
    </div>
    <button
      onclick={performSearch}
      disabled={isLoading}
      class="w-full sm:w-auto inline-flex items-center justify-center gap-2 bg-primary hover:bg-primary-dark text-primary-foreground px-8 py-3.5 rounded-xl sm:rounded-full font-semibold transition-colors mt-2 sm:mt-0 shrink-0 disabled:opacity-70 disabled:cursor-not-allowed"
    >
      {#if isLoading}
        <Loader2 class="w-5 h-5 animate-spin" /> {i18n.t('search.status.loading')}
      {:else}
        {i18n.t('page.search.btn')} <ChevronRight class="w-4 h-4" />
      {/if}
    </button>
  </div>

  <!-- Results Dashboard (Widget Style) -->
  {#if answered}
    {#if errorMessage}
      <div class="bg-destructive/10 border border-destructive/20 rounded-2xl p-6 text-center text-destructive animate-in fade-in slide-in-from-bottom-4 shadow-sm max-w-2xl mx-auto">
        <ShieldCheck class="w-8 h-8 mx-auto mb-3" />
        <h3 class="font-bold text-lg mb-2">{i18n.t('search.status.busy')}</h3>
        <p>{errorMessage}</p>
      </div>
    {:else if results.length === 0}
      <div class="text-center py-10 text-muted-foreground animate-in fade-in slide-in-from-bottom-4">
        {i18n.t('search.status.empty')}
      </div>
    {:else}
    <div class="flex flex-col gap-8 animate-in fade-in slide-in-from-bottom-8 duration-700">
      
      <!-- Top Destination Widget (Desktop Horizontal, Hidden Mobile) -->
      <div class="hidden md:flex bg-card rounded-2xl border border-border p-6 shadow-sm items-center gap-6">
        <div class="flex flex-col min-w-50 shrink-0">
          <div class="flex items-center gap-3 text-destructive">
            <span class="grid place-items-center w-12 h-12 rounded-xl bg-destructive/10"><MapPin class="w-6 h-6" /></span>
            <h3 class="font-display font-bold text-[1.25rem] text-foreground leading-tight">{i18n.t('page.search.dest.title')}</h3>
          </div>
        </div>
        
        <div class="flex-1 flex bg-linear-to-r from-muted/50 to-muted/10 rounded-xl p-5 border border-border items-center gap-6 justify-between">
          <div>
            <h4 class="font-semibold text-[1.125rem] text-foreground mb-1">{results[0].nama_objek}</h4>
            <p class="text-[0.9375rem] text-muted-foreground line-clamp-2 pr-4">{results[0].deskripsi || i18n.t('page.search.dest.desc')}</p>
          </div>
          <div class="shrink-0 flex gap-3">
            <a href="/safe-routes?destLat={results[0].geometri?.coordinates[1]}&destLng={results[0].geometri?.coordinates[0]}&destName={encodeURIComponent(results[0].nama_objek)}" class="inline-flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl bg-primary text-primary-foreground font-semibold hover:bg-primary-dark transition-colors shadow-sm">
              <Navigation2 class="w-5 h-5" /> {i18n.t('page.search.dest.btn1')}
            </a>
          </div>
        </div>
      </div>

      <!-- Main Answer Widget (Vertical List of Cards) -->
      <div class="flex flex-col gap-5">
        {#each results as result, idx}
        <a href="/safe-routes?destLat={result.geometri?.coordinates[1]}&destLng={result.geometri?.coordinates[0]}&destName={encodeURIComponent(result.nama_objek)}" class="block bg-card rounded-2xl border {idx === 0 ? 'border-primary ring-1 ring-primary/20' : 'border-border'} p-6 sm:p-8 shadow-sm relative overflow-hidden group hover:border-primary hover:shadow-md transition-all cursor-pointer text-left">
          {#if idx === 0}
            <div class="absolute top-0 left-0 w-1.5 h-full bg-primary"></div>
          {/if}
          <div class="flex flex-col sm:flex-row sm:items-center justify-between mb-5 gap-4">
            <div class="flex items-start sm:items-center gap-3 text-primary group-hover:text-primary-dark transition-colors">
              <span class="grid place-items-center w-10 h-10 rounded-xl bg-primary/10 group-hover:bg-primary group-hover:text-primary-foreground transition-colors shrink-0 mt-0.5 sm:mt-0"><Info class="w-5 h-5" /></span>
              <h3 class="font-display font-semibold text-[1.25rem] text-foreground leading-tight">{result.nama_objek}</h3>
            </div>
            <div class="text-left sm:text-right">
              <span class="inline-block px-3 py-1 bg-muted group-hover:bg-primary/10 group-hover:text-primary text-muted-foreground text-[0.75rem] font-medium rounded-md transition-colors">
                {i18n.t('search.result.category')}: {result.kategori}
              </span>
            </div>
          </div>
          <p class="text-[1rem] sm:text-[1.0625rem] text-muted-foreground leading-relaxed">
            {result.deskripsi || i18n.t('search.result.noDesc')}
          </p>
          <div class="mt-5 pt-5 border-t border-border/50 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3 text-[0.875rem]">
            <span class="text-muted-foreground flex items-center gap-2">
              <span class="inline-block w-2 h-2 rounded-full {idx === 0 ? 'bg-primary' : 'bg-muted-foreground/40'} shrink-0"></span>
              {i18n.t('search.result.capacity')}:
            </span>
            <span class="font-bold text-foreground bg-muted px-2 py-1 rounded-md ml-4 sm:ml-0">± {result.kapasitas} {i18n.t('search.result.people')}</span>
          </div>
        </a>
        {/each}
      </div>

    </div>
    {/if}
  {/if}
</div>
