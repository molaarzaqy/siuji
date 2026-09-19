import { useEffect, useState } from "react";

import Grid from "@mui/material/Grid";
import Icon from "@mui/material/Icon";

import SoftBox from "components/SoftBox";
import SoftTypography from "components/SoftTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import Footer from "examples/Footer";
import MiniStatisticsCard from "examples/Cards/StatisticsCards/MiniStatisticsCard";
import SalesTable from "examples/Tables/SalesTable";
import ReportsBarChart from "examples/Charts/BarCharts/ReportsBarChart";
import GradientLineChart from "examples/Charts/LineCharts/GradientLineChart";
import Globe from "examples/Globe";
import EventCalendar from "examples/Calendar";
import typography from "assets/theme/base/typography";
import breakpoints from "assets/theme/base/breakpoints";
import { apiRequest } from "services/api";

const emptyStats = {
  periods: [],
  sections: [],
  users: [],
  periodMeta: {},
  userMeta: {},
};

function Default() {
  const { values } = breakpoints;
  const { size } = typography;
  const [stats, setStats] = useState(emptyStats);

  useEffect(() => {
    let mounted = true;

    Promise.all([
      apiRequest("/periods/?page=1&limit=100"),
      apiRequest("/sections/"),
      apiRequest("/users/?page=1&limit=1000"),
    ])
      .then(([periodResponse, sectionResponse, userResponse]) => {
        if (!mounted) return;
        setStats({
          periods: periodResponse.data || [],
          sections: sectionResponse.data || [],
          users: userResponse.data || [],
          periodMeta: periodResponse.meta || {},
          userMeta: userResponse.meta || {},
        });
      })
      .catch(() => {
        // Dashboard tetap tampil dengan nilai kosong jika data belum tersedia.
      });

    return () => { mounted = false; };
  }, []);

  const { periods, sections, users, periodMeta, userMeta } = stats;
  const totalPeriods = periodMeta.total_datas ?? periods.length;
  const totalUsers = userMeta.total_datas ?? users.length;
  const totalParticipants = users.filter((user) => user.role === "participant").length;
  const publishedPeriods = periods.filter((period) => period.status === "published");
  const activePeriod = publishedPeriods.find((period) => {
    const now = Date.now();
    return new Date(period.start_time).getTime() <= now && new Date(period.end_time).getTime() >= now;
  }) || publishedPeriods[0];
  const statusCounts = ["draft", "published", "closed"].map((status) => periods.filter((period) => period.status === status).length);
  const periodRows = periods.slice(0, 4).map((period) => ({
    period: period.title,
    status: period.status,
    schedule: `${period.month} ${period.year}`,
    start: new Date(period.start_time).toLocaleDateString("id-ID"),
  }));
  const calendarEvents = periods.map((period) => ({
    title: period.title,
    start: period.start_time,
    end: period.end_time,
    className: period.status === "published" ? "success" : period.status === "closed" ? "dark" : "warning",
  }));

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <SoftBox py={3}>
        <Grid container>
          <Grid item xs={12} lg={7}>
            <SoftBox mb={3} p={1}>
              <SoftTypography variant={window.innerWidth < values.sm ? "h3" : "h2"} textTransform="capitalize" fontWeight="bold">
                dashboard admin SIUJI
              </SoftTypography>
              <SoftTypography variant="button" color="text" fontWeight="regular">
                Ringkasan pengelolaan ujian, section, dan peserta.
              </SoftTypography>
            </SoftBox>

            <Grid container>
              <Grid item xs={12}>
                <Globe display={{ xs: "none", md: "block" }} position="absolute" top="10%" right={0} mt={{ xs: -12, lg: 1 }} mr={{ xs: 0, lg: 10 }} canvasStyle={{ marginTop: "3rem" }} />
              </Grid>
            </Grid>

            <Grid container spacing={3}>
              <Grid item xs={12} sm={5}>
                <SoftBox mb={3}>
                  <MiniStatisticsCard title={{ text: "total periode", fontWeight: "bold" }} count={totalPeriods} percentage={{ color: "success", text: "ujian" }} icon={{ color: "info", component: "event" }} />
                </SoftBox>
                <MiniStatisticsCard title={{ text: "total section", fontWeight: "bold" }} count={sections.length} percentage={{ color: "success", text: "materi" }} icon={{ color: "info", component: "view_module" }} />
              </Grid>
              <Grid item xs={12} sm={5}>
                <SoftBox mb={3}>
                  <MiniStatisticsCard title={{ text: "total pengguna", fontWeight: "bold" }} count={totalUsers} percentage={{ color: "success", text: "akun" }} icon={{ color: "info", component: "people" }} />
                </SoftBox>
                <SoftBox mb={3}>
                  <MiniStatisticsCard title={{ text: "peserta ujian", fontWeight: "bold" }} count={totalParticipants} percentage={{ color: "success", text: "participant" }} icon={{ color: "info", component: "school" }} />
                </SoftBox>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs={12} md={10} lg={7}>
            <Grid item xs={12} lg={10}>
              <SoftBox mb={3} position="relative">
                <SalesTable title="Periode ujian terbaru" rows={periodRows.length ? periodRows : [{ period: "Belum ada periode", status: "-", schedule: "-", start: "-" }]} />
              </SoftBox>
            </Grid>
          </Grid>

          <Grid container spacing={3}>
            <Grid item xs={12} lg={5}>
              <ReportsBarChart
                title="status periode ujian"
                description={activePeriod ? <>Periode aktif: <strong>{activePeriod.title}</strong></> : "Belum ada periode aktif"}
                chart={{ labels: ["Draft", "Published", "Closed"], datasets: { label: "Periode", data: statusCounts } }}
                items={[
                  { icon: { color: "primary", component: "event" }, label: "semua periode", progress: { content: totalPeriods, percentage: totalPeriods ? 100 : 0 } },
                  { icon: { color: "info", component: "play_circle" }, label: "published", progress: { content: statusCounts[1], percentage: totalPeriods ? Math.round((statusCounts[1] / totalPeriods) * 100) : 0 } },
                  { icon: { color: "warning", component: "view_module" }, label: "section", progress: { content: sections.length, percentage: 100 } },
                  { icon: { color: "error", component: "people" }, label: "peserta", progress: { content: totalParticipants, percentage: totalUsers ? Math.round((totalParticipants / totalUsers) * 100) : 0 } },
                ]}
              />
            </Grid>
            <Grid item xs={12} lg={7}>
              <GradientLineChart
                title="ringkasan periode"
                description={
                  <SoftBox display="flex" alignItems="center">
                    <SoftBox fontSize={size.lg} color="success" mb={0.3} mr={0.5} lineHeight={0}><Icon sx={{ fontWeight: "bold" }}>insights</Icon></SoftBox>
                    <SoftTypography variant="button" color="text" fontWeight="medium">Status periode ujian SIUJI</SoftTypography>
                  </SoftBox>
                }
                chart={{ labels: ["Draft", "Published", "Closed"], datasets: [{ label: "Periode", color: "info", data: statusCounts }] }}
              />
            </Grid>
          </Grid>
          <Grid item xs={12} mt={3}>
            <EventCalendar
              header={{ title: "Kalender periode ujian", date: activePeriod ? `Periode aktif: ${activePeriod.title}` : "Jadwal ujian SIUJI" }}
              initialView="dayGridMonth"
              initialDate={activePeriod?.start_time || undefined}
              events={calendarEvents}
            />
          </Grid>
        </Grid>
      </SoftBox>
      <Footer />
    </DashboardLayout>
  );
}

export default Default;
