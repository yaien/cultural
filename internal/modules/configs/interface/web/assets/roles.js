document.addEventListener("alpine:init", () => {
  Alpine.data("members", () => ({
    members: [],
    loading: true,
    init() {
      console.log("Component 38654706270 initialized");
    },
    message: "Hello from component 38654706270!",
  }));
});
