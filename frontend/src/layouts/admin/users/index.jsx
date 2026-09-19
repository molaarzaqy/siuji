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
import Skeleton from "@mui/material/Skeleton";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
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

const roleMeta = {
  admin: { label: "Admin", color: "dark" },
  participant: { label: "Participant", color: "info" },
};

function Users() {
  const [users, setUsers] = useState([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [page, setPage] = useState(1);
  const [meta, setMeta] = useState({});
  const rowsPerPage = 10;

  // Detail dialog
  const [detailOpen, setDetailOpen] = useState(false);
  const [detailData, setDetailData] = useState(null);
  const [detailLoading, setDetailLoading] = useState(false);

  // ── Data loading ────────────────────────────────────────
  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const response = await apiRequest(`/users/?page=${page}&limit=${rowsPerPage}`);
      setUsers(response.data || []);
      setMeta(response.meta || {});
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    load();
  }, [load]);

  // ── Search filter (client-side on loaded page) ──────────
  const filtered = useMemo(() => {
    if (!search.trim()) return users;
    const query = search.toLowerCase();
    return users.filter((user) =>
      `${user.name || ""} ${user.email || ""} ${user.nim || ""} ${user.university || ""}`.toLowerCase().includes(query)
    );
  }, [users, search]);

  // ── Stats ───────────────────────────────────────────────
  const counts = useMemo(() => ({
    total: meta.total_datas ?? users.length,
    admin: users.filter((u) => u.role === "admin").length,
    participant: users.filter((u) => u.role === "participant").length,
  }), [users, meta]);

  const totalPages = meta.total_pages || Math.ceil((meta.total_datas || users.length) / rowsPerPage);

  // ── View detail ─────────────────────────────────────────
  const viewDetail = async (user) => {
    setDetailLoading(true);
    setDetailOpen(true);
    try {
      const response = await apiRequest(`/users/${user.public_id}`);
      setDetailData(response.data);
    } catch (requestError) {
      setError(requestError.message);
      setDetailOpen(false);
    } finally {
      setDetailLoading(false);
    }
  };

  // ── Delete user ─────────────────────────────────────────
  const remove = async (user) => {
    const identifier = user.name || user.email;
    if (!window.confirm(`Hapus user "${identifier}"?\n\nPerhatian: User yang masih memiliki data ujian mungkin tidak dapat dihapus.`)) return;
    try {
      await apiRequest(`/users/${user.public_id}`, { method: "DELETE" });
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
        <SoftBox mb={3}>
          <SoftTypography variant="h4" fontWeight="bold">Manajemen User</SoftTypography>
          <SoftTypography variant="button" color="text" fontWeight="regular">
            Kelola akun admin dan participant yang terdaftar.
          </SoftTypography>
        </SoftBox>

        {/* Stats cards */}
        <Grid container spacing={3} mb={3}>
          <Grid item xs={12} sm={4}>
            <MiniStatisticsCard
              title={{ text: "total user", fontWeight: "bold" }}
              count={counts.total}
              percentage={{ color: "info", text: "akun" }}
              icon={{ color: "info", component: "people" }}
            />
          </Grid>
          <Grid item xs={12} sm={4}>
            <MiniStatisticsCard
              title={{ text: "admin", fontWeight: "bold" }}
              count={counts.admin}
              percentage={{ color: "dark", text: "pengelola" }}
              icon={{ color: "dark", component: "admin_panel_settings" }}
            />
          </Grid>
          <Grid item xs={12} sm={4}>
            <MiniStatisticsCard
              title={{ text: "participant", fontWeight: "bold" }}
              count={counts.participant}
              percentage={{ color: "success", text: "peserta ujian" }}
              icon={{ color: "success", component: "school" }}
            />
          </Grid>
        </Grid>

        {/* Main card */}
        <Card>
          <SoftBox p={3} display="flex" gap={2} alignItems="center" flexWrap="wrap">
            <SoftBox flex={1} minWidth={240}>
              <SoftInput
                icon={{ component: "search", direction: "left" }}
                placeholder="Cari nama, email, NIM, atau universitas..."
                value={search}
                onChange={(event) => setSearch(event.target.value)}
              />
            </SoftBox>
            <Tooltip title="Muat ulang data">
              <IconButton onClick={load} disabled={loading}>
                <Icon>refresh</Icon>
              </IconButton>
            </Tooltip>
          </SoftBox>

          {error && (
            <SoftBox mx={3} mb={2} p={2} borderRadius="md" bgColor="error">
              <SoftTypography variant="button" color="white">{error}</SoftTypography>
            </SoftBox>
          )}

          <Divider />

          {loading ? (
            <SoftBox p={3}>
              <Skeleton height={64} />
              <Skeleton height={64} />
              <Skeleton height={64} />
            </SoftBox>
          ) : (
            <TableContainer sx={{ overflowX: "auto", borderRadius: 2 }}>
              <Table sx={{ minWidth: 720, display: "table", ...globalTableStyle }}>
                <TableHead sx={{ display: "table-header-group" }}>
                  <TableRow sx={{ display: "table-row" }}>
                    <TableCell>Nama</TableCell>
                    <TableCell>Email</TableCell>
                    <TableCell>Role</TableCell>
                    <TableCell>NIM / Universitas</TableCell>
                    <TableCell align="right">Aksi</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {filtered.map((user) => {
                    const role = roleMeta[user.role] || { label: user.role || "-", color: "secondary" };
                    return (
                      <TableRow hover key={user.public_id}>
                        <TableCell>
                          <SoftTypography variant="button" fontWeight="bold">
                            {user.name || "-"}
                          </SoftTypography>
                        </TableCell>
                        <TableCell>
                          <SoftTypography variant="caption" color="text">{user.email}</SoftTypography>
                        </TableCell>
                        <TableCell>
                          <SoftBadge color={role.color} variant="gradient" size="sm">{role.label}</SoftBadge>
                        </TableCell>
                        <TableCell>
                          <SoftTypography variant="caption" display="block">{user.nim || "-"}</SoftTypography>
                          <SoftTypography variant="caption" color="text">{user.university || "-"}</SoftTypography>
                        </TableCell>
                        <TableCell align="right">
                          <SoftBox display="flex" justifyContent="flex-end" gap={0.5}>
                            <Tooltip title="Lihat detail">
                              <IconButton size="small" color="info" onClick={() => viewDetail(user)}>
                                <Icon fontSize="small">visibility</Icon>
                              </IconButton>
                            </Tooltip>
                            <Tooltip title="Hapus user">
                              <IconButton size="small" color="error" onClick={() => remove(user)}>
                                <Icon fontSize="small">delete</Icon>
                              </IconButton>
                            </Tooltip>
                          </SoftBox>
                        </TableCell>
                      </TableRow>
                    );
                  })}
                  {!filtered.length && (
                    <TableRow>
                      <TableCell colSpan={5}>
                        <SoftBox py={6} textAlign="center">
                          <Icon color="disabled" sx={{ fontSize: 48 }}>people_outline</Icon>
                          <SoftTypography variant="h6" mt={1}>Belum ada user</SoftTypography>
                          <SoftTypography variant="button" color="text">
                            User participant dibuat melalui menu Peserta.
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
                Halaman {page} dari {totalPages} ({counts.total} user)
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

      {/* ── Detail User Dialog ─────────────────────────────── */}
      <Dialog open={detailOpen} onClose={() => setDetailOpen(false)} fullWidth maxWidth="sm">
        <DialogTitle>Detail User</DialogTitle>
        <DialogContent dividers>
          {detailLoading ? (
            <SoftBox p={3}>
              <Skeleton height={40} />
              <Skeleton height={40} />
              <Skeleton height={40} />
            </SoftBox>
          ) : detailData && (
            <Grid container spacing={2}>
              <Grid item xs={12}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Nama</SoftTypography>
                <SoftTypography variant="button" display="block">{detailData.name || "-"}</SoftTypography>
              </Grid>
              <Grid item xs={12}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Email</SoftTypography>
                <SoftTypography variant="button" display="block">{detailData.email || "-"}</SoftTypography>
              </Grid>
              <Grid item xs={6}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Role</SoftTypography>
                <SoftBox mt={0.5}>
                  <SoftBadge
                    color={roleMeta[detailData.role]?.color || "secondary"}
                    variant="gradient"
                    size="sm"
                  >
                    {roleMeta[detailData.role]?.label || detailData.role || "-"}
                  </SoftBadge>
                </SoftBox>
              </Grid>
              <Grid item xs={6}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">NIM</SoftTypography>
                <SoftTypography variant="button" display="block">{detailData.nim || "-"}</SoftTypography>
              </Grid>
              <Grid item xs={12}>
                <SoftTypography variant="caption" color="text" fontWeight="bold">Universitas</SoftTypography>
                <SoftTypography variant="button" display="block">{detailData.university || "-"}</SoftTypography>
              </Grid>
              {detailData.created_at && (
                <Grid item xs={12}>
                  <SoftTypography variant="caption" color="text" fontWeight="bold">Terdaftar</SoftTypography>
                  <SoftTypography variant="button" display="block">
                    {new Date(detailData.created_at).toLocaleString("id-ID")}
                  </SoftTypography>
                </Grid>
              )}
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

export default Users;
