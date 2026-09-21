<script lang="ts">
  import { Search, Layers, X, Navigation, ShieldCheck, MapPin, Sparkles } from "lucide-svelte";
  import KeludMapView from "$lib/components/KeludMapView.svelte";
  import { cn } from "$lib/utils";
  import { i18n } from "$lib/i18n.svelte";

  import { onMount } from "svelte";

  let LAYERS = $derived([
    { id: "hazard", label: i18n.t('layer.hazard'), color: "#EA580C" },
    { id: "shelters", label: i18n.t('layer.shelter'), color: "#16A34A" },
    { id: "health", label: i18n.t('layer.health'), color: "#0284C7" },
    { id: "lahar", label: i18n.t('layer.lahar'), color: "#92400E" },
  ]);

  let LEGEND = $derived([
    { label: i18n.t('legend.hazard'), color: "#EA580C", type: "fill" },
    { label: i18n.t('legend.shelter'), color: "#16A34A", type: "dot" },
    { label: i18n.t('legend.health'), color: "#0284C7", type: "dot" },
    { label: i18n.t('legend.lahar'), color: "#92400E", type: "line" },
  ]);

  let active: Record<string, boolean> = $state({ hazard: true, shelters: true, health: true, lahar: false });
  let selected = $state(false); 
  let selectedData: any = $state(null);
  let rekomendasiFaskes: any = $state(null);
  let potensiList: any[] = $state([]);

  onMount(async () => {
    try {
      const res = await fetch("/api/potensi");
      const data = await res.json();
      if (data.status === "success") {
        potensiList = data.data;
      }
    } catch (e) {
      console.error("Gagal mengambil data potensi bencana:", e);
    }
  });

  function handleMarkerClick(potensi: any) {
    selectedData = potensi;
    selected = true;

    // Rekomendasi Faskes Terdekat
    rekomendasiFaskes = null;
    if (!potensi.kategori.includes("Kesehatan")) {
      const faskesList = potensiList.filter(p => p.kategori.includes("Kesehatan"));
      if (faskesList.length > 0) {
        const pCoord = JSON.parse(potensi.geometri).coordinates;
        let minDist = Infinity;
        let closest = null;
        for (const f of faskesList) {
          const fCoord = JSON.parse(f.geometri).coordinates;
          const d = Math.pow(pCoord[0] - fCoord[0], 2) + Math.pow(pCoord[1] - fCoord[1], 2);
          if (d < minDist) {
            minDist = d;
            closest = f;
          }
        }
        rekomendasiFaskes = closest;
      }
    }
  }

</script>


{#snippet LayerToggles()}
  <div class="space-y-1">
    {#each LAYERS as l}
      <button
        onclick={() => active[l.id] = !active[l.id]}
        class="w-full flex items-center gap-2.5 px-2 py-2 rounded-lg hover:bg-muted text-left transition-colors"
      >
        <span
          class={cn(
            "grid place-items-center w-5 h-5 rounded-md border shrink-0 transition-colors",
            active[l.id] ? "bg-primary border-primary" : "border-border bg-card"
          )}
        >
          {#if active[l.id]}<span class="w-2.5 h-2.5 rounded-sm bg-primary-foreground"></span>{/if}
        </span>
        <span class="w-2.5 h-2.5 rounded-full shrink-0" style="background: {l.color}"></span>
        <span class="text-[0.875rem] text-foreground">{l.label}</span>
      </button>
    {/each}
  </div>
{/snippet}

{#snippet LegendList()}
  <ul class="space-y-2.5">
    {#each LEGEND as item}
      <li class="flex items-center gap-2.5 text-[0.8125rem] text-muted-foreground">
        {#if item.type === "fill"}
          <span class="w-4 h-4 rounded" style="background: {item.color}; opacity: 0.4"></span>
        {/if}
        {#if item.type === "line"}
          <span class="w-4 h-1 rounded-full" style="background: {item.color}"></span>
        {/if}
        {#if item.type === "dot"}
          <span class="w-3 h-3 rounded-full" style="background: {item.color}"></span>
        {/if}
        {item.label}
      </li>
    {/each}
  </ul>
{/snippet}

<div class="mx-auto max-w-360 px-4 lg:px-8 py-8">
  <header class="mb-5">
    <h1 class="font-display font-semibold text-[1.5rem] text-foreground">{i18n.t('page.risk.title')}</h1>
    <p class="mt-1 text-[0.9375rem] text-muted-foreground">{i18n.t('page.risk.desc')}</p>
  </header>

  <div class="flex flex-col-reverse lg:grid lg:grid-cols-[320px_1fr] gap-6 lg:gap-4">
    <!-- Mobile & Desktop control rail -->
    <aside class="flex flex-col gap-4">

      <div class="p-4 rounded-xl bg-card border border-border">
        <div class="flex items-center gap-2 mb-3">
          <Layers class="w-4 h-4 text-primary" />
          <h3 class="font-display font-semibold text-[0.9375rem] text-foreground">{i18n.t('page.risk.layers')}</h3>
        </div>
        {@render LayerToggles()}
      </div>

      <div class="p-4 rounded-xl bg-card border border-border">
        <h3 class="font-display font-semibold text-[0.9375rem] text-foreground mb-3">{i18n.t('page.risk.legend')}</h3>
        {@render LegendList()}
      </div>
    </aside>

    <!-- Map canvas -->
    <div class="relative w-full aspect-9/16 lg:aspect-auto lg:h-160 rounded-xl overflow-hidden border border-border">
      <KeludMapView
        showHazard={active.hazard}
        showRoute={false}
        showShelters={active.shelters}
        showHealth={active.health}
        showLahar={active.lahar}
        potensiData={potensiList}
        onMarkerClick={handleMarkerClick}
        destLat={0}
        destLng={0}
      />


      <!-- Floating info popover -->
      {#if selected && selectedData}
        <!-- Mobile backdrop -->
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div class="fixed inset-0 z-40 bg-black/40 backdrop-blur-sm lg:hidden" onclick={() => selected = false}></div>

        <div class="fixed inset-0 m-auto h-fit w-[calc(100vw-2rem)] max-w-85 z-50 lg:absolute lg:inset-auto lg:top-4 lg:left-4 lg:m-0 lg:w-[320px] lg:z-10 p-5 rounded-xl bg-card border border-border shadow-[0_8px_24px_-8px_rgba(15,23,42,0.18)]">
          <div class="flex items-start justify-between">
            <div class="flex flex-col items-start gap-3">
              <span class="grid place-items-center w-10 h-10 rounded-lg bg-primary/10 text-primary"><ShieldCheck class="w-5 h-5" /></span>
              <div>
                <h3 class="font-display font-semibold text-[1.125rem] text-foreground leading-snug">{selectedData.nama_objek}</h3>
                <p class="mt-1 text-[0.8125rem] text-muted-foreground flex items-center gap-1.5"><MapPin class="w-3.5 h-3.5" /> {selectedData.alamat_dusun || i18n.t('page.risk.popup.subtitle')}</p>
              </div>
            </div>
            <button onclick={() => selected = false} class="grid place-items-center w-8 h-8 rounded-md bg-muted/50 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"><X class="w-4 h-4" /></button>
          </div>

          <div class="mt-5 space-y-4">
            <div>
              <p class="text-[0.75rem] font-semibold text-muted-foreground mb-2">{i18n.t('page.risk.popup.facilities')} & Deskripsi</p>
              <p class="text-[0.8125rem] text-foreground leading-relaxed">
                {selectedData.deskripsi || "Tidak ada deskripsi detail untuk lokasi ini."}
              </p>
            </div>
            
            <div class="grid grid-cols-2 gap-3">
              <div class="p-2.5 rounded-lg bg-muted/40 border border-border/50">
                <p class="text-[0.75rem] text-muted-foreground">{i18n.t('page.risk.popup.capacity')}</p>
                <p class="mt-0.5 font-semibold text-[0.875rem] text-foreground">{selectedData.kapasitas_orang ? selectedData.kapasitas_orang + (i18n.locale === 'id' ? " orang" : " people") : "-"}</p>
              </div>
              <div class="p-2.5 rounded-lg bg-muted/40 border border-border/50">
                <p class="text-[0.75rem] text-muted-foreground">{i18n.locale === 'id' ? 'Kategori' : 'Category'}</p>
                <p class="mt-0.5 font-semibold text-[0.875rem] text-foreground">{selectedData.kategori}</p>
              </div>
            </div>

            {#if rekomendasiFaskes}
              <div class="p-3 rounded-lg bg-sky-500/10 border border-sky-500/20">
                <p class="text-[0.75rem] font-semibold text-sky-600 mb-1">
                   <Sparkles class="w-3 h-3 inline-block mr-0.5 -mt-0.5" /> AI Rekomendasi Faskes
                </p>
                <p class="text-[0.8125rem] text-foreground font-medium">{rekomendasiFaskes.nama_objek}</p>
                <p class="text-[0.75rem] text-muted-foreground mt-0.5"><MapPin class="w-3 h-3 inline-block mr-0.5" /> {rekomendasiFaskes.alamat_dusun}</p>
              </div>
            {/if}
          </div>

          <a href="/safe-routes?destName={encodeURIComponent(selectedData.nama_objek)}&destLat={JSON.parse(selectedData.geometri).coordinates[1]}&destLng={JSON.parse(selectedData.geometri).coordinates[0]}" class="w-full mt-5 h-11 inline-flex items-center justify-center rounded-xl bg-primary text-primary-foreground font-semibold text-[0.875rem] hover:bg-primary-dark transition-colors">
            <Navigation class="w-4 h-4 mr-2" /> {i18n.t('page.risk.popup.route')}
          </a>
        </div>
      {/if}

    </div>
  </div>
</div>
