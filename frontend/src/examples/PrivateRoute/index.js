import { Navigate, useLocation } from "react-router-dom";
import PropTypes from "prop-types";
import { useAuth } from "context/auth/AuthContext";

function PrivateRoute({ children }) {
  const { user, loading } = useAuth();
  const location = useLocation();
  if (loading) return null;
  if (!user) return <Navigate to="/authentication/sign-in/basic" replace state={{ from: location }} />;
  return children;
}

export function PublicRoute({ children }) {
  const { user, loading } = useAuth();

  if (loading) return null;
  if (user) return <Navigate to={user.role === "participant" ? "/participant/periods" : "/dashboards/default"} replace />;
  return children;
}

export function RoleRoute({ children, roles }) {
  const { user, loading } = useAuth();
  if (loading) return null;
  if (!user) return <Navigate to="/authentication/sign-in/basic" replace />;
  if (roles?.length && !roles.includes(user.role)) return <Navigate to="/authentication/error/404" replace />;
  return children;
}

PrivateRoute.propTypes = { children: PropTypes.node.isRequired };
PublicRoute.propTypes = { children: PropTypes.node.isRequired };
RoleRoute.propTypes = { children: PropTypes.node.isRequired, roles: PropTypes.arrayOf(PropTypes.string) };
export default PrivateRoute;
