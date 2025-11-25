import { Navigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";
import type { ReactNode } from "react";

interface ProtectedRouteProps {
  children: ReactNode;
  roles?: string[];
}

interface JwtPayload {
  role: string;
  user_id: string;
  exp: number;
}

export const ProtectedRoute = ({ children, roles }: ProtectedRouteProps) => {
  const token = localStorage.getItem("token");

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  try {
    const decoded = jwtDecode<JwtPayload>(token);

    // eslint-disable-next-line
    const currentTime = Date.now() / 1000;

    if (decoded.exp < currentTime) {
      return <Navigate to="/login" replace />;
    }

    // Проверка роли
    if (roles && !roles.includes(decoded.role)) {
      return <Navigate to="/dashboard" replace />;
    }

    return <>{children}</>;
  } catch {
    return <Navigate to="/login" replace />;
  }
};
