import { readableColor } from "polished";
import type * as MonacoEditor from "monaco-editor";

declare var require: any;
declare var monaco: typeof MonacoEditor;
declare var Coloris: any;

// Coloris Docs  https://github.com/mdbassit/Coloris/blob/main/README.md#customizing-the-color-picker
function ColorPicker() {
  Coloris({ formatToggle: true, alpha: false, clearButton: true });
}

class ReadableColor extends HTMLElement {
  connectedCallback() {
    const color = this.getAttribute("color");
    if (color) {
      this.style.color = readableColor(color);
    }

    this.infect(this.childNodes);
  }

  infect(nodes: NodeListOf<ChildNode>) {
    nodes.forEach((node) => {
      if (node.nodeType === Node.ELEMENT_NODE) {
        const element = node as HTMLElement;
        element.style.color = this.style.color;
        this.infect(element.childNodes);
      }
    });
  }
}

class Editor extends HTMLElement {
  mount: HTMLDivElement;

  constructor() {
    super();
    this.mount = document.createElement("div");
  }

  connectedCallback() {
    this.style.display = "block";
    this.style.height = "100%";
    this.style.width = "100%";

    this.appendChild(this.mount);
    const observer = new ResizeObserver(() => {
      this.mount.style.height = this.clientHeight * 0.99 + "px";
    });

    observer.observe(this);

    require.config({ paths: { vs: "/assets/static/dashboard/dist/monaco" } });

    require(["vs/editor/editor.main"], () => {
      const editor = monaco.editor.create(this.mount, {
        value: this.getAttribute("value") || "",
        language: this.getAttribute("language") || "html",
        minimap: { enabled: false },
        automaticLayout: true,
      });

      editor.onDidChangeModelContent(() => {
        const detail = {
          value: editor.getValue(),
          language: editor.getModel()?.getLanguageId(),
          editor: this.getAttribute("editor"),
        };
        this.dispatchEvent(new CustomEvent("input", { detail, bubbles: true }));
      });
    });
  }
}

export function HoverPlay() {
  document.querySelectorAll<HTMLElement>("[data-hover-play]").forEach((el) => {
    const video = el.querySelector("video");
    if (!video) return;
    el.addEventListener("mouseenter", () => video.play());
    el.addEventListener("mouseleave", () => video.pause());
  });
}

export function init() {
  customElements.get("x-readable-color") || customElements.define("x-readable-color", ReadableColor);
  customElements.get("x-editor") || customElements.define("x-editor", Editor);
  ColorPicker();
  HoverPlay();
}

init();

document.addEventListener("htmx:load", init);
