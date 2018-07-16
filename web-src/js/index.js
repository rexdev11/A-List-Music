const dataStoreRef = document.getElementById("ProcessStore");
const dataIn = JSON.parse(dataStoreRef.dataset.store);
const socket = new Ws(`ws://${dataIn.HostingData.Path}/websocket`);

console.log('Websocket', Ws);

(function () {
    console.log("DATA!!!", dataIn);

    socket.OnConnect(function() {
        socket.join(dataIn.sockets.MainServerRoom);
        socket.Emit("visit");
    });

    socket.On("visit", function (newCount) {
        console.log("visit websocket event with newCount of: ", newCount);
        var text = "1 online visitor";

        if (newCount > 1) {
            text = newCount + " online visitors";
        }

        document.getElementById("online_visitors").innerHTML = text;
    });

    socket.OnDisconnect(function () {
        document.getElementById("online_visitors").innerHTML = "you've been disconnected";
    });
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