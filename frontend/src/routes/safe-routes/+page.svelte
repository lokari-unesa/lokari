<script lang="ts">
  import {
    Search,
    Navigation,
    MapPin,
    Clock,
    Route as RouteIcon,
    CheckCircle2,
    AlertTriangle,
    Sparkles,
    ArrowRight,
    Loader2,
    Map as MapIcon,
    Crosshair,
    ShieldCheck,
    X,
    Share2,
  } from "lucide-svelte";
  import KeludMapView from "$lib/components/KeludMapView.svelte";
  import StatusBadge from "$lib/components/StatusBadge.svelte";
  import { i18n } from "$lib/i18n.svelte";
  import { getGeoCoordinates, dedupeCoords } from "$lib/geo";
  import { isQueryAllowed, ANOMALY_MESSAGE } from "$lib/queryGuard";
  import { page } from "$app/stores";
  import { onMount } from "svelte";

  // Search state
  let searchQuery = $state("");
  let isLocating = $state(false);
  let originName = $state("Lokasi Saat Ini (Manual)");

  // Coordinates
  let originLat = $state(-7.8688);
  let originLng = $state(112.1427);
  let destLat = $state(Number($page.url.searchParams.get("destLat")) || 0);
  let destLng = $state(Number($page.url.searchParams.get("destLng")) || 0);
  let destName = $state($page.url.searchParams.get("destName") || "");

  // Derived state
  let searched = $derived(destLat !== 0 && destLng !== 0);

  // Shelters
  let shelters = $state<any[]>([]);
  let selectedShelterId = $state("");
  let destKategori = $state("");

  // Map Popup State
  let selected = $state(false); 
  let selectedData: any = $state(null);
  let rekomendasiFaskes: any = $state(null);

  function onMarkerClick(potensi: any) {
    selectedData = potensi;
    selected = true;

    rekomendasiFaskes = null;
    if (!potensi.kategori.includes("Kesehatan")) {
      const faskesList = shelters.filter(p => p.kategori.includes("Kesehatan"));
      if (faskesList.length > 0) {
        const pCoord = getGeoCoordinates(potensi.geometri);
        let minDist = Infinity;
        let closest = null;
        for (const f of faskesList) {
          const fCoord = getGeoCoordinates(f.geometri);
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

  // Toast State
  let toastMessage = $state("");
  let toastVisible = $state(false);

  function showToast(msg: string) {
    toastMessage = msg;
    toastVisible = true;
    setTimeout(() => {
      toastVisible = false;
    }, 3000);
  }

  onMount(async () => {
    fetchShelters();
    if (searched) {
      // If accessed via search page, try to locate user instantly
      getUserLocation();
    }
  });

  function getUserLocation() {
    isLocating = true;
    originName = "Mendeteksi...";
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          originLat = position.coords.latitude;
          originLng = position.coords.longitude;
          originName = "Lokasi Anda (Akurat)";
          isLocating = false;
          sortShelters();
        },
        (error) => {
          originName = "Desa Jarak (Default)";
          isLocating = false;
        },
        { enableHighAccuracy: true, timeout: 5000 },
      );
    } else {
      originName = "Desa Jarak (Default)";
      isLocating = false;
    }
  }

  async function fetchShelters() {
    try {
      const res = await fetch("/api/potensi");
      const data = await res.json();
      if (data.data) {
        shelters = data.data.map((s: any) => {
          let geom =
            typeof s.geometri === "string"
              ? JSON.parse(s.geometri)
              : s.geometri;
          return { ...s, lat: geom.coordinates[1], lng: geom.coordinates[0] };
        });
        sortShelters();

        // Match selection if passed via URL
        if (destName || (destLat !== 0 && destLng !== 0)) {
          const matched = shelters.find((s) => 
            s.nama_objek === destName || 
            (Math.abs(s.lat - destLat) < 0.0001 && Math.abs(s.lng - destLng) < 0.0001)
          );
          if (matched) {
            selectedShelterId = matched.id_potensi;
            destKategori = matched.kategori || "";
            destLat = matched.lat;
            destLng = matched.lng;
            destName = matched.nama_objek;
          }
        }
      }
    } catch (e) {}
  }

  function sortShelters() {
    shelters = shelters
      .map((s) => {
        return { ...s, dist: calcDistance(originLat, originLng, s.lat, s.lng) };
      })
      .sort((a, b) => a.dist - b.dist);
  }

  function handleShelterChange(e: Event) {
    const id = (e.target as HTMLSelectElement).value;
    const s = shelters.find((x) => x.id_potensi === id);
    if (s) {
      destLat = s.lat;
      destLng = s.lng;
      destName = s.nama_objek;
      destKategori = s.kategori || "";
    }
  }

  let isAiSearching = $state(false);

  async function submitAISearch() {
    if (!searchQuery.trim()) return;

    // Filter Defensif: modul bersama $lib/queryGuard (validasi sebenarnya di backend)
    if (!isQueryAllowed(searchQuery)) {
      alert(ANOMALY_MESSAGE);
      return;
    }

    // Koreksi Typo Singkatan (Pre-processing) agar Vektor Embedding tidak bingung
    let cleanQuery = searchQuery.toLowerCase();
    const typoMap: Record<string, string> = {
      msjd: "masjid",
      mshl: "musholla",
      mshola: "musholla",
      mshlla: "musholla",
      sklh: "sekolah",
      rs: "rumah sakit",
      puskes: "puskesmas",
      bale: "balai",
      pnt: "panti",
      lpngn: "lapangan",
    };

    // Ganti kata-kata singkatan dengan kata aslinya
    for (const [typo, correct] of Object.entries(typoMap)) {
      cleanQuery = cleanQuery.replace(
        new RegExp(`\\b${typo}\\b`, "g"),
        correct,
      );
    }

    isAiSearching = true;
    try {
      const res = await fetch("/api/search", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          query: cleanQuery,
          lat: originLat,
          lng: originLng,
        }),
      });
      const data = await res.json();

      if (data.data && data.data.length > 0) {
        let bestMatch = data.data[0];

        // RERANKING LOGIC: Jika user mencari sesuatu yang "dekat", kita akan mengandalkan
        // kecerdasan AI untuk menebak kategorinya (dari result #1), lalu kita pilih
        // secara murni berdasarkan jarak fisik terdakat dari kategori tersebut!
        if (cleanQuery.includes("dekat") || cleanQuery.includes("paling dekat")) {
          const targetKategori = bestMatch.kategori;
          
          // Cari semua shelter LOKAL yang kategorinya sama dengan tebakan AI
          let sameCategoryShelters = shelters.filter(s => s.kategori === targetKategori);
          
          // Tambahan pengaman: Jika kategori terlalu luas (misal "Tempat Ibadah"), 
          // dan user secara spesifik mengetik "musholla", jangan kasih "masjid".
          if (cleanQuery.includes("musholla")) {
            sameCategoryShelters = sameCategoryShelters.filter(s => s.nama_objek.toLowerCase().includes("musholla") || s.nama_objek.toLowerCase().includes("mushola"));
          } else if (cleanQuery.includes("masjid")) {
            sameCategoryShelters = sameCategoryShelters.filter(s => s.nama_objek.toLowerCase().includes("masjid"));
          }

          // Karena array `shelters` sudah diurutkan berdasarkan jarak fisik terdekat (dist) di `sortShelters()`,
          // kita cukup mengambil elemen PERTAMA dari array yang sudah difilter ini!
          if (sameCategoryShelters.length > 0) {
            bestMatch = sameCategoryShelters[0];
          }
        }

        // Cari objek yang cocok di list dropdown (menggunakan id_potensi karena nama bisa duplikat/mirip)
        const matched = shelters.find(
          (s) => s.id_potensi === bestMatch.id_potensi || s.nama_objek === bestMatch.nama_objek
        );
        if (matched) {
          selectedShelterId = matched.id_potensi;
          // Update destinasi agar peta langsung me-render ulang
          destLat = matched.lat;
          destLng = matched.lng;
          destName = matched.nama_objek;
          destKategori = matched.kategori || "";
        } else {
          alert(i18n.t('page.evac.ai.notFoundLocal'));
        }
      } else {
        alert(data.message || i18n.t('page.evac.ai.noResult'));
      }
    } catch (e) {
      console.error(e);
      alert(i18n.t('page.evac.ai.error'));
    } finally {
      isAiSearching = false;
    }
  }

  function calcDistance(
    lat1: number,
    lon1: number,
    lat2: number,
    lon2: number,
  ) {
    if (!lat2 || !lon2) return 0;
    const R = 6371;
    const dLat = ((lat2 - lat1) * Math.PI) / 180;
    const dLon = ((lon2 - lon1) * Math.PI) / 180;
    const a =
      Math.sin(dLat / 2) * Math.sin(dLat / 2) +
      Math.cos((lat1 * Math.PI) / 180) *
        Math.cos((lat2 * Math.PI) / 180) *
        Math.sin(dLon / 2) *
        Math.sin(dLon / 2);
    return R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  }

  let realDistanceKm = $state<number | null>(null);
  let realEtaMinutes = $state<number | null>(null);
  let routeStatus = $state<"safe" | "fallback" | "danger">("safe");
  let routeCoords = $state<any[]>([]);

  type RouteResult = {
    distanceKm: number | null;
    etaMinutes: number | null;
    status: "safe" | "fallback" | "danger";
    coords: any[];
  };

  // Cache hasil rute per pasangan koordinat (dibulatkan 5 desimal ≈ 1 m)
  // agar $effect tidak me-refetch rute yang sama berulang kali.
  const routeCache = new Map<string, RouteResult>();

  function fallbackRoute(): RouteResult {
    return {
      distanceKm: null,
      etaMinutes: null,
      status: "fallback",
      coords: [[originLat, originLng], [destLat, destLng]],
    };
  }

  function applyRouteResult(result: RouteResult) {
    realDistanceKm = result.distanceKm;
    realEtaMinutes = result.etaMinutes;
    routeStatus = result.status;
    routeCoords = result.coords;
  }

  $effect(() => {
    if (destLat === 0 && destLng === 0) return;

    const key = `${originLat.toFixed(5)},${originLng.toFixed(5)}|${destLat.toFixed(5)},${destLng.toFixed(5)}`;
    const cached = routeCache.get(key);
    if (cached) {
      applyRouteResult(cached);
      return;
    }

    fetch(`/api/route?start=${originLng},${originLat}&end=${destLng},${destLat}`)
      .then(r => r.json())
      .then(data => {
        if (data.routes && data.routes.length > 0) {
          const route = data.routes[0];
          // Server sudah mengembalikan polyline lengkap (termasuk endpoint) —
          // buang titik duplikat, dan jaga-jaga jika endpoint tidak persis cocok.
          const poly = dedupeCoords((route.geometry.coordinates as any[]).map((c: any[]) => [c[1], c[0]]));
          if (poly.length === 0) {
            poly.push([originLat, originLng]);
          }
          const first = poly[0];
          const last = poly[poly.length - 1];
          if (Math.abs(first[0] - originLat) > 1e-5 || Math.abs(first[1] - originLng) > 1e-5) {
            poly.unshift([originLat, originLng]);
          }
          if (Math.abs(last[0] - destLat) > 1e-5 || Math.abs(last[1] - destLng) > 1e-5) {
            poly.push([destLat, destLng]);
          }
          const result: RouteResult = {
            distanceKm: route.distance / 1000,
            etaMinutes: Math.round(route.duration / 60),
            status: route.status || (route.isSafe ? "safe" : "danger"),
            coords: poly,
          };
          routeCache.set(key, result);
          applyRouteResult(result);
        } else {
          const result = fallbackRoute();
          routeCache.set(key, result);
          applyRouteResult(result);
        }
      })
      .catch(e => {
        console.error("Gagal menarik rute:", e);
        const result = fallbackRoute();
        routeCache.set(key, result);
        applyRouteResult(result);
      });
  });

  let distanceKm = $derived(realDistanceKm !== null ? realDistanceKm : calcDistance(originLat, originLng, destLat, destLng));
  // Tanpa ETA dari routing engine, tampilkan "—", bukan angka karangan.
  let etaMinutes = $derived<number | null>(realEtaMinutes);

  function formatDistance(distKm: number) {
    const meters = Math.round(distKm * 1000);
    if (meters < 1000) {
      return `${meters} Meter`;
    }
    return `${distKm.toFixed(1)} KM`;
  }
</script>

{#snippet Metric(Icon: any, label: string, value: string)}
  <div class="p-3 rounded-lg bg-muted/50">
    <p class="text-[0.75rem] text-muted-foreground flex items-center gap-1">
      <Icon class="w-3.5 h-3.5" />
      {label}
    </p>
    <p class="font-display font-bold text-[1.25rem] text-foreground mt-0.5">
      {value}
    </p>
  </div>
{/snippet}

<div class="mx-auto max-w-360 px-4 lg:px-8 py-8">
  <div class="mb-6">
    <h1 class="font-display font-semibold text-[1.5rem] text-foreground">
      {i18n.t("page.evac.title")}
    </h1>
    <p class="mt-1 text-[0.9375rem] text-muted-foreground">
      {i18n.t("page.evac.desc")}
    </p>
  </div>

  <!-- Search panel -->
  <div
    class="p-5 rounded-xl bg-card border border-border {!searched
      ? 'mb-40'
      : ''}"
  >
    <div class="grid grid-cols-1 md:grid-cols-2 gap-3 items-end">
      <!-- Lokasi Anda (Manual & GPS) -->
      <div>
        <label
          for="input-origin"
          class="block text-[0.8125rem] font-semibold text-foreground mb-1.5"
          >{i18n.t("page.evac.from")}</label
        >
        <div class="flex gap-2">
          <div
            class="flex-1 flex items-center gap-2 h-11 px-3.5 rounded-xl border border-input bg-muted/40 focus-within:ring-1 focus-within:ring-primary transition-all overflow-hidden"
          >
            <input
              id="input-origin"
              type="text"
              bind:value={originName}
              oninput={() => {
                // Jika diketik manual, set titik awal ke tengah Desa Jarak (Simulasi) agar tidak 71km
                originLat = -7.8688;
                originLng = 112.1427;
                sortShelters();
              }}
              class="w-full bg-transparent text-[0.9375rem] outline-none placeholder:text-muted-foreground"
              placeholder="Ketik lokasi Anda (Manual)..."
            />
          </div>
          <button
            onclick={getUserLocation}
            class="shrink-0 flex items-center justify-center w-11 h-11 bg-primary/10 text-primary hover:bg-primary hover:text-primary-foreground rounded-xl transition-colors shadow-sm"
            title="Gunakan GPS Saat Ini"
          >
            {#if isLocating}
              <Loader2 class="w-5 h-5 animate-spin" />
            {:else}
              <Crosshair class="w-5 h-5" />
            {/if}
          </button>
        </div>
      </div>

      <!-- Posko atau Titik Kumpul (Dropdown) -->
      <div>
        <label
          for="select-dest"
          class="block text-[0.8125rem] font-semibold text-foreground mb-1.5"
          >{i18n.t("page.evac.to")}</label
        >
        <div
          class="flex items-center gap-2 h-11 px-3.5 rounded-xl border border-input bg-muted/40 focus-within:ring-1 focus-within:ring-primary transition-all w-full overflow-hidden"
        >
          <Navigation class="w-4 h-4 text-muted-foreground shrink-0" />
          <select
            id="select-dest"
            bind:value={selectedShelterId}
            onchange={handleShelterChange}
            class="w-full min-w-0 bg-transparent text-[0.9375rem] outline-none text-foreground cursor-pointer"
          >
            <option value="" disabled selected class="bg-card text-foreground"
              >Pilih Tempat Rekomendasi...</option
            >
            {#if shelters.length === 0}
              <option value="" disabled class="bg-card text-foreground"
                >Sedang memuat data...</option
              >
            {/if}
            {#each shelters as s}
              <option value={s.id_potensi} class="bg-card text-foreground"
                >{s.nama_objek} ({s.dist
                  ? s.dist.toFixed(1) + " km"
                  : ""})</option
              >
            {/each}
          </select>
        </div>
      </div>
    </div>

    <!-- Natural language search -->
    <div class="mt-4 pt-4 border-t border-border">
      <div class="flex flex-col gap-2.5">
        <div class="flex items-center gap-2">
          <span
            class="grid place-items-center w-7 h-7 rounded-md bg-secondary/10 text-secondary shrink-0"
            ><Sparkles class="w-3.5 h-3.5" /></span
          >
          <p class="text-[0.8125rem] font-semibold text-foreground">
            {i18n.t("page.evac.ai.title")}
          </p>
        </div>
        <div
          class="flex items-center gap-2 h-11 px-3.5 rounded-xl border border-input bg-muted/40 w-full overflow-hidden focus-within:ring-1 focus-within:ring-primary transition-all"
        >
          <input
            bind:value={searchQuery}
            onkeydown={(e) => e.key === "Enter" && submitAISearch()}
            class="w-full min-w-0 bg-transparent text-[0.9375rem] outline-none placeholder:text-muted-foreground text-ellipsis"
            placeholder={i18n.t("page.evac.ai.placeholder")}
          />
          <button
            onclick={submitAISearch}
            disabled={isAiSearching}
            class="text-primary shrink-0 hover:text-primary-dark transition-colors disabled:opacity-50"
          >
            {#if isAiSearching}
              <Loader2 class="w-5 h-5 animate-spin" />
            {:else}
              <ArrowRight class="w-5 h-5" />
            {/if}
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Result -->
  {#if searched}
    <div class="mt-6 grid lg:grid-cols-[1fr_360px] gap-4">
      <div
        class="relative h-105 lg:h-130 rounded-xl"
      >
        <KeludMapView
          showHazard={true}
          showRoute={true}
          showShelters={true}
          highlight={true}
          {destLat}
          {destLng}
          {destKategori}
          {originLat}
          {originLng}
          {onMarkerClick}
          potensiData={shelters}
          {routeCoords}
        />

        <!-- Floating info popover -->
        {#if selected && selectedData}
          <!-- Mobile backdrop -->
          <!-- svelte-ignore a11y_click_events_have_key_events -->
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <div class="fixed inset-0 z-40 bg-black/40 backdrop-blur-sm lg:hidden" onclick={() => selected = false}></div>

          <div class="fixed inset-0 m-auto h-fit w-[calc(100vw-2rem)] max-w-85 z-50 lg:absolute lg:inset-auto lg:top-4 lg:left-4 lg:m-0 lg:w-[320px] lg:z-10 p-5 rounded-xl bg-card border border-border shadow-[0_8px_24px_-8px_rgba(15,23,42,0.18)] max-h-[calc(100%-2rem)] overflow-y-auto">
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
                <p class="mt-0.5 font-semibold text-[0.875rem] text-foreground">{selectedData.kapasitas_orang ? selectedData.kapasitas_orang + " orang" : "-"}</p>
              </div>
              <div class="p-2.5 rounded-lg bg-muted/40 border border-border/50">
                <p class="text-[0.75rem] text-muted-foreground">Kategori</p>
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

            <button onclick={() => {
              selectedShelterId = selectedData.id_potensi;
              handleShelterChange({ target: { value: selectedData.id_potensi } } as any);
              selected = false;
            }} class="w-full mt-5 h-11 inline-flex items-center justify-center rounded-xl bg-primary text-primary-foreground font-semibold text-[0.875rem] hover:bg-primary-dark transition-colors">
              <Navigation class="w-4 h-4 mr-2" /> {i18n.t("page.evac.route.setTarget")}
            </button>
          </div>
        {/if}
      </div>

      <div class="flex flex-col gap-4">
        <div class="p-5 rounded-xl bg-card border border-border">
          <div class="flex items-center gap-2 mb-1">
            <RouteIcon class="w-4 h-4 text-primary" />
            <span class="text-[0.8125rem] font-semibold text-primary"
              >{i18n.t("page.evac.route.title")}</span
            >
          </div>
          <h2
            class="font-display font-semibold text-[1.125rem] text-foreground leading-tight"
          >
            {destName}
          </h2>

          <div class="mt-4 flex flex-col gap-3">
            {@render Metric(
              RouteIcon,
              i18n.t("page.evac.route.dist"),
              formatDistance(distanceKm),
            )}
            {@render Metric(
              Clock,
              i18n.t("page.evac.route.est"),
              etaMinutes !== null ? `${etaMinutes} ${i18n.t("page.evac.unit.min")}` : "—",
            )}
          </div>

          {#if routeStatus === "safe"}
            <div class="mt-4 flex items-center gap-2 p-3 rounded-lg bg-safe/10 border border-safe/20">
              <CheckCircle2 class="w-5 h-5 text-safe" />
              <span class="text-[0.875rem] font-medium text-safe">{i18n.t("page.evac.route.safe")}</span>
            </div>
          {:else if routeStatus === "fallback"}
            <div class="mt-4 flex items-start gap-2 p-3 rounded-lg bg-warning/10 border border-warning/20">
              <AlertTriangle class="w-5 h-5 text-warning shrink-0 mt-0.5" />
              <div class="flex flex-col">
                <span class="text-[0.875rem] font-bold text-warning">Mode Rute Standar</span>
                <span class="text-[0.75rem] text-warning/90 mt-0.5">Sistem penghindar zona merah sedang *offline*. Rute ini adalah rute terdekat standar. Harap perhatikan sekeliling Anda.</span>
              </div>
            </div>
          {:else}
            <div class="mt-4 flex items-start gap-2 p-3 rounded-lg bg-destructive/10 border border-destructive/20">
              <AlertTriangle class="w-5 h-5 text-destructive shrink-0 mt-0.5" />
              <div class="flex flex-col">
                <span class="text-[0.875rem] font-bold text-destructive">Peringatan: Rute Darurat!</span>
                <span class="text-[0.75rem] text-destructive/90 mt-0.5">Rute ini terpaksa melintasi Zona Rawan Bencana karena tidak ada jalan memutar yang aman. Harap tingkatkan kewaspadaan.</span>
              </div>
            </div>
          {/if}

          <div class="mt-4 flex flex-col gap-2">
            <a
              href="https://www.google.com/maps/dir/?api=1&origin={originLat},{originLng}&destination={destLat},{destLng}"
              target="_blank"
              class="w-full inline-flex items-center justify-center h-11 rounded-xl bg-primary text-primary-foreground font-semibold text-[0.875rem] hover:bg-primary-dark transition-colors"
            >
              {i18n.t("page.evac.route.openMap")}
            </a>
            <button
              onclick={() => {
                const url = `${window.location.origin}/safe-routes?destLat=${destLat}&destLng=${destLng}&destName=${encodeURIComponent(destName)}`;
                navigator.clipboard.writeText(url);
                showToast(i18n.t("page.evac.route.shareMsg"));
              }}
              class="w-full inline-flex items-center justify-center h-11 rounded-xl bg-muted/40 text-foreground font-semibold text-[0.875rem] hover:bg-muted transition-colors border border-border"
            >
              <Share2 class="w-4 h-4 mr-2 text-muted-foreground" /> {i18n.t("page.evac.route.shareBtn")}
            </button>
          </div>
        </div>

        <!-- Alternative -->
        <div class="p-4 rounded-xl bg-card border border-border">
          <div class="flex items-center gap-2 mb-3">
            <AlertTriangle class="w-4 h-4 text-warning" />
            <span class="text-[0.8125rem] font-semibold text-foreground"
              >{i18n.t("page.evac.alt.title")}</span
            >
          </div>
          <p
            class="font-display font-semibold text-[0.9375rem] text-foreground"
          >
            Posko Cadangan Terdekat
          </p>
          <div
            class="mt-3 flex flex-col gap-2.5 text-[0.8125rem] text-muted-foreground"
          >
            <span class="flex items-center gap-2"
              ><RouteIcon class="w-4 h-4 text-primary" />
              Bukaan rute ke posko lain dapat dilihat di Peta Bahaya.</span
            >
            <div class="mt-0.5">
              <StatusBadge
                level={routeStatus === "safe" ? "aman" : routeStatus === "fallback" ? "waspada" : "awas"}
                label={routeStatus === "safe" ? i18n.t("page.evac.alt.safe") : routeStatus === "fallback" ? "Standar" : "Ekstra Waspada"}
                size="sm"
                pulse={false}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>

<!-- Custom Toast Notification -->
{#if toastVisible}
  <div class="fixed top-8 left-1/2 -translate-x-1/2 z-100 animate-in fade-in slide-in-from-top-5 duration-300 pointer-events-none w-[90vw] max-w-100">
    <div class="flex items-start sm:items-center gap-3 px-5 py-4 bg-primary text-primary-foreground rounded-2xl shadow-xl shadow-primary/20 border border-primary-dark">
      <CheckCircle2 class="w-5 h-5 shrink-0 mt-0.5 sm:mt-0" />
      <span class="font-medium text-[0.9375rem] leading-snug">{toastMessage}</span>
    </div>
  </div>
{/if}
