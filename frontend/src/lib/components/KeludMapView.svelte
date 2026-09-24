<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { cn } from "$lib/utils";
  import { getGeoCoordinates } from "$lib/geo";
  import desaJarakGeoJSON from "$lib/data/desa_jarak.json";

  let {
    showHazard = true,
    showRoute = false,
    showShelters = true,
    showHealth = false,
    showLahar = false,
    highlight = false,
    class: className = "",
    compact = false,
    destLat = -7.971,
    destLng = 112.213,
    destKategori = "",
    originLat = -7.952,
    originLng = 112.235,
    potensiData = [],
    onMarkerClick,
    routeCoords = []
  } = $props<{
    showHazard?: boolean;
    showRoute?: boolean;
    showShelters?: boolean;
    showHealth?: boolean;
    showLahar?: boolean;
    highlight?: boolean;
    class?: string;
    compact?: boolean;
    destLat?: number;
    destLng?: number;
    destKategori?: string;
    originLat?: number;
    originLng?: number;
    potensiData?: any[];
    onMarkerClick?: (potensi: any) => void;
    routeCoords?: any[];
  }>();

  let mapElement: HTMLElement;
  let map: any;
  let L: any;
  
  let routeLayer: any;
  let originMarker: any;
  let destMarker: any;
  let dynamicMarkers: any[] = []; // Untuk penanda dinamis dari potensiData
  let hazardGroup: any;
  let laharLayer: any;
  let desaJarakLayer: any;

  // Koordinat Gunung Kelud
  const keludLat = -7.9333;
  const keludLng = 112.3083;

  let mapReady = $state(false);

  onMount(async () => {
    if (typeof window === 'undefined') return;

    // Dynamic import to prevent SSR issues in SvelteKit
    L = await import('leaflet');
    await import('leaflet/dist/leaflet.css');

    // Initialize Map
    const centerLat = showRoute ? (originLat + destLat) / 2 : originLat;
    const centerLng = showRoute ? (originLng + destLng) / 2 : originLng;
    
    map = L.map(mapElement, {
      zoomControl: false,
      attributionControl: false
    }).setView([centerLat, centerLng], 13);

    // Google Satellite Hybrid (Satelit + Jalan + Label)
    L.tileLayer('https://mt1.google.com/vt/lyrs=y&x={x}&y={y}&z={z}', {
      maxZoom: 20,
      attribution: '&copy; Google Maps'
    }).addTo(map);

    // Tandai peta siap — ini akan memicu $effect pertama kali
    mapReady = true;
  });

  // Custom DivIcon factories (dipanggil setiap kali karena Leaflet butuh instance baru)
  function createOriginIcon() {
    return L.divIcon({
      className: 'bg-transparent border-0',
      html: `<div class="relative w-6 h-6 rounded-full bg-slate-900 border-2 border-white shadow-md flex items-center justify-center">
               <div class="absolute w-full h-full rounded-full border-2 border-slate-900 animate-ping opacity-75"></div>
               <div class="w-2 h-2 rounded-full bg-white"></div>
             </div>`,
      iconSize: [24, 24],
      iconAnchor: [12, 12]
    });
  }

  function createDestIcon(kategori?: string) {
    let color = 'bg-green-700'; // Default Posko / Titik Kumpul
    let iconHTML = 'P';

    if (kategori) {
      if (kategori.includes('Kesehatan')) {
        color = 'bg-sky-600';
        iconHTML = '🏥';
      } else if (kategori.includes('Ibadah')) {
        color = 'bg-teal-600'; 
        iconHTML = '🕌';
      } else if (kategori.includes('Pendidikan')) {
        color = 'bg-amber-600';
        iconHTML = '🎓';
      } else if (kategori.includes('Balai') || kategori.includes('Gedung')) {
        color = 'bg-indigo-600';
        iconHTML = '🏛️';
      } else {
        iconHTML = '📍';
      }
    }

    return L.divIcon({
      className: 'bg-transparent border-0',
      html: `<div class="w-9 h-9 rounded-full ${color} border-2 border-white shadow-lg text-white font-bold text-sm flex items-center justify-center" style="border-bottom-right-radius: 0; transform: rotate(-45deg);">
               <div style="transform: rotate(45deg);">${iconHTML}</div>
             </div>`,
      iconSize: [36, 36],
      iconAnchor: [18, 36]
    });
  }

  function createDynamicIcon(kategori: string) {
    let color = 'bg-emerald-600'; // Default Posko / Titik Kumpul
    let iconHTML = '📍';

    if (kategori.includes('Kesehatan')) {
      color = 'bg-sky-600';
      iconHTML = '🏥';
    } else if (kategori.includes('Ibadah')) {
      color = 'bg-teal-600'; 
      iconHTML = '🕌';
    } else if (kategori.includes('Pendidikan')) {
      color = 'bg-amber-600';
      iconHTML = '🎓';
    } else if (kategori.includes('Balai') || kategori.includes('Gedung')) {
      color = 'bg-indigo-600';
      iconHTML = '🏛️';
    }

    return L.divIcon({
      className: 'bg-transparent border-0',
      html: `<div class="w-8 h-8 rounded-full ${color} border-2 border-white shadow-md text-white text-sm flex items-center justify-center">
               ${iconHTML}
             </div>`,
      iconSize: [32, 32],
      iconAnchor: [16, 16]
    });
  }

  // REAKTIVITAS UTAMA: $effect ini melacak perubahan originLat, originLng, destLat, destLng
  // dan otomatis menggambar ulang marker + rute ketika koordinat berubah.
  $effect(() => {
    if (!mapReady || !map || !L) return;

    // Baca semua koordinat (ini membuat Svelte melacak dependensinya)
    const oLat = originLat;
    const oLng = originLng;
    const dLat = destLat;
    const dLng = destLng;
    const dKat = destKategori;
    const wantRoute = showRoute;
    const wantShelters = showShelters;
    const wantHealth = showHealth;
    const wantHazard = showHazard;
    const wantLahar = showLahar;
    const pd = potensiData;
    const rCoords = routeCoords;

    // --- Hazard Zone (KRB I, II, III - Lebih Presisi) ---
    if (wantHazard) {
      if (!hazardGroup) {
        // KRB III (Zona Merah - Sangat Berbahaya) - Radius 5km
        const krb3 = L.circle([keludLat, keludLng], {
          color: '#B91C1C', // Red-700
          fillColor: '#EF4444',
          fillOpacity: 0.35,
          weight: 2,
          dashArray: '4, 6',
          radius: 5000
        });

        // KRB II (Zona Oranye - Berbahaya) - Radius 10km
        const krb2 = L.circle([keludLat, keludLng], {
          color: '#C2410C', // Orange-700
          fillColor: '#F97316',
          fillOpacity: 0.15,
          weight: 2,
          dashArray: '4, 6',
          radius: 10000
        });

        // KRB I (Zona Kuning - Waspada) - Radius 15km
        const krb1 = L.circle([keludLat, keludLng], {
          color: '#A16207', // Yellow-700
          fillColor: '#EAB308',
          fillOpacity: 0.08,
          weight: 2,
          dashArray: '5, 10',
          radius: 15000
        });

        const marker = L.marker([keludLat, keludLng], { icon: L.divIcon({
          className: 'bg-transparent border-0',
          html: `<div class="flex flex-col items-center">
                   <div class="w-4 h-4 rounded-full bg-slate-800 border-2 border-red-500 shadow-[0_0_15px_rgba(239,68,68,0.8)] animate-pulse"></div>
                   <div class="text-[10px] font-bold text-white px-2 py-0.5 mt-1 rounded bg-black/60 backdrop-blur border border-white/20 shadow-md">G. KELUD</div>
                 </div>`,
          iconSize: [60, 40],
          iconAnchor: [30, 8]
        }) });
        hazardGroup = L.layerGroup([krb1, krb2, krb3, marker]).addTo(map);
      } else if (!map.hasLayer(hazardGroup)) {
        hazardGroup.addTo(map);
      }
    } else {
      if (hazardGroup && map.hasLayer(hazardGroup)) {
        map.removeLayer(hazardGroup);
      }
    }

    // --- Lahar Path (Presisi 3 Sungai Utama) ---
    if (wantLahar) {
      if (!laharLayer) {
        // Multi-branch lahar simulation (Menyebar radial ke segala arah, Worst-Case Scenario)
        // 1. Kali Ngobo (Barat menuju Plosoklaten)
        const branchNgobo = [[keludLat, keludLng], [-7.932, 112.295], [-7.933, 112.280], [-7.928, 112.260], [-7.915, 112.245], [-7.905, 112.220]];
        // 2. Kali Konto (Utara menuju Kasembon)
        const branchKonto = [[keludLat, keludLng], [-7.915, 112.312], [-7.890, 112.315], [-7.870, 112.310], [-7.850, 112.305], [-7.830, 112.308]];
        // 3. Kali Bladak / Putih (Selatan menuju Blitar)
        const branchBladak = [[keludLat, keludLng], [-7.950, 112.290], [-7.970, 112.270], [-7.990, 112.255], [-8.020, 112.240], [-8.050, 112.225]];
        // 4. Kali Petung (Barat Laut menuju Puncu / Kepung)
        const branchPetung = [[keludLat, keludLng], [-7.920, 112.290], [-7.900, 112.270], [-7.880, 112.250], [-7.860, 112.235]];
        // 5. Kali Gedog (Tenggara menuju Wlingi)
        const branchGedog = [[keludLat, keludLng], [-7.940, 112.320], [-7.955, 112.335], [-7.970, 112.345], [-7.990, 112.360]];
        // 6. Kali Jari / Semen (Selatan Barat Daya menuju Garum)
        const branchJari = [[keludLat, keludLng], [-7.945, 112.285], [-7.965, 112.255], [-7.985, 112.230], [-8.005, 112.205]];
        // 7. Kali Kuncir (Timur Laut menuju Ngantang)
        const branchKuncir = [[keludLat, keludLng], [-7.920, 112.325], [-7.905, 112.340], [-7.890, 112.360], [-7.875, 112.380]];
        // 8. Celah Timur (Arah Pujon)
        const branchTimur = [[keludLat, keludLng], [-7.933, 112.330], [-7.930, 112.350], [-7.925, 112.370], [-7.920, 112.390]];
        
        laharLayer = L.layerGroup([
          L.polyline(branchNgobo, { color: '#9a3412', weight: 6, opacity: 0.6, lineCap: 'round', lineJoin: 'round' }),
          L.polyline(branchKonto, { color: '#9a3412', weight: 5, opacity: 0.6, lineCap: 'round', lineJoin: 'round' }),
          L.polyline(branchBladak, { color: '#9a3412', weight: 7, opacity: 0.6, lineCap: 'round', lineJoin: 'round' }),
          L.polyline(branchPetung, { color: '#9a3412', weight: 5, opacity: 0.5, lineCap: 'round', lineJoin: 'round' }),
          L.polyline(branchGedog, { color: '#9a3412', weight: 5, opacity: 0.5, lineCap: 'round', lineJoin: 'round' }),
          L.polyline(branchJari, { color: '#9a3412', weight: 6, opacity: 0.6, lineCap: 'round', lineJoin: 'round' }),
          L.polyline(branchKuncir, { color: '#9a3412', weight: 4, opacity: 0.4, lineCap: 'round', lineJoin: 'round' }),
          L.polyline(branchTimur, { color: '#9a3412', weight: 4, opacity: 0.4, lineCap: 'round', lineJoin: 'round' })
        ]).addTo(map);
      } else if (!map.hasLayer(laharLayer)) {
        laharLayer.addTo(map);
      }
    } else {
      if (laharLayer && map.hasLayer(laharLayer)) {
        map.removeLayer(laharLayer);
      }
    }

    // --- Batas Desa Jarak ---
    if (!desaJarakLayer) {
      desaJarakLayer = L.geoJSON(desaJarakGeoJSON, {
        style: {
          color: '#E11D48', // Rose-600 (mirip merah di gambar)
          weight: 2,
          fillColor: '#E11D48',
          fillOpacity: 0.15, // Cukup transparan agar satelit terlihat
          dashArray: '4, 6'
        }
      }).bindTooltip("Wilayah Desa Jarak", { permanent: false, sticky: true }).addTo(map);
    } else if (!map.hasLayer(desaJarakLayer)) {
      desaJarakLayer.addTo(map);
    }

    // Bersihkan marker dinamis lama
    dynamicMarkers.forEach(m => map.removeLayer(m));
    dynamicMarkers = [];

    // --- Dynamic Markers (Dari potensiData) ---
    if (pd && pd.length > 0) {
      pd.forEach((p: any) => {
        const isFaskes = p.kategori && p.kategori.includes('Kesehatan');
        
        // Cek toggle yang aktif
        if ((wantHealth && isFaskes) || (wantShelters && !isFaskes)) {
          let lat, lng;
          if (typeof p.lat === "number" && typeof p.lng === "number") {
            lat = p.lat;
            lng = p.lng;
          } else {
            const [gLng, gLat] = getGeoCoordinates(p.geometri);
            lat = gLat;
            lng = gLng;
          }

          if (lat && lng) {
            // Jangan gambar marker dinamis jika ini adalah titik tujuan (destMarker akan digambar di atasnya)
            const isDestination = Math.abs(lat - dLat) < 0.0001 && Math.abs(lng - dLng) < 0.0001;
            
            if (!isDestination) {
              const m = L.marker([lat, lng], { icon: createDynamicIcon(p.kategori) }).addTo(map);
              if (onMarkerClick) {
                m.on('click', () => onMarkerClick(p));
              }
              dynamicMarkers.push(m);
            }
          }
        }
      });
    }

    // --- Origin Marker ---
    if (highlight) {
      if (originMarker) {
        originMarker.setLatLng([oLat, oLng]);
      } else {
        originMarker = L.marker([oLat, oLng], { icon: createOriginIcon() }).addTo(map);
      }
    } else if (originMarker) {
      map.removeLayer(originMarker);
      originMarker = null;
    }

    // --- Destination Marker ---
    if (showShelters && dLat !== 0 && dLng !== 0) {
      const handleDestClick = () => {
        if (onMarkerClick && pd) {
          const matched = pd.find((p: any) => {
            const [pLng, pLat] = getGeoCoordinates(p.geometri);
            return Math.abs(pLat - dLat) < 0.0001 && Math.abs(pLng - dLng) < 0.0001;
          });
          if (matched) onMarkerClick(matched);
        }
      };

      if (destMarker) {
        destMarker.setLatLng([dLat, dLng]);
        destMarker.setIcon(createDestIcon(dKat));
        destMarker.off('click');
        destMarker.on('click', handleDestClick);
      } else {
        destMarker = L.marker([dLat, dLng], { icon: createDestIcon(dKat) }).addTo(map);
        destMarker.on('click', handleDestClick);
      }
    } else if (destMarker) {
      map.removeLayer(destMarker);
      destMarker = null;
    }

    // --- Route Line ---
    if (wantRoute && dLat !== 0 && dLng !== 0) {
      if (routeLayer) {
        map.removeLayer(routeLayer);
        routeLayer = null;
      }
      if (rCoords && rCoords.length > 0) {
        routeLayer = L.polyline(rCoords, {
          color: '#16A34A',
          weight: 6,
          dashArray: '1, 15',
          lineCap: 'round',
          lineJoin: 'round'
        }).addTo(map);
        
        map.fitBounds(routeLayer.getBounds(), { padding: [50, 50] });
      }
    } else if (routeLayer) {
      map.removeLayer(routeLayer);
      routeLayer = null;
    }
  });

  onDestroy(() => {
    if (map) {
      map.remove();
    }
  });
</script>

<div class={cn("relative w-full h-full overflow-hidden rounded-xl border border-border z-0", className)}>
  <!-- Map Container -->
  <div bind:this={mapElement} class="w-full h-full bg-[#111827]"></div>
  
  {#if !compact}
    <div class="absolute bottom-2.5 left-2.5 flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-card/90 border border-border text-[0.6875rem] font-medium text-muted-foreground backdrop-blur z-1000 shadow-sm pointer-events-none">
      <span class="w-2 h-2 rounded-full bg-primary animate-pulse shadow-[0_0_8px_rgba(22,163,74,0.8)]"></span> Peta Satelit Nyata
    </div>
  {/if}
</div>

<style>
  /* Timpa Z-index bawaan Leaflet agar tidak menutupi Navbar/Dropdown kita */
  :global(.leaflet-container) {
    z-index: 1 !important;
    font-family: inherit;
    background: #111827 !important; /* Warna dasar gelap untuk lautan/ruang kosong */
  }
  
  /* Trik Menggelapkan Peta Satelit agar sesuai dengan Tema Gelap (Tanpa merusak warna asli) */
  :global(.leaflet-tile-pane) {
    filter: brightness(0.65) contrast(1.1) saturate(1.1);
  }

  :global(.leaflet-pane) {
    z-index: 1 !important;
  }
  :global(.leaflet-top),
  :global(.leaflet-bottom) {
    z-index: 2 !important;
  }
</style>
