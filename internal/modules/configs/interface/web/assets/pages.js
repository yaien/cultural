import { EditorView, basicSetup } from "codemirror";
import { css } from "@codemirror/lang-css";
import { html } from "@codemirror/lang-html";

document.addEventListener("alpine:init", () => {
  Alpine.data("pages", ({ url, filepath }) => ({
    filepath: filepath,
    url: url,
    pages: {},
    form: {},
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
      files: 4,
      editor: 5,
      styles: 6,
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

    async create() {
      try {
        this.loading = true
        const res = await fetch(`/dashboard/api/pages`, {
          method: "POST",
          headers: { "content-type": "application/json" },
          body: JSON.stringify(this.form),
        });

        if (!res.ok) {
          const data = await res.json()
          throw Error(data.error)
        }

        const page = { ...this.form, styles: "", body: "" }
        this.pages = { ...this.pages, [page.name]: page }
        this.select(page.name)
      } catch (err) {
        this.console.log(err)
      } finally {
        this.loading = false
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

  Alpine.data("mirror", (data, mode = "html") => ({
    init() {
      const theme = EditorView.theme({
        "&": {
          height: "100%",
          "margin-bottom": "1rem",
        },
        "&.cm-focused": {
          outline: "none",
        },
      });

      const listener = EditorView.updateListener.of(({ docChanged, state }) => {
        if (!docChanged) return;
        const content = state.doc.toString();
        this.$dispatch("update", content);
      });

      new EditorView({
        doc: this.format(data),
        extensions: [basicSetup, theme, this.extension(), listener],
        parent: this.$el,
      });
    },

    format(data) {
      return typeof data === "string" ? data : "";
    },

    extension() {
      switch (mode) {
        case "css":
          return css();
        case "html":
          return html();
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
      await Promise.all([this.fetchFonts(), this.fetchCurrent()]);
      this.ready = true;
    },

    async fetchFonts() {
      this.loading = true;
      const res = await fetch(`/dashboard/api/fonts?family=${this.family}&limit=${this.limit}&offset=${this.offset}`);
      const fonts = await res.json();
      fonts.forEach((font) => this.load(font));
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
      this.loading = true;
      const res = await fetch("/dashboard/api/fonts/config", {
        method: "PUT",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({
          ...this.current,
          [this.selected.tag]: this.selected.font,
        }),
      });

      if (res.ok) {
        await this.fetchCurrent();
        this.$dispatch("updated");
        this.state = this.states.initial;
      }

      this.loading = false;
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
      return Object.keys(this.current).length == 0;
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
  }));
});
