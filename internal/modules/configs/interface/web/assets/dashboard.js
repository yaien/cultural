document.addEventListener("alpine:init", () => {
  Alpine.data("dashboard", () => ({
    sidebarIsOpen: false,
    hideSidebar() {
      this.sidebarIsOpen = false;
    },
    toggleSidebar() {
      this.sidebarIsOpen = !this.sidebarIsOpen;
    },
  }));
});
