document.addEventListener("alpine:init", () => {
  Alpine.data("roles", () => ({
    roles: [],
    loading: true,
    init() {
      this.fetch();
    },
    async fetch() {
      try {
        this.loading = true;
        const res = await fetch("/dashboard/api/roles");
        this.roles = (await res.json()) || [];
      } catch (error) {
        console.error(error);
      }
    },
  }));
});
