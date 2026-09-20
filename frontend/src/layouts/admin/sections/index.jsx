import { useCallback, useEffect, useMemo, useState } from "react";
import PropTypes from "prop-types";
import { DndContext, PointerSensor, closestCenter, useSensor, useSensors } from "@dnd-kit/core";
import { SortableContext, arrayMove, horizontalListSortingStrategy, useSortable, verticalListSortingStrategy } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
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
import SoftSelect from "components/SoftSelect";
import SoftTypography from "components/SoftTypography";
import MDEditor from "components/SoftEditor";
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

const emptySection = { title: "", section_type: "listening" };
const emptyQuestion = { question: "", passage: "", audio: null, image: null };
const typeMeta = {
  listening: { label: "Listening", color: "info" },
  structure: { label: "Structure", color: "warning" },
  reading: { label: "Reading", color: "success" },
};

function SortableItem({ id, children }) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id });

  return (
    <SoftBox
      ref={setNodeRef}
      sx={{
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.55 : 1,
        zIndex: isDragging ? 1 : "auto",
      }}
    >
      {children({ attributes, listeners })}
    </SoftBox>
  );
}

SortableItem.propTypes = {
  id: PropTypes.string.isRequired,
  children: PropTypes.func.isRequired,
};

function Sections() {
  const [sections, setSections] = useState([]);
  const [search, setSearch] = useState("");
  const [typeFilter, setTypeFilter] = useState("all");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [page, setPage] = useState(1);
  const rowsPerPage = 10;

  // Section dialog
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editId, setEditId] = useState(null);
  const [form, setForm] = useState(emptySection);
  const [saving, setSaving] = useState(false);

  // Question dialog
  const [detail, setDetail] = useState(null);
  const [targetSection, setTargetSection] = useState(null);
  const [questionOpen, setQuestionOpen] = useState(false);
  const [questionForm, setQuestionForm] = useState(emptyQuestion);
  const [editingQuestion, setEditingQuestion] = useState(null);
  const [optionText, setOptionText] = useState("");
  const [editingOptionId, setEditingOptionId] = useState(null);
  const [editOptionText, setEditOptionText] = useState("");
  const [previewOpen, setPreviewOpen] = useState(false);
  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 8 } }));

  // ── Data loading ────────────────────────────────────────
  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const response = await apiRequest("/sections/");
      setSections(response.data || []);
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const refreshDetail = async (sectionId = detail?.public_id) => {
    if (!sectionId) return;
    try {
      const response = await apiRequest(`/sections/${sectionId}`);
      setDetail(response.data);
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  // ── Filtering & pagination ──────────────────────────────
  const filtered = useMemo(() => {
    return sections.filter((section) => {
      const matchType = typeFilter === "all" || section.section_type === typeFilter;
      const matchSearch = `${section.title} ${section.section_type}`.toLowerCase().includes(search.toLowerCase());
      return matchType && matchSearch;
    });
  }, [sections, search, typeFilter]);

  useEffect(() => {
    setPage(1);
  }, [search, typeFilter]);

  const totalPages = Math.ceil(filtered.length / rowsPerPage);
  const paginated = filtered.slice((page - 1) * rowsPerPage, page * rowsPerPage);

  const counts = useMemo(() => ({
    all: sections.length,
    listening: sections.filter((s) => s.section_type === "listening").length,
    structure: sections.filter((s) => s.section_type === "structure").length,
    reading: sections.filter((s) => s.section_type === "reading").length,
  }), [sections]);

  // ── Section CRUD ────────────────────────────────────────
  const submitSection = async (event) => {
    event.preventDefault();
    setSaving(true);
    setError("");
    try {
      await apiRequest(editId ? `/sections/${editId}` : "/sections/", {
        method: editId ? "PUT" : "POST",
        body: JSON.stringify(form),
      });
      setDialogOpen(false);
      setEditId(null);
      setForm(emptySection);
      await load();
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setSaving(false);
    }
  };

  const removeSection = async (section) => {
    if (!window.confirm(`Hapus section "${section.title}"?`)) return;
    try {
      await apiRequest(`/sections/${section.public_id}`, { method: "DELETE" });
      await load();
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  // ── Question CRUD ───────────────────────────────────────
  const submitQuestion = async (event) => {
    event.preventDefault();
    setSaving(true);
    const body = new FormData();
    body.append("question", questionForm.question);
    if (questionForm.passage) body.append("passage", questionForm.passage);
    if (questionForm.audio) body.append("audio", questionForm.audio);
    if (questionForm.image) body.append("image", questionForm.image);

    try {
      if (editingQuestion) {
        await apiRequest(`/questions/${editingQuestion.public_id}`, { method: "PUT", body });
      } else {
        await apiRequest(`/sections/${targetSection.public_id}/questions`, { method: "POST", body });
      }
      setQuestionOpen(false);
      setQuestionForm(emptyQuestion);
      setEditingQuestion(null);
      await refreshDetail(targetSection?.public_id || detail?.public_id);
    } catch (requestError) {
      setError(requestError.message);
    } finally {
      setSaving(false);
    }
  };

  const removeQuestion = async (question) => {
    if (!window.confirm("Hapus soal ini?")) return;
    try {
      await apiRequest(`/questions/${question.public_id}`, { method: "DELETE" });
      await refreshDetail();
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  const reorderQuestion = async (questions, from, to) => {
    const next = arrayMove(questions, from, to);
    setDetail((prev) => ({ ...prev, questions: next }));
    try {
      await apiRequest(`/sections/${detail.public_id}/questions/reorder`, {
        method: "PUT",
        body: JSON.stringify({ question_public_ids: next.map((q) => q.public_id) }),
      });
    } catch (requestError) {
      setError(requestError.message);
      await refreshDetail();
    }
  };

  // ── Options CRUD ────────────────────────────────────────
  const addOption = async (question) => {
    if (!optionText.trim()) return;
    try {
      await apiRequest(`/questions/${question.public_id}/options`, {
        method: "POST",
        body: JSON.stringify({ option_text: optionText.trim() }),
      });
      setOptionText("");
      await refreshDetail();
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  const deleteOption = async (option) => {
    if (!window.confirm("Hapus pilihan jawaban ini?")) return;
    try {
      await apiRequest(`/options/${option.public_id}`, { method: "DELETE" });
      await refreshDetail();
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  const updateOption = async (option) => {
    if (!editOptionText.trim()) return;
    try {
      await apiRequest(`/options/${option.public_id}`, {
        method: "PUT",
        body: JSON.stringify({ option_text: editOptionText.trim() }),
      });
      setEditingOptionId(null);
      setEditOptionText("");
      await refreshDetail();
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  const setAnswer = async (question, option) => {
    try {
      await apiRequest(`/questions/${question.public_id}/answer-key`, {
        method: "PUT",
        body: JSON.stringify({ option_public_id: option.public_id }),
      });
      await refreshDetail();
    } catch (requestError) {
      setError(requestError.message);
    }
  };

  const reorderOptions = async (question, from, to) => {
    const next = arrayMove(question.options || [], from, to);
    // Optimistic update
    setDetail((prev) => ({
      ...prev,
      questions: (prev.questions || []).map((q) =>
        q.public_id === question.public_id ? { ...q, options: next } : q
      ),
    }));
    try {
      await apiRequest(`/questions/${question.public_id}/options/reorder`, {
        method: "PUT",
        body: JSON.stringify({ option_public_ids: next.map((o) => o.public_id) }),
      });
    } catch (requestError) {
      setError(requestError.message);
      await refreshDetail();
    }
  };

  const handleQuestionDragEnd = ({ active, over }) => {
    if (!over || active.id === over.id) return;
    const questions = detail?.questions || [];
    const from = questions.findIndex((question) => question.public_id === active.id);
    const to = questions.findIndex((question) => question.public_id === over.id);
    if (from !== -1 && to !== -1) reorderQuestion(questions, from, to);
  };

  const handleOptionDragEnd = (question) => ({ active, over }) => {
    if (!over || active.id === over.id) return;
    const options = question.options || [];
    const from = options.findIndex((option) => option.public_id === active.id);
    const to = options.findIndex((option) => option.public_id === over.id);
    if (from !== -1 && to !== -1) reorderOptions(question, from, to);
  };

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox pt={6} pb={3}>
        {/* Header */}
        <SoftBox mb={3} display="flex" justifyContent="space-between" alignItems={{ xs: "flex-start", sm: "center" }} flexDirection={{ xs: "column", sm: "row" }} gap={2}>
          <SoftBox>
            <SoftTypography variant="h4" fontWeight="bold">Bank Section dan Soal</SoftTypography>
            <SoftTypography variant="button" color="text" fontWeight="regular">
              Kelola section, soal, pilihan jawaban, dan answer key.
            </SoftTypography>
          </SoftBox>
          <SoftButton variant="gradient" color="info" onClick={() => { setEditId(null); setForm(emptySection); setDialogOpen(true); }}>
            <Icon sx={{ mr: 1 }}>add</Icon>Tambah section
          </SoftButton>
        </SoftBox>

        {/* Stats cards */}
        <Grid container spacing={3} mb={3}>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard title={{ text: "total section", fontWeight: "bold" }} count={counts.all} percentage={{ color: "info", text: "semua tipe" }} icon={{ color: "info", component: "view_module" }} />
          </Grid>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard title={{ text: "listening", fontWeight: "bold" }} count={counts.listening} percentage={{ color: "info", text: "audio" }} icon={{ color: "info", component: "headphones" }} />
          </Grid>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard title={{ text: "structure", fontWeight: "bold" }} count={counts.structure} percentage={{ color: "warning", text: "grammar" }} icon={{ color: "warning", component: "text_fields" }} />
          </Grid>
          <Grid item xs={12} sm={6} lg={3}>
            <MiniStatisticsCard title={{ text: "reading", fontWeight: "bold" }} count={counts.reading} percentage={{ color: "success", text: "comprehension" }} icon={{ color: "success", component: "menu_book" }} />
          </Grid>
        </Grid>

        {/* Table card */}
        <Card>
          <SoftBox p={3}>
            <Grid container spacing={2} alignItems="center">
              <Grid item xs={12} md={5}>
                <SoftInput icon={{ component: "search", direction: "left" }} placeholder="Cari section..." value={search} onChange={(event) => setSearch(event.target.value)} />
              </Grid>
              <Grid item xs={12} md={7} display="flex" justifyContent={{ xs: "flex-start", md: "flex-end" }} gap={1} flexWrap="wrap">
                {[["all", "Semua"], ["listening", "Listening"], ["structure", "Structure"], ["reading", "Reading"]].map(([value, label]) => (
                  <Chip key={value} label={label} color={typeFilter === value ? "info" : "default"} variant={typeFilter === value ? "filled" : "outlined"} onClick={() => setTypeFilter(value)} />
                ))}
                <Tooltip title="Muat ulang data">
                  <IconButton onClick={load} disabled={loading}>
                    <Icon>refresh</Icon>
                  </IconButton>
                </Tooltip>
              </Grid>
            </Grid>
          </SoftBox>

          {error && (
            <SoftBox mx={3} mb={2} p={2} borderRadius="md" bgColor="error">
              <SoftTypography variant="button" color="white">{error}</SoftTypography>
            </SoftBox>
          )}

          <Divider />

          {loading ? (
            <SoftBox p={3}>
              <Skeleton height={62} />
              <Skeleton height={62} />
              <Skeleton height={62} />
            </SoftBox>
          ) : (
            <TableContainer sx={{ overflowX: "auto", borderRadius: 2 }}>
              <Table sx={{ minWidth: 680, display: "table", ...globalTableStyle }}>
                <TableHead sx={{ display: "table-header-group" }}>
                  <TableRow sx={{ display: "table-row" }}>
                    <TableCell>Section</TableCell>
                    <TableCell>Tipe</TableCell>
                    <TableCell>Diperbarui</TableCell>
                    <TableCell align="right">Aksi</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {paginated.map((section) => (
                    <TableRow hover key={section.public_id}>
                      <TableCell>
                        <SoftTypography variant="button" fontWeight="bold">{section.title}</SoftTypography>
                      </TableCell>
                      <TableCell>
                        <SoftBadge color={typeMeta[section.section_type]?.color || "secondary"} variant="gradient" size="sm">
                          {typeMeta[section.section_type]?.label || section.section_type}
                        </SoftBadge>
                      </TableCell>
                      <TableCell>
                        <SoftTypography variant="caption" color="text">
                          {section.updated_at ? new Date(section.updated_at).toLocaleString("id-ID") : "-"}
                        </SoftTypography>
                      </TableCell>
                      <TableCell align="right">
                        <Tooltip title="Kelola soal">
                          <IconButton size="small" color="info" onClick={() => { setTargetSection(section); refreshDetail(section.public_id); }}>
                            <Icon>quiz</Icon>
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Edit section">
                          <IconButton size="small" color="warning" onClick={() => { setEditId(section.public_id); setForm({ title: section.title, section_type: section.section_type }); setDialogOpen(true); }}>
                            <Icon>edit</Icon>
                          </IconButton>
                        </Tooltip>
                        <Tooltip title="Hapus section">
                          <IconButton size="small" color="error" onClick={() => removeSection(section)}>
                            <Icon>delete</Icon>
                          </IconButton>
                        </Tooltip>
                      </TableCell>
                    </TableRow>
                  ))}
                  {!filtered.length && (
                    <TableRow>
                      <TableCell colSpan={4}>
                        <SoftBox py={6} textAlign="center">
                          <Icon color="disabled" sx={{ fontSize: 48 }}>inventory_2</Icon>
                          <SoftTypography variant="h6" mt={1}>Belum ada section</SoftTypography>
                          <SoftTypography variant="button" color="text">
                            Tambahkan section pertama dengan tombol di atas.
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
                Menampilkan {(page - 1) * rowsPerPage + 1} hingga {Math.min(page * rowsPerPage, filtered.length)} dari {filtered.length}
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

      {/* ── Section Form Dialog ────────────────────────────── */}
      <Dialog open={dialogOpen} onClose={() => !saving && setDialogOpen(false)} fullWidth maxWidth="sm">
        <SoftBox component="form" onSubmit={submitSection}>
          <DialogTitle>{editId ? "Edit Section" : "Tambah Section"}</DialogTitle>
          <DialogContent dividers>
            <SoftTypography component="label" variant="caption" fontWeight="bold">Judul section</SoftTypography>
            <SoftInput sx={{ mt: 1 }} value={form.title} onChange={(event) => setForm({ ...form, title: event.target.value })} required />
            <SoftBox mb={1} mt={2} lineHeight={0} display="inline-block">
              <SoftTypography component="label" variant="caption" fontWeight="bold">
                Tipe section
              </SoftTypography>
            </SoftBox>
            <SoftSelect
              placeholder="Pilih tipe section..."
              options={Object.entries(typeMeta).map(([value, meta]) => ({ value, label: meta.label }))}
              value={{ value: form.section_type, label: typeMeta[form.section_type]?.label || "" }}
              onChange={(option) => setForm({ ...form, section_type: option ? option.value : "" })}
            />
            <SoftTypography variant="caption" color="text" display="block" mt={1}>
              Pilih kategori materi yang akan dikelola.
            </SoftTypography>
          </DialogContent>
          <DialogActions>
            <SoftButton color="secondary" onClick={() => setDialogOpen(false)}>Batal</SoftButton>
            <SoftButton type="submit" variant="gradient" color="info" disabled={saving}>
              {saving ? "Menyimpan..." : "Simpan"}
            </SoftButton>
          </DialogActions>
        </SoftBox>
      </Dialog>

      {/* ── Question Detail Dialog ─────────────────────────── */}
      <Dialog open={Boolean(detail)} onClose={() => setDetail(null)} fullWidth maxWidth="lg">
        <DialogTitle>
          <SoftBox display="flex" justifyContent="space-between" alignItems="center">
            <SoftBox>
              <SoftTypography variant="h5" fontWeight="bold">{detail?.title}</SoftTypography>
              <SoftTypography variant="caption" color="text">{detail?.questions?.length || 0} soal</SoftTypography>
            </SoftBox>
            <SoftBox display="flex" gap={1}>
              <SoftButton variant="outlined" color="dark" onClick={() => setPreviewOpen(true)}>
                <Icon sx={{ mr: 1 }}>visibility</Icon>Preview
              </SoftButton>
              <SoftButton variant="gradient" color="info" onClick={() => { setTargetSection(detail); setEditingQuestion(null); setQuestionForm(emptyQuestion); setQuestionOpen(true); }}>
                <Icon sx={{ mr: 1 }}>add</Icon>Tambah soal
              </SoftButton>
            </SoftBox>
          </SoftBox>
        </DialogTitle>
        <DialogContent dividers>
          <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleQuestionDragEnd}>
            <SortableContext items={(detail?.questions || []).map((question) => question.public_id)} strategy={verticalListSortingStrategy}>
              {(detail?.questions || []).map((question, index) => (
                <SortableItem key={question.public_id} id={question.public_id}>
                  {({ attributes, listeners }) => (
                    <SoftBox py={2} borderBottom="1px solid" borderColor="grey.200">
              <SoftBox display="flex" justifyContent="space-between" gap={2} alignItems="flex-start">
                <SoftBox display="flex" alignItems="flex-start" gap={1}>
                  <IconButton
                    size="small"
                    {...attributes}
                    {...listeners}
                    aria-label="Geser soal"
                    sx={{ cursor: "grab", touchAction: "none" }}
                  >
                    <Icon fontSize="small">open_with</Icon>
                  </IconButton>
                  <SoftTypography variant="button" fontWeight="bold">
                    {index + 1}. {question.question}
                  </SoftTypography>
                </SoftBox>
                <SoftBox display="flex" gap={0.5} flexShrink={0}>
                  <Tooltip title="Edit soal">
                    <IconButton size="small" color="warning" onClick={() => { setTargetSection(detail); setEditingQuestion(question); setQuestionForm({ question: question.question, passage: question.passage || "", audio: null, image: null }); setQuestionOpen(true); }}>
                      <Icon fontSize="small">edit</Icon>
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Hapus soal">
                    <IconButton size="small" color="error" onClick={() => removeQuestion(question)}>
                      <Icon fontSize="small">delete</Icon>
                    </IconButton>
                  </Tooltip>
                </SoftBox>
              </SoftBox>

              {question.passage && (
                <SoftBox mt={1} p={2} bgColor="grey.100" borderRadius="md">
                  <div dangerouslySetInnerHTML={{ __html: question.passage }} />
                </SoftBox>
              )}

              {/* Options */}
              <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleOptionDragEnd(question)}>
                <SortableContext items={(question.options || []).map((option) => option.public_id)} strategy={horizontalListSortingStrategy}>
                  <SoftBox mt={2} display="flex" gap={1} flexWrap="wrap">
                    {(question.options || []).map((option) => (
                      <SortableItem key={option.public_id} id={option.public_id}>
                        {({ attributes: optionAttributes, listeners: optionListeners }) => (
                          <SoftBox
                            display="flex"
                            alignItems="center"
                            gap={0.5}
                            px={1.5}
                            py={0.75}
                            borderRadius="md"
                            bgColor={question.correct_option_public_id === option.public_id ? "success" : "grey.100"}
                            sx={{ transition: "all 0.2s" }}
                          >
                            <IconButton
                              size="small"
                              {...optionAttributes}
                              {...optionListeners}
                              aria-label="Geser opsi jawaban"
                              sx={{ cursor: "grab", touchAction: "none" }}
                            >
                              <Icon sx={{ fontSize: 16 }}>open_with</Icon>
                            </IconButton>
                    {editingOptionId === option.public_id ? (
                      <SoftBox display="flex" alignItems="center" gap={1}>
                        <SoftTypography variant="button" color="text" fontWeight="bold">
                          {option.label}:
                        </SoftTypography>
                        <SoftInput 
                          size="small"
                          value={editOptionText}
                          onChange={(e) => setEditOptionText(e.target.value)}
                          autoFocus
                          onKeyDown={(e) => {
                            if (e.key === "Enter") updateOption(option);
                            if (e.key === "Escape") setEditingOptionId(null);
                          }}
                        />
                        <IconButton size="small" color="success" onClick={() => updateOption(option)}>
                          <Icon sx={{ fontSize: 16 }}>check</Icon>
                        </IconButton>
                        <IconButton size="small" color="secondary" onClick={() => setEditingOptionId(null)}>
                          <Icon sx={{ fontSize: 16 }}>close</Icon>
                        </IconButton>
                      </SoftBox>
                    ) : (
                      <>
                        <SoftButton
                          size="small"
                          color={question.correct_option_public_id === option.public_id ? "success" : "light"}
                          onClick={() => setAnswer(question, option)}
                        >
                          {option.label}: {option.option_text}
                        </SoftButton>
                        <SoftBox display="flex" gap={0}>
                          <IconButton size="small" onClick={() => { setEditingOptionId(option.public_id); setEditOptionText(option.option_text); }}>
                            <Icon sx={{ fontSize: 14 }}>edit</Icon>
                          </IconButton>
                          <IconButton size="small" onClick={() => deleteOption(option)}>
                            <Icon sx={{ fontSize: 14 }}>delete</Icon>
                          </IconButton>
                        </SoftBox>
                      </>
                    )}
                          </SoftBox>
                        )}
                      </SortableItem>
                    ))}
                  </SoftBox>
                </SortableContext>
              </DndContext>

              {/* Add option inline */}
              <SoftBox mt={1} display="flex" gap={1} alignItems="center">
                <SoftInput
                  placeholder="Tulis pilihan jawaban..."
                  value={detail?.activeQuestionId === question.public_id ? optionText : ""}
                  onFocus={() => setDetail((prev) => ({ ...prev, activeQuestionId: question.public_id }))}
                  onChange={(event) => { setDetail((prev) => ({ ...prev, activeQuestionId: question.public_id })); setOptionText(event.target.value); }}
                />
                <SoftButton size="small" color="info" variant="gradient" onClick={() => addOption(question)}>
                  Tambah opsi
                </SoftButton>
              </SoftBox>
                    </SoftBox>
                  )}
                </SortableItem>
              ))}
            </SortableContext>
          </DndContext>

          {!detail?.questions?.length && (
            <SoftBox py={6} textAlign="center">
              <Icon color="disabled" sx={{ fontSize: 48 }}>quiz</Icon>
              <SoftTypography variant="h6" mt={1}>Belum ada soal</SoftTypography>
              <SoftTypography variant="button" color="text">
                Tambahkan soal pertama untuk section ini.
              </SoftTypography>
            </SoftBox>
          )}
        </DialogContent>
        <DialogActions>
          <SoftButton color="light" onClick={() => setDetail(null)}>Tutup</SoftButton>
        </DialogActions>
      </Dialog>

      {/* ── Preview Dialog ─────────────────────────────────── */}
      <Dialog open={previewOpen} onClose={() => setPreviewOpen(false)} fullWidth maxWidth="md">
        <DialogTitle>Preview Section: {detail?.title}</DialogTitle>
        <DialogContent dividers sx={{ backgroundColor: "#f8f9fa" }}>
          {(detail?.questions || []).map((question, index) => (
            <Card key={question.public_id} sx={{ mb: 3 }}>
              <SoftBox p={3}>
                <SoftTypography variant="button" fontWeight="bold">
                  {question.number || index + 1}. {question.question}
                </SoftTypography>

                {question.passage && (
                  <SoftBox mt={2} p={2} bgColor="grey.100" borderRadius="md">
                    <div dangerouslySetInnerHTML={{ __html: question.passage }} />
                  </SoftBox>
                )}

                {question.audio_url && (
                  <SoftBox mt={2}>
                    <audio controls src={question.audio_url} style={{ width: "100%" }} />
                  </SoftBox>
                )}

                {question.image_url && (
                  <SoftBox mt={2}>
                    <img src={question.image_url} alt="Media soal" style={{ maxWidth: "100%", maxHeight: 280, borderRadius: 8 }} />
                  </SoftBox>
                )}

                <SoftBox mt={2} display="flex" flexDirection="column" gap={1}>
                  {(question.options || []).map((option) => (
                    <SoftBox 
                      key={option.public_id} 
                      display="flex" 
                      alignItems="center" 
                      gap={1}
                      p={1.5}
                      borderRadius="md"
                      sx={{ border: "1px solid", borderColor: "grey.300", opacity: question.correct_option_public_id === option.public_id ? 1 : 0.7, backgroundColor: question.correct_option_public_id === option.public_id ? "success.main" : "transparent" }}
                    >
                      <input type="radio" disabled checked={question.correct_option_public_id === option.public_id} />
                      <SoftTypography variant="button" color={question.correct_option_public_id === option.public_id ? "white" : "text"}>
                        {option.label}. {option.option_text}
                      </SoftTypography>
                    </SoftBox>
                  ))}
                </SoftBox>
              </SoftBox>
            </Card>
          ))}
        </DialogContent>
        <DialogActions>
          <SoftButton color="light" onClick={() => setPreviewOpen(false)}>Tutup Preview</SoftButton>
        </DialogActions>
      </Dialog>

      {/* ── Question Form Dialog ───────────────────────────── */}
      <Dialog open={questionOpen} onClose={() => !saving && setQuestionOpen(false)} fullWidth maxWidth="md">
        <SoftBox component="form" onSubmit={submitQuestion}>
          <DialogTitle>
            {editingQuestion ? "Edit Soal" : "Tambah Soal"} — {targetSection?.title}
          </DialogTitle>
          <DialogContent dividers>
            <SoftTypography component="label" variant="caption" fontWeight="bold">Teks soal</SoftTypography>
            <SoftInput
              sx={{ mt: 1 }}
              multiline
              rows={4}
              value={questionForm.question}
              onChange={(event) => setQuestionForm({ ...questionForm, question: event.target.value })}
              required
            />
            <SoftTypography component="label" variant="caption" fontWeight="bold" display="block" mt={2} mb={1}>
              Passage (opsional)
            </SoftTypography>
            <SoftBox sx={{ border: "1px solid", borderColor: "grey.300", borderRadius: 1 }}>
              {questionOpen && (
                <MDEditor
                  key={editingQuestion ? editingQuestion.public_id : 'new'}
                  initialValue={questionForm.passage}
                  value={(html) => setQuestionForm({ ...questionForm, passage: html })}
                />
              )}
            </SoftBox>
            <Grid container spacing={2} mt={1}>
              <Grid item xs={12} md={6}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Audio</SoftTypography>
                <SoftInput
                  type="file"
                  inputProps={{ accept: "audio/*" }}
                  sx={{ mt: 1 }}
                  onChange={(event) => setQuestionForm({ ...questionForm, audio: event.target.files[0] || null })}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <SoftTypography component="label" variant="caption" fontWeight="bold">Gambar</SoftTypography>
                <SoftInput
                  type="file"
                  inputProps={{ accept: "image/*" }}
                  sx={{ mt: 1 }}
                  onChange={(event) => setQuestionForm({ ...questionForm, image: event.target.files[0] || null })}
                />
              </Grid>
            </Grid>
          </DialogContent>
          <DialogActions>
            <SoftButton color="secondary" onClick={() => setQuestionOpen(false)}>Batal</SoftButton>
            <SoftButton type="submit" variant="gradient" color="info" disabled={saving}>
              {saving ? "Menyimpan..." : "Simpan soal"}
            </SoftButton>
          </DialogActions>
        </SoftBox>
      </Dialog>
    </DashboardLayout>
  );
}

export default Sections;
