class XToggler extends HTMLElement {
    private untoggled = "<slot></slot>";
    private toggled = "<slot name='toggled'></slot>";

    connectedCallback() {
        const target = this.getAttribute("target");
        if (!target) {
            return;
        }

        const el = document.querySelector(target);
        if (!el) {
            return;
        }

        const token = this.getAttribute("toggle") || "";
        const root = this.attachShadow({ mode: "open" });
        root.innerHTML = this.untoggled;

        this.addEventListener("click", () => {
            const toggled = el.classList.toggle(token);
            root.innerHTML = toggled ? this.toggled : this.untoggled;
        });
    }
}

customElements.define("x-toggler", XToggler);
