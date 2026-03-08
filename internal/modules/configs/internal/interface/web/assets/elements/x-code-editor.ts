import loader from "@monaco-editor/loader";
import Alpine from "alpinejs";
import type Monaco from "monaco-editor";

Alpine.data("monaco", ({ language, source = "" }: { language: string; source: string }) => ({
    loading: true,
    height: "300px",
    async init() {
        this.height = this.$root ? `${this.$root.clientHeight * 0.5}px` : "300px";

        const monaco = (await loader.init()) as typeof Monaco;
        const editor = monaco.editor.create(this.$root, {
            value: source,
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
            this.$dispatch("input", { value });
        });

        this.loading = false;
    },
}));
