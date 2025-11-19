import { EditorView, basicSetup } from 'codemirror'
import { javascript } from '@codemirror/lang-javascript'
import { css } from '@codemirror/lang-css'


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
            styles: 5,
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
        async render(options = { reset: true }) {
            const res = await fetch("/dashboard/api/render", {
                method: "POST",
                headers: { "content-type": "application/json" },
                body: JSON.stringify({ type: "page", page: this.data }),
            });

            const data = await res.json();

            if (options.reset) {
                this.srcdoc = data.html;
                return;
            }

            this.$refs.iframe.contentDocument.documentElement.innerHTML = data.html;


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

        async edit(content, scope = "body") {
            if (this.data) {
                this.data[scope] = content;
                await this.render({ reset: false });
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

    Alpine.data("mirror", (data, mode = "json") => ({

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
                const content = state.doc.toString();
                this.dispatch(content);
            });

            new EditorView({
                doc: this.format(data),
                extensions: [
                    basicSetup,
                    theme,
                    this.extension(),
                    listener
                ],
                parent: this.$el,
            });
        },

        format(data) {
            switch (mode) {
                case "css":
                    return typeof data === 'string' ? data : '';
                case "json":
                default:
                    return JSON.stringify(data, null, 2);
            }
        },

        extension() {
            switch (mode) {
                case "css":
                    return css();
                case "json":
                    return javascript();
                default:
                    return javascript();
            }
        },

        dispatch(content) {
            try {
                switch (mode) {
                    case "json":
                        const parsedData = JSON.parse(content);
                        this.$dispatch("update", parsedData);
                        break
                    default:
                        this.$dispatch("update", content);
                }
            } catch (e) {
                console.error("Invalid JSON", e);
            }
        },
    }));

    Alpine.data("fonts", () => ({
        current: null,
        fonts: [],
        limit: 20,
        offset: 0,
        family: "",
        ready: false,
        loading: false,
        state: 0,
        selected: { font: null, tag: "" },
        states: {
            initial: 0,
            browsing: 1,
            configuring: 2,
        },

        async init() {
            await Promise.all([this.fetchFonts(), this.fetchCurrent()])
            this.ready = true;
        },

        async fetchFonts() {
            this.loading = true;
            const res = await fetch(`/dashboard/api/fonts?family=${this.family}&limit=${this.limit}&offset=${this.offset}`)
            const fonts = await res.json();
            fonts.forEach(font => this.load(font))
            this.fonts = this.fonts.concat(fonts);
            this.loading = false;
        },

        async fetchCurrent() {
            const res = await fetch("/dashboard/api/fonts/config");
            this.current = (await res.json()) ?? {};
        },

        async load(font) {
            const face = new FontFace(font.family, `url("${font.files.regular}")`, {
                weight: "normal",
                display: 'swap'
            });
            const loaded = await face.load()
            document.fonts.add(loaded);
        },

        async search(family) {
            this.family = family;
            this.offset = 0;
            this.fonts = [];
            await this.fetchFonts();
        },

        async scroll(event) {
            if (event.target.scrollTop + event.target.clientHeight >= event.target.scrollHeight - 100) {
                this.offset += this.limit;
                await this.fetchFonts();
            }
        },

        async add() {
            this.loading = true;
            const res = await fetch("/dashboard/api/fonts/config", {
                method: "PUT",
                headers: { "content-type": "application/json" },
                body: JSON.stringify({
                    ...this.current,
                    [this.selected.tag]: this.selected.font
                }),
            });

            if (res.ok) {
                await this.fetchCurrent();
                this.$dispatch("updated")
                this.state = this.states.initial;
            }

            this.loading = false;
        },

        style(family) {
            return { "font-family": `"${family}", sans-serif` }
        },

        on(state) {
            return this.state == state
        },

        set(state) {
            this.state = state
        },

        select(font, tag = "", state = this.states.configuring) {
            this.selected = { font: font, tag: tag };
            this.state = state;
        },

        empty() {
            return Object.keys(this.current).length == 0;
        }
    }))

});


