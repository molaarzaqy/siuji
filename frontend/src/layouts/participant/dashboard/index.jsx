import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import Icon from "@mui/material/Icon";
import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";
import SoftButton from "components/SoftButton";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import MiniStatisticsCard from "examples/Cards/StatisticsCards/MiniStatisticsCard";
import { useAuth } from "context/auth/AuthContext";
import { apiRequest } from "services/api";
import typography from "assets/theme/base/typography";
import DefaultInfoCard from "examples/Cards/InfoCards/DefaultInfoCard";

function ParticipantDashboard() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [periods, setPeriods] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    apiRequest("/participant/periods")
      .then((response) => {
        if (mounted) setPeriods(response.data || []);
      })
      .catch(() => {
        // Abaikan error untuk dashboard
      })
      .finally(() => {
        if (mounted) setLoading(false);
      });
    return () => { mounted = false; };
  }, []);

  const totalExams = periods.length;
  const completedExams = periods.filter(p => p.status === "completed").length;
  const pendingExams = totalExams - completedExams;

  const firstName = user?.name ? user.name.split(" ")[0] : "Peserta";

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox py={3}>
        <SoftBox mb={3}>
          <Grid container spacing={3}>
            {/* Welcome Banner */}
            <Grid item xs={12}>
              <Card sx={{ 
                  position: "relative", 
                  overflow: "hidden", 
                  backgroundImage: "linear-gradient(310deg, #141727, #3A416F)", 
                  color: "#fff" 
                }}>
                <SoftBox p={4} position="relative" zIndex={2}>
                  <SoftTypography variant="h3" color="white" fontWeight="bold" mb={1}>
                    Halo, {firstName}! 👋
                  </SoftTypography>
                  <SoftTypography variant="body2" color="white" opacity={0.8} mb={3}>
                    Selamat datang di Dashboard Ujian Anda. Siapkan diri Anda dan raih hasil terbaik!
                  </SoftTypography>
                  <SoftButton variant="gradient" color="info" onClick={() => navigate("/participant/periods")}>
                    Mulai Ujian Sekarang
                  </SoftButton>
                </SoftBox>
                {/* Decorative Elements */}
                <Icon sx={{ 
                    position: "absolute", 
                    right: "-5%", 
                    top: "-10%", 
                    fontSize: "20rem !important", 
                    color: "rgba(255, 255, 255, 0.05)",
                    transform: "rotate(-15deg)",
                    zIndex: 1
                  }}>
                  workspace_premium
                </Icon>
              </Card>
            </Grid>
          </Grid>
        </SoftBox>

        {/* Statistics Cards */}
        <SoftBox mb={3}>
          <Grid container spacing={3}>
            <Grid item xs={12} sm={4}>
              <MiniStatisticsCard
                title={{ text: "Total Ujian" }}
                count={loading ? "-" : totalExams}
                icon={{ color: "info", component: "assignment" }}
              />
            </Grid>
            <Grid item xs={12} sm={4}>
              <MiniStatisticsCard
                title={{ text: "Ujian Selesai" }}
                count={loading ? "-" : completedExams}
                icon={{ color: "success", component: "check_circle" }}
              />
            </Grid>
            <Grid item xs={12} sm={4}>
              <MiniStatisticsCard
                title={{ text: "Ujian Menunggu" }}
                count={loading ? "-" : pendingExams}
                icon={{ color: "warning", component: "hourglass_empty" }}
              />
            </Grid>
          </Grid>
        </SoftBox>

        {/* Quick Actions */}
        <Grid container spacing={3}>
          <Grid item xs={12} md={6}>
            <SoftBox onClick={() => navigate("/participant/periods")} sx={{ cursor: "pointer", transition: "transform 0.2s", "&:hover": { transform: "translateY(-5px)" } }}>
              <DefaultInfoCard
                icon="event_note"
                title="Daftar Ujian Tersedia"
                description="Lihat semua periode ujian yang ditugaskan kepada Anda, periksa jadwal, dan kerjakan ujian yang sedang aktif."
                value="Lihat Detail"
              />
            </SoftBox>
          </Grid>
          <Grid item xs={12} md={6}>
            <SoftBox onClick={() => navigate("/participant/periods")} sx={{ cursor: "pointer", transition: "transform 0.2s", "&:hover": { transform: "translateY(-5px)" } }}>
              <DefaultInfoCard
                icon="workspace_premium"
                title="Riwayat & Sertifikat"
                description="Lihat skor dari ujian yang telah diselesaikan dan unduh sertifikat kelulusan Anda."
                value="Lihat Hasil"
              />
            </SoftBox>
          </Grid>
        </Grid>
      </SoftBox>
      <Footer />
    </DashboardLayout>
  );
}

export default ParticipantDashboard;
