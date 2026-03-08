import loader from "@monaco-editor/loader";
import type Monaco from "monaco-editor";

class XCodeEditor extends HTMLElement {
  async connectedCallback() {
    const language = this.getAttribute("language") || "javascript";
    const value = this.getAttribute("value") || "";
    const monaco = (await loader.init()) as typeof Monaco;

    this.innerHTML = "";
    this.style.height = this.parentElement ? `${this.parentElement.clientHeight * 0.99}px` : "300px";
    this.style.width = "100%";
    this.style.display = "block";

    const editor = monaco.editor.create(this, {
      value,
      language,
      automaticLayout: true,
      minimap: { enabled: false },
      lineNumbersMinChars: 1,
      scrollbar: {
        vertical: "hidden",
      },
    });
    editor.onDidChangeModelContent(() => {
      const value = editor.getValue();
      this.dispatchEvent(new CustomEvent("change", { detail: { value } }));
    });
  }
}

customElements.define("x-code-editor", XCodeEditor);
