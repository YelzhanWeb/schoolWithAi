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
import { uploadApi } from "../../api/upload"; // Импортируем API

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
  const handleImageUpload = async (file: File): Promise<string> => {
    try {
      const url = await uploadApi.uploadFile(file, "lesson");
      return url;
    } catch (error) {
      console.error("Image upload failed", error);
      alert("Ошибка загрузки изображения");
      throw error;
    }
  };

  return (
    <div className="border rounded-lg overflow-hidden bg-white">
      <MDXEditor
        markdown={value}
        onChange={onChange}
        placeholder={placeholder}
        plugins={[
          headingsPlugin(),
          listsPlugin(),
          quotePlugin(),
          thematicBreakPlugin(),
          linkPlugin(),
          linkDialogPlugin(),
          // Подключаем загрузчик
          imagePlugin({
            imageUploadHandler: handleImageUpload,
          }),
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
          markdownShortcutPlugin(),
          diffSourcePlugin({ viewMode: "rich-text" }),
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
                <InsertImage /> {/* Кнопка вставки картинки */}
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
