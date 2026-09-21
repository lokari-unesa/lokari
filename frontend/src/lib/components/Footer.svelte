<script lang="ts">
  import { MapPin } from "lucide-svelte";
  import { i18n } from "$lib/i18n.svelte";
</script>

<!-- The top border uses primary color -->
<footer class="bg-[#121826] text-white border-t-2 border-primary mt-auto">
  <div class="mx-auto max-w-360 px-4 lg:px-8 pt-16 pb-6">
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-12 lg:gap-8 pb-12">
      
      <!-- Column 1: Brand -->
      <div class="flex flex-col">
        <a href="/" class="flex items-center gap-3 mb-4 group" aria-label="Beranda">
          <img src="/logo.webp" alt="LOKARI Logo" class="h-10 w-auto object-contain transition-transform group-hover:scale-105" />
          <h2 class="font-display text-[1.75rem] font-bold text-white tracking-wide uppercase shadow-primary/20 text-shadow-sm">
            LOKARI
          </h2>
        </a>
        <p class="text-[0.9375rem] text-[#94a3b8] leading-relaxed max-w-sm mb-6">
          {i18n.t('footer.brand.desc')}
        </p>
        <div class="flex items-center gap-3">
          <a href="https://instagram.com/lokari.unesa" target="_blank" rel="noopener noreferrer" aria-label="Instagram" class="grid place-items-center w-10 h-10 rounded-full bg-white text-[#121826] hover:bg-primary hover:text-white transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-instagram"><rect width="20" height="20" x="2" y="2" rx="5" ry="5"/><path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"/><line x1="17.5" x2="17.51" y1="6.5" y2="6.5"/></svg>
          </a>
        </div>
      </div>

      <!-- Column 2: Useful Links -->
      <div class="flex flex-col">
        <h3 class="text-[1.125rem] font-bold text-white mb-6">
          {i18n.t('footer.links')}
        </h3>
        <ul class="space-y-4 text-[0.9375rem] text-[#94a3b8] flex flex-col items-start">
          <li><a href="/danger-map" class="hover:text-primary transition-colors flex items-center gap-2">» {i18n.t('footer.links.risk')}</a></li>
          <li><a href="/safe-routes" class="hover:text-primary transition-colors flex items-center gap-2">» {i18n.t('footer.links.evac')}</a></li>
          <li><a href="/news" class="hover:text-primary transition-colors flex items-center gap-2">» {i18n.t('footer.links.info')}</a></li>
          <li><a href="/about" class="hover:text-primary transition-colors flex items-center gap-2">» {i18n.t('footer.company.about')}</a></li>
        </ul>
      </div>

      <!-- Column 3: Kontak Darurat (WhatsApp) -->
      <div class="flex flex-col">
        <h3 class="text-[1.125rem] font-bold text-white mb-6">
          {i18n.t('footer.company.contact')}
        </h3>
        <ul class="space-y-4 text-[0.9375rem] text-[#94a3b8] flex flex-col items-start">
          <li>
            <a href="https://wa.me/6285604122031" target="_blank" rel="noopener noreferrer" class="hover:text-primary transition-colors flex items-center text-left gap-3 group">
              <span class="grid place-items-center w-9 h-9 shrink-0 rounded-full bg-[#1e293b] group-hover:bg-[#25D366] text-white transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51h-.57c-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z"/></svg>
              </span>
              <span class="flex-1">{i18n.t('footer.phone.kades')}</span>
            </a>
          </li>
          <li>
            <a href="https://wa.me/6281216928483" target="_blank" rel="noopener noreferrer" class="hover:text-primary transition-colors flex items-center text-left gap-3 group">
              <span class="grid place-items-center w-9 h-9 shrink-0 rounded-full bg-[#1e293b] group-hover:bg-[#25D366] text-white transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="currentColor"><path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51h-.57c-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z"/></svg>
              </span>
              <span class="flex-1">{i18n.t('footer.phone.siaga')}</span>
            </a>
          </li>
        </ul>
      </div>

      <!-- Column 4: Jam Operasional -->
      <div class="flex flex-col">
        <h3 class="text-[1.125rem] font-bold text-white mb-6">
          {i18n.t('footer.hours.title')}
        </h3>
        <ul class="space-y-4 text-[0.9375rem] w-full max-w-70 md:max-w-none">
          <li class="flex justify-between border-b border-[#334155]/60 pb-2">
            <span class="text-[#94a3b8]">{i18n.t('footer.hours.d1')}</span>
            <span class="text-white font-semibold">{i18n.t('footer.hours.t1')}</span>
          </li>
          <li class="flex justify-between border-b border-[#334155]/60 pb-2">
            <span class="text-[#94a3b8]">{i18n.t('footer.hours.d2')}</span>
            <span class="text-white font-semibold">{i18n.t('footer.hours.t2')}</span>
          </li>
          <li class="flex justify-between pt-1">
            <span class="text-[#94a3b8]">{i18n.t('footer.hours.d3')}</span>
            <span class="text-red-400 font-semibold">{i18n.t('footer.hours.t3')}</span>
          </li>
        </ul>
      </div>

    </div>

    <!-- Middle Section: Address & Maps -->
    <div class="py-10 border-t border-b border-[#334155]/60 grid lg:grid-cols-2 gap-8 items-center text-left">
      <div class="flex flex-col lg:flex-row items-start gap-4">
        <MapPin class="w-6 h-6 text-[#ef4444] shrink-0 mt-0.5" />
        <div>
          <h4 class="text-[1rem] font-bold text-white mb-2">{i18n.t('footer.address.title')}</h4>
          <p class="text-[0.9375rem] text-[#94a3b8] leading-relaxed max-w-md">
            {i18n.t('footer.address.val')}
          </p>
        </div>
      </div>
      
      <!-- Fake Maps Embed as an image or styled div matching the screenshot -->
      <div class="w-full h-40 bg-[#e5e7eb] rounded-xl overflow-hidden relative shadow-inner">
        <!-- Using a placeholder pattern that looks like a map or just an iframe -->
        <iframe 
          title="Peta Balai Desa Jarak"
          src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d15814.773822152862!2d112.2215033!3d-7.9255745!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e78f4b000000001%3A0x6b8f376483fb3cf4!2sKantor%20Desa%20Jarak!5e0!3m2!1sid!2sid!4v1700000000000!5m2!1sid!2sid" 
          width="100%" 
          height="100%" 
          style="border:0;" 
          allowfullscreen={false} 
          loading="lazy" 
          referrerpolicy="no-referrer-when-downgrade">
        </iframe>
        <a href="https://maps.google.com/maps?q=Kantor+Desa+Jarak,+Kecamatan+Plosoklaten,+Kabupaten+Kediri,+Jawa+Timur" target="_blank" rel="noopener noreferrer" class="absolute top-2 left-2 px-3 py-1.5 bg-white text-blue-600 text-xs font-bold rounded shadow flex items-center gap-1 hover:text-blue-700">
          Buka di Maps <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path><polyline points="15 3 21 3 21 9"></polyline><line x1="10" y1="14" x2="21" y2="3"></line></svg>
        </a>
      </div>
    </div>

    <!-- Bottom Copyright Bar -->
    <div class="pt-6 flex flex-col md:flex-row items-center justify-between gap-4 text-[0.875rem] text-[#94a3b8] text-left md:text-center">
      <p>{i18n.t('footer.copy')}</p>
    </div>

  </div>
</footer>
