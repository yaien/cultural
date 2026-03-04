function NavToggler() {
  document.querySelectorAll<HTMLElement>("[data-toggle]").forEach((el) => {
    el.addEventListener("click", () => {
      document.querySelector(el.dataset.toggle ?? "body")?.classList.toggle(el.dataset.class ?? "open");
    });
  });
}

export function init() {
  NavToggler();
}

init();

document.addEventListener("htmx:load", init);
