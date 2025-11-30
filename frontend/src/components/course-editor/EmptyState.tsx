import React from "react";

export const EmptyState: React.FC = () => {
  return (
    <div className="h-full flex flex-col items-center justify-center text-gray-400">
      <div className="text-6xl mb-4">👈</div>
      <p className="text-xl">Выберите урок или тест в меню слева</p>
    </div>
  );
};
