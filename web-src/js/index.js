const dataStoreRef = document.getElementsByTagName("body")[0];
const dataIn = JSON.parse(dataStoreRef.dataset.go.in);


(function () {

})();

console.log('Initializing UI, workers and cache');
if ('serviceWorker' in navigator) {
    window.addEventListener('load', async function () {
        // todo test for initial and installed

        let serviceWorker = navigator.serviceWorker;

        const registration = await navigator.serviceWorker
            .register('/alist-service', {})
            .catch( function (error) {
                console.log('Registration failed:', error);
            });

        navigator.serviceWorker.ready.then(function (evt) {
            console.log("wtf", evt);
        });

        registration.onstatechange = function (evt) {
            if (evt === "active") {
                serviceWorker = registration.active;
            }
            if (evt === "waiting") {
                serviceWorker = registration.waiting;
            }
            if (evt === "installing") {
                serviceWorker = registration.installing;
            }
        }
    });
}