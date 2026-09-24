<script lang="ts">
	import './layout.css';
	import Navbar from '$lib/components/Navbar.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import { onMount } from 'svelte';
	import { PUBLIC_VAPID_KEY } from '$env/static/public';

	let { children } = $props();

	onMount(async () => {
		if ('serviceWorker' in navigator && 'PushManager' in window) {
			try {
				const registration = await navigator.serviceWorker.register('/sw.js');
				
				if (Notification.permission === 'default') {
					// Memicu native pop-up seperti kompas.com setelah 2 detik
					setTimeout(async () => {
						const permission = await Notification.requestPermission();
						if (permission === 'granted') {
							subscribeToPush(registration);
						}
					}, 2000); 
				} else if (Notification.permission === 'granted') {
					subscribeToPush(registration);
				}
			} catch(e) {
				console.error("SW / Push Error:", e);
			}
		}
	});

	async function subscribeToPush(registration: ServiceWorkerRegistration) {
		try {
			const existing = await registration.pushManager.getSubscription();
			if (existing) return;

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

			const applicationServerKey = urlB64ToUint8Array(PUBLIC_VAPID_KEY);
			const subscription = await registration.pushManager.subscribe({
				userVisibleOnly: true,
				applicationServerKey
			});

			const subJSON = subscription.toJSON();
			await fetch('http://localhost:5181/api/subscribe', {
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
