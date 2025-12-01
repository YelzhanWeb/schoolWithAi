import { Outlet } from "react-router-dom";
import { StudentHeader } from "./StudentHeader";

export const StudentLayout = () => {
  return (
    <div className="min-h-screen bg-[#F5F6FA]">
      {/* Единая шапка */}
      <StudentHeader />

      {/* Контейнер для страниц (Outlet рендерит вложенные роуты) */}
      <div className="pt-16">
        <Outlet />
      </div>
    </div>
  );
};
