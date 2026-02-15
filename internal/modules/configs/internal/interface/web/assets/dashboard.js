document.addEventListener("alpine:init", () => {
  Alpine.data("dashboard", () => ({
    sidebar: false,
    toast: null,

    bindings: {
      [":class"]() {
        return { open: this.sidebar };
      },

      ["@sidebar"](event) {
        this.sidebar = event.detail;
      },

      ["@toast"](event) {
        cleanup = () => {
          this.toast.leave = true;
          setTimeout(() => {
            this.toast = null;
          }, 500);
        };

        // If no message is provided, close the current toast
        if (!event.detail.message && this.toast != null) {
          clearTimeout(this.toast.timeout);
          cleanup();
          return;
        }

        // Show new toast message
        this.toast = event.detail;
        this.toast.timeout = setTimeout(cleanup, 5000);
      },
    },
  }));
});
