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
            this.ready = true
        },
        async fetch() {
            const res = await fetch("/dashboard/api/pages")
            this.pages = await res.json()
            this.select("index")
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