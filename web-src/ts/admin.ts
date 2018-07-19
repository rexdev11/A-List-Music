let fileInputBusy = true;
let manifest = {};
let currentFile = {};
const MediaTracks  = new Set();
const titleInputRef = document.getElementById("SongTitleInput");
const descInputRef = document.getElementById("SongDescriptionInput");
const fileSelectRef = document.getElementById("FileSelectInput");
const fileSelectLabel = document.getElementById("FileInputLabel");

// Prep new files for transfer
function onFileSelected() {
    if (!fileSelectRef && !fileSelectRef.files) {
        return;
    }
    fileInputBusy = true;

    // Input only selects 1 file, I think...
    currentFile = fileSelectRef.files[0];
    fileSelectLabel.innerText = !!titleInputRef.value
        ? titleInputRef.value
        : 'No Title For File...';
    fileInputBusy = false;
}

function titleChangeHandler(evt) {
    fileSelectLabel.innerText = fileSelectRef.files > 0
        ? evt.value
        : fileSelectLabel.innerText;
}

function initMediaTracks() {
    const RoomName = RoomPrefix + "sound_manifest";
    socket.To(RoomName).Emit("get", mediaResultHandler);

    function mediaResultHandler(evt, err) {
        if (err) {
            console.log('Error: Manifest Result Handler');
            return void 0;
        }
        manifest = evt.value;
        setMedia()
    }

    function setMedia() {
        for (let entry of manifest) {
            MediaTracks.set(entry.id, {
                htmlItem: TEMPLATES.track({
                    artist: entry.artist,
                    name: entry.name,
                    location: entry.location,
                    duration: entry.duration,
                    encoding: entry.length,
                    size: entry.size,
                    url: entry.uri
                }),
                meta: entry
            });
        }
    }
}

function onSubmit() {
    Event.preventDefault();

    const RoomName = "file_upload";
    const description = descInputRef.value;
    const title = titleInputRef.value;

    if (fileInputBusy) {
        alert("Opps! Try submitting again.");
    } else {
        socket.On(RoomName, function() {
            currentFile['meta'] = {
                clientId: localStorage.getItem("client_id"),
                desc: description,
                title: title
            };
            socket.Emit('upload_file', currentFile);
        });
    }
    console.log("handler hit", currentFile);
}

function initializeAdminPanel() {
    const bodyRef = document.getElementsByTagName('body')[0];
    const DATA_IN = JSON.parse(bodyRef.dataset.in);

    console.log(DATA_IN);
}


(function () {
    console.log("AdminJS");
    initializeAdminPanel();
})();