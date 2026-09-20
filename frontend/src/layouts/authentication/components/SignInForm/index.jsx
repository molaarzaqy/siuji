import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import Alert from "@mui/material/Alert";
import IconButton from "@mui/material/IconButton";
import Switch from "@mui/material/Switch";
import Visibility from "@mui/icons-material/Visibility";
import VisibilityOff from "@mui/icons-material/VisibilityOff";
import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";
import SoftInput from "components/SoftInput";
import SoftButton from "components/SoftButton";
import { useAuth } from "context/auth/AuthContext";

function SignInForm() {
  const navigate = useNavigate();
  const location = useLocation();
  const { signIn, loading } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [rememberMe, setRememberMe] = useState(true);
  const [error, setError] = useState("");

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError("");
    if (!email || !password) return setError("Email dan password wajib diisi.");
    try {
      const authenticatedUser = await signIn(email, password);
      const destination = location.state?.from?.pathname || (authenticatedUser?.role === "participant" ? "/participant/dashboard" : "/dashboards/default");
      navigate(destination, { replace: true });
    } catch (signInError) {
      setError(signInError.message);
    }
  };

  return (
    <SoftBox component="form" role="form" onSubmit={handleSubmit}>
      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
      <SoftBox mb={2}>
        <SoftTypography component="label" variant="caption" fontWeight="bold" display="block" mb={1}>Email</SoftTypography>
        <SoftInput type="email" placeholder="Email" value={email} onChange={(event) => setEmail(event.target.value)} autoComplete="email" required />
      </SoftBox>
      <SoftBox mb={2}>
        <SoftTypography component="label" variant="caption" fontWeight="bold" display="block" mb={1}>Password</SoftTypography>
        <SoftInput type="password" placeholder="Password" value={password} onChange={(event) => setPassword(event.target.value)} autoComplete="current-password" required />
      </SoftBox>
      <SoftBox display="flex" alignItems="center">
        <Switch checked={rememberMe} onChange={() => setRememberMe((value) => !value)} />
        <SoftTypography variant="button" fontWeight="regular" onClick={() => setRememberMe((value) => !value)} sx={{ cursor: "pointer", userSelect: "none" }}>&nbsp;&nbsp;Remember me</SoftTypography>
      </SoftBox>
      <SoftBox mt={4} mb={1}>
        <SoftButton type="submit" variant="gradient" color="info" fullWidth disabled={loading}>{loading ? "Signing in..." : "Sign in"}</SoftButton>
      </SoftBox>
      <SoftBox mt={3} textAlign="center">
        <SoftTypography variant="button" color="text" fontWeight="regular">Lupa password? <SoftTypography component={Link} to="/authentication/reset-password/basic" variant="button" color="info" fontWeight="medium" textGradient>Reset password</SoftTypography></SoftTypography>
      </SoftBox>
    </SoftBox>
  );
}

export default SignInForm;
