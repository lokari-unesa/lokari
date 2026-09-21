<script lang="ts">
  import { Mountain, Waves, DoorOpen, CloudRain, AlertTriangle, Filter, X } from "lucide-svelte";
  import InfoCard from "$lib/components/InfoCard.svelte";
  import { cn } from "$lib/utils";
  import { i18n } from "$lib/i18n.svelte";

  import { onMount } from "svelte";

  const FILTERS = ["all", "volcano", "lahar", "evac", "weather", "warning"];

  const categoryToSlug: Record<string, string> = {
    "Gunung Api": "volcano",
    "Lahar": "lahar",
    "Evakuasi": "evac",
    "Cuaca": "weather",
    "Peringatan": "warning"
  };

  // Mapping kategori ke Ikon dan Warna
  const categoryMeta: Record<string, any> = {
    volcano: { icon: Mountain, accent: "bg-warning/10 text-warning" },
    lahar: { icon: Waves, accent: "bg-[#92400E]/10 text-[#92400E]" },
    evac: { icon: DoorOpen, accent: "bg-primary/10 text-primary" },
    weather: { icon: CloudRain, accent: "bg-secondary text-white" },
    warning: { icon: AlertTriangle, accent: "bg-destructive/10 text-destructive" },
  };

  let items = $state<any[]>([]);
  let isLoading = $state(true);

  onMount(async () => {
    try {
      const res = await fetch("/api/news");
      const data = await res.json();
      if (data.status === "success" && data.data) {
        items = data.data.map((n: any) => {
          const catSlug = categoryToSlug[n.category] || "volcano";
          const meta = categoryMeta[catSlug] || categoryMeta["volcano"];
          return {
            category: catSlug,
            title: n.title,
            summary: n.summary,
            date: n.created_at || i18n.t('info.date.now'),
            source: n.source,
            icon: meta.icon,
            accent: meta.accent
          };
        });
      }
    } catch (e) {
      console.error("Gagal menarik berita:", e);
    } finally {
      isLoading = false;
    }
  });

  let filter = $state("all");
  let filtered = $derived(filter === "all" ? items : items.filter((i) => i.category === filter));
  let showModal = $state(false);
  let selectedNews = $state<any>(null);
</script>

<div class="mx-auto max-w-360 px-4 lg:px-8 py-8 overflow-x-hidden">
  <div class="mb-6">
    <h1 class="font-display font-semibold text-[1.5rem] text-foreground">{i18n.t('page.info.title')}</h1>
    <p class="mt-1 text-[0.9375rem] text-muted-foreground">{i18n.t('page.info.desc')}</p>
  </div>

  <!-- Desktop Filters -->
  <div class="hidden md:flex flex-wrap gap-2 mb-6">
    {#each FILTERS as f}
      <button
        onclick={() => filter = f}
        class={cn(
          "px-4 h-9 rounded-full text-[0.875rem] font-medium border transition-colors",
          filter === f ? "bg-primary text-primary-foreground border-primary" : "bg-card text-foreground border-border hover:bg-muted"
        )}
      >
        {i18n.t(`info.filter.${f}` as any)}
      </button>
    {/each}
  </div>

  <!-- Mobile Filter Button -->
  <div class="md:hidden mb-6">
    <button onclick={() => showModal = true} class="flex items-center gap-2 h-11 px-4 rounded-xl bg-card border border-border font-semibold text-[0.875rem] text-foreground w-full justify-between hover:bg-muted transition-colors">
      <span>{i18n.t('page.info.filter')}: <span class="text-primary font-bold">{i18n.t(`info.filter.${filter}` as any)}</span></span>
      <Filter class="w-4 h-4 text-muted-foreground" />
    </button>
  </div>

  <!-- Mobile Filter Modal -->
  {#if showModal}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm md:hidden" onclick={() => showModal = false}></div>
    <div class="fixed inset-0 m-auto h-fit w-[calc(100vw-2rem)] max-w-85 z-50 md:hidden p-5 rounded-xl bg-card border border-border shadow-xl">
      <div class="flex items-center justify-between mb-4">
        <h3 class="font-display font-semibold text-[1.125rem] text-foreground">{i18n.t('page.info.filter')}</h3>
        <button onclick={() => showModal = false} class="grid place-items-center w-8 h-8 rounded-md bg-muted/50 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"><X class="w-4 h-4" /></button>
      </div>
      <div class="flex flex-col gap-2">
        {#each FILTERS as f}
          <button
            onclick={() => { filter = f; showModal = false; }}
            class={cn(
              "w-full text-left px-4 h-11 rounded-lg text-[0.9375rem] font-medium transition-colors",
              filter === f ? "bg-primary text-primary-foreground" : "bg-muted/40 text-foreground hover:bg-muted"
            )}
          >
            {i18n.t(`info.filter.${f}` as any)}
          </button>
        {/each}
      </div>
    </div>
  {/if}

  <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
    {#if isLoading}
      <div class="col-span-full py-16 text-center text-muted-foreground text-[0.9375rem]">Memuat berita cerdas AI...</div>
    {:else}
      {#each filtered as item}
        <InfoCard {...item} onclick={() => selectedNews = item} />
      {/each}
    {/if}
  </div>

  {#if !isLoading && filtered.length === 0}
    <div class="py-16 text-center text-muted-foreground text-[0.9375rem]">{i18n.t('page.info.empty')}</div>
  {/if}

  <!-- Detail Modal -->
  {#if selectedNews}
    {@const ModalIcon = selectedNews.icon}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm" onclick={() => selectedNews = null}></div>
    <div class="fixed inset-0 m-auto h-fit w-[calc(100vw-2rem)] max-w-2xl z-50 p-6 rounded-2xl bg-card border border-border shadow-xl animate-in fade-in zoom-in-95 duration-200">
      <div class="flex flex-col h-full max-h-[85vh]">
        <div class="flex items-center justify-between mb-4 shrink-0">
          <div class="flex items-center gap-3">
             <span class={cn("grid place-items-center w-10 h-10 rounded-xl", selectedNews.accent || "bg-primary/10 text-primary")}>
                <ModalIcon class="w-5 h-5" strokeWidth={2} />
             </span>
             <span class="text-[0.8125rem] font-semibold text-muted-foreground">{i18n.t(`info.filter.${selectedNews.category}` as any)}</span>
          </div>
          <button onclick={() => selectedNews = null} class="grid place-items-center w-8 h-8 rounded-md bg-muted/50 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
            <X class="w-4 h-4" />
          </button>
        </div>
        
        <div class="overflow-y-auto pr-2 custom-scrollbar pb-4">
           <h2 class="font-display font-semibold text-[1.25rem] text-foreground leading-snug mb-4">{selectedNews.title}</h2>
           <p class="text-[1rem] text-muted-foreground leading-relaxed whitespace-pre-wrap">{selectedNews.summary}</p>
        </div>

        <div class="mt-4 pt-4 border-t border-border flex items-center justify-between text-[0.875rem] text-muted-foreground shrink-0">
           <span class="font-medium text-foreground">{selectedNews.date}</span>
           <span>{i18n.t('status.source')} {selectedNews.source}</span>
        </div>
      </div>
    </div>
  {/if}
</div>
