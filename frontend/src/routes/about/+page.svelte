<script lang="ts">
  import { Mountain, Cpu, Sparkles, Navigation, Users, Code, PenTool, Landmark, Handshake, Database, MapPin, Building, ShieldCheck } from "lucide-svelte";
  import { i18n } from "$lib/i18n.svelte";

  const members = [
    { name: "Aida Rotila", ig: "https://www.instagram.com/ar_aidaa_/", pos: "object-[center_30%]" },
    { name: "Aisyah Nur Fadila", ig: "https://www.instagram.com/aisyahfadilaa_/", pos: "object-[center_30%]" },
    { name: "Asman", ig: "https://www.instagram.com/asman.id/", pos: "object-[center_30%]" },
    { name: "Dammar Sanggalie", ig: "https://www.instagram.com/dammarsanggalie/", pos: "object-center" },
    { name: "Ivan Yudhistira", ig: "#", pos: "object-[center_75%]" }, // Wajah di bawah, langit luas di atas
    { name: "Izza Mutia Aqila", ig: "https://www.instagram.com/mutiaaq17/", pos: "object-center" },
    { name: "Muhammad Anwar Ismail", ig: "https://www.instagram.com/jinwoostream/", pos: "object-[center_35%]" },
    { name: "MUHAMMAD HIKMAL KURNIAH", ig: "https://www.instagram.com/hikmal_krn/", pos: "object-[center_30%]" },
    { name: "Musa Chang Li Panjaitan", ig: "https://www.instagram.com/musachangli_/", pos: "object-[center_30%]" },
    { name: "Orlinezar Raveda Ar-Rizky", ig: "https://www.instagram.com/ravarrizky_/", pos: "object-top" },
    { name: "Ramzi Pandu Ma'arif ", ig: "#", pos: "object-[center_30%]" },
    { name: "Revalia Cantika Sari Setiawan", ig: "https://www.instagram.com/revaliacntikaa/", pos: "object-[center_20%]" },
    { name: "Sabrina Isna Riyani", ig: "https://www.instagram.com/sbrn_nanaa/", pos: "object-center" },
    { name: "Stefi Febianova", ig: "https://www.instagram.com/whtvr.fn/", pos: "object-[center_30%]" },
    { name: "Yulian Devy Ryfka Prastiwi ", ig: "https://www.instagram.com/yulianprst/", pos: "object-[center_30%]" }
  ];
</script>

<div class="mx-auto max-w-360 px-4 lg:px-8 py-10 overflow-x-hidden">
  <!-- Intro -->
  <section class="text-center mb-16">
    <img src="/logo.webp" alt="LOKARI Logo" class="w-20 h-20 mx-auto mb-6 object-contain" />
    <h1 class="font-display font-extrabold text-[2.25rem] md:text-[2.75rem] text-foreground tracking-tight">{i18n.t('page.about.title')}</h1>
    <p class="mt-4 text-[1.0625rem] md:text-[1.125rem] text-muted-foreground max-w-3xl mx-auto leading-relaxed">
      {i18n.t('page.about.desc')}
    </p>
  </section>

  <!-- Inovasi Teknologi -->
  <section class="mt-16">
    <div class="flex items-center gap-3 mb-6">
      <span class="grid place-items-center w-10 h-10 rounded-xl bg-secondary/10 text-secondary"><Cpu class="w-5 h-5" /></span>
      <h2 class="font-display font-bold text-[1.5rem] md:text-[1.75rem] text-foreground">{i18n.t('page.about.innovations')}</h2>
    </div>
    <div class="grid md:grid-cols-3 gap-4">
      <div class="flex flex-col p-6 rounded-2xl bg-card border border-border shadow-[0_8px_24px_-8px_rgba(15,23,42,0.06)] hover:border-primary/30 transition-colors">
        <span class="grid place-items-center w-12 h-12 rounded-xl bg-primary/10 text-primary mb-4"><Database class="w-6 h-6" /></span>
        <h3 class="font-display font-bold text-[1.125rem] text-foreground mb-2">{i18n.t('page.about.inv.zero')}</h3>
        <p class="text-[0.9375rem] text-muted-foreground leading-relaxed">{i18n.t('page.about.inv.zero.desc')}</p>
      </div>
      <div class="flex flex-col p-6 rounded-2xl bg-card border border-border shadow-[0_8px_24px_-8px_rgba(15,23,42,0.06)] hover:border-warning/30 transition-colors">
        <span class="grid place-items-center w-12 h-12 rounded-xl bg-warning/10 text-warning mb-4"><Sparkles class="w-6 h-6" /></span>
        <h3 class="font-display font-bold text-[1.125rem] text-foreground mb-2">{i18n.t('page.about.inv.nlp')}</h3>
        <p class="text-[0.9375rem] text-muted-foreground leading-relaxed">{i18n.t('page.about.inv.nlp.desc')}</p>
      </div>
      <div class="flex flex-col p-6 rounded-2xl bg-card border border-border shadow-[0_8px_24px_-8px_rgba(15,23,42,0.06)] hover:border-destructive/30 transition-colors">
        <span class="grid place-items-center w-12 h-12 rounded-xl bg-destructive/10 text-destructive mb-4"><Navigation class="w-6 h-6" /></span>
        <h3 class="font-display font-bold text-[1.125rem] text-foreground mb-2">{i18n.t('page.about.inv.semantic')}</h3>
        <p class="text-[0.9375rem] text-muted-foreground leading-relaxed">{i18n.t('page.about.inv.semantic.desc')}</p>
      </div>
    </div>
  </section>

  <!-- Kolaborasi Tim -->
  <section class="mt-16">
    <div class="flex flex-col mb-6">
      <span class="grid place-items-center w-12 h-12 rounded-xl bg-primary/10 text-primary mb-3"><Users class="w-6 h-6" /></span>
      <h2 class="font-display font-bold text-[1.5rem] md:text-[1.75rem] text-foreground">{i18n.t('page.about.collab')}</h2>
    </div>
    <p class="text-[1rem] text-muted-foreground mb-6">{i18n.t('page.about.collab.desc')}</p>

    <!-- Foto Bersama -->
    <div class="mb-8 rounded-2xl overflow-hidden border border-border bg-card shadow-sm">
      <img src="/assets/foto-bersama.webp" alt="Tim LOKARI Universitas Negeri Surabaya" class="w-full h-auto max-h-125 object-cover" />
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-10">
      <a href="https://ti.ft.unesa.ac.id/" target="_blank" rel="noopener noreferrer" class="flex flex-col p-6 rounded-2xl bg-muted/30 border border-border hover:border-primary/40 transition-colors">
        <div class="flex items-center gap-3 mb-3">
          <Code class="w-5 h-5 text-primary" />
          <h3 class="font-semibold text-[1.0625rem] text-foreground">{i18n.t('page.about.collab.ti')}</h3>
        </div>
        <p class="text-[0.9375rem] text-muted-foreground leading-relaxed">{i18n.t('page.about.collab.ti.desc')}</p>
      </a>
      <a href="https://si.ft.unesa.ac.id/" target="_blank" rel="noopener noreferrer" class="flex flex-col p-6 rounded-2xl bg-muted/30 border border-border hover:border-primary/40 transition-colors">
        <div class="flex items-center gap-3 mb-3">
          <PenTool class="w-5 h-5 text-secondary" />
          <h3 class="font-semibold text-[1.0625rem] text-foreground">{i18n.t('page.about.collab.si')}</h3>
        </div>
        <p class="text-[0.9375rem] text-muted-foreground leading-relaxed">{i18n.t('page.about.collab.si.desc')}</p>
      </a>
      <a href="https://s1an.fisipol.unesa.ac.id/" target="_blank" rel="noopener noreferrer" class="flex flex-col p-6 rounded-2xl bg-muted/30 border border-border md:col-span-2 lg:col-span-1 hover:border-primary/40 transition-colors">
        <div class="flex items-center gap-3 mb-3">
          <Landmark class="w-5 h-5 text-warning" />
          <h3 class="font-semibold text-[1.0625rem] text-foreground">{i18n.t('page.about.collab.an')}</h3>
        </div>
        <p class="text-[0.9375rem] text-muted-foreground leading-relaxed">{i18n.t('page.about.collab.an.desc')}</p>
      </a>
    </div>

    <!-- Foto Anggota Grid -->
    <div class="pt-8 border-t border-border">
      <h3 class="font-display font-semibold text-[1.125rem] text-foreground mb-6">{i18n.t('page.about.team')}</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each members as m}
          <svelte:element this={m.ig !== '#' ? 'a' : 'div'} href={m.ig !== '#' ? m.ig : undefined} target={m.ig !== '#' ? '_blank' : undefined} rel={m.ig !== '#' ? 'noopener noreferrer' : undefined} class="flex flex-col items-center gap-4 group p-6 rounded-2xl bg-card border border-border shadow-sm hover:border-primary/40 hover:shadow-md transition-all">
            <div class="w-40 h-40 sm:w-48 sm:h-48 rounded-full overflow-hidden border-4 border-muted group-hover:border-primary transition-colors">
              <img 
                src={`/assets/anggota/${m.name}.webp`} 
                alt={m.name} 
                class={`w-full h-full object-cover ${m.pos || 'object-center'}`} 
                loading="lazy" 
              />
            </div>
            <span class="text-[1.0625rem] md:text-[1.125rem] text-center font-display font-semibold text-foreground leading-tight px-2">{m.name}</span>
          </svelte:element>
        {/each}
      </div>
    </div>
  </section>

  <!-- Kemitraan & Sumber Data -->
  <div class="mt-16 grid lg:grid-cols-2 gap-8 mb-10 items-stretch">
    <section class="flex flex-col">
      <div class="flex items-center gap-3 mb-4">
        <span class="grid place-items-center w-10 h-10 rounded-xl bg-accent text-primary"><Handshake class="w-5 h-5" /></span>
        <h2 class="font-display font-bold text-[1.25rem] text-foreground">{i18n.t('page.about.partners')}</h2>
      </div>
      <p class="text-[0.9375rem] text-muted-foreground mb-4">{i18n.t('page.about.partners.desc')}</p>
      <div class="flex flex-col gap-3">
        <a href="https://sagasitas.org/" target="_blank" rel="noopener noreferrer" class="flex items-center gap-4 p-4 rounded-xl border border-border bg-card hover:border-primary/40 transition-colors">
          <img src="/assets/sagasitas.webp" alt="Sagasitas Indonesia" class="w-12 h-12 object-contain rounded-md" />
          <span class="font-medium text-[0.9375rem] text-foreground">Yayasan Sagasitas Indonesia</span>
        </a>
        <a href="https://desajarak.com/" target="_blank" rel="noopener noreferrer" class="flex items-center gap-4 p-4 rounded-xl border border-border bg-card hover:border-primary/40 transition-colors">
          <img src="/assets/kediri.svg" alt="Pemerintah Kabupaten Kediri" class="w-12 h-12 object-contain" />
          <span class="font-medium text-[0.9375rem] text-foreground">Pemerintah Desa Jarak, Kediri</span>
        </a>
      </div>
    </section>

    <section class="flex flex-col">
      <div class="flex items-center gap-3 mb-4">
        <span class="grid place-items-center w-10 h-10 rounded-xl bg-destructive/10 text-destructive"><Database class="w-5 h-5" /></span>
        <h2 class="font-display font-bold text-[1.25rem] text-foreground">{i18n.t('page.about.sources')}</h2>
      </div>
      <p class="text-[0.9375rem] text-muted-foreground mb-4">{i18n.t('page.about.sources.desc')}</p>
      <div class="flex flex-col gap-3">
        <a href="https://magma.esdm.go.id/" target="_blank" rel="noopener noreferrer" class="flex items-center gap-3 p-4 rounded-xl border border-border bg-card hover:border-destructive/40 transition-colors">
          <ShieldCheck class="w-5 h-5 text-destructive shrink-0" />
          <div class="flex flex-col">
            <span class="font-medium text-[0.9375rem] text-foreground">PVMBG MAGMA</span>
            <span class="text-[0.8125rem] text-muted-foreground">Pusat Vulkanologi dan Mitigasi Bencana Geologi (ESDM)</span>
          </div>
        </a>
        <a href="https://www.bmkg.go.id/" target="_blank" rel="noopener noreferrer" class="flex items-center gap-3 p-4 rounded-xl border border-border bg-card hover:border-destructive/40 transition-colors">
          <ShieldCheck class="w-5 h-5 text-destructive shrink-0" />
          <div class="flex flex-col">
            <span class="font-medium text-[0.9375rem] text-foreground">BMKG</span>
            <span class="text-[0.8125rem] text-muted-foreground">Badan Meteorologi, Klimatologi, dan Geofisika</span>
          </div>
        </a>
        <a href="https://eonet.gsfc.nasa.gov/" target="_blank" rel="noopener noreferrer" class="flex items-center gap-3 p-4 rounded-xl border border-border bg-card hover:border-destructive/40 transition-colors">
          <ShieldCheck class="w-5 h-5 text-destructive shrink-0" />
          <div class="flex flex-col">
            <span class="font-medium text-[0.9375rem] text-foreground">NASA EONET</span>
            <span class="text-[0.8125rem] text-muted-foreground">Earth Observatory Natural Event Tracker (Satelit Global)</span>
          </div>
        </a>
      </div>
    </section>
  </div>
</div>

<style>
  @keyframes marquee {
    0% { transform: translateX(0); }
    100% { transform: translateX(-50%); }
  }
  .animate-marquee {
    animation: marquee 30s linear infinite;
  }
  .animate-marquee:hover {
    animation-play-state: paused;
  }
</style>
