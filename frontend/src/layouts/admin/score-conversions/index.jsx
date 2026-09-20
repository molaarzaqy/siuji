import { useCallback, useEffect, useMemo, useState } from "react";
import Card from "@mui/material/Card";
import Chip from "@mui/material/Chip";
import Divider from "@mui/material/Divider";
import Grid from "@mui/material/Grid";
import Icon from "@mui/material/Icon";
import IconButton from "@mui/material/IconButton";
import MenuItem from "@mui/material/MenuItem";
import Skeleton from "@mui/material/Skeleton";
import TextField from "@mui/material/TextField";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Tooltip from "@mui/material/Tooltip";
import SoftBox from "components/SoftBox";
import SoftButton from "components/SoftButton";
import SoftInput from "components/SoftInput";
import SoftSelect from "components/SoftSelect";
import SoftTypography from "components/SoftTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import MiniStatisticsCard from "examples/Cards/StatisticsCards/MiniStatisticsCard";
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
  },
  "& thead": { display: "table-header-group" },
  "& tbody": { display: "table-row-group" },
  "& tr": { display: "table-row" },
  "& th, & td": { display: "table-cell" }
};

const sectionLabels = {
  listening: "Listening",
  structure: "Structure",
  reading: "Reading",
};

function ScoreConversions() {
  const [type, setType] = useState("listening");
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [bulkSaving, setBulkSaving] = useState(false);
  const [message, setMessage] = useState("");
  const [messageType, setMessageType] = useState("error");

  // ── Data loading ────────────────────────────────────────
  const load = useCallback(async () => {
    setLoading(true);
    setMessage("");
    try {
      const response = await apiRequest(`/score-conversions/?section_type=${type}`);
      setRows((response.data || []).map((row) => ({ ...row })));
    } catch (requestError) {
      setMessage(requestError.message);
      setMessageType("error");
    } finally {
      setLoading(false);
    }
  }, [type]);

  useEffect(() => {
    load();
  }, [load]);

  // ── Stats ───────────────────────────────────────────────
  const stats = useMemo(() => ({
    total: rows.length,
    minScore: rows.length ? Math.min(...rows.map((r) => Number(r.scaled_score) || 0)) : 0,
    maxScore: rows.length ? Math.max(...rows.map((r) => Number(r.scaled_score) || 0)) : 0,
  }), [rows]);

  // ── Row editing ─────────────────────────────────────────
  const updateRow = (id, field, value) => {
    setRows((current) =>
      current.map((row) => (row.id === id ? { ...row, [field]: value } : row))
    );
  };

  // ── Validate duplicates ─────────────────────────────────
  const hasDuplicates = useMemo(() => {
    const counts = {};
    for (const row of rows) {
      const key = Number(row.correct_count);
      if (counts[key]) return true;
      counts[key] = true;
    }
    return false;
  }, [rows]);

  // ── Save single row ────────────────────────────────────
  const saveRow = async (row) => {
    setSaving(true);
    setMessage("");
    try {
      await apiRequest(`/score-conversions/${row.id}`, {
        method: "PUT",
        body: JSON.stringify({
          section_type: row.section_type,
          correct_count: Number(row.correct_count),
          scaled_score: Number(row.scaled_score),
        }),
      });
      showSuccess("Konversi nilai berhasil disimpan.");
    } catch (requestError) {
      setMessage(requestError.message);
      setMessageType("error");
    } finally {
      setSaving(false);
    }
  };

  // ── Add single row ─────────────────────────────────────
  const addRow = async () => {
    setSaving(true);
    setMessage("");
    try {
      await apiRequest("/score-conversions/", {
        method: "POST",
        body: JSON.stringify({ section_type: type, correct_count: 0, scaled_score: 0 }),
      });
      await load();
      showSuccess("Mapping baru berhasil ditambahkan.");
    } catch (requestError) {
      setMessage(requestError.message);
      setMessageType("error");
    } finally {
      setSaving(false);
    }
  };

  // ── Delete single row ──────────────────────────────────
  const deleteRow = async (row) => {
    if (!window.confirm(`Hapus mapping jawaban benar ${row.correct_count} → nilai ${row.scaled_score}?`)) return;
    try {
      await apiRequest(`/score-conversions/${row.id}`, { method: "DELETE" });
      await load();
      showSuccess("Mapping berhasil dihapus.");
    } catch (requestError) {
      setMessage(requestError.message);
      setMessageType("error");
    }
  };

  // ── Bulk save ───────────────────────────────────────────
  const bulkSave = async () => {
    if (hasDuplicates) {
      setMessage("Terdapat duplikasi jumlah jawaban benar. Perbaiki sebelum menyimpan.");
      setMessageType("error");
      return;
    }
    setBulkSaving(true);
    setMessage("");
    try {
      await apiRequest("/score-conversions/bulk", {
        method: "POST",
        body: JSON.stringify({
          conversions: rows.map((row) => ({
            section_type: type,
            correct_count: Number(row.correct_count),
            scaled_score: Number(row.scaled_score),
          })),
        }),
      });
      showSuccess("Semua mapping berhasil disimpan secara bulk.");
      await load();
    } catch (requestError) {
      setMessage(requestError.message);
      setMessageType("error");
    } finally {
      setBulkSaving(false);
    }
  };

  const showSuccess = (text) => {
    setMessage(text);
    setMessageType("success");
  };

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox pt={6} pb={3}>
        {/* Header */}
        <SoftBox mb={3}>
          <SoftTypography variant="h4" fontWeight="bold">Konversi Nilai</SoftTypography>
          <SoftTypography variant="button" color="text" fontWeight="regular">
            Atur mapping jawaban benar ke nilai skala TOEFL.
          </SoftTypography>
        </SoftBox>

        {/* Stats cards */}
        <Grid container spacing={3} mb={3}>
          <Grid item xs={12} sm={4}>
            <MiniStatisticsCard
              title={{ text: `mapping ${sectionLabels[type]}`, fontWeight: "bold" }}
              count={stats.total}
              percentage={{ color: "info", text: "baris" }}
              icon={{ color: "info", component: "grid_on" }}
            />
          </Grid>
          <Grid item xs={12} sm={4}>
            <MiniStatisticsCard
              title={{ text: "nilai minimum", fontWeight: "bold" }}
              count={stats.minScore}
              percentage={{ color: "warning", text: "skala terendah" }}
              icon={{ color: "warning", component: "trending_down" }}
            />
          </Grid>
          <Grid item xs={12} sm={4}>
            <MiniStatisticsCard
              title={{ text: "nilai maksimum", fontWeight: "bold" }}
              count={stats.maxScore}
              percentage={{ color: "success", text: "skala tertinggi" }}
              icon={{ color: "success", component: "trending_up" }}
            />
          </Grid>
        </Grid>

        {/* Main card */}
        <Card>
          <SoftBox p={3}>
            <Grid container spacing={2} alignItems="center">
              <Grid item xs={12} md={4}>
                <SoftBox mb={1} ml={0.5} lineHeight={0} display="inline-block">
                  <SoftTypography component="label" variant="caption" fontWeight="bold">Section</SoftTypography>
                </SoftBox>
                <SoftSelect
                  placeholder="Pilih section..."
                  options={Object.entries(sectionLabels).map(([value, label]) => ({value, label}))}
                  value={{value: type, label: sectionLabels[type]}}
                  onChange={(option) => setType(option ? option.value : "")}
                />
              </Grid>
              <Grid item xs={12} md={8} display="flex" justifyContent={{ xs: "flex-start", md: "flex-end" }} gap={1} flexWrap="wrap">
                <SoftButton variant="gradient" color="info" onClick={addRow} disabled={saving}>
                  <Icon sx={{ mr: 1 }}>add</Icon>Tambah mapping
                </SoftButton>
                <SoftButton variant="gradient" color="success" onClick={bulkSave} disabled={bulkSaving || !rows.length}>
                  <Icon sx={{ mr: 1 }}>save</Icon>
                  {bulkSaving ? "Menyimpan..." : "Simpan semua (bulk)"}
                </SoftButton>
              </Grid>
            </Grid>
          </SoftBox>

          {/* Duplicate warning */}
          {hasDuplicates && (
            <SoftBox mx={3} mb={2} p={2} borderRadius="md" bgColor="warning">
              <SoftTypography variant="button" color="white">
                ⚠️ Terdapat duplikasi jumlah jawaban benar. Perbaiki sebelum menyimpan bulk.
              </SoftTypography>
            </SoftBox>
          )}

          {/* Messages */}
          {message && (
            <SoftBox mx={3} mb={2} p={2} borderRadius="md" bgColor={messageType}>
              <SoftTypography variant="button" color="white">{message}</SoftTypography>
            </SoftBox>
          )}

          <Divider />

          {/* Section type chips for quick switch */}
          <SoftBox px={3} pt={2} display="flex" gap={1}>
            {Object.entries(sectionLabels).map(([value, label]) => (
              <Chip
                key={value}
                label={label}
                color={type === value ? "info" : "default"}
                variant={type === value ? "filled" : "outlined"}
                onClick={() => setType(value)}
              />
            ))}
          </SoftBox>

          {/* Table */}
          {loading ? (
            <SoftBox p={3}>
              <Skeleton height={60} />
              <Skeleton height={60} />
              <Skeleton height={60} />
            </SoftBox>
          ) : (
            <TableContainer sx={{ overflowX: "auto", borderRadius: 2 }}>
              <Table sx={{ minWidth: 720, display: "table", ...globalTableStyle }}>
                <TableHead sx={{ display: "table-header-group" }}>
                  <TableRow sx={{ display: "table-row" }}>
                    <TableCell sx={{ width: "10%" }}>#</TableCell>
                    <TableCell>Jawaban benar</TableCell>
                    <TableCell>Nilai skala</TableCell>
                    <TableCell align="right">Aksi</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {rows.map((row, index) => (
                    <TableRow key={row.id}>
                      <TableCell>
                        <SoftTypography variant="caption" fontWeight="bold" color="text">
                          {index + 1}
                        </SoftTypography>
                      </TableCell>
                      <TableCell>
                        <SoftInput
                          type="number"
                          value={row.correct_count}
                          onChange={(event) => updateRow(row.id, "correct_count", event.target.value)}
                        />
                      </TableCell>
                      <TableCell>
                        <SoftInput
                          type="number"
                          value={row.scaled_score}
                          onChange={(event) => updateRow(row.id, "scaled_score", event.target.value)}
                        />
                      </TableCell>
                      <TableCell align="right">
                        <SoftBox display="flex" justifyContent="flex-end" gap={0.5}>
                          <Tooltip title="Simpan baris ini">
                            <span>
                              <SoftButton size="small" variant="gradient" color="info" onClick={() => saveRow(row)} disabled={saving}>
                                Simpan
                              </SoftButton>
                            </span>
                          </Tooltip>
                          <Tooltip title="Hapus mapping">
                            <IconButton size="small" color="error" onClick={() => deleteRow(row)}>
                              <Icon fontSize="small">delete</Icon>
                            </IconButton>
                          </Tooltip>
                        </SoftBox>
                      </TableCell>
                    </TableRow>
                  ))}
                  {!rows.length && (
                    <TableRow>
                      <TableCell colSpan={4}>
                        <SoftBox py={6} textAlign="center">
                          <Icon color="disabled" sx={{ fontSize: 48 }}>grid_off</Icon>
                          <SoftTypography variant="h6" mt={1}>Belum ada mapping nilai</SoftTypography>
                          <SoftTypography variant="button" color="text">
                            Tambahkan konfigurasi untuk {sectionLabels[type]}.
                          </SoftTypography>
                        </SoftBox>
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </Card>
      </SoftBox>
      <Footer />
    </DashboardLayout>
  );
}

export default ScoreConversions;
