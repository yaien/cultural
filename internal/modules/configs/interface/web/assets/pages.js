document.addEventListener("alpine:init", () => {
    Alpine.data("pages", ({ url }) => ({
        url: url,
        pages: {},
        data: null,
        page: null,
        ready: false,
        loading: false,
        async init() {
            await this.fetch()
            this.select("index")
            this.ready = true
        },
        async fetch() {
            const res = await fetch("/dashboard/api/pages")
            this.pages = await res.json()
        },
        async savePage() {
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

        select(page) {
            this.page = page
            this.data = { ...this.pages[page], path: page }
        },
        get pageUrl() {
            let path = this.data.path == "index" ? "" : "/" + this.data.path
            return this.url + path
        },
        get pageUrlIsDisabled() {
            return this.page == "index"
        }


    }))
})