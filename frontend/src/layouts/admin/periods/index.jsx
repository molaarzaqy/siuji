import { useCallback, useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import PropTypes from "prop-types";
import Card from "@mui/material/Card";
import Chip from "@mui/material/Chip";
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
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import TextField from "@mui/material/TextField";
import Tooltip from "@mui/material/Tooltip";
import SoftBadge from "components/SoftBadge";
import SoftBox from "components/SoftBox";
import SoftButton from "components/SoftButton";
import SoftInput from "components/SoftInput";
import SoftPagination from "components/SoftPagination";
import SoftTypography from "components/SoftTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import MiniStatisticsCard from "examples/Cards/StatisticsCards/MiniStatisticsCard";
import DefaultDoughnutChart from "examples/Charts/DoughnutCharts/DefaultDoughnutChart";
import { apiRequest } from "services/api";

const globalTableStyle = {
  "& th": { 
    backgroundColor: "grey.100", 
    color: "secondary.main", 
    fontSize: "0.75rem", 
    fontWeight: 700, 
    textTransform: "uppercase", 
    letterSpacing: "0.05em",
    borderBottom: "1px solid",
    borderColor: "grey.200",
    py: 1.5,
    px: 3
  }, 
  "& td": { 
    px: 3, 
    py: 2.5, 
    verticalAlign: "middle",
    borderBottom: "1px solid",
    borderColor: "grey.200"
  },
  "& tbody tr:hover": {
    backgroundColor: "grey.50"
  }
};

const initialForm = { title: "", month: "", year: new Date().getFullYear(), status: "draft", start_time: "", end_time: "", min_passing_grade: 0, max_passing_grade: 677, certificate_template: null };
const toRFC3339 = (value) => (value ? new Date(value).toISOString() : "");
const statusMeta = {
  draft: { label: "Draft", color: "warning" },
  published: { label: "Published", color: "success" },
  closed: { label: "Closed", color: "dark" },
};

function formatDate(value, withTime = false) {
  if (!value) return "-";
  return new Date(value).toLocaleString("id-ID", withTime ? { dateStyle: "medium", timeStyle: "short" } : { dateStyle: "medium" });
}

function StatusBadge({ status }) {
  const meta = statusMeta[status] || { label: status || "Unknown", color: "secondary" };
  return <SoftBadge color={meta.color} variant="gradient" size="sm">{meta.label}</SoftBadge>;
}

StatusBadge.propTypes = { status: PropTypes.string };

function Periods() {
  const [periods, setPeriods] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("all");
  const [dialogOpen, setDialogOpen] = useState(false);
  const [saving, setSaving] = useState(false);
  const [form, setForm] = useState(initialForm);
  const [page, setPage] = useState(1);
  const rowsPerPage = 10;
  const navigate = useNavigate();

  const loadPeriods = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const response = await apiRequest("/periods/?page=1&limit=100&sort=-created_at");
      setPeriods(response.data || []);
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { loadPeriods(); }, [loadPeriods]);

  const counts = useMemo(() => ({
    all: periods.length,
    published: periods.filter((period) => period.status === "published").length,
    draft: periods.filter((period) => period.status === "draft").length,
    closed: periods.filter((period) => period.status === "closed").length,
  }), [periods]);

  const activePeriod = periods.find((period) => period.status === "published" && new Date(period.start_time) <= new Date() && new Date(period.end_time) >= new Date());
  const filteredPeriods = periods.filter((period) => {
    const matchesStatus = statusFilter === "all" || period.status === statusFilter;
    const query = search.trim().toLowerCase();
    return matchesStatus && (!query || `${period.title} ${period.month} ${period.year}`.toLowerCase().includes(query));
  });

  useEffect(() => {
    setPage(1);
  }, [search, statusFilter]);

  const totalPages = Math.ceil(filteredPeriods.length / rowsPerPage);
  const paginatedPeriods = filteredPeriods.slice((page - 1) * rowsPerPage, page * rowsPerPage);

  const updateField = (field) => (event) => {
    const value = field === "certificate_template" ? event.target.files[0] : event.target.value;
    setForm((current) => ({ ...current, [field]: value }));
  };

  const handleCreate = async (event) => {
    event.preventDefault();
    if (!form.certificate_template) return setError("Template sertifikat wajib dipilih.");
    if (new Date(form.end_time) <= new Date(form.start_time)) return setError("Waktu selesai harus setelah waktu mulai.");
    setSaving(true);
    setError("");
    const payload = new FormData();
    Object.entries(form).forEach(([key, value]) => {
      if (key === "certificate_template") payload.append(key, value);
      else if (key === "start_time" || key === "end_time") payload.append(key, toRFC3339(value));
      else payload.append(key, value);
    });
    try {
      await apiRequest("/periods/", { method: "POST", body: payload });
      setDialogOpen(false);
      setForm(initialForm);
      await loadPeriods();
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (period) => {
    if (!window.confirm(`Hapus periode "${period.title}"?`)) return;
    try {
      await apiRequest(`/periods/${period.public_id}`, { method: "DELETE" });
      await loadPeriods();
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox pt={6} pb={3}>
        <SoftBox mb={3} display="flex" justifyContent="space-between" alignItems={{ xs: "flex-start", sm: "center" }} flexDirection={{ xs: "column", sm: "row" }} gap={2}>
          <SoftBox>
            <SoftTypography variant="h4" fontWeight="bold">Manajemen Periode Ujian</SoftTypography>
            <SoftTypography variant="button" color="text" fontWeight="regular">Kelola jadwal dan status periode ujian SIUJI dari satu tempat.</SoftTypography>
          </SoftBox>
          <SoftButton variant="gradient" color="info" onClick={() => { setError(""); setDialogOpen(true); }}><Icon sx={{ mr: 1 }}>add</Icon>Tambah periode</SoftButton>
        </SoftBox>

        <Grid container spacing={3} mb={3}>
          <Grid item xs={12} lg={4}>
            <DefaultDoughnutChart
              title="Komposisi Status Periode"
              description="Perbandingan jumlah periode berdasarkan status"
              chart={{
                labels: ["Published", "Draft", "Closed"],
                datasets: {
                  label: "Jumlah",
                  backgroundColors: ["success", "warning", "dark"],
                  data: [counts.published || 0, counts.draft || 0, counts.closed || 0],
                },
              }}
            />
          </Grid>
          <Grid item xs={12} lg={8}>
            <Grid container spacing={3}>
              <Grid item xs={12} sm={6}><MiniStatisticsCard title={{ text: "semua periode", fontWeight: "bold" }} count={counts.all} percentage={{ color: "info", text: "total" }} icon={{ color: "info", component: "event" }} /></Grid>
              <Grid item xs={12} sm={6}><MiniStatisticsCard title={{ text: "published", fontWeight: "bold" }} count={counts.published} percentage={{ color: "success", text: "siap ujian" }} icon={{ color: "success", component: "check_circle" }} /></Grid>
              <Grid item xs={12} sm={6}><MiniStatisticsCard title={{ text: "draft", fontWeight: "bold" }} count={counts.draft} percentage={{ color: "warning", text: "perlu disiapkan" }} icon={{ color: "warning", component: "edit_calendar" }} /></Grid>
              <Grid item xs={12} sm={6}><MiniStatisticsCard title={{ text: "periode aktif", fontWeight: "bold" }} count={activePeriod ? 1 : 0} percentage={{ color: "success", text: activePeriod ? "sedang berjalan" : "tidak ada" }} icon={{ color: "dark", component: "schedule" }} /></Grid>
            </Grid>
          </Grid>
        </Grid>

        <Card>
          <SoftBox p={3} pb={1}>
            <Grid container spacing={2} alignItems="center">
              <Grid item xs={12} md={6}><SoftInput icon={{ component: "search", direction: "left" }} placeholder="Cari judul, bulan, atau tahun..." value={search} onChange={(event) => setSearch(event.target.value)} /></Grid>
              <Grid item xs={12} md={6} display="flex" justifyContent={{ xs: "flex-start", md: "flex-end" }} gap={1} flexWrap="wrap">
                {[["all", "Semua"], ["published", "Published"], ["draft", "Draft"], ["closed", "Closed"]].map(([value, label]) => <Chip key={value} label={label} color={statusFilter === value ? "info" : "default"} variant={statusFilter === value ? "filled" : "outlined"} onClick={() => setStatusFilter(value)} />)}
                <Tooltip title="Muat ulang data"><IconButton onClick={loadPeriods} disabled={loading}><Icon>refresh</Icon></IconButton></Tooltip>
              </Grid>
            </Grid>
          </SoftBox>
          {error && <SoftBox mx={3} mb={2} p={2} borderRadius="md" bgColor="error"><SoftTypography variant="button" color="white">{error}</SoftTypography></SoftBox>}
          <Divider />
          {loading ? <SoftBox p={3}>{[1, 2, 3].map((item) => <Skeleton key={item} height={64} />)}</SoftBox> : (
            <TableContainer sx={{ width: "100%", overflowX: "auto", borderRadius: 2 }}>
              <Table sx={{ width: "100%", minWidth: { xs: 860, md: "100%" }, tableLayout: "fixed", display: "table", ...globalTableStyle }}>
                <colgroup><col style={{ width: "28%" }} /><col style={{ width: "14%" }} /><col style={{ width: "30%" }} /><col style={{ width: "16%" }} /><col style={{ width: "12%" }} /></colgroup>
                <TableHead sx={{ display: "table-header-group" }}>
                  <TableRow sx={{ display: "table-row" }}><TableCell align="left">Periode</TableCell><TableCell align="center">Status</TableCell><TableCell align="left">Jadwal</TableCell><TableCell align="center">Passing grade</TableCell><TableCell align="right">Aksi</TableCell></TableRow></TableHead>
                <TableBody>
                  {paginatedPeriods.map((period) => <TableRow hover key={period.public_id}>
                    <TableCell><SoftTypography variant="button" fontWeight="bold" color="dark" sx={{ display: "block", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{period.title}</SoftTypography><SoftTypography variant="caption" display="block" color="text">{period.month} {period.year}</SoftTypography></TableCell>
                    <TableCell align="center"><StatusBadge status={period.status} /></TableCell>
                    <TableCell><SoftTypography variant="caption" display="block" color="text" sx={{ whiteSpace: "nowrap" }}>Mulai: {formatDate(period.start_time, true)}</SoftTypography><SoftTypography variant="caption" display="block" color="text" sx={{ whiteSpace: "nowrap", mt: 0.5 }}>Selesai: {formatDate(period.end_time, true)}</SoftTypography></TableCell>
                    <TableCell align="center"><SoftTypography variant="button" fontWeight="medium" sx={{ whiteSpace: "nowrap" }}>{period.min_passing_grade ?? 0} - {period.max_passing_grade ?? 0}</SoftTypography></TableCell>
                    <TableCell align="right"><SoftBox display="flex" justifyContent="flex-end" gap={0.5}><Tooltip title="Kelola & Edit Periode"><IconButton size="small" color="info" onClick={() => navigate(`/admin/periods/${period.public_id}`)}><Icon fontSize="small">edit</Icon></IconButton></Tooltip><Tooltip title="Hapus periode"><IconButton size="small" color="error" onClick={() => handleDelete(period)}><Icon fontSize="small">delete</Icon></IconButton></Tooltip></SoftBox></TableCell>
                  </TableRow>)}
                  {!filteredPeriods.length && <TableRow><TableCell colSpan={5}><SoftBox py={6} textAlign="center"><Icon color="disabled" sx={{ fontSize: 48 }}>event_busy</Icon><SoftTypography variant="h6" mt={1}>Belum ada periode yang sesuai</SoftTypography><SoftTypography variant="button" color="text">Coba ubah filter atau tambahkan periode baru.</SoftTypography></SoftBox></TableCell></TableRow>}
                </TableBody>
              </Table>
            </TableContainer>
          )}
          {totalPages > 1 && (
            <SoftBox display="flex" justifyContent="space-between" alignItems="center" p={3} borderTop="1px solid" borderColor="grey.200">
              <SoftTypography variant="button" color="secondary" fontWeight="regular">
                Menampilkan {((page - 1) * rowsPerPage) + 1} hingga {Math.min(page * rowsPerPage, filteredPeriods.length)} dari {filteredPeriods.length} entri
              </SoftTypography>
              <SoftPagination variant="gradient" color="info">
                {page > 1 && (
                  <SoftPagination item onClick={() => setPage(page - 1)}>
                    <Icon sx={{ fontWeight: "bold" }}>chevron_left</Icon>
                  </SoftPagination>
                )}
                {totalPages > 6 ? (
                  <SoftBox width="5rem" mx={1}>
                    <SoftInput
                      inputProps={{ type: "number", min: 1, max: totalPages }}
                      value={page}
                      onChange={(e) => {
                        const val = Number(e.target.value);
                        if (val >= 1 && val <= totalPages) setPage(val);
                      }}
                    />
                  </SoftBox>
                ) : (
                  Array.from({ length: totalPages }, (_, i) => i + 1).map((p) => (
                    <SoftPagination item key={p} onClick={() => setPage(p)} active={page === p}>
                      {p}
                    </SoftPagination>
                  ))
                )}
                {page < totalPages && (
                  <SoftPagination item onClick={() => setPage(page + 1)}>
                    <Icon sx={{ fontWeight: "bold" }}>chevron_right</Icon>
                  </SoftPagination>
                )}
              </SoftPagination>
            </SoftBox>
          )}
        </Card>
      </SoftBox>
      <Footer />

      <Dialog open={dialogOpen} onClose={() => !saving && setDialogOpen(false)} fullWidth maxWidth="md">
        <DialogTitle>Tambah periode ujian</DialogTitle>
        <SoftBox component="form" onSubmit={handleCreate}>
          <DialogContent><Grid container spacing={2}>
            <Grid item xs={12}><SoftInput placeholder="Judul periode" value={form.title} onChange={updateField("title")} required /></Grid>
            <Grid item xs={12} md={6}><SoftInput placeholder="Bulan, contoh: September" value={form.month} onChange={updateField("month")} required /></Grid>
            <Grid item xs={12} md={6}><SoftInput type="number" placeholder="Tahun" value={form.year} onChange={updateField("year")} required /></Grid>
            <Grid item xs={12} md={4}><TextField select fullWidth label="Status" value={form.status} onChange={updateField("status")}><MenuItem value="draft">Draft</MenuItem><MenuItem value="published">Published</MenuItem><MenuItem value="closed">Closed</MenuItem></TextField></Grid>
            <Grid item xs={12} md={4}><TextField fullWidth label="Mulai" type="datetime-local" value={form.start_time} onChange={updateField("start_time")} InputLabelProps={{ shrink: true }} required /></Grid>
            <Grid item xs={12} md={4}><TextField fullWidth label="Selesai" type="datetime-local" value={form.end_time} onChange={updateField("end_time")} InputLabelProps={{ shrink: true }} required /></Grid>
            <Grid item xs={12} md={6}><SoftInput type="number" placeholder="Nilai minimum lulus" value={form.min_passing_grade} onChange={updateField("min_passing_grade")} /></Grid>
            <Grid item xs={12} md={6}><SoftInput type="number" placeholder="Nilai maksimum" value={form.max_passing_grade} onChange={updateField("max_passing_grade")} /></Grid>
            <Grid item xs={12}><TextField fullWidth type="file" label="Template sertifikat" onChange={updateField("certificate_template")} InputLabelProps={{ shrink: true }} required /></Grid>
          </Grid></DialogContent>
          <DialogActions sx={{ px: 3, pb: 3 }}><SoftButton color="light" onClick={() => setDialogOpen(false)} disabled={saving}>Batal</SoftButton><SoftButton type="submit" variant="gradient" color="info" disabled={saving}>{saving ? "Menyimpan..." : "Simpan periode"}</SoftButton></DialogActions>
        </SoftBox>
      </Dialog>
    </DashboardLayout>
  );
}

export default Periods;
