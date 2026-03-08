import Alpine from "alpinejs";

Alpine.data("preview", (src: string) => ({
    src,
    srcdoc: "",
    loading: false,

    async init() {
        this.loading = true;
        const res = await fetch(this.src);
        this.srcdoc = await res.text();
        this.loading = false;
    },

    async render() {
        this.loading = true;
        const res = await fetch(this.src);
        const iframe = this.$refs.iframe as HTMLIFrameElement;
        if (iframe && iframe.contentDocument) {
            const doc = iframe.contentDocument.open();
            doc.write(await res.text());
            doc.close();
        }
        this.loading = false;
    },
}));
