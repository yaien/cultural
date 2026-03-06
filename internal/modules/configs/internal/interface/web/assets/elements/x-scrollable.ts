class XScrollable extends HTMLElement {
    connectedCallback() {
        if (!this.parentElement) return;

        this.style.overflow = "auto";
        this.style.scrollBehavior = "smooth";
        this.style.display = "block";

        const scale = parseFloat(this.getAttribute("scale") || "1");

        const observer = new ResizeObserver(() => {
            if (!this.parentElement) return;
            this.style.maxHeight = this.parentElement.clientHeight * scale + "px";
        });

        observer.observe(this.parentElement);
    }
}

customElements.define("x-scrollable", XScrollable);
