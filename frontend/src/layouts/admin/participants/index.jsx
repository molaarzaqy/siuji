import { useCallback, useEffect, useMemo, useState } from "react";
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
import SoftSelect from "components/SoftSelect";
import SoftTypography from "components/SoftTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import MiniStatisticsCard from "examples/Cards/StatisticsCards/MiniStatisticsCard";
import { apiRequest } from "services/api";
import * as XLSX from "xlsx";

const statusMeta = {
  registered: { label: "Registered", color: "info" },
  started: { label: "Started", color: "warning" },
  completed: { label: "Completed", color: "success" },
};

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

const emptyParticipant = { name: "", email: "", nim: "", university: "" };

function Participants() {
  const [periods, setPeriods] = useState([]);
  const [periodId, setPeriodId] = useState("");
  const [participants, setParticipants] = useState([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(false);
  const [exporting, setExporting] = useState(false);
  const [error, setError] = useState("");
  const [page, setPage] = useState(1);
  const [meta, setMeta] = useState({});
  const rowsPerPage = 10;

  // Import
  const [file, setFile] = useState(null);
  const [importing, setImporting] = useState(false);
  const [importResult, setImportResult] = useState(null);

  // Add participant dialog
  const [addOpen, setAddOpen] = useState(false);
  const [addForm, setAddForm] = useState(emptyParticipant);
  const [saving, setSaving] = useState(false);

  // Detail dialog
  const [detailOpen, setDetailOpen] = useState(false);
  const [detailData, setDetailData] = useState(null);

  // Edit dialog
  const [editOpen, setEditOpen] = useState(false);
  const [editForm, setEditForm] = useState(emptyParticipant);
  const [editUserId, setEditUserId] = useState(null);

  // ── Load periods ────────────────────────────────────────
  useEffect(() => {
    apiRequest("/periods/?page=1&limit=100")
      .then((response) => {
        const data = response.data || [];
        setPeriods(data);
        if (data[0]) setPeriodId(data[0].public_id);
      })
      .catch((requestError) => setError(requestError.message));
  }, []);

  // ── Load participants ───────────────────────────────────
  const load = useCallback(async () => {
    if (!periodId) return;
    setLoading(true);
    setError("");
    try {
      const response = await apiRequest(
        `/periods/${periodId}/participants?page=${page}&limit=${rowsPerPage}&filter=${encodeURIComponent(search)}`
      );
      setParticipants(response.data || []);
      setMeta(response.meta || {});
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setLoading(false);
    }
  }, [periodId, page, search]);

  useEffect(() => {
    load();
  }, [load]);

  useEffect(() => {
    setPage(1);
  }, [periodId, search]);

  // ── Derived stats ───────────────────────────────────────
  const counts = useMemo(() => ({
    total: meta.total_datas ?? participants.length,
    registered: participants.filter((p) => p.status === "registered").length,
    started: participants.filter((p) => p.status === "started").length,
    completed: participants.filter((p) => p.status === "completed").length,
  }), [participants, meta]);

  const totalPages = meta.total_pages || Math.ceil((meta.total_datas || participants.length) / rowsPerPage);

  // ── Import Excel ────────────────────────────────────────
  const importFile = async () => {
    if (!file || !periodId) return;
    setImporting(true);
    setError("");
    setImportResult(null);
    const body = new FormData();
    body.append("file", file);
    try {
      const response = await apiRequest(`/periods/${periodId}/participants/import`, {
        method: "POST",
        body,
      });
      setImportResult(response.data);
      setFile(null);
      await load();
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setImporting(false);
    }
  };

  // ── Export Excel ──────────────────────────────────────────
  const exportExcel = async () => {
    if (!periodId) return;
    setExporting(true);
    setError("");
    try {
      const response = await apiRequest(`/periods/${periodId}/participants?page=1&limit=10000`);
      const data = response.data || [];
      if (!data.length) {
        setError("Tidak ada data peserta untuk diekspor.");
        return;
      }

      // Convert data to Excel format
      const worksheetData = data.map((p) => ({
        "Nama": p.user?.name || "",
        "Email": p.user?.email || "",
        "NIM": p.user?.nim || "",
        "Universitas": p.user?.university || "",
        "Status": p.status || "",
        "Nilai": p.score ?? ""
      }));

      const worksheet = XLSX.utils.json_to_sheet(worksheetData);
      const workbook = XLSX.utils.book_new();
      XLSX.utils.book_append_sheet(workbook, worksheet, "Laporan Peserta");
      
      // Auto-size columns slightly
      worksheet["!cols"] = [
        { wch: 25 }, { wch: 30 }, { wch: 15 }, { wch: 30 }, { wch: 15 }, { wch: 10 }
      ];

      XLSX.writeFile(workbook, `Laporan_Peserta.xlsx`);
    } catch (requestError) {
      setError("Gagal mengekspor data: " + requestError.message);
    } finally {
      setExporting(false);
    }
  };

  // ── Add participant ─────────────────────────────────────
  const submitAdd = async (event) => {
    event.preventDefault();
    setSaving(true);
    setError("");
    try {
      await apiRequest(`/periods/${periodId}/participants`, {
        method: "POST",
        body: JSON.stringify(addForm),
      });
      setAddOpen(false);
      setAddForm(emptyParticipant);
      await load();
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setSaving(false);
    }
  };

  // ── Edit participant ────────────────────────────────────
  const openEdit = (participant) => {
    setEditUserId(participant.user?.public_id || participant.public_id);
    setEditForm({
      name: participant.user?.name || "",
      email: participant.user?.email || "",
      nim: participant.user?.nim || "",
      university: participant.user?.university || "",
    });
    setEditOpen(true);
  };

  const submitEdit = async (event) => {
    event.preventDefault();
    setSaving(true);
    setError("");
    try {
      await apiRequest(`/periods/${periodId}/participants/${editUserId}`, {
        method: "PUT",
        body: JSON.stringify(editForm),
      });
      setEditOpen(false);
      await load();
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setSaving(false);
    }
  };

  // ── View detail ─────────────────────────────────────────
  const viewDetail = async (participant) => {
    const userId = participant.user?.public_id || participant.public_id;
    try {
      const response = await apiRequest(`/periods/${periodId}/participants/${userId}`);
      setDetailData(response.data);
      setDetailOpen(true);
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  // ── Delete participant ──────────────────────────────────
  const removeParticipant = async (participant) => {
    const name = participant.user?.name || participant.user?.email || "peserta ini";
    if (!window.confirm(`Hapus ${name} dari periode ini?`)) return;
    const userId = participant.user?.public_id || participant.public_id;
    try {
      await apiRequest(`/periods/${periodId}/participants/${userId}`, {
        method: "DELETE",
      });
      await load();
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox pt={6} pb={3}>
        {/* Header */}
        <SoftBox mb={3} display="flex" justifyContent="space-between" alignItems={{ xs: "flex-start", sm: "center" }} flexDirection={{ xs: "column", sm: "row" }} gap={2}>
          <SoftBox>
            <SoftTypography variant="h4" fontWeight="bold">Manajemen Peserta</SoftTypography>
            <SoftTypography variant="button" color="text" fontWeight="regular">
              Assign peserta ke periode dan import data secara massal.
            </SoftTypography>
          </SoftBox>
          <SoftBox display="flex" gap={1}>
            <SoftButton variant="gradient" color="success" onClick={exportExcel} disabled={!periodId || exporting}>
              <Icon sx={{ mr: 1 }}>download</Icon>{exporting ? "Mengekspor..." : "Export Laporan"}
            </SoftButton>
            <SoftButton variant="gradient" color="info" onClick={() => { setAddForm(emptyParticipant); setAddOpen(true); }} disabled={!periodId}>
              <Icon sx={{ mr: 1 }}>person_add</Icon>Tambah peserta
            </SoftButton>
          </SoftBox>
        </SoftBox>

        {/* Stats cards */}
        <Grid container spacing={3} mb={3}>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard title={{ text: "total peserta", fontWeight: "bold" }} count={counts.total} percentage={{ color: "info", text: "terdaftar" }} icon={{ color: "info", component: "people" }} />
          </Grid>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard title={{ text: "registered", fontWeight: "bold" }} count={counts.registered} percentage={{ color: "info", text: "belum mulai" }} icon={{ color: "info", component: "how_to_reg" }} />
          </Grid>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard title={{ text: "started", fontWeight: "bold" }} count={counts.started} percentage={{ color: "warning", text: "mengerjakan" }} icon={{ color: "warning", component: "play_circle" }} />
          </Grid>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard title={{ text: "completed", fontWeight: "bold" }} count={counts.completed} percentage={{ color: "success", text: "selesai" }} icon={{ color: "success", component: "check_circle" }} />
          </Grid>
        </Grid>

        {/* Main card */}
        <Card>
          <SoftBox p={3}>
            <Grid container spacing={2} alignItems="center">
              <Grid item xs={12} md={5}>
                <SoftBox mb={1} ml={0.5} lineHeight={0} display="inline-block">
                  <SoftTypography component="label" variant="caption" fontWeight="bold">
                    Periode Ujian
                  </SoftTypography>
                </SoftBox>
                <SoftSelect
                  placeholder="Pilih periode ujian..."
                  options={periods.map((period) => ({
                    value: period.public_id,
                    label: `${period.title} — ${period.month} ${period.year}`
                  }))}
                  value={
                    periodId
                      ? {
                          value: periodId,
                          label: periods
                            .filter(p => p.public_id === periodId)
                            .map(p => `${p.title} — ${p.month} ${p.year}`)[0]
                        }
                      : null
                  }
                  onChange={(option) => setPeriodId(option ? option.value : "")}
                />
              </Grid>
              <Grid item xs={12} md={7}>
                <SoftBox mb={1} ml={0.5} lineHeight={0} display="inline-block">
                  <SoftTypography component="label" variant="caption" fontWeight="bold">
                    Cari Peserta
                  </SoftTypography>
                </SoftBox>
                <SoftInput icon={{ component: "search", direction: "left" }} placeholder="Cari berdasarkan nama, email, atau institusi..." value={search} onChange={(event) => setSearch(event.target.value)} />
              </Grid>
              <Grid item xs={12}>
                <SoftBox display="flex" gap={2} alignItems="center" flexWrap="wrap">
                  <input
                    accept=".xlsx"
                    style={{ display: "none" }}
                    id="excel-file-upload"
                    type="file"
                    onChange={(event) => setFile(event.target.files[0] || null)}
                  />
                  <label htmlFor="excel-file-upload" style={{ margin: 0 }}>
                    <SoftButton variant="outlined" color="dark" component="span">
                      <Icon sx={{ mr: 1 }}>attach_file</Icon>
                      {file ? "Ganti File" : "Pilih File Excel"}
                    </SoftButton>
                  </label>
                  
                  {file && (
                    <SoftTypography variant="button" color="text" fontWeight="medium">
                      {file.name}
                    </SoftTypography>
                  )}

                  <SoftButton variant="gradient" color="info" onClick={importFile} disabled={!file || importing}>
                    <Icon sx={{ mr: 1 }}>upload_file</Icon>
                    {importing ? "Mengimpor..." : "Import Excel"}
                  </SoftButton>
                  
                  {!file && (
                    <SoftTypography variant="caption" color="text">
                      Format: Name, Email, NIM, University
                    </SoftTypography>
                  )}

                  <SoftButton 
                    variant="text" 
                    color="info" 
                    component="a" 
                    href="/template_peserta.xlsx" 
                    download="template_peserta.xlsx"
                  >
                    <Icon sx={{ mr: 1 }}>download</Icon>
                    Download Template
                  </SoftButton>
                </SoftBox>
              </Grid>
            </Grid>
          </SoftBox>

          {/* Messages */}
          {error && (
            <SoftBox mx={3} mb={2} p={2} borderRadius="md" bgColor="error">
              <SoftTypography variant="button" color="white">{error}</SoftTypography>
            </SoftBox>
          )}
          {importResult && (
            <SoftBox mx={3} mb={2} p={2} borderRadius="md" bgColor="success">
              <SoftTypography variant="button" color="white">
                Berhasil: {importResult.total_imported}, dilewati: {importResult.total_skipped}
              </SoftTypography>
            </SoftBox>
          )}

          <Divider />

          {/* Table */}
          {loading ? (
            <SoftBox p={3}>
              <Skeleton height={64} />
              <Skeleton height={64} />
              <Skeleton height={64} />
            </SoftBox>
          ) : (
            <TableContainer sx={{ overflowX: "auto", borderRadius: 2 }}>
              <Table sx={{ minWidth: 760, display: "table", ...globalTableStyle }}>
                <TableHead sx={{ display: "table-header-group" }}>
                  <TableRow sx={{ display: "table-row" }}>
                    <TableCell>Peserta</TableCell>
                    <TableCell>Email</TableCell>
                    <TableCell>NIM / Universitas</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell>Nilai</TableCell>
                    <TableCell align="right">Aksi</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {participants.map((participant) => {
                    const pStatus = statusMeta[participant.status] || { label: participant.status || "-", color: "secondary" };
                    return (
                      <TableRow hover key={participant.public_id}>
                        <TableCell>
                          <SoftTypography variant="button" fontWeight="bold">
                            {participant.user?.name || "-"}
                          </SoftTypography>
                        </TableCell>
                        <TableCell>
                          <SoftTypography variant="caption" color="text">
                            {participant.user?.email || "-"}
                          </SoftTypography>
                        </TableCell>
                        <TableCell>
                          <SoftTypography variant="caption" display="block">{participant.user?.nim || "-"}</SoftTypography>
                          <SoftTypography variant="caption" color="text">{participant.user?.university || "-"}</SoftTypography>
                        </TableCell>
                        <TableCell>
                          <SoftBadge color={pStatus.color} variant="gradient" size="sm">{pStatus.label}</SoftBadge>
                        </TableCell>
                        <TableCell>
                          <SoftTypography variant="button" fontWeight="medium">
                            {participant.score ?? "-"}
                          </SoftTypography>
                        </TableCell>
                        <TableCell align="right">
                          <SoftBox display="flex" justifyContent="flex-end" gap={0.5}>
                            <Tooltip title="Lihat detail">
                              <IconButton size="small" color="info" onClick={() => viewDetail(participant)}>
                                <Icon fontSize="small">visibility</Icon>
                              </IconButton>
                            </Tooltip>
                            <Tooltip title="Edit peserta">
                              <IconButton size="small" color="warning" onClick={() => openEdit(participant)}>
                                <Icon fontSize="small">edit</Icon>
                              </IconButton>
                            </Tooltip>
                            <Tooltip title="Hapus peserta">
                              <IconButton size="small" color="error" onClick={() => removeParticipant(participant)}>
                                <Icon fontSize="small">delete</Icon>
                              </IconButton>
                            </Tooltip>
                          </SoftBox>
                        </TableCell>
                      </TableRow>
                    );
                  })}
                  {!participants.length && (
                    <TableRow>
                      <TableCell colSpan={6}>
                        <SoftBox py={6} textAlign="center">
                          <Icon color="disabled" sx={{ fontSize: 48 }}>people_outline</Icon>
                          <SoftTypography variant="h6" mt={1}>Belum ada peserta</SoftTypography>
                          <SoftTypography variant="button" color="text">
                            Pilih periode dan tambahkan peserta melalui form atau import Excel.
                          </SoftTypography>
                        </SoftBox>
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </TableContainer>
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <SoftBox display="flex" justifyContent="space-between" alignItems="center" p={3} borderTop="1px solid" borderColor="grey.200">
              <SoftTypography variant="button" color="secondary" fontWeight="regular">
                Halaman {page} dari {totalPages} ({counts.total} peserta)
              </SoftTypography>
              <SoftPagination variant="gradient" color="info">
                {page > 1 && (
                  <SoftPagination item onClick={() => setPage(page - 1)}>
                    <Icon sx={{ fontWeight: "bold" }}>chevron_left</Icon>
                  </SoftPagination>
                )}
                {Array.from({ length: Math.min(totalPages, 6) }, (_, i) => i + 1).map((p) => (
                  <SoftPagination item key={p} onClick={() => setPage(p)} active={page === p}>{p}</SoftPagination>
                ))}
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

      {/* ── Add Participant Dialog ─────────────────────────── */}
      <Dialog open={addOpen} onClose={() => !saving && setAddOpen(false)} fullWidth maxWidth="sm">
        <SoftBox component="form" onSubmit={submitAdd}>
          <DialogTitle>Tambah Peserta</DialogTitle>
          <DialogContent dividers>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Nama</SoftTypography>
                <SoftInput value={addForm.name} onChange={(e) => setAddForm({ ...addForm, name: e.target.value })} required />
              </Grid>
              <Grid item xs={12}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Email</SoftTypography>
                <SoftInput type="email" value={addForm.email} onChange={(e) => setAddForm({ ...addForm, email: e.target.value })} required />
              </Grid>
              <Grid item xs={12} md={6}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">NIM</SoftTypography>
                <SoftInput value={addForm.nim} onChange={(e) => setAddForm({ ...addForm, nim: e.target.value })} required />
              </Grid>
              <Grid item xs={12} md={6}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Universitas</SoftTypography>
                <SoftInput value={addForm.university} onChange={(e) => setAddForm({ ...addForm, university: e.target.value })} required />
              </Grid>
            </Grid>
            <SoftTypography variant="caption" color="text" display="block" mt={2}>
              Password peserta akan dibuat otomatis berdasarkan NIM.
            </SoftTypography>
          </DialogContent>
          <DialogActions>
            <SoftButton color="secondary" onClick={() => setAddOpen(false)}>Batal</SoftButton>
            <SoftButton type="submit" variant="gradient" color="info" disabled={saving}>
              {saving ? "Menyimpan..." : "Simpan"}
            </SoftButton>
          </DialogActions>
        </SoftBox>
      </Dialog>

      {/* ── Edit Participant Dialog ────────────────────────── */}
      <Dialog open={editOpen} onClose={() => !saving && setEditOpen(false)} fullWidth maxWidth="sm">
        <SoftBox component="form" onSubmit={submitEdit}>
          <DialogTitle>Edit Peserta</DialogTitle>
          <DialogContent dividers>
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Nama</SoftTypography>
                <SoftInput value={editForm.name} onChange={(e) => setEditForm({ ...editForm, name: e.target.value })} required />
              </Grid>
              <Grid item xs={12}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Email</SoftTypography>
                <SoftInput type="email" value={editForm.email} onChange={(e) => setEditForm({ ...editForm, email: e.target.value })} required />
              </Grid>
              <Grid item xs={12} md={6}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">NIM</SoftTypography>
                <SoftInput value={editForm.nim} onChange={(e) => setEditForm({ ...editForm, nim: e.target.value })} />
              </Grid>
              <Grid item xs={12} md={6}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Universitas</SoftTypography>
                <SoftInput value={editForm.university} onChange={(e) => setEditForm({ ...editForm, university: e.target.value })} />
              </Grid>
            </Grid>
          </DialogContent>
          <DialogActions>
            <SoftButton color="secondary" onClick={() => setEditOpen(false)}>Batal</SoftButton>
            <SoftButton type="submit" variant="gradient" color="info" disabled={saving}>
              {saving ? "Menyimpan..." : "Simpan perubahan"}
            </SoftButton>
          </DialogActions>
        </SoftBox>
      </Dialog>

      {/* ── Detail Participant Dialog ──────────────────────── */}
      <Dialog open={detailOpen} onClose={() => setDetailOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>Detail Peserta</DialogTitle>
        <DialogContent dividers>
          {detailData && (
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Nama</SoftTypography>
                <SoftTypography variant="button" display="block">{detailData.user?.name || detailData.name || "-"}</SoftTypography>
              </Grid>
              <Grid item xs={12}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Email</SoftTypography>
                <SoftTypography variant="button" display="block">{detailData.user?.email || detailData.email || "-"}</SoftTypography>
              </Grid>
              <Grid item xs={6}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">NIM</SoftTypography>
                <SoftTypography variant="button" display="block">{detailData.user?.nim || detailData.nim || "-"}</SoftTypography>
              </Grid>
              <Grid item xs={6}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Universitas</SoftTypography>
                <SoftTypography variant="button" display="block">{detailData.user?.university || detailData.university || "-"}</SoftTypography>
              </Grid>
              <Grid item xs={6}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Status</SoftTypography>
                <SoftBox mt={0.5}>
                  <SoftBadge
                    color={statusMeta[detailData.status]?.color || "secondary"}
                    variant="gradient"
                    size="sm"
                  >
                    {statusMeta[detailData.status]?.label || detailData.status || "-"}
                  </SoftBadge>
                </SoftBox>
              </Grid>
              <Grid item xs={6}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Nilai</SoftTypography>
                <SoftTypography variant="h5" display="block">{detailData.score ?? "-"}</SoftTypography>
              </Grid>
            </Grid>
          )}
        </DialogContent>
        <DialogActions>
          <SoftButton color="light" onClick={() => setDetailOpen(false)}>Tutup</SoftButton>
        </DialogActions>
      </Dialog>
    </DashboardLayout>
  );
}

export default Participants;
