import { useCallback, useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Card from "@mui/material/Card";
import Chip from "@mui/material/Chip";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import FormControlLabel from "@mui/material/FormControlLabel";
import Grid from "@mui/material/Grid";
import Icon from "@mui/material/Icon";
import Radio from "@mui/material/Radio";
import RadioGroup from "@mui/material/RadioGroup";
import Skeleton from "@mui/material/Skeleton";
import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import Tooltip from "@mui/material/Tooltip";
import IconButton from "@mui/material/IconButton";
import SoftBox from "components/SoftBox";
import SoftButton from "components/SoftButton";
import SoftTypography from "components/SoftTypography";
import SoftSnackbar from "components/SoftSnackbar";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import CustomAudioPlayer from "components/CustomAudioPlayer";
import { apiRequest } from "services/api";

function ParticipantExam() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [session, setSession] = useState(null);
  const [answers, setAnswers] = useState({});
  const [flagged, setFlagged] = useState({});
  const [deadline, setDeadline] = useState(null);
  const [timeLeft, setTimeLeft] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [starting, setStarting] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [saveStatus, setSaveStatus] = useState("idle");
  const [activeSection, setActiveSection] = useState(0);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [snackbarState, setSnackbarState] = useState({ open: false, color: "info", icon: "notifications", title: "", content: "" });

  const requestSequence = useRef({});
  const autoSubmitted = useRef(false);
  const timerNotified = useRef(false);

  // ── Helper to show snackbar ─────────────────────────────
  const showSnackbar = useCallback((color, icon, title, content) => {
    setSnackbarState({ open: true, color, icon, title, content });
  }, []);
  const closeSnackbar = () => setSnackbarState(prev => ({ ...prev, open: false }));

  // ── Load period detail ──────────────────────────────────
  useEffect(() => {
    apiRequest(`/participant/periods/${id}`)
      .then((response) => {
        setSession(response.data);
        if (response.data?.end_time) setDeadline(new Date(response.data.end_time));
        if (response.data?.participant_status === "completed") {
          navigate(`/participant/results/${id}`, { replace: true });
        }
      })
      .catch((requestError) => setError(requestError.message))
      .finally(() => setLoading(false));
  }, [id, navigate]);

  // ── Submit handler ──────────────────────────────────────
  const submit = useCallback(
    async (automatic = false) => {
      if (submitting || autoSubmitted.current) return;
      if (automatic) autoSubmitted.current = true;
      setSubmitting(true);
      try {
        await apiRequest(`/participant/periods/${id}/submit`, { method: "POST" });
        navigate(`/participant/results/${id}`);
      } catch (requestError) {
        setError(requestError.message);
        autoSubmitted.current = false;
      } finally {
        setSubmitting(false);
      }
    },
    [id, navigate, submitting]
  );

  // ── Timer ───────────────────────────────────────────────
  useEffect(() => {
    if (!deadline || !session?.sections) return undefined;
    const tick = () => {
      const remaining = Math.max(0, deadline.getTime() - Date.now());
      setTimeLeft(remaining);
      
      // Auto submit at 0
      if (remaining === 0) {
        submit(true);
      }
      
      // Toast notification exactly at 5 minutes
      if (remaining > 0 && remaining <= 300000 && !timerNotified.current) {
        timerNotified.current = true;
        showSnackbar("warning", "timer", "Waktu Kritis!", "Waktu tersisa kurang dari 5 menit. Harap segera selesaikan jawaban Anda!");
      }
    };
    tick();
    const timer = window.setInterval(tick, 1000);
    return () => window.clearInterval(timer);
  }, [deadline, session, submit, showSnackbar]);

  // ── Fullscreen API ──────────────────────────────────────
  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };
    document.addEventListener("fullscreenchange", handleFullscreenChange);
    return () => document.removeEventListener("fullscreenchange", handleFullscreenChange);
  }, []);

  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen().catch((err) => {
        console.error(`Error attempting to enable fullscreen: ${err.message}`);
      });
    } else {
      if (document.exitFullscreen) {
        document.exitFullscreen();
      }
    }
  };

  // ── Anti-Cheat: Visibility Change ───────────────────────
  useEffect(() => {
    if (!session?.sections || submitting) return undefined;
    
    const handleVisibilityChange = () => {
      if (document.hidden) {
        showSnackbar("error", "warning", "Peringatan Anti-Cheat!", "Berpindah tab atau window terdeteksi. Harap fokus pada ujian!");
      }
    };
    
    document.addEventListener("visibilitychange", handleVisibilityChange);
    return () => document.removeEventListener("visibilitychange", handleVisibilityChange);
  }, [session, submitting, showSnackbar]);

  // ── Start exam ──────────────────────────────────────────
  const start = async () => {
    setStarting(true);
    setError("");
    try {
      const response = await apiRequest(`/participant/periods/${id}/start`, { method: "POST" });
      setSession(response.data);
      if (response.data?.end_time) setDeadline(new Date(response.data.end_time));
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setStarting(false);
    }
  };

  // ── Answer handler ──────────────────────────────────────
  const answer = async (questionId, optionId) => {
    const sequence = (requestSequence.current[questionId] || 0) + 1;
    requestSequence.current[questionId] = sequence;
    setAnswers((current) => ({ ...current, [questionId]: optionId }));
    setSaveStatus("saving");
    try {
      await apiRequest(`/participant/periods/${id}/answers`, {
        method: "POST",
        body: JSON.stringify({ question_public_id: questionId, option_public_id: optionId }),
      });
      if (requestSequence.current[questionId] === sequence) setSaveStatus("saved");
    } catch (requestError) {
      if (requestSequence.current[questionId] === sequence) setSaveStatus("failed");
      setError(requestError.message);
    }
  };

  // ── Toggle Flag ─────────────────────────────────────────
  const toggleFlag = (questionId) => {
    setFlagged((prev) => ({ ...prev, [questionId]: !prev[questionId] }));
  };

  // ── Format time ─────────────────────────────────────────
  const formatTime = (milliseconds) => {
    if (milliseconds === null) return "--:--:--";
    const seconds = Math.floor(milliseconds / 1000);
    const h = String(Math.floor(seconds / 3600)).padStart(2, "0");
    const m = String(Math.floor((seconds % 3600) / 60)).padStart(2, "0");
    const s = String(seconds % 60).padStart(2, "0");
    return `${h}:${m}:${s}`;
  };

  // ── Derived data ────────────────────────────────────────
  const sections = session?.sections || [];
  const currentSection = sections[activeSection];
  const allQuestions = sections.flatMap((s) => s.questions || []);
  const totalQuestions = allQuestions.length;
  const answeredCount = allQuestions.filter((q) => answers[q.public_id]).length;
  const isUrgent = timeLeft !== null && timeLeft <= 300000; // 5 minutes
  const isVeryUrgent = timeLeft !== null && timeLeft <= 60000; // 1 minute

  // ── Loading state ───────────────────────────────────────
  if (loading) {
    return (
      <DashboardLayout>
        <DashboardNavbar />
        <SoftBox pt={6} pb={3}>
          <Skeleton variant="text" width={300} height={40} />
          <Skeleton variant="text" width={200} height={24} sx={{ mt: 1 }} />
          <Skeleton variant="rounded" height={400} sx={{ mt: 3 }} />
        </SoftBox>
      </DashboardLayout>
    );
  }

  // ── Pre-start view ──────────────────────────────────────
  if (!session?.sections) {
    return (
      <DashboardLayout>
        <DashboardNavbar />
        <SoftBox pt={6} pb={3}>
          <Card>
            <SoftBox p={5} textAlign="center">
              <Icon color="info" sx={{ fontSize: 64 }}>quiz</Icon>
              <SoftTypography variant="h4" fontWeight="bold" mt={2}>
                {session?.title || "Pengerjaan Ujian"}
              </SoftTypography>
              <SoftTypography variant="button" color="text" display="block" mt={1}>
                Timer akan mengikuti waktu selesai periode dari server.
              </SoftTypography>
              <SoftTypography variant="button" color="text" display="block" mt={0.5}>
                Pastikan koneksi internet Anda stabil sebelum memulai.
              </SoftTypography>

              {error && (
                <SoftBox mt={3} p={2} borderRadius="md" bgColor="error">
                  <SoftTypography variant="button" color="white">{error}</SoftTypography>
                </SoftBox>
              )}

              <SoftButton
                variant="gradient"
                color="info"
                size="large"
                sx={{ mt: 4, px: 6 }}
                onClick={start}
                disabled={starting}
              >
                <Icon sx={{ mr: 1 }}>{starting ? "hourglass_top" : "play_arrow"}</Icon>
                {starting ? "Menyiapkan..." : "Mulai Ujian"}
              </SoftButton>
            </SoftBox>
          </Card>
        </SoftBox>
        <Footer />
      </DashboardLayout>
    );
  }

  // ── Exam view ───────────────────────────────────────────
  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox pt={6} pb={3}>
        {/* Header with timer */}
        <SoftBox
          display="flex"
          justifyContent="space-between"
          alignItems={{ xs: "flex-start", md: "center" }}
          flexDirection={{ xs: "column", md: "row" }}
          gap={2}
          mb={3}
        >
          <SoftBox>
            <SoftTypography variant="h4" fontWeight="bold">
              {session?.title || "Pengerjaan Ujian"}
            </SoftTypography>
            <SoftTypography variant="button" color="text">
              Jawaban tersimpan otomatis ke server.
            </SoftTypography>
          </SoftBox>
          <SoftBox display="flex" alignItems="center" gap={1} flexWrap="wrap">
            {/* Fullscreen Button */}
            <Tooltip title={isFullscreen ? "Keluar Mode Fokus" : "Masuk Mode Fokus"}>
              <IconButton 
                onClick={toggleFullscreen} 
                color="info" 
                size="medium"
                sx={{ bgcolor: "info.main", color: "white", "&:hover": { bgcolor: "info.dark" } }}
              >
                <Icon>{isFullscreen ? "fullscreen_exit" : "fullscreen"}</Icon>
              </IconButton>
            </Tooltip>
            {/* Timer chip */}
            <Chip
              icon={<Icon>{isVeryUrgent ? "warning" : isUrgent ? "timer_off" : "timer"}</Icon>}
              label={formatTime(timeLeft)}
              color={isUrgent ? "error" : "info"}
              variant="filled"
              sx={{ 
                fontWeight: "bold", 
                fontSize: "1rem", 
                height: 40,
                animation: isVeryUrgent ? "pulse 1s infinite" : "none",
                "@keyframes pulse": {
                  "0%": { opacity: 1 },
                  "50%": { opacity: 0.5 },
                  "100%": { opacity: 1 }
                }
              }}
            />
            {/* Save status */}
            <Chip
              label={
                saveStatus === "saving" ? "Menyimpan..." :
                saveStatus === "saved" ? "Tersimpan ✓" :
                saveStatus === "failed" ? "Gagal ✗" : "—"
              }
              color={saveStatus === "failed" ? "error" : "default"}
              size="small"
            />
            {/* Progress */}
            <Chip
              label={`${answeredCount}/${totalQuestions} dijawab`}
              color={answeredCount === totalQuestions ? "success" : "default"}
              size="small"
            />
          </SoftBox>
        </SoftBox>

        {/* Error */}
        {error && (
          <SoftBox mb={3} p={2} borderRadius="md" bgColor="error">
            <SoftTypography variant="button" color="white">{error}</SoftTypography>
          </SoftBox>
        )}

        <Grid container spacing={3}>
          {/* Main content */}
          <Grid item xs={12} md={9}>
            {/* Section tabs */}
            {sections.length > 1 && (
              <Card sx={{ mb: 3 }}>
                <Tabs
                  value={activeSection}
                  onChange={(_, newValue) => setActiveSection(newValue)}
                  variant="scrollable"
                  scrollButtons="auto"
                  sx={{ px: 2 }}
                >
                  {sections.map((section, index) => (
                    <Tab
                      key={section.public_id}
                      label={section.title}
                      icon={<Icon>{index <= activeSection ? "check_circle" : "radio_button_unchecked"}</Icon>}
                      iconPosition="start"
                    />
                  ))}
                </Tabs>
              </Card>
            )}

            {/* Questions */}
            {currentSection && (
              <Card>
                <SoftBox p={3}>
                  <SoftTypography variant="h5" fontWeight="bold" mb={2}>
                    {currentSection.title}
                  </SoftTypography>

                  {(currentSection.questions || []).map((question, index) => (
                    <SoftBox
                      key={question.public_id}
                      py={3}
                      borderBottom="1px solid"
                      borderColor="grey.200"
                      id={`question-${question.public_id}`}
                    >
                      <SoftBox display="flex" justifyContent="space-between" alignItems="flex-start" gap={2}>
                        <SoftTypography variant="button" fontWeight="bold">
                          {question.number || index + 1}. {question.question}
                        </SoftTypography>
                        <Tooltip title={flagged[question.public_id] ? "Hapus tanda ragu-ragu" : "Tandai ragu-ragu"}>
                          <IconButton
                            size="small"
                            color={flagged[question.public_id] ? "warning" : "default"}
                            onClick={() => toggleFlag(question.public_id)}
                          >
                            <Icon>{flagged[question.public_id] ? "flag" : "outlined_flag"}</Icon>
                          </IconButton>
                        </Tooltip>
                      </SoftBox>

                      {question.passage && (
                        <SoftBox mt={2} p={2} bgColor="grey.100" borderRadius="md">
                          <div dangerouslySetInnerHTML={{ __html: question.passage }} />
                        </SoftBox>
                      )}

                      {question.audio_url && (
                        <SoftBox mt={2}>
                          <CustomAudioPlayer src={question.audio_url} />
                        </SoftBox>
                      )}

                      {question.image_url && (
                        <SoftBox mt={2}>
                          <img
                            src={question.image_url}
                            alt="Media soal"
                            style={{ maxWidth: "100%", maxHeight: 280, borderRadius: 8 }}
                          />
                        </SoftBox>
                      )}

                      <RadioGroup
                        value={answers[question.public_id] || ""}
                        onChange={(event) => answer(question.public_id, event.target.value)}
                        sx={{ mt: 2 }}
                      >
                        {(question.options || []).map((option) => (
                          <FormControlLabel
                            key={option.public_id}
                            value={option.public_id}
                            control={<Radio />}
                            label={`${option.label}. ${option.option_text}`}
                            sx={{
                              p: 1,
                              borderRadius: 1,
                              mb: 0.5,
                              transition: "background-color 0.2s",
                              backgroundColor: answers[question.public_id] === option.public_id ? "rgba(33, 150, 243, 0.08)" : "transparent",
                              "&:hover": { backgroundColor: "rgba(0,0,0,0.04)" },
                            }}
                          />
                        ))}
                      </RadioGroup>
                    </SoftBox>
                  ))}
                </SoftBox>
              </Card>
            )}

            {/* Submit button */}
            <SoftBox mt={3} display="flex" justifyContent="flex-end">
              <SoftButton
                variant="gradient"
                color="error"
                size="large"
                onClick={() => setConfirmOpen(true)}
                disabled={submitting || timeLeft === 0}
              >
                <Icon sx={{ mr: 1 }}>send</Icon>
                {submitting ? "Mengirim..." : "Submit Ujian"}
              </SoftButton>
            </SoftBox>
          </Grid>

          {/* Sidebar: question navigator */}
          <Grid item xs={12} md={3}>
            <Card sx={{ position: "sticky", top: 100 }}>
              <SoftBox p={2}>
                <SoftTypography variant="button" fontWeight="bold" display="block" mb={1}>
                  Navigasi Soal
                </SoftTypography>
                <SoftTypography variant="caption" color="text" display="block" mb={2}>
                  {answeredCount} dari {totalQuestions} terjawab
                </SoftTypography>

                {sections.map((section, sIdx) => (
                  <SoftBox key={section.public_id} mb={2}>
                    <SoftTypography variant="caption" fontWeight="bold" color="text" display="block" mb={1}>
                      {section.title}
                    </SoftTypography>
                    <SoftBox display="flex" flexWrap="wrap" gap={0.5}>
                      {(section.questions || []).map((question, qIdx) => {
                        const isAnswered = Boolean(answers[question.public_id]);
                        const isCurrent = sIdx === activeSection;
                        return (
                          <SoftBox
                            key={question.public_id}
                            component="button"
                            type="button"
                            onClick={() => {
                              setActiveSection(sIdx);
                              setTimeout(() => {
                                document.getElementById(`question-${question.public_id}`)?.scrollIntoView({ behavior: "smooth", block: "center" });
                              }, 100);
                            }}
                            sx={{
                              width: 36,
                              height: 36,
                              borderRadius: 1,
                              border: isCurrent ? "2px solid" : "1px solid",
                              borderColor: flagged[question.public_id] ? "warning.main" : isAnswered ? "success.main" : isCurrent ? "info.main" : "grey.300",
                              backgroundColor: flagged[question.public_id] ? "warning.main" : isAnswered ? "success.main" : "transparent",
                              color: (flagged[question.public_id] || isAnswered) ? "#fff" : "text.primary",
                              cursor: "pointer",
                              fontWeight: "bold",
                              fontSize: 12,
                              display: "flex",
                              alignItems: "center",
                              justifyContent: "center",
                              transition: "all 0.2s",
                              "&:hover": { transform: "scale(1.1)", boxShadow: 2 },
                            }}
                          >
                            {question.number || qIdx + 1}
                          </SoftBox>
                        );
                      })}
                    </SoftBox>
                  </SoftBox>
                ))}
              </SoftBox>
            </Card>
          </Grid>
        </Grid>
      </SoftBox>
      <Footer />

      {/* ── Submit Confirmation Dialog ─────────────────────── */}
      <Dialog open={confirmOpen} onClose={() => setConfirmOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Konfirmasi Submit Ujian</DialogTitle>
        <DialogContent dividers>
          <SoftBox textAlign="center" py={2}>
            <Icon color={answeredCount < totalQuestions ? "warning" : "success"} sx={{ fontSize: 48 }}>
              {answeredCount < totalQuestions ? "warning" : "check_circle"}
            </Icon>
            <SoftTypography variant="h6" mt={2}>
              {answeredCount < totalQuestions
                ? `Anda baru menjawab ${answeredCount} dari ${totalQuestions} soal.`
                : `Semua ${totalQuestions} soal sudah dijawab.`}
            </SoftTypography>
            <SoftTypography variant="button" color="text" display="block" mt={1}>
              Jawaban tidak dapat diubah setelah submit. Yakin ingin mengirim sekarang?
            </SoftTypography>

            {/* Per-section breakdown */}
            <SoftBox mt={3}>
              {sections.map((section) => {
                const sectionQuestions = section.questions || [];
                const sectionAnswered = sectionQuestions.filter((q) => answers[q.public_id]).length;
                return (
                  <SoftBox key={section.public_id} display="flex" justifyContent="space-between" px={4} py={0.5}>
                    <SoftTypography variant="caption" color="text">{section.title}</SoftTypography>
                    <SoftTypography variant="caption" fontWeight="bold" color={sectionAnswered === sectionQuestions.length ? "success" : "warning"}>
                      {sectionAnswered}/{sectionQuestions.length}
                    </SoftTypography>
                  </SoftBox>
                );
              })}
            </SoftBox>
          </SoftBox>
        </DialogContent>
        <DialogActions sx={{ px: 3, pb: 3 }}>
          <SoftButton color="secondary" onClick={() => setConfirmOpen(false)}>
            Kembali ke ujian
          </SoftButton>
          <SoftButton
            variant="gradient"
            color="error"
            onClick={() => { setConfirmOpen(false); submit(false); }}
            disabled={submitting}
          >
            <Icon sx={{ mr: 1 }}>send</Icon>
            {submitting ? "Mengirim..." : "Submit sekarang"}
          </SoftButton>
        </DialogActions>
      </Dialog>
      
      {/* ── Global Exam Snackbar ───────────────────────────── */}
      <SoftSnackbar
        color={snackbarState.color}
        icon={snackbarState.icon}
        title={snackbarState.title}
        content={snackbarState.content}
        dateTime="Sekarang"
        open={snackbarState.open}
        close={closeSnackbar}
      />
    </DashboardLayout>
  );
}

export default ParticipantExam;
