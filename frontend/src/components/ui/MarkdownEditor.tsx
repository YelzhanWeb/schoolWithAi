import React from "react";
import {
  MDXEditor,
  headingsPlugin,
  listsPlugin,
  quotePlugin,
  thematicBreakPlugin,
  markdownShortcutPlugin,
  linkPlugin,
  linkDialogPlugin,
  imagePlugin,
  tablePlugin,
  codeBlockPlugin,
  codeMirrorPlugin,
  diffSourcePlugin,
  toolbarPlugin,
  UndoRedo,
  BoldItalicUnderlineToggles,
  BlockTypeSelect,
  CreateLink,
  InsertImage,
  InsertTable,
  ListsToggle,
  Separator,
  InsertCodeBlock,
} from "@mdxeditor/editor";
import "@mdxeditor/editor/style.css";

interface MarkdownEditorProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
}

export const MarkdownEditor: React.FC<MarkdownEditorProps> = ({
  value,
  onChange,
  placeholder = "Начните писать контент урока...",
}) => {
  return (
    <div className="border rounded-lg overflow-hidden bg-white">
      <MDXEditor
        markdown={value}
        onChange={onChange}
        placeholder={placeholder}
        plugins={[
          // Основные возможности форматирования
          headingsPlugin(),
          listsPlugin(),
          quotePlugin(),
          thematicBreakPlugin(),
          linkPlugin(),
          linkDialogPlugin(),
          imagePlugin(),
          tablePlugin(),
          codeBlockPlugin({ defaultCodeBlockLanguage: "javascript" }),
          codeMirrorPlugin({
            codeBlockLanguages: {
              javascript: "JavaScript",
              typescript: "TypeScript",
              python: "Python",
              html: "HTML",
              css: "CSS",
              json: "JSON",
            },
          }),

          // Горячие клавиши Markdown
          markdownShortcutPlugin(),

          // Возможность переключиться на исходный код
          diffSourcePlugin({ viewMode: "rich-text" }),

          // Панель инструментов
          toolbarPlugin({
            toolbarContents: () => (
              <>
                <UndoRedo />
                <Separator />
                <BoldItalicUnderlineToggles />
                <Separator />
                <BlockTypeSelect />
                <Separator />
                <CreateLink />
                <InsertImage />
                <Separator />
                <ListsToggle />
                <Separator />
                <InsertTable />
                <InsertCodeBlock />
              </>
            ),
          }),
        ]}
        contentEditableClassName="prose max-w-none p-4 min-h-[400px] focus:outline-none"
      />
    </div>
  );
};
