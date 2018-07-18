text = "1 online visitor";

function setWS() {
    if (document.location && document.location.protocol) {

        // Set up socket URL
        const scheme = document.location.protocol
            === "https:"
                ? "wss"
                : "ws"
    } else {
        const port = document.location.port
            ? (":" + document.location.port)
            : ""
    }

    const wsURL = scheme + "://" + dataIn.Paths.Host + port + "/websocket";
        console.log("wsURL", wsURL);
    return wsURL;
}

async function runAdminSocket() {
        const socket = new Ws(wsURL);
        const ServerRoom = DATA_IN.ServerRoomName + ':Main';

        socket.OnConnect( function() {
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


async function runIndexSocket() {
    console.log('Websocket', Ws);
    console.log("DATA!!!", dataIn);

    socket.OnConnect(function() {
        socket.Join(dataIn.sockets.MainServerRoom);
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

    if ('serviceWorker' in navigator) {

        window.addEventListener('/load', async function () {
            // todo test for initial and installed

            let serviceWorker = navigator.serviceWorker;

            const registration = await navigator.serviceWorker
                .register('/alist-service', {
                    foo: 'dataBar'
                }).catch(function (error) {
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
}