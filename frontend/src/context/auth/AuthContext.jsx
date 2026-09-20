import { createContext, useContext, useEffect, useMemo, useState } from "react";
import PropTypes from "prop-types";
import { apiRequest } from "services/api";

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    apiRequest("/auth/me")
      .then((response) => { if (mounted) setUser(response.data || null); })
      .catch(() => { if (mounted) setUser(null); })
      .finally(() => { if (mounted) setLoading(false); });
    return () => { mounted = false; };
  }, []);

  const signIn = async (email, password) => {
    setLoading(true);
    try {
      await apiRequest("/auth/login", { method: "POST", body: JSON.stringify({ email, password }) });
      const response = await apiRequest("/auth/me");
      const authenticatedUser = response.data || null;
      setUser(authenticatedUser);
      return authenticatedUser;
    } finally {
      setLoading(false);
    }
  };

  const signOut = async () => {
    setLoading(true);
    try {
      await apiRequest("/auth/logout", { method: "POST" });
    } finally {
      setUser(null);
      setLoading(false);
    }
  };

  const forgotPassword = async (email) => {
    try {
      await apiRequest("/auth/forgot-password", { method: "POST", body: JSON.stringify({ email }) });
      return true;
    } catch (err) {
      throw err;
    }
  };

  const verifyOTP = async (email, otp) => {
    try {
      const res = await apiRequest("/auth/verify-otp", { method: "POST", body: JSON.stringify({ email, otp }) });
      return res.data; // Should return { reset_token: "..." }
    } catch (err) {
      throw err;
    }
  };

  const updatePassword = async (token, new_password) => {
    try {
      await apiRequest("/auth/reset-password", { method: "POST", body: JSON.stringify({ token, new_password }) });
      return true;
    } catch (err) {
      throw err;
    }
  };

  const value = useMemo(() => ({ user, loading, signIn, signOut, forgotPassword, verifyOTP, updatePassword }), [user, loading]);
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within an AuthProvider");
  return context;
};

AuthProvider.propTypes = { children: PropTypes.node.isRequired };
