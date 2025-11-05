import { EditorView, basicSetup } from 'codemirror'
import { javascript } from '@codemirror/lang-javascript'


document.addEventListener("alpine:init", () => {
    Alpine.data("pages", ({ url }) => ({
        url: url,
        pages: {},
        data: null,
        page: null,
        ready: false,
        loading: false,
        editing: false,
        state: 0,
        srcdoc: "",
        states: {
            initial: 0,
            create: 1,
            fonts: 2,
            delete: 3,
            editor: 4,
        },

        async init() {
            await this.fetch();
            this.select("index");
            await this.render();
            this.ready = true;
        },
        async fetch() {
            const res = await fetch("/dashboard/api/pages");
            this.pages = await res.json();
        },
        async render() {
            const res = await fetch("/dashboard/api/render", {
                method: "POST",
                headers: { "content-type": "application/json" },
                body: JSON.stringify({ type: "page", body: this.data }),
            });

            const data = await res.json();

            this.srcdoc = data.html;
        },
        async update() {
            this.loading = true;
            const res = await fetch(`/dashboard/api/pages/${this.page}`, {
                method: "PUT",
                headers: { "content-type": "application/json" },
                body: JSON.stringify(this.data),
            });

            if (!res.ok) {
                alert("Error saving page");
                return;
            }

            const updated = await res.json();

            await this.fetch();
            this.select(updated.path);

            this.loading = false;
        },

        async edit(body) {
            if (this.data) {
                this.data.body = body;
                await this.render();
            }
        },

        set(state) {
            this.state = state;
        },

        on(state) {
            return this.state == state;
        },

        select(page) {
            this.page = page;
            this.data = { ...this.pages[page], path: page };
        },
        get pageUrl() {
            if (!this.data) return "";
            let path = this.data.path == "index" ? "" : "/" + this.data.path;
            return this.url + path;
        },
        get pageIsIndex() {
            return this.page == "index";
        },
    }));

    Alpine.data("mirror", (data) => ({

        init() {
            const theme = EditorView.theme({
                "&": {
                    height: "100%",
                    "margin-bottom": "1rem",
                },
                "&.cm-focused": {
                    outline: "none",
                },
            })

            const listener = EditorView.updateListener.of(({ docChanged, state }) => {
                if (!docChanged) return
                try {
                    const content = state.doc.toString();
                    const parsedData = JSON.parse(content);
                    this.$dispatch("update", parsedData);
                } catch (e) {
                    console.error("Invalid JSON", e);
                }
            });

            new EditorView({
                doc: JSON.stringify(data, null, 2),
                extensions: [
                    basicSetup,
                    theme,
                    javascript(),
                    listener
                ],
                parent: this.$el,
            });
        },
    }));
});


