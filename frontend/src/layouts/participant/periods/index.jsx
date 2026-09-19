import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Card from "@mui/material/Card";
import Grid from "@mui/material/Grid";
import Icon from "@mui/material/Icon";
import Skeleton from "@mui/material/Skeleton";
import SoftBadge from "components/SoftBadge";
import SoftBox from "components/SoftBox";
import SoftButton from "components/SoftButton";
import SoftTypography from "components/SoftTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import { apiRequest } from "services/api";

const statusMeta = {
  registered: { label: "Belum Mulai", color: "info", icon: "how_to_reg" },
  started: { label: "Sedang Berjalan", color: "warning", icon: "play_circle" },
  completed: { label: "Selesai", color: "success", icon: "check_circle" },
};

function formatDate(value) {
  if (!value) return "-";
  return new Date(value).toLocaleString("id-ID", { dateStyle: "medium", timeStyle: "short" });
}

function ParticipantPeriods() {
  const [periods, setPeriods] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    apiRequest("/participant/periods")
      .then((response) => setPeriods(response.data || []))
      .catch((requestError) => setError(requestError.message))
      .finally(() => setLoading(false));
  }, []);

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox pt={6} pb={3}>
        {/* Header */}
        <SoftBox mb={3}>
          <SoftTypography variant="h4" fontWeight="bold">Ujian Saya</SoftTypography>
          <SoftTypography variant="button" color="text" fontWeight="regular">
            Daftar periode ujian yang ditugaskan kepada Anda.
          </SoftTypography>
        </SoftBox>

        {/* Error */}
        {error && (
          <SoftBox mb={3} p={2} borderRadius="md" bgColor="error">
            <SoftTypography variant="button" color="white">{error}</SoftTypography>
          </SoftBox>
        )}

        {/* Loading skeleton */}
        {loading ? (
          <Grid container spacing={3}>
            {[1, 2, 3].map((item) => (
              <Grid item xs={12} md={6} lg={4} key={item}>
                <Card>
                  <SoftBox p={3}>
                    <Skeleton variant="circular" width={40} height={40} />
                    <Skeleton variant="text" width="70%" height={32} sx={{ mt: 2 }} />
                    <Skeleton variant="text" width="50%" height={20} sx={{ mt: 1 }} />
                    <Skeleton variant="text" width="60%" height={20} sx={{ mt: 0.5 }} />
                    <Skeleton variant="rounded" height={40} sx={{ mt: 3 }} />
                  </SoftBox>
                </Card>
              </Grid>
            ))}
          </Grid>
        ) : (
          <Grid container spacing={3}>
            {periods.map((assignment) => {
              const period = assignment.period || assignment;
              const status = statusMeta[assignment.status] || statusMeta.registered;
              const now = Date.now();
              const started = new Date(period.start_time).getTime() <= now;
              const ended = new Date(period.end_time).getTime() < now;
              const isCompleted = assignment.status === "completed";

              return (
                <Grid item xs={12} md={6} lg={4} key={assignment.public_id}>
                  <Card sx={{ height: "100%", display: "flex", flexDirection: "column", transition: "transform 0.2s, box-shadow 0.2s", "&:hover": { transform: "translateY(-4px)", boxShadow: 6 } }}>
                    <SoftBox p={3} display="flex" flexDirection="column" flex={1}>
                      {/* Top badge */}
                      <SoftBox display="flex" justifyContent="space-between" alignItems="flex-start">
                        <SoftBox
                          display="flex"
                          alignItems="center"
                          justifyContent="center"
                          width={44}
                          height={44}
                          borderRadius="md"
                          bgColor="info"
                        >
                          <Icon sx={{ color: "#fff", fontSize: 24 }}>{status.icon}</Icon>
                        </SoftBox>
                        <SoftBadge color={status.color} variant="gradient" size="sm">
                          {status.label}
                        </SoftBadge>
                      </SoftBox>

                      {/* Content */}
                      <SoftTypography variant="h5" fontWeight="bold" mt={2}>
                        {period.title}
                      </SoftTypography>

                      <SoftBox mt={2} flex={1}>
                        <SoftBox display="flex" alignItems="center" gap={1} mb={1}>
                          <Icon color="action" sx={{ fontSize: 18 }}>schedule</Icon>
                          <SoftTypography variant="caption" color="text">
                            Mulai: {formatDate(period.start_time)}
                          </SoftTypography>
                        </SoftBox>
                        <SoftBox display="flex" alignItems="center" gap={1} mb={1}>
                          <Icon color="action" sx={{ fontSize: 18 }}>event</Icon>
                          <SoftTypography variant="caption" color="text">
                            Selesai: {formatDate(period.end_time)}
                          </SoftTypography>
                        </SoftBox>
                        {!started && (
                          <SoftBox display="flex" alignItems="center" gap={1}>
                            <Icon color="warning" sx={{ fontSize: 18 }}>info</Icon>
                            <SoftTypography variant="caption" color="warning">
                              Ujian belum dimulai
                            </SoftTypography>
                          </SoftBox>
                        )}
                        {ended && !isCompleted && (
                          <SoftBox display="flex" alignItems="center" gap={1}>
                            <Icon color="error" sx={{ fontSize: 18 }}>warning</Icon>
                            <SoftTypography variant="caption" color="error">
                              Waktu ujian sudah berakhir
                            </SoftTypography>
                          </SoftBox>
                        )}
                      </SoftBox>

                      {/* Action button */}
                      <SoftButton
                        fullWidth
                        variant="gradient"
                        color={isCompleted ? "success" : "info"}
                        sx={{ mt: 3 }}
                        onClick={() => navigate(isCompleted ? `/participant/results/${period.public_id}` : `/participant/exam/${period.public_id}`)}
                      >
                        <Icon sx={{ mr: 1 }}>{isCompleted ? "assessment" : "play_arrow"}</Icon>
                        {isCompleted ? "Lihat Hasil" : "Detail / Mulai"}
                      </SoftButton>
                    </SoftBox>
                  </Card>
                </Grid>
              );
            })}

            {/* Empty state */}
            {!periods.length && (
              <Grid item xs={12}>
                <Card>
                  <SoftBox p={6} textAlign="center">
                    <Icon color="disabled" sx={{ fontSize: 64 }}>assignment_late</Icon>
                    <SoftTypography variant="h5" mt={2}>Belum Ada Ujian</SoftTypography>
                    <SoftTypography variant="button" color="text" mt={1}>
                      Belum ada periode ujian yang ditugaskan kepada Anda. Hubungi admin jika Anda merasa ini tidak benar.
                    </SoftTypography>
                  </SoftBox>
                </Card>
              </Grid>
            )}
          </Grid>
        )}
      </SoftBox>
      <Footer />
    </DashboardLayout>
  );
}

export default ParticipantPeriods;
