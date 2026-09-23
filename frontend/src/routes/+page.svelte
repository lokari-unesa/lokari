<script lang="ts">
  import { onMount } from "svelte";
  import {
    Map,
    Route,
    Tent,
    BellRing,
    Mountain,
    Waves,
    DoorOpen,
    ChevronRight,
    ShieldCheck,
    Activity,
    CloudRain,
    X,
  } from "lucide-svelte";
  import KeludMapView from "$lib/components/KeludMapView.svelte";
  import StatusBadge from "$lib/components/StatusBadge.svelte";
  import QuickActionCard from "$lib/components/QuickActionCard.svelte";
  import InfoCard from "$lib/components/InfoCard.svelte";
  import { i18n } from "$lib/i18n.svelte";
  import { cn } from "$lib/utils";

  let alertData = $state({
    sumber: "Sistem Lokal",
    pesan: i18n.t('status.desc') as string,
    waktu: i18n.t('info.date.now') as string,
    kategori: "warning",
    judul: i18n.t('status.title') as string
  });
  let isLoadingAlert = $state(true);

  let potensiData = $state<any[]>([]);
  let newsItems = $state<any[]>([]);
  let selectedNews = $state<any>(null);

  // Mapping kategori ke Ikon dan Warna untuk InfoCard
  const categoryMeta: Record<string, any> = {
    volcano: { icon: Mountain, accent: "bg-warning/10 text-warning" },
    lahar: { icon: Waves, accent: "bg-[#92400E]/10 text-[#92400E]" },
    evac: { icon: DoorOpen, accent: "bg-primary/10 text-primary" },
    weather: { icon: CloudRain, accent: "bg-secondary text-white" },
    warning: { icon: Activity, accent: "bg-destructive/10 text-destructive" },
  };

  const categoryToSlug: Record<string, string> = {
    "Gunung Api": "volcano",
    "volcano": "volcano",
    "Lahar": "lahar",
    "lahar": "lahar",
    "Evakuasi": "evac",
    "evac": "evac",
    "Cuaca": "weather",
    "weather": "weather",
    "Peringatan": "warning",
    "warning": "warning"
  };

  onMount(async () => {
    // 1. Fetch Potensi Data untuk Map
    try {
      const res = await fetch("/api/potensi");
      const data = await res.json();
      if (data.status === "success" && data.data) {
        potensiData = data.data;
      }
    } catch (e) {
      console.error("Gagal menarik data potensi:", e);
    }

    // 2. Fetch News Data untuk Card Status & List Berita
    try {
      const res = await fetch("/api/news");
      const data = await res.json();
      if (data.status === "success" && data.data && data.data.length > 0) {
        // Ambil berita pertama untuk Status Card utama
        const topNews = data.data[0];
        alertData = {
          sumber: topNews.source,
          pesan: topNews.summary,
          waktu: topNews.created_at || i18n.t('info.date.now'),
          kategori: categoryToSlug[topNews.category] || topNews.category,
          judul: topNews.title
        };

        // Simpan 3 berita terbaru untuk bottom section
        newsItems = data.data.slice(0, 3).map((n: any) => {
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
      console.error("Gagal menarik data berita AI:", e);
    } finally {
      isLoadingAlert = false;
    }
  });
</script>

<div>
  <!-- Hero -->
  <section class="bg-card border-b border-border">
    <div class="mx-auto max-w-360 px-4 lg:px-8">
      <div
        class="grid lg:grid-cols-[55%_45%] gap-8 items-center pt-6 pb-12 lg:py-16"
      >
        <div>
          <div
            class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-accent text-primary text-[0.8125rem] font-semibold"
          >
            <ShieldCheck class="w-4 h-4" /> {i18n.t('hero.badge')}
          </div>
          <h1
            class="mt-5 font-display font-extrabold text-[2rem] lg:text-[2.75rem] leading-[1.1] text-foreground text-balance"
          >
            {i18n.t('hero.title')}
          </h1>
          <p
            class="mt-5 text-[1.0625rem] text-muted-foreground leading-relaxed max-w-xl"
          >
            {i18n.t('hero.subtitle')}
          </p>
          <div class="mt-7 flex flex-col sm:flex-row gap-3">
            <a
              href="/danger-map"
              class="inline-flex items-center justify-center h-12 px-6 rounded-xl bg-primary text-primary-foreground font-semibold text-[0.9375rem] hover:bg-primary-dark transition-colors"
            >
              <Map class="w-5 h-5 mr-2" /> {i18n.t('hero.btn.map')}
            </a>
            <a
              href="/safe-routes"
              class="inline-flex items-center justify-center h-12 px-6 rounded-xl bg-card border border-border text-foreground font-semibold text-[0.9375rem] hover:border-primary hover:text-primary transition-colors"
            >
              <Route class="w-5 h-5 mr-2" /> {i18n.t('hero.btn.route')}
            </a>
          </div>
        </div>

        <div class="relative w-full aspect-9/16 lg:aspect-auto lg:h-105 mt-6 lg:mt-0 rounded-2xl overflow-hidden shadow-lg">
          <KeludMapView
            showHazard={true}
            showShelters={true}
            showRoute={false}
            showHealth={true}
            destLat={0}
            destLng={0}
            {potensiData}
          />
          <div
            class="absolute top-3 right-3 flex items-center gap-2 px-3 py-1.5 rounded-lg bg-card/90 border border-border backdrop-blur text-[0.75rem] font-medium text-muted-foreground"
          >
            <span class="w-2 h-2 rounded-full bg-warning"></span> {i18n.t('hero.map.danger')}
            &nbsp;·&nbsp;
            <span class="w-2 h-2 rounded-full bg-primary"></span> {i18n.t('hero.map.shelter')}
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Status card -->
  <section class="mx-auto max-w-360 px-4 lg:px-8 mt-10">
    <div class="flex flex-col p-6 lg:p-7 rounded-2xl bg-card border border-warning/30 shadow-md">
      <!-- Baris 1: Icon -->
      <div class="mb-4">
        <span class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-warning/20">
          <Activity class="w-6 h-6 text-warning" />
        </span>
      </div>
      
      <!-- Baris 2: Judul Section -->
      <p class="text-[0.875rem] font-bold uppercase tracking-wider text-warning mb-3">
        {i18n.t('status.title')}
      </p>

      <!-- Baris 3: Label Waspada -->
      <div class="mb-4">
        <StatusBadge
          level="waspada"
          label={alertData.kategori === 'warning' || alertData.kategori === 'volcano' ? "WARNING" : "INFO"}
          size="lg"
          pulse={true}
        />
      </div>

      <!-- Baris 4: Isi Berita AI -->
      <p class="text-[1rem] leading-relaxed text-foreground font-medium mb-6 whitespace-pre-wrap">
        {#if isLoadingAlert}
          <span class="inline-block w-4 h-4 border-2 border-primary border-t-transparent rounded-full animate-spin mr-2 align-middle"></span> Memuat analisis AI...
        {:else}
          {alertData.pesan}
        {/if}
      </p>

      <!-- Baris 5 & 6: Footer Info (Waktu & Sumber) -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 pt-5 border-t border-border">
        <p class="text-[0.8125rem] text-muted-foreground">
          <span class="font-semibold text-foreground">{i18n.t('status.updated')}</span> {alertData.waktu}
        </p>
        <p class="text-[0.8125rem] text-muted-foreground">
          <span class="font-semibold text-foreground">{i18n.t('status.source')}</span> {alertData.sumber}
        </p>
      </div>
    </div>
  </section>

  <!-- Quick actions -->
  <section class="mx-auto max-w-360 px-4 lg:px-8 mt-12">
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <QuickActionCard
        to="/danger-map"
        icon={Map}
        title={i18n.t('qa.risk.title')}
        description={i18n.t('qa.risk.desc')}
        accent="primary"
      />
      <QuickActionCard
        to="/safe-routes"
        icon={Route}
        title={i18n.t('qa.route.title')}
        description={i18n.t('qa.route.desc')}
        accent="secondary"
      />
      <QuickActionCard
        to="/danger-map"
        icon={Tent}
        title={i18n.t('qa.shelter.title')}
        description={i18n.t('qa.shelter.desc')}
        accent="safe"
      />
      <QuickActionCard
        to="/news"
        icon={BellRing}
        title={i18n.t('qa.info.title')}
        description={i18n.t('qa.info.desc')}
        accent="warning"
      />
    </div>
  </section>

  <!-- Emergency info -->
  <section class="mx-auto max-w-360 px-4 lg:px-8 mt-14 mb-16">
    <div class="flex items-end justify-between mb-6">
      <div>
        <h2 class="font-display font-semibold text-[1.5rem] text-foreground">
          {i18n.t('info.title')}
        </h2>
        <p class="mt-1 text-[0.9375rem] text-muted-foreground">
          {i18n.t('info.desc')}
        </p>
      </div>
      <a
        href="/news"
        class="hidden sm:inline-flex items-center gap-1 text-[0.875rem] font-semibold text-primary hover:gap-1.5 transition-all"
      >
        {i18n.t('info.btn')} <ChevronRight class="w-4 h-4" />
      </a>
    </div>
    <div class="grid md:grid-cols-3 gap-4">
      {#if isLoadingAlert}
        <div class="col-span-full py-12 text-center text-muted-foreground text-[0.9375rem]">Memuat berita cerdas AI...</div>
      {:else if newsItems.length > 0}
        {#each newsItems as item}
          <InfoCard {...item} onclick={() => selectedNews = item} />
        {/each}
      {:else}
        <!-- Fallback jika belum ada data AI -->
        <InfoCard
          category="volcano"
          title="Sistem Siaga"
          summary="Belum ada peringatan darurat saat ini."
          date="Baru Saja"
          source="Sistem LOKARI"
          icon={Mountain}
          accent="bg-warning/10 text-warning"
          href="/news"
        />
      {/if}
    </div>
  </section>

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
             <span class={cn("inline-flex items-center justify-center w-10 h-10 rounded-xl", selectedNews.accent || "bg-primary/10 text-primary")}>
                <ModalIcon class="w-5 h-5" strokeWidth={2} />
             </span>
             <span class="text-[0.8125rem] font-semibold text-muted-foreground">{i18n.t(`info.filter.${selectedNews.category}` as any)}</span>
          </div>
          <button onclick={() => selectedNews = null} class="inline-flex items-center justify-center w-8 h-8 rounded-md bg-muted/50 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
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
