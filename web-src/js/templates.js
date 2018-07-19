const TEMPLATES = {
    track(props) {
        return `
                <div class="list-item">
                    <audio class="track">
                        <div class="left">
                            <label class="label-artist">
                                ${props.artist}
                            </label>
                            <label class="label-name">
                                ${props.name}
                            </label>
                        </div>
                        <div class="center">
                            <label class="label-location">
                                ${props.location}
                            </label>
                            <input role="media" class="scroll" type="range" />
                        </div>
                        <div class="right">
                            <i class="dragable-handle">x</i>
                        </div>
                    </audio>
                </div>
            `
    }
};