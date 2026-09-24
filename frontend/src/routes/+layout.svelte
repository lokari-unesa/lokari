<script lang="ts">
	import './layout.css';
	import Navbar from '$lib/components/Navbar.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import { onMount } from 'svelte';
	import { env } from '$env/dynamic/public';

	let { children } = $props();

	// Cegah penumpukan listener fallback bila subscribeToPush dipanggil berulang
	// saat masih gagal (mis. hot-reload saat halaman terbuka).
	let pushRetryArmed = false;

	onMount(async () => {
		if ('serviceWorker' in navigator && 'PushManager' in window) {
			try {
				await navigator.serviceWorker.register('/sw.js');
				// Tunggu SW aktif dulu — hindari AbortError saat kunjungan pertama
				const registration = await navigator.serviceWorker.ready;

				if (Notification.permission === 'default') {
					// Memicu native pop-up seperti kompas.com setelah 2 detik
					setTimeout(async () => {
						const permission = await Notification.requestPermission();
						if (permission === 'granted') {
							await subscribeToPush(registration);
						}
					}, 2000);
				} else if (Notification.permission === 'granted') {
					await subscribeToPush(registration);
				}
			} catch(e) {
				console.error("SW / Push Error:", e);
			}
		}
	});

	// Dipanggil ketika permission sudah "granted". Percobaan pertama `subscribe`
	// bisa gagal sementara tepat setelah user menekan Izinkan (AbortError "push
	// service error") — makanya di-retry dengan backoff, plus fallback otomatis
	// saat tab difokus agar tidak perlu reload manual.
	async function subscribeToPush(registration: ServiceWorkerRegistration) {
		try {
			const vapidKey = env.PUBLIC_VAPID_KEY;
			if (!vapidKey) {
				console.warn("VAPID Key belum dikonfigurasi, notifikasi tidak aktif.");
				return;
			}

			const urlB64ToUint8Array = (base64String: string) => {
				const padding = '='.repeat((4 - base64String.length % 4) % 4);
				const base64 = (base64String + padding).replace(/\-/g, '+').replace(/_/g, '/');
				const rawData = window.atob(base64);
				const outputArray = new Uint8Array(rawData.length);
				for (let i = 0; i < rawData.length; ++i) {
					outputArray[i] = rawData.charCodeAt(i);
				}
				return outputArray;
			};

			// Pakai subscription browser yang sudah ada; kalau belum, buat baru
			// dengan retry untuk error transien.
			let subscription: PushSubscription | null = await registration.pushManager.getSubscription();
			if (!subscription) {
				const MAX_ATTEMPTS = 3;
				for (let attempt = 1; attempt <= MAX_ATTEMPTS && !subscription; attempt++) {
					try {
						subscription = await registration.pushManager.subscribe({
							userVisibleOnly: true,
							applicationServerKey: urlB64ToUint8Array(vapidKey)
						});
					} catch (e) {
						console.warn(`[Push] subscribe percobaan ${attempt}/${MAX_ATTEMPTS} gagal:`, e);
						if (attempt < MAX_ATTEMPTS) {
							// Backoff 1.2s → 2.4s — beri waktu push service browser selesai.
							await new Promise((r) => setTimeout(r, 1200 * attempt));
						}
					}
				}
			}

			if (!subscription) {
				console.error("[Push] Gagal subscribe setelah beberapa percobaan — mencoba lagi saat tab difokus.");
				if (!pushRetryArmed) {
					pushRetryArmed = true;
					const retry = async () => {
						pushRetryArmed = false;
						window.removeEventListener('focus', retry);
						await subscribeToPush(registration);
					};
					window.addEventListener('focus', retry);
				}
				return;
			}

			const subJSON = subscription.toJSON();
			// Relatif /api/* agar konsisten dengan halaman lain:
			//   - dev manual   → vite proxy (vite.config.ts)
			//   - build/prod   → hooks.server.ts (BACKEND_URL)
			await fetch('/api/subscribe', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					endpoint: subscription.endpoint,
					keys: {
						p256dh: subJSON.keys?.p256dh,
						auth: subJSON.keys?.auth
					}
				})
			});
			console.log("Berhasil mendaftar Notifikasi LOKARI!");
		} catch(e) {
			console.error("Gagal subscribe Push:", e);
		}
	}
</script>

<div class="min-h-screen flex flex-col bg-background">
  <Navbar />
  <main class="flex-1">
    {@render children()}
  </main>
  <Footer />
</div>
