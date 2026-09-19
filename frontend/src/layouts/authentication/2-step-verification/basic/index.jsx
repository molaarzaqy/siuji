/**
=========================================================
* Soft UI Dashboard PRO React - v4.0.3
=========================================================

* Product Page: https://www.creative-tim.com/product/soft-ui-dashboard-pro-react
* Copyright 2024 Creative Tim (https://www.creative-tim.com)

Coded by www.creative-tim.com

 =========================================================

* The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
*/

import { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";

// @mui material components
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import Alert from "@mui/material/Alert";

// Soft UI Dashboard PRO React components
import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";
import SoftInput from "components/SoftInput";
import SoftButton from "components/SoftButton";

// Soft UI Dashboard PRO React example components
import DefaultNavbar from "examples/Navbars/DefaultNavbar";
import PageLayout from "examples/LayoutContainers/PageLayout";

// Soft UI Dashboard PRO React page layout routes
import pageRoutes from "page.routes";

import { useAuth } from "context/auth/AuthContext";
import Footer from "layouts/authentication/components/Footer";

function Basic() {
  const [otp, setOtp] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();
  const location = useLocation();
  const { verifyOTP } = useAuth();
  
  const email = location.state?.email || "";

  if (!email) {
    navigate("/authentication/reset-password/basic");
    return null;
  }

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (otp.length !== 6) {
      setError("Please enter a 6-digit OTP code.");
      return;
    }
    setError("");
    setLoading(true);

    try {
      const data = await verifyOTP(email, otp);
      if (data && data.reset_token) {
        navigate("/authentication/reset-password/update", { state: { token: data.reset_token } });
      } else {
        setError("Invalid response from server.");
      }
    } catch (err) {
      setError(err.message || "Invalid OTP code.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <PageLayout background="light">
      <DefaultNavbar
        routes={pageRoutes}
        action={{
          type: "external",
          route: "https://www.creative-tim.com/product/soft-ui-dashboard-pro-react",
          label: "buy now",
        }}
      />
      <Grid container spacing={3} justifyContent="center" sx={{ minHeight: "75vh" }}>
        <Grid item xs={10} md={6} lg={4}>
          <SoftBox mt={32} mb={3} px={{ xs: 0, lg: 2 }}>
            <Card>
              <SoftBox pt={3} px={3} pb={1} textAlign="center">
                <SoftTypography variant="h4" fontWeight="bold" textGradient>
                  2-Step Verification
                </SoftTypography>
                <SoftTypography variant="body2" color="text">
                  Enter the 6-digit code sent to {email}
                </SoftTypography>
              </SoftBox>
              <SoftBox p={3}>
                <SoftBox component="form" role="form" onSubmit={handleSubmit}>
                  {error && (
                    <SoftBox mb={2}>
                      <Alert severity="error">{error}</Alert>
                    </SoftBox>
                  )}
                  <SoftBox mb={2} display="flex" justifyContent="center">
                    <SoftInput
                      size="large"
                      placeholder="------"
                      inputProps={{ maxLength: 6, style: { textAlign: "center", letterSpacing: "8px", fontSize: "1.5rem" } }}
                      value={otp}
                      onChange={(e) => setOtp(e.target.value.replace(/\D/g, ""))}
                      required
                    />
                  </SoftBox>
                  <SoftBox mt={5} mb={1}>
                    <SoftButton variant="gradient" color="dark" size="large" fullWidth type="submit" disabled={loading}>
                      {loading ? "Verifying..." : "Verify Code"}
                    </SoftButton>
                  </SoftBox>
                </SoftBox>
              </SoftBox>
            </Card>
          </SoftBox>
        </Grid>
      </Grid>
      <Footer />
    </PageLayout>
  );
}

export default Basic;
