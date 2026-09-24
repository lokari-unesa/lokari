export const translations = {
  id: {
    // Nav
    'nav.home': 'Beranda',
    'nav.riskMap': 'Peta Bahaya',
    'nav.evacuation': 'Jalur Evakuasi',
    'nav.info': 'Kabar Kelud',
    'nav.about': 'Tentang',
    'nav.emergency': 'Kontak Darurat',
    
    // Hero
    'hero.badge': 'Sistem Pemantauan Cerdas Desa Jarak, Kediri',
    'hero.title': 'Peta Pantau Gunung Kelud',
    'hero.subtitle': 'Aplikasi untuk melihat bahaya letusan Gunung Kelud dan aliran lahar di Kali Ngobo secara langsung.',
    'hero.btn.map': 'Lihat Peta Bahaya',
    'hero.btn.route': 'Cari Jalan Selamat',
    'hero.map.danger': 'Zona Bahaya',
    'hero.map.shelter': 'Tempat Aman',
    
    // Status
    'status.title': 'Kondisi Gunung Kelud Saat Ini',
    'status.alert': 'WASPADA',
    'status.desc': 'Tetap waspada dan ikuti petunjuk dari desa.',
    'status.updated': 'Diperbarui:',
    'status.source': 'Sumber:',
    'status.sourceVal': 'PVMBG MAGMA',
    'status.monitoring': 'MEMANTAU',
    'status.normal': 'Gunung Kelud berstatus Normal (Level I). Tidak ada aktivitas yang membahayakan — warga dapat beraktivitas seperti biasa, namun tetap waspada.',
    'status.waspada': 'Gunung Kelud berstatus Waspada (Level II). Jauhi aliran sungai (Kali Ngobo) dan pantau terus informasi resmi PVMBG serta arahan aparat desa.',
    'status.siaga': 'Gunung Kelud berstatus Siaga (Level III). Siapkan tas siaga, hindari area rawan lahar, dan ikuti arahan petugas.',
    'status.awas': 'Gunung Kelud berstatus Awas (Level IV). Segera menuju tempat aman terdekat dan ikuti instruksi BPBD serta aparat desa.',
    'status.nodata': 'Belum ada data status terbaru dari PVMBG. Sistem tetap memantau aktivitas Gunung Kelud.',
    
    // Quick Actions
    'qa.risk.title': 'Peta Bahaya',
    'qa.risk.desc': 'Lihat daerah mana saja yang tidak aman dari lahar.',
    'qa.route.title': 'Jalan Selamat',
    'qa.route.desc': 'Cari jalan tercepat menuju tempat aman terdekat.',
    'qa.shelter.title': 'Tempat Mengungsi',
    'qa.shelter.desc': 'Cari balai desa dan lokasi aman berkumpul.',
    'qa.info.title': 'Kabar Penting',
    'qa.info.desc': 'Baca berita terbaru yang mudah dipahami.',
    
    // Info Cards
    'info.title': 'Kabar Gunung Kelud & Pesan Penting',
    'info.desc': 'Rangkuman keadaan dari pemerintah yang sudah disederhanakan.',
    'info.btn': 'Semua Berita',
    'info.c1.cat': 'Gunung Api',
    'info.c1.title': 'Kondisi Gunung Kelud',
    'info.c1.desc': 'Penjelasan otomatis tentang keadaan gunung saat ini.',
    'info.c2.cat': 'Lahar',
    'info.c2.title': 'Bahaya Lahar Kali Ngobo',
    'info.c2.desc': 'Peringatan lahar dingin bila hujan deras di puncak kelud.',
    'info.c3.cat': 'Evakuasi',
    'info.c3.title': 'Keadaan Jalan Mengungsi',
    'info.c3.desc': 'Kabar jalan mana yang aman dilewati warga Dusun Simbar Kidul.',
    
    // UI & Pages
    'ui.search': 'Cari',
    'ui.search.info': 'Cari Info',
    'ui.lightMode': 'Terang',
    'ui.darkMode': 'Gelap',
    'ui.language': 'Bahasa:',
    
    'page.evac.title': 'Cari Rute Selamat',
    'page.evac.desc': 'Temukan jalan selamat yang paling dekat dengan lokasi Anda.',
    'page.evac.from': 'Lokasi Anda',
    'page.evac.to': 'Posko atau Titik Kumpul',
    'page.evac.btn': 'Cari Jalan',
    'page.evac.ai.title': 'Pencarian Berbasis AI',
    'page.evac.ai.placeholder': 'Ketik lokasi Anda...',
    'page.evac.ai.noResult': 'Tidak ditemukan rekomendasi untuk lokasi tersebut.',
    'page.evac.ai.notFoundLocal': 'Tempat ditemukan oleh AI, tapi tidak ada di data posko lokal LOKARI.',
    'page.evac.ai.error': 'Terjadi kesalahan saat memproses pencarian cerdas.',
    'page.evac.route.openMap': 'Buka Google Maps Real',
    'page.evac.route.shareBtn': 'Bagikan Rute',
    'page.evac.route.shareMsg': 'Tautan rute evakuasi berhasil disalin! Anda bisa membagikannya sekarang.',
    'page.evac.route.setTarget': 'Jadikan Tujuan Evakuasi',
    'page.evac.unit.min': 'menit',
    'page.evac.route.title': 'Rute Evakuasi Terdekat',
    'page.evac.route.subtitle': 'Titik Kumpul Desa Jarak',
    'page.evac.route.dist': 'Jarak',
    'page.evac.route.est': 'Estimasi',
    'page.evac.route.safe': 'Aman digunakan',
    'page.evac.route.detail': 'Lihat Detail Rute',
    'page.evac.route.map': 'Buka di Peta',
    'page.evac.alt.title': 'Rute Alternatif',
    'page.evac.alt.subtitle': 'Posko Balai Desa Jarak',
    'page.evac.alt.safe': 'Aman',
    'page.evac.map.highlight': 'Rute evakuasi disorot hijau',
    
    'page.risk.title': 'Peta Pantau Bahaya',
    'page.risk.desc': 'Jelajahi batas zona bahaya dan tempat pengungsian.',
    'page.risk.layers': 'Lapisan Peta',
    'page.risk.legend': 'Keterangan Peta',
    'page.risk.hint': 'Ketuk pin pada peta untuk melihat detail lokasi.',
    'page.risk.popup.title': 'Posko Pengungsian Desa Jarak',
    'page.risk.popup.subtitle': 'Desa Jarak',
    'page.risk.popup.facilities': 'Fasilitas',
    'page.risk.popup.capacity': 'Kapasitas',
    'page.risk.popup.status': 'Status',
    'page.risk.popup.active': 'Aktif',
    'page.risk.popup.route': 'Lihat Rute',
    'layer.hazard': 'Zona Risiko',
    'layer.shelter': 'Posko & Titik Kumpul',
    'layer.health': 'Fasilitas Kesehatan',
    'layer.lahar': 'Potensi Lahar',
    'legend.hazard': 'Zona Risiko Tinggi',
    'legend.shelter': 'Posko & Titik Kumpul',
    'legend.health': 'Fasilitas Kesehatan',
    'legend.lahar': 'Potensi Lahar',
    
    'page.info.title': 'Kabar Kelud & Desa',
    'page.info.desc': 'Kumpulan berita resmi tentang keadaan Gunung Kelud.',
    'page.info.filter': 'Pilih Kategori',
    'page.info.empty': 'Belum ada informasi untuk kategori ini.',
    'info.filter.all': 'Semua',
    'info.filter.volcano': 'Gunung Api',
    'info.filter.lahar': 'Lahar',
    'info.filter.evac': 'Evakuasi',
    'info.filter.weather': 'Cuaca',
    'info.filter.warning': 'Peringatan',
    'info.card.readmore': 'Baca selengkapnya',
    'info.date.now': 'Baru Saja',
    
    'page.about.title': 'Tentang LOKARI',
    'page.about.desc': 'Platform WebGIS mitigasi bencana erupsi Gunung Kelud berbasis Zero-Admin dan Spatial AI, dirancang dengan antarmuka Active and Simplified khusus untuk masyarakat Desa Jarak.',
    'page.about.innovations': 'Inovasi Teknologi Utama',
    'page.about.inv.zero': 'Zero-Admin Otomatis',
    'page.about.inv.zero.desc': 'Pembaruan data aktivitas vulkanik secara real-time langsung dari server tanpa memerlukan pengelolaan manual oleh perangkat desa.',
    'page.about.inv.nlp': 'Spatial AI (NLP Summarizer)',
    'page.about.inv.nlp.desc': 'Mengonversi data teknis gunung api yang rumit menjadi kalimat peringatan dini dan instruksi keselamatan yang mudah dipahami warga.',
    'page.about.inv.semantic': 'Pencarian Rute Pintar (Semantic Search)',
    'page.about.inv.semantic.desc': 'Memungkinkan warga mencari rute evakuasi teraman dan posko terdekat menggunakan pertanyaan dalam bahasa sehari-hari.',
    'page.about.team': 'Tim Pengembang',
    'page.about.collab': 'Kolaborasi Multidisiplin Tim UNESA',
    'page.about.collab.desc': 'Platform ini dibangun oleh 15 mahasiswa Universitas Negeri Surabaya (UNESA) lintas disiplin ilmu:',
    'page.about.collab.ti': 'Teknik Informatika',
    'page.about.collab.ti.desc': 'Pengembangan arsitektur backend, basis data spasial (PostgreSQL/PostGIS), dan integrasi Kecerdasan Buatan (AI).',
    'page.about.collab.si': 'Sistem Informasi',
    'page.about.collab.si.desc': 'Perancangan antarmuka pengguna (UI/UX) interaktif dan pengujian kegunaan antarmuka peta.',
    'page.about.collab.an': 'Administrasi Negara',
    'page.about.collab.an.desc': 'Penyelarasan kebijakan publik, komunikasi birokrasi, dan evaluasi dampak bagi masyarakat desa.',
    'page.about.partners': 'Kemitraan Strategis',
    'page.about.partners.desc': 'LOKARI merupakan luaran dari program Studi Independen Mobilitas Akademik yang bermitra dengan:',
    'page.about.sources': 'Sumber Data Resmi',
    'page.about.sources.desc': 'Data peringatan dini dan pemantauan aktivitas pada platform ini bersumber dari lembaga resmi pemerintah dan internasional:',
    
    // Search Page (Simple Language)
    'page.search.title': 'Pencarian Informasi',
    'page.search.desc': 'Ketik apa saja yang ingin Anda ketahui tentang evakuasi, jalan aman, atau posko terdekat.',
    'page.search.placeholder': 'Cth. Posko Terdekat',
    'page.search.btn': 'Cari',
    'page.search.ans.title': 'Jawaban',
    'page.search.ans.desc': 'Saat terjadi erupsi, segera menuju balai desa. Jangan mendekati sungai (Kali Ngobo) karena bahaya lahar dingin. Jangan lupa bawa surat penting dan obat-obatan.',
    'page.search.route.title': 'Perkiraan Jalan',
    'page.search.route.dist': 'Jarak',
    'page.search.route.time': 'Waktu',
    'page.search.source.title': 'Sumber Asli',
    'page.search.source.desc': 'Sesuai dengan arahan resmi dari BMKG dan Pemerintah Desa Jarak.',
    'page.search.source.badge': 'Paling Baru',
    'page.search.dest.title': 'Tempat Aman Terdekat',
    'page.search.dest.name': 'Titik Kumpul Lapangan Desa',
    'page.search.dest.desc': 'Lapangan Utama Desa Jarak, Area Bebas Lahar.',
    'page.search.dest.btn1': 'Tunjukkan Jalan',
    'page.search.dest.btn2': 'Buka Peta',
    
    // Search Results
    'search.status.busy': 'Layanan AI Sedang Sibuk',
    'search.status.empty': 'Tidak menemukan lokasi yang cocok dengan query semantik Anda. Coba bahasa atau frasa lain.',
    'search.status.loading': 'Sedang Mencari...',
    'search.result.category': 'Kategori',
    'search.result.noDesc': 'Tidak ada deskripsi detail untuk lokasi ini.',
    'search.result.capacity': 'Kapasitas Penampungan',
    'search.result.people': 'Orang',

    // New Footer
    'footer.brand.desc': 'LOKARI Platform mitigasi yang menyediakan Layanan Pemantauan Bencana Erupsi dan Peta Evakuasi Desa Jarak.',
    'footer.links': 'Tautan Penting',
    'footer.links.risk': 'Peta Bahaya',
    'footer.links.evac': 'Cari Rute Selamat',
    'footer.links.info': 'Berita Terkini',
    'footer.company': 'Informasi',
    'footer.company.about': 'Tentang Kami',
    'footer.company.contact': 'Kontak Desa',
    'footer.hours.title': 'Jam Operasional',
    'footer.hours.d1': 'Senin - Kamis',
    'footer.hours.t1': '08.00 - 13.00 WIB',
    'footer.hours.d2': 'Jumat',
    'footer.hours.t2': '08.00 - 11.00 WIB',
    'footer.hours.d3': 'Sabtu & Minggu',
    'footer.hours.t3': 'Libur',
    'footer.phone.kades': 'Bapak Moh. Toha (Kepala Desa)',
    'footer.phone.siaga': 'Bapak Bayan Heri (Mobil Siaga)',
    'footer.address.title': 'Alamat Balai Desa',
    'footer.address.val': 'Desa Jarak, Kecamatan Plosoklaten, Kabupaten Kediri, Jawa Timur 64175, Indonesia',
    'footer.copy': 'Copyrights © 2026. All rights reserved by LOKARI Desa Jarak',
    
    // Errors
    'error.quota.alert.src': 'Sistem LOKARI (Peringatan Standar)',
    'error.quota.alert.msg': 'Maaf, limit harian sistem AI cerdas kami saat ini telah habis. Silakan hubungi aparat desa. Sebagai peringatan standar, pantau terus arahan resmi dari NASA dan BPBD setempat.',
    'error.quota.search.msg': 'Maaf, jatah limit pencarian AI untuk bulan ini telah habis. Silakan gunakan fitur Peta Bahaya sementara waktu atau hubungi Bapak Kades.'
  },
  en: {
    // Nav
    'nav.home': 'Home',
    'nav.riskMap': 'Danger Map',
    'nav.evacuation': 'Evacuation Route',
    'nav.info': 'Kelud News',
    'nav.about': 'About',
    'nav.emergency': 'Emergency',
    
    // Hero
    'hero.badge': 'Smart Monitoring System of Jarak Village, Kediri',
    'hero.title': 'Mount Kelud Monitoring Map',
    'hero.subtitle': 'Application to view Mount Kelud eruption hazards and Kali Ngobo cold lava flows in real-time.',
    'hero.btn.map': 'View Danger Map',
    'hero.btn.route': 'Find Safe Route',
    'hero.map.danger': 'Danger Zone',
    'hero.map.shelter': 'Safe Place',
    
    // Status
    'status.title': 'Current Mount Kelud Condition',
    'status.alert': 'WARNING',
    'status.desc': 'Stay alert and follow village instructions.',
    'status.updated': 'Updated:',
    'status.source': 'Source:',
    'status.sourceVal': 'PVMBG MAGMA',
    'status.monitoring': 'MONITORING',
    'status.normal': 'Mount Kelud is at NORMAL level (Level I). No activity endangers residents — normal activities may resume, but stay alert.',
    'status.waspada': 'Mount Kelud is at WATCH level (Level II). Stay away from rivers (Kali Ngobo) and keep monitoring official information from PVMBG and village officials.',
    'status.siaga': 'Mount Kelud is at ALERT level (Level III). Prepare an emergency bag, avoid lahar-prone areas, and follow instructions from officials.',
    'status.awas': 'Mount Kelud is at DANGER level (Level IV). Immediately head to the nearest safe place and follow BPBD and village officials\' instructions.',
    'status.nodata': 'No latest status data from PVMBG yet. The system continues monitoring Mount Kelud.',
    
    // Quick Actions
    'qa.risk.title': 'Danger Map',
    'qa.risk.desc': 'See which areas are unsafe from cold lava.',
    'qa.route.title': 'Safe Route',
    'qa.route.desc': 'Find the fastest path to the nearest safe place.',
    'qa.shelter.title': 'Evacuation Place',
    'qa.shelter.desc': 'Find village halls and safe gathering locations.',
    'qa.info.title': 'Important News',
    'qa.info.desc': 'Read the latest easy-to-understand news.',
    
    // Info Cards
    'info.title': 'Kelud News & Important Messages',
    'info.desc': 'Simplified situation summary from the government.',
    'info.btn': 'All News',
    'info.c1.cat': 'Volcano',
    'info.c1.title': 'Mount Kelud Condition',
    'info.c1.desc': 'Automated explanation of the current volcano condition.',
    'info.c2.cat': 'Lahar',
    'info.c2.title': 'Kali Ngobo Lahar Danger',
    'info.c2.desc': 'Cold lava warning when heavy rain occurs at the peak.',
    'info.c3.cat': 'Evacuation',
    'info.c3.title': 'Evacuation Route Condition',
    'info.c3.desc': 'Updates on safe roads for Simbar Kidul residents.',
    
    // UI & Pages
    'ui.search': 'Search',
    'ui.search.info': 'Find Info',
    'ui.lightMode': 'Light',
    'ui.darkMode': 'Dark',
    'ui.language': 'Lang:',
    
    'page.evac.title': 'Find Safe Route',
    'page.evac.desc': 'Find the closest safe path from your location.',
    'page.evac.from': 'Your Location',
    'page.evac.to': 'Shelter or Meeting Point',
    'page.evac.btn': 'Find Path',
    'page.evac.ai.title': 'AI-Based Search',
    'page.evac.ai.placeholder': 'Type your location...',
    'page.evac.ai.noResult': 'No recommendation found for that location.',
    'page.evac.ai.notFoundLocal': 'Location found by AI, but not present in LOKARI local shelter data.',
    'page.evac.ai.error': 'An error occurred while processing smart search.',
    'page.evac.route.openMap': 'Open Google Maps',
    'page.evac.route.shareBtn': 'Share Route',
    'page.evac.route.shareMsg': 'Evacuation route link copied! You can share it now.',
    'page.evac.route.setTarget': 'Set as Evacuation Destination',
    'page.evac.unit.min': 'mins',
    'page.evac.route.title': 'Nearest Evacuation Route',
    'page.evac.route.subtitle': 'Jarak Village Meeting Point',
    'page.evac.route.dist': 'Distance',
    'page.evac.route.est': 'ETA',
    'page.evac.route.safe': 'Safe to use',
    'page.evac.route.detail': 'View Route Details',
    'page.evac.route.map': 'Open in Map',
    'page.evac.alt.title': 'Alternative Route',
    'page.evac.alt.subtitle': 'Jarak Village Hall',
    'page.evac.alt.safe': 'Safe',
    'page.evac.map.highlight': 'Evacuation route highlighted in green',
    
    'page.risk.title': 'Danger Monitor Map',
    'page.risk.desc': 'Explore hazard zones and evacuation shelters.',
    'page.risk.layers': 'Map Layers',
    'page.risk.legend': 'Legend',
    'page.risk.hint': 'Tap a pin on the map to see location details.',
    'page.risk.popup.title': 'Jarak Village Evacuation Post',
    'page.risk.popup.subtitle': 'Jarak Village',
    'page.risk.popup.facilities': 'Facilities',
    'page.risk.popup.capacity': 'Capacity',
    'page.risk.popup.status': 'Status',
    'page.risk.popup.active': 'Active',
    'page.risk.popup.route': 'View Route',
    'layer.hazard': 'Risk Zone',
    'layer.shelter': 'Shelters & Assembly Points',
    'layer.health': 'Health Facilities',
    'layer.lahar': 'Lahar Potential',
    'legend.hazard': 'High Risk Zone',
    'legend.shelter': 'Shelters & Assembly Points',
    'legend.health': 'Health Facilities',
    'legend.lahar': 'Lahar Potential',
    
    'page.info.title': 'Kelud & Village News',
    'page.info.desc': 'Collection of official news about Mount Kelud.',
    'page.info.filter': 'Select Category',
    'page.info.empty': 'No information available for this category yet.',
    'info.filter.all': 'All',
    'info.filter.volcano': 'Volcano',
    'info.filter.lahar': 'Lahar',
    'info.filter.evac': 'Evacuation',
    'info.filter.weather': 'Weather',
    'info.filter.warning': 'Warning',
    'info.card.readmore': 'Read more',
    'info.date.now': 'Just Now',
    
    'page.about.title': 'About LOKARI',
    'page.about.desc': 'A WebGIS platform for Mount Kelud eruption mitigation based on Zero-Admin and Spatial AI, designed with an Active and Simplified interface specifically for the Jarak Village community.',
    'page.about.innovations': 'Key Technological Innovations',
    'page.about.inv.zero': 'Zero-Admin Automation',
    'page.about.inv.zero.desc': 'Real-time updates of volcanic activity data directly from the server without requiring manual management by village officials.',
    'page.about.inv.nlp': 'Spatial AI (NLP Summarizer)',
    'page.about.inv.nlp.desc': 'Converts complex technical volcano data into easy-to-understand early warning sentences and safety instructions for residents.',
    'page.about.inv.semantic': 'Smart Route Search (Semantic Search)',
    'page.about.inv.semantic.desc': 'Allows residents to find the safest evacuation routes and nearest shelters using everyday language.',
    'page.about.team': 'Development Team',
    'page.about.collab': 'UNESA Multidisciplinary Team',
    'page.about.collab.desc': 'This platform was built by 15 students from Surabaya State University (UNESA) across 3 disciplines:',
    'page.about.collab.ti': 'Informatics Engineering',
    'page.about.collab.ti.desc': 'Development of backend architecture, spatial database (PostgreSQL/PostGIS), and Artificial Intelligence (AI) integration.',
    'page.about.collab.si': 'Information Systems',
    'page.about.collab.si.desc': 'Interactive user interface (UI/UX) design and usability testing of the map interface.',
    'page.about.collab.an': 'Public Administration',
    'page.about.collab.an.desc': 'Alignment of public policy, bureaucratic communication, and community impact evaluation.',
    'page.about.partners': 'Strategic Partnerships',
    'page.about.partners.desc': 'LOKARI is the output of an Independent Academic Mobility Study program in partnership with:',
    'page.about.sources': 'Official Data Sources',
    'page.about.sources.desc': 'Early warning and activity monitoring data on this platform are sourced from official government and international agencies:',
    
    // Search Page (Simple Language)
    'page.search.title': 'Information Search',
    'page.search.desc': 'Type anything you want to know about evacuation, safe routes, or nearest shelters.',
    'page.search.placeholder': 'e.g. Nearest shelter',
    'page.search.btn': 'Search',
    'page.search.ans.title': 'Answer',
    'page.search.ans.desc': 'During an eruption, immediately head to the village hall. Do not approach the river (Kali Ngobo) due to cold lava danger. Don\'t forget to bring important documents and medicines.',
    'page.search.route.title': 'Route Estimate',
    'page.search.route.dist': 'Distance',
    'page.search.route.time': 'Time',
    'page.search.source.title': 'Official Source',
    'page.search.source.desc': 'In accordance with official directives from BMKG and Jarak Village Government.',
    'page.search.source.badge': 'Up to Date',
    'page.search.dest.title': 'Nearest Safe Place',
    'page.search.dest.name': 'Village Field Gathering Point',
    'page.search.dest.desc': 'Main Field of Jarak Village, Lahar-Free Area.',
    'page.search.dest.btn1': 'Show the Way',
    'page.search.dest.btn2': 'Open Map',
    
    // Search Results
    'search.status.busy': 'AI Service is Busy',
    'search.status.empty': 'Could not find a location that matches your semantic query. Try different words.',
    'search.status.loading': 'Searching...',
    'search.result.category': 'Category',
    'search.result.noDesc': 'No detailed description available for this location.',
    'search.result.capacity': 'Shelter Capacity',
    'search.result.people': 'People',

    // New Footer
    'footer.brand.desc': 'LOKARI Platform providing Disaster Monitoring and Evacuation Maps for Jarak Village.',
    'footer.links': 'Useful Links',
    'footer.links.risk': 'Danger Map',
    'footer.links.evac': 'Find Safe Route',
    'footer.links.info': 'Latest News',
    'footer.company': 'Company',
    'footer.company.about': 'About Us',
    'footer.company.contact': 'Village Contact',
    'footer.hours.title': 'Operational Hours',
    'footer.hours.d1': 'Monday - Thursday',
    'footer.hours.t1': '08.00 - 13.00 WIB',
    'footer.hours.d2': 'Friday',
    'footer.hours.t2': '08.00 - 11.00 WIB',
    'footer.hours.d3': 'Saturday & Sunday',
    'footer.hours.t3': 'Closed',
    'footer.phone.kades': 'Mr. Moh. Toha (Village Head)',
    'footer.phone.siaga': 'Mr. Bayan Heri (Emergency Car)',
    'footer.address.title': 'Village Hall Address',
    'footer.address.val': 'Jl. Raya Jarak, Jarak Village, Plosoklaten Dist., Kediri Regency, East Java, Indonesia',
    'footer.copy': 'Copyrights © 2026. All rights reserved by LOKARI Desa Jarak',
    
    // Errors
    'error.quota.alert.src': 'LOKARI System (Standard Warning)',
    'error.quota.alert.msg': 'Sorry, the daily limit for our smart AI system has been exhausted. Please contact village officials. As a standard warning, keep monitoring official directives from NASA and the local BPBD.',
    'error.quota.search.msg': 'Sorry, the AI search limit quota for this month has been exhausted. Please use the Danger Map feature temporarily or contact the Village Head.'
  }
};

type Locale = keyof typeof translations;
type TranslationKey = keyof typeof translations.id;

function createI18n() {
  let locale = $state<Locale>('id');

  // Load from localStorage if available (client side only)
  if (typeof window !== 'undefined') {
    const saved = localStorage.getItem('lokari-locale') as Locale;
    if (saved && translations[saved]) {
      locale = saved;
    }
  }

  return {
    get locale() {
      return locale;
    },
    set locale(value: Locale) {
      locale = value;
      if (typeof window !== 'undefined') {
        localStorage.setItem('lokari-locale', value);
        document.documentElement.lang = value;
      }
    },
    t: (key: TranslationKey) => {
      return translations[locale][key] || key;
    },
    toggle: () => {
      locale = locale === 'id' ? 'en' : 'id';
      if (typeof window !== 'undefined') {
        localStorage.setItem('lokari-locale', locale);
        document.documentElement.lang = locale;
      }
    }
  };
}

export const i18n = createI18n();
