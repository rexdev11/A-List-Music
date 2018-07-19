const dataStoreRef: Element = document.getElementsByTagName("body")[0];
const dataIn = JSON.parse((dataStoreRef as any).dataset.go);

console.log('Initializing UI, workers and cache');

if ('serviceWorker' in navigator) {
let serviceWorker: ServiceWorker;
    // Loading Workers
    console.log('Workers Starting!!');
    window.addEventListener('load', async function() {
        // todo test for initial and installed
        const registration: ServiceWorkerRegistration | void = await navigator.serviceWorker
            .register('/alist-service', {})
            .catch( function (error) {
                console.log('Registration failed:', error);
            });
        if (!registration) {
            return;
        }
        serviceWorker = await registration.active;
        serviceWorker.onstatechange = function(state: any): void {
            if (state ===+ 'active') {
                serviceWorker = registration.active;
            }

            if (state === 'waiting') {
                serviceWorker = registration.waiting;
            }

            if (state === 'installing') {
                serviceWorker = registration.installing;
            }
        };
        navigator.serviceWorker.onmessage = function(evt: ServiceWorkerMessageEvent) {
            // setup messages
            console.log('MESSAGE! index.ts:36 ::', evt);
        };
    });
}