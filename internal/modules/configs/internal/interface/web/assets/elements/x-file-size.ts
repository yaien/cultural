import { filesize } from "filesize";

class XFileSize extends HTMLElement {
    connectedCallback() {
        const size = this.getAttribute("size");
        if (!size) return;
        this.innerHTML = filesize(size);
    }
}

customElements.define("x-file-size", XFileSize);
