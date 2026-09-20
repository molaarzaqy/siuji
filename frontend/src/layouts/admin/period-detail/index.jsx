import { useCallback, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import Card from "@mui/material/Card";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import Divider from "@mui/material/Divider";
import Grid from "@mui/material/Grid";
import Icon from "@mui/material/Icon";
import IconButton from "@mui/material/IconButton";
import MenuItem from "@mui/material/MenuItem";
import Skeleton from "@mui/material/Skeleton";
import TextField from "@mui/material/TextField";
import Tooltip from "@mui/material/Tooltip";
import SoftBadge from "components/SoftBadge";
import SoftBox from "components/SoftBox";
import SoftButton from "components/SoftButton";
import SoftInput from "components/SoftInput";
import SoftTypography from "components/SoftTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import MiniStatisticsCard from "examples/Cards/StatisticsCards/MiniStatisticsCard";
import { apiRequest } from "services/api";

const statusMeta = {
  draft: { label: "Draft", color: "warning" },
  published: { label: "Published", color: "success" },
  closed: { label: "Closed", color: "dark" },
};

const toRFC3339 = (value) => (value ? new Date(value).toISOString() : "");

function toLocalDatetime(rfc3339) {
  if (!rfc3339) return "";
  const date = new Date(rfc3339);
  const offset = date.getTimezoneOffset();
  const local = new Date(date.getTime() - offset * 60000);
  return local.toISOString().slice(0, 16);
}

function formatDate(value) {
  if (!value) return "-";
  return new Date(value).toLocaleString("id-ID", { dateStyle: "medium", timeStyle: "short" });
}

function PeriodDetail() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [period, setPeriod] = useState(null);
  const [sections, setSections] = useState([]);
  const [available, setAvailable] = useState([]);
  const [selected, setSelected] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  // Edit dialog
  const [editOpen, setEditOpen] = useState(false);
  const [editForm, setEditForm] = useState({});
  const [saving, setSaving] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const [periodResponse, sectionResponse] = await Promise.all([
        apiRequest(`/periods/${id}`),
        apiRequest("/sections/"),
      ]);
      const detail = periodResponse.data;
      setPeriod(detail);
      setSections(detail.sections || []);
      setAvailable(sectionResponse.data || []);
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  // ── Section management ──────────────────────────────────
  const addSection = async () => {
    if (!selected) return;
    try {
      await apiRequest(`/periods/${id}/sections`, {
        method: "POST",
        body: JSON.stringify({ 
          section_public_id: selected,
          position: sections.length + 1
        }),
      });
      setSelected("");
      await load();
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  const removeSection = async (section) => {
    if (!window.confirm(`Hapus "${section.title}" dari periode ini?`)) return;
    try {
      await apiRequest(`/periods/${id}/sections/${section.section_public_id}`, {
        method: "DELETE",
      });
      await load();
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  const reorder = async (from, to) => {
    const next = [...sections];
    const [item] = next.splice(from, 1);
    next.splice(to, 0, item);
    setSections(next);
    try {
      await apiRequest(`/periods/${id}/sections/reorder`, {
        method: "PUT",
        body: JSON.stringify({
          section_public_ids: next.map((section) => section.section_public_id),
        }),
      });
    } catch (requestError) {
      setError(requestError.message);
      await load();
    }
  };

  // ── Edit period ─────────────────────────────────────────
  const openEdit = () => {
    setEditForm({
      title: period.title || "",
      month: period.month || "",
      year: period.year || new Date().getFullYear(),
      status: period.status || "draft",
      start_time: toLocalDatetime(period.start_time),
      end_time: toLocalDatetime(period.end_time),
      min_passing_grade: period.min_passing_grade ?? 0,
      max_passing_grade: period.max_passing_grade ?? 677,
      certificate_template: null,
    });
    setEditOpen(true);
  };

  const updateEditField = (field) => (event) => {
    const value = field === "certificate_template" ? event.target.files[0] : event.target.value;
    setEditForm((current) => ({ ...current, [field]: value }));
  };

  const submitEdit = async (event) => {
    event.preventDefault();
    if (new Date(editForm.end_time) <= new Date(editForm.start_time)) {
      setError("Waktu selesai harus setelah waktu mulai.");
      return;
    }
    setSaving(true);
    setError("");
    const payload = new FormData();
    Object.entries(editForm).forEach(([key, value]) => {
      if (key === "certificate_template") {
        if (value) payload.append(key, value);
      } else if (key === "start_time" || key === "end_time") {
        payload.append(key, toRFC3339(value));
      } else {
        payload.append(key, value);
      }
    });
    try {
      await apiRequest(`/periods/${id}`, { method: "PUT", body: payload });
      setEditOpen(false);
      await load();
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setSaving(false);
    }
  };

  // ── Derived data ────────────────────────────────────────
  const meta = statusMeta[period?.status] || { label: period?.status || "Unknown", color: "secondary" };
  const availableToAdd = available.filter(
    (item) => !sections.some((section) => section.section_public_id === item.public_id)
  );

  // ── Loading skeleton ────────────────────────────────────
  if (loading) {
    return (
      <DashboardLayout>
        <DashboardNavbar />
        <SoftBox pt={6} pb={3}>
          <Skeleton variant="text" width={200} height={40} />
          <Skeleton variant="text" width={300} height={24} sx={{ mt: 1 }} />
          <Grid container spacing={3} mt={1}>
            {[1, 2, 3].map((item) => (
              <Grid item xs={12} sm={4} key={item}>
                <Skeleton variant="rounded" height={100} />
              </Grid>
            ))}
          </Grid>
          <Skeleton variant="rounded" height={300} sx={{ mt: 3 }} />
        </SoftBox>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox pt={6} pb={3}>
        {/* Header */}
        <SoftBox display="flex" justifyContent="space-between" alignItems={{ xs: "flex-start", sm: "center" }} flexDirection={{ xs: "column", sm: "row" }} gap={2} mb={3}>
          <SoftBox display="flex" alignItems="center" gap={2}>
            <SoftButton variant="outlined" color="secondary" size="small" onClick={() => navigate("/admin/periods")}>
              <Icon sx={{ mr: 0.5 }}>arrow_back</Icon>Kembali
            </SoftButton>
            <SoftBox>
              <SoftTypography variant="h4" fontWeight="bold">
                {period?.title || "Detail Periode"}
              </SoftTypography>
              <SoftTypography variant="button" color="text" fontWeight="regular">
                Konfigurasi section dan struktur ujian.
              </SoftTypography>
            </SoftBox>
          </SoftBox>
          <SoftButton variant="gradient" color="warning" onClick={openEdit}>
            <Icon sx={{ mr: 1 }}>edit</Icon>Edit periode
          </SoftButton>
        </SoftBox>

        {/* Error */}
        {error && (
          <SoftBox mb={3} p={2} borderRadius="md" bgColor="error">
            <SoftTypography variant="button" color="white">{error}</SoftTypography>
          </SoftBox>
        )}

        {/* Stats cards */}
        <Grid container spacing={3} mb={3}>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard
              title={{ text: "status", fontWeight: "bold" }}
              count={meta.label}
              percentage={{ color: meta.color, text: period?.month + " " + period?.year }}
              icon={{ color: meta.color, component: "flag" }}
            />
          </Grid>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard
              title={{ text: "jumlah section", fontWeight: "bold" }}
              count={sections.length}
              percentage={{ color: "info", text: "terstruktur" }}
              icon={{ color: "info", component: "view_module" }}
            />
          </Grid>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard
              title={{ text: "waktu mulai", fontWeight: "bold" }}
              count={formatDate(period?.start_time)}
              percentage={{ color: "success", text: "" }}
              icon={{ color: "success", component: "schedule" }}
            />
          </Grid>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard
              title={{ text: "passing grade", fontWeight: "bold" }}
              count={`${period?.min_passing_grade ?? 0} - ${period?.max_passing_grade ?? 0}`}
              percentage={{ color: "dark", text: "rentang nilai" }}
              icon={{ color: "dark", component: "grading" }}
            />
          </Grid>
        </Grid>

        {/* Period summary card */}
        <Card sx={{ mb: 3 }}>
          <SoftBox p={3}>
            <SoftBox display="flex" justifyContent="space-between" alignItems="center">
              <SoftTypography variant="h6" fontWeight="bold">Ringkasan Periode</SoftTypography>
              <SoftBadge color={meta.color} variant="gradient" size="sm">{meta.label}</SoftBadge>
            </SoftBox>
            <Divider />
            <Grid container spacing={2} mt={0}>
              <Grid item xs={12} md={3}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Judul</SoftTypography>
                <SoftTypography variant="button" display="block">{period?.title}</SoftTypography>
              </Grid>
              <Grid item xs={6} md={2}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Periode</SoftTypography>
                <SoftTypography variant="button" display="block">{period?.month} {period?.year}</SoftTypography>
              </Grid>
              <Grid item xs={6} md={2}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Mulai</SoftTypography>
                <SoftTypography variant="button" display="block">{formatDate(period?.start_time)}</SoftTypography>
              </Grid>
              <Grid item xs={6} md={2}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Selesai</SoftTypography>
                <SoftTypography variant="button" display="block">{formatDate(period?.end_time)}</SoftTypography>
              </Grid>
              <Grid item xs={6} md={3}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Passing Grade</SoftTypography>
                <SoftTypography variant="button" display="block">
                  {period?.min_passing_grade ?? 0} — {period?.max_passing_grade ?? 0}
                </SoftTypography>
              </Grid>
            </Grid>
          </SoftBox>
        </Card>

        {/* Sections management */}
        <Card>
          <SoftBox p={3}>
            <SoftBox display="flex" justifyContent="space-between" alignItems="center">
              <SoftBox>
                <SoftTypography variant="h6" fontWeight="bold">Section dalam Periode</SoftTypography>
                <SoftTypography variant="caption" color="text">
                  {sections.length} section terdaftar. Gunakan tombol panah untuk mengubah urutan.
                </SoftTypography>
              </SoftBox>
            </SoftBox>
          </SoftBox>
          <Divider />

          {/* Section list */}
          <SoftBox p={3}>
            {sections.length > 0 ? (
              sections.map((section, index) => (
                <SoftBox
                  key={section.section_public_id}
                  display="flex"
                  alignItems="center"
                  justifyContent="space-between"
                  p={2}
                  mb={1}
                  borderRadius="md"
                  bgColor="grey.100"
                  sx={{ transition: "all 0.2s", "&:hover": { bgColor: "grey.200", transform: "translateX(4px)" } }}
                >
                  <SoftBox display="flex" alignItems="center" gap={2}>
                    <SoftBox
                      display="flex"
                      alignItems="center"
                      justifyContent="center"
                      width={36}
                      height={36}
                      borderRadius="md"
                      bgColor="info"
                      color="white"
                    >
                      <SoftTypography variant="button" fontWeight="bold" color="white">
                        {index + 1}
                      </SoftTypography>
                    </SoftBox>
                    <SoftBox>
                      <SoftTypography variant="button" fontWeight="bold">{section.title}</SoftTypography>
                      <SoftTypography variant="caption" display="block" color="text">
                        {section.section_type || "section"}
                      </SoftTypography>
                    </SoftBox>
                  </SoftBox>
                  <SoftBox display="flex" alignItems="center" gap={0.5}>
                    <Tooltip title="Pindah ke atas">
                      <span>
                        <IconButton size="small" disabled={index === 0} onClick={() => reorder(index, index - 1)}>
                          <Icon>arrow_upward</Icon>
                        </IconButton>
                      </span>
                    </Tooltip>
                    <Tooltip title="Pindah ke bawah">
                      <span>
                        <IconButton size="small" disabled={index === sections.length - 1} onClick={() => reorder(index, index + 1)}>
                          <Icon>arrow_downward</Icon>
                        </IconButton>
                      </span>
                    </Tooltip>
                    <Tooltip title="Hapus dari periode">
                      <IconButton size="small" color="error" onClick={() => removeSection(section)}>
                        <Icon>delete</Icon>
                      </IconButton>
                    </Tooltip>
                  </SoftBox>
                </SoftBox>
              ))
            ) : (
              <SoftBox py={6} textAlign="center">
                <Icon color="disabled" sx={{ fontSize: 48 }}>view_module</Icon>
                <SoftTypography variant="h6" mt={1}>Belum ada section</SoftTypography>
                <SoftTypography variant="button" color="text">
                  Tambahkan section dari dropdown di bawah.
                </SoftTypography>
              </SoftBox>
            )}
          </SoftBox>

          <Divider />

          {/* Add section */}
          <SoftBox p={3} display="flex" gap={2} alignItems="center" flexWrap="wrap">
            <TextField
              select
              label="Pilih section"
              value={selected}
              onChange={(event) => setSelected(event.target.value)}
              sx={{ minWidth: 280 }}
            >
              {availableToAdd.length > 0 ? (
                availableToAdd.map((section) => (
                  <MenuItem key={section.public_id} value={section.public_id}>
                    {section.title} — {section.section_type}
                  </MenuItem>
                ))
              ) : (
                <MenuItem disabled>Semua section sudah ditambahkan</MenuItem>
              )}
            </TextField>
            <SoftButton variant="gradient" color="info" onClick={addSection} disabled={!selected}>
              <Icon sx={{ mr: 1 }}>add</Icon>Tambahkan section
            </SoftButton>
          </SoftBox>
        </Card>
      </SoftBox>
      <Footer />

      {/* ── Edit Period Dialog ─────────────────────────────── */}
      <Dialog open={editOpen} onClose={() => !saving && setEditOpen(false)} fullWidth maxWidth="md">
        <SoftBox component="form" onSubmit={submitEdit}>
          <DialogTitle>Edit Periode Ujian</DialogTitle>
          <DialogContent dividers>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Judul periode</SoftTypography>
                <SoftInput value={editForm.title || ""} onChange={updateEditField("title")} required />
              </Grid>
              <Grid item xs={12} md={6}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Bulan</SoftTypography>
                <SoftInput value={editForm.month || ""} onChange={updateEditField("month")} placeholder="contoh: September" required />
              </Grid>
              <Grid item xs={12} md={6}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Tahun</SoftTypography>
                <SoftInput type="number" value={editForm.year || ""} onChange={updateEditField("year")} required />
              </Grid>
              <Grid item xs={12} md={4}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Status</SoftTypography>
                <TextField select fullWidth value={editForm.status || "draft"} onChange={updateEditField("status")} sx={{ mt: 0.5 }}>
                  <MenuItem value="draft">Draft</MenuItem>
                  <MenuItem value="published">Published</MenuItem>
                  <MenuItem value="closed">Closed</MenuItem>
                </TextField>
              </Grid>
              <Grid item xs={12} md={4}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Mulai</SoftTypography>
                <TextField fullWidth type="datetime-local" value={editForm.start_time || ""} onChange={updateEditField("start_time")} InputLabelProps={{ shrink: true }} sx={{ mt: 0.5 }} required />
              </Grid>
              <Grid item xs={12} md={4}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Selesai</SoftTypography>
                <TextField fullWidth type="datetime-local" value={editForm.end_time || ""} onChange={updateEditField("end_time")} InputLabelProps={{ shrink: true }} sx={{ mt: 0.5 }} required />
              </Grid>
              <Grid item xs={12} md={6}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Nilai minimum lulus</SoftTypography>
                <SoftInput type="number" value={editForm.min_passing_grade ?? 0} onChange={updateEditField("min_passing_grade")} />
              </Grid>
              <Grid item xs={12} md={6}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Nilai maksimum</SoftTypography>
                <SoftInput type="number" value={editForm.max_passing_grade ?? 0} onChange={updateEditField("max_passing_grade")} />
              </Grid>
              <Grid item xs={12}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">
                  Template sertifikat (opsional, kosongkan jika tidak ingin mengubah)
                </SoftTypography>
                <TextField fullWidth type="file" onChange={updateEditField("certificate_template")} InputLabelProps={{ shrink: true }} sx={{ mt: 0.5 }} />
              </Grid>
            </Grid>
          </DialogContent>
          <DialogActions sx={{ px: 3, pb: 3 }}>
            <SoftButton color="secondary" onClick={() => setEditOpen(false)} disabled={saving}>Batal</SoftButton>
            <SoftButton type="submit" variant="gradient" color="info" disabled={saving}>
              {saving ? "Menyimpan..." : "Simpan perubahan"}
            </SoftButton>
          </DialogActions>
        </SoftBox>
      </Dialog>
    </DashboardLayout>
  );
}

export default PeriodDetail;
