import { readableColor } from "polished";
import type * as MonacoEditor from "monaco-editor";

declare var require: any;
declare var monaco: typeof MonacoEditor;
declare var Coloris: any;

// Coloris Docs  https://github.com/mdbassit/Coloris/blob/main/README.md#customizing-the-color-picker
function ColorPicker() {
  Coloris({ formatToggle: true, alpha: false });
}

function ReadableColor() {
  document.querySelectorAll<HTMLElement>("[data-readable-color]").forEach((el) => {
    const color = el.dataset.readableColor;
    if (color) {
      el.style.color = readableColor(color);
    }
  });
}

function Monaco() {
  const container = document.querySelector<HTMLElement>("[data-editor-container]");
  if (!container) {
    console.error("No container [data-editor-container] found for Monaco Editor");
    return;
  }

  require.config({ paths: { vs: "/assets/static/dashboard/dist/monaco" } });

  const install = (el: HTMLElement) => {
    const observer = new ResizeObserver(() => {
      el.style.height = container.clientHeight * 0.85 + "px";
    });

    observer.observe(container);

    const editor = monaco.editor.create(el, {
      value: el.dataset.value || "",
      language: el.dataset.language || "html",
      minimap: { enabled: false },
      automaticLayout: true,
    });

    editor.onDidChangeModelContent(() => {
      const detail = {
        value: editor.getValue(),
        language: editor.getModel()?.getLanguageId(),
        editor: el.dataset.editor,
      };
      el.dispatchEvent(new CustomEvent("editor-change", { detail }));

      console.debug("Editor content changed", detail);
    });
  };

  require(["vs/editor/editor.main"], () => {
    document.querySelectorAll<HTMLElement>("[data-editor]").forEach(install);
  });
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
  ColorPicker();
  ReadableColor();
  Monaco();
  HoverPlay();
}

init();

document.addEventListener("htmx:load", init);
