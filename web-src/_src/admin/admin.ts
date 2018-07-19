import {TEMPLATES} from "../templates/templates";
import {MediaMeta} from "./admin.model";

declare const window: Window;
declare const Event: Event;
declare const dataIn: any;
declare const socket: any;

let fileInputBusy = true;
let manifest: {
    [id: string]: MediaMeta;
} = {};
let currentFile = {};

export interface InputEle extends Element {

}

const MediaTracks  = new Map<string, {meta: MediaMeta, htmlItem: (props: any) => {}}>();
const titleInputRef: Element = document.getElementById("SongTitleInput");
const descInputRef: Element = document.getElementById("SongDescriptionInput");
const fileSelectRef: Element = document.getElementById("FileSelectInput");
const fileSelectLabel: HTMLElement = document.getElementById("FileInputLabel");

// Prep new files for transfer
function onFileSelected() {
    if (!fileSelectRef && !fileSelectRef.hasAttribute('files')) {
        return;
    }
    fileInputBusy = true;

    // Input only selects 1 file, I think...
    currentFile = (fileSelectRef as any).files[0];
    fileSelectLabel.innerText = titleInputRef.hasAttribute("value")
        ? (titleInputRef as any).value
        : 'No Title For File...';
    fileInputBusy = false;
}

function titleChangeHandler(evt): void {
    fileSelectLabel.innerText = (fileSelectRef as any).files.length > 0
        ? evt.value
        : fileSelectLabel.innerText;
}

function initMediaTracks(): void {
    const RoomName: string = dataIn.RoomPrefix + "sound_manifest";
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
        for (let entryId in manifest) {
            MediaTracks.set(entryId, {
                htmlItem: TEMPLATES.track({
                    artist: manifest[entryId].artist,
                    name: manifest[entryId].name,
                    location: manifest[entryId].location,
                    duration: manifest[entryId].duration,
                    encoding: manifest[entryId].length,
                    size: manifest[entryId].size,
                    url: manifest[entryId].url
                }),
                meta: manifest[entryId]
            });
        }
    }
}

function onSubmit() {
    Event.preventDefault();

    const RoomName = "file_upload";
    const description = (descInputRef as any ).value;
    const title = (titleInputRef as any).value;

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