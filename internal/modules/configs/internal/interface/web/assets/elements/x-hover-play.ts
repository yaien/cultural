class XVideoHoverPlay extends HTMLElement {
    connectedCallback() {
        this.addEventListener("mouseover", () => {
            const video = this.querySelector("video");
            if (!video) return;
            video.play();
        });

        this.addEventListener("mouseout", () => {
            const video = this.querySelector("video");
            if (!video) return;
            video.pause();
        });
    }
}

customElements.define("x-hover-play", XVideoHoverPlay);
