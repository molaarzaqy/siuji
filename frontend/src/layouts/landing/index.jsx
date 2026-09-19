import React from "react";
import { Link } from "react-router-dom";
import Container from "@mui/material/Container";
import Grid from "@mui/material/Grid";
import Icon from "@mui/material/Icon";
import Card from "@mui/material/Card";

import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";
import SoftButton from "components/SoftButton";
import PageLayout from "examples/LayoutContainers/PageLayout";

function LandingPage() {
  return (
    <PageLayout background="white">
      {/* ── Navbar ── */}
      <SoftBox
        position="absolute"
        top={0}
        left={0}
        width="100%"
        zIndex={3}
        px={3}
        py={2}
      >
        <Container>
          <SoftBox display="flex" justifyContent="space-between" alignItems="center">
            <SoftTypography variant="h5" fontWeight="bold" color="dark">
              SIUJI
            </SoftTypography>
            <SoftBox display="flex" gap={2}>
              <SoftButton
                component={Link}
                to="/authentication/sign-in/basic"
                variant="gradient"
                color="info"
                size="small"
              >
                Masuk / Login
              </SoftButton>
            </SoftBox>
          </SoftBox>
        </Container>
      </SoftBox>

      {/* ── Hero Section ── */}
      <SoftBox
        minHeight="100vh"
        width="100%"
        display="flex"
        alignItems="center"
        justifyContent="center"
        sx={{
          background: "linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%)",
          position: "relative",
          overflow: "hidden"
        }}
      >
        {/* Decorative Shapes */}
        <SoftBox
          position="absolute"
          top="-10%"
          right="-5%"
          width="50%"
          height="70%"
          borderRadius="50%"
          sx={{ background: "rgba(33, 150, 243, 0.05)", filter: "blur(60px)" }}
        />
        <SoftBox
          position="absolute"
          bottom="-10%"
          left="-5%"
          width="40%"
          height="60%"
          borderRadius="50%"
          sx={{ background: "rgba(233, 30, 99, 0.05)", filter: "blur(60px)" }}
        />

        <Container position="relative" zIndex={2}>
          <Grid container spacing={3} alignItems="center">
            <Grid item xs={12} lg={6}>
              <SoftBox mb={3}>
                <SoftTypography
                  variant="h1"
                  fontWeight="bold"
                  color="dark"
                  textGradient={false}
                  sx={{
                    fontSize: { xs: "2.5rem", md: "3.5rem" },
                    lineHeight: 1.2,
                    mb: 2,
                  }}
                >
                  Platform Ujian <br />
                  <SoftTypography
                    component="span"
                    variant="h1"
                    color="info"
                    textGradient
                    sx={{ fontSize: "inherit" }}
                  >
                    Digital & Terpercaya
                  </SoftTypography>
                </SoftTypography>
                <SoftTypography variant="body1" color="text" mb={4} pr={{ lg: 5 }}>
                  Tingkatkan integritas evaluasi Anda dengan antarmuka ujian yang intuitif, cepat, responsif, dan didesain khusus untuk memberikan pengalaman Computer Based Test (CBT) terbaik.
                </SoftTypography>
                <SoftBox display="flex" gap={2}>
                  <SoftButton
                    component={Link}
                    to="/authentication/sign-in/basic"
                    variant="gradient"
                    color="info"
                    size="large"
                  >
                    Mulai Ujian Sekarang
                  </SoftButton>
                </SoftBox>
              </SoftBox>
            </Grid>
            <Grid item xs={12} lg={6} sx={{ display: { xs: "none", lg: "block" } }}>
              <SoftBox
                position="relative"
                width="100%"
                height="100%"
                display="flex"
                justifyContent="center"
              >
                {/* Dashboard Mockup Illustration */}
                <Card
                  sx={{
                    width: "90%",
                    height: "400px",
                    background: "white",
                    borderRadius: "xl",
                    boxShadow: "xl",
                    overflow: "hidden",
                    position: "relative",
                  }}
                >
                  <SoftBox p={2} borderBottom="1px solid #eee" display="flex" gap={1}>
                    <SoftBox width="10px" height="10px" borderRadius="50%" bgColor="error" />
                    <SoftBox width="10px" height="10px" borderRadius="50%" bgColor="warning" />
                    <SoftBox width="10px" height="10px" borderRadius="50%" bgColor="success" />
                  </SoftBox>
                  <SoftBox p={3}>
                    <Grid container spacing={2}>
                      <Grid item xs={4}>
                        <SoftBox height="100px" borderRadius="lg" bgColor="grey-100" />
                      </Grid>
                      <Grid item xs={8}>
                        <SoftBox height="20px" width="80%" borderRadius="md" bgColor="grey-200" mb={2} />
                        <SoftBox height="15px" width="60%" borderRadius="md" bgColor="grey-100" mb={1} />
                        <SoftBox height="15px" width="90%" borderRadius="md" bgColor="grey-100" />
                      </Grid>
                    </Grid>
                    <SoftBox mt={3} height="150px" borderRadius="lg" bgColor="info" opacity={0.1} />
                  </SoftBox>
                </Card>
              </SoftBox>
            </Grid>
          </Grid>
        </Container>
      </SoftBox>

      {/* ── Features Section ── */}
      <SoftBox py={8} bgColor="white">
        <Container>
          <SoftBox textAlign="center" mb={6}>
            <SoftTypography variant="h2" fontWeight="bold" color="dark" mb={2}>
              Mengapa Memilih SIUJI?
            </SoftTypography>
            <SoftTypography variant="body2" color="text">
              Platform modern yang dibangun untuk menunjang segala kebutuhan tes dan evaluasi berskala besar.
            </SoftTypography>
          </SoftBox>

          <Grid container spacing={4}>
            <Grid item xs={12} md={4}>
              <Card sx={{ p: 4, height: "100%", textAlign: "center", boxShadow: "md" }}>
                <SoftBox
                  width="64px"
                  height="64px"
                  mx="auto"
                  mb={3}
                  borderRadius="lg"
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  bgColor="info"
                  variant="gradient"
                >
                  <Icon sx={{ color: "white", fontSize: "2rem !important" }}>security</Icon>
                </SoftBox>
                <SoftTypography variant="h5" fontWeight="bold" mb={2}>
                  Keamanan Tinggi
                </SoftTypography>
                <SoftTypography variant="body2" color="text">
                  Dilengkapi fitur pencegahan kecurangan, token sesi ujian, dan keamanan data peserta terenkripsi penuh.
                </SoftTypography>
              </Card>
            </Grid>
            <Grid item xs={12} md={4}>
              <Card sx={{ p: 4, height: "100%", textAlign: "center", boxShadow: "md" }}>
                <SoftBox
                  width="64px"
                  height="64px"
                  mx="auto"
                  mb={3}
                  borderRadius="lg"
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  bgColor="success"
                  variant="gradient"
                >
                  <Icon sx={{ color: "white", fontSize: "2rem !important" }}>speed</Icon>
                </SoftBox>
                <SoftTypography variant="h5" fontWeight="bold" mb={2}>
                  Hasil Real-time
                </SoftTypography>
                <SoftTypography variant="body2" color="text">
                  Penilaian dikalkulasi secara instan setelah peserta menyelesaikan ujian, lengkap dengan fitur konversi nilai.
                </SoftTypography>
              </Card>
            </Grid>
            <Grid item xs={12} md={4}>
              <Card sx={{ p: 4, height: "100%", textAlign: "center", boxShadow: "md" }}>
                <SoftBox
                  width="64px"
                  height="64px"
                  mx="auto"
                  mb={3}
                  borderRadius="lg"
                  display="flex"
                  alignItems="center"
                  justifyContent="center"
                  bgColor="warning"
                  variant="gradient"
                >
                  <Icon sx={{ color: "white", fontSize: "2rem !important" }}>devices</Icon>
                </SoftBox>
                <SoftTypography variant="h5" fontWeight="bold" mb={2}>
                  Responsif & Ringan
                </SoftTypography>
                <SoftTypography variant="body2" color="text">
                  Antarmuka yang dioptimalkan untuk berbagai perangkat. Ujian dapat dikerjakan lancar di desktop, tablet, maupun ponsel.
                </SoftTypography>
              </Card>
            </Grid>
          </Grid>
        </Container>
      </SoftBox>

      {/* ── Footer ── */}
      <SoftBox py={4} bgColor="light">
        <Container>
          <SoftBox display="flex" justifyContent="space-between" alignItems="center" flexWrap="wrap">
            <SoftTypography variant="body2" color="secondary">
              &copy; {new Date().getFullYear()} SIUJI Platform. All rights reserved.
            </SoftTypography>
            <SoftBox display="flex" gap={2}>
              <SoftTypography component="a" href="#" variant="body2" color="secondary">
                Bantuan
              </SoftTypography>
              <SoftTypography component="a" href="#" variant="body2" color="secondary">
                Privasi
              </SoftTypography>
            </SoftBox>
          </SoftBox>
        </Container>
      </SoftBox>
    </PageLayout>
  );
}

export default LandingPage;
