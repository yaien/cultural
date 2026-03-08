import "@melloware/coloris/dist/coloris.css";
import coloris from "@melloware/coloris";
import { readableColor } from "polished";

class XColorPicker extends HTMLElement {
  color: string = "black";
  swatches: string[] = [];

  connectedCallback() {
    this.color = this.getAttribute("color") || "black";
    this.swatches = JSON.parse(this.getAttribute("swatches") || "[]");

    this.style.backgroundColor = this.color;
    this.style.color = readableColor(this.color);

    coloris.init();

    setTimeout(() => this.init(), 100);
  }

  init() {
    const input = this.querySelector<HTMLInputElement>("input[data-coloris]");
    if (!input) return;
    coloris({ el: input, themeMode: "auto", swatches: this.swatches });

    input.style.color = this.style.color;
    input.style.cursor = "pointer";

    input.addEventListener("input", () => {
      const color = input.value;
      this.style.backgroundColor = color;
      this.style.color = readableColor(color);
      input.style.color = this.style.color;
    });
  }
}

customElements.define("x-color-picker", XColorPicker);
