class XFont extends HTMLElement {
    async connectedCallback() {
        const family = this.getAttribute("family");
        const url = this.getAttribute("url");

        if (!family || !url) return;

        this.style.opacity = "0";
        this.style.transition = "opacity 0.3s ease-in-out";

        const face = new FontFace(family, `url("${url}")`, {
            weight: "normal",
            display: "swap",
        });
        const loaded = await face.load();
        document.fonts.add(loaded);
        this.style.fontFamily = `"${family}", sans-serif`;
        this.style.opacity = "1";
    }
}

customElements.define("x-font", XFont);
