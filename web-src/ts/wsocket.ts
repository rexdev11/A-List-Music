declare var  socket: any;
declare var Ws: any;
declare var scheme: any;
declare var DATA_IN: any;
let port;
let scheme;
let wsURL;
let text = "1 online visitor";

export function setWS() {
    if (document.location && document.location.protocol) {

    // Set up socket URL
    scheme = document.location.protocol
        === "https:"
            ? "wss"
            : "ws"
    } else {
        port = document.location.port
            ? (":" + document.location.port)
            : ""
    }

    const wsURL = scheme + "://" + dataIn.Paths.Host + port + "/websocket";
        console.log("wsURL", wsURL);
    return wsURL;
}

export async function runAdminSocket() {
    const socket = new Ws(wsURL);
    const ServerRoom = DATA_IN.ServerRoomName + ':Main';

    socket.OnConnect(function() {
        socket.Join(ServerRoom);

        // update the rest of connected clients, including "myself" when "my" connection is 100% ready.
        socket.Emit("visit", {

        });
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

    });
}


export async function runIndexSocket() {
    console.log('Websocket', Ws);
    console.log("DATA!!!", dataIn);

    socket.OnConnect(function() {
        socket.Join(DATA_IN.ServerRoomName);
        socket.Emit("visit");
    });

    socket.On("visit", function (newCount) {
        if (newCount > 1) {
            text = newCount + " online visitors";
        }

        document.getElementById("online_visitors").innerHTML = text;
    });

    socket.OnDisconnect(function () {
        document.getElementById("online_visitors").innerHTML = "you've been disconnected";
    });

    console.log('Initializing UI, workers and cache');

    let worker: ServiceWorker;
    if ('serviceWorker' in navigator) {

        window.addEventListener('/load', async function () {
            // todo test for initial and installed

            let serviceWorkerContainer: ServiceWorkerContainer = navigator.serviceWorker;

            const registration: ServiceWorkerRegistration | void = <ServiceWorkerRegistration> await serviceWorkerContainer
                .register('/alist-service')
                .then(registration => {
                    if (registration.installing) {
                        worker = registration.installing;
                    }
                    if (registration.active) {
                        worker = registration.active;
                    }
                    if (registration.waiting) {
                        worker = registration.waiting;
                    }
                    worker.addEventListener('statechange', (state: any) => {
                        if (state === 'installing'){
                            console.log('installing on reg worker');
                            // worker.waitUntil(worker.preCache(['footer', 'sounds']))
                            worker.postMessage({
                                type: 'TEST'
                            });
                        }
                    });
            });
        });
    }
}