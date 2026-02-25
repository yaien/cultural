import { readableColor } from "polished";
import { filesize } from "filesize";

require.config({ paths: { vs: "/assets/static/dashboard/dist/monaco" } });

document.addEventListener("alpine:init", () => {
    Alpine.data("pages", ({ url, filepath }) => ({
        filepath: filepath,
        url: url,
        draft: null,
        model: { map: "", key: "", value: {} },
        ready: false,
        loading: false,
        state: 0,
        states: {
            initial: 0,
            create: 1,
            fonts: 2,
            delete: 3,
            files: 4,
            editor: 5,
            styles: 6,
            colors: 7,
            script: 8,
            publish: 9,
        },

        async init() {
            await this.fetch();
            await this.select("index", "pages");
            this.ready = true;
        },

        async fetch() {
            const res = await fetch("/dashboard/api/draft");
            this.draft = await res.json();
        },

        async select(key, map) {
            this.model = { map, key, value: { ...this.draft[map][key] } };
            await this.render();
        },

        async render(options = { reset: true, wait: 0 }) {
            this.loading = true;

            const res = await fetch("/dashboard/api/render", {
                method: "POST",
                headers: { "content-type": "application/json" },
                body: JSON.stringify({
                    map: this.model.map,
                    key: this.model.key,
                    pages: this.draft.pages,
                    layouts: this.draft.layouts,
                    emails: this.draft.emails,
                    colors: this.draft.colors,
                    fonts: this.draft.fonts,
                }),
            });

            const data = await res.json();

            this.loading = false;

            if (options.wait) {
                this.srcdoc = null;
                await new Promise((resolve) => setTimeout(resolve, options.wait));
            }

            if (options.reset) {
                this.srcdoc = data.html;
                return;
            }

            this.$refs.iframe.contentDocument.documentElement.innerHTML = data.html;
        },

        async update({ toast, draft, model, wait, reset } = { toast: true, reset: false, wait: 0, draft: null, model: null }) {
            this.loading = true;

            if (draft) {
                this.draft = draft;
            }

            if (model) {
                this.model = model;
            }



            const res = await fetch("/dashboard/api/draft", {
                method: "PUT",
                headers: { "content-type": "application/json" },
                body: JSON.stringify({
                    pages: this.draft.pages,
                    layouts: this.draft.layouts,
                    emails: this.draft.emails,
                    colors: this.draft.colors,
                    fonts: this.draft.fonts,
                }),
            });

            if (!res.ok) {
                this.$dispatch("toast", { message: "Error inesperado al guardar sitios", variant: "danger" });
                return;
            }

            if (toast) {
                this.$dispatch("toast", { message: "Cambios guardados" });
            }

            this.loading = false;

            await this.render({ reset, wait });
        },

        set(state) {
            this.state = state;
        },

        on(state) {
            return this.state == state;
        },

        get deleteable() {
            if (!this.model) return false;
            return (this.model.map == "layouts" && this.model.key != "default") || (this.model.map == "pages" && this.model.key != "index");
        },
        get forWeb() {
            if (!this.model) return false;
            return this.model.map == "pages" || this.model.map == "layouts";
        },
    }));

    Alpine.data("basic", ({ draft, model }) => ({
        draft: draft,
        model: model,

        get pageUrl() {
            if (this.model.map != "pages") return "";
            let path = this.model.key == "index" ? "" : "/" + this.model.key;
            return this.url + path;
        },
        get pageIsIndex() {
            return this.model.map == "pages" && this.model.key == "index";
        },
        get layoutIsDefault() {
            return this.model.map == "layouts" && this.model.key == "default";
        },

        submit() {
            this.draft[this.model.map][this.model.value.name] = this.model.value;
            if (this.model.key != this.model.value.name) {
                delete this.draft[this.model.map][this.model.key];
                this.model.key = this.model.value.name;
            }
            this.$dispatch("update", { draft: this.draft, model: this.model, toast: true });
        },

        changeTemplate(key) {
            this.model.key = key;
            this.model.value = this.draft[this.model.map][key];
            this.$dispatch("model", this.model);
        },

        changeMap(map) {
            switch (map) {
                case "pages":
                    this.model.map = "pages";
                    this.model.key = "index";
                    this.model.value = this.draft.pages.index;
                    this.$dispatch("model", this.model);
                    break;

                case "emails":
                    this.model.map = "emails";
                    this.model.key = "invitation";
                    this.model.value = this.draft.emails.invitation;
                    this.$dispatch("model", this.model);
                    break;

                case "layouts":
                    this.model.map = "layouts";
                    this.model.key = "default";
                    this.model.value = this.draft.layouts.default;
                    this.$dispatch("model", this.model);
                    break;
            }
        },
    }));

    Alpine.data("create", ({ draft, model }) => ({
        draft: draft,
        model: model,
        form: {
            map: "pages",
            value: {},
        },

        submit() {
            this.draft[this.form.map][this.form.value.name] = this.form.value;
            this.model.map = this.form.map;
            this.model.key = this.form.value.name;
            this.model.value = this.form.value;
            this.$dispatch("update", { draft: this.draft, model: this.model, toast: true });
            this.$dispatch("submitted");
        },

        get disabled() {
            return this.form && this.draft[this.form.map][this.form.value.name];
        },
    }));

    Alpine.data("remove", ({ draft, model }) => ({
        draft: draft,
        model: model,

        remove() {
            delete this.draft[this.model.map][this.model.key];
            switch (this.model.map) {
                case "layouts":
                    if (this.draft.layouts?.length) {
                        this.model.key = this.draft.layouts[0];
                        this.model.value = this.draft.layouts[this.model.key];
                        break;
                    }
                default:
                    this.model.key = "index";
                    this.model.value = this.draft.pages.index;
                    break;
            }
            this.$dispatch("update", { draft: this.draft, model: this.model, toast: true });
            this.$dispatch("removed");
        },
    }));

    Alpine.data("monaco", ({ draft, model, language, reset }) => ({
        draft: draft,
        model: model,
        language: language,

        init() {
            require(["vs/editor/editor.main"], () => {


                const container = document.querySelector("[data-editor]");
                const observer = new ResizeObserver(() => {
                    this.$el.style.height = container.clientHeight * 0.85 + "px";
                })

                observer.observe(container);


                monaco.editor.setTheme("vs");

                const editor = monaco.editor.create(this.$el, {
                    value: this.doc(),
                    language: this.language,
                    minimap: { enabled: false },
                    automaticLayout: true,
                });

                editor.getModel().onDidChangeContent(event => {
                    const content = editor.getValue();
                    if (content != this.doc()) {
                        this.dispach(content);
                    }
                })


            });
        },

        dispach(content) {
            switch (this.language) {
                case "css":
                    this.model.value.styles = content;
                    break;
                case "html":
                    this.model.value.body = content;
                    break;
                case "javascript":
                    this.model.value.script = content;
                    break;
            }

            this.draft[this.model.map][this.model.key] = this.model.value;
            this.$dispatch("update", { draft: this.draft, model: this.model, toast: false, reset });
        },

        doc() {
            switch (this.language) {
                case "css":
                    return this.model.value.styles || "";
                case "html":
                    return this.model.value.body || "";
                case "javascript":
                    return this.model.value.script || "";
            }
        },
    }));

    Alpine.data("fonts", ({ draft }) => ({
        draft: { ...draft },
        fonts: [],
        limit: 10,
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
            await Promise.all([this.fetchFonts(), this.initFonts()]);
            this.ready = true;
        },

        async initFonts() {
            if (this.draft.fonts) {
                for (const tag in this.draft.fonts) {
                    const font = this.draft.fonts[tag];
                    await this.load(font);
                }
            }
        },

        async fetchFonts() {
            this.loading = true;
            const res = await fetch(
                `/dashboard/api/fonts?family=${this.family}&limit=${this.limit}&offset=${this.offset}`,
            );
            const fonts = await res.json();
            if (fonts) {
                fonts.forEach((font) => this.load(font));
                this.fonts = this.fonts.concat(fonts);
            }
            this.loading = false;
        },

        async load(font) {
            const face = new FontFace(font.family, `url("${font.files.regular}")`, {
                weight: "normal",
                display: "swap",
            });
            const loaded = await face.load();
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
            if (!this.draft.fonts) {
                this.draft.fonts = {};
            }
            this.draft.fonts[this.selected.tag] = this.selected.font;
            this.$dispatch("update", { draft: this.draft, toast: true });
            this.set(this.states.initial);
        },

        style(family) {
            return { "font-family": `"${family}", sans-serif` };
        },

        on(state) {
            return this.state == state;
        },

        set(state) {
            this.state = state;
        },

        select(font, tag = "", state = this.states.configuring) {
            this.selected = { font: font, tag: tag };
            this.state = state;
        },

        empty() {
            return Object.keys(this.draft.fonts || {}).length == 0;
        },
    }));

    Alpine.data("files", () => ({
        files: [],
        ready: false,
        loading: false,
        selected: null,
        data: {},
        state: 0,
        states: {
            initial: 0,
            edit: 1,
        },

        async init() {
            await this.fetch();
            this.ready = true;
        },

        async fetch() {
            this.loading = true;
            const res = await fetch("/dashboard/api/files");
            this.files = (await res.json()) || [];
            this.loading = false;
        },

        async upload(event) {
            this.loading = true;

            for (const file of event.target.files) {
                const data = new FormData();
                data.append("file", file);

                const res = await fetch("/dashboard/api/files", {
                    method: "POST",
                    body: data,
                });

                if (res.ok) {
                    this.files.push(await res.json());
                } else {
                    alert("Error uploading file");
                }
            }

            this.loading = false;
            this.$dispatch("toast", { message: "Archivo(s) subido(s)" });
            this.$dispatch("updated");
        },

        async update() {
            try {
                this.loading = true;
                const res = await fetch(`/dashboard/api/files/${this.selected.name}`, {
                    method: "PUT",
                    body: JSON.stringify({ newName: this.data.name }),
                    headers: { "Content-Type": "application/json" },
                });

                if (!res.ok) {
                    const data = await res.json();
                    throw new Error(data.error);
                }

                this.selected.name = this.data.name;
                this.set(this.states.initial);
                this.$dispatch("toast", { message: "Archivo actualizado" });
                this.$dispatch("updated");
            } catch (error) {
                console.error(error);
            } finally {
                this.loading = false;
            }
        },

        async drop() {
            try {
                this.loading = true;
                const res = await fetch(`/dashboard/api/files/${this.selected.name}`, {
                    method: "DELETE",
                });

                if (!res.ok) {
                    const data = await res.json();
                    throw new Error(data.error);
                }

                this.files = this.files.filter((file) => file.id != this.selected.id);
                this.set(this.states.initial);
                this.$dispatch("toast", { message: "Archivo eliminado" });
                this.$dispatch("updated");
            } catch (error) {
                console.error(error);
            } finally {
                this.loading = false;
            }
        },

        select(file) {
            this.selected = file;
            this.data = { name: file.name };
            this.state = this.states.edit;
        },

        on(state) {
            return this.state == state;
        },

        set(state) {
            this.state = state;
        },

        enter($event) {
            const video = $event.target.querySelector("video");
            video?.play();
        },

        leave($event) {
            const video = $event.target.querySelector("video");
            video?.pause();
        },

        filesize(size) {
            return filesize(size);
        },
    }));

    Alpine.data("colors", ({ draft }) => ({
        draft: draft,
        init() {
            Coloris({
                formatToggle: true,
                alpha: false,
            });
        },
        readable(color) {
            return readableColor(color, "#000", "#fff");
        },
        changeKey(oldKey, newKey) {
            if (oldKey == newKey) return;
            this.draft.colors[newKey] = this.draft.colors[oldKey];
            delete this.draft.colors[oldKey];
            this.$dispatch("update", { draft: this.draft, toast: false });
        },
        changeColor(key, color) {
            this.draft.colors[key] = color;
            if (color === "") {
                delete this.draft.colors[key];
            }
            this.$dispatch("update", { draft: this.draft, toast: false });
        },
        add() {
            let index = 1;
            let key = "color-" + index;
            while (this.draft.colors[key]) {
                key = "color-" + index++;
            }
            this.draft.colors[key] = "#000000";
            this.$dispatch("update", { draft: this.draft, toast: false });
        },
    }));

    Alpine.data("publish", () => ({
        loading: false,

        async publish() {
            this.loading = true;

            const res = await fetch("/dashboard/api/draft/commit", {
                method: "POST",
                headers: { "content-type": "application/json" },
            });

            if (!res.ok) {
                this.$dispatch("toast", { message: "Error inesperado al publicar sitios", variant: "danger" });
                return;
            }

            this.loading = false;

            this.$dispatch("toast", { message: "La configuracion ha sido publicada", variant: "success" });
            this.$dispatch("published");
        },
    }));
});
