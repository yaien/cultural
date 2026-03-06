class XFileInput extends HTMLElement {
    connectedCallback() {
        const name = this.getAttribute("name") || "file";
        const accept = this.getAttribute("accept") || "*/*";
        const multiple = this.getAttribute("multiple");

        this.innerHTML = `
          ${this.innerHTML}
          <input hidden name=${name} accept=${accept} ${multiple !== null ? "multiple" : ""} type="file" />
          <div class="progress-bar" hidden>
            <div class="progress-bar-value"></div>
          </div>
          `;

        const input = this.querySelector("input")!;
        const progress = this.querySelector<HTMLElement>(".progress-bar")!;
        const value = this.querySelector<HTMLElement>(".progress-bar-value")!;

        this.addEventListener("click", () => input.click());

        this.closest("form")?.addEventListener("htmx:xhr:progress", (ev) => {
            const custom = ev as CustomEvent<{ loaded: number; total: number }>;
            const percent = (custom.detail.loaded / custom.detail.total) * 100;

            progress.hidden = false;
            value.style.width = `${percent}%`;
            if (percent >= 100) progress.hidden = true;
        });
    }
}

customElements.define("x-file-input", XFileInput);
