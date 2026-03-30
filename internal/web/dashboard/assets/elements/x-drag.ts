import Alpine from "alpinejs";

Alpine.data("drag", () => ({
  init() {
    this.$el.addEventListener("dragover", (ev) => {
      ev.preventDefault();
      this.$el.classList.add("dragover");
    });

    this.$el.addEventListener("dragstart", (e) => {
      const id = (e.target as Element).id;
      e.dataTransfer?.setData("text/plain", id);
    });

    this.$el.addEventListener("drop", (ev) => {
      ev.preventDefault();
      this.$el.classList.remove("dragover");
      const id = ev.dataTransfer?.getData("text/plain");
      if (id) {
        this.$dispatch("dropped", { id });
      }
    });
  },
}));
