
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
        states: {
            initial: 0,
            create: 1,
            fonts: 2,
            delete: 3,
            editor: 4
        },

        async init() {
            await this.fetch()
            this.select("index")
            this.$refs.textarea
            this.ready = true
        },
        async fetch() {
            const res = await fetch("/dashboard/api/pages")
            this.pages = await res.json()
        },
        async update() {
            this.loading = true
            const res = await fetch(`/dashboard/api/pages/${this.page}`, {
                method: "PUT",
                headers: { "content-type": "application/json" },
                body: JSON.stringify(this.data)
            })

            if (!res.ok) {
                alert("Error saving page")
                return
            }

            const updated = await res.json()

            await this.fetch()
            this.select(updated.path)

            this.loading = false
        },

        async edit(body) {
            this.data.body = body
            await this.update()
            this.$refs.iframe.contentWindow.location.reload()
        },

        set(state) {
            this.state = state
        },

        on(state) {
            return this.state == state
        },

        select(page) {
            this.page = page
            this.data = { ...this.pages[page], path: page }
        },
        get pageUrl() {
            let path = this.data.path == "index" ? "" : "/" + this.data.path
            return this.url + path
        },
        get pageIsIndex() {
            return this.page == "index"
        },
        get pageBody() {
            return JSON.stringify(this.data.body, null, 2)
        }


    }))


    Alpine.data("mirror", (data) => ({
        editor: null,
        init() {
            console.log("Initializing CodeMirror", this.$refs.textarea)
            this.editor = CodeMirror.fromTextArea(this.$refs.textarea, {
                mode: { name: "javascript", json: true },
                tabSize: 2,
                parent: this.$el
            });

            this.editor.setValue(JSON.stringify(data, null, 2))

            this.editor.on("change", () => {
                try {
                    let data = JSON.parse(this.editor.getValue())
                    this.$dispatch("update", data)
                } catch (e) {
                    console.error("Invalid JSON", e)
                }
            })

        }
    }))
})