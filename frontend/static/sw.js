self.addEventListener('push', function(event) {
    if (event.data) {
        try {
            const payload = event.data.json();
            
            const options = {
                body: payload.body || 'Tidak ada deskripsi.',
                icon: '/icons/lokari-logo.png', // Fallback icon
                badge: '/icons/lokari-badge.png',
                vibrate: [200, 100, 200, 100, 200, 100, 200], // SOS vibration
                requireInteraction: true,
                data: {
                    url: payload.url || '/'
                }
            };
            
            event.waitUntil(
                self.registration.showNotification(payload.title || 'Peringatan LOKARI', options)
            );
        } catch(e) {
            console.error("Push payload error:", e);
        }
    }
});

self.addEventListener('notificationclick', function(event) {
    event.notification.close();
    if (event.notification.data && event.notification.data.url) {
        event.waitUntil(
            clients.openWindow(event.notification.data.url)
        );
    }
});
