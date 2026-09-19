import { useEffect, useState, useRef } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Card from "@mui/material/Card";
import Divider from "@mui/material/Divider";
import Grid from "@mui/material/Grid";
import Icon from "@mui/material/Icon";
import LinearProgress from "@mui/material/LinearProgress";
import Skeleton from "@mui/material/Skeleton";
import SoftBadge from "components/SoftBadge";
import SoftBox from "components/SoftBox";
import SoftButton from "components/SoftButton";
import SoftTypography from "components/SoftTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import CertificateTemplate from "components/CertificateTemplate";
import { useAuth } from "context/auth/AuthContext";
import { apiRequest } from "services/api";
import html2canvas from "html2canvas";
import jsPDF from "jspdf";

function ParticipantResult() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [isGeneratingPdf, setIsGeneratingPdf] = useState(false);
  const certificateRef = useRef(null);

  useEffect(() => {
    apiRequest(`/participant/periods/${id}/result`)
      .then((response) => setResult(response.data))
      .catch((requestError) => setError(requestError.message))
      .finally(() => setLoading(false));
  }, [id]);

  const isPassed = result?.status === "passed";
  const maxPossibleScore = 677; // TOEFL max score

  const handleDownloadPDF = async () => {
    if (!certificateRef.current) return;
    setIsGeneratingPdf(true);
    try {
      const element = certificateRef.current;
      const canvas = await html2canvas(element, { scale: 2, useCORS: true });
      const imgData = canvas.toDataURL("image/png");
      const pdf = new jsPDF("l", "pt", "a4");
      
      const pdfWidth = pdf.internal.pageSize.getWidth();
      const pdfHeight = pdf.internal.pageSize.getHeight();
      
      pdf.addImage(imgData, "PNG", 0, 0, pdfWidth, pdfHeight);
      pdf.save(`Sertifikat_${user?.name || "Peserta"}_${result?.period_title || "SIUJI"}.pdf`);
    } catch (err) {
      console.error("Gagal membuat PDF:", err);
      setError("Gagal membuat PDF Sertifikat.");
    } finally {
      setIsGeneratingPdf(false);
    }
  };

  // ── Loading state ───────────────────────────────────────
  if (loading) {
    return (
      <DashboardLayout>
        <DashboardNavbar />
        <SoftBox pt={6} pb={3}>
          <Card>
            <SoftBox p={4}>
              <Skeleton variant="text" width={200} height={40} />
              <Skeleton variant="text" width={300} height={24} sx={{ mt: 1 }} />
              <Skeleton variant="circular" width={120} height={120} sx={{ mx: "auto", mt: 4 }} />
              <Skeleton variant="text" width="60%" sx={{ mx: "auto", mt: 3 }} />
              <Grid container spacing={3} mt={2}>
                {[1, 2, 3].map((item) => (
                  <Grid item xs={12} md={4} key={item}>
                    <Skeleton variant="rounded" height={120} />
                  </Grid>
                ))}
              </Grid>
            </SoftBox>
          </Card>
        </SoftBox>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox pt={6} pb={3}>
        {/* Error */}
        {error && (
          <SoftBox mb={3} p={2} borderRadius="md" bgColor="error">
            <SoftTypography variant="button" color="white">{error}</SoftTypography>
          </SoftBox>
        )}

        {/* Back button */}
        <SoftBox mb={3}>
          <SoftButton variant="outlined" color="secondary" size="small" onClick={() => navigate("/participant/periods")}>
            <Icon sx={{ mr: 0.5 }}>arrow_back</Icon>Kembali ke daftar ujian
          </SoftButton>
        </SoftBox>

        {/* Main result card */}
        <Card>
          <SoftBox p={4}>
            {/* Header */}
            <SoftBox textAlign="center" mb={4}>
              <SoftTypography variant="h4" fontWeight="bold">Hasil Ujian</SoftTypography>
              <SoftTypography variant="h6" color="text" mt={1}>
                {result?.period_title || "-"}
              </SoftTypography>
            </SoftBox>

            {/* Score circle */}
            <SoftBox textAlign="center" mb={4}>
              <SoftBox
                display="inline-flex"
                alignItems="center"
                justifyContent="center"
                width={160}
                height={160}
                borderRadius="50%"
                sx={{
                  background: isPassed
                    ? "linear-gradient(135deg, #43a047, #66bb6a)"
                    : "linear-gradient(135deg, #e53935, #ef5350)",
                  boxShadow: isPassed
                    ? "0 8px 32px rgba(67, 160, 71, 0.3)"
                    : "0 8px 32px rgba(229, 57, 53, 0.3)",
                }}
              >
                <SoftBox textAlign="center">
                  <SoftTypography variant="h1" color="white" fontWeight="bold" sx={{ lineHeight: 1.2 }}>
                    {result?.final_score ?? "-"}
                  </SoftTypography>
                  <SoftTypography variant="caption" color="white" sx={{ opacity: 0.9 }}>
                    / {maxPossibleScore}
                  </SoftTypography>
                </SoftBox>
              </SoftBox>

              <SoftBox mt={3}>
                <SoftBadge
                  color={isPassed ? "success" : "error"}
                  variant="gradient"
                  size="lg"
                >
                  {isPassed ? "LULUS" : "TIDAK LULUS"}
                </SoftBadge>
              </SoftBox>

              <SoftTypography variant="button" color="text" display="block" mt={2}>
                Nilai minimum kelulusan: <strong>{result?.min_passing_grade ?? "-"}</strong>
              </SoftTypography>
            </SoftBox>

            <Divider />

            {/* Section breakdowns */}
            <SoftTypography variant="h6" fontWeight="bold" mt={3} mb={2}>
              Rincian Nilai per Section
            </SoftTypography>

            <Grid container spacing={3}>
              {(result?.section_breakdowns || []).map((section) => {
                const scorePercent = maxPossibleScore > 0
                  ? Math.round((section.scaled_score / maxPossibleScore) * 100)
                  : 0;
                const sectionColors = {
                  listening: "info",
                  structure: "warning",
                  reading: "success",
                };
                const color = sectionColors[section.section_type] || "info";

                return (
                  <Grid item xs={12} md={4} key={section.section_title}>
                    <Card variant="outlined" sx={{ height: "100%" }}>
                      <SoftBox p={3}>
                        {/* Section header */}
                        <SoftBox display="flex" justifyContent="space-between" alignItems="center" mb={2}>
                          <SoftBox display="flex" alignItems="center" gap={1}>
                            <SoftBox
                              display="flex"
                              alignItems="center"
                              justifyContent="center"
                              width={36}
                              height={36}
                              borderRadius="md"
                              bgColor={color}
                            >
                              <Icon sx={{ color: "#fff", fontSize: 20 }}>
                                {section.section_type === "listening" ? "headphones" :
                                 section.section_type === "structure" ? "text_fields" : "menu_book"}
                              </Icon>
                            </SoftBox>
                            <SoftTypography variant="button" fontWeight="bold">
                              {section.section_title}
                            </SoftTypography>
                          </SoftBox>
                        </SoftBox>

                        {/* Score */}
                        <SoftTypography variant="h3" fontWeight="bold" textAlign="center" my={2}>
                          {section.scaled_score}
                        </SoftTypography>

                        {/* Progress bar */}
                        <LinearProgress
                          variant="determinate"
                          value={Math.min(scorePercent, 100)}
                          color={color}
                          sx={{ height: 8, borderRadius: 4 }}
                        />

                        {/* Details */}
                        <SoftBox display="flex" justifyContent="space-between" mt={2}>
                          <SoftBox>
                            <SoftTypography variant="caption" color="text">Jawaban benar</SoftTypography>
                            <SoftTypography variant="button" fontWeight="bold" display="block">
                              {section.correct_count}
                            </SoftTypography>
                          </SoftBox>
                          <SoftBox textAlign="right">
                            <SoftTypography variant="caption" color="text">Persentase</SoftTypography>
                            <SoftTypography variant="button" fontWeight="bold" display="block">
                              {scorePercent}%
                            </SoftTypography>
                          </SoftBox>
                        </SoftBox>
                      </SoftBox>
                    </Card>
                  </Grid>
                );
              })}
            </Grid>

            {/* Certificate Generation Button (always available if passed) */}
            {isPassed && (
              <>
                <Divider sx={{ mt: 4 }} />
                <SoftBox textAlign="center" mt={4} p={3} bgColor="grey.100" borderRadius="lg">
                  <Icon color="success" sx={{ fontSize: 48 }}>workspace_premium</Icon>
                  <SoftTypography variant="h5" fontWeight="bold" mt={2}>
                    Selamat! Anda Berhak Mendapat Sertifikat
                  </SoftTypography>
                  <SoftTypography variant="button" color="text" display="block" mt={1}>
                    Sertifikat kelulusan Anda telah tersedia dan dapat diunduh langsung.
                  </SoftTypography>
                  <SoftButton
                    variant="gradient"
                    color="success"
                    size="large"
                    onClick={handleDownloadPDF}
                    disabled={isGeneratingPdf}
                    sx={{ mt: 3, px: 6 }}
                  >
                    <Icon sx={{ mr: 1 }}>{isGeneratingPdf ? "hourglass_empty" : "download"}</Icon>
                    {isGeneratingPdf ? "Membuat PDF..." : "Unduh e-Sertifikat (PDF)"}
                  </SoftButton>
                </SoftBox>
                
                {/* Hidden template for PDF generation */}
                <CertificateTemplate 
                  ref={certificateRef}
                  participantName={user?.name || "Peserta"}
                  periodTitle={result?.period_title || "Periode Ujian"}
                  score={result?.final_score || 0}
                  dateStr={new Date().toLocaleDateString("id-ID", { day: "numeric", month: "long", year: "numeric" })}
                />
              </>
            )}

            {/* Not passed message */}
            {!isPassed && result && (
              <>
                <Divider sx={{ mt: 4 }} />
                <SoftBox textAlign="center" mt={4} p={3} bgColor="grey.100" borderRadius="lg">
                  <Icon color="warning" sx={{ fontSize: 48 }}>sentiment_dissatisfied</Icon>
                  <SoftTypography variant="h6" mt={2}>
                    Nilai Anda belum memenuhi batas kelulusan
                  </SoftTypography>
                  <SoftTypography variant="button" color="text" display="block" mt={1}>
                    Nilai minimum kelulusan adalah {result.min_passing_grade}. Anda mendapat {result.final_score}.
                    Hubungi admin untuk informasi lebih lanjut.
                  </SoftTypography>
                </SoftBox>
              </>
            )}
          </SoftBox>
        </Card>
      </SoftBox>
      <Footer />
    </DashboardLayout>
  );
}

export default ParticipantResult;
